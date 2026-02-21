# Changelog

Все изменения проекта документируются в этом файле.

## [1.29.0] - 2026-02-21

### Добавлено

- **Запрет не-ASCII** (1.8) — правило `style.forbid_unicode`: запрет не-ASCII символов в ключах и строковых значениях (строгие ASCII-конфиги).

## [1.28.0] - 2026-02-21

### Добавлено

- **Колонка в ошибках длины строки** (1.7) — при превышении `max_line_length` в `pkg.Error` заполняется поле `Column` (позиция, превысившая лимит). Выводится в форматах compact (`file:line:col: message`) и severity; SARIF использует `startColumn`.

## [1.27.0] - 2026-02-21

### Добавлено

- **Минимум пустых строк между блоками** (1.6) — параметр `style.min_empty_lines_between_blocks` (0 или 1): при 1 — требовать хотя бы одну пустую строку между топ-уровневыми блоками; при 0 — проверка отключена. Гибче, чем только `require_empty_line_between_blocks`.

## [1.26.0] - 2026-02-21

### Добавлено

- **Явное правило для табов** (1.5) — `style.forbid_tabs`: отдельное правило стиля для запрета табуляции. При включении проверка выполняется в CheckStyle; иначе — в check_common_errors (обратная совместимость). Конфиг `configs/strict.yaml` включает `forbid_tabs: true`.

## [1.25.0] - 2026-02-21

### Добавлено

