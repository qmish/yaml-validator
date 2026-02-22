package validator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckConsistency_Inconsistent(t *testing.T) {
	dir := t.TempDir()
	f1 := filepath.Join(dir, "a.yaml")
	f2 := filepath.Join(dir, "b.yaml")
	_ = os.WriteFile(f1, []byte("version: 1.0\n"), 0644)
	_ = os.WriteFile(f2, []byte("version: 2.0\n"), 0644)

	errs := CheckConsistency([]string{f1, f2}, []string{"version"})
	if len(errs) != 1 {
		t.Fatalf("Expected 1 error, got %d", len(errs))
	}
	if errs[0].Type != "InconsistentValue" {
		t.Errorf("Expected InconsistentValue, got %s", errs[0].Type)
	}
}

func TestCheckConsistency_Consistent(t *testing.T) {
	dir := t.TempDir()
	f1 := filepath.Join(dir, "a.yaml")
	f2 := filepath.Join(dir, "b.yaml")
	_ = os.WriteFile(f1, []byte("version: 1.0\n"), 0644)
	_ = os.WriteFile(f2, []byte("version: 1.0\n"), 0644)

	errs := CheckConsistency([]string{f1, f2}, []string{"version"})
	if len(errs) != 0 {
		t.Errorf("Expected no errors, got %d: %v", len(errs), errs)
	}
}
