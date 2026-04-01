// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"go/ast"
	"reflect"
	"testing"
)

var commonRangesTests = []struct {
	x   string
	y   string
	out []rangePair
}{
	{"", "", nil},
	{"x", "x", []rangePair{{0, 0, 1}}},
	{"x", "xy", []rangePair{{0, 0, 1}}},
	{"xy", "x", []rangePair{{0, 0, 1}}},
	{"x", "yx", []rangePair{{0, 1, 1}}},
	{"yx", "x", []rangePair{{1, 0, 1}}},
	{"zx", "zx", []rangePair{{0, 0, 2}}},
	{"zx", "zxy", []rangePair{{0, 0, 2}}},
	{"zxy", "zx", []rangePair{{0, 0, 2}}},
	{"zx", "zyx", []rangePair{{0, 0, 1}, {1, 2, 1}}},
	{"zyx", "zx", []rangePair{{0, 0, 1}, {2, 1, 1}}},
	{"zx", "wx", []rangePair{{1, 1, 1}}},
	{"zx", "wxy", []rangePair{{1, 1, 1}}},
	{"zxy", "wx", []rangePair{{1, 1, 1}}},
	{"zx", "wyx", []rangePair{{1, 2, 1}}},
	{"zyx", "wx", []rangePair{{2, 1, 1}}},
	{"a", "b", nil},
}

func TestTrimCommentsEscape(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{`hello # comment`, "hello"},
		{`"hello # world"`, `"hello # world"`},
		{`"esc\"ape" # comment`, `"esc\"ape"`},
		{`'esc\'ape' # comment`, `'esc\'ape'`},
		{"no comment", "no comment"},
	}
	for _, tt := range tests {
		got := trimComments(tt.in)
		if got != tt.want {
			t.Errorf("trimComments(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestImportPath(t *testing.T) {
	tests := []struct {
		value string
		want  string
	}{
		{`"fmt"`, "fmt"},
		{`"net/http"`, "net/http"},
		{`bad`, ""},
	}
	for _, tt := range tests {
		spec := &ast.ImportSpec{Path: &ast.BasicLit{Value: tt.value}}
		got := importPath(spec)
		if got != tt.want {
			t.Errorf("importPath(%q) = %q, want %q", tt.value, got, tt.want)
		}
	}
}

func TestCommonRanges(t *testing.T) {
	for _, tt := range commonRangesTests {
		list := commonRanges(tt.x, tt.y)
		if !reflect.DeepEqual(list, tt.out) {
			t.Errorf("commonRanges(%q, %q) = %v, want %v", tt.x, tt.y, list, tt.out)
		}
	}
}
