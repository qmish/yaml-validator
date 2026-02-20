package plugins

import (
	"fmt"

	"yaml-validator/internal/parser"
	"yaml-validator/internal/validator"
	"yaml-validator/pkg"

	"gopkg.in/yaml.v3"
)

// K8sValidator проверяет Kubernetes-специфичные поля в манифестах
type K8sValidator struct{}

func init() {
	validator.RegisterPlugin("kubernetes", &K8sValidator{})
}

// Name возвращает имя плагина
func (k *K8sValidator) Name() string {
	return "KubernetesValidator"
}

// Validate проверяет манифест на соответствие требованиям Kubernetes
func (k *K8sValidator) Validate(node *yaml.Node) []pkg.Error {
	var errors []pkg.Error

	root := parser.GetRootMapping(node)
	if root == nil || root.Kind != yaml.MappingNode {
		return errors
	}

	fields := buildFlatMap(root, "")
	apiVersion, hasAPI := fields["apiVersion"]
	kind, hasKind := fields["kind"]

	if !hasAPI && !hasKind {
		return errors
	}

	if hasAPI && !hasKind {
		errors = append(errors, pkg.Error{
			Type:    "K8sMissingKind",
			Message: "Kubernetes manifest has apiVersion but missing required field 'kind'",
			Path:    "kind",
		})
	}
	if hasKind && !hasAPI {
		errors = append(errors, pkg.Error{
			Type:    "K8sMissingAPIVersion",
			Message: "Kubernetes manifest has kind but missing required field 'apiVersion'",
			Path:    "apiVersion",
		})
	}

	if apiVersion != "" && kind != "" {
		metadata, ok := fields["metadata"]
		if !ok || metadata == "" {
			errors = append(errors, pkg.Error{
				Type:    "K8sMissingMetadata",
				Message: "Kubernetes manifest must have 'metadata' section",
				Path:    "metadata",
			})
		}
		if _, hasName := fields["metadata.name"]; metadata != "" && !hasName {
			errors = append(errors, pkg.Error{
				Type:    "K8sMissingName",
				Message: "Kubernetes manifest metadata must have 'name' field",
				Path:    "metadata.name",
			})
		}
	}

	if kind != "" {
		validKinds := map[string]bool{
			"Pod": true, "Deployment": true, "Service": true, "ConfigMap": true,
			"Secret": true, "Namespace": true, "Ingress": true, "Job": true,
			"CronJob": true, "DaemonSet": true, "StatefulSet": true,
		}
		if !validKinds[kind] {
			errors = append(errors, pkg.Error{
				Type:    "K8sUnknownKind",
				Message: fmt.Sprintf("Unknown Kubernetes resource kind: '%s'", kind),
				Path:    "kind",
			})
		}
	}

	return errors
}

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
