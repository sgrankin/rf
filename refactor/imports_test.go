package refactor

import (
	"go/ast"
	"go/token"
	"testing"
)

func TestImportName(t *testing.T) {
	// Named import
	named := &ast.ImportSpec{
		Name: &ast.Ident{Name: "foo"},
		Path: &ast.BasicLit{Value: `"example.com/foo"`},
	}
	if got := importName(named); got != "foo" {
		t.Errorf("importName(named) = %q, want %q", got, "foo")
	}

	// Unnamed import
	unnamed := &ast.ImportSpec{
		Path: &ast.BasicLit{Value: `"example.com/bar"`},
	}
	if got := importName(unnamed); got != "" {
		t.Errorf("importName(unnamed) = %q, want %q", got, "")
	}
}

func TestImportPath(t *testing.T) {
	good := &ast.ImportSpec{Path: &ast.BasicLit{Value: `"fmt"`}}
	if got := importPath(good); got != "fmt" {
		t.Errorf("importPath(good) = %q, want %q", got, "fmt")
	}

	bad := &ast.ImportSpec{Path: &ast.BasicLit{Value: `bad`}}
	if got := importPath(bad); got != "" {
		t.Errorf("importPath(bad) = %q, want %q", got, "")
	}
}

func TestDeclImports(t *testing.T) {
	imp := &ast.GenDecl{
		Tok: token.IMPORT,
		Specs: []ast.Spec{
			&ast.ImportSpec{Path: &ast.BasicLit{Value: `"fmt"`}},
			&ast.ImportSpec{Path: &ast.BasicLit{Value: `"os"`}},
		},
	}
	if !declImports(imp, "fmt") {
		t.Error("declImports should find fmt")
	}
	if declImports(imp, "io") {
		t.Error("declImports should not find io")
	}

	// Non-import decl
	nonImp := &ast.GenDecl{Tok: token.VAR}
	if declImports(nonImp, "fmt") {
		t.Error("declImports should return false for VAR decl")
	}
}

func TestMatchLen(t *testing.T) {
	tests := []struct {
		x, y string
		want int
	}{
		{"fmt", "fmt", 3},
		{"fmt", "foo", 1},
		{"fmt", "os", 0},
		{"github.com/a", "github.com/b", 11},
		// Different pathKinds
		{"fmt", "github.com/a", -1},
		{"cmd/go", "fmt", -1},
	}
	for _, tt := range tests {
		if got := matchLen(tt.x, tt.y); got != tt.want {
			t.Errorf("matchLen(%q, %q) = %d, want %d", tt.x, tt.y, got, tt.want)
		}
	}
}

func TestPathKind(t *testing.T) {
	tests := []struct {
		x    string
		want int
	}{
		{"fmt", 0},
		{"cmd/go", 1},
		{"github.com/foo", 2},
		{"golang.org/x/tools", 2},
	}
	for _, tt := range tests {
		if got := pathKind(tt.x); got != tt.want {
			t.Errorf("pathKind(%q) = %d, want %d", tt.x, got, tt.want)
		}
	}
}
