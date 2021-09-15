package stutter

import (
	"bytes"
	"flag"
	"go/ast"
	"go/printer"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
)

type CheckType int

const (
	UnknownIssue         CheckType = 0
	PackageNameIssue     CheckType = 1
	FieldNameIssue       CheckType = 2
	FuncNamePackageIssue CheckType = 3
	FuncNameStructIssue  CheckType = 4
)

type Context struct {
	strictMode *bool
}

func (c *Context) isStrict() bool {
	return *c.strictMode
}

func NewAnalyzer() *analysis.Analyzer {
	flagSet := *flag.NewFlagSet("stutter", flag.ExitOnError)

	config := &Context{
		strictMode: flagSet.Bool("strict", false, "Will work in strict mode"),
	}

	return &analysis.Analyzer{
		Name: "stutter",
		Doc:  "checks for stuttering",
		Run: func(p *analysis.Pass) (interface{}, error) {
			return runCheck(config, p)
		},
		Flags: flagSet,
	}
}

type Visitor struct {
	config      *Context
	packageName string
	detected    []Issue
}

func (v *Visitor) check(parentName string, childName string, kind CheckType, node ast.Node) {
	sanChildName := strings.ToLower(childName)
	sanParentName := strings.ToLower(parentName)

	switch kind {
	case FuncNamePackageIssue:
		// Function main can have same name as package
		if childName == "main" {
			return
		}
		// NewSomething in package something is allowed for nonstrict mode
		if !v.config.isStrict() && strings.HasPrefix(sanChildName, "new") {
			return
		}
	}

	var detected bool
	if v.config.isStrict() {
		detected = strings.Contains(sanChildName, sanParentName)
	} else {
		detected = strings.HasPrefix(sanChildName, sanParentName)
	}

	if detected {
		v.detected = append(v.detected, Issue{
			pos:       node.Pos(),
			node:      node,
			kind:      kind,
			substring: parentName,
		})
	}
}

func (v *Visitor) Visit(node ast.Node) ast.Visitor {
	switch node := node.(type) {
	case *ast.TypeSpec:
		v.visitTypeSpec(node)
		return nil
	case *ast.FuncDecl:
		v.visitFuncDecl(node)
		return nil
	case *ast.GenDecl:
		if node.Tok == token.IMPORT {
			return nil
		}
	case *ast.ValueSpec:
		for _, ident := range node.Names {
			v.check(v.packageName, ident.Name, PackageNameIssue, ident)
		}
		return nil
	}

	return v
}

func (v *Visitor) visitTypeSpec(node *ast.TypeSpec) {
	typeName := node.Name.Name
	// Verify for: typename vs package
	v.check(v.packageName, typeName, PackageNameIssue, node.Name)

	if strukt, ok := node.Type.(*ast.StructType); ok {
		for _, field := range strukt.Fields.List {
			for _, fieldName := range field.Names {
				// Verify for: fields vs struct name
				v.check(typeName, fieldName.Name, FieldNameIssue, fieldName)
			}
		}
	}
}

func (v *Visitor) visitFuncDecl(node *ast.FuncDecl) {
	// Verify for: func name vs package
	v.check(v.packageName, node.Name.Name, FuncNamePackageIssue, node.Name)

	structName := extractStructNameFromFunc(node)
	// Verify for: func name vs struct name
	if structName != "" {
		v.check(structName, node.Name.Name, FuncNameStructIssue, node.Name)
	}
}

func extractStructNameFromFunc(funcDecl *ast.FuncDecl) string {
	if funcDecl.Recv == nil {
		return ""
	}

	for _, f := range funcDecl.Recv.List {
		switch node := f.Type.(type) {
		case *ast.StarExpr:
			if ident, ok := node.X.(*ast.Ident); ok {
				return ident.Name
			}
		case *ast.Ident:
			return node.Name
		}
	}

	return ""
}

func runCheck(c *Context, pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		visitor := &Visitor{
			config:      c,
			packageName: strings.ToLower(file.Name.Name),
		}

		ast.Walk(visitor, file)

		for _, issue := range visitor.detected {
			var buf bytes.Buffer
			printer.Fprint(&buf, pass.Fset, issue.node)
			pass.Reportf(issue.pos, IssueKindToString(issue.kind), buf.String(), issue.substring)
		}
	}
	return nil, nil
}
