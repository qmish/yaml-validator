package reporter

import (
	"encoding/json"
	"testing"

	"yaml-validator/pkg"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateGitLabCodeQualityReport(t *testing.T) {
	results := []FileResult{
		{
			File:   "config.yaml",
			Errors: []pkg.Error{{Type: "SyntaxError", Message: "invalid YAML", Line: 3}},
		},
	}
	report, err := GenerateGitLabCodeQualityReport(results)
	require.NoError(t, err)

	var entries []map[string]interface{}
	require.NoError(t, json.Unmarshal(report, &entries))
	require.Len(t, entries, 1)
	assert.Equal(t, "yaml-validator", entries[0]["engine_name"])
	assert.Equal(t, "SyntaxError", entries[0]["check_name"])
	assert.Equal(t, "invalid YAML", entries[0]["description"])
}
