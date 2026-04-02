package refactor

import (
	"go/ast"
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
