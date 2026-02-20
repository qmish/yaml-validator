package validator

import (
	"fmt"

	"yaml-validator/internal/parser"
	"yaml-validator/pkg"

	"gopkg.in/yaml.v3"
)

// CheckIntegrity проверяет логическую целостность: обязательные поля и типы данных
func CheckIntegrity(node *yaml.Node, requiredFields []string) []pkg.Error {
	var errors []pkg.Error

	if len(requiredFields) == 0 {
		return errors
	}

	root := parser.GetRootMapping(node)
	if root == nil || root.Kind != yaml.MappingNode {
		return errors
	}

	fields := buildFieldMap(root, "")
	for _, field := range requiredFields {
		if _, exists := fields[field]; !exists {
			errors = append(errors, pkg.Error{
				Type:    "MissingRequiredField",
				Message: fmt.Sprintf("Required field '%s' is missing", field),
				Path:    field,
			})
		}
	}

	return errors
}

// buildFieldMap строит карту всех полей и их путей из YAML
func buildFieldMap(node *yaml.Node, prefix string) map[string]string {
	result := make(map[string]string)
	if node == nil || node.Kind != yaml.MappingNode {
		return result
	}
	for i := 0; i < len(node.Content); i += 2 {
		if i+1 >= len(node.Content) {
			break
		}
		key := node.Content[i].Value
		fullPath := key
		if prefix != "" {
			fullPath = prefix + "." + key
		}
		result[fullPath] = fullPath
		val := node.Content[i+1]
		if val.Kind == yaml.MappingNode {
			for k, v := range buildFieldMap(val, fullPath) {
				result[k] = v
			}
		}
	}
	return result
}
