package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchFileProfile(t *testing.T) {
	tests := []struct {
		pattern  string
		filePath string
		want     bool
	}{
		{"**/k8s/**", "deploy/k8s/deployment.yaml", true},
		{"**/k8s/**", "k8s/service.yaml", true},
		{"**/k8s/**", "base/app.yaml", false},
		{"*docker-compose*.yaml", "docker-compose.yaml", true},
		{"*docker-compose*.yaml", "docker-compose.prod.yaml", true},
		{"*docker-compose*.yaml", "other.yaml", false},
		{"**/k8s/**", filepath.Join("a", "b", "k8s", "c.yaml"), true},
	}
	for _, tt := range tests {
		got := MatchFileProfile(tt.pattern, tt.filePath)
		assert.Equalf(t, tt.want, got, "MatchFileProfile(%q, %q)", tt.pattern, tt.filePath)
	}
}

func TestConfigForFile_NoProfiles(t *testing.T) {
	base := DefaultConfig()
	got := ConfigForFile(base, "k8s/deploy.yaml")
	assert.Equal(t, base, got)
}
