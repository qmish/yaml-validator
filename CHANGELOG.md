# Changelog

Все изменения проекта документируются в этом файле.

## [1.65.0] - 2026-02-22

### Добавлено

- **Docker-образ GHCR (8.7)** — workflow docker-publish: push тегов v* → ghcr.io/qmish/yaml-validator. docs/DISTRIBUTION.md.

## [1.64.0] - 2026-02-22

### Добавлено

- **Плагин OpenAPI (8.10)** — OpenAPIValidator: openapi/swagger, info.title, info.version, paths. plugins: 79.1%.

## [1.63.0] - 2026-02-22

### Добавлено

- **Формат SonarQube Generic (8.9)** — вывод `-o sonarqube`: JSON SonarQube Generic Issue Import (rules, issues). reporter: 77.5%.

## [1.62.0] - 2026-02-22

### Добавлено

- **Документация** — README (CodeClimate, LSP, Terraform, консистентность, плагины, Jenkins); docs/Jenkinsfile-example (Go 1.24, find+xargs, --jobs 4); docs/EXAMPLES.md (Terraform, консистентность).

## [1.61.0] - 2026-02-22

### Добавлено

- **ROADMAP** — новые идеи в разделе 8: 8.7 Docker-образ, 8.8 VS Code расширение, 8.9 SonarQube Generic, 8.10 плагин OpenAPI.

## [1.60.0] - 2026-02-22

### Добавлено

- **Makefile** — цель `make coverage-check`: тесты с покрытием и проверка порога 40% (как в CI).

### Изменено

- **go mod tidy** — обновление и очистка зависимостей.

## [1.59.1] - 2026-02-22

### Исправлено

- **CI (golangci-lint)** — errcheck: `_ = conn.Notify`, `_ = cmd.Help`; gosimple: убраны лишние nil-проверки перед `len()`; ineffassign: `var newPath/keyPath` вместо неэффективного присваивания. Убраны устаревшие линтеры deadcode, varcheck из `.golangci.yml`.

## [1.59.0] - 2026-02-20

### Добавлено

- **Покрытие в CI (8.6)** — шаг проверки минимального порога покрытия (40%) и загрузка `coverage.out` в артефакты workflow.

## [1.58.0] - 2026-02-20

### Добавлено

- **Проверка консистентности между файлами (8.5)** — опция `consistency.enabled` и `consistency.paths` в конфиге. При валидации нескольких файлов сравниваются значения по dot-путям. Ошибка `InconsistentValue` при расхождении. `parser.GetValueAtPath`, `configs/consistency.yaml`.

## [1.57.0] - 2026-02-20

### Добавлено

- **Плагин Terraform (8.4)** — TerraformValidator: проверка tfvars.yaml и YAML в контексте Terraform. Правило `TerraformVarNameHyphen`: имена переменных Terraform используют underscore, не hyphen. Конфиг `configs/terraform.yaml`.

## [1.56.0] - 2026-02-20

### Добавлено

- **Формат CodeClimate (8.3)** — вывод `-o codeclimate`: NDJSON в формате Code Climate Engine Specification. Каждая строка — JSON-объект issue (type, check_name, description, categories, location, severity, fingerprint). Интеграция с Code Climate, codeclimate analyze.

## [1.55.0] - 2026-02-20

### Добавлено

- **Параллельная валидация (8.2)** — флаг `--jobs N`: валидация нескольких файлов параллельно. Ускорение на больших проектах с множеством YAML-файлов. По умолчанию `--jobs 1` (последовательная валидация).

## [1.54.0] - 2026-02-20

### Добавлено

- **LSP (8.1)** — подкоманда `yaml-validator lsp`: Language Server Protocol для IDE. Работа в режиме stdio: `textDocument/didOpen`, `textDocument/didChange`, `textDocument/publishDiagnostics`. Подсветка ошибок в реальном времени. Зависимости: go.lsp.dev/protocol, go.lsp.dev/jsonrpc2.

