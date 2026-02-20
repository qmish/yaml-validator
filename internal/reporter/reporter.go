package reporter

import (
	"encoding/json"
	"fmt"
	"strings"

	"yaml-validator/pkg"
)

// Report содержит результат валидации
type Report struct {
	File   string      `json:"file"`
	Valid  bool        `json:"valid"`
	Errors []pkg.Error `json:"errors"`
}

// Format определяет формат вывода отчёта
type Format string

const (
	FormatJSON          Format = "json"
	FormatJUnit         Format = "junit"
	FormatHumanReadable Format = "human"
)

// GenerateJSONReport создаёт JSON-отчёт для интеграции с CI/CD
func GenerateJSONReport(file string, errors []pkg.Error) ([]byte, error) {
	report := Report{
		File:   file,
		Valid:  len(errors) == 0,
		Errors: errors,
	}
	return json.MarshalIndent(report, "", "  ")
}

// GenerateJUnitReport формирует отчёт в формате JUnit для Jenkins, GitLab CI
func GenerateJUnitReport(file string, errors []pkg.Error) ([]byte, error) {
	var sb strings.Builder

	valid := len(errors) == 0
	status := "success"
	if !valid {
		status = "failure"
	}

	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	sb.WriteString(fmt.Sprintf("\n<testsuites tests=\"1\" failures=\"%d\" errors=\"0\">\n", len(errors)))
	sb.WriteString(fmt.Sprintf("  <testsuite name=\"yaml-validator\" tests=\"1\" failures=\"%d\" errors=\"0\">\n", len(errors)))
	sb.WriteString(fmt.Sprintf("    <testcase name=\"%s\" status=\"%s\">\n", escapeXML(file), status))

	for _, e := range errors {
		sb.WriteString("      <failure type=\"" + escapeXML(e.Type) + "\" message=\"" + escapeXML(e.Message) + "\"/>\n")
	}

	sb.WriteString("    </testcase>\n")
	sb.WriteString("  </testsuite>\n")
	sb.WriteString("</testsuites>\n")

	return []byte(sb.String()), nil
}

func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}

// PrintHumanReadable выводит результаты в консоль в человекочитаемом формате
func PrintHumanReadable(file string, errors []pkg.Error) {
	if len(errors) == 0 {
		fmt.Printf("✓ %s: valid\n", file)
		return
	}

	fmt.Printf("✗ %s: %d error(s)\n", file, len(errors))
	for i, e := range errors {
		if e.Line > 0 {
			fmt.Printf("  %d. [%s] %s (line %d)\n", i+1, e.Type, e.Message, e.Line)
		} else if e.Path != "" {
			fmt.Printf("  %d. [%s] %s (path: %s)\n", i+1, e.Type, e.Message, e.Path)
		} else {
			fmt.Printf("  %d. [%s] %s\n", i+1, e.Type, e.Message)
		}
	}
}
