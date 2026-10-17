// Package graph builds a dependency graph from module imports and
// generates data for visualisation as an interactive SVG.
package graph

import (
	"strings"

	"github.com/forgekit/docforge/internal/types"
)

// Builder constructs dependency graphs from project modules.
type Builder struct {
	// Nodes maps module names to their node data.
	Nodes map[string]*Node
	// Edges is the list of dependency edges.
	Edges []Edge
}

// Node is a node in the dependency graph.
type Node struct {
	ID     string `json:"id"`
	Label  string `json:"label"`
	Degree int    `json:"degree"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
}

// Edge represents a dependency from one module to another.
type Edge struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// NewBuilder creates an empty graph builder.
func NewBuilder() *Builder {
	return &Builder{
		Nodes: make(map[string]*Node),
	}
}

// Build constructs a dependency graph from the project's modules.
func Build(project types.Project) *Builder {
	b := NewBuilder()

	// Add all modules as nodes
	for _, m := range project.Modules {
		b.addNode(m.Name)
	}

	// Add edges based on imports
	for _, m := range project.Modules {
		for _, imp := range m.Imports {
			target := resolveImport(imp, project)
			if target != "" && target != m.Name {
				b.addNode(target)
				b.addEdge(m.Name, target)
			}
		}
	}

	// Layout nodes in a circular arrangement
	b.layout()

	return b
}

func (b *Builder) addNode(id string) {
	if _, exists := b.Nodes[id]; !exists {
		b.Nodes[id] = &Node{
			ID:    id,
			Label: id,
		}
	}
}

func (b *Builder) addEdge(from, to string) {
	b.Edges = append(b.Edges, Edge{From: from, To: to})
	b.Nodes[from].Degree++
	b.Nodes[to].Degree++
}

// layout positions nodes in a circle.
func (b *Builder) layout() {
	count := len(b.Nodes)
	if count == 0 {
		return
	}

	const cx, cy, radius = 400.0, 200.0, 150.0
	i := 0
	for _, node := range b.Nodes {
		angle := float64(i) / float64(count) * 2 * 3.14159265358979
		node.X = cx + radius*angle
		node.Y = cy + radius*angle
		i++
	}
}

// ToGraphData converts the builder to the types.GraphData format
// used by the generator.
func (b *Builder) ToGraphData() types.GraphData {
	data := types.GraphData{}
	for _, node := range b.Nodes {
		data.Nodes = append(data.Nodes, types.GraphNode{
			ID:     node.ID,
			Label:  node.Label,
			Degree: node.Degree,
		})
	}
	for _, e := range b.Edges {
		data.Edges = append(data.Edges, types.GraphEdge{
			From: e.From,
			To:   e.To,
		})
	}
	return data
}

// NodeCount returns the number of nodes in the graph.
func (b *Builder) NodeCount() int {
	return len(b.Nodes)
}

// EdgeCount returns the number of edges in the graph.
func (b *Builder) EdgeCount() int {
	return len(b.Edges)
}

// HasNode checks whether a node with the given ID exists.
func (b *Builder) HasNode(id string) bool {
	_, exists := b.Nodes[id]
	return exists
}

// HasEdge checks whether an edge from→to exists.
func (b *Builder) HasEdge(from, to string) bool {
	for _, e := range b.Edges {
		if e.From == from && e.To == to {
			return true
		}
	}
	return false
}

// GetNode returns the node with the given ID, or nil.
func (b *Builder) GetNode(id string) *Node {
	return b.Nodes[id]
}

// GetDependents returns modules that depend on the given module.
func (b *Builder) GetDependents(module string) []string {
	var deps []string
	for _, e := range b.Edges {
		if e.To == module {
			deps = append(deps, e.From)
		}
	}
	return deps
}

// GetDependencies returns modules that the given module depends on.
func (b *Builder) GetDependencies(module string) []string {
	var deps []string
	for _, e := range b.Edges {
		if e.From == module {
			deps = append(deps, e.To)
		}
	}
	return deps
}

// resolveImport tries to match an import path to a module name in the project.
func resolveImport(imp string, project types.Project) string {
	parts := strings.Split(imp, "/")
	for i := len(parts) - 1; i >= 0; i-- {
		candidate := parts[i]
		if candidate == "" || candidate == "." || candidate == ".." {
			continue
		}
		for _, m := range project.Modules {
			// Try exact match
			if m.Name == candidate {
				return m.Name
			}
			// Try stripping extension from module name
			base := m.Name
			for _, ext := range []string{".ts", ".tsx", ".py", ".js"} {
				if strings.HasSuffix(base, ext) {
					base = strings.TrimSuffix(base, ext)
				}
			}
			if base == candidate {
				return m.Name
			}
		}
	}
	return ""
}
