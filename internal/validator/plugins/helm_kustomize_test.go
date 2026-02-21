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

func TestHelmKustomizeValidator_HelmChart(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "Chart.yaml")
	content := `apiVersion: v2
name: my-chart
version: 0.1.0
`
	require.NoError(t, os.WriteFile(f, []byte(content), 0644))

	cfg := config.DefaultConfig()
	cfg.Rules.CheckIntegrity = false
	cfg.Rules.RequiredFields = nil

	errors, err := validator.Validate(f, cfg)
	require.NoError(t, err)
	for _, e := range errors {
		if e.Type == "HelmChartMissingAPIVersion" || e.Type == "HelmChartMissingName" || e.Type == "HelmChartMissingVersion" {
			t.Errorf("Helm Chart.yaml with all required fields should pass, got %s", e.Type)
		}
	}
}

func TestHelmKustomizeValidator_HelmChartMissingVersion(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "Chart.yaml")
	content := `apiVersion: v2
name: my-chart
`
	require.NoError(t, os.WriteFile(f, []byte(content), 0644))

	cfg := config.DefaultConfig()
	cfg.Rules.CheckIntegrity = false
	cfg.Rules.RequiredFields = nil

	errors, err := validator.Validate(f, cfg)
	require.NoError(t, err)
	var helmErrs []string
	for _, e := range errors {
		if e.Type == "HelmChartMissingVersion" {
			helmErrs = append(helmErrs, e.Message)
		}
	}
	require.Len(t, helmErrs, 1)
	assert.Contains(t, helmErrs[0], "version")
}

func TestHelmKustomizeValidator_Kustomize(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "kustomization.yaml")
	content := `apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - deployment.yaml
`
	require.NoError(t, os.WriteFile(f, []byte(content), 0644))

	cfg := config.DefaultConfig()
	cfg.Rules.CheckIntegrity = false
	cfg.Rules.RequiredFields = nil

	errors, err := validator.Validate(f, cfg)
	require.NoError(t, err)
	for _, e := range errors {
		if e.Type == "KustomizeMissingResources" {
			t.Errorf("Kustomize with resources should pass, got %s", e.Type)
		}
	}
}

func TestHelmKustomizeValidator_KustomizeMissingResources(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "kustomization.yaml")
	content := `apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
`
	require.NoError(t, os.WriteFile(f, []byte(content), 0644))

	cfg := config.DefaultConfig()
	cfg.Rules.CheckIntegrity = false
	cfg.Rules.RequiredFields = nil

	errors, err := validator.Validate(f, cfg)
	require.NoError(t, err)
	var kustErrs []string
	for _, e := range errors {
		if e.Type == "KustomizeMissingResources" {
			kustErrs = append(kustErrs, e.Message)
		}
	}
	require.Len(t, kustErrs, 1)
	assert.Contains(t, kustErrs[0], "resources")
}

func TestHelmKustomizeValidator_SkipsK8s(t *testing.T) {
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
		if e.Type == "HelmChartMissingVersion" || e.Type == "KustomizeMissingResources" {
			t.Errorf("Helm/Kustomize plugin should skip K8s files, got %s", e.Type)
		}
	}
}
