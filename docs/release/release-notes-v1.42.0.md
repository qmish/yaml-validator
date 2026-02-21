# Изменения в v1.42.0

## v1.42.0

- **Запрет значений по умолчанию** (5.4) — добавлено правило `forbid_default_values` и `default_values` в конфиг: запрет явного указания ключей со значением по умолчанию (например, для K8s: `imagePullPolicy: Always`, `restartPolicy: Always`, `terminationGracePeriodSeconds: 30`). Пример конфига: `configs/k8s-forbid-defaults.yaml`.
