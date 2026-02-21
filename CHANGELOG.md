# Changelog

Все изменения проекта документируются в этом файле.

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
