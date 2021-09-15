# gostutter
Simple linter for golang for stuttering detection (repeating of names) in code.



| Package name  | Code          | Relaxed mode  | Strict mode |
| --------- | ------------- | ------------- | ------------- |
| foo | `func NewFoo() *Foo`  | Allowed | Forbidden  |
| foo | `func HandleFoo()`  | Allowed | Forbidden  |
| foo | `type ConfigFoo struct`  | Allowed | Forbidden  |
| foo | `func fooHandle()`  | Forbidden | Forbidden  |
| foo | `type FooConfig struct`  | Forbidden | Forbidden  |
| - | ```type Config struct { config int }```  | Forbidden | Forbidden  |
 
 

### How to run
```
go run ./cmd/lint/main.go --  ./...
```

Output
```
pkg/stutter/analyzer.go:54:6: function name "runStutterCheck"  contains name of package "stutter"
pkg/stutter/analyzer.go:78:19: function name "checkStutter"  contains name of package "stutter"
pkg/stutter/analyzer.go:172:6: function name "stutteringDemo"  contains name of package "stutter"
pkg/stutter/analyzer.go:176:6: type name "Stutter" contains name of package "stutter"
pkg/stutter/analyzer.go:177:2: field name "stutter" contains name of structure "Stutter"
```

### Strict mode
GoStutter has strict feature that checks for any substring in functions, struct fields or global variable names. To start with strict mode, just add strict parameter:
```
go run ./cmd/lint/main.go --strict=true  --  ./...
```
