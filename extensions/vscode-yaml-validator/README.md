# YAML Validator (VS Code Extension)

Расширение VS Code для валидации YAML через [yaml-validator](https://github.com/qmish/yaml-validator) LSP.

## Требования

- [yaml-validator](https://github.com/qmish/yaml-validator) в PATH (`go install github.com/qmish/yaml-validator@latest` или скачать релиз)

## Установка (разработка)

```bash
cd extensions/vscode-yaml-validator
npm install
npm run compile
```

В VS Code: Run Extension (F5).

## Конфигурация

- `yamlValidator.serverPath` — путь к yaml-validator (по умолчанию: `yaml-validator`)
