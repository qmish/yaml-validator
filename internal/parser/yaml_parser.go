package parser

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// ParseFileMulti читает и парсит мультидокументный YAML-файл, возвращает срез узлов (по одному на документ)
func ParseFileMulti(filename string) ([]*yaml.Node, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	dec := yaml.NewDecoder(bytes.NewReader(data))
	var nodes []*yaml.Node
	for {
		var node yaml.Node
		err := dec.Decode(&node)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		nodes = append(nodes, &node)
	}
	if len(nodes) == 0 {
		return nil, fmt.Errorf("no YAML documents in file")
	}
	return nodes, nil
}

// NodeInfo содержит информацию об узле YAML
type NodeInfo struct {
	Node  *yaml.Node
	Path  string
	Value interface{}
}

// ParseBytesMulti парсит мультидокументный YAML из байтов (для LSP, несохранённые буферы)
func ParseBytesMulti(data []byte) ([]*yaml.Node, error) {
	dec := yaml.NewDecoder(bytes.NewReader(data))
	var nodes []*yaml.Node
	for {
		var node yaml.Node
		err := dec.Decode(&node)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		nodes = append(nodes, &node)
	}
	if len(nodes) == 0 {
		return nil, fmt.Errorf("no YAML documents")
	}
	return nodes, nil
}

// ParseFile читает и парсит YAML-файл
func ParseFile(filename string) (*yaml.Node, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var node yaml.Node
	err = yaml.Unmarshal(data, &node)
	if err != nil {
		return nil, err
	}
	return &node, nil
}

// Traverse рекурсивно обходит структуру YAML и вызывает callback для каждого узла
func Traverse(node *yaml.Node, path string, callback func(NodeInfo)) {
	if node == nil {
		return
	}
	info := NodeInfo{Node: node, Path: path, Value: node.Value}
	callback(info)
	for _, child := range node.Content {
		childPath := path
		if path != "" {
			childPath = path + "." + child.Value
		} else if child.Value != "" {
			childPath = child.Value
		}
		Traverse(child, childPath, callback)
	}
}

// TraverseMappings обходит все узлы и для каждого MappingNode вызывает callback с ключами и путями
func TraverseMappings(node *yaml.Node, path string, callback func(node *yaml.Node, path string)) {
	if node == nil {
		return
	}
	if node.Kind == yaml.MappingNode {
		callback(node, path)
		for i := 0; i < len(node.Content); i += 2 {
			if i+1 < len(node.Content) {
				key := node.Content[i].Value
				var newPath string
				if path != "" {
					newPath = path + "." + key
				} else {
					newPath = key
				}
				TraverseMappings(node.Content[i+1], newPath, callback)
			}
		}
	} else if node.Kind == yaml.SequenceNode {
		for j, child := range node.Content {
			newPath := fmt.Sprintf("%s[%d]", path, j)
			TraverseMappings(child, newPath, callback)
		}
	} else if node.Kind == yaml.DocumentNode {
		for _, child := range node.Content {
			TraverseMappings(child, path, callback)
		}
	} else {
		for _, child := range node.Content {
			TraverseMappings(child, path, callback)
		}
	}
}

// TraverseSequences обходит дерево и для каждого SequenceNode вызывает callback(seqNode, path).
// path — путь к массиву (напр. "spec.template.spec.containers").
func TraverseSequences(node *yaml.Node, path string, callback func(seqNode *yaml.Node, seqPath string)) {
	if node == nil {
		return
	}
	if node.Kind == yaml.SequenceNode {
		callback(node, path)
		for j, child := range node.Content {
			newPath := fmt.Sprintf("%s[%d]", path, j)
			TraverseSequences(child, newPath, callback)
		}
		return
	}
	if node.Kind == yaml.MappingNode {
		for i := 0; i < len(node.Content); i += 2 {
			if i+1 < len(node.Content) {
				key := node.Content[i].Value
				var newPath string
				if path != "" {
					newPath = path + "." + key
				} else {
					newPath = key
				}
				TraverseSequences(node.Content[i+1], newPath, callback)
			}
		}
		return
	}
	if node.Kind == yaml.DocumentNode {
		for _, child := range node.Content {
			TraverseSequences(child, path, callback)
		}
		return
	}
	for _, child := range node.Content {
		TraverseSequences(child, path, callback)
	}
}

// GetRootMapping возвращает корневой MappingNode из документа
func GetRootMapping(node *yaml.Node) *yaml.Node {
	if node == nil {
		return nil
	}
	if node.Kind == yaml.DocumentNode && len(node.Content) > 0 {
		return node.Content[0]
	}
	if node.Kind == yaml.MappingNode {
		return node
	}
	return nil
}

// GetValueAtPath извлекает скалярное значение по dot-пути (например, "metadata.name").
// Возвращает строковое представление и true, если найдено; иначе "", false.
func GetValueAtPath(root *yaml.Node, path string) (string, bool) {
	node := GetRootMapping(root)
	if node == nil || path == "" {
		return "", false
	}
	parts := splitPath(path)
	for _, key := range parts {
		if node == nil || node.Kind != yaml.MappingNode {
			return "", false
		}
		found := false
		for i := 0; i < len(node.Content); i += 2 {
			if i+1 >= len(node.Content) {
				break
			}
			if node.Content[i].Value == key {
				node = node.Content[i+1]
				found = true
				break
			}
		}
		if !found {
			return "", false
		}
	}
	if node == nil {
		return "", false
	}
	if node.Kind == yaml.ScalarNode {
		return node.Value, true
	}
	return "", false
}

func splitPath(path string) []string {
	if path == "" {
		return nil
	}
	var parts []string
	for _, s := range strings.Split(path, ".") {
		if s != "" {
			parts = append(parts, s)
		}
	}
	return parts
}
