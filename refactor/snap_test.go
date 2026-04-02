package refactor

import "testing"

func TestImportMapLookup(t *testing.T) {
	im := importMap{"vendor/x": "x"}
	if got := im.Lookup("vendor/x"); got != "x" {
		t.Errorf("Lookup(vendor/x) = %q, want %q", got, "x")
	}
	if got := im.Lookup("other"); got != "other" {
		t.Errorf("Lookup(other) = %q, want %q", got, "other")
	}
}

func TestStringList(t *testing.T) {
	got := stringList("a", []string{"b", "c"}, "d")
	want := []string{"a", "b", "c", "d"}
	if len(got) != len(want) {
		t.Fatalf("stringList() = %v, want %v", got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("stringList()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestStringListPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for invalid type")
		}
	}()
	stringList(42)
}

func TestStringListEmpty(t *testing.T) {
	got := stringList()
	if len(got) != 0 {
		t.Errorf("stringList() = %v, want empty", got)
	}
}
