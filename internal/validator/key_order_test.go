package validator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"yaml-validator/internal/parser"
	"yaml-validator/pkg"
)

func TestCheckKeyOrdering(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "out_of_order.yaml")
	content := []byte("---\nz: 1\nb: 2\na: 3\n")
	require.NoError(t, os.WriteFile(f, content, 0644))

	node, err := parser.ParseFile(f)
	require.NoError(t, err)

	errors := CheckKeyOrdering(node)
	require.GreaterOrEqual(t, len(errors), 1)
	var keyOrderErrs []pkg.Error
	for _, e := range errors {
		if e.Type == "KeyOrdering" {
			keyOrderErrs = append(keyOrderErrs, e)
		}
	}
	assert.GreaterOrEqual(t, len(keyOrderErrs), 1)
	assert.Contains(t, keyOrderErrs[0].Message, "should be before")
}

func TestCheckKeyOrdering_Ok(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "ordered.yaml")
	content := []byte("---\na: 1\nb: 2\nz: 3\n")
	require.NoError(t, os.WriteFile(f, content, 0644))

	node, err := parser.ParseFile(f)
	require.NoError(t, err)

	errors := CheckKeyOrdering(node)
	assert.Empty(t, errors)
}

func TestCheckKeyOrderConfigurable(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "k8s.yaml")
	// spec before metadata - wrong order for K8s
	content := []byte("apiVersion: apps/v1\nkind: Deployment\nspec: {}\nmetadata:\n  name: foo\n")
	require.NoError(t, os.WriteFile(f, content, 0644))

	node, err := parser.ParseFile(f)
	require.NoError(t, err)

	order := []string{"apiVersion", "kind", "metadata", "spec"}
	errors := CheckKeyOrderConfigurable(node, order)
	require.NotEmpty(t, errors)
	assert.Contains(t, errors[0].Message, "metadata")
	assert.Contains(t, errors[0].Message, "spec")
}

func TestCheckKeyOrderConfigurable_Ok(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "k8s.yaml")
	content := []byte("apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: foo\nspec: {}\n")
	require.NoError(t, os.WriteFile(f, content, 0644))

	node, err := parser.ParseFile(f)
	require.NoError(t, err)

	order := []string{"apiVersion", "kind", "metadata", "spec"}
	errors := CheckKeyOrderConfigurable(node, order)
	assert.Empty(t, errors)
}

func TestCheckMaxKeyNameLength(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "long_keys.yaml")
	content := []byte("short: 1\nvery_long_key_name_that_exceeds_limit: 2\n")
	require.NoError(t, os.WriteFile(f, content, 0644))

	node, err := parser.ParseFile(f)
	require.NoError(t, err)

	errors := CheckMaxKeyNameLength(node, 20)
	require.Len(t, errors, 1)
	assert.Equal(t, "MaxKeyNameLength", errors[0].Type)
}

func TestCheckMaxKeyNameLength_Ok(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "keys.yaml")
	content := []byte("a: 1\nb: 2\n")
	require.NoError(t, os.WriteFile(f, content, 0644))

	node, err := parser.ParseFile(f)
	require.NoError(t, err)

	errors := CheckMaxKeyNameLength(node, 50)
	assert.Empty(t, errors)
}