- **Валидация по содержимому** — `parser.ParseBytesMulti`, `validator.ValidateFromContent`, `CheckStyleContent`, `CheckCommonErrorsContent`, `CheckSyntaxContent`, `FilterInlineIgnoreContent` для LSP (несохранённые буферы).

## [1.53.0] - 2026-02-20

### Исправлено

- Форматирование в `internal/validator/style.go`: лишний отступ у блока `RequireCommentsIndented`.

### Изменено

- **Рефакторинг** — устранено дублирование в `cmd/validator.go`: общая логика валидации вынесена в `validateAndReportFiles`.
- **Плагины** — `buildFlatMap` и `findChild` перенесены в `internal/validator/plugins/utils.go`.
- **CI** — добавлены `go vet`, `golangci-lint` и расчёт покрытия тестов (`-coverprofile=coverage.out`).
- **Makefile** — цель `make lint` (go vet + golangci-lint).

### Добавлено

- **golangci-lint** — конфиг `.golangci.yml` (errcheck, govet, ineffassign, staticcheck).
- **ROADMAP** — раздел «8. Будущие идеи»: LSP, параллельная валидация, CodeClimate, плагин Terraform, проверка консистентности, покрытие в CI.
- **Документация** — `docs/FAQ.md`, `docs/EXAMPLES.md`, таблица документации в README.

## [1.52.0] - 2026-02-21

### Добавлено

- **Мультидокумент YAML** (6.7) — `parser.ParseFileMulti` парсит все документы. Каждый документ валидируется отдельно, ошибки содержат `DocumentIndex`, вывод — `file#docN`. JsonSchema/K8sSchema — только первый документ.

## [1.51.0] - 2026-02-21

### Добавлено

- **Регрессионные тесты** (7.3) — `testdata/regression/` с эталонными YAML и golden-файлами. Пакет `internal/regression` сравнивает compact-вывод с эталоном. CI-шаг в `.github/workflows/yaml-validation.yml`.

## [1.50.0] - 2026-02-21

### Добавлено

- **Пакеты для дистрибутивов** (7.2) — документ `docs/DISTRIBUTION.md`: инструкции для .deb, .rpm, Homebrew, Chocolatey. Скрипты `scripts/packaging/build-deb.sh`, `build-rpm.sh`, формула Homebrew и nuspec для Chocolatey.

## [1.49.0] - 2026-02-21

### Добавлено

- **Watch-режим** (7.1) — флаг `--watch`: при изменении файлов автоматически перезапускать валидацию. Debounce 300 мс. Выход по Ctrl+C.

## [1.48.0] - 2026-02-21

### Добавлено

- **Схема по URL** (6.6) — `schema_path` поддерживает URL (`https://...`), опция `schema_cache_dir` кэширует загруженные схемы локально.

## [1.47.0] - 2026-02-21

### Добавлено

- **Плагин Helm/Kustomize** (6.5) — базовая валидация Helm Chart.yaml (apiVersion, name, version) и Kustomize kustomization.yaml (apiVersion с kustomize, resources или bases). Не реагирует на K8s-манифесты.

## [1.46.0] - 2026-02-21

### Добавлено

- **Плагин Ansible** (6.4) — проверка структуры Ansible playbook: hosts или name в каждом play, tasks (module/include/role), roles (role spec должен иметь ключ `role`). Пропускает K8s и docker-compose файлы.

## [1.45.0] - 2026-02-21

### Добавлено

- **Предупреждения (severity)** (5.7) — правило `rule_severity` в конфиге: map Type → «error»|«warning». Поле `Severity` в `pkg.Error`. Exit code: 0 — OK, 1 — есть ошибки, 2 — только предупреждения. Форматы severity и github-annotations выводят [WARN] / ::warning для warnings. Пример: `configs/rule-severity-warnings.yaml`.

## [1.44.0] - 2026-02-21

### Добавлено

- **Регулярки для ключей/значений** (5.6) — правило `key_value_patterns` с `path`, `pattern` и `target` (keys/values): проверка соответствия ключей или значений регулярному выражению. Например, `metadata.name` по K8s DNS-шаблону. Пример: `configs/k8s-name-patterns.yaml`.

