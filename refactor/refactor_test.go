package refactor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	dir := t.TempDir()
	r, err := New(dir)
	if err != nil {
		t.Fatal(err)
	}
	if r.ModRoot() != "" {
		t.Errorf("ModRoot() = %q, want empty", r.ModRoot())
	}
	if r.ModPath() != "" {
		t.Errorf("ModPath() = %q, want empty", r.ModPath())
	}
}

func TestNewNotDir(t *testing.T) {
	f, err := os.CreateTemp("", "rf-test")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove(f.Name())

	_, err = New(f.Name())
	if err == nil {
		t.Fatal("expected error for non-directory")
	}
}

func TestNewNotExist(t *testing.T) {
	_, err := New("/nonexistent/path/that/does/not/exist")
	if err == nil {
		t.Fatal("expected error for nonexistent path")
	}
}

func TestPkgDir(t *testing.T) {
	dir := t.TempDir()
	r, err := New(dir)
	if err != nil {
		t.Fatal(err)
	}
	abs, _ := filepath.Abs(dir)
	r.modRoot = abs
	r.modPath = "example.com/mod"

	tests := []struct {
		pkg     string
		wantDir string
		wantErr bool
	}{
		{"example.com/mod", abs, false},
		{"example.com/mod/sub", filepath.Join(abs, "sub"), false},
		{"example.com/mod/a/b", filepath.Join(abs, "a/b"), false},
		{"other.com/pkg", "", true},
		{"example.com/modextra", "", true},
	}
	for _, tt := range tests {
		d, err := r.PkgDir(tt.pkg)
		if (err != nil) != tt.wantErr {
			t.Errorf("PkgDir(%q) error = %v, wantErr %v", tt.pkg, err, tt.wantErr)
			continue
		}
		if d != tt.wantDir {
			t.Errorf("PkgDir(%q) = %q, want %q", tt.pkg, d, tt.wantDir)
		}
	}
}

func TestCutLast(t *testing.T) {
	tests := []struct {
		s, sep     string
		before     string
		after      string
		ok         bool
	}{
		{"a.b.c", ".", "a.b", "c", true},
		{"abc", ".", "abc", "", false},
		{"a/b/c", "/", "a/b", "c", true},
	}
	for _, tt := range tests {
		before, after, ok := cutLast(tt.s, tt.sep)
		if before != tt.before || after != tt.after || ok != tt.ok {
			t.Errorf("cutLast(%q, %q) = (%q, %q, %v), want (%q, %q, %v)",
				tt.s, tt.sep, before, after, ok, tt.before, tt.after, tt.ok)
		}
	}
}
