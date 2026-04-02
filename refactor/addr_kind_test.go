package refactor

import "testing"

func TestItemKindString(t *testing.T) {
	tests := []struct {
		kind ItemKind
		want string
	}{
		{ItemNotFound, "not found"},
		{ItemFile, "file"},
		{ItemDir, "dir"},
		{ItemConst, "const"},
		{ItemType, "type"},
		{ItemVar, "var"},
		{ItemFunc, "func"},
		{ItemField, "field"},
		{ItemMethod, "method"},
		{ItemPos, "text"},
		{ItemPkg, "pkg"},
		{ItemKind(99), "???"},
	}
	for _, tt := range tests {
		if got := tt.kind.String(); got != tt.want {
			t.Errorf("ItemKind(%d).String() = %q, want %q", tt.kind, got, tt.want)
		}
	}
}
