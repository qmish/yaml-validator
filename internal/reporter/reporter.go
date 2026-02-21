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
	FormatSARIF         Format = "sarif"
	FormatCheckstyle    Format = "checkstyle" // XML для Jenkins, SonarQube
	FormatHumanReadable Format = "human"
	FormatCompact       Format = "compact" // file:line: message (ESLint-style, для редакторов)
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
	nErr, nWarn := 0, 0
	for _, e := range errors {
		if e.Severity == "warning" {
			nWarn++
		} else {
			nErr++
		}
	}
	if nWarn > 0 && nErr > 0 {
		fmt.Printf("✗ %s: %d error(s), %d warning(s)\n", file, nErr, nWarn)
	} else if nWarn > 0 {
		fmt.Printf("⚠ %s: %d warning(s)\n", file, nWarn)
	} else {
		fmt.Printf("✗ %s: %d error(s)\n", file, nErr)
	}
	for i, e := range errors {
		sev := ""
		if e.Severity == "warning" {
			sev = " [WARN]"
		}
		if e.Line > 0 {
			fmt.Printf("  %d. [%s]%s %s (line %d)\n", i+1, e.Type, sev, e.Message, e.Line)
		} else if e.Path != "" {
			fmt.Printf("  %d. [%s]%s %s (path: %s)\n", i+1, e.Type, sev, e.Message, e.Path)
		} else {
			fmt.Printf("  %d. [%s]%s %s\n", i+1, e.Type, sev, e.Message)
		}
		if e.Suggestion != "" {
			fmt.Printf("      To fix: %s\n", e.Suggestion)
		}
	}
}

// GenerateCompactReport возвращает вывод в формате file:line[:col]: message (ESLint-style).
// Для мультидокумента (DocumentIndex > 0) используется file#docN.
func GenerateCompactReport(file string, errors []pkg.Error) string {
	var sb strings.Builder
	for _, e := range errors {
		displayFile := file
		if e.DocumentIndex > 0 {
			displayFile = fmt.Sprintf("%s#doc%d", file, e.DocumentIndex)
		}
		if e.Line > 0 {
			if e.Column > 0 {
				sb.WriteString(fmt.Sprintf("%s:%d:%d: %s\n", displayFile, e.Line, e.Column, e.Message))
			} else {
				sb.WriteString(fmt.Sprintf("%s:%d: %s\n", displayFile, e.Line, e.Message))
			}
		} else {
			sb.WriteString(fmt.Sprintf("%s: %s\n", displayFile, e.Message))
		}
	}
	return sb.String()
}

// PrintCompact выводит по одной строке на ошибку в формате file:line: message (ESLint-style).
func PrintCompact(file string, errors []pkg.Error) {
	fmt.Print(GenerateCompactReport(file, errors))
}

// GenerateGitHubAnnotations возвращает вывод в формате GitHub Actions ::error/::warning file=...,line=...::
func GenerateGitHubAnnotations(file string, errors []pkg.Error) string {
	var sb strings.Builder
	for _, e := range errors {
		displayFile := file
		if e.DocumentIndex > 0 {
			displayFile = fmt.Sprintf("%s#doc%d", file, e.DocumentIndex)
		}
		msg := strings.ReplaceAll(e.Message, "%", "%25")
		level := "error"
		if e.Severity == "warning" {
			level = "warning"
		}
		if e.Line > 0 {
			sb.WriteString(fmt.Sprintf("::%s file=%s,line=%d::%s\n", level, displayFile, e.Line, msg))
		} else {
			sb.WriteString(fmt.Sprintf("::%s file=%s::%s\n", level, displayFile, msg))
		}
	}
	return sb.String()
}

// PrintGitHubAnnotations выводит ошибки в формате GitHub Actions ::error file=...,line=...::
func PrintGitHubAnnotations(file string, errors []pkg.Error) {
	fmt.Print(GenerateGitHubAnnotations(file, errors))
}

// GenerateSeverityReport возвращает вывод [ERROR] file:line: message / [WARN] для скриптов и CI.
func GenerateSeverityReport(file string, errors []pkg.Error) string {
	var sb strings.Builder
	for _, e := range errors {
		displayFile := file
		if e.DocumentIndex > 0 {
			displayFile = fmt.Sprintf("%s#doc%d", file, e.DocumentIndex)
		}
		sev := "[ERROR]"
		if e.Severity == "warning" {
			sev = "[WARN]"
		}
		if e.Line > 0 {
			if e.Column > 0 {
				sb.WriteString(fmt.Sprintf("%s %s:%d:%d: %s\n", sev, displayFile, e.Line, e.Column, e.Message))
			} else {
				sb.WriteString(fmt.Sprintf("%s %s:%d: %s\n", sev, displayFile, e.Line, e.Message))
			}
		} else {
			sb.WriteString(fmt.Sprintf("%s %s: %s\n", sev, displayFile, e.Message))
		}
	}
	return sb.String()
}

// PrintSeverity выводит ошибки в формате [ERROR] file:line: message.
func PrintSeverity(file string, errors []pkg.Error) {
	fmt.Print(GenerateSeverityReport(file, errors))
}
