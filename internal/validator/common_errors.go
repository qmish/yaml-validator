package validator

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"yaml-validator/pkg"
)

// CheckCommonErrors ищет распространённые ошибки: табуляции (если skipTabCheck false), длинные строки, чувствительные данные.
// При skipTabCheck true проверка табов выполняется в CheckStyle (style.forbid_tabs).
func CheckCommonErrors(filename string, maxLineLength int, sensitivePatterns []string, skipTabCheck bool) []pkg.Error {
	var errors []pkg.Error

	data, err := os.ReadFile(filename)
	if err != nil {
		return errors
	}

	lines := strings.Split(string(data), "\n")

	for i, line := range lines {
		lineNum := i + 1

		if !skipTabCheck && strings.Contains(line, "\t") {
			errors = append(errors, pkg.Error{
				Type:       "TabInsteadOfSpaces",
				Message:    "Line contains tabs; use spaces for indentation",
				Suggestion: "replace tabs with spaces",
				Line:       lineNum,
			})
		}

		if maxLineLength > 0 && len(line) > maxLineLength {
			col := maxLineLength + 1 // первая позиция, превысившая лимит
			errors = append(errors, pkg.Error{
				Type:       "LineTooLong",
				Message:    fmt.Sprintf("Line exceeds max length of %d characters (%d)", maxLineLength, len(line)),
				Suggestion: "shorten line or split across multiple lines",
				Line:       lineNum,
				Column:     col,
			})
		}

		for _, pattern := range sensitivePatterns {
			keyPattern := regexp.MustCompile(`(?m)^\s*` + regexp.QuoteMeta(pattern) + `\s*:`)
			if keyPattern.MatchString(strings.ToLower(line)) {
				errors = append(errors, pkg.Error{
					Type:    "SensitiveData",
					Message: fmt.Sprintf("Possible sensitive data: field matching '%s' detected", pattern),
					Line:    lineNum,
				})
			}
		}
	}

	return errors
}
