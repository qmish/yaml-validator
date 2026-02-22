package reporter

import (
	"encoding/json"
	"path/filepath"

	"yaml-validator/pkg"
)

// codeclimateIssue — формат issue для Code Climate Engine (SPEC.md)
type codeclimateIssue struct {
	Type              string           `json:"type"`
	CheckName         string           `json:"check_name"`
	Description       string           `json:"description"`
	Categories        []string         `json:"categories"`
	Location          codeclimateLoc   `json:"location"`
	RemediationPoints int              `json:"remediation_points,omitempty"`
	Severity          string           `json:"severity,omitempty"`
	Fingerprint       string           `json:"fingerprint,omitempty"`
	Content           *codeclimateBody `json:"content,omitempty"`
}

type codeclimateLoc struct {
	Path string            `json:"path"`
	Lines *codeclimateLines `json:"lines,omitempty"`
}

type codeclimateLines struct {
	Begin int `json:"begin"`
	End   int `json:"end"`
}

type codeclimateBody struct {
	Body string `json:"body"`
}

// GenerateCodeClimateReport создаёт отчёт в формате Code Climate (NDJSON — по одному issue на строку).
// Каждая строка — валидный JSON-объект. Пути относительно basePath (пустая строка = как есть).
func GenerateCodeClimateReport(results []FileResult, basePath string) ([]byte, error) {
	var out []byte
	for _, r := range results {
		path := r.File
		if basePath != "" {
			if rel, err := filepath.Rel(basePath, r.File); err == nil {
				path = filepath.ToSlash(rel)
			}
		}

		for _, e := range r.Errors {
			line := e.Line
			if line <= 0 {
				line = 1
			}
			issue := codeclimateIssue{
				Type:        "issue",
				CheckName:   e.Type,
				Description: e.Message,
				Categories:  []string{categoryForType(e.Type)},
				Location: codeclimateLoc{
					Path:  path,
					Lines: &codeclimateLines{Begin: line, End: line},
				},
				RemediationPoints: 50000,
				Severity:          severityCodeClimateForType(e),
				Fingerprint:       fingerprint(r.File, line, e.Type, e.Message),
			}
			if e.Suggestion != "" {
				issue.Content = &codeclimateBody{Body: "To fix: " + e.Suggestion}
			}
			b, err := json.Marshal(issue)
			if err != nil {
				return nil, err
			}
			out = append(out, b...)
			out = append(out, '\n')
		}
	}
	return out, nil
}

func categoryForType(t string) string {
	switch t {
	case "SyntaxError", "DuplicateKeys", "EmptyFile", "ReadError":
		return "Bug Risk"
	case "SensitiveData":
		return "Security"
	case "LineTooLong", "TrailingSpaces", "TrailingDots", "RequireNewlineAtEof",
		"ForbidConsecutiveEmptyLines", "RequireCommentsIndented", "RequireQuotedKeys",
		"RequireQuotedValues", "IndentSpaces", "ForbidTabs", "ForbidUnicode", "ForbidBOM":
		return "Style"
	case "MissingRequiredField", "CheckKeyOrdering", "KeyOrder", "ForbidDefaultValues",
		"UniqueListFields", "KeyValuePatterns":
		return "Clarity"
	default:
		return "Style"
	}
}

func severityCodeClimateForType(e pkg.Error) string {
	if e.Severity == "warning" {
		return "minor"
	}
	switch e.Type {
	case "SyntaxError", "EmptyFile", "ReadError", "SensitiveData":
		return "critical"
	case "DuplicateKeys", "MissingRequiredField":
		return "major"
	default:
		return "major"
	}
}

