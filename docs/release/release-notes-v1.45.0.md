# Изменения в v1.45.0

## v1.45.0

- **Предупреждения (severity)** (5.7) — добавлено правило `rule_severity` в конфиге: map Type → «error»|«warning». Поле `Severity` в `pkg.Error`. Exit code: 0 — OK, 1 — есть ошибки, 2 — только предупреждения. Форматы severity и github-annotations выводят [WARN] / ::warning для warnings. Пример: `configs/rule-severity-warnings.yaml`.
