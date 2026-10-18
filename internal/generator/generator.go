// Package generator builds the static HTML documentation site from
// parsed modules.  It produces an index page, one page per module,
// a CSS stylesheet, a JavaScript search index, and a dependency graph
// visualization.
package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/forgekit/docforge/internal/types"
)

// Generator creates the static documentation site.
type Generator struct {
	OutputDir string
}

// New creates a Generator that writes to the given output directory.
func New(outputDir string) *Generator {
	return &Generator{OutputDir: outputDir}
}

// Generate writes the full static site for the given project.
func (g *Generator) Generate(project types.Project) error {
	// Ensure output directory exists
	if err := os.MkdirAll(g.OutputDir, 0755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	// Generate CSS
	if err := g.writeFile("styles.css", CSS); err != nil {
		return err
	}

	// Generate JS
	if err := g.writeFile("search.js", SearchJS); err != nil {
		return err
	}
	if err := g.writeFile("graph.js", GraphJS); err != nil {
		return err
	}

	// Generate search index
	searchIndex := BuildSearchIndex(project)
	searchJSON, err := json.MarshalIndent(searchIndex, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal search index: %w", err)
	}
	if err := g.writeFile("search-index.json", string(searchJSON)); err != nil {
		return err
	}

	// Generate index page
	indexHTML := RenderIndex(project)
	if err := g.writeFile("index.html", indexHTML); err != nil {
		return err
	}

	// Generate module pages
	for _, module := range project.Modules {
		html := RenderModule(module)
		filename := fmt.Sprintf("%s.html", sanitizeFilename(module.Name))
		if err := g.writeFile(filename, html); err != nil {
			return err
		}
	}

	// Generate dependency graph data
	graphData := BuildGraphData(project)
	graphJSON, err := json.MarshalIndent(graphData, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal graph data: %w", err)
	}
	if err := g.writeFile("graph-data.json", string(graphJSON)); err != nil {
		return err
	}

	return nil
}

// writeFile writes a file to the output directory.
func (g *Generator) writeFile(filename, content string) error {
	path := filepath.Join(g.OutputDir, filename)
	return os.WriteFile(path, []byte(content), 0644)
}

// sanitizeFilename converts a module name to a safe filename.
func sanitizeFilename(name string) string {
	name = strings.ReplaceAll(name, ".", "-")
	name = strings.ReplaceAll(name, "/", "-")
	name = strings.ReplaceAll(name, " ", "-")
	return strings.ToLower(name)
}

// BuildSearchIndex creates a flat list of searchable entries from the project.
func BuildSearchIndex(project types.Project) []SearchEntry {
	var entries []SearchEntry
	for _, module := range project.Modules {
		for _, doc := range module.Documents {
			entries = append(entries, SearchEntry{
				Name:        doc.Name,
				Module:      module.Name,
				Kind:        string(doc.Kind),
				Description: doc.Description,
				URL:         fmt.Sprintf("%s.html#%s", sanitizeFilename(module.Name), doc.Name),
			})
		}
	}
	return entries
}

// SearchEntry is a single entry in the search index.
type SearchEntry struct {
	Name        string `json:"name"`
	Module      string `json:"module"`
	Kind        string `json:"kind"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

// BuildGraphData constructs the dependency graph from module imports.
func BuildGraphData(project types.Project) types.GraphData {
	data := types.GraphData{}
	nodeSet := make(map[string]bool)
	moduleIndex := make(map[string]bool)

	for _, m := range project.Modules {
		moduleIndex[m.Name] = true
	}

	for _, m := range project.Modules {
		if !nodeSet[m.Name] {
			data.Nodes = append(data.Nodes, types.GraphNode{
				ID:     m.Name,
				Label:  m.Name,
				Module: m.Name,
				Degree: 0,
			})
			nodeSet[m.Name] = true
		}

		for _, imp := range m.Imports {
			// Resolve import to a module name
			target := resolveImport(imp, moduleIndex)
			if target == "" {
				continue
			}

			if !nodeSet[target] {
				data.Nodes = append(data.Nodes, types.GraphNode{
					ID:     target,
					Label:  target,
					Module: target,
					Degree: 0,
				})
				nodeSet[target] = true
			}

			data.Edges = append(data.Edges, types.GraphEdge{
				From: m.Name,
				To:   target,
			})

			// Increment degree
			for i := range data.Nodes {
				if data.Nodes[i].ID == m.Name || data.Nodes[i].ID == target {
					data.Nodes[i].Degree++
				}
			}
		}
	}

	return data
}

// resolveImport tries to match an import path to a module name.
func resolveImport(imp string, moduleIndex map[string]bool) string {
	parts := strings.Split(imp, "/")
	for i := len(parts) - 1; i >= 0; i-- {
		candidate := parts[i]
		if candidate == "" || candidate == "." || candidate == ".." {
			continue
		}
		// Try matching directly
		if moduleIndex[candidate] {
			return candidate
		}
		// Try matching with common extensions
		for _, ext := range []string{".ts", ".tsx", ".py", ".js"} {
			if moduleIndex[candidate+ext] {
				return candidate + ext
			}
		}
	}
	return ""
}

// RenderIndex generates the HTML for the index/landing page.
func RenderIndex(project types.Project) string {
	totalDocs := 0
	documented := 0
	for _, m := range project.Modules {
		for _, d := range m.Documents {
			totalDocs++
			if d.IsDocumented {
				documented++
			}
		}
	}

	var moduleList strings.Builder
	for _, m := range project.Modules {
		docCount := len(m.Documents)
		moduleList.WriteString(fmt.Sprintf(
			`<li><a href="%s.html">%s</a> <span class="badge">%d docs</span></li>`,
			sanitizeFilename(m.Name), m.Name, docCount,
		))
	}

	coverage := 0.0
	if totalDocs > 0 {
		coverage = float64(documented) / float64(totalDocs) * 100
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>%s — Documentation</title>
<link rel="stylesheet" href="styles.css">
</head>
<body>
<nav class="sidebar">
  <h1 class="logo">%s</h1>
  <input type="text" id="search" placeholder="Search docs..." class="search-input">
  <div id="search-results"></div>
  <ul class="module-list">%s</ul>
</nav>
<main class="content">
  <h1>%s</h1>
  <p class="version">v%s</p>
  <div class="coverage">
    <h3>Documentation Coverage</h3>
    <div class="progress-bar"><div class="progress-fill" style="width: %.1f%%"></div></div>
    <p>%d / %d symbols documented (%.1f%%)</p>
  </div>
  <h2>Modules</h2>
  <ul class="module-grid">%s</ul>
  <h2>Dependency Graph</h2>
  <div id="graph-container"></div>
</main>
<script src="search.js"></script>
<script src="graph.js"></script>
</body>
</html>`, project.Name, project.Name, moduleList.String(), project.Name, project.Version, coverage, documented, totalDocs, coverage, moduleList.String())
}

