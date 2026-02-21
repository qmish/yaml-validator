package validator

import (
	"fmt"
	"os"
	"strings"

	"yaml-validator/internal/config"
	"yaml-validator/pkg"
)

// CheckStyle проверяет правила стиля (document-start/end, trailing spaces/dots, newline at EOF, consecutive empty lines, comment indentation, quoted keys, indent step)
func CheckStyle(filename string, opts config.StyleOptions) []pkg.Error {
	if !opts.RequireDocumentStart && !opts.ForbidTrailingSpaces && !opts.ForbidTrailingDots && !opts.RequireNewlineAtEof &&
		!opts.ForbidConsecutiveEmptyLines && !opts.RequireEmptyLineBetweenBlocks && opts.MinEmptyLinesBetweenBlocks <= 0 && !opts.RequireDocumentEnd && !opts.RequireCommentsIndented && !opts.RequireQuotedKeys && opts.IndentSpaces <= 0 && !opts.ForbidTabs {
		return nil
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil
	}

	content := string(data)
	lines := strings.Split(content, "\n")

	var errors []pkg.Error

	if opts.RequireDocumentStart && len(lines) > 0 {
		firstLine := strings.TrimSpace(lines[0])
		if firstLine != "" && firstLine != "---" {
			errors = append(errors, pkg.Error{
				Type:       "DocumentStart",
				Message:    "document should start with '---'",
				Suggestion: "add '---' at the start of the file",
				Line:       1,
			})
		}
	}

	if opts.ForbidTrailingSpaces {
		for i, line := range lines {
			if len(line) > 0 && (line[len(line)-1] == ' ' || line[len(line)-1] == '\t') {
				errors = append(errors, pkg.Error{
					Type:       "TrailingSpaces",
					Message:    "trailing spaces at end of line",
					Suggestion: "remove trailing spaces",
					Line:       i + 1,
				})
			}
		}
	}

	if opts.ForbidTrailingDots {
		for i, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed != "" && trimmed != "..." && strings.HasSuffix(trimmed, ".") {
				errors = append(errors, pkg.Error{
					Type:       "TrailingDots",
					Message:    "forbid trailing dots at end of line",
					Suggestion: "remove trailing dot",
					Line:       i + 1,
				})
			}
		}
	}

	if opts.ForbidTabs {
		for i, line := range lines {
			if strings.Contains(line, "\t") {
				errors = append(errors, pkg.Error{
					Type:       "TabInsteadOfSpaces",
					Message:    "line contains tabs; use spaces for indentation",
					Suggestion: "replace tabs with spaces",
					Line:       i + 1,
				})
			}
		}
	}

	if opts.RequireNewlineAtEof && len(data) > 0 {
		if data[len(data)-1] != '\n' {
			errors = append(errors, pkg.Error{
				Type:       "NewlineAtEof",
				Message:    "missing newline at end of file",
				Suggestion: "add newline at end of file",
				Line:       len(lines),
			})
		}
	}

	if opts.ForbidConsecutiveEmptyLines {
		emptyCount := 0
		for i, line := range lines {
			if strings.TrimSpace(line) == "" {
				emptyCount++
				if emptyCount > 1 {
					errors = append(errors, pkg.Error{
						Type:       "ConsecutiveEmptyLines",
						Message:    "more than one consecutive empty line",
						Suggestion: "leave at most one empty line",
						Line:       i + 1,
					})
				}
			} else {
				emptyCount = 0
			}
		}
	}

	if opts.RequireEmptyLineBetweenBlocks || opts.MinEmptyLinesBetweenBlocks >= 1 {
		var lastTopLevelKeyLine int
		for i, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" || strings.HasPrefix(trimmed, "#") || trimmed == "---" || trimmed == "..." {
				continue
			}
			indent := len(line) - len(strings.TrimLeft(line, " \t"))
			if indent != 0 {
				continue
			}
			keyPart := trimmed
			if idx := strings.Index(trimmed, "#"); idx >= 0 {
				keyPart = strings.TrimSpace(trimmed[:idx])
			}
			if idx := strings.Index(keyPart, ":"); idx > 0 {
				if lastTopLevelKeyLine > 0 {
					emptyCount := 0
					for j := lastTopLevelKeyLine; j < i; j++ {
						if strings.TrimSpace(lines[j]) == "" {
							emptyCount++
						}
					}
					if emptyCount < 1 {
						msg := "at least one empty line required between top-level blocks"
						if opts.RequireEmptyLineBetweenBlocks {
							msg = "exactly one empty line required between top-level blocks"
						}
						errors = append(errors, pkg.Error{
							Type:       "EmptyLineBetweenBlocks",
							Message:    msg,
							Suggestion: "insert one empty line between blocks",
							Line:       i + 1,
						})
					}
				}
				lastTopLevelKeyLine = i + 1
			}
		}
	}

	if opts.RequireDocumentEnd && len(lines) > 0 {
		lastNonEmptyLine := 0
		for i := len(lines) - 1; i >= 0; i-- {
			if strings.TrimSpace(lines[i]) != "" {
				lastNonEmptyLine = i + 1
				if strings.TrimSpace(lines[i]) != "..." {
					errors = append(errors, pkg.Error{
						Type:       "DocumentEnd",
						Message:    "document should end with '...'",
						Suggestion: "add '...' at end of document",
						Line:       lastNonEmptyLine,
					})
				}
				break
			}
		}
	}

		if opts.RequireCommentsIndented {
		lastIndent := 0
		for i, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" {
				continue
			}
			if strings.HasPrefix(trimmed, "#") {
				commentIndent := len(line) - len(strings.TrimLeft(line, " \t"))
				if lastIndent > 0 && commentIndent == 0 {
					errors = append(errors, pkg.Error{
						Type:       "CommentIndentation",
						Message:    "comment should be indented to match the block",
						Suggestion: "indent comment to match block",
						Line:       i + 1,
					})
				}
				continue
			}
			lastIndent = len(line) - len(strings.TrimLeft(line, " \t"))
		}
	}

	if opts.IndentSpaces == 2 || opts.IndentSpaces == 4 {
		for i, line := range lines {
			if strings.TrimSpace(line) == "" {
				continue
			}
			if strings.Contains(line, "\t") {
				continue // TabInsteadOfSpaces в style.forbid_tabs или check_common_errors
			}
			indent := len(line) - len(strings.TrimLeft(line, " "))
			if indent > 0 && indent%opts.IndentSpaces != 0 {
				errors = append(errors, pkg.Error{
					Type:       "IndentSpaces",
					Message:    fmt.Sprintf("indent should be multiple of %d spaces, got %d", opts.IndentSpaces, indent),
					Suggestion: fmt.Sprintf("use %d-space indent", opts.IndentSpaces),
					Line:       i + 1,
				})
			}
		}
	}

	if opts.RequireQuotedKeys {
		for i, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" || strings.HasPrefix(trimmed, "#") || trimmed == "---" || trimmed == "..." {
				continue
			}
			if strings.HasPrefix(strings.TrimLeft(line, " \t"), "- ") {
				continue
			}
			idx := strings.Index(line, ":")
			if idx < 0 {
				continue
			}
			keyPart := strings.TrimSpace(line[:idx])
			if keyPart == "" {
				continue
			}
			if keyPart[0] != '"' && keyPart[0] != '\'' {
				errors = append(errors, pkg.Error{
					Type:       "QuotedKeys",
					Message:    "mapping key should be quoted",
					Suggestion: "wrap key in double or single quotes",
					Line:       i + 1,
				})
			}
		}
	}

	return errors
}
