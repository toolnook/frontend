// Package types defines the core data structures used across DocForge.
//
// A Document represents a single documented symbol (function, method, class,
// or type) parsed from source code. A Module groups all documents found in
// a single source file. A Project is the top-level container for all modules.
package types

import "fmt"

// Visibility of a documented symbol.
type Visibility string

const (
	VisibilityPublic  Visibility = "public"
	VisibilityPrivate  Visibility = "private"
	VisibilityUnknown  Visibility = "unknown"
)

// Kind of a documented symbol.
type SymbolKind string

const (
	KindFunction SymbolKind = "function"
	KindMethod   SymbolKind = "method"
	KindClass    SymbolKind = "class"
	KindType     SymbolKind = "type"
	KindVariable SymbolKind = "variable"
)

// Parameter describes a function or method parameter.
type Parameter struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

// Example is a usage example extracted from an @example annotation.
type Example struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

// Document represents a single documented symbol in source code.
type Document struct {
	Name        string       `json:"name"`
	Kind        SymbolKind   `json:"kind"`
	Visibility  Visibility   `json:"visibility"`
	Module      string       `json:"module"`
	Signature   string       `json:"signature"`
	Description string       `json:"description"`
	Parameters  []Parameter  `json:"parameters"`
	Returns     string       `json:"returns"`
	ReturnsDesc string       `json:"returnsDescription"`
	Examples    []Example    `json:"examples"`
	Line        int           `json:"line"`
	IsDocumented bool         `json:"isDocumented"`
	Tags        []string      `json:"tags"`
}

// Module represents a single source file and all the documents it contains.
type Module struct {
	Name      string     `json:"name"`
	Path      string     `json:"path"`
	Language  string     `json:"language"` // "typescript", "python", etc.
	Documents []Document `json:"documents"`
	Imports   []string   `json:"imports"` // module paths imported by this file
}

// Project is the top-level container for all parsed modules.
type Project struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Modules     []Module `json:"modules"`
	Description string   `json:"description"`
}

// CoverageReport describes documentation coverage for a project.
type CoverageReport struct {
	Total     int               `json:"total"`
	Documented int              `json:"documented"`
	Missing   []MissingDocument `json:"missing"`
	Percent   float64           `json:"percent"`
}

// MissingDocument is a public symbol that lacks documentation.
type MissingDocument struct {
	Module string `json:"module"`
	Name   string `json:"name"`
	Kind   SymbolKind `json:"kind"`
	Line   int    `json:"line"`
}

// GraphEdge represents a dependency relationship between modules.
type GraphEdge struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// GraphData is the dependency graph for the project.
type GraphData struct {
	Nodes []GraphNode `json:"nodes"`
	Edges []GraphEdge `json:"edges"`
}

// GraphNode is a node in the dependency graph.
type GraphNode struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	Module   string `json:"module"`
	Degree   int    `json:"degree"`
}

// String returns a human-readable representation of a Document.
func (d Document) String() string {
	return fmt.Sprintf("%s (%s) in %s:%d", d.Name, d.Kind, d.Module, d.Line)
}

// HasParam checks if a parameter with the given name exists.
func (d Document) HasParam(name string) bool {
	for _, p := range d.Parameters {
		if p.Name == name {
			return true
		}
	}
	return false
}

// AddTag adds a tag if it doesn't already exist.
func (d *Document) AddTag(tag string) {
	for _, t := range d.Tags {
		if t == tag {
			return
		}
	}
	d.Tags = append(d.Tags, tag)
}

// GetModule returns the module with the given name, or nil.
func (p *Project) GetModule(name string) *Module {
	for i := range p.Modules {
		if p.Modules[i].Name == name {
			return &p.Modules[i]
		}
	}
	return nil
}

// AllDocuments returns every document across all modules.
func (p *Project) AllDocuments() []Document {
	var docs []Document
	for _, m := range p.Modules {
		docs = append(docs, m.Documents...)
	}
	return docs
}

// PublicDocuments returns only documents marked as public.
func (p *Project) PublicDocuments() []Document {
	var docs []Document
	for _, d := range p.AllDocuments() {
		if d.Visibility == VisibilityPublic {
			docs = append(docs, d)
		}
	}
	return docs
}
