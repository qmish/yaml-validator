package validator

import (
	"fmt"
	"strings"

	"yaml-validator/internal/config"
	"yaml-validator/internal/parser"
	"yaml-validator/pkg"

	"gopkg.in/yaml.v3"
)

// CheckUniqueListFields проверяет, что элементы массивов (по выбранному полю) не повторяются (5.5).
// Например, spec.template.spec.containers[].name — имена контейнеров должны быть уникальны.
func CheckUniqueListFields(node *yaml.Node, rules []config.UniqueListField) []pkg.Error {
	if len(rules) == 0 {
		return nil
	}
	pathToField := make(map[string]string)
	for _, r := range rules {
		if r.Path != "" && r.Field != "" {
			pathToField[r.Path] = r.Field
		}
	}
	if len(pathToField) == 0 {
		return nil
	}

	var errors []pkg.Error

	parser.TraverseSequences(node, "", func(seqNode *yaml.Node, path string) {
		field, ok := pathToField[path]
		if !ok {
			return
		}
		if seqNode.Kind != yaml.SequenceNode || len(seqNode.Content) < 2 {
			return
		}

		seen := make(map[string]int) // value -> first index
		for idx, elem := range seqNode.Content {
			if elem.Kind != yaml.MappingNode {
				continue
			}
			val := getMappingFieldValue(elem, field)
			if val == "" {
				continue
			}
			if firstIdx, exists := seen[val]; exists {
				line := elem.Line
				if line <= 0 {
					line = seqNode.Content[firstIdx].Line
				}
				col := elem.Column
				errors = append(errors, pkg.Error{
					Type:       "DuplicateListElement",
					Message:    fmt.Sprintf("duplicate value '%s' in list at path '%s' (field '%s')", val, path, field),
					Suggestion: fmt.Sprintf("ensure unique values for field '%s' in %s", field, path),
					Line:       line,
					Column:     col,
				})
			} else {
				seen[val] = idx
			}
		}
	})

	return errors
}

// getMappingFieldValue возвращает строковое значение поля в маппинге.
func getMappingFieldValue(mapping *yaml.Node, field string) string {
	if mapping == nil || mapping.Kind != yaml.MappingNode || field == "" {
		return ""
	}
	field = strings.TrimSpace(field)
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		keyNode := mapping.Content[i]
		if keyNode.Value == field {
			return scalarValueString(mapping.Content[i+1])
		}
	}
	return ""
}
