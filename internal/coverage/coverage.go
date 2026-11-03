// Package coverage calculates documentation coverage — the percentage
// of public symbols that have documentation annotations.
package coverage

import (
	"fmt"

	"github.com/crafstack/docforge/internal/types"
)

// Calculate computes the documentation coverage report for a project.
//
// Only public symbols are counted (private symbols, identified by a
// leading underscore, are excluded).  A symbol is "documented" if it
// has at least one annotation (@doc, @param, @returns, @example) or a
// non-empty description.
func Calculate(project types.Project) types.CoverageReport {
	report := types.CoverageReport{}

	for _, module := range project.Modules {
		for _, doc := range module.Documents {
			// Skip private symbols
			if doc.Visibility == types.VisibilityPrivate {
				continue
			}

			report.Total++

			if doc.IsDocumented || doc.Description != "" {
				report.Documented++
			} else {
				report.Missing = append(report.Missing, types.MissingDocument{
					Module: module.Name,
					Name:   doc.Name,
					Kind:   doc.Kind,
					Line:   doc.Line,
				})
			}
		}
	}

	if report.Total > 0 {
		report.Percent = float64(report.Documented) / float64(report.Total) * 100
	}

	return report
}

// FormatReport returns a human-readable coverage report as a string.
func FormatReport(report types.CoverageReport) string {
	result := fmt.Sprintf(
		"Documentation Coverage: %d/%d (%.1f%%)\n",
		report.Documented, report.Total, report.Percent,
	)

	if len(report.Missing) > 0 {
		result += "\nMissing documentation:\n"
		for _, m := range report.Missing {
			result += fmt.Sprintf("  - %s (%s) in %s:%d\n", m.Name, m.Kind, m.Module, m.Line)
		}
	}

	return result
}

// IsFullyDocumented returns true if all public symbols have documentation.
func IsFullyDocumented(report types.CoverageReport) bool {
	return report.Total > 0 && report.Documented == report.Total
}

// MissingCount returns the number of undocumented public symbols.
func MissingCount(report types.CoverageReport) int {
	return len(report.Missing)
}

// CoverageForModule calculates coverage for a single module.
func CoverageForModule(module types.Module) (total int, documented int) {
	for _, doc := range module.Documents {
		if doc.Visibility == types.VisibilityPrivate {
			continue
		}
		total++
		if doc.IsDocumented || doc.Description != "" {
			documented++
		}
	}
	return
}
