package refactor

import (
	"go/types"
	"testing"
)

func TestEvalScopeUnhandledType(t *testing.T) {
	// evalScope should return ItemNotFound for unhandled object types
	// (e.g. *types.Builtin, *types.Label) instead of crashing.
	scope := types.Universe
	// "len" is a *types.Builtin in the universe scope.
	item := evalScope(scope, "len")
	if item.Kind != ItemNotFound {
		t.Errorf("evalScope(universe, 'len') = %v, want ItemNotFound", item.Kind)
	}
}

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
