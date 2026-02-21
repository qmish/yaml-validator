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

func TestCheckStyle_ForbidTrailingDots(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "doc.yaml")
	require.NoError(t, os.WriteFile(f, []byte("key: value.\n"), 0644))

	opts := config.StyleOptions{ForbidTrailingDots: true}
	errors := CheckStyle(f, opts)
	require.Len(t, errors, 1)
	assert.Equal(t, "TrailingDots", errors[0].Type)
}

func TestCheckStyle_ForbidTrailingDots_Ok(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "doc.yaml")
	require.NoError(t, os.WriteFile(f, []byte("key: value\n---\n...\n"), 0644))

	opts := config.StyleOptions{ForbidTrailingDots: true}
	errors := CheckStyle(f, opts)
	assert.Empty(t, errors)
}

func TestCheckStyle_ForbidTabs(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "doc.yaml")
	require.NoError(t, os.WriteFile(f, []byte("key:\n\tvalue: 1\n"), 0644))

	opts := config.StyleOptions{ForbidTabs: true}
	errors := CheckStyle(f, opts)
	require.Len(t, errors, 1)
	assert.Equal(t, "TabInsteadOfSpaces", errors[0].Type)
}

func TestCheckStyle_ForbidTabs_Ok(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "doc.yaml")
	require.NoError(t, os.WriteFile(f, []byte("key:\n  value: 1\n"), 0644))

	opts := config.StyleOptions{ForbidTabs: true}
	errors := CheckStyle(f, opts)
	assert.Empty(t, errors)
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

func TestCheckStyle_ConsecutiveEmptyLines(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "doc.yaml")
	require.NoError(t, os.WriteFile(f, []byte("a: 1\n\n\nb: 2\n"), 0644))

	opts := config.StyleOptions{ForbidConsecutiveEmptyLines: true}
	errors := CheckStyle(f, opts)
	require.Len(t, errors, 1)
	assert.Equal(t, "ConsecutiveEmptyLines", errors[0].Type)
	assert.Equal(t, 3, errors[0].Line)
}

func TestCheckStyle_ConsecutiveEmptyLines_Ok(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "doc.yaml")
	require.NoError(t, os.WriteFile(f, []byte("a: 1\n\nb: 2\n"), 0644))

	opts := config.StyleOptions{ForbidConsecutiveEmptyLines: true}
	errors := CheckStyle(f, opts)
	assert.Empty(t, errors)
}

func TestCheckStyle_RequireEmptyLineBetweenBlocks(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "doc.yaml")
	// два топ-уровневых ключа без пустой строки между ними
	require.NoError(t, os.WriteFile(f, []byte("apiVersion: v1\nkind: Pod\nmetadata:\n  name: x\n"), 0644))

	opts := config.StyleOptions{RequireEmptyLineBetweenBlocks: true}
	errors := CheckStyle(f, opts)
	require.Len(t, errors, 2) // kind без разрыва после apiVersion, metadata без разрыва после kind
	assert.Equal(t, "EmptyLineBetweenBlocks", errors[0].Type)
}

func TestCheckStyle_MinEmptyLinesBetweenBlocks(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "doc.yaml")
	require.NoError(t, os.WriteFile(f, []byte("a: 1\nb: 2\n"), 0644))

	opts := config.StyleOptions{MinEmptyLinesBetweenBlocks: 1}
	errors := CheckStyle(f, opts)
	require.Len(t, errors, 1)
	assert.Equal(t, "EmptyLineBetweenBlocks", errors[0].Type)
}

func TestCheckStyle_MinEmptyLinesBetweenBlocks_Zero(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "doc.yaml")
	require.NoError(t, os.WriteFile(f, []byte("a: 1\nb: 2\n"), 0644))

	opts := config.StyleOptions{MinEmptyLinesBetweenBlocks: 0}
	errors := CheckStyle(f, opts)
	assert.Empty(t, errors)
}

func TestCheckStyle_RequireEmptyLineBetweenBlocks_Ok(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "doc.yaml")
	require.NoError(t, os.WriteFile(f, []byte("apiVersion: v1\n\nkind: Pod\n\nmetadata:\n  name: x\n"), 0644))

	opts := config.StyleOptions{RequireEmptyLineBetweenBlocks: true}
	errors := CheckStyle(f, opts)
	assert.Empty(t, errors)
}

