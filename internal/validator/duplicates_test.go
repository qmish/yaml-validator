package validator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"yaml-validator/internal/parser"
)

func TestCheckDuplicates(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "dup.yaml")
	content := `
apiVersion: v1
kind: Pod
metadata:
  name: a
  name: b
`
	require.NoError(t, os.WriteFile(f, []byte(content), 0644))

	node, err := parser.ParseFile(f)
	require.NoError(t, err)

	errors := CheckDuplicates(node)
	require.Len(t, errors, 1)
	assert.Equal(t, "DuplicateKey", errors[0].Type)
	assert.Contains(t, errors[0].Message, "name")
}

func TestCheckDuplicates_NoDuplicates(t *testing.T) {
	var node yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte("a: 1\nb: 2"), &node))

	errors := CheckDuplicates(&node)
	assert.Empty(t, errors)
}
