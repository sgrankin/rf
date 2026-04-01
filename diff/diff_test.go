// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package diff

import "testing"

const (
	oldName = "a/b/c"
	newName = "d/e/f"
	oldText = "abc\ndef\nghi\n"
	newText = "ABC\ndef\nGHI\n"
	want    = "diff a/b/c d/e/f\n--- a/b/c\n+++ d/e/f\n@@ -1,3 +1,3 @@\n-abc\n+ABC\n def\n-ghi\n+GHI\n"
)

func TestDiff(t *testing.T) {
	out, err := Diff(oldName, []byte(oldText), newName, []byte(newText))
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != want {
		t.Errorf("Diff: have:\n%s", out)
		t.Errorf("Diff: want:\n%s", want)
	}
}

func TestDiffIdentical(t *testing.T) {
	text := []byte("abc\ndef\n")
	out, err := Diff("a", text, "b", text)
	if err != nil {
		t.Fatal(err)
	}
	if out != nil {
		t.Errorf("Diff identical: want nil, got:\n%s", out)
	}
}

func TestDiffEmpty(t *testing.T) {
	out, err := Diff("a", nil, "b", []byte("abc\n"))
	if err != nil {
		t.Fatal(err)
	}
	if out == nil {
		t.Fatal("Diff empty vs non-empty: want non-nil")
	}
}

func TestDiffSingleLine(t *testing.T) {
	out, err := Diff("a", []byte("x\n"), "b", []byte("y\n"))
	if err != nil {
		t.Fatal(err)
	}
	if out == nil {
		t.Fatal("Diff single line: want non-nil")
	}
}

func TestDiffNoNewline(t *testing.T) {
	out, err := Diff("a", []byte("x"), "b", []byte("y"))
	if err != nil {
		t.Fatal(err)
	}
	if out == nil {
		t.Fatal("Diff no newline: want non-nil")
	}
}

func TestDiffLarger(t *testing.T) {
	old := []byte("line1\nline2\nline3\nline4\nline5\n")
	new := []byte("line1\nLINE2\nline3\nline4\nLINE5\n")
	out, err := Diff("old.txt", old, "new.txt", new)
	if err != nil {
		t.Fatal(err)
	}
	if out == nil {
		t.Fatal("Diff larger: want non-nil")
	}
}
