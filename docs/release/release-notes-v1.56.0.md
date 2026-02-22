# Изменения в v1.56.0

## v1.56.0 — Формат CodeClimate (8.3)

### Добавлено

- **Флаг `-o codeclimate`** — вывод в формате Code Climate Engine (NDJSON). Каждая строка — отдельный JSON-объект issue.
- Совместимость со спецификацией [Code Climate Engine SPEC](https://github.com/codeclimate/platform/blob/master/spec/analyzers/SPEC.md).
- Поля: type, check_name, description, categories, location (path, lines), severity, fingerprint, remediation_points.

### Использование

```bash
yaml-validator validate **/*.yaml -o codeclimate
```

Для интеграции с Code Climate или `codeclimate analyze` направьте вывод в файл или pipe.
