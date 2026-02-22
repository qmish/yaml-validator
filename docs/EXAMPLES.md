# Примеры использования

Практические сценарии использования yaml-validator.

---

## Быстрый старт

```bash
# Одна команда
yaml-validator validate config.yaml

# Несколько файлов
yaml-validator validate *.yaml

# Тихий режим (только OK или количество ошибок)
yaml-validator validate config.yaml -q
```

---

## Сценарий: Kubernetes-манифесты

### Базовая проверка (apiVersion, kind, metadata)

```bash
yaml-validator validate deployment.yaml -c configs/default.yaml
```

### Строгий порядок ключей (metadata перед spec)

```bash
yaml-validator validate deployment.yaml -c configs/k8s-key-order.yaml
```

### Полная валидация по OpenAPI схеме K8s

```bash
yaml-validator validate deployment.yaml -c configs/k8s-strict.yaml
```

### Папка с манифестами + автовыбор конфига

В `yaml-validator.yaml`:

```yaml
file_profiles:
  - pattern: "**/manifests/**"
    config: "configs/k8s-strict.yaml"
```

```bash
yaml-validator validate manifests/*.yaml
```

---

## Сценарий: docker-compose

```bash
yaml-validator validate docker-compose.yaml -c configs/docker-compose.yaml
```

Проверяет: наличие `image` или `build` у каждого сервиса.

---

## Сценарий: CI/CD

### GitHub Actions (простой)

```yaml
- run: go install github.com/qmish/yaml-validator@latest
- run: yaml-validator validate "**/*.yml" "**/*.yaml"
```

### GitHub Code Scanning (SARIF)

```yaml
- run: yaml-validator validate "**/*.yaml" -o sarif > results.sarif
- uses: github/codeql-action/upload-sarif@v3
  with:
    sarif_file: results.sarif
```

### GitLab CI (Code Quality)

```yaml
validate-yaml:
  script:
    - yaml-validator validate "**/*.yml" "**/*.yaml" -o gitlab > gl-code-quality-report.json
  artifacts:
    reports:
      codequality: gl-code-quality-report.json
```

### Jenkins (Checkstyle для SonarQube)

```bash
yaml-validator validate "**/*.yaml" -o checkstyle > checkstyle-result.xml
```

---

## Сценарий: Автофикс перед коммитом

```bash
# Исправить trailing spaces, newline at EOF, лишние пустые строки
yaml-validator validate . --fix
```

В pre-commit: добавить `--fix` в аргументы хука (если хук поддерживает).

---

## Сценарий: Watch-режим при разработке

```bash
yaml-validator validate config.yaml --watch
```

При сохранении файла валидация перезапускается автоматически.

---

## Сценарий: Миграция с yamllint

1. Создайте конфиг по таблице соответствия: [docs/YAMLLINT_MIGRATION.md](YAMLLINT_MIGRATION.md)
2. Запустите:

```bash
yaml-validator validate **/*.yaml -c yaml-validator.yaml
```

---

## Сценарий: Список правил для скриптов

```bash
yaml-validator rules list -o json
yaml-validator rules list -o yaml
```

---

## Сценарий: Ansible playbooks

Конфиг с включённым плагином Ansible (проверка структуры playbook):

```bash
yaml-validator validate playbook.yml -c configs/ansible.yaml
```

(Если есть `configs/ansible.yaml` — см. docs по плагинам.)
