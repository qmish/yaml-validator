package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"yaml-validator/internal/config"
	"yaml-validator/internal/fixer"
	"yaml-validator/internal/logger"
	"yaml-validator/internal/reporter"
	"yaml-validator/internal/validator"
	_ "yaml-validator/internal/validator/plugins" // регистрация плагинов
)

var version = "1.0.0" // переопределяется при сборке: -ldflags "-X yaml-validator/cmd.version=v1.0.0"

var (
	configPath string
	outputFmt  string
	verbose    bool
	logJSON    bool
	quiet      bool
	fix        bool
)

var rootCmd = &cobra.Command{
	Use:   "yaml-validator",
	Short: "YAML validation tool for DevOps",
	Long:  "A CLI tool for validating YAML files: syntax, duplicates, integrity, and common errors",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Коды выхода: 0 — OK, 1 — есть ошибки, 2 — только предупреждения (для будущего 5.7)
const (
	ExitOK       = 0
	ExitErrors   = 1
	ExitWarnings = 2
)

var validateCmd = &cobra.Command{
	Use:   "validate [file]",
	Short: "Validate YAML file",
	Long:  "Validate YAML files. Exit codes: 0 — OK, 1 — validation errors, 2 — only warnings (reserved for future)",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if verbose {
			logger.SetLevel(logger.LevelDebug)
		}
		if logJSON {
			logger.SetJSON(true)
		}

		baseCfg := loadConfig()
		logger.Debug("Loaded config: check_syntax=%v, check_duplicates=%v", baseCfg.Rules.CheckSyntax, baseCfg.Rules.CheckDuplicates)
		exitCode := ExitOK
		hasErrors := false
		totalErrors := 0
		totalWarnings := 0
		var sarifResults []reporter.FileResult
		var gitlabResults []reporter.FileResult
		var checkstyleResults []reporter.FileResult

		for _, file := range args {
			logger.Debug("Validating file: %s", file)
			absPath, err := filepath.Abs(file)
			if err != nil {
				absPath = file
			}

			cfg := config.ConfigForFile(baseCfg, absPath)
			if fix {
				// При --fix всегда применяем fixable-правила (trailing spaces, newline at EOF, consecutive empty lines)
				fixCfg := *cfg
				fixCfg.Rules.Style.ForbidTrailingSpaces = true
				fixCfg.Rules.Style.RequireNewlineAtEof = true
				fixCfg.Rules.Style.ForbidConsecutiveEmptyLines = true
				fixRes, fixErr := fixer.FixFile(absPath, &fixCfg)
				if fixErr != nil {
					logger.Error("Fix failed for %s: %v", file, fixErr)
					fmt.Fprintf(os.Stderr, "Fix failed for %s: %v\n", file, fixErr)
					hasErrors = true
					continue
				}
				if fixRes.Modified && !quiet {
					fmt.Printf("Fixed %s: %v\n", file, fixRes.Applied)
				}
			}

			errors, err := validator.Validate(absPath, cfg)
			if err != nil {
				logger.Error("Validation failed for %s: %v", file, err)
				fmt.Fprintf(os.Stderr, "Error validating %s: %v\n", file, err)
				hasErrors = true
				continue
			}
			logger.Debug("File %s: %d issue(s)", file, len(errors))

			// 5.7: считаем errors и warnings
			for _, e := range errors {
				if e.Severity == "warning" {
					totalWarnings++
				} else {
					totalErrors++
				}
			}

			if outputFmt == "sarif" || outputFmt == "gitlab" || outputFmt == "checkstyle" {
				switch outputFmt {
				case "sarif":
					sarifResults = append(sarifResults, reporter.FileResult{File: absPath, Errors: errors})
				case "gitlab":
					gitlabResults = append(gitlabResults, reporter.FileResult{File: file, Errors: errors})
				case "checkstyle":
					checkstyleResults = append(checkstyleResults, reporter.FileResult{File: absPath, Errors: errors})
				}
			}
			if !quiet {
				if outputFmt != "sarif" && outputFmt != "gitlab" && outputFmt != "checkstyle" {
					switch outputFmt {
					case "json":
						report, _ := reporter.GenerateJSONReport(file, errors)
						fmt.Println(string(report))
					case "junit":
						report, _ := reporter.GenerateJUnitReport(file, errors)
						fmt.Println(string(report))
					case "compact":
						reporter.PrintCompact(file, errors)
					case "github-annotations":
						reporter.PrintGitHubAnnotations(file, errors)
					case "severity":
						reporter.PrintSeverity(file, errors)
					default:
						reporter.PrintHumanReadable(file, errors)
					}
				}
			}

			if totalErrors > 0 {
				hasErrors = true
			}
		}

		// 5.7: exit code 0 OK, 1 errors (или parse/fix failure), 2 только warnings
		if hasErrors || totalErrors > 0 {
			exitCode = ExitErrors
		} else if totalWarnings > 0 {
			exitCode = ExitWarnings
		}

		if outputFmt == "sarif" && len(sarifResults) > 0 {
			report, _ := reporter.GenerateSARIFReport(version, sarifResults)
			fmt.Println(string(report))
		}
		if outputFmt == "gitlab" && len(gitlabResults) > 0 {
			report, _ := reporter.GenerateGitLabCodeQualityReport(gitlabResults)
			fmt.Println(string(report))
		}
		if outputFmt == "checkstyle" && len(checkstyleResults) > 0 {
			report, _ := reporter.GenerateCheckstyleReport(checkstyleResults)
			fmt.Println(string(report))
		}

		if quiet {
			// В режиме quiet — только итог (OK / N errors / M warnings); для машинных форматов итог не выводим
			machineFormat := outputFmt == "sarif" || outputFmt == "gitlab" || outputFmt == "checkstyle" || outputFmt == "json" || outputFmt == "junit"
			if !machineFormat {
				if totalErrors > 0 || totalWarnings > 0 {
					parts := []string{}
					if totalErrors > 0 {
						if totalErrors == 1 {
							parts = append(parts, "1 error")
						} else {
							parts = append(parts, fmt.Sprintf("%d errors", totalErrors))
						}
					}
					if totalWarnings > 0 {
						if totalWarnings == 1 {
							parts = append(parts, "1 warning")
						} else {
							parts = append(parts, fmt.Sprintf("%d warnings", totalWarnings))
						}
					}
					fmt.Println(strings.Join(parts, ", "))
				} else {
					fmt.Println("OK")
				}
			}
		}

		os.Exit(exitCode) // 0 OK, 1 errors, 2 warnings (reserved)
	},
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Create default configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		content := `rules:
  check_syntax: true
  check_duplicates: true
  check_integrity: true
  check_common_errors: true
  check_key_ordering: false
  key_order: []
  max_key_name_length: 0
  forbid_default_values: false
  default_values: {}
  unique_list_fields: []
  key_value_patterns: []
  rule_severity: {}
  inline_ignore: false
  style:
    require_document_start: false
    forbid_trailing_spaces: false
    forbid_trailing_dots: false
    require_newline_at_eof: false
    forbid_consecutive_empty_lines: false
    forbid_tabs: false
    forbid_unicode: false
    forbid_bom: false
    require_empty_line_between_blocks: false
    min_empty_lines_between_blocks: 0
    require_document_end: false
    require_comments_indented: false
    require_quoted_keys: false
    require_quoted_values: false
    indent_spaces: 0
  json_schema:
    enabled: false
    schema_path: ""
    schema_cache_dir: ""
  k8s_schema:
    enabled: false
    version: master
    strict: false
    cache_dir: ""
    ignore_missing_schemas: true
  required_fields:
    - apiVersion
    - kind
    - metadata.name
  max_line_length: 200
  sensitive_patterns:
    - password
    - secret
    - token
    - key
`
		filename := "yaml-validator.yaml"
		if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created configuration file: %s\n", filename)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Run: func(c *cobra.Command, args []string) {
		fmt.Fprintf(c.OutOrStdout(), "yaml-validator version %s\n", version)
	},
}

