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

	if len(rules.KeyOrder) > 0 {
		allErrors = append(allErrors, CheckKeyOrderConfigurable(node, rules.KeyOrder)...)
	} else if rules.CheckKeyOrdering {
		allErrors = append(allErrors, CheckKeyOrdering(node)...)
	}

	if rules.MaxKeyNameLength > 0 {
		allErrors = append(allErrors, CheckMaxKeyNameLength(node, rules.MaxKeyNameLength)...)
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

	if rules.Style.RequireDocumentStart || rules.Style.ForbidTrailingSpaces || rules.Style.ForbidTrailingDots ||
		rules.Style.RequireNewlineAtEof || rules.Style.ForbidConsecutiveEmptyLines || rules.Style.RequireEmptyLineBetweenBlocks || rules.Style.RequireDocumentEnd ||
		rules.Style.RequireCommentsIndented || rules.Style.RequireQuotedKeys || rules.Style.IndentSpaces > 0 {
		allErrors = append(allErrors, CheckStyle(filename, rules.Style)...)
	}

	if rules.Style.RequireQuotedValues {
		allErrors = append(allErrors, CheckQuotedValues(node)...)
	}

	if rules.JsonSchema.Enabled && rules.JsonSchema.SchemaPath != "" {
		allErrors = append(allErrors, CheckJsonSchema(filename, rules.JsonSchema)...)
	}

	if rules.K8sSchema.Enabled {
		allErrors = append(allErrors, CheckK8sSchema(filename, rules.K8sSchema)...)
	}

	// Запуск зарегистрированных плагинов
	allErrors = append(allErrors, RunPlugins(node)...)

	if rules.InlineIgnore {
		allErrors = FilterInlineIgnore(filename, allErrors)
	}

	return allErrors, nil
}
