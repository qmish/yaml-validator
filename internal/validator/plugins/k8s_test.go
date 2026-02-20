package plugins

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"yaml-validator/internal/parser"
)

func TestK8sValidator_ValidManifest(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "k8s.yaml")
	content := `
apiVersion: v1
kind: Pod
metadata:
  name: test-pod
`
	require.NoError(t, os.WriteFile(f, []byte(content), 0644))

	node, err := parser.ParseFile(f)
	require.NoError(t, err)

	k := &K8sValidator{}
	errors := k.Validate(node)
	assert.Empty(t, errors)
}

func TestK8sValidator_MissingKind(t *testing.T) {
	var node yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte("apiVersion: v1\nmetadata:\n  name: x\n"), &node))

	k := &K8sValidator{}
	errors := k.Validate(&node)
	require.Len(t, errors, 1)
	assert.Equal(t, "K8sMissingKind", errors[0].Type)
}

func TestK8sValidator_UnknownKind(t *testing.T) {
	var node yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte("apiVersion: v1\nkind: UnknownResource\nmetadata:\n  name: x\n"), &node))

	k := &K8sValidator{}
	errors := k.Validate(&node)
	require.Len(t, errors, 1)
	assert.Equal(t, "K8sUnknownKind", errors[0].Type)
}

func TestK8sValidator_NonK8sDocument(t *testing.T) {
	var node yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte("services:\n  web:\n    image: nginx\n"), &node))

	k := &K8sValidator{}
	errors := k.Validate(&node)
	assert.Empty(t, errors)
}
