// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"golang.org/x/tools/txtar"
	"rsc.io/rf/diff"
	"rsc.io/rf/refactor"
)

func TestMain(m *testing.M) {
	if os.Getenv("TEST_MAIN") == "rf" {
		main()
		return
	}
	os.Exit(m.Run())
}

// runRF invokes the test binary as the rf command with the given args,
// working directory, and stdin. Returns stdout, stderr, and exit code.
// If GOCOVERDIR is set in the environment, subprocess coverage data
// is written there automatically.
func runRF(t *testing.T, dir string, stdin string, args ...string) (string, string, int) {
	t.Helper()
	cmd := exec.Command(os.Args[0], args...)
	cmd.Env = append(os.Environ(), "TEST_MAIN=rf")
	cmd.Dir = dir
	cmd.Stdin = strings.NewReader(stdin)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			exitCode = ee.ExitCode()
		} else {
			t.Fatalf("failed to run rf: %v", err)
		}
	}
	return stdout.String(), stderr.String(), exitCode
}

var trimCommentsTests = []struct {
	in  string
	out string
}{
	{"hello", "hello"},
	{"hello # comment", "hello"},
	{"hello '#' world", "hello '#' world"},
	{`hello "#" world`, `hello "#" world`},
	{`hello \# world`, `hello \`},
	{"", ""},
	{"# all comment", ""},
	{`"hello # world"`, `"hello # world"`},
	{"`hello # world`", "`hello # world`"},
}

func TestTrimComments(t *testing.T) {
	for _, tt := range trimCommentsTests {
		out := trimComments(tt.in)
		if out != tt.out {
			t.Errorf("trimComments(%q) = %q, want %q", tt.in, out, tt.out)
		}
	}
}

var readLineTests = []struct {
	in  string
	out []string
	err error
}{
	{
		in:  "cmd x y",
		out: []string{"cmd x y"},
	},
	{
		in:  "cmd x \\\ny",
		out: []string{"cmd x \ny"},
	},
	{
		in:  "cmd x \\ # hello\ny",
		out: []string{"cmd x \ny"},
	},
	{
		in:  "cmd x y\n",
		out: []string{"cmd x y"},
	},
	{
		in:  "cmd (\nx y\n)\n",
		out: []string{"cmd (\nx y\n)"},
	},
	{
		in:  "cmd {\nx y\n}\n",
		out: []string{"cmd {\nx y\n}"},
	},
	{
		in:  "cmd It\\'s not a failure\n",
		out: []string{"cmd It's not a failure"},
	},
	{
		in:  "cmd 'It\\'s not a failure'\n",
		out: []string{"cmd 'It\\'s not a failure'"},
	},
	{
		in:  "cmd It\\\"s not a failure\n",
		out: []string{"cmd It\"s not a failure"},
	},
	{
		in:  "cmd \"It\\\"s not a failure\"\n",
		out: []string{"cmd \"It\\\"s not a failure\""},
	},
	{
		in:  "cmd It's a failure\n",
		err: fmt.Errorf("newline in '-quoted string"),
	},
	{
		in:  "cmd It\"s a failure\n",
		err: fmt.Errorf("newline in \"-quoted string"),
	},
}

func TestReadLine(t *testing.T) {
	for _, tt := range readLineTests {
		var out []string
		var err error
		text := tt.in
		for text != "" && err == nil {
			var line string
			line, text, err = readLine(text)
			if line != "" {
				out = append(out, line)
			}
		}
		if !reflect.DeepEqual(out, tt.out) || fmt.Sprint(err) != fmt.Sprint(tt.err) {
			t.Errorf("input:\n%s\nreadLine => %q, %v, want %q, %v",
				tt.in, out, err, tt.out, tt.err)
		}
	}
}

var updateTestData = flag.Bool("u", false, "update testdata instead of failing")
var flagKeep = flag.Bool("keep", false, "keep temporary work directories")

func TestRun(t *testing.T) {
	files, err := filepath.Glob("testdata/*.txt")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no test cases")
	}

	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			t.Log(file)
			ar, err := txtar.ParseFile(file)
			if err != nil {
				t.Fatal(err)
			}
			var dir string
			if *flagKeep {
				name := filepath.Base(file)
				name = strings.TrimSuffix(name, filepath.Ext(name))
				dir, err = os.MkdirTemp("", "rf-"+name)
				if err != nil {
					t.Fatal("creating work directory:", err)
				}
				t.Log("dir:", dir)
			} else {
				dir = t.TempDir()
			}
			if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module m\n"), 0666); err != nil {
				t.Fatal(err)
			}
			var wantStdout, wantStderr txtar.File
			for _, file := range ar.Files {
				if file.Name == "stdout" {
					wantStdout = file
					continue
				}
				if file.Name == "stderr" {
					wantStderr = file
					continue
				}
				targ := filepath.Join(dir, file.Name)
				if err := os.MkdirAll(filepath.Dir(targ), 0777); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(targ, file.Data, 0666); err != nil {
					t.Fatal(err)
				}
			}

			// Process flags in the comment.
			var flags flags
			flagSet := flag.NewFlagSet(filepath.Base(file), flag.ContinueOnError)
			flags.register(flagSet)
			lines := bytes.Split(ar.Comment, []byte("\n"))
			var script strings.Builder
			for _, line := range lines {
				if len(line) > 0 && line[0] == '-' {
					// Flags
					if err := flagSet.Parse(strings.Fields(string(line))); err != nil {
						t.Fatal(err)
					}
				} else {
					script.Write(line)
					script.WriteByte('\n')
				}
			}

			var stdout, stderr bytes.Buffer
			defer func() {
				// Flush stderr to the test log on panic.
				if stderr.Len() == 0 {
					return
				}
				if err := recover(); err != nil {
					s := stderr.String()
					for len(s) > 0 && s != "\n" {
						var line string
						line, s, _ = strings.Cut(s, "\n")
						t.Logf("stderr: %s", line)
					}
					panic(err)
				}
			}()
			rf, err := refactor.New(dir)
			if err != nil {
				t.Fatal(err)
			}
			rf.Stdout = &stdout
			rf.Stderr = &stderr
			flags.apply(rf)
			rf.ShowDiff = true
			if err := run(rf, script.String()); err != nil {
				fmt.Fprintf(rf.Stderr, "%v\n", err)
			}

			if *updateTestData {
				stderrChanged := updateFile(ar, "stderr", stderr.Bytes())
				stdoutChanged := updateFile(ar, "stdout", stdout.Bytes())
				if stdoutChanged || stderrChanged {
					if err := os.WriteFile(file, txtar.Format(ar), 0666); err != nil {
						t.Fatal(err)
					}
					t.Log("updated")
				}
				return
			}

			cmp := func(name string, have, want []byte) {
				have = trimSpace(have)
				want = trimSpace(want)
				if !bytes.Equal(have, want) {
					t.Errorf("%s:\n%s", name, have)
					t.Errorf("want:\n%s", want)
					d, err := diff.Diff("want", want, "have", have)
					if err == nil {
						t.Errorf("diff of diffs:\n%s", d)
					}
				}
			}
			cmp("stderr", stderr.Bytes(), wantStderr.Data)
			cmp("stdout", stdout.Bytes(), wantStdout.Data)
		})
	}
}