// RuleDescriptor — правило для вывода rules list (машиночитаемый список)
type RuleDescriptor struct {
	ID          string `json:"id" yaml:"id"`
	Description string `json:"description" yaml:"description"`
	ConfigKey   string `json:"config_key" yaml:"config_key"`
}

var builtinRules = []RuleDescriptor{
	{ID: "check_syntax", Description: "Проверка синтаксиса YAML", ConfigKey: "rules.check_syntax"},
	{ID: "check_duplicates", Description: "Запрет дублирующихся ключей", ConfigKey: "rules.check_duplicates"},
	{ID: "check_integrity", Description: "Наличие обязательных полей (required_fields)", ConfigKey: "rules.check_integrity"},
	{ID: "check_common_errors", Description: "Табы, длина строки, чувствительные данные", ConfigKey: "rules.check_common_errors"},
	{ID: "check_key_ordering", Description: "Порядок ключей (алфавитный)", ConfigKey: "rules.check_key_ordering"},
	{ID: "key_order", Description: "Настраиваемый порядок ключей (key_order)", ConfigKey: "rules.key_order"},
	{ID: "max_key_name_length", Description: "Максимальная длина имени ключа", ConfigKey: "rules.max_key_name_length"},
	{ID: "forbid_default_values", Description: "Запрет ключей со значением по умолчанию (default_values)", ConfigKey: "rules.forbid_default_values"},
	{ID: "unique_list_fields", Description: "Уникальность элементов массива по полю (path, field)", ConfigKey: "rules.unique_list_fields"},
	{ID: "key_value_patterns", Description: "Регулярки для ключей/значений (path, pattern, target)", ConfigKey: "rules.key_value_patterns"},
	{ID: "rule_severity", Description: "Severity правил: Type -> error|warning (5.7)", ConfigKey: "rules.rule_severity"},
	{ID: "inline_ignore", Description: "Комментарии для отключения правил в YAML", ConfigKey: "rules.inline_ignore"},
	{ID: "require_document_start", Description: "Требовать --- в начале документа", ConfigKey: "rules.style.require_document_start"},
	{ID: "forbid_trailing_spaces", Description: "Запрет пробелов в конце строки", ConfigKey: "rules.style.forbid_trailing_spaces"},
	{ID: "forbid_trailing_dots", Description: "Запрет точек в конце строки", ConfigKey: "rules.style.forbid_trailing_dots"},
	{ID: "require_newline_at_eof", Description: "Перевод строки в конце файла", ConfigKey: "rules.style.require_newline_at_eof"},
	{ID: "forbid_consecutive_empty_lines", Description: "Запрет нескольких пустых строк подряд", ConfigKey: "rules.style.forbid_consecutive_empty_lines"},
	{ID: "require_empty_line_between_blocks", Description: "Пустая строка между топ-уровневыми блоками", ConfigKey: "rules.style.require_empty_line_between_blocks"},
	{ID: "min_empty_lines_between_blocks", Description: "Минимум пустых строк между блоками", ConfigKey: "rules.style.min_empty_lines_between_blocks"},
	{ID: "require_document_end", Description: "Требовать ... в конце документа", ConfigKey: "rules.style.require_document_end"},
	{ID: "require_comments_indented", Description: "Комментарии с отступом блока", ConfigKey: "rules.style.require_comments_indented"},
	{ID: "require_quoted_keys", Description: "Ключи в кавычках", ConfigKey: "rules.style.require_quoted_keys"},
	{ID: "require_quoted_values", Description: "Строковые значения в кавычках", ConfigKey: "rules.style.require_quoted_values"},
	{ID: "indent_spaces", Description: "Шаг отступов (2 или 4 пробела)", ConfigKey: "rules.style.indent_spaces"},
	{ID: "forbid_tabs", Description: "Запрет табуляции", ConfigKey: "rules.style.forbid_tabs"},
	{ID: "forbid_unicode", Description: "Запрет не-ASCII символов", ConfigKey: "rules.style.forbid_unicode"},
	{ID: "forbid_bom", Description: "Запрет BOM в начале файла", ConfigKey: "rules.style.forbid_bom"},
	{ID: "max_line_length", Description: "Максимальная длина строки", ConfigKey: "rules.max_line_length"},
	{ID: "json_schema", Description: "Валидация по JSON Schema", ConfigKey: "rules.json_schema"},
	{ID: "k8s_schema", Description: "Валидация по схеме Kubernetes", ConfigKey: "rules.k8s_schema"},
	{ID: "docker_compose", Description: "Плагин: image или build у сервисов", ConfigKey: "plugins"},
}

