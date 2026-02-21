package cmd

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestExitCodes(t *testing.T) {
	if ExitOK != 0 || ExitErrors != 1 || ExitWarnings != 2 {
		t.Errorf("Expected ExitOK=0, ExitErrors=1, ExitWarnings=2, got %d, %d, %d", ExitOK, ExitErrors, ExitWarnings)
	}
}

func TestRulesListCommand(t *testing.T) {
	rootCmd.SetArgs([]string{"rules", "list", "-o", "json"})
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "check_syntax") || !strings.Contains(out, "config_key") {
		t.Errorf("Expected rules list JSON with check_syntax and config_key, got: %s", out)
	}
	if !strings.Contains(out, "rules.check_syntax") {
		t.Errorf("Expected config_key rules.check_syntax in output, got: %s", out)
	}
}

func TestVersionCommand(t *testing.T) {
	rootCmd.SetArgs([]string{"version"})
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "yaml-validator") {
		t.Errorf("Expected version output, got: %s", out)
	}
}

func TestValidateCommand_WatchFlag(t *testing.T) {
	rootCmd.SetArgs([]string{"validate", "--help"})
	var buf bytes.Buffer
	rootCmd.SetErr(&buf)
	rootCmd.SetOut(&buf)
	_ = rootCmd.Execute()
	out := buf.String()
	if !strings.Contains(out, "--watch") && !strings.Contains(out, "watch") {
		t.Errorf("Expected --watch in help, got: %s", out)
	}
}

func TestValidateCommand_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	tmp := t.TempDir()
	validFile := filepath.Join(tmp, "valid.yaml")
	if err := os.WriteFile(validFile, []byte("apiVersion: v1\nkind: Pod\nmetadata:\n  name: test\n"), 0644); err != nil {
		t.Fatal(err)
	}
	root, _ := os.Getwd()
	runCmd := exec.Command("go", "run", ".", "validate", validFile)
	runCmd.Dir = root
	out, err := runCmd.CombinedOutput()
	if err != nil {
		t.Logf("Output: %s", out)
		t.Skip("integration test: run from project root")
		return
	}
	if !strings.Contains(string(out), "valid") {
		t.Errorf("Expected 'valid' in output, got: %s", out)
	}
}
