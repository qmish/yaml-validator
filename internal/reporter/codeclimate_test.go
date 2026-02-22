package reporter

import (
	"strings"
	"testing"

	"yaml-validator/pkg"
)

func TestGenerateCodeClimateReport(t *testing.T) {
	results := []FileResult{
		{
			File: "/code/config.yaml",
			Errors: []pkg.Error{
				{Type: "TrailingSpaces", Message: "trailing spaces", Line: 5},
				{Type: "SyntaxError", Message: "syntax error", Line: 10},
			},
		},
	}
	out, err := GenerateCodeClimateReport(results, "/code")
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(lines))
	}
	for i, line := range lines {
		if line == "" {
			continue
		}
		if !strings.Contains(line, `"type":"issue"`) {
			t.Errorf("Line %d: expected type issue, got %s", i+1, line)
		}
		if !strings.Contains(line, "config.yaml") {
			t.Errorf("Line %d: expected path, got %s", i+1, line)
		}
	}
}

func TestGenerateCodeClimateReport_Empty(t *testing.T) {
	results := []FileResult{{File: "a.yaml", Errors: nil}}
	out, err := GenerateCodeClimateReport(results, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 0 {
		t.Errorf("Expected empty output, got %q", string(out))
	}
}
