package validator

import (
	"os"
	"regexp"
	"strings"

	"yaml-validator/pkg"
)

// InlineIgnoreRule префикс комментария для отключения правил (в духе yamllint)
const InlineIgnoreRule = "yaml-validator"

// parsedIgnore: для строки N отключены правила; пустой set = отключить все для этой строки
type parsedIgnore map[int]map[string]bool

var (
	reDisableLine = regexp.MustCompile(`#\s*` + InlineIgnoreRule + `\s+disable-line`)
	reDisableNext = regexp.MustCompile(`#\s*` + InlineIgnoreRule + `\s+disable-next-line`)
	reRule        = regexp.MustCompile(`rule:(\S+)`)
)

// parseInlineIgnore читает файл и собирает по строкам отключённые правила.
func parseInlineIgnore(filename string) (parsedIgnore, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return parseInlineIgnoreContent(data), nil
}

// parseInlineIgnoreContent то же по содержимому (для LSP). Формат: # yaml-validator disable-line [rule:RuleType]
func parseInlineIgnoreContent(data []byte) parsedIgnore {
	lines := strings.Split(string(data), "\n")
	out := make(parsedIgnore)

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		// disable-line — эта строка
		if reDisableLine.MatchString(trimmed) {
			rules := extractRules(trimmed)
			lineNum := i + 1
			if _, ok := out[lineNum]; !ok {
				out[lineNum] = rules
			}
		}

		// disable-next-line — следующая строка
		if reDisableNext.MatchString(trimmed) {
			rules := extractRules(trimmed)
			lineNum := i + 2
			if lineNum <= len(lines) {
				if _, ok := out[lineNum]; !ok {
					out[lineNum] = rules
				}
			}
		}
	}
	return out
}

func extractRules(line string) map[string]bool {
	matches := reRule.FindAllStringSubmatch(line, -1)
	if len(matches) == 0 {
		return nil // nil = disable all rules for this line
	}
	rules := make(map[string]bool)
	for _, m := range matches {
		if len(m) > 1 && m[1] != "" {
			rules[m[1]] = true
		}
	}
	return rules
}

// isDisabled возвращает true, если ошибка на строке line с типом ruleType отключена через комментарий.
func (p parsedIgnore) isDisabled(line int, ruleType string) bool {
	set, ok := p[line]
	if !ok {
		return false
	}
	if set == nil || len(set) == 0 {
		return true
	}
	return set[ruleType]
}

// FilterInlineIgnore удаляет из списка ошибок те, что отключены комментариями в файле.
func FilterInlineIgnore(filename string, errors []pkg.Error) []pkg.Error {
	ignore, err := parseInlineIgnore(filename)
	if err != nil || len(ignore) == 0 {
		return errors
	}
	return filterInlineIgnoreWith(ignore, errors)
}

// FilterInlineIgnoreContent то же по содержимому (для LSP)
func FilterInlineIgnoreContent(data []byte, errors []pkg.Error) []pkg.Error {
	ignore := parseInlineIgnoreContent(data)
	if len(ignore) == 0 {
		return errors
	}
	return filterInlineIgnoreWith(ignore, errors)
}

func filterInlineIgnoreWith(ignore parsedIgnore, errors []pkg.Error) []pkg.Error {

	var out []pkg.Error
	for _, e := range errors {
		if e.Line <= 0 || !ignore.isDisabled(e.Line, e.Type) {
			out = append(out, e)
		}
	}
	return out
}
