package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/toolnook/frontend/internal/types"
)

func TestGenerateCreatesFiles(t *testing.T) {
	dir := t.TempDir()
	gen := New(dir)

	project := types.Project{
		Name:    "test-project",
		Version: "1.0.0",
		Modules: []types.Module{
			{
				Name:     "utils.ts",
				Language: "typescript",
				Documents: []types.Document{
					{Name: "add", Kind: types.KindFunction, IsDocumented: true, Description: "Adds two numbers"},
					{Name: "subtract", Kind: types.KindFunction, IsDocumented: false},
				},
			},
		},
	}

	err := gen.Generate(project)
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}

	// Check that key files exist
	expectedFiles := []string{"index.html", "styles.css", "search.js", "graph.js", "search-index.json", "graph-data.json", "utils-ts.html"}
	for _, f := range expectedFiles {
		if _, err := os.Stat(filepath.Join(dir, f)); os.IsNotExist(err) {
			t.Errorf("expected file %s to exist", f)
		}
	}
}

func TestRenderIndexContainsProjectName(t *testing.T) {
	project := types.Project{
		Name:    "my-lib",
		Version: "2.0.0",
		Modules: []types.Module{
			{Name: "index.ts", Documents: []types.Document{{Name: "foo"}}},
		},
	}
	html := RenderIndex(project)
	if !contains(html, "my-lib") {
		t.Error("expected index HTML to contain project name 'my-lib'")
	}
	if !contains(html, "2.0.0") {
		t.Error("expected index HTML to contain version '2.0.0'")
	}
}

func TestRenderModuleContainsSymbolName(t *testing.T) {
	module := types.Module{
		Name:     "calc.ts",
		Language: "typescript",
		Documents: []types.Document{
			{Name: "calculate", Kind: types.KindFunction, IsDocumented: true, Description: "Calculates stuff"},
		},
	}
	html := RenderModule(module)
	if !contains(html, "calculate") {
		t.Error("expected module HTML to contain symbol name 'calculate'")
	}
}

func TestBuildSearchIndex(t *testing.T) {
	project := types.Project{
		Modules: []types.Module{
			{
				Name: "a.ts",
				Documents: []types.Document{
					{Name: "func1", Kind: types.KindFunction, Description: "Does thing 1"},
					{Name: "func2", Kind: types.KindFunction, Description: "Does thing 2"},
				},
			},
			{
				Name: "b.ts",
				Documents: []types.Document{
					{Name: "Class1", Kind: types.KindClass},
				},
			},
		},
	}
	entries := BuildSearchIndex(project)
	if len(entries) != 3 {
		t.Errorf("expected 3 search entries, got %d", len(entries))
	}
	if entries[0].Name != "func1" {
		t.Errorf("expected first entry name 'func1', got '%s'", entries[0].Name)
	}
}

func TestBuildGraphData(t *testing.T) {
	project := types.Project{
		Modules: []types.Module{
			{Name: "a.ts", Imports: []string{"./b", "./c"}},
			{Name: "b.ts", Imports: []string{"./c"}},
			{Name: "c.ts"},
		},
	}
	data := BuildGraphData(project)
	if len(data.Nodes) != 3 {
		t.Errorf("expected 3 nodes, got %d", len(data.Nodes))
	}
	// a.ts depends on b and c; b.ts depends on c
	if len(data.Edges) < 2 {
		t.Errorf("expected at least 2 edges, got %d", len(data.Edges))
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || indexOf(s, substr) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
