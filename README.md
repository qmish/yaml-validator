# YAML Validator

[![CI](https://github.com/qmish/yaml-validator/actions/workflows/yaml-validation.yml/badge.svg)](https://github.com/qmish/yaml-validator/actions/workflows/yaml-validation.yml)
[![Release](https://img.shields.io/github/v/release/qmish/yaml-validator)](https://github.com/qmish/yaml-validator/releases)
[![Go 1.21+](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev/)

Инструмент для проверки YAML-файлов на Go, предназначенный для DevOps.

## Возможности

- **Валидация синтаксиса** — проверка корректности YAML, пустые документы
- **Поиск дубликатов** — дублирующиеся ключи на одном уровне
- **Логическая целостность** — обязательные поля (apiVersion, kind, metadata.name для Kubernetes), типы данных
- **Распространённые ошибки** — табуляции вместо пробелов, слишком длинные строки (>200 символов), чувствительные данные (password, token и т.д.)
- **Плагины** — расширяемость через `ValidatorPlugin` (встроенный плагин для Kubernetes)
- **Логирование** — verbose режим, JSON-логи для ELK/Loki

## Установка

```bash
go build -o yaml-validator .
# или (с Make): make build
```

## Использование

```bash
# Проверка файла (человекочитаемый вывод)
yaml-validator validate config.yaml

# Вывод в формате JSON
yaml-validator validate config.yaml -o json

# Вывод в формате JUnit (для Jenkins, GitLab CI)
yaml-validator validate config.yaml -o junit

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

- **K8s-манифесты** — используйте конфиг по умолчанию (проверка apiVersion, kind, metadata.name).
- **docker-compose и другой YAML** — используйте `-c configs/docker-compose.yaml` (без обязательных K8s-полей).

```yaml
rules:
  check_syntax: true
  check_duplicates: true
  check_integrity: true
  check_common_errors: true
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

Встроенный плагин `kubernetes` проверяет K8s-манифесты (apiVersion, kind, metadata).

## Тестирование

```bash
go test ./... -v
# Короткий прогон: go test ./... -short
# С покрытием: make test-coverage
```

См. [CONTRIBUTING.md](CONTRIBUTING.md) для участия в разработке.

## Интеграция с CI/CD

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

Проект включает workflow `.github/workflows/yaml-validation.yml` с матрицей Go 1.21/1.22 и сборкой Docker.

```yaml
- name: Validate YAML
  run: |
    ./yaml-validator validate **/*.yml **/*.yaml
```

### Pre-commit hook

```bash
cp .githooks/pre-commit .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
```

### Docker

```bash
# Сборка образа
docker build -t yaml-validator .

# Запуск с монтированием директории
docker run --rm -v $(pwd):/workspace yaml-validator validate config.yaml

# Docker Compose (профиль test)
docker compose --profile test run --rm validator
```
