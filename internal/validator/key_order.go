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
					Type:    "KeyOrdering",
					Message: fmt.Sprintf("key '%s' should be before '%s' (lexicographic order)", next, key),
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
