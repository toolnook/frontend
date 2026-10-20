package types

import "testing"

func TestDocumentHasParam(t *testing.T) {
	doc := Document{
		Name:       "add",
		Parameters: []Parameter{{Name: "a"}, {Name: "b"}},
	}
	if !doc.HasParam("a") {
		t.Error("expected HasParam('a') to return true")
	}
	if doc.HasParam("c") {
		t.Error("expected HasParam('c') to return false")
	}
}

func TestDocumentAddTag(t *testing.T) {
	doc := Document{Name: "test"}
	doc.AddTag("api")
	doc.AddTag("api") // duplicate
	doc.AddTag("public")
	if len(doc.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(doc.Tags))
	}
}

func TestProjectGetModule(t *testing.T) {
	p := Project{
		Modules: []Module{{Name: "users.ts"}, {Name: "orders.ts"}},
	}
	m := p.GetModule("users.ts")
	if m == nil {
		t.Fatal("expected to find module users.ts")
	}
	if m.Name != "users.ts" {
		t.Errorf("expected name 'users.ts', got '%s'", m.Name)
	}
	if p.GetModule("nonexistent") != nil {
		t.Error("expected nil for nonexistent module")
	}
}

func TestProjectAllDocuments(t *testing.T) {
	p := Project{
		Modules: []Module{
			{Name: "a.ts", Documents: []Document{{Name: "f1"}, {Name: "f2"}}},
			{Name: "b.ts", Documents: []Document{{Name: "f3"}}},
		},
	}
	docs := p.AllDocuments()
	if len(docs) != 3 {
		t.Errorf("expected 3 documents, got %d", len(docs))
	}
}

func TestProjectPublicDocuments(t *testing.T) {
	p := Project{
		Modules: []Module{
			{Name: "a.ts", Documents: []Document{
				{Name: "pub", Visibility: VisibilityPublic},
				{Name: "_priv", Visibility: VisibilityPrivate},
			}},
		},
	}
	docs := p.PublicDocuments()
	if len(docs) != 1 {
		t.Errorf("expected 1 public document, got %d", len(docs))
	}
	if docs[0].Name != "pub" {
		t.Errorf("expected 'pub', got '%s'", docs[0].Name)
	}
}

func TestDocumentString(t *testing.T) {
	doc := Document{Name: "sum", Kind: KindFunction, Module: "math.ts", Line: 10}
	s := doc.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}
