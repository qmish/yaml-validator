package fixer

import (
	"os"
	"path/filepath"
	"testing"

	"yaml-validator/internal/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFixFile_TrailingSpaces(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "test.yaml")
	require.NoError(t, os.WriteFile(f, []byte("key: value   \nother: x\t\n"), 0644))

	cfg := &config.Config{Rules: config.ValidationRules{
		Style: config.StyleOptions{ForbidTrailingSpaces: true},
	}}
	res, err := FixFile(f, cfg)
	require.NoError(t, err)
	assert.True(t, res.Modified)
	assert.Contains(t, res.Applied, "TrailingSpaces")

	data, _ := os.ReadFile(f)
	assert.Equal(t, "key: value\nother: x\n", string(data))
}

func TestFixFile_NewlineAtEof(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "test.yaml")
	require.NoError(t, os.WriteFile(f, []byte("key: value"), 0644))

	cfg := &config.Config{Rules: config.ValidationRules{
		Style: config.StyleOptions{RequireNewlineAtEof: true},
	}}
	res, err := FixFile(f, cfg)
	require.NoError(t, err)
	assert.True(t, res.Modified)
	assert.Contains(t, res.Applied, "NewlineAtEof")

	data, _ := os.ReadFile(f)
	assert.True(t, len(data) > 0 && data[len(data)-1] == '\n')
}

func TestFixFile_ConsecutiveEmptyLines(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "test.yaml")
	require.NoError(t, os.WriteFile(f, []byte("a: 1\n\n\n\nb: 2\n"), 0644))

	cfg := &config.Config{Rules: config.ValidationRules{
		Style: config.StyleOptions{ForbidConsecutiveEmptyLines: true},
	}}
	res, err := FixFile(f, cfg)
	require.NoError(t, err)
	assert.True(t, res.Modified)
	assert.Contains(t, res.Applied, "ConsecutiveEmptyLines")

	data, _ := os.ReadFile(f)
	assert.Equal(t, "a: 1\n\nb: 2\n", string(data))
}

func TestFixFile_NoChange(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "test.yaml")
	content := "key: value\n"
	require.NoError(t, os.WriteFile(f, []byte(content), 0644))

	cfg := &config.Config{Rules: config.ValidationRules{
		Style: config.StyleOptions{
			ForbidTrailingSpaces:       true,
			RequireNewlineAtEof:        true,
			ForbidConsecutiveEmptyLines: true,
		},
	}}
	res, err := FixFile(f, cfg)
	require.NoError(t, err)
	assert.False(t, res.Modified)
	assert.Empty(t, res.Applied)

	data, _ := os.ReadFile(f)
	assert.Equal(t, content, string(data))
}
