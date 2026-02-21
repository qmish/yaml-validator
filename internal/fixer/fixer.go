package fixer

import (
	"os"
	"strings"

	"yaml-validator/internal/config"
)

// FixResult результат автоисправления
type FixResult struct {
	Modified bool   // файл был изменён
	Applied  []string // применённые исправления (например "TrailingSpaces", "NewlineAtEof")
}

// FixFile применяет автоисправления к файлу согласно конфигу.
// Поддерживает: trailing spaces, newline at EOF, consecutive empty lines.
func FixFile(path string, cfg *config.Config) (FixResult, error) {
	var result FixResult
	if cfg == nil {
		cfg = config.DefaultConfig()
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return result, err
	}

	content := string(data)
	lines := strings.Split(content, "\n")

	// Нормализуем: Split по \n даёт N элементов для N-1 переводов.
	// Последний элемент может быть пустым, если файл заканчивается \n.
	hasTrailingNewline := len(data) > 0 && data[len(data)-1] == '\n'

	opts := cfg.Rules.Style
	modified := false

	// 1. Trailing spaces
	if opts.ForbidTrailingSpaces {
		for i, line := range lines {
			trimmed := strings.TrimRight(line, " \t")
			if len(trimmed) != len(line) {
				lines[i] = trimmed
				modified = true
				if !contains(result.Applied, "TrailingSpaces") {
					result.Applied = append(result.Applied, "TrailingSpaces")
				}
			}
		}
	}

	// 2. Consecutive empty lines
	if opts.ForbidConsecutiveEmptyLines {
		var newLines []string
		prevEmpty := false
		for _, line := range lines {
			empty := strings.TrimSpace(line) == ""
			if empty && prevEmpty {
				modified = true
				if !contains(result.Applied, "ConsecutiveEmptyLines") {
					result.Applied = append(result.Applied, "ConsecutiveEmptyLines")
				}
				continue
			}
			prevEmpty = empty
			newLines = append(newLines, line)
		}
		lines = newLines
	}

	if !modified && !(opts.RequireNewlineAtEof && len(data) > 0 && !hasTrailingNewline) {
		return result, nil
	}

	result.Modified = true
	newContent := strings.Join(lines, "\n")

	// 3. Newline at EOF
	if opts.RequireNewlineAtEof && len(newContent) > 0 && !strings.HasSuffix(newContent, "\n") {
		newContent += "\n"
		if !contains(result.Applied, "NewlineAtEof") {
			result.Applied = append(result.Applied, "NewlineAtEof")
		}
	}
	if err := os.WriteFile(path, []byte(newContent), 0644); err != nil {
		return result, err
	}
	return result, nil
}

func contains(s []string, x string) bool {
	for _, v := range s {
		if v == x {
			return true
		}
	}
	return false
}
