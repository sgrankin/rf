// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package diff

import (
	"os"
	"testing"
)

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

func TestDiffBothEmpty(t *testing.T) {
	out, err := Diff("a", []byte{}, "b", []byte{})
	if err != nil {
		t.Fatal(err)
	}
	if out != nil {
		t.Errorf("Diff both empty: want nil, got:\n%s", out)
	}
}

func TestDiffBothNil(t *testing.T) {
	out, err := Diff("a", nil, "b", nil)
	if err != nil {
		t.Fatal(err)
	}
	if out != nil {
		t.Errorf("Diff both nil: want nil, got:\n%s", out)
	}
}

func TestWriteTempFileErrorBadDir(t *testing.T) {
	// Set TMPDIR to a non-existent directory to trigger os.CreateTemp failure (lines 59-60).
	t.Setenv("TMPDIR", "/nonexistent-dir-for-test")
	old := []byte("hello\n")
	new := []byte("world\n")
	_, err := Diff("a", old, "b", new)
	if err == nil {
		t.Fatal("expected error when TMPDIR is invalid, got nil")
	}
}

func TestDiffCommandNotFound(t *testing.T) {
	// Set PATH to empty so the diff command cannot be found (lines 33-34).
	t.Setenv("PATH", "")
	old := []byte("hello\n")
	new := []byte("world\n")
	_, err := Diff("a", old, "b", new)
	if err == nil {
		t.Fatal("expected error when diff command is not found, got nil")
	}
}

func TestWriteTempFileErrorSecondCall(t *testing.T) {
	// Use a temp dir with limited space: create the dir, write the first file,
	// then make the dir read-only so the second writeTempFile fails (lines 27-28).
	dir := t.TempDir()
	t.Setenv("TMPDIR", dir)

	// Verify that Diff works with this TMPDIR first.
	old := []byte("hello\n")
	new := []byte("world\n")
	out, err := Diff("a", old, "b", new)
	if err != nil {
		t.Fatal(err)
	}
	if out == nil {
		t.Fatal("expected non-nil diff output")
	}

	// Now remove the dir entirely and verify we get an error.
	if err := os.RemoveAll(dir); err != nil {
		t.Fatal(err)
	}

	_, err = Diff("a", old, "b", new)
	if err == nil {
		t.Fatal("expected error when TMPDIR is removed, got nil")
	}
}
