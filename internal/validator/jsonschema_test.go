package validator

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"yaml-validator/internal/config"
)

func TestCheckJsonSchema_Disabled(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "config.yaml")
	require.NoError(t, os.WriteFile(f, []byte("name: test\nversion: 1.0\n"), 0644))

	opts := config.JsonSchemaOptions{Enabled: false}
	errors := CheckJsonSchema(f, opts)
	assert.Empty(t, errors)
}

func TestCheckJsonSchema_EmptyPath(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "config.yaml")
	require.NoError(t, os.WriteFile(f, []byte("name: test\nversion: 1.0\n"), 0644))

	opts := config.JsonSchemaOptions{Enabled: true, SchemaPath: ""}
	errors := CheckJsonSchema(f, opts)
	assert.Empty(t, errors)
}

func TestCheckJsonSchema_Valid(t *testing.T) {
	tmp := t.TempDir()
	schemaPath := filepath.Join(tmp, "schema.json")
	schemaContent := `{"$schema":"https://json-schema.org/draft/2020-12/schema","type":"object","required":["name","version"],"properties":{"name":{"type":"string"},"version":{"type":"string"}}}`
	require.NoError(t, os.WriteFile(schemaPath, []byte(schemaContent), 0644))

	f := filepath.Join(tmp, "config.yaml")
	require.NoError(t, os.WriteFile(f, []byte("name: myapp\nversion: 1.0.0\n"), 0644))

	opts := config.JsonSchemaOptions{Enabled: true, SchemaPath: schemaPath}
	errors := CheckJsonSchema(f, opts)
	assert.Empty(t, errors)
}

func TestCheckJsonSchema_Invalid(t *testing.T) {
	tmp := t.TempDir()
	schemaPath := filepath.Join(tmp, "schema.json")
	schemaContent := `{"$schema":"https://json-schema.org/draft/2020-12/schema","type":"object","required":["name","version"],"properties":{"name":{"type":"string"},"version":{"type":"string"}}}`
	require.NoError(t, os.WriteFile(schemaPath, []byte(schemaContent), 0644))

	f := filepath.Join(tmp, "config.yaml")
	require.NoError(t, os.WriteFile(f, []byte("name: myapp\n"), 0644))

	opts := config.JsonSchemaOptions{Enabled: true, SchemaPath: schemaPath}
	errors := CheckJsonSchema(f, opts)
	require.NotEmpty(t, errors)
	assert.Equal(t, "JsonSchema", errors[0].Type)
}

func TestCheckJsonSchema_ValidFromURL(t *testing.T) {
	schemaContent := `{"$schema":"https://json-schema.org/draft/2020-12/schema","type":"object","required":["name","version"],"properties":{"name":{"type":"string"},"version":{"type":"string"}}}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(schemaContent))
	}))
	defer srv.Close()

	tmp := t.TempDir()
	f := filepath.Join(tmp, "config.yaml")
	require.NoError(t, os.WriteFile(f, []byte("name: myapp\nversion: 1.0.0\n"), 0644))

	opts := config.JsonSchemaOptions{Enabled: true, SchemaPath: srv.URL + "/schema.json", SchemaCacheDir: filepath.Join(tmp, "cache")}
	errors := CheckJsonSchema(f, opts)
	assert.Empty(t, errors)
	// второй вызов — из кэша
	errors = CheckJsonSchema(f, opts)
	assert.Empty(t, errors)
}

func TestCheckJsonSchema_InvalidFromURL(t *testing.T) {
	schemaContent := `{"$schema":"https://json-schema.org/draft/2020-12/schema","type":"object","required":["name","version"],"properties":{"name":{"type":"string"},"version":{"type":"string"}}}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(schemaContent))
	}))
	defer srv.Close()

	tmp := t.TempDir()
	f := filepath.Join(tmp, "config.yaml")
	require.NoError(t, os.WriteFile(f, []byte("name: myapp\n"), 0644))

	opts := config.JsonSchemaOptions{Enabled: true, SchemaPath: srv.URL + "/schema.json"}
	errors := CheckJsonSchema(f, opts)
	require.NotEmpty(t, errors)
	assert.Equal(t, "JsonSchema", errors[0].Type)
}
