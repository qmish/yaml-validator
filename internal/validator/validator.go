package validator

import (
	"gopkg.in/yaml.v3"

	"yaml-validator/internal/config"
	"yaml-validator/internal/parser"
	"yaml-validator/pkg"
)

// Validate выполняет полную валидацию YAML-файла согласно конфигурации.
// Для мультидокументных файлов каждый документ валидируется отдельно, ошибки содержат DocumentIndex.
func Validate(filename string, cfg *config.Config) ([]pkg.Error, error) {
	var allErrors []pkg.Error

	if cfg == nil {
		cfg = config.DefaultConfig()
	}

	rules := cfg.Rules

	nodes, err := parser.ParseFileMulti(filename)
	if err != nil {
		if rules.CheckSyntax {
			allErrors = append(allErrors, CheckSyntax(filename)...)
		}
		return allErrors, err
	}

	// Файл-уровень: один раз
	if rules.CheckCommonErrors {
		maxLen := rules.MaxLineLength
		if maxLen <= 0 {
			maxLen = 200
		}
		patterns := rules.SensitivePatterns
		if len(patterns) == 0 {
			patterns = config.DefaultConfig().Rules.SensitivePatterns
		}
		allErrors = append(allErrors, CheckCommonErrors(filename, maxLen, patterns, rules.Style.ForbidTabs)...)
	}
	if rules.Style.RequireDocumentStart || rules.Style.ForbidTrailingSpaces || rules.Style.ForbidTrailingDots ||
		rules.Style.RequireNewlineAtEof || rules.Style.ForbidConsecutiveEmptyLines || rules.Style.RequireEmptyLineBetweenBlocks || rules.Style.MinEmptyLinesBetweenBlocks >= 1 || rules.Style.RequireDocumentEnd ||
		rules.Style.RequireCommentsIndented || rules.Style.RequireQuotedKeys || rules.Style.IndentSpaces > 0 || rules.Style.ForbidTabs || rules.Style.ForbidUnicode || rules.Style.ForbidBOM {
		allErrors = append(allErrors, CheckStyle(filename, rules.Style)...)
	}

	// Уровень документа: каждый документ отдельно
	for i, node := range nodes {
		idx := 0
		if len(nodes) > 1 {
			idx = i + 1
		}
		docErrs := validateNode(node, rules, cfg)
		for _, e := range docErrs {
			e.DocumentIndex = idx
			allErrors = append(allErrors, e)
		}
	}

	// JsonSchema и K8sSchema — читают файл, проверяют первый документ
	if rules.JsonSchema.Enabled && rules.JsonSchema.SchemaPath != "" {
		jsErrs := CheckJsonSchema(filename, rules.JsonSchema)
		for _, e := range jsErrs {
			if len(nodes) > 1 {
				e.DocumentIndex = 1
			}
			allErrors = append(allErrors, e)
		}
	}
	if rules.K8sSchema.Enabled {
		k8sErrs := CheckK8sSchema(filename, rules.K8sSchema)
		for _, e := range k8sErrs {
			if len(nodes) > 1 {
				e.DocumentIndex = 1
			}
			allErrors = append(allErrors, e)
		}
	}

	if rules.InlineIgnore {
		allErrors = FilterInlineIgnore(filename, allErrors)
	}

	for i := range allErrors {
		if allErrors[i].Severity != "" {
			continue
		}
		if sv, ok := rules.RuleSeverity[allErrors[i].Type]; ok && sv == "warning" {
			allErrors[i].Severity = "warning"
		} else {
			allErrors[i].Severity = "error"
		}
	}

	return allErrors, nil
}

// ValidateFromContent валидирует YAML по содержимому в памяти (для LSP). JsonSchema и K8sSchema не выполняются.
func ValidateFromContent(uri string, content []byte, cfg *config.Config) []pkg.Error {
	var allErrors []pkg.Error
	if cfg == nil {
		cfg = config.DefaultConfig()
	}
	rules := cfg.Rules

	nodes, err := parser.ParseBytesMulti(content)
	if err != nil {
		if rules.CheckSyntax {
			allErrors = append(allErrors, CheckSyntaxContent(content)...)
		}
		return allErrors
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
		allErrors = append(allErrors, CheckCommonErrorsContent(content, maxLen, patterns, rules.Style.ForbidTabs)...)
	}
	if rules.Style.RequireDocumentStart || rules.Style.ForbidTrailingSpaces || rules.Style.ForbidTrailingDots ||
		rules.Style.RequireNewlineAtEof || rules.Style.ForbidConsecutiveEmptyLines || rules.Style.RequireEmptyLineBetweenBlocks || rules.Style.MinEmptyLinesBetweenBlocks >= 1 || rules.Style.RequireDocumentEnd ||
		rules.Style.RequireCommentsIndented || rules.Style.RequireQuotedKeys || rules.Style.IndentSpaces > 0 || rules.Style.ForbidTabs || rules.Style.ForbidUnicode || rules.Style.ForbidBOM {
		allErrors = append(allErrors, CheckStyleContent(content, rules.Style)...)
	}

	for i, node := range nodes {
		idx := 0
		if len(nodes) > 1 {
			idx = i + 1
		}
		docErrs := validateNode(node, rules, cfg)
		for _, e := range docErrs {
			e.DocumentIndex = idx
			allErrors = append(allErrors, e)
		}
	}

	if rules.InlineIgnore {
		allErrors = FilterInlineIgnoreContent(content, allErrors)
	}

	for i := range allErrors {
		if allErrors[i].Severity != "" {
			continue
		}
		if sv, ok := rules.RuleSeverity[allErrors[i].Type]; ok && sv == "warning" {
			allErrors[i].Severity = "warning"
		} else {
			allErrors[i].Severity = "error"
		}
	}
	return allErrors
}

func validateNode(node *yaml.Node, rules config.ValidationRules, cfg *config.Config) []pkg.Error {
	var errs []pkg.Error
	if rules.CheckDuplicates {
		errs = append(errs, CheckDuplicates(node)...)
	}
	if len(rules.KeyOrder) > 0 {
		errs = append(errs, CheckKeyOrderConfigurable(node, rules.KeyOrder)...)
	} else if rules.CheckKeyOrdering {
		errs = append(errs, CheckKeyOrdering(node)...)
	}
	if rules.MaxKeyNameLength > 0 {
		errs = append(errs, CheckMaxKeyNameLength(node, rules.MaxKeyNameLength)...)
	}
	if rules.ForbidDefaultValues && len(rules.DefaultValues) > 0 {
		errs = append(errs, CheckForbidDefaultValues(node, rules.DefaultValues)...)
	}
	if len(rules.UniqueListFields) > 0 {
		errs = append(errs, CheckUniqueListFields(node, rules.UniqueListFields)...)
	}
	if len(rules.KeyValuePatterns) > 0 {
		errs = append(errs, CheckKeyValuePatterns(node, rules.KeyValuePatterns)...)
	}
	if rules.CheckIntegrity {
		fields := rules.RequiredFields
		if len(fields) == 0 {
			fields = config.DefaultConfig().Rules.RequiredFields
		}
		errs = append(errs, CheckIntegrity(node, fields)...)
	}
	if rules.Style.RequireQuotedValues {
		errs = append(errs, CheckQuotedValues(node)...)
	}
	errs = append(errs, RunPlugins(node)...)
	return errs
}
