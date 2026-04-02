package refactor

import (
	"strings"
	"testing"
)

func TestConfigString(t *testing.T) {
	c := Config{BuildTags: []string{"linux", "amd64"}}
	if s := c.String(); s != "linux,amd64" {
		t.Errorf("String() = %q, want %q", s, "linux,amd64")
	}
}

func TestConfigFlagsEnvs(t *testing.T) {
	tests := []struct {
		name      string
		tags      []string
		wantFlags []string
		wantEnvs  []string
	}{
		{
			name:     "goos",
			tags:     []string{"linux"},
			wantEnvs: []string{"GOOS=linux"},
		},
		{
			name:     "goarch",
			tags:     []string{"amd64"},
			wantEnvs: []string{"GOARCH=amd64"},
		},
		{
			name:     "cgo",
			tags:     []string{"cgo"},
			wantEnvs: []string{"CGO_ENABLED=1"},
		},
		{
			name:     "nocgo",
			tags:     []string{"!cgo"},
			wantEnvs: []string{"CGO_ENABLED=0"},
		},
		{
			name:      "race",
			tags:      []string{"race"},
			wantFlags: []string{"-race"},
		},
		{
			name:      "custom tag",
			tags:      []string{"mytag"},
			wantFlags: []string{"-tags=mytag"},
		},
		{
			name:      "combined",
			tags:      []string{"linux", "amd64", "cgo", "race", "mytag"},
			wantFlags: []string{"-race", "-tags=mytag"},
			wantEnvs:  []string{"GOOS=linux", "GOARCH=amd64", "CGO_ENABLED=1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Config{BuildTags: tt.tags}
			flags, envs, err := c.flagsEnvs()
			if err != nil {
				t.Fatal(err)
			}
			if len(flags) != len(tt.wantFlags) {
				t.Errorf("flags = %v, want %v", flags, tt.wantFlags)
			} else {
				for i, f := range flags {
					if f != tt.wantFlags[i] {
						t.Errorf("flags[%d] = %q, want %q", i, f, tt.wantFlags[i])
					}
				}
			}
			if len(envs) != len(tt.wantEnvs) {
				t.Errorf("envs = %v, want %v", envs, tt.wantEnvs)
			} else {
				for i, e := range envs {
					if e != tt.wantEnvs[i] {
						t.Errorf("envs[%d] = %q, want %q", i, e, tt.wantEnvs[i])
					}
				}
			}
		})
	}
}

func TestConfigFlagsEnvsConflict(t *testing.T) {
	c := Config{BuildTags: []string{"linux", "darwin"}}
	_, _, err := c.flagsEnvs()
	if err == nil {
		t.Fatal("expected error for conflicting GOOS values")
	}
	if !strings.Contains(err.Error(), "conflicting") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestConfigFlagsEnvsConflictGOARCH(t *testing.T) {
	c := Config{BuildTags: []string{"amd64", "arm64"}}
	_, _, err := c.flagsEnvs()
	if err == nil {
		t.Fatal("expected error for conflicting GOARCH values")
	}
	if !strings.Contains(err.Error(), "conflicting") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestConfigFlagsEnvsConflictCgo(t *testing.T) {
	c := Config{BuildTags: []string{"cgo", "!cgo"}}
	_, _, err := c.flagsEnvs()
	if err == nil {
		t.Fatal("expected error for conflicting CGO_ENABLED values")
	}
	if !strings.Contains(err.Error(), "conflicting") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestConfigFlagsEnvsDuplicateGOOS(t *testing.T) {
	c := Config{BuildTags: []string{"linux", "linux"}}
	_, envs, err := c.flagsEnvs()
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, e := range envs {
		if e == "GOOS=linux" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected 1 GOOS=linux env, got %d in %v", count, envs)
	}
}

func TestNewConfigs(t *testing.T) {
	cs := NewConfigs("linux", "amd64")
	if len(cs.c) != 1 {
		t.Fatalf("len(c) = %d, want 1", len(cs.c))
	}
	if s := cs.c[0].String(); s != "linux,amd64" {
		t.Errorf("String() = %q", s)
	}
}

func TestConfigsCross(t *testing.T) {
	cs1 := NewConfigs("linux").Plus(NewConfigs("darwin"))
	cs2 := NewConfigs("amd64").Plus(NewConfigs("arm64"))
	cs3 := cs1.Cross(cs2)
	if len(cs3.c) != 4 {
		t.Errorf("Cross: got %d configs, want 4", len(cs3.c))
	}
}

func TestConfigsPlus(t *testing.T) {
	cs1 := NewConfigs("linux")
	cs2 := NewConfigs("darwin")
	cs3 := cs1.Plus(cs2)
	if len(cs3.c) != 2 {
		t.Errorf("Plus: got %d configs, want 2", len(cs3.c))
	}
}
