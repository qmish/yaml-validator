package validator

import (
	"yaml-validator/pkg"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"yaml-validator/internal/config"
)

func TestCheckKeyValuePatterns_Keys(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "deploy.yaml")
	yaml := `apiVersion: v1
kind: Deployment
metadata:
  name: app
spec:
  template:
    spec:
      containers:
        - name: app
          image: nginx
`
	require.NoError(t, os.WriteFile(f, []byte(yaml), 0644))

	cfg := config.DefaultConfig()
	cfg.Rules.KeyValuePatterns = []config.KeyValuePattern{
		{Path: "metadata", Pattern: "^[a-z0-9]([-a-z0-9]*[a-z0-9])?$", Target: "keys"},
	}

	errors, err := Validate(f, cfg)
	require.NoError(t, err)
	assert.Empty(t, errors) // name matches pattern
}

func TestCheckKeyValuePatterns_Keys_Mismatch(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "bad.yaml")
	yaml := `apiVersion: v1
kind: ConfigMap
metadata:
  BadName: value
  name: ok
`
	require.NoError(t, os.WriteFile(f, []byte(yaml), 0644))

	cfg := config.DefaultConfig()
	cfg.Rules.KeyValuePatterns = []config.KeyValuePattern{
		{Path: "metadata", Pattern: "^[a-z][a-z0-9-]*$", Target: "keys"},
	}

	errors, err := Validate(f, cfg)
	require.NoError(t, err)
	var keyErrs []pkg.Error
	for _, e := range errors {
		if e.Type == "KeyPatternMismatch" {
			keyErrs = append(keyErrs, e)
		}
	}
	require.Len(t, keyErrs, 1)
	assert.Contains(t, keyErrs[0].Message, "BadName")
}

func TestCheckKeyValuePatterns_Values(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "deploy.yaml")
	yaml := `apiVersion: v1
kind: Deployment
metadata:
  name: my-app
`
	require.NoError(t, os.WriteFile(f, []byte(yaml), 0644))

	cfg := config.DefaultConfig()
	cfg.Rules.KeyValuePatterns = []config.KeyValuePattern{
		{Path: "metadata.name", Pattern: "^[a-z0-9]([-a-z0-9]*[a-z0-9])?$", Target: "values"},
	}

	errors, err := Validate(f, cfg)
	require.NoError(t, err)
	assert.Empty(t, errors)
}

func TestCheckKeyValuePatterns_Values_Mismatch(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "bad.yaml")
	yaml := `apiVersion: v1
kind: ConfigMap
metadata:
  name: Bad_Name
`
	require.NoError(t, os.WriteFile(f, []byte(yaml), 0644))

	cfg := config.DefaultConfig()
	cfg.Rules.KeyValuePatterns = []config.KeyValuePattern{
		{Path: "metadata.name", Pattern: "^[a-z0-9]([-a-z0-9]*[a-z0-9])?$", Target: "values"},
	}

	errors, err := Validate(f, cfg)
	require.NoError(t, err)
	var valErrs []pkg.Error
	for _, e := range errors {
		if e.Type == "ValuePatternMismatch" {
			valErrs = append(valErrs, e)
		}
	}
	require.Len(t, valErrs, 1)
	assert.Contains(t, valErrs[0].Message, "Bad_Name")
}
