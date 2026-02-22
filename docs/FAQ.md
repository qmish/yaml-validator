# Часто задаваемые вопросы (FAQ)

## Общее

### Чем yaml-validator отличается от yamllint?

yaml-validator — инструмент на Go с фокусом на DevOps: валидация Kubernetes-манифестов, docker-compose, интеграция с GitHub/GitLab CI, SARIF, pre-commit. yamllint — Python-инструмент с широким набором правил стиля. Подробное сравнение: [docs/COMPARISON.md](COMPARISON.md). Таблица миграции правил yamllint → yaml-validator: [docs/YAMLLINT_MIGRATION.md](YAMLLINT_MIGRATION.md).

### Поддерживается ли мультидокументный YAML (---)?

Да, начиная с v1.52.0. Каждый документ в файле валидируется отдельно. Ошибки содержат `DocumentIndex` (file#doc2 и т.п.) для указания номера документа.

### Как использовать yaml-validator как линтер в IDE?

См. [docs/ide.md](ide.md) — примеры для VS Code (tasks, Run on Save), JetBrains (External Tools), Sublime, Vim.

---

## Конфигурация

### Где искать конфиг по умолчанию?

По порядку: `yaml-validator.yaml`, `.yaml-validator.yaml`, `configs/default.yaml`. Если ни один не найден — используется встроенная конфигурация по умолчанию.

### Как применить разный конфиг к разным файлам?

Используйте `file_profiles` в yaml-validator.yaml:

```yaml
file_profiles:
  - pattern: "**/k8s/**"
    config: "configs/k8s-strict.yaml"
  - pattern: "*docker-compose*.yaml"
    config: "configs/docker-compose.yaml"
```

### Как отключить правило для одной строки?

Inline-игнор через комментарий:
- `# yaml-validator disable-line rule:LineTooLong`
- `# yaml-validator disable-next-line rule:TrailingSpaces`

---

## Валидация

### Exit codes — что они означают?

- `0` — всё в порядке (нет ошибок и предупреждений)
- `1` — есть ошибки (или ошибка парсинга/фикса)
- `2` — только предупреждения (reserved)

### Как включить проверку по Kubernetes OpenAPI схеме?

Конфиг `configs/k8s-strict.yaml` или:

```yaml
rules:
  k8s_schema:
    enabled: true
    version: master  # или конкретная версия: 1.28, 1.29
    ignore_missing_schemas: true
```

Требуется [kubeconform](https://github.com/yannh/kubeconform) (автоматически подтягивается как зависимость).

### Как валидировать по произвольной JSON Schema?

```yaml
rules:
  json_schema:
    enabled: true
    schema_path: "path/to/schema.json"
    schema_cache_dir: ".schema-cache"  # для schema_path: https://...
```

---

## CI/CD

### Как использовать в GitHub Actions?

```yaml
- uses: actions/checkout@v4
- uses: qmish/yaml-validator-action@v1  # или скачать бинарник
- run: yaml-validator validate "**/*.yml" "**/*.yaml"
```

Для Code Scanning (SARIF): [docs/code-scanning-example.yml](code-scanning-example.yml).

### Как добавить pre-commit хук?

```yaml
# .pre-commit-config.yaml
repos:
  - repo: https://github.com/qmish/yaml-validator
    rev: v1.52.0
    hooks:
      - id: yaml-validator
```

Подробнее: [docs/pre-commit.md](pre-commit.md).

---

## Разработка

### Как запустить линтер и тесты?

```bash
make lint          # go vet + golangci-lint
make test          # go test
make test-coverage # покрытие, coverage.html
```

### Как добавить свой плагин?

См. [docs/PLUGINS.md](PLUGINS.md). Интерфейс: `Name() string`, `Validate(node *yaml.Node) []Error`. Регистрация через `validator.RegisterPlugin`.
