# Миграция с yamllint на yaml-validator

Таблица соответствия правил и пример конвертации конфигурации.

---

## Таблица соответствия правил

| yamllint | yaml-validator | Примечание |
|----------|----------------|------------|
| `document-start` | `rules.style.require_document_start: true` | Требовать `---` в начале |
| `document-end` | `rules.style.require_document_end: true` | Требовать `...` в конце |
| `trailing-spaces` | `rules.style.forbid_trailing_spaces: true` | Автофикс: `--fix` |
| `line-length` | `rules.max_line_length: N` | По умолчанию 200 (yamllint: 80) |
| `indentation` | `rules.style.indent_spaces: 2` или `4` | Шаг отступов; табы: `forbid_tabs` |
| `key-duplicates` | `rules.check_duplicates: true` | Включено по умолчанию |
| `empty-lines` | `rules.style.forbid_consecutive_empty_lines: true` | Автофикс: `--fix` |
| `comments` (indentation) | `rules.style.require_comments_indented: true` | Отступ комментариев |
| `key-ordering` | `rules.check_key_ordering: true` или `rules.key_order: [...]` | `key_order` — приоритетный порядок |
| `quoted-strings` | `rules.style.require_quoted_keys: true` | Ключи в кавычках |
| `quoted-values` | `rules.style.require_quoted_values: true` | Строковые значения в кавычках |
| tabs (в indentation) | `rules.style.forbid_tabs: true` | Отдельное правило |
| new-line-at-end-of-file | `rules.style.require_newline_at_eof: true` | Автофикс: `--fix` |
| — | `rules.check_integrity` + `required_fields` | Нет в yamllint: apiVersion, kind, metadata.name |
| — | `rules.sensitive_patterns` | Нет в yamllint: password, token, secret |
| — | `rules.check_common_errors` | Табы, длина строки, чувствительные данные |

### Правила yamllint без прямого соответствия

- `anchors` (aliases, anchors) — yaml-validator не проверяет
- `braces`, `brackets` — flow mappings `{}`, `[]`
- `colons` — пробелы после `:`
- `hyphens` — пробелы после `-`
- `truthy` — значения true/false/yes/no
- `comments` (require-starting-space и др.) — частично

---

## Пример конвертации

### .yamllint (yamllint)

```yaml
extends: default

rules:
  document-start: enable
  document-end: enable
  line-length:
    max: 120
  indentation:
    spaces: 2
    indent-sequences: true
  trailing-spaces: enable
  key-duplicates: enable
  empty-lines:
    max: 1
  comments:
    require-starting-space: true
  quoted-strings:
    quote-type: any
  quoted-values:
    quote-type: any
```

### yaml-validator.yaml

```yaml
rules:
  check_syntax: true
  check_duplicates: true
  check_integrity: true
  check_common_errors: true
  max_line_length: 120
  style:
    require_document_start: true
    require_document_end: true
    forbid_trailing_spaces: true
    forbid_consecutive_empty_lines: true
    require_comments_indented: true
    require_quoted_keys: true
    require_quoted_values: true
    indent_spaces: 2
    forbid_tabs: true
    require_newline_at_eof: true
```

---

## Команды

| yamllint | yaml-validator |
|----------|----------------|
| `yamllint .` | `yaml-validator validate "**/*.yaml" "**/*.yml"` |
| `yamllint -f parsable file.yaml` | `yaml-validator validate file.yaml -o compact` |
| `yamllint -f github file.yaml` | `yaml-validator validate file.yaml -o github-annotations` |

---

## pre-commit

yamllint:
```yaml
- repo: https://github.com/adrienverge/yamllint
  rev: v1.35.1
  hooks:
    - id: yamllint
```

yaml-validator:
```yaml
- repo: https://github.com/qmish/yaml-validator
  rev: v1.41.0
  hooks:
    - id: yaml-validator
```
