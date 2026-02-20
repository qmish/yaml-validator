package validator

import (
	"sync"

	"yaml-validator/pkg"

	"gopkg.in/yaml.v3"
)

// ValidatorPlugin интерфейс для пользовательских валидаторов
type ValidatorPlugin interface {
	Name() string
	Validate(node *yaml.Node) []pkg.Error
}

var (
	plugins   = make(map[string]ValidatorPlugin)
	pluginsMu sync.RWMutex
)

// RegisterPlugin регистрирует плагин валидации
func RegisterPlugin(name string, plugin ValidatorPlugin) {
	pluginsMu.Lock()
	defer pluginsMu.Unlock()
	plugins[name] = plugin
}

// UnregisterPlugin удаляет плагин из регистрации
func UnregisterPlugin(name string) {
	pluginsMu.Lock()
	defer pluginsMu.Unlock()
	delete(plugins, name)
}

// GetPlugins возвращает копию списка зарегистрированных плагинов
func GetPlugins() map[string]ValidatorPlugin {
	pluginsMu.RLock()
	defer pluginsMu.RUnlock()
	result := make(map[string]ValidatorPlugin, len(plugins))
	for k, v := range plugins {
		result[k] = v
	}
	return result
}

// RunPlugins запускает все зарегистрированные плагины
func RunPlugins(node *yaml.Node) []pkg.Error {
	var errors []pkg.Error
	for _, plugin := range GetPlugins() {
		errors = append(errors, plugin.Validate(node)...)
	}
	return errors
}