## [1.43.0] - 2026-02-21

### Добавлено

- **Уникальность элементов списка** (5.5) — правило `unique_list_fields` с `path` и `field`: проверка, что элементы массива (по выбранному полю) не повторяются. Например, `spec.template.spec.containers[].name` — имена контейнеров должны быть уникальны. Пример: `configs/k8s-unique-containers.yaml`.

## [1.42.0] - 2026-02-21

### Добавлено

- **Запрет значений по умолчанию** (5.4) — правило `forbid_default_values` и `default_values`: запрет ключей со значением по умолчанию (для K8s: imagePullPolicy: Always, restartPolicy: Always, terminationGracePeriodSeconds: 30). Пример: `configs/k8s-forbid-defaults.yaml`.

## [1.41.0] - 2026-02-21

### Добавлено

- **Миграция с yamllint** (4.7) — документ `docs/YAMLLINT_MIGRATION.md`: таблица соответствия правил yamllint ↔ yaml-validator и пример конвертации `.yamllint` в `yaml-validator.yaml`, а также сравнение команд и pre-commit хуков.

## [1.40.0] - 2026-02-21

### Добавлено

- **Документация по всем правилам** (4.6) — документ `docs/RULES.md`: список правил, ключи конфига, опции и примеры (основные правила, стиль, схемы, file_profiles, типы ошибок).

## [1.39.0] - 2026-02-21

### Добавлено

- **Конфиг по имени файла** (4.5) — `file_profiles` в yaml-validator.yaml: автовыбор профиля по маске (напр. `**/k8s/**` → k8s-strict, `*docker-compose*.yaml` → docker-compose). Пример: `configs/file-profiles.yaml`.

## [1.38.0] - 2026-02-21

### Добавлено

- **Автофикс (--fix)** (4.4) — флаг `--fix` для автоисправления: trailing spaces, newline at EOF, consecutive empty lines. Пакет `internal/fixer`.

## [1.37.0] - 2026-02-21

### Добавлено

- **VS Code / IDE** (3.8) — документация `docs/ide.md` по использованию как линтера: VS Code (tasks с problemMatcher, Run on Save), JetBrains (External Tools), Sublime, Vim.

## [1.36.0] - 2026-02-21

### Добавлено

- **pre-commit (official)** (3.7) — hook подключается из репозитория: `repo: https://github.com/qmish/yaml-validator`. Добавлен `.pre-commit-hooks.yaml`, `docs/pre-commit.md`.

## [1.35.0] - 2026-02-21

### Добавлено

- **GitLab CI** (3.5) — готовый фрагмент job для валидации YAML в `docs/gitlab-ci-example.yml`. Варианты: простая валидация (severity) и с отчётом Code Quality (gl-code-quality-report.json).
- **Jenkinsfile** (3.6) — пример пайплайна Jenkins с шагом валидации YAML в `docs/Jenkinsfile-example`.

## [1.34.0] - 2026-02-21

### Добавлено

- **Список правил (машиночитаемый)** (2.7) — подкоманда `rules list` с выводом в JSON или YAML (`-o json`, `-o yaml`) для скриптов и документации. Список правил с id, description, config_key.

## [1.33.0] - 2026-02-21

### Добавлено

- **Режим --quiet** (2.6) — флаг `-q`/`--quiet`: минимальный вывод — только итог (OK или N errors). Для человекочитаемых форматов (human, compact, severity, github-annotations).

## [1.32.0] - 2026-02-21

### Добавлено

- **Разные exit code** (2.5) — явные коды выхода: 0 — OK, 1 — есть ошибки, 2 — только предупреждения (зарезервировано для 5.7). Описание в справке `validate`.

## [1.31.0] - 2026-02-21

### Добавлено

- **Checkstyle XML** (2.4) — формат вывода `-o checkstyle` для Jenkins, SonarQube. Совместим с форматом Checkstyle XML.

## [1.30.0] - 2026-02-21

### Добавлено

- **Проверка BOM** (1.9) — правило `style.forbid_bom`: ошибка при наличии UTF-8 BOM в начале файла.

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
