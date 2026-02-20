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
| **yaml-validator** (наш) | Go | Синтаксис + дубликаты + K8s-поля + **K8s OpenAPI-схема** (опционально) + расширяемость | MIT |

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
| Схема K8s API      | ✅ полная              | ❌ только базовые поля      |
| apiVersion, kind   | ✅ по схеме            | ✅ плагин K8s                |
| Типы полей        | ✅ (int, string, …)   | ❌ не проверяются           |
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
| Стиль (отступы, длина)    | ✅ много | ❌                   | ✅ базово      |
| Обязательные поля K8s     | ❌       | через схему          | ✅             |
| Схема K8s (типы, API)     | ❌       | ✅                   | ✅ (опционально, kubeconform) |
| Чувствительные данные     | ❌       | ❌                   | ✅             |
| Плагины / свои проверки   | ❌       | ❌                   | ✅             |
| JSON/JUnit отчёт          | огранич. | разное               | ✅             |
| Один бинарник (Go)        | ❌       | ✅                   | ✅             |
| Любой YAML (не только K8s)| ✅       | ❌                   | ✅             |

---

## 7. Выводы

### Сильные стороны yaml-validator

1. **Один бинарник на Go** — не нужен Python или внешний runtime, удобно в CI и контейнерах.
2. **Расширяемость** — плагины (`ValidatorPlugin`), встроенный K8s-плагин.
3. **Гибкий конфиг** — разные профили (K8s vs docker-compose), свои обязательные поля и паттерны.
4. **Отчёты для CI** — JSON и JUnit из коробки.
5. **Доп. проверки** — чувствительные поля, табы, длина строки.
6. **Лицензия MIT** — проще использование в корпоративных проектах по сравнению с GPL.

### Что можно усилить (по сравнению с другими)

1. **Правила стиля** — по образцу yamllint: document-start/end, quoted-strings, комментарии, key-ordering.
2. **K8s-схема** — реализовано в v1.2.0: опциональная проверка по полной OpenAPI/схеме через kubeconform (`rules.k8s_schema.enabled`, конфиг `configs/k8s-strict.yaml`).
3. **Правила стиля** — реализовано в v1.3.0: document-start, trailing spaces, newline at EOF (конфиг `style:`, `configs/strict.yaml`).
4. **Inline-игнор** — реализовано в v1.4.0: `# yaml-validator disable-line rule:XXX`, `disable-next-line` (включается через `rules.inline_ignore`).
4. **Больше форматов вывода** — например, SARIF для GitHub Advanced Security.

При необходимости этот документ можно вынести в репозиторий (например, в `docs/COMPARISON.md`) и обновлять по мере развития инструмента.
