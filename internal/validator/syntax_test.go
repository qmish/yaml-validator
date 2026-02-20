package validator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckSyntax_ValidFile(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "valid.yaml")
	require.NoError(t, os.WriteFile(f, []byte("apiVersion: v1\nkind: Pod\n"), 0644))

	errors := CheckSyntax(f)
	assert.Empty(t, errors)
}

func TestCheckSyntax_EmptyFile(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "empty.yaml")
	require.NoError(t, os.WriteFile(f, []byte{}, 0644))

	errors := CheckSyntax(f)
	require.Len(t, errors, 1)
	assert.Equal(t, "EmptyFile", errors[0].Type)
}

func TestCheckSyntax_InvalidYAML(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "invalid.yaml")
	require.NoError(t, os.WriteFile(f, []byte("key: [unclosed\n"), 0644))

	errors := CheckSyntax(f)
	require.Len(t, errors, 1)
	assert.Equal(t, "SyntaxError", errors[0].Type)
}
