# Изменения в v1.63.0

## v1.63.0 — Формат SonarQube Generic (8.9)

### Добавлено

- **Флаг `-o sonarqube`** — вывод в формате SonarQube Generic Issue Import (JSON).
- Совместимость с `sonar.externalIssuesReportPaths`.
- Поля: `rules` (id, name, description, engineId, type, severity), `issues` (engineId, ruleId, primaryLocation с message, filePath, textRange).

### Использование

```bash
yaml-validator validate **/*.yaml -o sonarqube > sonarqube-report.json
```

В `sonar-project.properties`:
```properties
sonar.externalIssuesReportPaths=sonarqube-report.json
```

### Покрытие тестами

reporter: 77.5%
