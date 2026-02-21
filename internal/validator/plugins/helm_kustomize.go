package plugins

import (
	"strings"

	"yaml-validator/internal/validator"
	"yaml-validator/pkg"

	"gopkg.in/yaml.v3"
)

// HelmKustomizeValidator проверяет Helm Chart.yaml и Kustomize kustomization.yaml
type HelmKustomizeValidator struct{}

func init() {
	validator.RegisterPlugin("helm-kustomize", &HelmKustomizeValidator{})
}

// Name возвращает имя плагина
func (h *HelmKustomizeValidator) Name() string {
	return "HelmKustomizeValidator"
}

// Validate проверяет Helm Chart.yaml (apiVersion, name, version) и Kustomize (apiVersion, kind, resources)
func (h *HelmKustomizeValidator) Validate(node *yaml.Node) []pkg.Error {
	var errors []pkg.Error

	root := getRootNode(node)
	if root == nil || root.Kind != yaml.MappingNode {
		return errors
	}

	fields := buildFlatMap(root, "")

	// Helm Chart.yaml: apiVersion, name, version
	if isHelmChart(fields) {
		errors = append(errors, validateHelmChart(root, fields)...)
		return errors
	}

	// Kustomize kustomization
	if isKustomize(fields) {
		errors = append(errors, validateKustomize(root, fields)...)
	}

	return errors
}

func getRootNode(node *yaml.Node) *yaml.Node {
	if node == nil {
		return nil
	}
	if node.Kind == yaml.DocumentNode && len(node.Content) > 0 {
		return node.Content[0]
	}
	return node
}

func isHelmChart(fields map[string]string) bool {
	_, hasName := fields["name"]
	apiVersion, hasAPI := fields["apiVersion"]
	// Chart.yaml: name, apiVersion (v1 or v2); нет kind, metadata
	_, hasKind := fields["kind"]
	_, hasMetadata := fields["metadata"]
	if hasKind || hasMetadata {
		return false // K8s или Kustomize
	}
	if hasName && hasAPI {
		return strings.HasPrefix(apiVersion, "v1") || strings.HasPrefix(apiVersion, "v2")
	}
	return false
}

func isKustomize(fields map[string]string) bool {
	apiVersion, hasAPI := fields["apiVersion"]
	kind, hasKind := fields["kind"]
	if !hasAPI || !hasKind {
		return false
	}
	return strings.Contains(apiVersion, "kustomize") && kind == "Kustomization"
}

func validateHelmChart(root *yaml.Node, fields map[string]string) []pkg.Error {
	var errors []pkg.Error

	if _, ok := fields["apiVersion"]; !ok {
		errors = append(errors, pkg.Error{
			Type:       "HelmChartMissingAPIVersion",
			Message:    "Chart.yaml must have 'apiVersion' (v1 or v2)",
			Suggestion: "add apiVersion: v2 to Chart.yaml",
			Path:       "apiVersion",
		})
	}
	if _, ok := fields["name"]; !ok {
		errors = append(errors, pkg.Error{
			Type:       "HelmChartMissingName",
			Message:    "Chart.yaml must have 'name'",
			Suggestion: "add name: <chart-name> to Chart.yaml",
			Path:       "name",
		})
	}
	if _, ok := fields["version"]; !ok {
		errors = append(errors, pkg.Error{
			Type:       "HelmChartMissingVersion",
			Message:    "Chart.yaml must have 'version'",
			Suggestion: "add version: 0.1.0 to Chart.yaml",
			Path:       "version",
		})
	}

	return errors
}

func validateKustomize(root *yaml.Node, fields map[string]string) []pkg.Error {
	var errors []pkg.Error

	// buildFlatMap не добавляет SequenceNode в fields — проверяем через findChild
	hasResources := findChild(root, "resources") != nil
	hasBases := findChild(root, "bases") != nil
	if !hasResources && !hasBases {
		errors = append(errors, pkg.Error{
			Type:       "KustomizeMissingResources",
			Message:    "kustomization must have 'resources' (or 'bases' for legacy)",
			Suggestion: "add resources: [...] to kustomization.yaml",
			Path:       "resources",
		})
	}

	return errors
}
