package validator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"yaml-validator/internal/config"
	"yaml-validator/pkg"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

// CheckJsonSchema проверяет YAML по произвольной JSON Schema.
// Включается через rules.json_schema.enabled, путь к схеме — rules.json_schema.schema_path.
func CheckJsonSchema(filename string, opts config.JsonSchemaOptions) []pkg.Error {
	if !opts.Enabled || opts.SchemaPath == "" {
		return nil
	}

	absPath, err := filepath.Abs(opts.SchemaPath)
	if err != nil {
		return []pkg.Error{{
			Type:    "JsonSchemaInit",
			Message: fmt.Sprintf("invalid schema path: %v", err),
		}}
	}
	schemaURL := "file:///" + filepath.ToSlash(absPath)

	compiler := jsonschema.NewCompiler()
	schema, err := compiler.Compile(schemaURL)
	if err != nil {
		return []pkg.Error{{
			Type:    "JsonSchemaInit",
			Message: fmt.Sprintf("failed to load schema: %v", err),
		}}
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return []pkg.Error{{
			Type:    "JsonSchemaRead",
			Message: err.Error(),
		}}
	}

	var instance interface{}
	if err := yaml.Unmarshal(data, &instance); err != nil {
		return []pkg.Error{{
			Type:    "JsonSchemaParse",
			Message: fmt.Sprintf("failed to parse YAML: %v", err),
		}}
	}

	if err := schema.Validate(instance); err != nil {
		ve, ok := err.(*jsonschema.ValidationError)
		if !ok {
			return []pkg.Error{{
				Type:    "JsonSchema",
				Message: err.Error(),
			}}
		}
		return extractJsonSchemaErrors(ve)
	}

	return nil
}

func extractJsonSchemaErrors(ve *jsonschema.ValidationError) []pkg.Error {
	path := strings.Join(ve.InstanceLocation, ".")
	if path == "" {
		path = "/"
	}
	if len(ve.Causes) == 0 {
		return []pkg.Error{{
			Type:    "JsonSchema",
			Message: ve.Error(),
			Path:    path,
		}}
	}
	var errors []pkg.Error
	for _, c := range ve.Causes {
		errors = append(errors, extractJsonSchemaErrors(c)...)
	}
	if len(errors) == 0 {
		errors = append(errors, pkg.Error{
			Type:    "JsonSchema",
			Message: ve.Error(),
			Path:    path,
		})
	}
	return errors
}