- **Поддержка JSON Schema** (6.3) — опциональная валидация по произвольной JSON Schema (не только K8s): `rules.json_schema.enabled: true`, `rules.json_schema.schema_path: путь/к/схеме.json`. Конфиг `configs/json-schema.yaml`, пример схемы `testdata/schemas/config-schema.json`. Используется [santhosh-tekuri/jsonschema](https://github.com/santhosh-tekuri/jsonschema).

## [1.24.0] - 2026-02-21

### Добавлено

- **Одна пустая строка между блоками** (1.3) — правило `style.require_empty_line_between_blocks`: требовать ровно одну пустую строку между топ-уровневыми блоками; конфиг `configs/strict.yaml` включает это правило.

## [1.23.0] - 2026-02-21

### Добавлено

- **Benchmark / performance** (4.3) — тесты производительности: `internal/validator/benchmark_test.go` (BenchmarkValidate_Small/Medium/Large), `internal/parser/benchmark_test.go` (BenchmarkParseFile_Small/Large); цель `make benchmark`.

## [1.22.0] - 2026-02-21

### Добавлено

- **indent_spaces** (1.4) — настраиваемый шаг отступов: `style.indent_spaces: 2` или `4`; проверка, что отступ кратен заданному числу пробелов.

## [1.21.0] - 2026-02-21

### Добавлено

- **Подсказка по исправлению** (4.1) — поле `Suggestion` в `pkg.Error`, вывод «To fix: …» в человекочитаемом формате. Подсказки для DocumentStart, TrailingSpaces, QuotedKeys, QuotedValues, LineTooLong, KeyOrdering и др.

## [1.20.0] - 2026-02-21

### Добавлено

- **Плагин docker-compose** — проверка: у каждого сервиса должно быть `image` или `build`. Работает только для файлов с секцией `services` (без `apiVersion`).

## [1.19.0] - 2026-02-21

### Добавлено

- **key_order** — приоритетный порядок ключей (напр. apiVersion, kind, metadata, spec) вместо алфавитного; конфиг `configs/k8s-key-order.yaml`.
- **max_key_name_length** — ограничение длины имён ключей.

## [1.18.0] - 2026-02-21

### Добавлено

- **quoted-values** — правило `style.require_quoted_values`: строковые значения должны быть в кавычках (как yamllint quoted-values).
- **trailing-dots** — правило `style.forbid_trailing_dots`: запрет точек в конце строк.
- **docs/PLUGINS.md** — документация по созданию плагинов (API, примеры).

## [1.17.0] - 2026-02-21

### Добавлено

- **Формат severity** — вывод `-o severity`: `[ERROR] file:line: message` для скриптов и CI.
- **Azure DevOps** — пример pipeline (`docs/azure-pipelines-example.yml`).
- **Bitbucket Pipelines** — пример pipeline (`docs/bitbucket-pipelines-example.yml`).
- **Standalone pre-commit hook** — `hooks/pre-commit`, подключаемый без pre-commit framework.

## [1.16.0] - 2026-02-21

### Добавлено

- **docs/ROADMAP.md** — чек-лист доработок (правила, форматы, интеграции, документация, плагины).

## [1.15.0] - 2026-02-21

### Добавлено

- **Тесты** — `TestGenerateGitHubAnnotations`, `TestGenerateGitHubAnnotations_Empty`; функция `GenerateGitHubAnnotations` для тестируемого вывода.

### Изменено

- **release.ps1** — команда `gh release create` теперь выводит путь к архиву и ко всем бинарникам (gh не принимает каталог).

## [1.14.0] - 2026-02-20

### Добавлено

- **Формат GitHub Annotations** — вывод `-o github-annotations` (`::error file=...,line=...::`) для GitHub Actions.
- **release.ps1** — артефакты релизов сохраняются в `docs/release/` (release-v*, *.tar.gz).

### Изменено

- Артефакты релизов перенесены в `docs/release/` (см. `.gitignore`).

## [1.13.0] - 2026-02-20

### Добавлено

- **Формат GitLab Code Quality** — вывод `-o gitlab` для GitLab CI (`gl-code-quality-report.json`), артефакт `reports.codequality`.
- **Makefile** — цель `validate` проверяет также `docker-compose.yaml` с профилем `configs/docker-compose.yaml`.
- **Пример workflow** — `docs/code-scanning-example.yml` для GitHub Code Scanning (SARIF + Advanced Security).

## [1.12.0] - 2026-02-20

### Добавлено

- **Правило стиля** — `style.require_quoted_keys`: ключи маппинга должны быть в кавычках в исходном файле (в духе yamllint quoted-strings).

## [1.11.0] - 2026-02-20

### Добавлено

- **Правило порядка ключей** — `check_key_ordering`: требовать лексикографический порядок ключей в маппингах (в духе yamllint key-ordering).

## [1.10.0] - 2026-02-20

### Добавлено

- **Правило стиля** — `style.require_comments_indented`: комментарии внутри отступного блока должны иметь отступ (в духе yamllint comment-indentation).

## [1.9.0] - 2026-02-20

### Добавлено

- **Колонка в отчётах** — в `pkg.Error` добавлено поле `Column`; формат `-o compact` выводит `file:line:col: message`, когда колонка известна (например, для DuplicateKey); SARIF использует `startColumn` в регионах.

## [1.8.0] - 2026-02-20

### Добавлено

- **Правило стиля** — `style.require_document_end`: требовать маркер `...` в конце файла (для много-документного YAML, в духе yamllint document-end).

## [1.7.0] - 2026-02-20

### Добавлено

- **Правило стиля** — `style.forbid_consecutive_empty_lines`: запрет более одной пустой строки подряд (в духе yamllint empty-lines).

## [1.6.0] - 2026-02-20

### Добавлено

- **Формат compact** — вывод `-o compact`: одна строка на ошибку в формате `file:line: message` (ESLint-style), удобно для редакторов и скриптов.

## [1.5.1] - 2026-02-20

### Исправлено

- CI: явная версия Go 1.24 в workflow (совпадение с go.mod и Dockerfile)
- Dockerfile: образ сборки обновлён до `golang:1.24-alpine` для совместимости с go.mod

### Документация

- CONTRIBUTING: уточнена версия Go (1.24+), блок про авторство вынесен в локальный файл по желанию

## [1.5.0] - 2026-02-20

### Добавлено

- **Формат SARIF** — вывод `-o sarif` для GitHub Code Scanning / Advanced Security. Один отчёт по всем переданным файлам в формате SARIF 2.1.0.

## [1.4.0] - 2026-02-20

### Добавлено

- **Inline-игнор**: отключение правил через комментарии в YAML. Формат: `# yaml-validator disable-line [rule:RuleType]` и `# yaml-validator disable-next-line [rule:...]`. Включается через `rules.inline_ignore: true`.

## [1.3.0] - 2026-02-20

### Добавлено

- **Правила стиля** (в духе yamllint): `style.require_document_start`, `style.forbid_trailing_spaces`, `style.require_newline_at_eof`. Включаются через конфиг.
- Конфиг `configs/strict.yaml` для строгой проверки стиля.

## [1.2.0] - 2026-02-20

### Добавлено

- **Проверка по полной OpenAPI/схеме Kubernetes** (опционально): типы полей, все стандартные ресурсы K8s. Включается через `rules.k8s_schema.enabled: true` в конфиге. Используется [kubeconform](https://github.com/yannh/kubeconform).
- Конфиг `configs/k8s-strict.yaml` для валидации K8s по схеме.
- Параметры `k8s_schema`: `version`, `strict`, `cache_dir`, `ignore_missing_schemas` (CRD).

## [1.0.1] - 2026-02-20

### Исправлено

- go.mod: версия Go исправлена с 1.25.6 на 1.21 (Docker-сборка)
- CI: валидация — не проверять configs как K8s-манифесты
- Makefile, docker-compose: корректные пути для validate
- release.ps1: добавлен -short к тестам для стабильности CI

## [1.0.0] - 2026-02-20

### Добавлено

- Валидация синтаксиса YAML
- Поиск дублирующихся ключей
- Проверка обязательных полей (Kubernetes)
- Распространённые ошибки: табуляции, длинные строки, чувствительные данные
- Плагины (интерфейс ValidatorPlugin, плагин Kubernetes)
- Логирование (verbose, JSON для ELK/Loki)
- CLI с Cobra: validate, config init, version
- Форматы отчётов: human, JSON, JUnit
- Docker, docker-compose, GitHub Actions CI
