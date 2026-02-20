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
