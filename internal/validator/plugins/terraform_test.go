package plugins

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestTerraformValidator_HyphenInKey(t *testing.T) {
	doc := `
region: us-east-1
instance-type: t3.micro
environment: prod
`
	var node yaml.Node
	if err := yaml.Unmarshal([]byte(doc), &node); err != nil {
		t.Fatal(err)
	}
	v := &TerraformValidator{}
	errs := v.Validate(&node)
	if len(errs) != 1 {
		t.Fatalf("Expected 1 error, got %d", len(errs))
	}
	if errs[0].Type != "TerraformVarNameHyphen" {
		t.Errorf("Expected TerraformVarNameHyphen, got %s", errs[0].Type)
	}
	if errs[0].Path != "instance-type" {
		t.Errorf("Expected path instance-type, got %s", errs[0].Path)
	}
}

func TestTerraformValidator_ValidKeys(t *testing.T) {
	doc := `
region: us-east-1
instance_type: t3.micro
environment: prod
`
	var node yaml.Node
	if err := yaml.Unmarshal([]byte(doc), &node); err != nil {
		t.Fatal(err)
	}
	v := &TerraformValidator{}
	errs := v.Validate(&node)
	if len(errs) != 0 {
		t.Errorf("Expected no errors, got %d: %v", len(errs), errs)
	}
}

func TestTerraformValidator_SkipsK8s(t *testing.T) {
	doc := `
apiVersion: v1
kind: Pod
metadata:
  name: test
`
	var node yaml.Node
	if err := yaml.Unmarshal([]byte(doc), &node); err != nil {
		t.Fatal(err)
	}
	v := &TerraformValidator{}
	errs := v.Validate(&node)
	if len(errs) != 0 {
		t.Errorf("Expected skip for K8s, got %d errors", len(errs))
	}
}
