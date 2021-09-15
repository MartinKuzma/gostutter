package stutter

import "go/token"

type Issue struct {
	pos       token.Pos
	node      interface{}
	kind      CheckType
	substring string
}

func IssueKindToString(kind CheckType) string {
	switch kind {
	case PackageNameIssue:
		return "declaration of \"%s\" contains name of package \"%s\""
	case FieldNameIssue:
		return "field name \"%s\" contains name of structure \"%s\""
	case FuncNamePackageIssue:
		return "function name \"%s\"  contains name of package \"%s\""
	case FuncNameStructIssue:
		return "function name \"%q\" contains name of structure \"%s\""
	default:
		return "stuttering detected"
	}
}
