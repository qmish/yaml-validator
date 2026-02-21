package validator

import (
	"yaml-validator/pkg"

	"gopkg.in/yaml.v3"
)

// CheckQuotedValues проверяет, что строковые значения (не ключи) в YAML в кавычках (как yamllint quoted-values).
func CheckQuotedValues(node *yaml.Node) []pkg.Error {
	var errors []pkg.Error
	traverseValueScalars(node, func(n *yaml.Node) {
		if n.Kind != yaml.ScalarNode {
			return
		}
		tag := n.ShortTag()
		if tag != "!!str" {
			return
		}
		if n.Style&(yaml.DoubleQuotedStyle|yaml.SingleQuotedStyle|yaml.LiteralStyle|yaml.FoldedStyle) != 0 {
			return
		}
		line := n.Line
		if line == 0 {
			line = 1
		}
		errors = append(errors, pkg.Error{
			Type:       "QuotedValues",
			Message:    "string value should be quoted",
			Suggestion: "wrap value in double or single quotes",
			Line:       line,
			Column:     n.Column,
		})
	})
	return errors
}

// traverseValueScalars обходит только скалярные значения (не ключи в маппингах).
func traverseValueScalars(node *yaml.Node, fn func(*yaml.Node)) {
	if node == nil {
		return
	}
	if node.Kind == yaml.ScalarNode {
		fn(node)
		return
	}
	if node.Kind == yaml.MappingNode {
		for i := 1; i < len(node.Content); i += 2 {
			traverseValueScalars(node.Content[i], fn)
		}
		return
	}
	for _, c := range node.Content {
		traverseValueScalars(c, fn)
	}
}
