package plugins

import (
	"yaml-validator/internal/parser"
	"yaml-validator/internal/validator"
	"yaml-validator/pkg"

	"gopkg.in/yaml.v3"
)

// DockerComposeValidator проверяет docker-compose специфичные поля
type DockerComposeValidator struct{}

func init() {
	validator.RegisterPlugin("docker-compose", &DockerComposeValidator{})
}

// Name возвращает имя плагина
func (d *DockerComposeValidator) Name() string {
	return "DockerComposeValidator"
}

// Validate проверяет docker-compose файл (версия, services, image/build)
func (d *DockerComposeValidator) Validate(node *yaml.Node) []pkg.Error {
	var errors []pkg.Error

	root := parser.GetRootMapping(node)
	if root == nil || root.Kind != yaml.MappingNode {
		return errors
	}

	fields := buildFlatMap(root, "")

	// Только для docker-compose: есть services, нет apiVersion (K8s)
	if _, hasServices := fields["services"]; !hasServices {
		return errors
	}
	if _, hasAPI := fields["apiVersion"]; hasAPI {
		return errors
	}

	// Опционально: version (для Compose v1)
	// Пропускаем — version устарел в Compose V2

	// services — обязательная секция
	servicesNode := findChild(root, "services")
	if servicesNode == nil || servicesNode.Kind != yaml.MappingNode {
		return errors
	}

	for i := 0; i < len(servicesNode.Content); i += 2 {
		if i+1 >= len(servicesNode.Content) {
			break
		}
		svcName := servicesNode.Content[i].Value
		svcNode := servicesNode.Content[i+1]
		if svcNode.Kind != yaml.MappingNode {
			continue
		}
		svcFields := buildFlatMap(svcNode, "")
		if _, hasImage := svcFields["image"]; !hasImage {
			if _, hasBuild := svcFields["build"]; !hasBuild {
				line := servicesNode.Content[i].Line
				if line <= 0 {
					line = svcNode.Line
				}
				errors = append(errors, pkg.Error{
					Type:       "DockerComposeServiceImage",
					Message:    "service '" + svcName + "' must have 'image' or 'build'",
					Suggestion: "add 'image: <name>' or 'build: <path>' to service",
					Path:    "services." + svcName,
					Line:    line,
					Column:  servicesNode.Content[i].Column,
				})
			}
		}
	}

	return errors
}

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
