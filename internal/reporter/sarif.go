package reporter

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"yaml-validator/pkg"
)

// FileResult результат валидации одного файла для SARIF
type FileResult struct {
	File   string
	Errors []pkg.Error
}

// sarifLog минимальная структура SARIF 2.1.0 для GitHub Code Scanning
type sarifLog struct {
	Schema  string   `json:"$schema"`
	Version string   `json:"version"`
	Runs    []sarifRun `json:"runs"`
}

type sarifRun struct {
	Tool    sarifTool `json:"tool"`
	Artifacts []sarifArtifact `json:"artifacts"`
	Results   []sarifResult   `json:"results"`
}

type sarifTool struct {
	Driver sarifDriver `json:"driver"`
}

type sarifDriver struct {
	Name           string              `json:"name"`
	Version        string              `json:"version"`
	InformationURI string              `json:"informationUri,omitempty"`
	Rules          []sarifRule         `json:"rules"`
}

type sarifRule struct {
	ID               string        `json:"id"`
	ShortDescription sarifMessage `json:"shortDescription"`
}

type sarifMessage struct {
	Text string `json:"text"`
}

type sarifArtifact struct {
	Location sarifLocation `json:"location"`
}

type sarifLocation struct {
	URI string `json:"uri"`
}

type sarifResult struct {
	RuleID   string          `json:"ruleId"`
	Level    string          `json:"level,omitempty"`
	Message  sarifMessage    `json:"message"`
	Locations []sarifResultLocation `json:"locations,omitempty"`
}

type sarifResultLocation struct {
	PhysicalLocation sarifPhysicalLocation `json:"physicalLocation"`
}

type sarifPhysicalLocation struct {
	ArtifactLocation sarifArtifactRef `json:"artifactLocation"`
	Region           *sarifRegion     `json:"region,omitempty"`
}

type sarifArtifactRef struct {
	Index int    `json:"index,omitempty"`
	URI   string `json:"uri,omitempty"`
}

type sarifRegion struct {
	StartLine   int `json:"startLine,omitempty"`
	StartColumn int `json:"startColumn,omitempty"`
}

// GenerateSARIFReport создаёт отчёт в формате SARIF 2.1.0 для GitHub Code Scanning.
// fileResults — список (файл, ошибки); toolVersion — версия инструмента.
func GenerateSARIFReport(toolVersion string, fileResults []FileResult) ([]byte, error) {
	if toolVersion == "" {
		toolVersion = "1.0.0"
	}

	ruleIDs := make(map[string]bool)
	var artifacts []sarifArtifact
	var results []sarifResult
	artifactIndexByFile := make(map[string]int)

	for i, fr := range fileResults {
		uri := toFileURI(fr.File)
		artifactIndexByFile[fr.File] = i
		artifacts = append(artifacts, sarifArtifact{Location: sarifLocation{URI: uri}})

		for _, e := range fr.Errors {
			ruleIDs[e.Type] = true
			r := sarifResult{
				RuleID:  e.Type,
				Level:   "error",
				Message: sarifMessage{Text: e.Message},
				Locations: []sarifResultLocation{{
					PhysicalLocation: sarifPhysicalLocation{
						ArtifactLocation: sarifArtifactRef{Index: i},
						Region:           regionFromError(e),
					},
				}},
			}
			results = append(results, r)
		}
	}

	var rules []sarifRule
	for id := range ruleIDs {
		rules = append(rules, sarifRule{
			ID:               id,
			ShortDescription: sarifMessage{Text: id},
		})
	}

	log := sarifLog{
		Schema:  "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
		Version: "2.1.0",
		Runs: []sarifRun{{
			Tool: sarifTool{
				Driver: sarifDriver{
					Name:           "yaml-validator",
					Version:        toolVersion,
					InformationURI: "https://github.com/qmish/yaml-validator",
					Rules:          rules,
				},
			},
			Artifacts: artifacts,
			Results:   results,
		}},
	}

	return json.MarshalIndent(log, "", "  ")
}

func toFileURI(path string) string {
	path = filepath.ToSlash(path)
	if !strings.HasPrefix(path, "/") && len(path) > 1 && path[1] != ':' {
		path = "/" + path
	}
	return "file://" + path
}

func regionFromError(e pkg.Error) *sarifRegion {
	if e.Line <= 0 {
		return nil
	}
	col := e.Column
	if col <= 0 {
		col = 1
	}
	return &sarifRegion{StartLine: e.Line, StartColumn: col}
}
