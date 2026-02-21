package validator

import (
	"fmt"

	"yaml-validator/internal/parser"
	"yaml-validator/pkg"

	"gopkg.in/yaml.v3"
)

// CheckKeyOrdering проверяет, что ключи в каждом маппинге отсортированы лексикографически.
func CheckKeyOrdering(node *yaml.Node) []pkg.Error {
	var errors []pkg.Error

	parser.TraverseMappings(node, "", func(mappingNode *yaml.Node, _ string) {
		if mappingNode.Kind != yaml.MappingNode {
			return
		}
		for i := 0; i+2 < len(mappingNode.Content); i += 2 {
			keyNode := mappingNode.Content[i]
			nextKeyNode := mappingNode.Content[i+2]
			key, next := keyNode.Value, nextKeyNode.Value
			if key > next {
				e := pkg.Error{
					Type:       "KeyOrdering",
					Message:    fmt.Sprintf("key '%s' should be before '%s' (lexicographic order)", next, key),
					Suggestion: "reorder keys lexicographically",
					Line:    nextKeyNode.Line,
					Column:  nextKeyNode.Column,
				}
				if e.Line <= 0 {
					e.Line = keyNode.Line
					e.Column = keyNode.Column
				}
				errors = append(errors, e)
			}
		}
	})

	return errors
}

// keyPriority возвращает приоритет ключа (меньше = выше). Ключи вне списка имеют приоритет len(order).
func keyPriority(key string, order []string) int {
	for i, k := range order {
		if k == key {
			return i
		}
	}
	return len(order)
}

// CheckKeyOrderConfigurable проверяет порядок ключей по заданному списку (напр. apiVersion, kind, metadata, spec).
func CheckKeyOrderConfigurable(node *yaml.Node, order []string) []pkg.Error {
	if len(order) == 0 {
		return nil
	}
	var errors []pkg.Error

	parser.TraverseMappings(node, "", func(mappingNode *yaml.Node, _ string) {
		if mappingNode.Kind != yaml.MappingNode {
			return
		}
		for i := 0; i+2 < len(mappingNode.Content); i += 2 {
			keyNode := mappingNode.Content[i]
			nextKeyNode := mappingNode.Content[i+2]
			key, next := keyNode.Value, nextKeyNode.Value
			pi, pj := keyPriority(key, order), keyPriority(next, order)
			if pi > pj {
				e := pkg.Error{
					Type:       "KeyOrderConfigurable",
					Message:    fmt.Sprintf("key '%s' should be before '%s'", next, key),
					Suggestion: "reorder keys as per config",
					Line:       nextKeyNode.Line,
					Column:     nextKeyNode.Column,
				}
				if e.Line <= 0 {
					e.Line = keyNode.Line
					e.Column = keyNode.Column
				}
				errors = append(errors, e)
			}
		}
	})

	return errors
}

// CheckMaxKeyNameLength проверяет, что длина имён ключей не превышает maxLen.
func CheckMaxKeyNameLength(node *yaml.Node, maxLen int) []pkg.Error {
	if maxLen <= 0 {
		return nil
	}
	var errors []pkg.Error

	parser.TraverseMappings(node, "", func(mappingNode *yaml.Node, _ string) {
		if mappingNode.Kind != yaml.MappingNode {
			return
		}
		for i := 0; i < len(mappingNode.Content); i += 2 {
			keyNode := mappingNode.Content[i]
			if len(keyNode.Value) > maxLen {
				line := keyNode.Line
				if line <= 0 {
					line = 1
				}
				errors = append(errors, pkg.Error{
					Type:       "MaxKeyNameLength",
					Message:    fmt.Sprintf("key name '%s' exceeds max length %d", keyNode.Value, maxLen),
					Suggestion: "shorten key name",
					Line:    line,
					Column:  keyNode.Column,
				})
			}
		}
	})

	return errors
}
