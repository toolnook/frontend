// Package parser reads annotated source files and extracts documentation
// from comment blocks using the @doc, @param, @returns, and @example
// convention.
//
// The parser is language-aware: TypeScript and Python are supported
// out of the box, with a registry pattern that allows adding more
// languages without modifying the core.
package parser

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/crafstack/docforge/internal/types"
)

// LanguageParser is the interface that each language-specific parser
// must implement.
type LanguageParser interface {
	// ParseFile reads a source file and returns a Module with all
	// documented and undocumented symbols.
	ParseFile(path string, content []byte) (types.Module, error)
	// Language returns the language name (e.g. "typescript").
	Language() string
	// Extensions returns the file extensions this parser handles.
	Extensions() []string
}

// Registry holds all registered language parsers.
type Registry struct {
	parsers []LanguageParser
	extMap  map[string]LanguageParser
}

// NewRegistry creates a registry with the built-in parsers
// (TypeScript and Python).
func NewRegistry() *Registry {
	r := &Registry{
		extMap: make(map[string]LanguageParser),
	}
	r.Register(&TypeScriptParser{})
	r.Register(&PythonParser{})
	return r
}

// Register adds a language parser to the registry.
func (r *Registry) Register(p LanguageParser) {
	r.parsers = append(r.parsers, p)
	for _, ext := range p.Extensions() {
		r.extMap[ext] = p
	}
}

// ParserForFile returns the language parser for the given file path,
// based on its extension. Returns nil if no parser matches.
func (r *Registry) ParserForFile(path string) LanguageParser {
	ext := strings.ToLower(filepath.Ext(path))
	return r.extMap[ext]
}

// ParseFile parses a single file using the appropriate language parser.
func (r *Registry) ParseFile(path string) (types.Module, error) {
	p := r.ParserForFile(path)
	if p == nil {
		return types.Module{}, ErrUnsupportedLanguage
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return types.Module{}, err
	}

	module, err := p.ParseFile(path, content)
	if err != nil {
		return types.Module{}, err
	}

	return module, nil
}

// ParseDir walks a directory recursively and parses all supported files.
func (r *Registry) ParseDir(dir string) ([]types.Module, error) {
	var modules []types.Module

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// Skip hidden directories and common non-source dirs
			name := info.Name()
			if name != dir && strings.HasPrefix(name, ".") {
				return filepath.SkipDir
			}
			if name == "node_modules" || name == "vendor" || name == "dist" || name == "build" {
				return filepath.SkipDir
			}
			return nil
		}

		p := r.ParserForFile(path)
		if p == nil {
			return nil
		}

		module, err := r.ParseFile(path)
		if err != nil {
			return nil // skip unparseable files
		}
		modules = append(modules, module)
		return nil
	})

	return modules, err
}

// DetectLanguage returns the language name for a file path, or "".
func (r *Registry) DetectLanguage(path string) string {
	p := r.ParserForFile(path)
	if p == nil {
		return ""
	}
	return p.Language()
}

// SupportedExtensions returns all file extensions the registry can parse.
func (r *Registry) SupportedExtensions() []string {
	exts := make([]string, 0, len(r.extMap))
	for ext := range r.extMap {
		exts = append(exts, ext)
	}
	return exts
}