// RenderModule generates the HTML for a single module page.
func RenderModule(module types.Module) string {
	var docList strings.Builder
	for _, doc := range module.Documents {
		var params strings.Builder
		for _, p := range doc.Parameters {
			params.WriteString(fmt.Sprintf(`<li><code>%s</code> <span class="type">%s</span> — %s</li>`, p.Name, p.Type, p.Description))
		}

		var examples strings.Builder
		for _, ex := range doc.Examples {
			examples.WriteString(fmt.Sprintf(`<div class="example"><p>%s</p><pre><code>%s</code></pre></div>`, ex.Description, ex.Code))
		}

		documentedClass := ""
		if !doc.IsDocumented {
			documentedClass = " missing-docs"
		}

		docList.WriteString(fmt.Sprintf(
			`<section class="doc%s" id="%s">
  <h3>%s <span class="kind %s">%s</span>%s</h3>
  <p class="signature"><code>%s</code></p>
  <p>%s</p>
  <ul class="params">%s</ul>
  <p class="returns"><strong>Returns:</strong> <span class="type">%s</span> %s</p>
  <div class="examples">%s</div>
</section>`,
			documentedClass, doc.Name, doc.Name, doc.Kind, doc.Kind,
			coverageBadge(doc.IsDocumented),
			doc.Signature, doc.Description,
			params.String(), doc.Returns, doc.ReturnsDesc,
			examples.String(),
		))
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>%s — Documentation</title>
<link rel="stylesheet" href="styles.css">
</head>
<body>
<nav class="sidebar">
  <h1 class="logo"><a href="index.html">%s</a></h1>
  <input type="text" id="search" placeholder="Search docs..." class="search-input">
  <div id="search-results"></div>
  <a href="index.html" class="back-link">← Back to index</a>
</nav>
<main class="content">
  <h1>%s</h1>
  <p class="language">%s</p>
  %s
</main>
<script src="search.js"></script>
</body>
</html>`, module.Name, module.Name, module.Name, module.Language, docList.String())
}

// coverageBadge returns an HTML badge indicating documentation status.
func coverageBadge(documented bool) string {
	if documented {
		return ` <span class="badge documented">documented</span>`
	}
	return ` <span class="badge missing">missing docs</span>`
}
