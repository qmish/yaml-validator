# YAML Validator

[![CI](https://github.com/qmish/yaml-validator/actions/workflows/yaml-validation.yml/badge.svg)](https://github.com/qmish/yaml-validator/actions/workflows/yaml-validation.yml)
[![Release](https://img.shields.io/github/v/release/qmish/yaml-validator)](https://github.com/qmish/yaml-validator/releases)
[![Go 1.24+](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)](https://go.dev/)

Инструмент для проверки YAML-файлов на Go, предназначенный для DevOps.

## Возможности

- **Валидация синтаксиса** — проверка корректности YAML, пустые документы
- **Поиск дубликатов** — дублирующиеся ключи на одном уровне
- **Логическая целостность** — обязательные поля (apiVersion, kind, metadata.name для Kubernetes), типы данных
- **Распространённые ошибки** — табуляции вместо пробелов, слишком длинные строки (>200 символов), чувствительные данные (password, token и т.д.)
- **Плагины** — расширяемость через `ValidatorPlugin` (встроенный плагин для Kubernetes)
- **K8s по схеме** — опциональная проверка по полной OpenAPI/JSON Schema Kubernetes (типы полей, все ресурсы), через [kubeconform](https://github.com/yannh/kubeconform)
- **Произвольная JSON Schema** — опциональная валидация по любой JSON Schema (не только K8s): `rules.json_schema.enabled: true`, `rules.json_schema.schema_path: путь/к/схеме.json`
- **Правила стиля** — document-start (`---`), document-end (`...`), запрет пробелов и точек в конце строки, перевод строки в конце файла, запрет нескольких пустых строк подряд, отступ у комментариев, кавычки для ключей и значений; порядок ключей (`check_key_ordering` или `key_order`), длина имён ключей (`max_key_name_length`)
- **Inline-игнор** — отключение правил через комментарии: `# yaml-validator disable-line rule:LineTooLong`, `# yaml-validator disable-next-line rule:TrailingSpaces`
- **Логирование** — verbose режим, JSON-логи для ELK/Loki

## Установка

```bash
go build -o yaml-validator .
# или: make build
# Linux/macOS кросс-сборка: ./build.sh
# Windows: .\build.ps1
```

## Использование

```bash
# Проверка файла (человекочитаемый вывод)
yaml-validator validate config.yaml

# Вывод в формате JSON
yaml-validator validate config.yaml -o json

# Вывод в формате JUnit (для Jenkins, GitLab CI)
yaml-validator validate config.yaml -o junit

# Вывод в формате SARIF (для GitHub Code Scanning)
yaml-validator validate config.yaml -o sarif > results.sarif

# Вывод в формате GitLab Code Quality (gl-code-quality-report.json)
yaml-validator validate config.yaml -o gitlab > gl-code-quality-report.json

# Компактный вывод file:line[:col]: message (для редакторов, парсеров; колонка выводится, когда известна)
yaml-validator validate config.yaml -o compact

# GitHub Actions annotations (::error file=...,line=...::)
yaml-validator validate config.yaml -o github-annotations

# Текстовый формат с severity [ERROR] file:line: message (для скриптов, CI)
yaml-validator validate config.yaml -o severity

# Автофикс (trailing spaces, newline at EOF, consecutive empty lines)
yaml-validator validate config.yaml --fix

# Использование конфигурационного файла
yaml-validator validate config.yaml -c configs/default.yaml

# Создание конфигурации по умолчанию
yaml-validator config init

# Версия
yaml-validator version

# Подробный лог (отладка)
yaml-validator validate config.yaml -v

# JSON-логи для CI (ELK, Loki)
yaml-validator validate config.yaml --log-json
```

## Конфигурация

По умолчанию: `configs/default.yaml` или `yaml-validator.yaml`.

- **K8s-манифесты** — конфиг по умолчанию (проверка apiVersion, kind, metadata.name). Для порядка ключей K8s (metadata перед spec): `-c configs/k8s-key-order.yaml`. Для полной схемы: `-c configs/k8s-strict.yaml`.
- **docker-compose и другой YAML** — `-c configs/docker-compose.yaml` (без обязательных K8s-полей).
- **Произвольная JSON Schema** — `-c configs/json-schema.yaml` с путём к схеме в `rules.json_schema.schema_path`.

```yaml
rules:
  check_syntax: true
  check_duplicates: true
  check_integrity: true
  check_common_errors: true
  check_key_ordering: false
  key_order: []
  max_key_name_length: 0
  required_fields:
    - apiVersion
    - kind
    - metadata.name
  inline_ignore: false
  style:
    require_document_start: false
    forbid_trailing_spaces: false
    require_newline_at_eof: false
    forbid_consecutive_empty_lines: false
    require_document_end: false
    require_comments_indented: false
    require_quoted_keys: false
    require_quoted_values: false
    forbid_trailing_dots: false
    indent_spaces: 0
  k8s_schema:
    enabled: false
    version: master
    ignore_missing_schemas: true
  max_line_length: 200
  sensitive_patterns:
    - password
    - secret
    - token
    - key
```

## Структура проекта

```
yaml-validator/
├── cmd/validator.go          # CLI с Cobra
├── internal/
│   ├── config/              # Конфигурация
│   ├── logger/              # Логирование
│   ├── parser/              # Парсер YAML
│   ├── validator/           # Валидаторы
│   │   └── plugins/         # Плагины (Kubernetes и др.)
│   └── reporter/            # Генератор отчётов
├── pkg/types.go             # Общие типы
├── configs/default.yaml     # Конфиг по умолчанию
├── build.ps1                # Кросс-платформенная сборка
├── release.ps1              # Скрипт релиза
└── testdata/                # Тестовые данные
```

## Расширяемость (плагины)

Интерфейс плагина:

```go
type ValidatorPlugin interface {
    Name() string
    Validate(node *yaml.Node) []Error
}

func RegisterPlugin(name string, plugin ValidatorPlugin)
```

Встроенные плагины: `kubernetes` (apiVersion, kind, metadata), `docker-compose` (image/build у сервисов).

## Тестирование

```bash
go test ./... -v
# Короткий прогон: go test ./... -short
# С покрытием: make test-coverage
```

См. [CONTRIBUTING.md](CONTRIBUTING.md) для участия в разработке. [Сравнение с другими инструментами](docs/COMPARISON.md) (yamllint, kubeval, kubeconform и др.).

## Интеграция с CI/CD

### GitHub Code Scanning (SARIF)

См. [docs/code-scanning-example.yml](docs/code-scanning-example.yml) — полный пример workflow.

```yaml
- name: Validate YAML and upload SARIF
  run: |
    ./yaml-validator validate **/*.yaml **/*.yml -o sarif > results.sarif
- uses: github/codeql-action/upload-sarif@v3
  with:
    sarif_file: results.sarif
```

### GitLab CI (Code Quality)

```yaml
validate-yaml:
  script:
    - yaml-validator validate "**/*.yml" "**/*.yaml" -o gitlab > gl-code-quality-report.json
  artifacts:
    reports:
      codequality: gl-code-quality-report.json
```

### Git pre-commit hook

```bash
#!/bin/bash
for file in $(git diff --cached --name-only | grep -E '\.(yml|yaml)$'); do
    yaml-validator validate "$file"
    if [ $? -ne 0 ]; then
        echo "Validation failed for $file"
        exit 1
    fi
done
```

### GitHub Actions

Проект включает workflow `.github/workflows/yaml-validation.yml` (Go 1.24) и сборку Docker.

```yaml
- name: Validate YAML
  run: |
    ./yaml-validator validate **/*.yml **/*.yaml
```

### GitLab CI

```yaml
validate-yaml:
  image: ghcr.io/qmish/yaml-validator:latest
  script:
    - yaml-validator validate "**/*.yml" "**/*.yaml"
  rules:
    - if: $CI_PIPELINE_SOURCE == "merge_request_event"
```

Или с бинарником (скачать релиз):

```yaml
validate-yaml:
  stage: test
  script:
    - curl -sSL https://github.com/qmish/yaml-validator/releases/download/v1.6.0/yaml-validator-v1.6.0-linux-amd64.tar.gz | tar xz -C /usr/local/bin
    - yaml-validator validate **/*.yml **/*.yaml
```

### Jenkins

```groovy
stage('Validate YAML') {
  steps {
    sh '''
      curl -sSL https://github.com/qmish/yaml-validator/releases/download/v1.6.0/yaml-validator-v1.6.0-linux-amd64.tar.gz | tar xz -C /tmp
      /tmp/yaml-validator validate **/*.yml **/*.yaml
    '''
  }
}
```

Или с Docker:

```groovy
stage('Validate YAML') {
  steps {
    sh 'docker run --rm -v ${WORKSPACE}:/workspace ghcr.io/qmish/yaml-validator:latest validate **/*.yml **/*.yaml'
  }
}
```

### VS Code / IDE

Использование как линтера: [docs/ide.md](docs/ide.md) — VS Code (tasks, Run on Save), JetBrains (External Tools), Sublime, Vim.

### Pre-commit (official)

yaml-validator — [official pre-commit hook](docs/pre-commit.md). Добавьте в `.pre-commit-config.yaml`:

```yaml
repos:
  - repo: https://github.com/qmish/yaml-validator
    rev: v1.35.0
    hooks:
      - id: yaml-validator
```

```bash
pre-commit install
```

### Docker

```bash
# Образ из GitHub Container Registry
docker pull ghcr.io/qmish/yaml-validator:latest

# Сборка локально
docker build -t yaml-validator .

# Запуск с монтированием директории
docker run --rm -v $(pwd):/workspace yaml-validator validate config.yaml

# Docker Compose (профиль test)
docker compose --profile test run --rm validator
```
