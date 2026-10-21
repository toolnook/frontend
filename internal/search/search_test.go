package search

import (
	"testing"

	"github.com/forgekit/docforge/internal/types"
)

func TestBuildIndex(t *testing.T) {
	project := types.Project{
		Modules: []types.Module{
			{
				Name: "utils.ts",
				Documents: []types.Document{
					{Name: "add", Kind: types.KindFunction, Description: "Add two numbers"},
					{Name: "subtract", Kind: types.KindFunction, Description: "Subtract two numbers"},
				},
			},
			{
				Name: "users.ts",
				Documents: []types.Document{
					{Name: "User", Kind: types.KindClass, Description: "User model"},
				},
			},
		},
	}
	idx := BuildIndex(project)
	if idx.EntryCount() != 3 {
		t.Errorf("expected 3 entries, got %d", idx.EntryCount())
	}
}

func TestSearchByName(t *testing.T) {
	idx := Index{
		Entries: []Entry{
			{Name: "add", Module: "utils.ts", Kind: "function"},
			{Name: "subtract", Module: "utils.ts", Kind: "function"},
			{Name: "User", Module: "users.ts", Kind: "class"},
		},
	}
	results := idx.Search("add")
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Name != "add" {
		t.Errorf("expected 'add', got '%s'", results[0].Name)
	}
}

func TestSearchByDescription(t *testing.T) {
	idx := Index{
		Entries: []Entry{
			{Name: "add", Description: "Add two numbers"},
			{Name: "concat", Description: "Concatenate two strings"},
		},
	}
	results := idx.Search("numbers")
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Name != "add" {
		t.Errorf("expected 'add', got '%s'", results[0].Name)
	}
}

func TestSearchEmptyQueryReturnsNil(t *testing.T) {
	idx := Index{
		Entries: []Entry{{Name: "foo"}},
	}
	if idx.Search("") != nil {
		t.Error("expected nil for empty query")
	}
	if idx.Search("   ") != nil {
		t.Error("expected nil for whitespace query")
	}
}

func TestSearchByKind(t *testing.T) {
	idx := Index{
		Entries: []Entry{
			{Name: "add", Kind: "function"},
			{Name: "User", Kind: "class"},
			{Name: "Status", Kind: "type"},
		},
	}
	results := idx.SearchByKind("function")
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestSearchByModule(t *testing.T) {
	idx := Index{
		Entries: []Entry{
			{Name: "add", Module: "utils.ts"},
			{Name: "User", Module: "users.ts"},
		},
	}
	results := idx.SearchByModule("utils.ts")
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestContains(t *testing.T) {
	idx := Index{
		Entries: []Entry{{Name: "add"}, {Name: "subtract"}},
	}
	if !idx.Contains("add") {
		t.Error("expected Contains('add') to return true")
	}
	if idx.Contains("multiply") {
		t.Error("expected Contains('multiply') to return false")
	}
}
