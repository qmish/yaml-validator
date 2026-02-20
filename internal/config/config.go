package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

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
	RequireDocumentStart bool `yaml:"require_document_start"` // требовать --- в начале файла
	ForbidTrailingSpaces bool `yaml:"forbid_trailing_spaces"` // запрет пробелов в конце строки
	RequireNewlineAtEof  bool `yaml:"require_newline_at_eof"` // требовать перевод строки в конце файла
	ForbidConsecutiveEmptyLines bool `yaml:"forbid_consecutive_empty_lines"` // запрет более одной пустой строки подряд
	RequireDocumentEnd  bool `yaml:"require_document_end"` // требовать ... в конце файла (много-документный YAML)
}

// ValidationRules определяет правила валидации
type ValidationRules struct {
	CheckSyntax       bool              `yaml:"check_syntax"`
	CheckDuplicates   bool              `yaml:"check_duplicates"`
	CheckIntegrity    bool              `yaml:"check_integrity"`
	CheckCommonErrors bool              `yaml:"check_common_errors"`
	InlineIgnore     bool              `yaml:"inline_ignore"` // разрешить отключение правил через комментарии в YAML
	Style             StyleOptions      `yaml:"style"`
	K8sSchema        K8sSchemaOptions   `yaml:"k8s_schema"`
	RequiredFields    []string          `yaml:"required_fields"`
	MaxLineLength     int               `yaml:"max_line_length"`
	SensitivePatterns []string          `yaml:"sensitive_patterns"`
}

// Config основная структура конфигурации
type Config struct {
	Rules ValidationRules `yaml:"rules"`
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
			K8sSchema:        K8sSchemaOptions{Enabled: false, Version: "master", IgnoreMissing: true},
			RequiredFields:    []string{"apiVersion", "kind", "metadata.name"},
			MaxLineLength:     200,
			SensitivePatterns: []string{"password", "secret", "token", "key"},
		},
	}
}
