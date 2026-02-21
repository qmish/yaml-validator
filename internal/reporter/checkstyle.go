package reporter

import (
	"fmt"
	"strings"
)

// GenerateCheckstyleReport создаёт отчёт в формате Checkstyle XML для Jenkins, SonarQube.
func GenerateCheckstyleReport(results []FileResult) ([]byte, error) {
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	sb.WriteString("\n<checkstyle version=\"yaml-validator\">\n")

	for _, r := range results {
		sb.WriteString(fmt.Sprintf("  <file name=\"%s\">\n", escapeXMLAttr(r.File)))
		for _, e := range r.Errors {
			line := e.Line
			if line <= 0 {
				line = 1
			}
			col := e.Column
			if col <= 0 {
				col = 1
			}
			sb.WriteString(fmt.Sprintf("    <error line=\"%d\" column=\"%d\" severity=\"error\" message=\"%s\" source=\"%s\"/>\n",
				line, col, escapeXMLAttr(e.Message), escapeXMLAttr(e.Type)))
		}
		sb.WriteString("  </file>\n")
	}

	sb.WriteString("</checkstyle>\n")
	return []byte(sb.String()), nil
}

func escapeXMLAttr(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}
