package validator

import (
	"yaml-validator/internal/config"
	"yaml-validator/internal/parser"
	"yaml-validator/pkg"
)

// Validate выполняет полную валидацию YAML-файла согласно конфигурации
func Validate(filename string, cfg *config.Config) ([]pkg.Error, error) {
	var allErrors []pkg.Error

	if cfg == nil {
		cfg = config.DefaultConfig()
	}

	rules := cfg.Rules

	if rules.CheckSyntax {
		allErrors = append(allErrors, CheckSyntax(filename)...)
	}

	node, err := parser.ParseFile(filename)
	if err != nil {
		return allErrors, err
	}

	if rules.CheckDuplicates {
		allErrors = append(allErrors, CheckDuplicates(node)...)
	}

	if rules.CheckIntegrity {
		fields := rules.RequiredFields
		if len(fields) == 0 {
			fields = config.DefaultConfig().Rules.RequiredFields
		}
		allErrors = append(allErrors, CheckIntegrity(node, fields)...)
	}

	if rules.CheckCommonErrors {
		maxLen := rules.MaxLineLength
		if maxLen <= 0 {
			maxLen = 200
		}
		patterns := rules.SensitivePatterns
		if len(patterns) == 0 {
			patterns = config.DefaultConfig().Rules.SensitivePatterns
		}
		allErrors = append(allErrors, CheckCommonErrors(filename, maxLen, patterns)...)
	}

	// Запуск зарегистрированных плагинов
	allErrors = append(allErrors, RunPlugins(node)...)

	return allErrors, nil
}
