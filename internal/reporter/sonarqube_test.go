package reporter

import (
	"encoding/json"
	"strings"
	"testing"

	"yaml-validator/pkg"
)

func TestGenerateSonarQubeGenericReport(t *testing.T) {
	results := []FileResult{
		{
			File: "/project/config.yaml",
			Errors: []pkg.Error{
				{Type: "TrailingSpaces", Message: "trailing spaces", Line: 5},
				{Type: "SyntaxError", Message: "syntax error", Line: 10},
			},
		},
	}
	out, err := GenerateSonarQubeGenericReport(results, "/project")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"issues"`) {
		t.Errorf("Expected issues in output")
	}
	if !strings.Contains(string(out), "yaml-validator") {
		t.Errorf("Expected engineId in output")
	}
	var report struct {
		Issues []interface{} `json:"issues"`
		Rules  []interface{} `json:"rules"`
	}
	if err := json.Unmarshal(out, &report); err != nil {
		t.Fatal(err)
	}
	if len(report.Issues) != 2 {
		t.Errorf("Expected 2 issues, got %d", len(report.Issues))
	}
	if len(report.Rules) < 1 {
		t.Errorf("Expected at least 1 rule, got %d", len(report.Rules))
	}
}
