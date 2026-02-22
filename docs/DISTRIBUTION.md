# Пакеты и дистрибуция

Инструкции и скрипты для установки yaml-validator из пакетов и сборки .deb, .rpm, Homebrew, Chocolatey.

## Docker / GitHub Container Registry

Официальный образ публикуется в GHCR при push тегов `v*` и при создании GitHub Release.

```bash
# Запуск (тег = версия, например v1.64.0)
docker run --rm -v $(pwd):/workspace ghcr.io/qmish/yaml-validator:latest validate config.yaml

# Указать версию
docker run --rm -v $(pwd):/workspace ghcr.io/qmish/yaml-validator:v1.64.0 validate **/*.yaml
```

Теги: `latest`, `v1.64.0`, `1.64` (major.minor).

---

## Скачивание релиза

Бинарники публикуются на [GitHub Releases](https://github.com/qmish/yaml-validator/releases).

```bash
# Linux (amd64)
curl -sSL https://github.com/qmish/yaml-validator/releases/download/v1.50.0/yaml-validator-v1.50.0-linux-amd64.tar.gz | tar xz -C /usr/local/bin --strip-components=1

# Linux (arm64)
curl -sSL https://github.com/qmish/yaml-validator/releases/download/v1.50.0/yaml-validator-v1.50.0-linux-arm64.tar.gz | tar xz -C /usr/local/bin --strip-components=1

# macOS (darwin)
curl -sSL https://github.com/qmish/yaml-validator/releases/download/v1.50.0/yaml-validator-v1.50.0-darwin-amd64.tar.gz | tar xz -C /usr/local/bin --strip-components=1
# Apple Silicon:
curl -sSL https://github.com/qmish/yaml-validator/releases/download/v1.50.0/yaml-validator-v1.50.0-darwin-arm64.tar.gz | tar xz -C /usr/local/bin --strip-components=1

# Windows — распаковать .tar.gz и добавить путь к exe в PATH
```

Архив содержит бинарники для всех платформ. Для установки одного файла скачайте нужный бинарник из Assets релиза.

---

## Debian / Ubuntu (.deb)

### Сборка пакета (fpm)

Установите [fpm](https://github.com/jordansissel/fpm): `gem install fpm`.

```bash
# Из корня репозитория после release.ps1 или build.sh
VERSION=v1.50.0
./scripts/packaging/build-deb.sh $VERSION
# Результат: dist/yaml-validator_1.50.0_amd64.deb
```

### Установка

```bash
sudo dpkg -i yaml-validator_1.50.0_amd64.deb
# Зависимости: libc6
```

---

## Red Hat / Fedora / CentOS (.rpm)

### Сборка пакета (fpm)

```bash
VERSION=v1.50.0
./scripts/packaging/build-rpm.sh $VERSION
# Результат: dist/yaml-validator-1.50.0-1.x86_64.rpm
```

### Установка

```bash
sudo rpm -i yaml-validator-1.50.0-1.x86_64.rpm
```

---

## Homebrew (macOS / Linux)

### Установка из tap (если есть)

```bash
brew tap qmish/yaml-validator
brew install yaml-validator
```

### Установка напрямую из URL

```bash
# macOS Intel
brew install --formula https://raw.githubusercontent.com/qmish/yaml-validator/main/docs/packaging/yaml-validator.rb

# Или вручную скачать и установить бинарник:
VERSION=v1.50.0
curl -sSL "https://github.com/qmish/yaml-validator/releases/download/${VERSION}/yaml-validator-${VERSION}-darwin-arm64.tar.gz" | tar xz -C /usr/local/bin --strip-components=1
chmod +x /usr/local/bin/yaml-validator
```

### Формула Homebrew

Файл формулы: `docs/packaging/yaml-validator.rb`. Для добавления в свой tap скопируйте в `Formula/yaml-validator.rb`.

---

## Chocolatey (Windows)

### Установка (если пакет опубликован)

```powershell
choco install yaml-validator
```

### Сборка пакета вручную

```powershell
# Из корня репозитория после release
VERSION=v1.50.0
.\scripts\packaging\build-chocolatey.ps1 -Version $VERSION
# Результат: dist/yaml-validator.1.50.0.nupkg
choco pack docs/packaging/yaml-validator.nuspec
```

### Установка из .nupkg

```powershell
choco install yaml-validator -s .
```

---

## Рекомендуемые версии URL

При ссылке на релизы подставляйте актуальную версию: https://github.com/qmish/yaml-validator/releases
