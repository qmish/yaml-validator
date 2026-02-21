# Standalone Git hooks

Подключайте без pre-commit framework.

## pre-commit

Проверяет YAML перед коммитом.

```bash
cp hooks/pre-commit .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
```

Требуется `yaml-validator` в PATH.
