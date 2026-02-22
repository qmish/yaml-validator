# Изменения в v1.65.0

## v1.65.0 — Docker-образ GHCR (8.7)

### Добавлено

- **Публикация в GHCR** — образ `ghcr.io/qmish/yaml-validator` собирается при push тегов `v*` и при создании GitHub Release.
- Теги: `latest`, `v1.65.0`, `1.65` (major.minor).
- Документация: docs/DISTRIBUTION.md — раздел Docker/GHCR.

### Использование

```bash
docker run --rm -v $(pwd):/workspace ghcr.io/qmish/yaml-validator:latest validate config.yaml
```
