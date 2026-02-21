package regression

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"yaml-validator/internal/config"
	"yaml-validator/internal/reporter"
	"yaml-validator/internal/validator"
)

// TestRegression проверяет, что вывод валидации не изменился для эталонных YAML.
// Golden-файлы: testdata/regression/*.yaml.golden содержат ожидаемый compact-вывод.
// При изменении поведения обновите .golden вручную.
func TestRegression(t *testing.T) {
	dir := filepath.Join("..", "..", "testdata", "regression")
	entries, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
	require.NoError(t, err)

	for _, yamlPath := range entries {
		if strings.HasSuffix(yamlPath, ".golden") {
			continue
		}
		t.Run(filepath.Base(yamlPath), func(t *testing.T) {
			goldenPath := yamlPath + ".golden"
			cfg := config.DefaultConfig()
			errors, _ := validator.Validate(yamlPath, cfg)

			// путь для вывода: testdata/regression/file.yaml (кросс-платформа)
			relPath := filepath.ToSlash(filepath.Join("testdata", "regression", filepath.Base(yamlPath)))

			actual := reporter.GenerateCompactReport(relPath, errors)

			goldenBytes, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatalf("golden file %s: %v", goldenPath, err)
			}
			expected := string(goldenBytes)

			assert.Equal(t, expected, actual,
				"Regression: output changed for %s. Update %s if intentional.", yamlPath, goldenPath)
		})
	}
}
