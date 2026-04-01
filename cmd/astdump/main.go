// Command astdump prints the AST of Go code snippets.
// It helps plan tests by showing how Go parses expressions, statements, and declarations.
//
// Usage:
//
//	go run ./cmd/astdump 'x + y'
//	go run ./cmd/astdump -stmt 'x = y'
//	go run ./cmd/astdump -file 'package p; func f() {}'
//	go run ./cmd/astdump -err 'func {'
//
// Flags:
//
//	-stmt	Parse as statement (wrap in func body)
//	-file	Parse as complete file
//	-err	Show parse errors (try expr, stmt, file)
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: astdump [-stmt|-file|-err] 'code'\n")
		os.Exit(1)
	}

	mode := "expr"
	code := os.Args[1]
	if len(os.Args) >= 3 {
		mode = strings.TrimPrefix(os.Args[1], "-")
		code = os.Args[2]
	}

	fset := token.NewFileSet()

	switch mode {
	case "err":
		// Try all modes and show errors
		fmt.Println("=== As expression ===")
		expr, err := parser.ParseExpr(code)
		if err != nil {
			fmt.Println("ERROR:", err)
		} else {
			ast.Print(fset, expr)
		}
		fmt.Println("\n=== As statement ===")
		src := "package p\nfunc _() {\n" + code + "\n}"
		f, err := parser.ParseFile(fset, "x.go", src, 0)
		if err != nil {
			fmt.Println("ERROR:", err)
		} else if len(f.Decls) > 0 {
			if fd, ok := f.Decls[0].(*ast.FuncDecl); ok && fd.Body != nil && len(fd.Body.List) > 0 {
				ast.Print(fset, fd.Body.List[0])
			}
		}
		fmt.Println("\n=== As file ===")
		fset2 := token.NewFileSet()
		f2, err := parser.ParseFile(fset2, "x.go", code, 0)
		if err != nil {
			fmt.Println("ERROR:", err)
		} else {
			ast.Print(fset2, f2)
		}

	case "file":
		f, err := parser.ParseFile(fset, "x.go", code, 0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
			os.Exit(1)
		}
		ast.Print(fset, f)

	case "stmt":
		src := "package p\nfunc _() {\n" + code + "\n}"
		f, err := parser.ParseFile(fset, "x.go", src, 0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
			os.Exit(1)
		}
		if len(f.Decls) > 0 {
			if fd, ok := f.Decls[0].(*ast.FuncDecl); ok && fd.Body != nil {
				for _, s := range fd.Body.List {
					ast.Print(fset, s)
					fmt.Println()
				}
			}
		}

	default: // expr
		expr, err := parser.ParseExpr(code)
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
			os.Exit(1)
		}
		ast.Print(fset, expr)
	}
	fmt.Println()
}
