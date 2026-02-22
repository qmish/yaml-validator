# VS Code и IDE: использование как линтера

Использование yaml-validator как линтера в редакторах. Формат `-o compact` выводит `file:line[:col]: message` (ESLint-style), подходящий для парсинга IDE. Для live-валидации при редактировании можно использовать `yaml-validator validate file.yaml -o compact --watch` (повторная проверка при сохранении).

## LSP (Language Server Protocol)

Начиная с v1.54.0, yaml-validator поддерживает LSP. Команда `yaml-validator lsp` запускает сервер в режиме stdio для интеграции с IDE (подсветка ошибок в реальном времени).

### VS Code

**Вариант A: Расширение YAML Validator** (рекомендуется)

В репозитории есть готовое расширение `extensions/vscode-yaml-validator`. Установка из исходников:

```bash
cd extensions/vscode-yaml-validator
npm install
npm run compile
```

Затем в VS Code: F5 для запуска в режиме Extension Development Host. Либо собрать .vsix: `vsce package` и установить через Extensions → Install from VSIX.

**Вариант B: Generic LSP-клиент**

Добавьте в настройки или `settings.json`:

```json
{
  "yaml-validator-lsp.serverPath": "yaml-validator",
  "yaml-validator-lsp.serverArgs": ["lsp"]
}
```

Или через generic LSP-клиент: команда `yaml-validator lsp`, transport stdio.

## VS Code

### Вариант 1: Task с Problem Matcher

Создайте `.vscode/tasks.json` в корне проекта:

```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "Validate YAML",
      "type": "shell",
      "command": "yaml-validator",
      "args": ["validate", "${file}", "-o", "compact"],
      "group": {
        "kind": "build",
        "isDefault": true
      },
      "presentation": {
        "reveal": "silent"
      },
      "problemMatcher": {
        "owner": "yaml-validator",
        "fileLocation": ["relative", "${workspaceFolder}"],
        "pattern": {
          "regexp": "^(.*):(\\d+):(\\d+): (.*)$",
          "file": 1,
          "line": 2,
          "column": 3,
          "message": 4
        }
      }
    }
  ]
}
```

Для формата без колонки (`file:line: message`):

```json
"pattern": {
  "regexp": "^(.*):(\\d+): (.*)$",
  "file": 1,
  "line": 2,
  "message": 3
}
```

Запуск: `Ctrl+Shift+B` (Windows/Linux) или `Cmd+Shift+B` (macOS), либо Terminal → Run Build Task.

### Вариант 2: Run on Save

Расширение [Run on Save](https://marketplace.visualstudio.com/items?itemName=emeraldwalk.RunOnSave) — запускать валидацию при сохранении. В `settings.json`:

```json
{
  "emeraldwalk.runonsave": {
    "commands": [
      {
        "match": "\\.ya?ml$",
        "cmd": "yaml-validator validate ${file} -o compact"
      }
    ]
  }
}
```

## JetBrains (IntelliJ IDEA, GoLand, PyCharm)

1. **File → Settings → Tools → External Tools**
2. Добавить инструмент:
   - **Name:** YAML Validator
   - **Program:** `yaml-validator` (или полный путь к бинарнику)
   - **Arguments:** `validate $FilePath$ -o compact`
   - **Working directory:** `$ProjectFileDir$`
3. Использовать: правый клик по YAML-файлу → External Tools → YAML Validator

Output появится в Run. Для отображения ошибок в Editor можно настроить «Output Filters» с regex `(.*):(\d+):(\d+): (.*)`.

## Sublime Text

Через [External Command](https://packagecontrol.io/packages/External%20Command) или shell command: `yaml-validator validate $file -o compact`. Парсинг вывода — через соответствующий плагин.

## Vim / Neovim

`:!yaml-validator validate % -o compact` — запуск вручную. Для интегрированного линтинга — ALE, null-ls или аналоги с поддержкой внешних linter’ов.
