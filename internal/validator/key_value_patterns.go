package validator

import (
	"fmt"
	"regexp"

	"yaml-validator/internal/config"
	"yaml-validator/internal/parser"
	"yaml-validator/pkg"

	"gopkg.in/yaml.v3"
)

// ruleWithRE — правило с скомпилированным regexp
type ruleWithRE struct {
	path string
	re   *regexp.Regexp
}

// CheckKeyValuePatterns проверяет, что ключи или значения соответствуют регулярным выражениям (5.6).
func CheckKeyValuePatterns(node *yaml.Node, rules []config.KeyValuePattern) []pkg.Error {
	if len(rules) == 0 {
		return nil
	}

	var keyReRules, valueReRules []ruleWithRE
	for _, r := range rules {
		if r.Pattern == "" {
			continue
		}
		re, err := regexp.Compile(r.Pattern)
		if err != nil {
			continue
		}
		rr := ruleWithRE{path: r.Path, re: re}
		switch r.Target {
		case "keys":
			keyReRules = append(keyReRules, rr)
		case "values", "":
			valueReRules = append(valueReRules, rr)
		}
	}

	var errors []pkg.Error

	// Keys: path = путь к маппингу. Проверяем каждый ключ в маппинге по path.
	if len(keyReRules) > 0 {
		parser.TraverseMappings(node, "", func(mappingNode *yaml.Node, path string) {
			if mappingNode.Kind != yaml.MappingNode {
				return
			}
			for _, rr := range keyReRules {
				if rr.path != "" && path != rr.path {
					continue
				}
				for i := 0; i < len(mappingNode.Content); i += 2 {
					keyNode := mappingNode.Content[i]
					if !rr.re.MatchString(keyNode.Value) {
						line := keyNode.Line
						if line <= 0 {
							line = 1
						}
						errors = append(errors, pkg.Error{
							Type:       "KeyPatternMismatch",
							Message:    fmt.Sprintf("key '%s' at path '%s' does not match pattern '%s'", keyNode.Value, path, rr.re.String()),
							Suggestion: fmt.Sprintf("key must match: %s", rr.re.String()),
							Line:       line,
							Column:     keyNode.Column,
						})
					}
				}
			}
		})
	}

	// Values: path = полный путь к скаляру (напр. metadata.name). Проверяем каждый скаляр.
	if len(valueReRules) > 0 {
		traverseKeyValues(node, "", func(keyPath string, keyNode, valNode *yaml.Node) {
			if valNode.Kind != yaml.ScalarNode {
				return
			}
			val := valNode.Value
			for _, rr := range valueReRules {
				if rr.path != "" && keyPath != rr.path {
					continue
				}
				if !rr.re.MatchString(val) {
					line := valNode.Line
					if line <= 0 {
						line = 1
					}
					errors = append(errors, pkg.Error{
						Type:       "ValuePatternMismatch",
						Message:    fmt.Sprintf("value '%s' at path '%s' does not match pattern '%s'", val, keyPath, rr.re.String()),
						Suggestion: fmt.Sprintf("value must match: %s", rr.re.String()),
						Line:       line,
						Column:     valNode.Column,
					})
				}
			}
		})
	}

	return errors
}

// traverseKeyValues обходит все пары ключ-значение и вызывает callback с полным путём.
func traverseKeyValues(node *yaml.Node, path string, callback func(keyPath string, keyNode, valNode *yaml.Node)) {
	if node == nil {
		return
	}
	if node.Kind == yaml.MappingNode {
		for i := 0; i+1 < len(node.Content); i += 2 {
			keyNode := node.Content[i]
			valNode := node.Content[i+1]
			keyPath := path
			if path != "" {
				keyPath = path + "." + keyNode.Value
			} else {
				keyPath = keyNode.Value
			}
			callback(keyPath, keyNode, valNode)
			traverseKeyValues(valNode, keyPath, callback)
		}
		return
	}
	if node.Kind == yaml.SequenceNode {
		for j, child := range node.Content {
			childPath := fmt.Sprintf("%s[%d]", path, j)
			traverseKeyValues(child, childPath, callback)
		}
		return
	}
	if node.Kind == yaml.DocumentNode {
		for _, child := range node.Content {
			traverseKeyValues(child, path, callback)
		}
	}
}
