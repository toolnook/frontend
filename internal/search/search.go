// Package search builds a search index from parsed documentation and
// provides a simple full-text search over the index.
package search

import (
	"strings"

	"github.com/toolnook/docforge/internal/types"
)

// Index is a searchable index of documentation entries.
type Index struct {
	Entries []Entry `json:"entries"`
}

// Entry is a single searchable documentation entry.
type Entry struct {
	Name        string `json:"name"`
	Module      string `json:"module"`
	Kind        string `json:"kind"`
	Description string `json:"description"`
	Line        int    `json:"line"`
}

// BuildIndex creates a search index from a project's documents.
func BuildIndex(project types.Project) Index {
	var entries []Entry
	for _, module := range project.Modules {
		for _, doc := range module.Documents {
			entries = append(entries, Entry{
				Name:        doc.Name,
				Module:      module.Name,
				Kind:        string(doc.Kind),
				Description: doc.Description,
				Line:        doc.Line,
			})
		}
	}
	return Index{Entries: entries}
}

// Search returns entries matching the query string.  Matching is
// case-insensitive and checks name, description, and module fields.
func (idx Index) Search(query string) []Entry {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return nil
	}

	var results []Entry
	for _, entry := range idx.Entries {
		if strings.Contains(strings.ToLower(entry.Name), query) ||
			strings.Contains(strings.ToLower(entry.Description), query) ||
			strings.Contains(strings.ToLower(entry.Module), query) {
			results = append(results, entry)
		}
	}
	return results
}

// SearchByKind returns entries of a specific kind.
func (idx Index) SearchByKind(kind string) []Entry {
	var results []Entry
	for _, entry := range idx.Entries {
		if entry.Kind == kind {
			results = append(results, entry)
		}
	}
	return results
}

// SearchByModule returns entries from a specific module.
func (idx Index) SearchByModule(module string) []Entry {
	var results []Entry
	for _, entry := range idx.Entries {
		if entry.Module == module {
			results = append(results, entry)
		}
	}
	return results
}

// EntryCount returns the total number of entries in the index.
func (idx Index) EntryCount() int {
	return len(idx.Entries)
}

// Contains checks whether an entry with the given name exists.
func (idx Index) Contains(name string) bool {
	for _, entry := range idx.Entries {
		if entry.Name == name {
			return true
		}
	}
	return false
}
