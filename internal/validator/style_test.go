package validator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"yaml-validator/internal/config"
)

func TestCheckStyle_RequireDocumentStart(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "doc.yaml")
	require.NoError(t, os.WriteFile(f, []byte("key: value\n"), 0644))

	opts := config.StyleOptions{RequireDocumentStart: true}
	errors := CheckStyle(f, opts)
	require.Len(t, errors, 1)
	assert.Equal(t, "DocumentStart", errors[0].Type)
}

func TestCheckStyle_RequireDocumentStart_Ok(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "doc.yaml")
	require.NoError(t, os.WriteFile(f, []byte("---\nkey: value\n"), 0644))

	opts := config.StyleOptions{RequireDocumentStart: true}
	errors := CheckStyle(f, opts)
	assert.Empty(t, errors)
}

func TestCheckStyle_TrailingSpaces(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "doc.yaml")
	require.NoError(t, os.WriteFile(f, []byte("key: value   \n"), 0644))

	opts := config.StyleOptions{ForbidTrailingSpaces: true}
	errors := CheckStyle(f, opts)
	require.Len(t, errors, 1)
	assert.Equal(t, "TrailingSpaces", errors[0].Type)
}

func TestCheckStyle_RequireNewlineAtEof(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "doc.yaml")
	require.NoError(t, os.WriteFile(f, []byte("key: value"), 0644))

	opts := config.StyleOptions{RequireNewlineAtEof: true}
	errors := CheckStyle(f, opts)
	require.Len(t, errors, 1)
	assert.Equal(t, "NewlineAtEof", errors[0].Type)
}

func TestCheckStyle_Disabled(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "doc.yaml")
	require.NoError(t, os.WriteFile(f, []byte("key: value"), 0644))

	opts := config.StyleOptions{}
	errors := CheckStyle(f, opts)
	assert.Empty(t, errors)
}
