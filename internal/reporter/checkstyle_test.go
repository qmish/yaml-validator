package reporter

import (
	"strings"
	"testing"

	"yaml-validator/pkg"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateCheckstyleReport(t *testing.T) {
	results := []FileResult{
		{
			File:   "config.yaml",
			Errors: []pkg.Error{{Type: "SyntaxError", Message: "invalid YAML", Line: 3, Column: 5}},
		},
		{
			File:   "deploy.yaml",
			Errors: []pkg.Error{},
		},
	}
	report, err := GenerateCheckstyleReport(results)
	require.NoError(t, err)

	s := string(report)
	assert.Contains(t, s, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>")
	assert.Contains(t, s, "<checkstyle version=\"yaml-validator\">")
	assert.Contains(t, s, "<file name=\"config.yaml\">")
	assert.Contains(t, s, `<error line="3" column="5" severity="error" message="invalid YAML" source="SyntaxError"/>`)
	assert.Contains(t, s, "<file name=\"deploy.yaml\">")
	assert.Contains(t, s, "</checkstyle>")
}

func TestGenerateCheckstyleReport_EmptyResults(t *testing.T) {
	report, err := GenerateCheckstyleReport(nil)
	require.NoError(t, err)
	s := string(report)
	assert.True(t, strings.HasPrefix(s, `<?xml version="1.0" encoding="UTF-8"?>`))
	assert.Contains(t, s, "</checkstyle>")
}
