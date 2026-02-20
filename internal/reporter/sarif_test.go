package reporter

import (
	"encoding/json"
	"testing"

	"yaml-validator/pkg"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateSARIFReport(t *testing.T) {
	results := []FileResult{
		{File: "/path/to/file.yaml", Errors: []pkg.Error{
			{Type: "LineTooLong", Message: "line too long", Line: 10},
		}},
	}
	report, err := GenerateSARIFReport("1.5.0", results)
	require.NoError(t, err)

	var log sarifLog
	require.NoError(t, json.Unmarshal(report, &log))
	assert.Equal(t, "2.1.0", log.Version)
	require.Len(t, log.Runs, 1)
	run := log.Runs[0]
	assert.Equal(t, "yaml-validator", run.Tool.Driver.Name)
	assert.Equal(t, "1.5.0", run.Tool.Driver.Version)
	require.Len(t, run.Artifacts, 1)
	assert.Contains(t, run.Artifacts[0].Location.URI, "file.yaml")
	require.Len(t, run.Results, 1)
	assert.Equal(t, "LineTooLong", run.Results[0].RuleID)
	require.NotEmpty(t, run.Results[0].Locations)
	assert.Equal(t, 10, run.Results[0].Locations[0].PhysicalLocation.Region.StartLine)
}

func TestGenerateSARIFReport_Empty(t *testing.T) {
	results := []FileResult{{File: "a.yaml", Errors: nil}}
	report, err := GenerateSARIFReport("1.0.0", results)
	require.NoError(t, err)
	var log sarifLog
	require.NoError(t, json.Unmarshal(report, &log))
	require.Len(t, log.Runs[0].Results, 0)
}
