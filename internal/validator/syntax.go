package validator

import (
	"os"

	"yaml-validator/pkg"

	"gopkg.in/yaml.v3"
)

// CheckSyntax проверяет синтаксис YAML: пустой документ и ошибки парсинга
func CheckSyntax(filename string) []pkg.Error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return []pkg.Error{{Type: "ReadError", Message: err.Error()}}
	}
	return CheckSyntaxContent(data)
}

// CheckSyntaxContent проверяет синтаксис по содержимому (для LSP)
func CheckSyntaxContent(data []byte) []pkg.Error {
	var errors []pkg.Error
	if len(data) == 0 {
		errors = append(errors, pkg.Error{
			Type:    "EmptyFile",
			Message: "File is empty",
		})
		return errors
	}

	var node yaml.Node
	err := yaml.Unmarshal(data, &node)
	if err != nil {
		return []pkg.Error{{Type: "SyntaxError", Message: err.Error()}}
	}

	if node.Kind == yaml.DocumentNode && len(node.Content) == 0 {
		errors = append(errors, pkg.Error{
			Type:    "EmptyDocument",
			Message: "YAML document is empty",
		})
	}

	return errors
}
