package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// ValidationRules определяет правила валидации
type ValidationRules struct {
	CheckSyntax       bool     `yaml:"check_syntax"`
	CheckDuplicates   bool     `yaml:"check_duplicates"`
	CheckIntegrity    bool     `yaml:"check_integrity"`
	CheckCommonErrors bool     `yaml:"check_common_errors"`
	RequiredFields    []string `yaml:"required_fields"`
	MaxLineLength     int      `yaml:"max_line_length"`
	SensitivePatterns []string `yaml:"sensitive_patterns"`
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
			RequiredFields:    []string{"apiVersion", "kind", "metadata.name"},
			MaxLineLength:     200,
			SensitivePatterns: []string{"password", "secret", "token", "key"},
		},
	}
}
