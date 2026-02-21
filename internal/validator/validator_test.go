package validator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"yaml-validator/internal/config"
)

func TestValidate_FullFlow(t *testing.T) {
	tmp := t.TempDir()
	validFile := filepath.Join(tmp, "valid.yaml")
	content := `
apiVersion: v1
kind: Pod
metadata:
  name: test
spec:
  containers: []
`
	require.NoError(t, os.WriteFile(validFile, []byte(content), 0644))

	cfg := config.DefaultConfig()
	errors, err := Validate(validFile, cfg)
	require.NoError(t, err)
	assert.Empty(t, errors)
}

func TestValidate_SyntaxError(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "bad.yaml")
	require.NoError(t, os.WriteFile(f, []byte("key: [unclosed\n"), 0644))

	cfg := config.DefaultConfig()
	errors, err := Validate(f, cfg)
	assert.Error(t, err)
	assert.NotEmpty(t, errors)
	assert.Equal(t, "SyntaxError", errors[0].Type)
}

func TestValidate_DuplicateKeys(t *testing.T) {
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

	cfg := config.DefaultConfig()
	errors, err := Validate(f, cfg)
	require.NoError(t, err)
	assert.NotEmpty(t, errors)
	types := make(map[string]bool)
	for _, e := range errors {
		types[e.Type] = true
	}
	assert.True(t, types["DuplicateKey"], "expected DuplicateKey error")
}

func TestValidate_RuleSeverity(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "long.yaml")
	// Строка длиннее 80 символов
	longLine := "key: value" + strings.Repeat("x", 80)
	require.NoError(t, os.WriteFile(f, []byte(longLine), 0644))

	cfg := config.DefaultConfig()
	cfg.Rules.CheckIntegrity = false
	cfg.Rules.RequiredFields = nil
	cfg.Rules.MaxLineLength = 80
	cfg.Rules.RuleSeverity = map[string]string{"LineTooLong": "warning"}

	errors, err := Validate(f, cfg)
	require.NoError(t, err)
	require.NotEmpty(t, errors)
	for _, e := range errors {
		if e.Type == "LineTooLong" {
			assert.Equal(t, "warning", e.Severity)
			return
		}
	}
	t.Fatal("expected LineTooLong error")
}
