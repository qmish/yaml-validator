package validator

import (
	"os"
	"strings"

	"yaml-validator/internal/config"
	"yaml-validator/pkg"
)

// CheckStyle проверяет правила стиля (document-start, trailing spaces, newline at EOF)
func CheckStyle(filename string, opts config.StyleOptions) []pkg.Error {
	if !opts.RequireDocumentStart && !opts.ForbidTrailingSpaces && !opts.RequireNewlineAtEof {
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
				Type:    "DocumentStart",
				Message: "document should start with '---'",
				Line:    1,
			})
		}
	}

	if opts.ForbidTrailingSpaces {
		for i, line := range lines {
			if len(line) > 0 && (line[len(line)-1] == ' ' || line[len(line)-1] == '\t') {
				errors = append(errors, pkg.Error{
					Type:    "TrailingSpaces",
					Message: "trailing spaces at end of line",
					Line:    i + 1,
				})
			}
		}
	}

	if opts.RequireNewlineAtEof && len(data) > 0 {
		if data[len(data)-1] != '\n' {
			errors = append(errors, pkg.Error{
				Type:    "NewlineAtEof",
				Message: "missing newline at end of file",
				Line:    len(lines),
			})
		}
	}

	return errors
}
