package plugins

import "gopkg.in/yaml.v3"

// buildFlatMap создаёт плоскую map путь->значение из YAML mapping node
func buildFlatMap(node *yaml.Node, prefix string) map[string]string {
	result := make(map[string]string)
	if node == nil || node.Kind != yaml.MappingNode {
		return result
	}
	for i := 0; i < len(node.Content); i += 2 {
		if i+1 >= len(node.Content) {
			break
		}
		key := node.Content[i].Value
		val := node.Content[i+1]
		fullPath := key
		if prefix != "" {
			fullPath = prefix + "." + key
		}
		if val.Kind == yaml.ScalarNode {
			result[fullPath] = val.Value
		} else if val.Kind == yaml.MappingNode {
			result[fullPath] = fullPath
			for k, v := range buildFlatMap(val, fullPath) {
				result[k] = v
			}
		}
	}
	return result
}

// findChild возвращает дочерний узел mapping по ключу
func findChild(node *yaml.Node, key string) *yaml.Node {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i < len(node.Content); i += 2 {
		if i+1 >= len(node.Content) {
			break
		}
		if node.Content[i].Value == key {
			return node.Content[i+1]
		}
	}
	return nil
}