// TestMainNoArgs tests that rf with no arguments prints usage and exits 2.
func TestMainNoArgs(t *testing.T) {
	_, stderr, code := runRF(t, t.TempDir(), "")
	if code != 2 {
		t.Errorf("exit code = %d, want 2", code)
	}
	if !strings.Contains(stderr, "usage:") {
		t.Errorf("stderr = %q, want usage message", stderr)
	}
}

// TestMainUnknownCommand tests that rf with an unknown command fails.
func TestMainUnknownCommand(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module m\n"), 0666); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "x.go"), []byte("package m\n"), 0666); err != nil {
		t.Fatal(err)
	}
	_, stderr, code := runRF(t, dir, "", "nosuchcommand x y")
	if code == 0 {
		t.Error("expected non-zero exit code")
	}
	if !strings.Contains(stderr, "unknown command") {
		t.Errorf("stderr = %q, want 'unknown command'", stderr)
	}
}

// TestMainDiffFlag tests that -diff flag works.
func TestMainDiffFlag(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module m\n"), 0666); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "x.go"), []byte("package m\n\nvar X int\n"), 0666); err != nil {
		t.Fatal(err)
	}
	stdout, _, code := runRF(t, dir, "", "-diff", "mv X Y")
	if code != 0 {
		t.Errorf("exit code = %d, want 0", code)
	}
	if !strings.Contains(stdout, "diff") {
		t.Errorf("stdout = %q, want diff output", stdout)
	}
}

func trimSpace(data []byte) []byte {
	lines := bytes.Split(data, []byte("\n"))
	for i, line := range lines {
		lines[i] = bytes.TrimRight(line, " ")
	}
	return bytes.Join(lines, []byte("\n"))
}

func updateFile(ar *txtar.Archive, name string, data []byte) bool {
	data = trimSpace(data)

	for i := range ar.Files {
		if file := &ar.Files[i]; file.Name == name {
			if len(data) == 0 {
				ar.Files = append(ar.Files[:i], ar.Files[i+1:]...)
				return true
			}

			if !bytes.Equal(file.Data, data) {
				file.Data = data
				return true
			}

			return false
		}
	}

	if len(data) != 0 {
		ar.Files = append(ar.Files, txtar.File{
			Name: name,
			Data: data,
		})
		return true
	}

	return false
}
