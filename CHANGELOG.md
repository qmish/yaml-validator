# Changelog

Все изменения проекта документируются в этом файле.

## [1.2.0] - 2025-02-20

### Добавлено

- **Проверка по полной OpenAPI/схеме Kubernetes** (опционально): типы полей, все стандартные ресурсы K8s. Включается через `rules.k8s_schema.enabled: true` в конфиге. Используется [kubeconform](https://github.com/yannh/kubeconform).
- Конфиг `configs/k8s-strict.yaml` для валидации K8s по схеме.
- Параметры `k8s_schema`: `version`, `strict`, `cache_dir`, `ignore_missing_schemas` (CRD).

## [1.0.1] - 2025-02-20

### Исправлено

- go.mod: версия Go исправлена с 1.25.6 на 1.21 (Docker-сборка)
- CI: валидация — не проверять configs как K8s-манифесты
- Makefile, docker-compose: корректные пути для validate
- release.ps1: добавлен -short к тестам для стабильности CI

## [1.0.0] - 2025-02-20

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
