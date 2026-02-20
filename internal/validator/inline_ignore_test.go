package validator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"yaml-validator/pkg"
)

func TestParseInlineIgnore_DisableLine(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "f.yaml")
	// Комментарий на строке 1 отключает правило для строки 1
	require.NoError(t, os.WriteFile(f, []byte("key: value   # yaml-validator disable-line rule:TrailingSpaces\n"), 0644))

	ignore, err := parseInlineIgnore(f)
	require.NoError(t, err)
	require.NotEmpty(t, ignore)
	assert.True(t, ignore.isDisabled(1, "TrailingSpaces"))
	assert.False(t, ignore.isDisabled(1, "LineTooLong"))
}

func TestParseInlineIgnore_DisableNextLine(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "f.yaml")
	content := "# yaml-validator disable-next-line rule:LineTooLong\nlong line here\n"
	require.NoError(t, os.WriteFile(f, []byte(content), 0644))

	ignore, err := parseInlineIgnore(f)
	require.NoError(t, err)
	assert.True(t, ignore.isDisabled(2, "LineTooLong"))
}

func TestParseInlineIgnore_DisableAllForLine(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "f.yaml")
	require.NoError(t, os.WriteFile(f, []byte("# yaml-validator disable-line\nx: 1\n"), 0644))

	ignore, err := parseInlineIgnore(f)
	require.NoError(t, err)
	assert.True(t, ignore.isDisabled(1, "AnyRule"))
}

func TestFilterInlineIgnore(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "f.yaml")
	require.NoError(t, os.WriteFile(f, []byte("line1\n# yaml-validator disable-line rule:TrailingSpaces\n"), 0644))

	errors := []pkg.Error{
		{Type: "TrailingSpaces", Line: 2, Message: "trailing"},
		{Type: "LineTooLong", Line: 2, Message: "long"},
	}
	filtered := FilterInlineIgnore(f, errors)
	require.Len(t, filtered, 1)
	assert.Equal(t, "LineTooLong", filtered[0].Type)
}
