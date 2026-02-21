# pre-commit (official)

yaml-validator поддерживается как official pre-commit hook — можно подключать напрямую из репозитория.

## Подключение

Добавьте в `.pre-commit-config.yaml`:

```yaml
repos:
  - repo: https://github.com/qmish/yaml-validator
    rev: v1.35.0  # укажите нужную версию (тег)
    hooks:
      - id: yaml-validator
```

Установка:

```bash
pre-commit install
pre-commit run --all-files  # проверить все файлы
```

## Требования

- [pre-commit](https://pre-commit.com/) — `pip install pre-commit`
- Go в PATH — pre-commit собирает yaml-validator при первом запуске

## Вариант без Go

Если Go не установлен, используйте бинарник:

```yaml
repos:
  - repo: local
    hooks:
      - id: yaml-validator
        name: yaml-validator
        entry: yaml-validator validate
        language: system
        types: [yaml]
```

Скачайте бинарник из [releases](https://github.com/qmish/yaml-validator/releases) и добавьте в PATH.
