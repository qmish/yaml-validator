package plugins

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"yaml-validator/internal/config"
	"yaml-validator/internal/validator"
)

func TestAnsibleValidator_PlayMissingHosts(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "playbook.yaml")
	content := `---
- tasks:
    - name: do something
      debug:
        msg: hello
`
	require.NoError(t, os.WriteFile(f, []byte(content), 0644))

	cfg := config.DefaultConfig()
	cfg.Rules.CheckIntegrity = false
	cfg.Rules.RequiredFields = nil

	errors, err := validator.Validate(f, cfg)
	require.NoError(t, err)
	var ansibleErrs []string
	for _, e := range errors {
		if e.Type == "AnsiblePlayMissingHosts" {
			ansibleErrs = append(ansibleErrs, e.Message)
		}
	}
	require.Len(t, ansibleErrs, 1)
	assert.Contains(t, ansibleErrs[0], "hosts")
}

func TestAnsibleValidator_PlayWithHosts(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "playbook.yaml")
	content := `---
- hosts: all
  tasks:
    - name: do something
      debug:
        msg: hello
`
	require.NoError(t, os.WriteFile(f, []byte(content), 0644))

	cfg := config.DefaultConfig()
	cfg.Rules.CheckIntegrity = false
	cfg.Rules.RequiredFields = nil

	errors, err := validator.Validate(f, cfg)
	require.NoError(t, err)
	for _, e := range errors {
		if e.Type == "AnsiblePlayMissingHosts" {
			t.Errorf("unexpected AnsiblePlayMissingHosts when hosts is present")
		}
	}
}

func TestAnsibleValidator_SkipsK8s(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "deploy.yaml")
	content := `apiVersion: v1
kind: Pod
metadata:
  name: test
`
	require.NoError(t, os.WriteFile(f, []byte(content), 0644))

	cfg := config.DefaultConfig()
	errors, err := validator.Validate(f, cfg)
	require.NoError(t, err)
	for _, e := range errors {
		if e.Type == "AnsiblePlayMissingHosts" || e.Type == "AnsibleTaskNoModule" {
			t.Errorf("Ansible plugin should skip K8s files, got %s", e.Type)
		}
	}
}
