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

func TestParseFileMulti(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "multi.yaml")
	content := `---
apiVersion: v1
kind: Pod
metadata:
  name: first
---
apiVersion: v1
kind: Pod
metadata:
  name: second
`
	require.NoError(t, os.WriteFile(f, []byte(content), 0644))

	nodes, err := ParseFileMulti(f)
	require.NoError(t, err)
	require.Len(t, nodes, 2)
	assert.Equal(t, yaml.DocumentNode, nodes[0].Kind)
	assert.Equal(t, yaml.DocumentNode, nodes[1].Kind)
}

func TestParseFileMulti_SingleDoc(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "single.yaml")
	content := `apiVersion: v1
kind: Pod
`
	require.NoError(t, os.WriteFile(f, []byte(content), 0644))

	nodes, err := ParseFileMulti(f)
	require.NoError(t, err)
	require.Len(t, nodes, 1)
}
