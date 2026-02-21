package validator

import (
	"os"
	"path/filepath"
	"testing"

	"yaml-validator/internal/parser"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckQuotedValues_Unquoted(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "test.yaml")
	err := os.WriteFile(f, []byte("name: foo\nvalue: bar\n"), 0644)
	require.NoError(t, err)

	node, err := parser.ParseFile(f)
	require.NoError(t, err)

	errors := CheckQuotedValues(node)
	assert.NotEmpty(t, errors)
	var hasQuotedValues bool
	for _, e := range errors {
		if e.Type == "QuotedValues" {
			hasQuotedValues = true
			break
		}
	}
	assert.True(t, hasQuotedValues)
}

func TestCheckQuotedValues_Empty(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "test.yaml")
	// count: 1 — число; port: 80 — число; строковых plain-значений нет
	err := os.WriteFile(f, []byte("count: 1\nport: 80\n"), 0644)
	require.NoError(t, err)

	node, err := parser.ParseFile(f)
	require.NoError(t, err)

	errors := CheckQuotedValues(node)
	assert.Empty(t, errors)
}
