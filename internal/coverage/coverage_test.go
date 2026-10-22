package coverage

import (
	"strings"
	"testing"

	"github.com/forgekit/docforge/internal/types"
)

func TestCalculateFullCoverage(t *testing.T) {
	project := types.Project{
		Modules: []types.Module{
			{
				Name: "utils.ts",
				Documents: []types.Document{
					{Name: "add", IsDocumented: true, Visibility: types.VisibilityPublic},
					{Name: "sub", IsDocumented: true, Visibility: types.VisibilityPublic},
				},
			},
		},
	}
	report := Calculate(project)
	if report.Total != 2 {
		t.Errorf("expected total 2, got %d", report.Total)
	}
	if report.Documented != 2 {
		t.Errorf("expected documented 2, got %d", report.Documented)
	}
	if report.Percent != 100.0 {
		t.Errorf("expected 100%%, got %.1f%%", report.Percent)
	}
	if !IsFullyDocumented(report) {
		t.Error("expected fully documented")
	}
	if MissingCount(report) != 0 {
		t.Errorf("expected 0 missing, got %d", MissingCount(report))
	}
}

func TestCalculatePartialCoverage(t *testing.T) {
	project := types.Project{
		Modules: []types.Module{
			{
				Name: "utils.ts",
				Documents: []types.Document{
					{Name: "add", IsDocumented: true, Visibility: types.VisibilityPublic},
					{Name: "sub", IsDocumented: false, Visibility: types.VisibilityPublic},
					{Name: "mul", IsDocumented: false, Visibility: types.VisibilityPublic},
				},
			},
		},
	}
	report := Calculate(project)
	if report.Total != 3 {
		t.Errorf("expected total 3, got %d", report.Total)
	}
	if report.Documented != 1 {
		t.Errorf("expected documented 1, got %d", report.Documented)
	}
	if MissingCount(report) != 2 {
		t.Errorf("expected 2 missing, got %d", MissingCount(report))
	}
	if IsFullyDocumented(report) {
		t.Error("expected not fully documented")
	}
}

func TestCalculateSkipsPrivate(t *testing.T) {
	project := types.Project{
		Modules: []types.Module{
			{
				Name: "utils.ts",
				Documents: []types.Document{
					{Name: "add", IsDocumented: true, Visibility: types.VisibilityPublic},
					{Name: "_internal", IsDocumented: false, Visibility: types.VisibilityPrivate},
				},
			},
		},
	}
	report := Calculate(project)
	if report.Total != 1 {
		t.Errorf("expected total 1 (private excluded), got %d", report.Total)
	}
}

func TestCalculateEmptyProject(t *testing.T) {
	report := Calculate(types.Project{})
	if report.Total != 0 {
		t.Errorf("expected total 0, got %d", report.Total)
	}
	if report.Percent != 0 {
		t.Errorf("expected 0%%, got %.1f%%", report.Percent)
	}
}

func TestFormatReport(t *testing.T) {
	report := types.CoverageReport{
		Total:      10,
		Documented: 8,
		Percent:    80.0,
		Missing: []types.MissingDocument{
			{Module: "utils.ts", Name: "undocumented", Kind: types.KindFunction, Line: 5},
		},
	}
	result := FormatReport(report)
	if !strings.Contains(result, "8/10") {
		t.Error("expected report to contain '8/10'")
	}
	if !strings.Contains(result, "80.0%") {
		t.Error("expected report to contain '80.0%'")
	}
	if !strings.Contains(result, "undocumented") {
		t.Error("expected report to contain 'undocumented'")
	}
}

func TestCoverageForModule(t *testing.T) {
	module := types.Module{
		Name: "test.ts",
		Documents: []types.Document{
			{Name: "a", IsDocumented: true, Visibility: types.VisibilityPublic},
			{Name: "b", IsDocumented: false, Visibility: types.VisibilityPublic},
			{Name: "_c", IsDocumented: false, Visibility: types.VisibilityPrivate},
		},
	}
	total, documented := CoverageForModule(module)
	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
	if documented != 1 {
		t.Errorf("expected documented 1, got %d", documented)
	}
}

func TestCalculateWithDescriptionButNoAnnotation(t *testing.T) {
	project := types.Project{
		Modules: []types.Module{
			{
				Name: "utils.ts",
				Documents: []types.Document{
					{Name: "add", IsDocumented: false, Description: "Adds two numbers", Visibility: types.VisibilityPublic},
				},
			},
		},
	}
	report := Calculate(project)
	// Description without annotations should still count as documented
	if report.Documented != 1 {
		t.Errorf("expected documented 1 (has description), got %d", report.Documented)
	}
}
