package validator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"yaml-validator/internal/config"
)

func TestCheckUniqueListFields(t *testing.T) {
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
        - name: app
          image: redis
        - name: sidecar
          image: busybox
`
	require.NoError(t, os.WriteFile(f, []byte(yaml), 0644))

	cfg := config.DefaultConfig()
	cfg.Rules.UniqueListFields = []config.UniqueListField{
		{Path: "spec.template.spec.containers", Field: "name"},
	}

	errors, err := Validate(f, cfg)
	require.NoError(t, err)
	require.Len(t, errors, 1)
	assert.Equal(t, "DuplicateListElement", errors[0].Type)
	assert.Contains(t, errors[0].Message, "duplicate value 'app'")
	assert.Contains(t, errors[0].Message, "containers")
}

func TestCheckUniqueListFields_Ok(t *testing.T) {
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
        - name: redis
          image: redis
`
	require.NoError(t, os.WriteFile(f, []byte(yaml), 0644))

	cfg := config.DefaultConfig()
	cfg.Rules.UniqueListFields = []config.UniqueListField{
		{Path: "spec.template.spec.containers", Field: "name"},
	}

	errors, err := Validate(f, cfg)
	require.NoError(t, err)
	assert.Empty(t, errors)
}

func TestCheckUniqueListFields_Disabled(t *testing.T) {
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
        - name: app
          image: redis
`
	require.NoError(t, os.WriteFile(f, []byte(yaml), 0644))

	cfg := config.DefaultConfig()
	cfg.Rules.UniqueListFields = nil

	errors, err := Validate(f, cfg)
	require.NoError(t, err)
	assert.Empty(t, errors)
}
