package plugins

import (
	"strings"

	"yaml-validator/internal/parser"
	"yaml-validator/internal/validator"
	"yaml-validator/pkg"

	"gopkg.in/yaml.v3"
)

// OpenAPIValidator проверяет YAML OpenAPI/Swagger спецификации
type OpenAPIValidator struct{}

func init() {
	validator.RegisterPlugin("openapi", &OpenAPIValidator{})
}

// Name возвращает имя плагина
func (o *OpenAPIValidator) Name() string {
	return "OpenAPIValidator"
}

// Validate проверяет OpenAPI/Swagger спецификацию (openapi/swagger, info, paths)
func (o *OpenAPIValidator) Validate(node *yaml.Node) []pkg.Error {
	var errors []pkg.Error

	root := parser.GetRootMapping(node)
	if root == nil || root.Kind != yaml.MappingNode {
		return errors
	}

	fields := buildFlatMap(root, "")

	// Пропускаем K8s, docker-compose, Ansible, Helm, Kustomize, Terraform
	if _, hasAPI := fields["apiVersion"]; hasAPI {
		if _, hasKind := fields["kind"]; hasKind {
			return errors
		}
	}
	if _, hasServices := fields["services"]; hasServices {
		return errors
	}
	if _, hasHosts := fields["hosts"]; hasHosts {
		return errors
	}
	if _, hasTasks := fields["tasks"]; hasTasks {
		return errors
	}
	if kind, hasKind := fields["kind"]; hasKind && kind == "Kustomization" {
		return errors
	}

	// OpenAPI 3: openapi, info, paths
	openapi, hasOpenAPI := fields["openapi"]
	// Swagger 2: swagger, info, paths
	swagger, hasSwagger := fields["swagger"]

	isOpenAPI3 := hasOpenAPI && (strings.HasPrefix(openapi, "3.") || openapi == "3.0")
	isSwagger2 := hasSwagger && (strings.HasPrefix(swagger, "2.") || swagger == "2.0")

	if !isOpenAPI3 && !isSwagger2 {
		return errors
	}

	// info (required)
	infoNode := findChild(root, "info")
	if infoNode == nil {
		errors = append(errors, pkg.Error{
			Type:       "OpenAPIMissingInfo",
			Message:    "OpenAPI/Swagger spec must have 'info' object",
			Suggestion: "add info: { title: ..., version: ... }",
			Path:       "info",
		})
	} else if infoNode.Kind == yaml.MappingNode {
		infoFields := buildFlatMap(infoNode, "")
		if _, ok := infoFields["title"]; !ok {
			errors = append(errors, pkg.Error{
				Type:       "OpenAPIMissingInfoTitle",
				Message:    "info must have 'title'",
				Suggestion: "add info.title",
				Path:       "info.title",
			})
		}
		if _, ok := infoFields["version"]; !ok {
			errors = append(errors, pkg.Error{
				Type:       "OpenAPIMissingInfoVersion",
				Message:    "info must have 'version'",
				Suggestion: "add info.version",
				Path:       "info.version",
			})
		}
	}

	// paths (required)
	if findChild(root, "paths") == nil {
		errors = append(errors, pkg.Error{
			Type:       "OpenAPIMissingPaths",
			Message:    "OpenAPI/Swagger spec must have 'paths' object",
			Suggestion: "add paths: { /path: {...} }",
			Path:       "paths",
		})
	}

	return errors
}
