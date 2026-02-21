package validator

import (
	"fmt"

	"yaml-validator/internal/parser"
	"yaml-validator/pkg"

	"gopkg.in/yaml.v3"
)

// CheckForbidDefaultValues проверяет, что ключи с заданными значениями по умолчанию не используются.
// Например, для K8s: imagePullPolicy: Always — значение по умолчанию, лишнее указание.
func CheckForbidDefaultValues(node *yaml.Node, defaults map[string]string) []pkg.Error {
	if len(defaults) == 0 {
		return nil
	}
	var errors []pkg.Error

	parser.TraverseMappings(node, "", func(mappingNode *yaml.Node, _ string) {
		if mappingNode.Kind != yaml.MappingNode {
			return
		}
		for i := 0; i+1 < len(mappingNode.Content); i += 2 {
			keyNode := mappingNode.Content[i]
			valNode := mappingNode.Content[i+1]
			key := keyNode.Value
			defaultVal, ok := defaults[key]
			if !ok {
				continue
			}
			actual := scalarValueString(valNode)
			if actual != "" && actual == defaultVal {
				line := keyNode.Line
				if line <= 0 {
					line = valNode.Line
				}
				col := keyNode.Column
				if col <= 0 {
					col = valNode.Column
				}
				errors = append(errors, pkg.Error{
					Type:       "ForbidDefaultValue",
					Message:    fmt.Sprintf("key '%s' has default value '%s', remove it", key, defaultVal),
					Suggestion: fmt.Sprintf("remove key '%s' (it is the default)", key),
					Line:       line,
					Column:     col,
				})
			}
		}
	})

	return errors
}

// scalarValueString возвращает строковое представление скалярного значения.
func scalarValueString(n *yaml.Node) string {
	if n == nil {
		return ""
	}
	if n.Kind == yaml.ScalarNode {
		return n.Value
	}
	return ""
}
