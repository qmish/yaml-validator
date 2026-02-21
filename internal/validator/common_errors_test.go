package validator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckCommonErrors_Tabs(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "tabs.yaml")
	require.NoError(t, os.WriteFile(f, []byte("key:\n\tvalue: 1\n"), 0644))

	errors := CheckCommonErrors(f, 200, nil, false)
	require.Len(t, errors, 1)
	assert.Equal(t, "TabInsteadOfSpaces", errors[0].Type)
}

func TestCheckCommonErrors_LineTooLong(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "long.yaml")
	line := "a: " + string(make([]byte, 250))
	require.NoError(t, os.WriteFile(f, []byte(line), 0644))

	errors := CheckCommonErrors(f, 200, nil, false)
	require.Len(t, errors, 1)
	assert.Equal(t, "LineTooLong", errors[0].Type)
}

func TestCheckCommonErrors_SensitiveData(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "secret.yaml")
	require.NoError(t, os.WriteFile(f, []byte("password: secret123\n"), 0644))

	errors := CheckCommonErrors(f, 200, []string{"password"}, false)
	require.Len(t, errors, 1)
	assert.Equal(t, "SensitiveData", errors[0].Type)
}
