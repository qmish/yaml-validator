package plugins

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestOpenAPIValidator_Valid(t *testing.T) {
	doc := `
openapi: "3.0.0"
info:
  title: My API
  version: 1.0.0
paths:
  /users:
    get:
      summary: List
`
	var node yaml.Node
	if err := yaml.Unmarshal([]byte(doc), &node); err != nil {
		t.Fatal(err)
	}
	v := &OpenAPIValidator{}
	errs := v.Validate(&node)
	if len(errs) != 0 {
		t.Errorf("Expected no errors, got %d: %v", len(errs), errs)
	}
}

func TestOpenAPIValidator_MissingVersion(t *testing.T) {
	doc := `
openapi: "3.0.0"
info:
  title: API
paths: {}
`
	var node yaml.Node
	if err := yaml.Unmarshal([]byte(doc), &node); err != nil {
		t.Fatal(err)
	}
	v := &OpenAPIValidator{}
	errs := v.Validate(&node)
	if len(errs) != 1 {
		t.Errorf("Expected 1 error, got %d: %v", len(errs), errs)
	}
	if len(errs) > 0 && errs[0].Type != "OpenAPIMissingInfoVersion" {
		t.Errorf("Expected OpenAPIMissingInfoVersion, got %s", errs[0].Type)
	}
}

func TestOpenAPIValidator_SkipsK8s(t *testing.T) {
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
	v := &OpenAPIValidator{}
	errs := v.Validate(&node)
	if len(errs) != 0 {
		t.Errorf("Expected skip for K8s, got %d errors", len(errs))
	}
}
