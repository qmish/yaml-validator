# Изменения в v1.58.0

## v1.58.0 — Проверка консистентности между файлами (8.5)

### Добавлено

- **Секция `consistency` в конфиге** — проверка совпадения значений по заданным путям между несколькими YAML-файлами.
- **consistency.enabled** — включить проверку.
- **consistency.paths** — список dot-путей для сравнения (напр. `version`, `metadata.annotations.app-version`).
- **Ошибка InconsistentValue** — при расхождении значений между файлами.

### Использование

```yaml
# configs/consistency.yaml
consistency:
  enabled: true
  paths:
    - version
    - apiVersion
```

```bash
yaml-validator validate file1.yaml file2.yaml file3.yaml -c configs/consistency.yaml
```
