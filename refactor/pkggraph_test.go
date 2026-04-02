package refactor

import (
	"strings"
	"testing"
)

func TestPkgGraphAddDuplicate(t *testing.T) {
	g := newPkgGraph("test")
	p := &Package{PkgPath: "example.com/foo", ID: "foo1"}
	g.add(p)
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for duplicate package")
		}
		if !strings.Contains(r.(string), "duplicate package path") {
			t.Fatalf("unexpected panic: %v", r)
		}
	}()
	p2 := &Package{PkgPath: "example.com/foo", ID: "foo2"}
	g.add(p2)
}

func TestPkgGraphMergeDisjoint(t *testing.T) {
	g1 := newPkgGraph("g1")
	g2 := newPkgGraph("g2")
	// Same package path, but disjoint files and imports — neither is a superset.
	p1 := &Package{
		PkgPath: "example.com/foo",
		ID:      "foo",
		Files:   []*File{{Name: "a.go"}},
		Imports: []string{"example.com/bar"},
	}
	p2 := &Package{
		PkgPath: "example.com/foo",
		ID:      "foo",
		Files:   []*File{{Name: "b.go"}},
		Imports: []string{"example.com/baz"},
	}
	g1.add(p1)
	g2.add(p2)
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for disjoint packages")
		}
		if !strings.Contains(r.(string), "disjoint") {
			t.Fatalf("unexpected panic: %v", r)
		}
	}()
	g1.merge(g2)
}

func TestPkgGraphMergeSuperset(t *testing.T) {
	g1 := newPkgGraph("g1")
	g2 := newPkgGraph("g2")
	// g2's version is a superset of g1's (has all files + extras).
	p1 := &Package{
		PkgPath: "example.com/foo",
		ID:      "foo",
		Files:   []*File{{Name: "a.go"}},
		Imports: []string{"example.com/bar"},
	}
	p2 := &Package{
		PkgPath: "example.com/foo",
		ID:      "foo",
		Files:   []*File{{Name: "a.go"}, {Name: "b.go"}},
		Imports: []string{"example.com/bar", "example.com/baz"},
	}
	g1.add(p1)
	g2.add(p2)
	g3 := g1.merge(g2)
	// g3 should have the superset version.
	merged := g3.byPath("example.com/foo")
	if len(merged.Files) != 2 {
		t.Errorf("merged has %d files, want 2", len(merged.Files))
	}
	if len(merged.Imports) != 2 {
		t.Errorf("merged has %d imports, want 2", len(merged.Imports))
	}
}

func TestPkgGraphFindCycleMissingImport(t *testing.T) {
	g := newPkgGraph("test")
	p := &Package{
		PkgPath: "example.com/foo",
		ID:      "foo",
		Imports: []string{"example.com/bar"}, // bar not in graph
	}
	g.add(p)
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for missing import")
		}
		if !strings.Contains(r.(string), "missing from the package graph") {
			t.Fatalf("unexpected panic: %v", r)
		}
	}()
	g.findCycle(false)
}

func TestPkgGraphFindCycleDiagnose(t *testing.T) {
	g := newPkgGraph("test")
	// Create a cycle: foo -> bar -> foo
	foo := &Package{PkgPath: "foo", ID: "foo", Imports: []string{"bar"}}
	bar := &Package{PkgPath: "bar", ID: "bar", Imports: []string{"foo"}}
	g.add(foo)
	g.add(bar)
	cycle := g.findCycle(true)
	if cycle == nil {
		t.Fatal("expected cycle, got nil")
	}
	s := cycle.String()
	if !strings.Contains(s, "foo") || !strings.Contains(s, "bar") {
		t.Errorf("cycle.Error() = %q, want to mention foo and bar", s)
	}
}

func TestPkgGraphVisitBottomUpCycle(t *testing.T) {
	g := newPkgGraph("test")
	foo := &Package{PkgPath: "foo", ID: "foo", Imports: []string{"bar"}}
	bar := &Package{PkgPath: "bar", ID: "bar", Imports: []string{"foo"}}
	g.add(foo)
	g.add(bar)
	err := g.visitBottomUp(func(p *Package) error {
		return nil
	})
	if err == nil {
		t.Fatal("expected error for cycle")
	}
}

func TestPkgGraphDump(t *testing.T) {
	g := newPkgGraph("test")
	p := &Package{
		PkgPath: "example.com/foo",
		ID:      "foo",
		Files:   []*File{{Name: "a.go"}},
		Imports: []string{"example.com/bar"},
	}
	g.add(p)
	var buf strings.Builder
	g.dump(&buf)
	s := buf.String()
	if !strings.Contains(s, "example.com/foo") {
		t.Errorf("dump missing package path: %s", s)
	}
	if !strings.Contains(s, "example.com/bar") {
		t.Errorf("dump missing import: %s", s)
	}
}
