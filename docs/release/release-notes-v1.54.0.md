# Изменения в v1.54.0

## v1.54.0 — LSP (Language Server Protocol)

### Добавлено

- **Подкоманда `lsp`** — режим работы как Language Server для IDE. Запуск: `yaml-validator lsp`. Использует stdio для обмена сообщениями (LSP over JSON-RPC 2.0).
- **Поддержка textDocument/didOpen, textDocument/didChange** — при открытии и изменении YAML-файла выполняется валидация.
- **textDocument/publishDiagnostics** — подсветка ошибок и предупреждений в редакторе в реальном времени.
- **Валидация по содержимому** — новые функции для работы с буфером в памяти (несохранённые файлы): `ParseBytesMulti`, `ValidateFromContent`, `CheckStyleContent`, `CheckCommonErrorsContent`, `CheckSyntaxContent`, `FilterInlineIgnoreContent`.

### Использование

```bash
yaml-validator lsp
```

Подключение в VS Code, Neovim и др. через LSP-клиент с командой `yaml-validator lsp` и transport stdio.

Документация: [docs/ide.md](../ide.md#lsp-language-server-protocol).
