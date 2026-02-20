# Contributing to yaml-validator

## Разработка

### Требования

- Go 1.24+ (см. go.mod)
- Make или PowerShell (для сборки)

### Клонирование и сборка

```bash
git clone https://github.com/your-org/yaml-validator.git
cd yaml-validator
go mod download
go build .
```

### Тестирование

```bash
go test ./... -v
```

### Добавление плагина

1. Создайте файл в `internal/validator/plugins/` (например, `custom.go`)
2. Реализуйте интерфейс `ValidatorPlugin`:

```go
package plugins

import (
    "yaml-validator/internal/validator"
    "yaml-validator/pkg"
    "gopkg.in/yaml.v3"
)

type CustomValidator struct{}

func init() {
    validator.RegisterPlugin("custom", &CustomValidator{})
}

func (c *CustomValidator) Name() string {
    return "CustomValidator"
}

func (c *CustomValidator) Validate(node *yaml.Node) []pkg.Error {
    var errors []pkg.Error
    // ваша логика
    return errors
}
```

3. Добавьте импорт в `cmd/validator.go`:
```go
_ "yaml-validator/internal/validator/plugins"
```

## Pull Requests

1. Создайте ветку от `main`
2. Добавьте тесты для новой функциональности
3. Убедитесь: `go test ./...` проходит
4. Опишите изменения в PR

## Версионирование

Проект использует [SemVer](https://semver.org/):
- **MAJOR** — несовместимые изменения API
- **MINOR** — новая функциональность (обратно совместимая)
- **PATCH** — исправления ошибок
