package validator

import (
	"os"

	"yaml-validator/pkg"

	"gopkg.in/yaml.v3"
)

// CheckSyntax проверяет синтаксис YAML: пустой документ и ошибки парсинга
func CheckSyntax(filename string) []pkg.Error {
	var errors []pkg.Error

	data, err := os.ReadFile(filename)
	if err != nil {
		errors = append(errors, pkg.Error{
			Type:    "ReadError",
			Message: err.Error(),
		})
		return errors
	}

	if len(data) == 0 {
		errors = append(errors, pkg.Error{
			Type:    "EmptyFile",
			Message: "File is empty",
		})
		return errors
	}

	var node yaml.Node
	err = yaml.Unmarshal(data, &node)
	if err != nil {
		errors = append(errors, pkg.Error{
			Type:    "SyntaxError",
			Message: err.Error(),
		})
		return errors
	}

	if node.Kind == yaml.DocumentNode && len(node.Content) == 0 {
		errors = append(errors, pkg.Error{
			Type:    "EmptyDocument",
			Message: "YAML document is empty",
		})
	}

	return errors
}
