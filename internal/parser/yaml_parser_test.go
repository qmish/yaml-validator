package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestParseFile(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "test.yaml")
	content := `apiVersion: v1
kind: Pod
metadata:
  name: test
`
	require.NoError(t, os.WriteFile(f, []byte(content), 0644))

	node, err := ParseFile(f)
	require.NoError(t, err)
	require.NotNil(t, node)
	assert.Equal(t, yaml.DocumentNode, node.Kind)
}

func TestParseFile_NotFound(t *testing.T) {
	_, err := ParseFile("/nonexistent/file.yaml")
	assert.Error(t, err)
}
