package validator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"yaml-validator/internal/config"
)

func TestCheckK8sSchema_Disabled(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "pod.yaml")
	require.NoError(t, os.WriteFile(f, []byte("apiVersion: v1\nkind: Pod\nmetadata:\n  name: x\n"), 0644))

	opts := config.K8sSchemaOptions{Enabled: false}
	errors := CheckK8sSchema(f, opts)
	assert.Empty(t, errors)
}

func TestCheckK8sSchema_InvalidType(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping K8s schema test (network) in short mode")
	}
	tmp := t.TempDir()
	f := filepath.Join(tmp, "bad.yaml")
	// replicas должен быть число, не строка — схема должна вернуть ошибку
	content := `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test
spec:
  replicas: "three"
  selector:
    matchLabels:
      app: test
  template:
    metadata:
      labels:
        app: test
    spec:
      containers: []
`
	require.NoError(t, os.WriteFile(f, []byte(content), 0644))

	opts := config.K8sSchemaOptions{
		Enabled:       true,
		Version:       "1.28.0",
		IgnoreMissing: true,
	}
	errors := CheckK8sSchema(f, opts)
	require.NotEmpty(t, errors, "expected schema validation error for replicas: \"three\"")
	assert.Equal(t, "K8sSchema", errors[0].Type)
}
