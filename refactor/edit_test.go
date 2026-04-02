package refactor

import (
	"go/token"
	"sort"
	"strings"
	"testing"
)

func TestEditsSort(t *testing.T) {
	q := edits{
		{pos: 10, end: 20, force: false},
		{pos: 5, end: 15, force: false},
		{pos: 10, end: 20, force: true},
		{pos: 10, end: 25, force: false},
	}
	sort.Sort(q)
	// pos=5 first, then pos=10 force=true, then two pos=10 force=false
	if q[0].pos != 5 {
		t.Errorf("q[0].pos = %d, want 5", q[0].pos)
	}
	if q[1].pos != 10 || !q[1].force {
		t.Errorf("q[1] = {pos:%d force:%v}, want {pos:10 force:true}", q[1].pos, q[1].force)
	}
	if q[2].force || q[3].force {
		t.Error("q[2] and q[3] should have force=false")
	}
}

func TestBufferBytes(t *testing.T) {
	text := []byte("hello world")
	b := &Buffer{pos: 1, end: token.Pos(1 + len(text)), old: text}
	if got := string(b.Bytes()); got != "hello world" {
		t.Errorf("Bytes() = %q, want %q", got, "hello world")
	}
}

func TestBufferEditReplace(t *testing.T) {
	text := []byte("hello world")
	b := &Buffer{pos: 1, end: token.Pos(1 + len(text)), old: text}
	// Replace "hello" (positions 1-6) with "goodbye"
	b.edit(1, 6, "goodbye", false)
	got := string(b.Bytes())
	if got != "goodbye world" {
		t.Errorf("Bytes() after edit = %q, want %q", got, "goodbye world")
	}
}

func TestBufferEditPanicEndGtBufEnd(t *testing.T) {
	text := []byte("hello")
	b := &Buffer{pos: 1, end: token.Pos(1 + len(text)), old: text}
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for invalid edit position")
		}
	}()
	b.edit(1, token.Pos(100), "x", false) // end > b.end
}

func TestBufferEditPanicEndLtPos(t *testing.T) {
	text := []byte("hello")
	b := &Buffer{pos: 1, end: token.Pos(1 + len(text)), old: text}
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for end < pos")
		}
	}()
	b.edit(5, 3, "x", false) // end < pos
}

func TestCreateFileParsePanic(t *testing.T) {
	fset := token.NewFileSet()
	cache := &buildCache{
		fset:       fset,
		files: make(map[string]*File),
	}
	r := &Refactor{cache: cache}
	cache.r = r
	s := &Snapshot{
		fset:  fset,
		r:     r,
		files: map[string]*File{},
		edits: map[string]*Edit{},
	}
	p := &Package{Dir: "/tmp/test", Name: "m"}
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for parse error")
		}
		if !strings.Contains(r.(string), "CreateFile parse") {
			t.Fatalf("unexpected panic: %v", r)
		}
	}()
	s.CreateFile(p, "bad.go", "not valid go {{{{")
}

func TestCreateFileDefaultText(t *testing.T) {
	fset := token.NewFileSet()
	cache := &buildCache{
		fset:       fset,
		files: make(map[string]*File),
	}
	r := &Refactor{cache: cache}
	cache.r = r
	s := &Snapshot{
		fset:  fset,
		r:     r,
		files: map[string]*File{},
		edits: map[string]*Edit{},
	}
	p := &Package{Dir: "/tmp/test", Name: "mypkg"}
	f := s.CreateFile(p, "new.go", "")
	if f == nil {
		t.Fatal("CreateFile returned nil")
	}
	if f.Name.Name != "mypkg" {
		t.Errorf("package name = %q, want %q", f.Name.Name, "mypkg")
	}
}

func TestCreateFileDuplicatePanic(t *testing.T) {
	fset := token.NewFileSet()
	cache := &buildCache{
		fset:       fset,
		files: make(map[string]*File),
	}
	r := &Refactor{cache: cache}
	cache.r = r
	s := &Snapshot{
		fset:  fset,
		r:     r,
		files: map[string]*File{},
		edits: map[string]*Edit{},
	}
	p := &Package{Dir: "/tmp/test", Name: "m"}
	s.CreateFile(p, "x.go", "package m\n")
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for duplicate file")
		}
		if !strings.Contains(r.(string), "created twice") {
			t.Fatalf("unexpected panic: %v", r)
		}
	}()
	s.CreateFile(p, "x.go", "package m\n")
}
