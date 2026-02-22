# Изменения в v1.57.0

## v1.57.0 — Плагин Terraform (8.4)

### Добавлено

- **Плагин Terraform** — проверка YAML в контексте Terraform (tfvars.yaml, backend config).
- **TerraformVarNameHyphen** — предупреждение: в Terraform имена переменных используют underscore (`instance_type`), а не hyphen (`instance-type`).
- **configs/terraform.yaml** — конфигурация для Terraform YAML (без обязательных K8s-полей).

### Использование

```bash
yaml-validator validate terraform.tfvars.yaml -c configs/terraform.yaml
```
