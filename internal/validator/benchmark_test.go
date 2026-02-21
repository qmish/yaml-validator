package validator

import (
	"os"
	"path/filepath"
	"testing"

	"yaml-validator/internal/config"
)

// generateLargeYAML creates YAML of ~size bytes for benchmark
func generateLargeYAML(size int) []byte {
	const block = `  key: value
  another: field
  nested:
    a: 1
    b: 2

`
	var b []byte
	for len(b) < size {
		b = append(b, block...)
	}
	return append([]byte("apiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: test\ndata:\n"), b...)
}

func BenchmarkValidate_Small(b *testing.B) {
	tmp := b.TempDir()
	f := filepath.Join(tmp, "small.yaml")
	content := []byte("apiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: test\ndata:\n  key: value\n")
	if err := os.WriteFile(f, content, 0644); err != nil {
		b.Fatal(err)
	}
	cfg := config.DefaultConfig()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Validate(f, cfg)
	}
}

func BenchmarkValidate_Medium(b *testing.B) {
	tmp := b.TempDir()
	f := filepath.Join(tmp, "medium.yaml")
	content := generateLargeYAML(10 * 1024)
	if err := os.WriteFile(f, content, 0644); err != nil {
		b.Fatal(err)
	}
	cfg := config.DefaultConfig()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Validate(f, cfg)
	}
}

func BenchmarkValidate_Large(b *testing.B) {
	tmp := b.TempDir()
	f := filepath.Join(tmp, "large.yaml")
	content := generateLargeYAML(100 * 1024)
	if err := os.WriteFile(f, content, 0644); err != nil {
		b.Fatal(err)
	}
	cfg := config.DefaultConfig()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Validate(f, cfg)
	}
}
