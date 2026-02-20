# Changelog

Все изменения проекта документируются в этом файле.

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
