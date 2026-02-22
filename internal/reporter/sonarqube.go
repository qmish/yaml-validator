package reporter

import (
	"encoding/json"
	"path/filepath"

	"yaml-validator/pkg"
)

const sonarEngineID = "yaml-validator"

// sonarReport — формат SonarQube Generic Issue Import
type sonarReport struct {
	Issues []sonarIssue `json:"issues"`
}

type sonarIssue struct {
	EngineID        string         `json:"engineId"`
	RuleID          string         `json:"ruleId"`
	Type            string         `json:"type"`
	Severity        string         `json:"severity"`
	PrimaryLocation sonarLocation  `json:"primaryLocation"`
}

type sonarLocation struct {
	Message  string       `json:"message"`
	FilePath string       `json:"filePath"`
	TextRange *sonarRange `json:"textRange,omitempty"`
}

type sonarRange struct {
	StartLine int `json:"startLine"`
	EndLine   int `json:"endLine,omitempty"`
}

// GenerateSonarQubeGenericReport создаёт отчёт в формате SonarQube Generic Issue Import.
// Пути относительно basePath (пустая строка = как есть).
func GenerateSonarQubeGenericReport(results []FileResult, basePath string) ([]byte, error) {
	var issues []sonarIssue
	ruleIDs := make(map[string]bool)

	for _, r := range results {
		path := r.File
		if basePath != "" {
			if rel, err := filepath.Rel(basePath, r.File); err == nil {
				path = filepath.ToSlash(rel)
			}
		}

		for _, e := range r.Errors {
			ruleIDs[e.Type] = true
			line := e.Line
			if line <= 0 {
				line = 1
			}
			issue := sonarIssue{
				EngineID: sonarEngineID,
				RuleID:   e.Type,
				Type:     sonarType(e.Type),
				Severity: sonarSeverity(e),
				PrimaryLocation: sonarLocation{
					Message:  e.Message,
					FilePath: path,
					TextRange: &sonarRange{StartLine: line, EndLine: line},
				},
			}
			issues = append(issues, issue)
		}
	}

	// SonarQube Generic требует rules — создаём из уникальных ruleId
	rules := make([]sonarRule, 0, len(ruleIDs))
	for id := range ruleIDs {
		rules = append(rules, sonarRule{
			ID:             id,
			Name:           id,
			Description:    id,
			EngineID:       sonarEngineID,
			Type:           sonarType(id),
			Severity:       sonarSeverityForType(id),
		})
	}

	report := struct {
		Issues []sonarIssue `json:"issues"`
		Rules  []sonarRule  `json:"rules"`
	}{
		Issues: issues,
		Rules:  rules,
	}
	return json.MarshalIndent(report, "", "  ")
}

type sonarRule struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	EngineID    string `json:"engineId"`
	Type        string `json:"type"`
	Severity    string `json:"severity"`
}

func sonarType(t string) string {
	switch t {
	case "SyntaxError", "DuplicateKeys", "EmptyFile", "ReadError":
		return "BUG"
	case "SensitiveData":
		return "VULNERABILITY"
	default:
		return "CODE_SMELL"
	}
}

func sonarSeverity(e pkg.Error) string {
	if e.Severity == "warning" {
		return "MINOR"
	}
	return sonarSeverityForType(e.Type)
}

func sonarSeverityForType(t string) string {
	switch t {
	case "SyntaxError", "EmptyFile", "ReadError", "SensitiveData":
		return "BLOCKER"
	case "DuplicateKeys", "MissingRequiredField":
		return "MAJOR"
	default:
		return "MAJOR"
	}
}
