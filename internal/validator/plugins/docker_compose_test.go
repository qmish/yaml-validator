package plugins

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"yaml-validator/internal/parser"
)

func TestDockerComposeValidator_Ok(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "docker-compose.yaml")
	content := []byte(`version: "3.8"
services:
  web:
    image: nginx:latest
  app:
    build: .
`)
	require.NoError(t, os.WriteFile(f, content, 0644))

	node, err := parser.ParseFile(f)
	require.NoError(t, err)

	validator := &DockerComposeValidator{}
	errors := validator.Validate(node)
	assert.Empty(t, errors)
}

func TestDockerComposeValidator_MissingImageOrBuild(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "docker-compose.yaml")
	content := []byte(`services:
  bad:
    environment:
      FOO: bar
`)
	require.NoError(t, os.WriteFile(f, content, 0644))

	node, err := parser.ParseFile(f)
	require.NoError(t, err)

	validator := &DockerComposeValidator{}
	errors := validator.Validate(node)
	require.Len(t, errors, 1)
	assert.Equal(t, "DockerComposeServiceImage", errors[0].Type)
	assert.Contains(t, errors[0].Message, "image")
	assert.Contains(t, errors[0].Message, "build")
}

func TestDockerComposeValidator_SkipsK8s(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "k8s.yaml")
	content := []byte(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: foo
spec:
  template:
    spec:
      containers:
        - name: app
          image: nginx
`)
	require.NoError(t, os.WriteFile(f, content, 0644))

	node, err := parser.ParseFile(f)
	require.NoError(t, err)

	validator := &DockerComposeValidator{}
	errors := validator.Validate(node)
	assert.Empty(t, errors)
}
