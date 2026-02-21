package parser

import (
	"os"
	"path/filepath"
	"testing"
)

// generateLargeYAML создаёт YAML размером ~size bytes
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

func BenchmarkParseFile_Small(b *testing.B) {
	tmp := b.TempDir()
	f := filepath.Join(tmp, "small.yaml")
	content := []byte("apiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: test\ndata:\n  key: value\n")
	if err := os.WriteFile(f, content, 0644); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ParseFile(f)
	}
}

func BenchmarkParseFile_Large(b *testing.B) {
	tmp := b.TempDir()
	f := filepath.Join(tmp, "large.yaml")
	content := generateLargeYAML(100 * 1024) // ~100 KB
	if err := os.WriteFile(f, content, 0644); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ParseFile(f)
	}
}
