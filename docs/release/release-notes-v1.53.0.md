# Изменения в v1.53.0

## v1.53.0 — Техдолг, ROADMAP, документация

### Исправлено

- Форматирование в `internal/validator/style.go`: лишний отступ у блока `RequireCommentsIndented`.

### Изменено

- **Рефакторинг** — устранено дублирование в `cmd/validator.go`: общая логика валидации вынесена в функцию `validateAndReportFiles`.
- **Плагины** — функции `buildFlatMap` и `findChild` перенесены в `internal/validator/plugins/utils.go`.
- **CI** — добавлены шаги `go vet`, `golangci-lint` и расчёт покрытия тестов.
- **Makefile** — цель `make lint` (go vet + golangci-lint).

### Добавлено

- **golangci-lint** — конфигурация `.golangci.yml` (errcheck, govet, ineffassign, staticcheck).
- **ROADMAP** — раздел «8. Будущие идеи»: LSP, параллельная валидация, CodeClimate, плагин Terraform, проверка консистентности между файлами, покрытие в CI.
- **Документация** — `docs/FAQ.md` (часто задаваемые вопросы), `docs/EXAMPLES.md` (примеры сценариев), таблица документации в README.
