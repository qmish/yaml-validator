package validator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"yaml-validator/internal/parser"
)

func TestCheckIntegrity_MissingFields(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "k8s.yaml")
	content := `
apiVersion: v1
kind: Pod
metadata:
  labels: {}
`
	require.NoError(t, os.WriteFile(f, []byte(content), 0644))

	node, err := parser.ParseFile(f)
	require.NoError(t, err)

	errors := CheckIntegrity(node, []string{"apiVersion", "kind", "metadata.name"})
	require.Len(t, errors, 1)
	assert.Equal(t, "MissingRequiredField", errors[0].Type)
	assert.Equal(t, "metadata.name", errors[0].Path)
}

func TestCheckIntegrity_AllPresent(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "valid.yaml")
	content := `
apiVersion: v1
kind: Pod
metadata:
  name: test
`
	require.NoError(t, os.WriteFile(f, []byte(content), 0644))

	node, err := parser.ParseFile(f)
	require.NoError(t, err)

	errors := CheckIntegrity(node, []string{"apiVersion", "kind", "metadata.name"})
	assert.Empty(t, errors)
}
