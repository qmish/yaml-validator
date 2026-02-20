package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "config.yaml")
	content := `
rules:
  check_syntax: true
  check_duplicates: false
  max_line_length: 100
`
	require.NoError(t, os.WriteFile(f, []byte(content), 0644))

	cfg, err := LoadConfig(f)
	require.NoError(t, err)
	assert.True(t, cfg.Rules.CheckSyntax)
	assert.False(t, cfg.Rules.CheckDuplicates)
	assert.Equal(t, 100, cfg.Rules.MaxLineLength)
}

func TestLoadConfig_NotFound(t *testing.T) {
	_, err := LoadConfig("/nonexistent/config.yaml")
	assert.Error(t, err)
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	assert.True(t, cfg.Rules.CheckSyntax)
	assert.True(t, cfg.Rules.CheckDuplicates)
	assert.NotEmpty(t, cfg.Rules.RequiredFields)
	assert.Contains(t, cfg.Rules.RequiredFields, "apiVersion")
}
