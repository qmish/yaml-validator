package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"yaml-validator/internal/config"
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

		cfg := loadConfig()
		logger.Debug("Loaded config: check_syntax=%v, check_duplicates=%v", cfg.Rules.CheckSyntax, cfg.Rules.CheckDuplicates)
		exitCode := ExitOK
		hasErrors := false
		var sarifResults []reporter.FileResult
		var gitlabResults []reporter.FileResult
		var checkstyleResults []reporter.FileResult

		for _, file := range args {
			logger.Debug("Validating file: %s", file)
			absPath, err := filepath.Abs(file)
			if err != nil {
				absPath = file
			}

			errors, err := validator.Validate(absPath, cfg)
			if err != nil {
				logger.Error("Validation failed for %s: %v", file, err)
				fmt.Fprintf(os.Stderr, "Error validating %s: %v\n", file, err)
				hasErrors = true
				continue
			}
			logger.Debug("File %s: %d error(s)", file, len(errors))

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

			if len(errors) > 0 {
				hasErrors = true
			}
		}

		if hasErrors {
			exitCode = ExitErrors
		}
		// ExitWarnings (2) — для будущего, когда будут предупреждения (5.7)

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
	validateCmd.Flags().BoolVar(&logJSON, "log-json", false, "Output logs in JSON format (for ELK, Loki)")

	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Configuration management",
	}
	configCmd.AddCommand(configInitCmd)

	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(versionCmd)
}

// Execute запускает корневую команду
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
