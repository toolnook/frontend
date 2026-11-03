// DocForge is a documentation-as-code generator that turns annotated source
// code into a static, searchable API reference site.
//
// Usage:
//
//	docforge scan ./src            Parse source files and output a summary
//	docforge build -o ./docs       Generate static HTML documentation
//	docforge serve -d ./docs       Serve docs locally on port 8090
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/crafstack/docforge/internal/coverage"
	"github.com/crafstack/docforge/internal/generator"
	"github.com/crafstack/docforge/internal/parser"
	"github.com/crafstack/docforge/internal/server"
	"github.com/crafstack/docforge/internal/types"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "scan":
		scanCmd(args)
	case "build":
		buildCmd(args)
	case "serve":
		serveCmd(args)
	case "version":
		fmt.Println("docforge v0.1.0")
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

// scanCmd parses source files and outputs a JSON summary of what was found.
func scanCmd(args []string) {
	fs := flag.NewFlagSet("scan", flag.ExitOnError)
	dir := fs.String("dir", ".", "Source directory to scan")
	format := fs.String("format", "text", "Output format: text or json")
	fs.Parse(args)

	registry := parser.NewRegistry()
	modules, err := registry.ParseDir(*dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning directory: %v\n", err)
		os.Exit(1)
	}

	project := types.Project{
		Name:    filepath.Base(getAbsPath(*dir)),
		Version: "0.1.0",
		Modules: modules,
	}

	if *format == "json" {
		printJSON(project)
	} else {
		printTextSummary(project)
	}
}

// buildCmd generates the static documentation site.
func buildCmd(args []string) {
	fs := flag.NewFlagSet("build", flag.ExitOnError)
	srcDir := fs.String("dir", ".", "Source directory to scan")
	outDir := fs.String("o", "./docs", "Output directory for generated docs")
	version := fs.String("version", "0.1.0", "Version tag for the documentation")
	name := fs.String("name", "", "Project name (defaults to directory name)")
	fs.Parse(args)

	registry := parser.NewRegistry()
	modules, err := registry.ParseDir(*srcDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning source: %v\n", err)
		os.Exit(1)
	}

	projectName := *name
	if projectName == "" {
		projectName = filepath.Base(getAbsPath(*srcDir))
	}

	project := types.Project{
		Name:    projectName,
		Version: *version,
		Modules: modules,
	}

	gen := generator.New(*outDir)
	if err := gen.Generate(project); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating docs: %v\n", err)
		os.Exit(1)
	}

	// Print coverage report
	coverageReport := coverage.Calculate(project)
	fmt.Println(coverage.FormatReport(coverageReport))
	fmt.Printf("\nDocumentation generated in %s\n", *outDir)
}

// serveCmd starts a local HTTP server to serve the generated docs.
func serveCmd(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	dir := fs.String("d", "./docs", "Directory containing generated docs")
	port := fs.Int("port", 8090, "Port to serve on")
	fs.Parse(args)

	absDir, err := filepath.Abs(*dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving path: %v\n", err)
		os.Exit(1)
	}

	if _, err := os.Stat(absDir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Directory not found: %s\n", absDir)
		os.Exit(1)
	}

	srv := server.New(absDir, *port)
	fmt.Printf("Serving docs from %s on http://localhost:%d\n", absDir, *port)
	fmt.Println("Press Ctrl+C to stop.")
	if err := srv.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}

// printUsage outputs the CLI help text.
func printUsage() {
	fmt.Println(`DocForge — documentation-as-code generator

Usage:
  docforge <command> [options]

Commands:
  scan    Parse source files and output a summary
  build   Generate static HTML documentation site
  serve   Serve generated docs locally
  version Print version
  help    Show this help

Options for 'scan':
  -dir <path>     Source directory to scan (default: .)
  -format <fmt>   Output format: text or json (default: text)

Options for 'build':
  -dir <path>     Source directory to scan (default: .)
  -o <path>       Output directory (default: ./docs)
  -version <tag>  Version tag for the docs (default: 0.1.0)
  -name <name>    Project name (defaults to directory name)

Options for 'serve':
  -d <path>       Directory to serve (default: ./docs)
  -port <number>  Port number (default: 8090)`)
}

// printJSON outputs the project as JSON.
func printJSON(project types.Project) {
	data, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(data))
}

// printTextSummary outputs a human-readable summary of scanned modules.
func printTextSummary(project types.Project) {
	fmt.Printf("Project: %s (v%s)\n\n", project.Name, project.Version)

	totalDocs := 0
	documented := 0

	for _, module := range project.Modules {
		fmt.Printf("  %s [%s] — %d symbols", module.Name, module.Language, len(module.Documents))
		if len(module.Imports) > 0 {
			fmt.Printf(" (%d imports)", len(module.Imports))
		}
		fmt.Println()

		for _, doc := range module.Documents {
			status := "✓"
			if !doc.IsDocumented {
				status = "✗"
			}
			fmt.Printf("    %s %s (%s) line %d\n", status, doc.Name, doc.Kind, doc.Line)
			totalDocs++
			if doc.IsDocumented {
				documented++
			}
		}
	}

	fmt.Printf("\nTotal: %d symbols, %d documented\n", totalDocs, documented)
}

// getAbsPath returns the absolute path, handling errors.
func getAbsPath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	// Use the last component as the project name
	parts := strings.Split(strings.TrimRight(abs, string(filepath.Separator)), string(filepath.Separator))
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return path
}
