# Изменения в v1.64.0

## v1.64.0 — Плагин OpenAPI (8.10)

### Добавлено

- **Плагин OpenAPI** — проверка YAML OpenAPI 3 / Swagger 2 спецификаций.
- Обязательные поля: `openapi` или `swagger`, `info` (title, version), `paths`.

### Использование

```bash
yaml-validator validate openapi.yaml -c configs/openapi.yaml
```

Создайте `configs/openapi.yaml` с `check_integrity: false` и `required_fields: []` (аналогично configs/docker-compose.yaml).

### Покрытие тестами

plugins: 79.1%
