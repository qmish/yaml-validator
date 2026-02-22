package plugins

import (
	"regexp"
	"strings"

	"yaml-validator/internal/parser"
	"yaml-validator/internal/validator"
	"yaml-validator/pkg"

	"gopkg.in/yaml.v3"
)

// TerraformValidator проверяет YAML, используемый в Terraform (tfvars.yaml, backend config)
type TerraformValidator struct{}

func init() {
	validator.RegisterPlugin("terraform", &TerraformValidator{})
}

// terraformVarNameRegex — допустимые символы в Terraform variable name
var terraformVarNameRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

// Name возвращает имя плагина
func (t *TerraformValidator) Name() string {
	return "TerraformValidator"
}

// Validate проверяет YAML в контексте Terraform (tfvars, backend config)
func (t *TerraformValidator) Validate(node *yaml.Node) []pkg.Error {
	var errors []pkg.Error

	root := parser.GetRootMapping(node)
	if root == nil || root.Kind != yaml.MappingNode {
		return errors
	}

	fields := buildFlatMap(root, "")

	// Пропускаем K8s, docker-compose, Ansible, Helm, Kustomize
	if _, hasAPI := fields["apiVersion"]; hasAPI {
		return errors
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
	// Helm Chart (name + apiVersion v1/v2) — helm plugin
	if name, hasName := fields["name"]; hasName {
		if apiVersion, hasAPI := fields["apiVersion"]; hasAPI {
			if strings.HasPrefix(apiVersion, "v1") || strings.HasPrefix(apiVersion, "v2") {
				return errors
			}
		}
		_ = name
	}
	if kind, hasKind := fields["kind"]; hasKind && kind == "Kustomization" {
		return errors
	}

	// Топ-уровневые ключи: в tfvars имена переменных — snake_case (underscore, не hyphen)
	for i := 0; i < len(root.Content); i += 2 {
		if i+1 >= len(root.Content) {
			break
		}
		keyNode := root.Content[i]
		key := keyNode.Value
		if key == "" {
			continue
		}
		if strings.Contains(key, "-") && !terraformVarNameRegex.MatchString(key) {
			suggested := strings.ReplaceAll(key, "-", "_")
			errors = append(errors, pkg.Error{
				Type:       "TerraformVarNameHyphen",
				Message:    "Terraform variable names use underscores, not hyphens: '" + key + "'",
				Suggestion: "use '" + suggested + "'",
				Path:       key,
				Line:       keyNode.Line,
				Column:     keyNode.Column,
			})
		}
	}

	return errors
}
