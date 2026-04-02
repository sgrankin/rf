package refactor

import (
	"go/token"
	"sort"
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
