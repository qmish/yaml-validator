# Изменения в v1.66.0

## v1.66.0 — VS Code расширение (8.8)

Покрытие тестами: **62.7%**

### Добавлено

- **VS Code расширение** — `extensions/vscode-yaml-validator`: LSP-клиент для `yaml-validator lsp`, подсветка ошибок YAML в реальном времени.
- Конфигурация `yamlValidator.serverPath` — путь к бинарнику yaml-validator.
- Документация docs/ide.md: вариант A (расширение) и B (generic LSP).

### Изменено

- **CI** — проверка покрытия: вместо `bc` используется `awk` (портируемость). golangci-lint: v1.64 → v2.10.1. Добавлен job `vscode-extension` (Node.js, npm install, npm run compile). Добавлен `.dockerignore` (исключены extensions/, testdata/, docs/release и др.).

### Установка расширения (разработка)

```bash
cd extensions/vscode-yaml-validator
npm install
npm run compile
```

В VS Code: F5 для запуска в режиме Extension Development Host.
