# Плагины yaml-validator

Описание API и примеры создания плагинов для расширения валидации YAML.

## Интерфейс плагина

```go
package validator

type ValidatorPlugin interface {
    Name() string
    Validate(node *yaml.Node) []pkg.Error
}

func RegisterPlugin(name string, plugin ValidatorPlugin)
```

- **Name()** — имя плагина (отображается в сообщениях об ошибках).
- **Validate(node *yaml.Node)** — проверка YAML-дерева, возвращает список ошибок или пустой слайс.

## Регистрация плагина

1. Создайте файл в `internal/validator/plugins/` (например, `custom.go`).
2. Реализуйте интерфейс `ValidatorPlugin`.
3. Вызовите `RegisterPlugin` в `init()`.

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
    // ваша логика обхода node.Content
    return errors
}
```

4. Добавьте импорт в `cmd/validator.go`:
```go
_ "yaml-validator/internal/validator/plugins"
```

## Обход YAML-дерева

Используйте `parser.TraverseMappings` или `parser.Traverse` из пакета `internal/parser`:

```go
import "yaml-validator/internal/parser"

// Обход всех маппингов
parser.TraverseMappings(node, "", func(mappingNode *yaml.Node, path string) {
    // path — полный путь к маппингу (например, "metadata")
})

// Обход всех узлов
parser.Traverse(node, "", func(info parser.NodeInfo) {
    info.Node   // *yaml.Node
    info.Path   // путь
    info.Value  // значение узла
})
```

## Структура ошибки

```go
pkg.Error{
    Type:    "RuleName",   // тип правила
    Message: "описание",   // текст для пользователя
    Path:    "metadata",   // путь в YAML (опционально)
    Line:    5,            // строка (0 = неизвестно)
    Column:  2,            // колонка (0 = неизвестно)
}
```

## Встроенные плагины

- **kubernetes** — проверка apiVersion, kind, metadata.name и структуры K8s-манифестов.
- **docker-compose** (v1.20.0) — проверка: у каждого сервиса должно быть `image` или `build`. Запускается только для файлов с секцией `services` (без `apiVersion`).
- **ansible** (v1.46.0) — проверка структуры Ansible playbook: hosts или name в каждом play, tasks (module/include/role), roles (role spec должен иметь ключ `role`). Пропускает K8s и docker-compose файлы.

## Пример: проверка обязательного поля

```go
func (c *CustomValidator) Validate(node *yaml.Node) []pkg.Error {
    root := parser.GetRootMapping(node)
    if root == nil || root.Kind != yaml.MappingNode {
        return nil
    }
    found := false
    for i := 0; i < len(root.Content); i += 2 {
        if i+1 < len(root.Content) && root.Content[i].Value == "myRequiredField" {
            found = true
            break
        }
    }
    if !found {
        return []pkg.Error{{
            Type:    "MissingField",
            Message: "required field 'myRequiredField' is missing",
            Path:    "myRequiredField",
        }}
    }
    return nil
}
```

См. [CONTRIBUTING.md](../CONTRIBUTING.md) для полного примера и инструкций по PR.
