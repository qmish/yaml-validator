package config

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// JsonSchemaOptions настройки проверки по произвольной JSON Schema
type JsonSchemaOptions struct {
	Enabled    bool   `yaml:"enabled"`
	SchemaPath string `yaml:"schema_path"` // путь к файлу схемы (JSON или YAML)
}

// K8sSchemaOptions настройки проверки по OpenAPI-схеме Kubernetes
type K8sSchemaOptions struct {
	Enabled    bool   `yaml:"enabled"`
	Version    string `yaml:"version"`     // например "1.28", "master"
	Strict     bool   `yaml:"strict"`      // отклонять недокументированные поля
	CacheDir   string `yaml:"cache_dir"`   // кэш схем (пусто = в памяти)
	IgnoreMissing bool `yaml:"ignore_missing_schemas"` // пропускать ресурсы без схемы (CRD)
}

// StyleOptions правила стиля (в духе yamllint)
type StyleOptions struct {
	RequireDocumentStart       bool `yaml:"require_document_start"`       // требовать --- в начале файла
	ForbidTrailingSpaces       bool `yaml:"forbid_trailing_spaces"`       // запрет пробелов в конце строки
	ForbidTrailingDots         bool `yaml:"forbid_trailing_dots"`         // запрет точек в конце строки
	RequireNewlineAtEof        bool `yaml:"require_newline_at_eof"`       // требовать перевод строки в конце файла
	ForbidConsecutiveEmptyLines     bool `yaml:"forbid_consecutive_empty_lines"`     // запрет более одной пустой строки подряд
	RequireEmptyLineBetweenBlocks   bool `yaml:"require_empty_line_between_blocks"`  // ровно одна пустая строка между топ-уровневыми блоками
	MinEmptyLinesBetweenBlocks      int  `yaml:"min_empty_lines_between_blocks"`     // минимум пустых строк между блоками (0 или 1); гибче, чем require_empty_line_between_blocks
	RequireDocumentEnd  bool `yaml:"require_document_end"` // требовать ... в конце файла (много-документный YAML)
	RequireCommentsIndented bool `yaml:"require_comments_indented"` // комментарии внутри блока должны иметь отступ (как в yamllint)
	RequireQuotedKeys    bool `yaml:"require_quoted_keys"`    // ключи маппинга должны быть в кавычках
	RequireQuotedValues  bool `yaml:"require_quoted_values"`  // строковые значения должны быть в кавычках
	IndentSpaces         int  `yaml:"indent_spaces"`          // шаг отступов (2 или 4 пробелов; 0 = отключено)
	ForbidTabs           bool `yaml:"forbid_tabs"`            // запрет табуляции (отдельное правило стиля; при true — ошибка TabInsteadOfSpaces)
	ForbidUnicode        bool `yaml:"forbid_unicode"`         // запрет не-ASCII в ключах и строках (строгие ASCII-конфиги)
	ForbidBOM            bool `yaml:"forbid_bom"`              // запрет BOM (Byte Order Mark) в начале файла
}

// UniqueListField — правило уникальности элементов массива по полю (5.5)
type UniqueListField struct {
	Path  string `yaml:"path"`  // путь к массиву, напр. "spec.template.spec.containers"
	Field string `yaml:"field"` // поле для проверки уникальности, напр. "name"
}

// ValidationRules определяет правила валидации
type ValidationRules struct {
	CheckSyntax       bool              `yaml:"check_syntax"`
	CheckDuplicates   bool              `yaml:"check_duplicates"`
	CheckIntegrity    bool              `yaml:"check_integrity"`
	CheckCommonErrors bool              `yaml:"check_common_errors"`
	CheckKeyOrdering   bool     `yaml:"check_key_ordering"`   // требовать лексикографический порядок ключей
	KeyOrder             []string          `yaml:"key_order"`               // приоритетный порядок ключей (напр. apiVersion, kind, metadata, spec)
	MaxKeyNameLength     int               `yaml:"max_key_name_length"`     // максимальная длина имён ключей (0 = отключено)
	ForbidDefaultValues  bool                `yaml:"forbid_default_values"`   // запрет ключей со значением по умолчанию (5.4)
	DefaultValues        map[string]string   `yaml:"default_values"`          // ключ -> значение по умолчанию
	UniqueListFields    []UniqueListField    `yaml:"unique_list_fields"`      // уникальность элементов по полю (5.5)
	InlineIgnore        bool                `yaml:"inline_ignore"`            // разрешить отключение правил через комментарии в YAML
	Style             StyleOptions      `yaml:"style"`
	JsonSchema        JsonSchemaOptions `yaml:"json_schema"`
	K8sSchema        K8sSchemaOptions   `yaml:"k8s_schema"`
	RequiredFields    []string          `yaml:"required_fields"`
	MaxLineLength     int               `yaml:"max_line_length"`
	SensitivePatterns []string          `yaml:"sensitive_patterns"`
}

// FileProfile — правило автовыбора конфига по маске пути к файлу
type FileProfile struct {
	Pattern string `yaml:"pattern"` // маска, напр. "**/k8s/**", "*docker-compose*.yaml"
	Config  string `yaml:"config"`  // путь к конфигу (configs/k8s-strict.yaml)
}

// Config основная структура конфигурации
type Config struct {
	Rules         ValidationRules `yaml:"rules"`
	FileProfiles  []FileProfile   `yaml:"file_profiles"` // автовыбор профиля по имени файла (4.5)
}

// MatchFileProfile проверяет, подходит ли путь к файлу под маску.
// Поддерживает: "**/X/**" — путь содержит /X/, "*name*" — filepath.Match по имени.
func MatchFileProfile(pattern, filePath string) bool {
	normalized := filepath.ToSlash(filepath.Clean(filePath))
	base := filepath.Base(filePath)

	if strings.HasPrefix(pattern, "**/") && strings.HasSuffix(pattern, "/**") {
		mid := pattern[3 : len(pattern)-3]
		return strings.Contains(normalized, "/"+mid+"/") || strings.HasPrefix(normalized, mid+"/")
	}
	if strings.HasPrefix(pattern, "**/") {
		suffix := pattern[3:]
		if ok, _ := filepath.Match(suffix, base); ok {
			return true
		}
		if strings.Contains(normalized, "/"+suffix) || strings.HasPrefix(normalized, suffix) {
			return true
		}
	}
	ok, _ := filepath.Match(pattern, base)
	return ok
}

// ConfigForFile возвращает конфиг для файла с учётом file_profiles.
// Первое совпадение pattern wins.
func ConfigForFile(base *Config, filePath string) *Config {
	if base == nil || len(base.FileProfiles) == 0 {
		return base
	}
	for _, p := range base.FileProfiles {
		if p.Pattern == "" || p.Config == "" {
			continue
		}
		if MatchFileProfile(p.Pattern, filePath) {
			if cfg, err := LoadConfig(p.Config); err == nil {
				return cfg
			}
		}
	}
	return base
}

// LoadConfig загружает конфигурацию из файла
func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	return &cfg, err
}

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() *Config {
	return &Config{
		Rules: ValidationRules{
			CheckSyntax:       true,
			CheckDuplicates:   true,
			CheckIntegrity:    true,
			CheckCommonErrors: true,
			Style:            StyleOptions{},
			JsonSchema:       JsonSchemaOptions{Enabled: false},
			K8sSchema:        K8sSchemaOptions{Enabled: false, Version: "master", IgnoreMissing: true},
			RequiredFields:    []string{"apiVersion", "kind", "metadata.name"},
			MaxLineLength:     200,
			SensitivePatterns: []string{"password", "secret", "token", "key"},
		},
	}
}
