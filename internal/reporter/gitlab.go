package reporter

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"path/filepath"
)

// gitlabCodeQualityEntry — одна запись для GitLab Code Quality report
type gitlabCodeQualityEntry struct {
	Description  string              `json:"description"`
	Fingerprint  string              `json:"fingerprint"`
	Severity     string              `json:"severity"`
	Location     gitlabLocation      `json:"location"`
	CheckName    string              `json:"check_name"`
	Categories   []string            `json:"categories"`
	EngineName   string              `json:"engine_name"`
	Type         string              `json:"type"`
	Content      *gitlabContent      `json:"content,omitempty"`
}

type gitlabLocation struct {
	Path  string       `json:"path"`
	Lines *gitlabLines `json:"lines"`
}

type gitlabLines struct {
	Begin int `json:"begin"`
	End   int `json:"end"`
}

type gitlabContent struct {
	Body string `json:"body"`
}

// GenerateGitLabCodeQualityReport создаёт отчёт в формате GitLab Code Quality (gl-code-quality-report.json)
func GenerateGitLabCodeQualityReport(results []FileResult) ([]byte, error) {
	var entries []gitlabCodeQualityEntry

	for _, r := range results {
		relPath := r.File
		if abs, err := filepath.Abs(r.File); err == nil {
			relPath = filepath.Base(abs)
		}

		for _, e := range r.Errors {
			line := e.Line
			if line <= 0 {
				line = 1
			}

			fp := fingerprint(r.File, line, e.Type, e.Message)

			entries = append(entries, gitlabCodeQualityEntry{
				Description: e.Message,
				Fingerprint: fp,
				Severity:    severityForType(e.Type),
				Location: gitlabLocation{
					Path:  relPath,
					Lines: &gitlabLines{Begin: line, End: line},
				},
				CheckName:  e.Type,
				Categories: []string{"Style"},
				EngineName: "yaml-validator",
				Type:       "issue",
				Content:    &gitlabContent{Body: e.Message},
			})
		}
	}

	return json.MarshalIndent(entries, "", "  ")
}

func fingerprint(file string, line int, rule, msg string) string {
	h := sha256.New()
	h.Write([]byte(file))
	h.Write([]byte{0})
	h.Write([]byte{byte(line >> 24), byte(line >> 16), byte(line >> 8), byte(line)})
	h.Write([]byte(rule))
	h.Write([]byte(msg))
	return hex.EncodeToString(h.Sum(nil))
}

func severityForType(t string) string {
	switch t {
	case "SyntaxError", "EmptyFile", "ReadError", "SensitiveData":
		return "critical"
	default:
		return "major"
	}
}