var rulesListOutputFmt string

var rulesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List validation rules (machine-readable)",
	Long:  "Output list of rules in JSON or YAML for scripts and documentation.",
	RunE: func(cmd *cobra.Command, args []string) error {
		switch rulesListOutputFmt {
		case "json":
			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "  ")
			return enc.Encode(builtinRules)
		case "yaml":
			return yaml.NewEncoder(cmd.OutOrStdout()).Encode(builtinRules)
		default:
			return yaml.NewEncoder(cmd.OutOrStdout()).Encode(builtinRules)
		}
	},
}

var rulesCmd = &cobra.Command{
	Use:   "rules",
	Short: "List and describe validation rules",
	Long:  "Subcommands: list — output rules in JSON/YAML for scripts and docs.",
}

func loadConfig() *config.Config {
	paths := []string{configPath}
	if configPath == "" {
		paths = []string{
			"yaml-validator.yaml",
			".yaml-validator.yaml",
			"configs/default.yaml",
		}
	}

	for _, p := range paths {
		if p == "" {
			continue
		}
		if cfg, err := config.LoadConfig(p); err == nil {
			return cfg
		}
	}

	return config.DefaultConfig()
}

func init() {
	validateCmd.Flags().StringVarP(&outputFmt, "output", "o", "human", "Output format: human, json, junit, sarif, checkstyle, compact, gitlab, github-annotations, severity")
	validateCmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to configuration file")
	validateCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose (debug) logging")
	validateCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Minimal output: only OK or N errors")
	validateCmd.Flags().BoolVar(&fix, "fix", false, "Auto-fix: trailing spaces, newline at EOF, consecutive empty lines")
	validateCmd.Flags().BoolVar(&logJSON, "log-json", false, "Output logs in JSON format (for ELK, Loki)")

	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Configuration management",
	}
	configCmd.AddCommand(configInitCmd)

	rulesListCmd.Flags().StringVarP(&rulesListOutputFmt, "output", "o", "yaml", "Output format: json, yaml")

	rulesCmd.AddCommand(rulesListCmd)

	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(rulesCmd)
	rootCmd.AddCommand(versionCmd)
}

// Execute запускает корневую команду
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
