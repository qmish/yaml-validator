package validator

import (
	"fmt"
	"strings"

	"yaml-validator/internal/parser"
	"yaml-validator/pkg"

	"gopkg.in/yaml.v3"
)

// CheckDuplicates ищет дублирующиеся ключи на одном уровне
func CheckDuplicates(node *yaml.Node) []pkg.Error {
	var errors []pkg.Error

	parser.TraverseMappings(node, "", func(mappingNode *yaml.Node, path string) {
		if mappingNode.Kind != yaml.MappingNode {
			return
		}
		keyTracker := make(map[string][]string)
		for i := 0; i < len(mappingNode.Content); i += 2 {
			if i+1 >= len(mappingNode.Content) {
				break
			}
			key := mappingNode.Content[i].Value
			keyPath := path
			if keyPath != "" {
				keyPath = path + "." + key
			} else {
				keyPath = key
			}

			if existing, exists := keyTracker[key]; exists {
				allPaths := append(existing, keyPath)
				keyNode := mappingNode.Content[i]
				e := pkg.Error{
					Type:    "DuplicateKey",
					Message: fmt.Sprintf("Duplicate key '%s' found in paths: %s",
						key, strings.Join(allPaths, ", ")),
					Path: keyPath,
				}
				if keyNode.Line > 0 {
					e.Line = keyNode.Line
					e.Column = keyNode.Column
				}
				errors = append(errors, e)
			} else {
				keyTracker[key] = []string{keyPath}
			}
		}
	})

	return errors
}
