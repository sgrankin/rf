package refactor

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"
	"testing"
)

func TestWalkNil(t *testing.T) {
	called := false
	Walk(nil, func(stack []ast.Node) {
		called = true
	})
	if called {
		t.Error("Walk(nil) should not call visitor")
	}
}

func TestWalkNilFile(t *testing.T) {
	called := false
	Walk((*ast.File)(nil), func(stack []ast.Node) {
		called = true
	})
	if called {
		t.Error("Walk(nil file) should not call visitor")
	}
}

func TestWalkNilBlockStmt(t *testing.T) {
	called := false
	Walk((*ast.BlockStmt)(nil), func(stack []ast.Node) {
		called = true
	})
	if called {
		t.Error("Walk(nil block) should not call visitor")
	}
}

func TestWalkPostNil(t *testing.T) {
	called := false
	WalkPost(nil, func(stack []ast.Node) {
		called = true
	})
	if called {
		t.Error("WalkPost(nil) should not call visitor")
	}
}

func TestWalkRangePostNil(t *testing.T) {
	called := false
	WalkRangePost(nil, 0, 100, func(stack []ast.Node) {
		called = true
	})
	if called {
		t.Error("WalkRangePost(nil) should not call visitor")
	}
}

func TestTextPanicMissingFile(t *testing.T) {
	fset := token.NewFileSet()
	f := fset.AddFile("test.go", -1, 100)
	_ = f
	s := &Snapshot{fset: fset, files: map[string]*File{}}
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic")
		}
		if !strings.Contains(r.(string), "file not found") {
			t.Fatalf("unexpected panic: %v", r)
		}
	}()
	s.Text(token.Pos(1), token.Pos(10))
}

func TestEvalPanicNoFileScope(t *testing.T) {
	fset := token.NewFileSet()
	f := fset.AddFile("test.go", -1, 100)
	_ = f
	syntax := &ast.File{Name: ast.NewIdent("m")}
	pkg := &Package{
		Files:     []*File{{Name: "test.go", Syntax: syntax}},
		TypesInfo: &types.Info{Scopes: map[ast.Node]*types.Scope{}},
	}
	s := &Snapshot{
		fset:  fset,
		files: map[string]*File{"test.go": {Name: "test.go", Text: []byte("package m\n")}},
		edits: map[string]*Edit{},
	}
	s.target = pkg
	s.packages = []*Package{pkg}
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for missing file scope")
		}
		if !strings.Contains(r.(string), "no file scope") {
			t.Fatalf("unexpected panic: %v", r)
		}
	}()
	s.Eval("nonExistentName")
}

func TestEditAtPanicMissingFile(t *testing.T) {
	fset := token.NewFileSet()
	f := fset.AddFile("test.go", -1, 100)
	_ = f
	s := &Snapshot{fset: fset, files: map[string]*File{}, edits: map[string]*Edit{}}
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic")
		}
		if !strings.Contains(r.(string), "file not found") {
			t.Fatalf("unexpected panic: %v", r)
		}
	}()
	s.editAt(token.Pos(1))
}
