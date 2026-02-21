package validator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"yaml-validator/internal/config"
)

func TestCheckForbidDefaultValues(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "deploy.yaml")
	yaml := `apiVersion: v1
kind: Deployment
metadata:
  name: app
spec:
  template:
    spec:
      restartPolicy: Always
      containers:
        - name: app
          image: nginx
          imagePullPolicy: Always
`
	require.NoError(t, os.WriteFile(f, []byte(yaml), 0644))

	cfg := config.DefaultConfig()
	cfg.Rules.ForbidDefaultValues = true
	cfg.Rules.DefaultValues = map[string]string{
		"imagePullPolicy": "Always",
		"restartPolicy":   "Always",
	}

	errors, err := Validate(f, cfg)
	require.NoError(t, err)
	require.Len(t, errors, 2)
	types := map[string]bool{}
	for _, e := range errors {
		assert.Equal(t, "ForbidDefaultValue", e.Type)
		types[e.Message] = true
	}
	assert.True(t, types["key 'imagePullPolicy' has default value 'Always', remove it"])
	assert.True(t, types["key 'restartPolicy' has default value 'Always', remove it"])
}

func TestCheckForbidDefaultValues_Ok(t *testing.T) {
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
          imagePullPolicy: IfNotPresent
`
	require.NoError(t, os.WriteFile(f, []byte(yaml), 0644))

	cfg := config.DefaultConfig()
	cfg.Rules.ForbidDefaultValues = true
	cfg.Rules.DefaultValues = map[string]string{
		"imagePullPolicy": "Always",
	}

	errors, err := Validate(f, cfg)
	require.NoError(t, err)
	assert.Empty(t, errors)
}

func TestCheckForbidDefaultValues_Disabled(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "deploy.yaml")
	yaml := `apiVersion: v1
kind: Deployment
metadata:
  name: app
spec:
  containers:
    - image: nginx
      imagePullPolicy: Always
`
	require.NoError(t, os.WriteFile(f, []byte(yaml), 0644))

	cfg := config.DefaultConfig()
	cfg.Rules.ForbidDefaultValues = false
	cfg.Rules.DefaultValues = map[string]string{"imagePullPolicy": "Always"}

	errors, err := Validate(f, cfg)
	require.NoError(t, err)
	assert.Empty(t, errors)
}
