# Сравнение с существующими решениями

Анализ популярных инструментов валидации YAML и сравнение с **yaml-validator** (данный проект).

---

## 1. Обзор решений

| Инструмент      | Язык   | Фокус              | Тип лицензии |
|-----------------|--------|--------------------|--------------|
| **yamllint**    | Python | Синтаксис + стиль  | GPL v3       |
| **kubeval**     | Go     | K8s JSON Schema    | Apache 2.0   |
| **kubeconform** | Go     | K8s OpenAPI        | Apache 2.0   |
| **kube-score**  | Go     | K8s best practices | —            |
| **yq**          | Go     | Запросы/формат     | MIT          |
| **yaml-validator** (наш) | Go | Синтаксис, дубликаты, K8s-поля, K8s OpenAPI (kubeconform), правила стиля, inline-игнор, отчёты JSON/JUnit/SARIF, плагины, Docker, CI (Go 1.24) | MIT |

---

## 2. yamllint

**Источник:** [yamllint](https://yamllint.readthedocs.io/), [GitHub](https://github.com/adrienverge/yamllint)

### Возможности

- Синтаксис YAML
- 20+ правил: отступы, дубликаты ключей (`key-duplicates`), длина строки (`line-length`), trailing spaces, кавычки, комментарии, document-start/end и др.
- Конфиг: `.yamllint`, `.yamllint.yaml`
- Формат вывода в стиле ESLint
- Отключение правил через комментарии в файле
- Интеграция: pre-commit, Ansible, редакторы

### Сравнение с yaml-validator

| Критерий              | yamllint           | yaml-validator                    |
|-----------------------|--------------------|-----------------------------------|
| Дубликаты ключей      | ✅ key-duplicates  | ✅ CheckDuplicates                |
| Длина строки          | ✅ line-length     | ✅ max_line_length                 |
| Табы vs пробелы       | ✅ indentation     | ✅ TabInsteadOfSpaces             |
| document-start / EOF  | ✅                | ✅ style (v1.3.0)                  |
| document-end (...)    | ✅                | ✅ style.require_document_end (v1.8.0) |
| Пустые строки подряд  | ✅ empty-lines    | ✅ style.forbid_consecutive_empty_lines (v1.7.0) |
| Trailing spaces       | ✅                | ✅ style.forbid_trailing_spaces    |
| Отступ комментариев   | ✅ comment-indentation | ✅ style.require_comments_indented |
| Порядок ключей        | ✅ key-ordering    | ✅ check_key_ordering, key_order (v1.19.0) |
| Ключи в кавычках      | ✅ quoted-strings  | ✅ style.require_quoted_keys        |
| Строковые значения в кавычках | ✅ quoted-values | ✅ style.require_quoted_values (v1.18.0) |
| Inline disable        | ✅ disable-line    | ✅ disable-line / disable-next-line (v1.4.0) |
| Синтаксис             | ✅                 | ✅ CheckSyntax                    |
| Обязательные поля     | ❌                 | ✅ apiVersion, kind, metadata.name |
| Чувствительные данные | ❌                 | ✅ password, token, secret        |
| Плагины / расширения  | ❌                 | ✅ ValidatorPlugin, K8s-плагин    |
| Вывод JSON/JUnit      | Ограниченно        | ✅ -o json, -o junit               |
| Язык                  | Python             | Go (один бинарник, без runtime)   |
| Лицензия              | GPL v3             | MIT                               |

**Итог:** yaml-validator покрывает базовые проверки стиля и синтаксиса плюс целостность K8s и расширяемость; yamllint богаче по правилам стиля (комментарии, document markers, quoted-strings и т.д.) и экосистеме (pre-commit, Ansible).

---

## 3. kubeval / kubeconform

**Источники:** [kubeval](https://github.com/instrumenta/kubeval), [kubeconform](https://github.com/yannh/kubeconform)

### Возможности

- Валидация манифестов против **официальных схем Kubernetes** (OpenAPI/JSON Schema)
- Проверка типов: `replicas` — число, а не строка
- Версии API (apiVersion) и совместимость
- Работа в CI без кластера

### Сравнение с yaml-validator

| Критерий           | kubeval / kubeconform | yaml-validator              |
|--------------------|------------------------|-----------------------------|
| Схема K8s API      | ✅ полная              | ✅ опционально (kubeconform, v1.2.0) |
| apiVersion, kind   | ✅ по схеме            | ✅ плагин K8s + схема       |
| Типы полей        | ✅ (int, string, …)   | ✅ при k8s_schema.enabled  |
| Синтаксис YAML     | Зависит от парсера    | ✅                          |
| Дубликаты ключей   | ❌                     | ✅                          |
| Другие форматы    | Только K8s             | ✅ любой YAML + конфиг      |

**Итог:** для строгой проверки K8s против API лучше kubeval/kubeconform; yaml-validator даёт быструю проверку структуры (обязательные поля, дубликаты) и подходит для любого YAML (docker-compose, конфиги).

---

## 4. kube-score

**Источник:** [kube-score](https://github.com/zegl/kube-score)

### Возможности

- Анализ K8s-манифестов на **best practices** и безопасность
- Проверки: resource requests/limits, readiness/liveness, PodSecurityPolicy, сеть и т.д.
- Не заменяет валидатор синтаксиса/схемы

### Сравнение

kube-score решает другую задачу (качество и безопасность), а не парсинг и базовую валидацию. С yaml-validator не конкурирует, может использоваться вместе.

---

## 5. yq

**Источник:** [yq](https://github.com/mikefarah/yq)

### Возможности

- Запросы и преобразование YAML (аналог jq для YAML)
- Форматирование, слияние файлов
- Валидация — побочная (ошибки парсинга при чтении)

### Сравнение

yq — инструмент обработки YAML, а не линтер. yaml-validator — именно валидатор с отчётами и правилами.

---

## 6. Сводная таблица

| Функция                    | yamllint | kubeval/kubeconform | yaml-validator |
|---------------------------|----------|----------------------|----------------|
| Синтаксис YAML            | ✅       | ✅                   | ✅             |
| Дубликаты ключей          | ✅       | ❌                   | ✅             |
| Стиль (document-start/end, trailing, newline EOF, пустые строки, длина, табы) | ✅ много | ❌ | ✅ (style + common_errors) |
| Inline disable правил     | ✅       | ❌                   | ✅ (v1.4.0)    |
| Обязательные поля K8s     | ❌       | через схему          | ✅             |
| Схема K8s (типы, API)     | ❌       | ✅                   | ✅ опц. (v1.2.0, kubeconform) |
| Чувствительные данные     | ❌       | ❌                   | ✅             |
| Плагины / свои проверки   | ❌       | ❌                   | ✅             |
| JSON/JUnit/SARIF отчёт    | огранич. | разное               | ✅ (v1.5.0 SARIF) |
| Один бинарник (Go)        | ❌       | ✅                   | ✅             |
| Любой YAML (не только K8s)| ✅       | ❌                   | ✅             |

---

## 7. Выводы

### Сильные стороны yaml-validator

1. **Один бинарник на Go** — не нужен Python или внешний runtime, удобно в CI и контейнерах.
2. **Расширяемость** — плагины (`ValidatorPlugin`), встроенный K8s-плагин.
3. **Гибкий конфиг** — разные профили (K8s vs docker-compose, strict), свои обязательные поля и паттерны.
4. **Отчёты для CI** — JSON, JUnit и **SARIF** (v1.5.0) для GitHub Code Scanning.
5. **K8s по схеме** (v1.2.0) — опциональная полная проверка через kubeconform (`configs/k8s-strict.yaml`).
6. **Правила стиля** (v1.3.0–v1.8.0) — document-start, document-end, trailing spaces, newline at EOF, запрет пустых строк подряд (`configs/strict.yaml`).
7. **Inline-игнор** (v1.4.0) — отключение правил через комментарии, как в yamllint.
8. **Лицензия MIT** — проще использование в корпоративных проектах по сравнению с GPL.
9. **Сборка и CI** (v1.5.1) — Docker на Go 1.24, GitHub Actions с явной версией Go, публикация образа в ghcr.io при релизе.

### Что можно усилить (по сравнению с другими)

1. **Правила стиля** — реализованы: document-start/end, trailing spaces, trailing dots, newline EOF, пустые строки подряд, отступ комментариев, порядок ключей, кавычки для ключей и значений (`require_quoted_keys`, `require_quoted_values`). Длина строки через `max_line_length`.
2. **Интеграции** — готовый pre-commit (`.pre-commit-config.yaml`), примеры в README для GitLab CI и Jenkins (Docker и бинарник).
3. **Доп. форматы** — реализовано: `-o compact`, SARIF, **GitLab Code Quality** (`-o gitlab`), **GitHub Annotations** (`-o github-annotations`), **severity** (`-o severity`: `[ERROR] file:line: message`). Колонка в compact/SARIF, когда доступна.

Документ обновлён по состоянию v1.19.0.
