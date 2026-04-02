package refactor

import (
	"go/token"
	"strings"
	"testing"
)

func TestDepsGraphAddEmptyFrom(t *testing.T) {
	g := &DepsGraph{Level: SymRefs}
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for empty from")
		}
		if !strings.Contains(r.(string), "empty from") {
			t.Fatalf("unexpected panic: %v", r)
		}
	}()
	pkg := &Package{PkgPath: "example.com/foo"}
	g.add(QualName{Pkg: pkg, Name: ""}, QualName{Pkg: pkg, Name: "bar"})
}

func TestDepsGraphAddEmptyTo(t *testing.T) {
	g := &DepsGraph{Level: SymRefs}
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for empty to")
		}
		if !strings.Contains(r.(string), "empty to") {
			t.Fatalf("unexpected panic: %v", r)
		}
	}()
	pkg := &Package{PkgPath: "example.com/foo"}
	g.add(QualName{Pkg: pkg, Name: "foo"}, QualName{Pkg: pkg, Name: ""})
}

func TestDepsGraphAddSamePkgFiltered(t *testing.T) {
	// At PkgRefs level, same-package references are still added
	// (names are zeroed). At lower levels, same-package refs are skipped.
	g := &DepsGraph{Level: PkgRefs - 1}
	pkg := &Package{PkgPath: "example.com/foo"}
	// Should not panic — just return without adding.
	g.add(QualName{Pkg: pkg, Name: "foo"}, QualName{Pkg: pkg, Name: "bar"})
}

func TestAddDepsNilTypesInfo(t *testing.T) {
	fset := token.NewFileSet()
	s := &Snapshot{fset: fset}
	g := &DepsGraph{Level: SymRefs}
	pkg := &Package{PkgPath: "example.com/foo"}
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for nil TypesInfo")
		}
		if !strings.Contains(r.(string), "no TypesInfo") {
			t.Fatalf("unexpected panic: %v", r)
		}
	}()
	s.addDeps(g, QualName{Pkg: pkg, Name: "Foo"}, pkg, nil)
}

func TestQualNameString(t *testing.T) {
	// Nil package.
	q := QualName{}
	if got := q.String(); got != "<q>" {
		t.Errorf("QualName{}.String() = %q, want %q", got, "<q>")
	}
	// With package.
	pkg := &Package{PkgPath: "example.com/foo"}
	q = QualName{Pkg: pkg, Name: "Bar"}
	if got := q.String(); got != "example.com/foo.Bar" {
		t.Errorf("QualName.String() = %q, want %q", got, "example.com/foo.Bar")
	}
}
