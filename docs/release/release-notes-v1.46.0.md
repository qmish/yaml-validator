# Изменения в v1.46.0

## v1.46.0

- **Плагин Ansible** (6.4) — добавлен плагин проверки структуры Ansible playbook: hosts или name в каждом play, tasks (module/include/role), roles (role spec должен иметь ключ `role`). Пропускает K8s и docker-compose файлы.
