#!/usr/bin/env bash
# Usage: ./astdump.sh 'go code expression or statement'
# Prints the AST that Go parses the input to.
# Useful for understanding what AST nodes rf patterns will produce.
cat > /tmp/astdump.go << 'GOEOF'
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
		fmt.Fprintf(os.Stderr, "usage: astdump <expr-or-stmt>\n")
		os.Exit(1)
	}
	input := strings.Join(os.Args[1:], " ")

	fset := token.NewFileSet()

	// Try as expression first
	expr, err := parser.ParseExpr(input)
	if err == nil {
		fmt.Println("Parsed as expression:")
		ast.Print(fset, expr)
		return
	}

	// Try as statement(s) in a function body
	src := fmt.Sprintf("package p\nfunc _() {\n%s\n}", input)
	f, err2 := parser.ParseFile(fset, "input.go", src, 0)
	if err2 == nil {
		body := f.Decls[0].(*ast.FuncDecl).Body.List
		if len(body) == 1 {
			fmt.Println("Parsed as statement:")
			ast.Print(fset, body[0])
		} else {
			fmt.Printf("Parsed as %d statements:\n", len(body))
			for i, s := range body {
				fmt.Printf("--- statement %d ---\n", i)
				ast.Print(fset, s)
			}
		}
		return
	}

	// Try as top-level declaration
	src2 := fmt.Sprintf("package p\n%s", input)
	f2, err3 := parser.ParseFile(fset, "input.go", src2, 0)
	if err3 == nil {
		for i, d := range f2.Decls {
			fmt.Printf("--- decl %d ---\n", i)
			ast.Print(fset, d)
		}
		return
	}

	fmt.Fprintf(os.Stderr, "as expr: %v\nas stmt: %v\nas decl: %v\n", err, err2, err3)
	os.Exit(1)
}
GOEOF
go run /tmp/astdump.go "$@"
