# Артефакты релизов

При запуске `.\release.ps1 v1.x.x` сюда помещаются:

- `release-v1.x.x/` — бинарники для linux/darwin/windows (amd64/arm64)
- `yaml-validator-v1.x.x.tar.gz` — архив бинарников

Эти файлы игнорируются git (см. `.gitignore`). Артефакты публикуются на [GitHub Releases](https://github.com/qmish/yaml-validator/releases).
