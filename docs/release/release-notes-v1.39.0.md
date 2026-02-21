# Изменения в v1.39.0

## v1.39.0

- **Конфиг по имени файла** (4.5) — в yaml-validator.yaml добавлена секция `file_profiles`: список `{pattern, config}` для автовыбора профиля по пути к файлу. Пример: `**/k8s/**` → configs/k8s-strict.yaml, `*docker-compose*.yaml` → configs/docker-compose.yaml. Референс: `configs/file-profiles.yaml`.