func TestCheckStyle_RequireDocumentEnd(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "doc.yaml")
	require.NoError(t, os.WriteFile(f, []byte("---\nkey: value\n"), 0644))

	opts := config.StyleOptions{RequireDocumentEnd: true}
	errors := CheckStyle(f, opts)
	require.Len(t, errors, 1)
	assert.Equal(t, "DocumentEnd", errors[0].Type)
}

func TestCheckStyle_RequireDocumentEnd_Ok(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "doc.yaml")
	require.NoError(t, os.WriteFile(f, []byte("---\nkey: value\n...\n"), 0644))

	opts := config.StyleOptions{RequireDocumentEnd: true}
	errors := CheckStyle(f, opts)
	assert.Empty(t, errors)
}

func TestCheckStyle_RequireCommentsIndented(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "doc.yaml")
	require.NoError(t, os.WriteFile(f, []byte("key:\n  sub: 1\n# unindented comment\n  other: 2\n"), 0644))

	opts := config.StyleOptions{RequireCommentsIndented: true}
	errors := CheckStyle(f, opts)
	require.Len(t, errors, 1)
	assert.Equal(t, "CommentIndentation", errors[0].Type)
	assert.Equal(t, 3, errors[0].Line)
}

func TestCheckStyle_RequireCommentsIndented_Ok(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "doc.yaml")
	require.NoError(t, os.WriteFile(f, []byte("key:\n  sub: 1\n  # indented comment\n  other: 2\n"), 0644))

	opts := config.StyleOptions{RequireCommentsIndented: true}
	errors := CheckStyle(f, opts)
	assert.Empty(t, errors)
}

func TestCheckStyle_ForbidUnicode(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "doc.yaml")
	require.NoError(t, os.WriteFile(f, []byte("key: значение\n"), 0644))

	opts := config.StyleOptions{ForbidUnicode: true}
	errors := CheckStyle(f, opts)
	require.NotEmpty(t, errors)
	assert.Equal(t, "ForbidUnicode", errors[0].Type)
}

func TestCheckStyle_ForbidUnicode_Ok(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "doc.yaml")
	require.NoError(t, os.WriteFile(f, []byte("key: value\nname: test\n"), 0644))

	opts := config.StyleOptions{ForbidUnicode: true}
	errors := CheckStyle(f, opts)
	assert.Empty(t, errors)
}

func TestCheckStyle_RequireQuotedKeys(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "doc.yaml")
	require.NoError(t, os.WriteFile(f, []byte("key: value\n"), 0644))

	opts := config.StyleOptions{RequireQuotedKeys: true}
	errors := CheckStyle(f, opts)
	require.Len(t, errors, 1)
	assert.Equal(t, "QuotedKeys", errors[0].Type)
	assert.Equal(t, 1, errors[0].Line)
}

func TestCheckStyle_RequireQuotedKeys_Ok(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "doc.yaml")
	require.NoError(t, os.WriteFile(f, []byte("\"key\": value\n"), 0644))

	opts := config.StyleOptions{RequireQuotedKeys: true}
	errors := CheckStyle(f, opts)
	assert.Empty(t, errors)
}

func TestCheckStyle_IndentSpaces(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "doc.yaml")
	// 3 spaces indent - not multiple of 2
	content := []byte("a: 1\nb:\n   3\n")
	require.NoError(t, os.WriteFile(f, content, 0644))

	opts := config.StyleOptions{IndentSpaces: 2}
	errors := CheckStyle(f, opts)
	require.NotEmpty(t, errors)
	assert.Equal(t, "IndentSpaces", errors[0].Type)
}

func TestCheckStyle_IndentSpaces_Ok(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "doc.yaml")
	content := []byte("a: 1\nb:\n  2\n    3\n")
	require.NoError(t, os.WriteFile(f, content, 0644))

	opts := config.StyleOptions{IndentSpaces: 2}
	errors := CheckStyle(f, opts)
	assert.Empty(t, errors)
}

func TestCheckStyle_Disabled(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "doc.yaml")
	require.NoError(t, os.WriteFile(f, []byte("key: value"), 0644))

	opts := config.StyleOptions{}
	errors := CheckStyle(f, opts)
	assert.Empty(t, errors)
}
