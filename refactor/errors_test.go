package refactor

import (
	"go/token"
	"go/types"
	"strings"
	"testing"
)

func TestErrorString(t *testing.T) {
	e := &Error{Pos: token.Position{Filename: "x.go", Line: 1, Column: 1, Offset: 0}, Msg: "bad"}
	if got := e.Error(); got != "x.go:1:1: bad" {
		t.Errorf("Error() = %q", got)
	}

	e2 := &Error{Msg: "no pos"}
	if got := e2.Error(); got != "no pos" {
		t.Errorf("Error() = %q", got)
	}
}

func TestErrorListEmpty(t *testing.T) {
	var l ErrorList
	if s := l.Error(); s != "no errors" {
		t.Errorf("Error() = %q, want %q", s, "no errors")
	}
	if l.Err() != nil {
		t.Error("Err() should be nil for empty list")
	}
}

func TestErrorListAddTypesError(t *testing.T) {
	var l ErrorList
	fset := token.NewFileSet()
	f := fset.AddFile("x.go", -1, 100)
	pos := f.Pos(10)

	// Primary error.
	l.Add(types.Error{Fset: fset, Pos: pos, Msg: "primary error"})
	// Secondary error (starts with tab).
	l.Add(types.Error{Fset: fset, Pos: pos, Msg: "\tsecondary detail"})

	s := l.Error()
	if !strings.Contains(s, "primary error") {
		t.Errorf("missing primary error in %q", s)
	}
	if !strings.Contains(s, "secondary detail") {
		t.Errorf("missing secondary error in %q", s)
	}
}

func TestErrorListDuplicateCollapse(t *testing.T) {
	var l ErrorList
	// Add the same message at 4 different positions — should collapse.
	for i := 0; i < 4; i++ {
		l.Add(&Error{
			Pos: token.Position{Filename: "x.go", Line: i + 1, Column: 1, Offset: i * 10},
			Msg: "repeated error",
		})
	}

	s := l.Error()
	if !strings.Contains(s, "× 4") {
		t.Errorf("expected collapse marker in %q", s)
	}
	// Should only appear once (collapsed).
	if strings.Count(s, "repeated error") != 1 {
		t.Errorf("expected single occurrence of collapsed error, got:\n%s", s)
	}
}

func TestErrorListSorting(t *testing.T) {
	var l ErrorList
	l.Add(&Error{Pos: token.Position{Filename: "b.go", Offset: 10}, Msg: "error b"})
	l.Add(&Error{Pos: token.Position{Filename: "a.go", Offset: 5}, Msg: "error a"})

	s := l.Error()
	aIdx := strings.Index(s, "error a")
	bIdx := strings.Index(s, "error b")
	if aIdx > bIdx {
		t.Errorf("expected error a before error b:\n%s", s)
	}
}
