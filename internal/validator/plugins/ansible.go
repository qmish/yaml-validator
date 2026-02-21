package plugins

import (
	"yaml-validator/internal/validator"
	"yaml-validator/pkg"

	"gopkg.in/yaml.v3"
)

// AnsibleValidator проверяет структуру Ansible playbook (hosts, tasks, vars)
type AnsibleValidator struct{}

func init() {
	validator.RegisterPlugin("ansible", &AnsibleValidator{})
}

// Name возвращает имя плагина
func (a *AnsibleValidator) Name() string {
	return "AnsibleValidator"
}

// Validate проверяет Ansible playbook: hosts в каждом play, структура tasks
func (a *AnsibleValidator) Validate(node *yaml.Node) []pkg.Error {
	var errors []pkg.Error

	if node == nil {
		return errors
	}

	// Root может быть DocumentNode -> один child (sequence или mapping)
	root := node
	if node.Kind == yaml.DocumentNode && len(node.Content) > 0 {
		root = node.Content[0]
	}

	// Ansible playbook: list of plays или single play (mapping с hosts)
	if root.Kind == yaml.SequenceNode {
		// playbook: список plays
		for idx, playNode := range root.Content {
			if playNode.Kind == yaml.MappingNode {
				playErrs := validatePlay(playNode, idx)
				errors = append(errors, playErrs...)
			}
		}
		return errors
	}

	if root.Kind == yaml.MappingNode {
		fields := buildFlatMap(root, "")
		// Не Ansible: K8s или docker-compose
		if _, hasAPI := fields["apiVersion"]; hasAPI {
			return errors
		}
		if _, hasServices := fields["services"]; hasServices {
			return errors
		}
		// single play
		if _, hasHosts := fields["hosts"]; hasHosts {
			errors = append(errors, validatePlay(root, 0)...)
		}
	}

	return errors
}

func validatePlay(playNode *yaml.Node, playIdx int) []pkg.Error {
	var errors []pkg.Error

	fields := buildFlatMap(playNode, "")
	_, hasHosts := fields["hosts"]
	_, hasName := fields["name"]

	// Ansible play требует hosts (или use 'name' для include_role/import_playbook)
	if !hasHosts && !hasName {
		line := playNode.Line
		errors = append(errors, pkg.Error{
			Type:       "AnsiblePlayMissingHosts",
			Message:    "play must have 'hosts' or 'name'",
			Suggestion: "add hosts: all or hosts: <group> to play",
			Path:       "hosts",
			Line:       line,
		})
	}

	// Проверка tasks: каждый task — mapping с module или include/role
	tasksNode := findChild(playNode, "tasks")
	if tasksNode != nil && tasksNode.Kind == yaml.SequenceNode {
		for idx, taskNode := range tasksNode.Content {
			if taskNode.Kind == yaml.MappingNode {
				taskErrs := validateTask(taskNode, idx)
				errors = append(errors, taskErrs...)
			}
		}
	}

	// Проверка roles: list of strings (role name) или role specs (role: name)
	rolesNode := findChild(playNode, "roles")
	if rolesNode != nil && rolesNode.Kind == yaml.SequenceNode {
		for _, roleNode := range rolesNode.Content {
			if roleNode.Kind == yaml.MappingNode {
				// role spec: { role: name, vars: {...} }
				roleFields := buildFlatMap(roleNode, "")
				if _, hasRole := roleFields["role"]; !hasRole {
					line := roleNode.Line
					errors = append(errors, pkg.Error{
						Type:       "AnsibleRoleMissingName",
						Message:    "role spec must have 'role' key",
						Suggestion: "add role: <role_name> to role entry",
						Path:       "roles",
						Line:       line,
					})
				}
			}
			// ScalarNode — role name as string, OK
		}
	}

	return errors
}

func validateTask(taskNode *yaml.Node, taskIdx int) []pkg.Error {
	var errors []pkg.Error
	fields := buildFlatMap(taskNode, "")
	_, hasName := fields["name"]
	_, hasInclude := fields["include"]
	_, hasIncludeRole := fields["include_role"]
	_, hasImportRole := fields["import_role"]
	_, hasImportTasks := fields["import_tasks"]
	_, hasIncludeTasks := fields["include_tasks"]
	_, hasRole := fields["role"]
	_, hasBlock := fields["block"]
	_, hasMeta := fields["meta"]
	hasModule := hasModuleKey(fields)

	// Task должен иметь name и/или module/include/role/block
	if !hasName && !hasInclude && !hasIncludeRole && !hasImportRole && !hasImportTasks && !hasIncludeTasks && !hasRole && !hasBlock && !hasMeta && !hasModule {
		line := taskNode.Line
		errors = append(errors, pkg.Error{
			Type:       "AnsibleTaskNoModule",
			Message:    "task should have 'name' and a module (e.g. debug, copy, shell) or include/role",
			Suggestion: "add name: <description> and a module key to task",
			Path:       "tasks",
			Line:       line,
		})
	}

	_ = taskIdx
	return errors
}

func hasModuleKey(fields map[string]string) bool {
	// Ansible modules: любая ключ, не входящий в зарезервированные
	reserved := map[string]bool{
		"name": true, "include": true, "include_role": true, "import_role": true,
		"import_tasks": true, "include_tasks": true, "role": true, "block": true,
		"meta": true, "when": true, "loop": true, "vars": true, "tags": true,
	}
	for k := range fields {
		if !reserved[k] && k != "" {
			return true // есть не-зарезервированный ключ — скорее всего module
		}
	}
	return false
}
