package reporter

import (
	"encoding/json"
	"testing"

	"yaml-validator/pkg"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateJSONReport(t *testing.T) {
	errors := []pkg.Error{
		{Type: "DuplicateKey", Message: "Duplicate key 'name'"},
	}
	report, err := GenerateJSONReport("config.yaml", errors)
	require.NoError(t, err)

	var r Report
	require.NoError(t, json.Unmarshal(report, &r))
	assert.Equal(t, "config.yaml", r.File)
	assert.False(t, r.Valid)
	assert.Len(t, r.Errors, 1)
	assert.Equal(t, "DuplicateKey", r.Errors[0].Type)
}

func TestGenerateJSONReport_Valid(t *testing.T) {
	report, err := GenerateJSONReport("valid.yaml", nil)
	require.NoError(t, err)

	var r Report
	require.NoError(t, json.Unmarshal(report, &r))
	assert.True(t, r.Valid)
	assert.Empty(t, r.Errors)
}

func TestGenerateJUnitReport(t *testing.T) {
	errors := []pkg.Error{
		{Type: "SyntaxError", Message: "invalid YAML"},
	}
	report, err := GenerateJUnitReport("test.yaml", errors)
	require.NoError(t, err)

	str := string(report)
	assert.Contains(t, str, `<testsuites`)
	assert.Contains(t, str, `failures="1"`)
	assert.Contains(t, str, `<failure`)
	assert.Contains(t, str, `SyntaxError`)
}

func TestEscapeXML(t *testing.T) {
	assert.Equal(t, "&lt;tag&gt;", escapeXML("<tag>"))
	assert.Equal(t, "&amp;", escapeXML("&"))
	assert.Equal(t, "&quot;quoted&quot;", escapeXML(`"quoted"`))
}

func TestGenerateCompactReport(t *testing.T) {
	errors := []pkg.Error{
		{Type: "SyntaxError", Message: "invalid YAML", Line: 3},
		{Type: "DuplicateKey", Message: "Duplicate key 'x'", Line: 5, Column: 2},
		{Type: "Other", Message: "no position", Line: 0},
	}
	out := GenerateCompactReport("app.yaml", errors)
	assert.Contains(t, out, "app.yaml:3: invalid YAML")
	assert.Contains(t, out, "app.yaml:5:2: Duplicate key 'x'")
	assert.Contains(t, out, "app.yaml: no position")
}

func TestGenerateCompactReport_Empty(t *testing.T) {
	out := GenerateCompactReport("valid.yaml", nil)
	assert.Empty(t, out)
}

func TestGenerateGitHubAnnotations(t *testing.T) {
	errors := []pkg.Error{
		{Type: "SyntaxError", Message: "invalid YAML", Line: 3},
		{Type: "Other", Message: "no position", Line: 0},
		{Type: "Percent", Message: "contains % sign", Line: 1},
	}
	out := GenerateGitHubAnnotations("app.yaml", errors)
	assert.Contains(t, out, "::error file=app.yaml,line=3::invalid YAML")
	assert.Contains(t, out, "::error file=app.yaml::no position")
	assert.Contains(t, out, "::error file=app.yaml,line=1::contains %25 sign")
}

func TestGenerateGitHubAnnotations_Empty(t *testing.T) {
	out := GenerateGitHubAnnotations("valid.yaml", nil)
	assert.Empty(t, out)
}

func TestGenerateSeverityReport(t *testing.T) {
	errors := []pkg.Error{
		{Type: "SyntaxError", Message: "invalid YAML", Line: 3},
		{Type: "DuplicateKey", Message: "Duplicate key 'x'", Line: 5, Column: 2},
		{Type: "Other", Message: "no position", Line: 0},
	}
	out := GenerateSeverityReport("app.yaml", errors)
	assert.Contains(t, out, "[ERROR] app.yaml:3: invalid YAML")
	assert.Contains(t, out, "[ERROR] app.yaml:5:2: Duplicate key 'x'")
	assert.Contains(t, out, "[ERROR] app.yaml: no position")
}

func TestGenerateSeverityReport_Empty(t *testing.T) {
	out := GenerateSeverityReport("valid.yaml", nil)
	assert.Empty(t, out)
}
