package validator

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"yaml-validator/internal/config"
	"yaml-validator/pkg"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

// httpURLLoader загружает JSON Schema по http/https URL с опциональным кэшем.
type httpURLLoader struct {
	cacheDir string
	client   *http.Client
}

func (l *httpURLLoader) Load(url string) (any, error) {
	cachePath := ""
	if l.cacheDir != "" {
		h := sha256.Sum256([]byte(url))
		cachePath = filepath.Join(l.cacheDir, fmt.Sprintf("%x.json", h[:8]))
		if data, err := os.ReadFile(cachePath); err == nil {
			return jsonschema.UnmarshalJSON(strings.NewReader(string(data)))
		}
	}
	resp, err := l.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, url)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	doc, err := jsonschema.UnmarshalJSON(strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}
	if cachePath != "" {
		_ = os.MkdirAll(l.cacheDir, 0755)
		_ = os.WriteFile(cachePath, data, 0644)
	}
	return doc, nil
}

// schemeURLLoader поддерживает file, http и https.
func newSchemeURLLoader(cacheDir string) jsonschema.URLLoader {
	loader := jsonschema.SchemeURLLoader{
		"file": jsonschema.FileLoader{},
	}
	if cacheDir != "" {
		_ = os.MkdirAll(cacheDir, 0755)
	}
	hl := &httpURLLoader{cacheDir: cacheDir, client: &http.Client{}}
	loader["http"] = hl
	loader["https"] = hl
	return loader
}

func isSchemaURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

// CheckJsonSchema проверяет YAML по произвольной JSON Schema.
// schema_path — путь к файлу или URL (https://...). При URL и schema_cache_dir схема кэшируется локально.
func CheckJsonSchema(filename string, opts config.JsonSchemaOptions) []pkg.Error {
	if !opts.Enabled || opts.SchemaPath == "" {
		return nil
	}

	var schemaURL string
	if isSchemaURL(opts.SchemaPath) {
		schemaURL = opts.SchemaPath
	} else {
		absPath, err := filepath.Abs(opts.SchemaPath)
		if err != nil {
			return []pkg.Error{{
				Type:    "JsonSchemaInit",
				Message: fmt.Sprintf("invalid schema path: %v", err),
			}}
		}
		schemaURL = "file:///" + filepath.ToSlash(absPath)
	}

	compiler := jsonschema.NewCompiler()
	if isSchemaURL(opts.SchemaPath) {
		compiler.UseLoader(newSchemeURLLoader(opts.SchemaCacheDir))
	}
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
