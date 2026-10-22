package graph

import (
	"testing"

	"github.com/forgekit/docforge/internal/types"
)

func TestBuildSimpleGraph(t *testing.T) {
	project := types.Project{
		Modules: []types.Module{
			{Name: "a.ts", Imports: []string{"./b", "./c"}},
			{Name: "b.ts", Imports: []string{"./c"}},
			{Name: "c.ts"},
		},
	}
	g := Build(project)
	if g.NodeCount() != 3 {
		t.Errorf("expected 3 nodes, got %d", g.NodeCount())
	}
	if g.EdgeCount() < 2 {
		t.Errorf("expected at least 2 edges, got %d", g.EdgeCount())
	}
}

func TestBuildNoDependencies(t *testing.T) {
	project := types.Project{
		Modules: []types.Module{
			{Name: "a.ts"},
			{Name: "b.ts"},
		},
	}
	g := Build(project)
	if g.NodeCount() != 2 {
		t.Errorf("expected 2 nodes, got %d", g.NodeCount())
	}
	if g.EdgeCount() != 0 {
		t.Errorf("expected 0 edges, got %d", g.EdgeCount())
	}
}

func TestHasNodeAndEdge(t *testing.T) {
	project := types.Project{
		Modules: []types.Module{
			{Name: "a.ts", Imports: []string{"./b"}},
			{Name: "b.ts"},
		},
	}
	g := Build(project)
	if !g.HasNode("a.ts") {
		t.Error("expected HasNode('a.ts') to return true")
	}
	if !g.HasNode("b.ts") {
		t.Error("expected HasNode('b.ts') to return true")
	}
	if !g.HasEdge("a.ts", "b.ts") {
		t.Error("expected HasEdge('a.ts', 'b.ts') to return true")
	}
	if g.HasEdge("b.ts", "a.ts") {
		t.Error("expected HasEdge('b.ts', 'a.ts') to return false")
	}
}

func TestGetDependents(t *testing.T) {
	project := types.Project{
		Modules: []types.Module{
			{Name: "a.ts", Imports: []string{"./c"}},
			{Name: "b.ts", Imports: []string{"./c"}},
			{Name: "c.ts"},
		},
	}
	g := Build(project)
	deps := g.GetDependents("c.ts")
	if len(deps) != 2 {
		t.Errorf("expected 2 dependents, got %d", len(deps))
	}
}

func TestGetDependencies(t *testing.T) {
	project := types.Project{
		Modules: []types.Module{
			{Name: "a.ts", Imports: []string{"./b", "./c"}},
			{Name: "b.ts"},
			{Name: "c.ts"},
		},
	}
	g := Build(project)
	deps := g.GetDependencies("a.ts")
	if len(deps) != 2 {
		t.Errorf("expected 2 dependencies, got %d", len(deps))
	}
}

func TestToGraphData(t *testing.T) {
	project := types.Project{
		Modules: []types.Module{
			{Name: "a.ts", Imports: []string{"./b"}},
			{Name: "b.ts"},
		},
	}
	g := Build(project)
	data := g.ToGraphData()
	if len(data.Nodes) != 2 {
		t.Errorf("expected 2 nodes, got %d", len(data.Nodes))
	}
	if len(data.Edges) != 1 {
		t.Errorf("expected 1 edge, got %d", len(data.Edges))
	}
}

func TestNodeDegree(t *testing.T) {
	project := types.Project{
		Modules: []types.Module{
			{Name: "a.ts", Imports: []string{"./c"}},
			{Name: "b.ts", Imports: []string{"./c"}},
			{Name: "c.ts"},
		},
	}
	g := Build(project)
	cNode := g.GetNode("c.ts")
	if cNode == nil {
		t.Fatal("expected node c.ts")
	}
	if cNode.Degree != 2 {
		t.Errorf("expected degree 2 for c.ts, got %d", cNode.Degree)
	}
}

func TestEmptyProject(t *testing.T) {
	g := Build(types.Project{})
	if g.NodeCount() != 0 {
		t.Errorf("expected 0 nodes, got %d", g.NodeCount())
	}
	if g.EdgeCount() != 0 {
		t.Errorf("expected 0 edges, got %d", g.EdgeCount())
	}
}
