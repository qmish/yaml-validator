# Правила валидации yaml-validator

Справочник правил, опций конфигурации и примеров. Машиночитаемый список: `yaml-validator rules list -o yaml`.

---

## Основные правила

| Правило | Ключ конфига | Описание | Пример |
|--------|----------------|----------|--------|
| **check_syntax** | `rules.check_syntax` | Проверка синтаксиса YAML (парсинг, пустые документы). | Включено по умолчанию. |
| **check_duplicates** | `rules.check_duplicates` | Запрет дублирующихся ключей на одном уровне. | `key: 1` и `key: 2` в одном блоке → ошибка. |
| **check_integrity** | `rules.check_integrity` | Наличие обязательных полей из `required_fields`. | По умолчанию: apiVersion, kind, metadata.name. |
| **check_common_errors** | `rules.check_common_errors` | Табы вместо пробелов, длина строки (`max_line_length`), чувствительные данные (password, secret, token, key). | `max_line_length: 200`. |
| **check_key_ordering** | `rules.check_key_ordering` | Требовать алфавитный порядок ключей. | Для K8s часто отключают в пользу `key_order`. |
| **key_order** | `rules.key_order` | Приоритетный порядок ключей (не алфавитный). | `[apiVersion, kind, metadata, spec]`. |
| **max_key_name_length** | `rules.max_key_name_length` | Максимальная длина имени ключа (0 = отключено). | `max_key_name_length: 50`. |
| **forbid_default_values** | `rules.forbid_default_values` | Запрет ключей со значением по умолчанию (K8s: imagePullPolicy: Always и т.п.). | `default_values: {imagePullPolicy: "Always"}`. |
| **default_values** | `rules.default_values` | Карта ключ → значение по умолчанию. Используется с `forbid_default_values`. | См. `configs/k8s-forbid-defaults.yaml`. |
| **unique_list_fields** | `rules.unique_list_fields` | Уникальность элементов массива по полю (path, field). | См. `configs/k8s-unique-containers.yaml`. |
| **key_value_patterns** | `rules.key_value_patterns` | Регулярки для ключей/значений (path, pattern, target: keys/values). | См. `configs/k8s-name-patterns.yaml`. |
| **inline_ignore** | `rules.inline_ignore` | Разрешить отключение правил через комментарии в YAML. | `# yaml-validator disable-line rule:LineTooLong`. |
| **required_fields** | `rules.required_fields` | Список обязательных полей (точечная нотация). | `[apiVersion, kind, metadata.name]`. |
| **max_line_length** | `rules.max_line_length` | Максимальная длина строки в символах. | `200`. |
| **sensitive_patterns** | `rules.sensitive_patterns` | Подстроки в ключах, при которых значение считается чувствительным. | `[password, secret, token, key]`. |

---

## Правила стиля (rules.style)

| Правило | Ключ | Описание | Пример |
|--------|------|----------|--------|
| **require_document_start** | `require_document_start` | Требовать `---` в начале документа. | yamllint document-start. |
| **forbid_trailing_spaces** | `forbid_trailing_spaces` | Запрет пробелов/табов в конце строки. | Автофикс: `--fix`. |
| **forbid_trailing_dots** | `forbid_trailing_dots` | Запрет точек в конце значения (`key: value.`). | |
| **require_newline_at_eof** | `require_newline_at_eof` | Требовать перевод строки в конце файла. | Автофикс: `--fix`. |
| **forbid_consecutive_empty_lines** | `forbid_consecutive_empty_lines` | Запрет более одной пустой строки подряд. | Автофикс: `--fix`. |
| **require_empty_line_between_blocks** | `require_empty_line_between_blocks` | Ровно одна пустая строка между топ-уровневыми блоками. | |
| **min_empty_lines_between_blocks** | `min_empty_lines_between_blocks` | Минимум пустых строк между блоками (0 или 1). | Гибче, чем только «ровно одна». |
| **require_document_end** | `require_document_end` | Требовать `...` в конце документа (много-документный YAML). | |
| **require_comments_indented** | `require_comments_indented` | Комментарии внутри блока должны иметь отступ блока. | |
| **require_quoted_keys** | `require_quoted_keys` | Ключи маппинга в кавычках. | |
| **require_quoted_values** | `require_quoted_values` | Строковые значения в кавычках. | yamllint quoted-values. |
| **indent_spaces** | `indent_spaces` | Шаг отступов (2 или 4 пробела; 0 = отключено). | |
| **forbid_tabs** | `forbid_tabs` | Запрет табуляции. | |
| **forbid_unicode** | `forbid_unicode` | Запрет не-ASCII в ключах и строках. | Строгие ASCII-конфиги. |
| **forbid_bom** | `forbid_bom` | Запрет BOM (Byte Order Mark) в начале файла. | |

---

## Схемы

| Правило | Ключ | Описание |
|--------|------|----------|
| **json_schema** | `rules.json_schema` | Валидация по произвольной JSON Schema. `enabled`, `schema_path`. |
| **k8s_schema** | `rules.k8s_schema` | Валидация по OpenAPI-схеме Kubernetes. `enabled`, `version`, `strict`, `cache_dir`, `ignore_missing_schemas`. |

---

## Плагины

- **docker_compose** — у каждого сервиса должно быть `image` или `build`. Включается автоматически для файлов с секцией `services` (без apiVersion).

---

## Автовыбор конфига (file_profiles)

В корневом конфиге (yaml-validator.yaml):

```yaml
file_profiles:
  - pattern: "**/k8s/**"
    config: "configs/k8s-strict.yaml"
  - pattern: "*docker-compose*.yaml"
    config: "configs/docker-compose.yaml"
```

Первое совпадение пути файла с `pattern` задаёт конфиг. Маски: `**/X/**`, `*name*.yaml`. Пример: `configs/file-profiles.yaml`.

---

## Пример конфига (фрагмент)

```yaml
rules:
  check_syntax: true
  check_duplicates: true
  check_integrity: true
  check_common_errors: true
  check_key_ordering: false
  key_order: [apiVersion, kind, metadata, spec]
  max_key_name_length: 0
  inline_ignore: true
  required_fields: [apiVersion, kind, metadata.name]
  max_line_length: 200
  style:
    require_document_start: false
    forbid_trailing_spaces: true
    forbid_trailing_dots: true
    require_newline_at_eof: true
    forbid_consecutive_empty_lines: true
    forbid_tabs: true
    forbid_unicode: false
    forbid_bom: true
  json_schema:
    enabled: false
    schema_path: ""
  k8s_schema:
    enabled: false
    version: master
    ignore_missing_schemas: true
```

---

## Типы ошибок (Type)

При валидации в отчётах фигурируют, в частности: `SyntaxError`, `DuplicateKey`, `MissingRequiredField`, `TabInsteadOfSpaces`, `LineTooLong`, `SensitiveData`, `KeyOrdering`, `KeyOrderConfigurable`, `MaxKeyNameLength`, `ForbidDefaultValue`, `DuplicateListElement`, `KeyPatternMismatch`, `ValuePatternMismatch`, `DocumentStart`, `TrailingSpaces`, `TrailingDots`, `NewlineAtEof`, `ConsecutiveEmptyLines`, `EmptyLineBetweenBlocks`, `DocumentEnd`, `CommentIndentation`, `QuotedKeys`, `QuotedValues`, `IndentSpaces`, `ForbidUnicode`, `ForbidBOM`, `DockerComposeServiceImage`, `JsonSchema`.

Список правил в коде: `cmd/validator.go` (builtinRules), вывод: `yaml-validator rules list`.
