package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Level уровень логирования
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

var (
	currentLevel = LevelInfo
	useJSON      = false
)

// SetLevel устанавливает минимальный уровень логирования
func SetLevel(level Level) {
	currentLevel = level
}

// SetJSON включает вывод логов в формате JSON
func SetJSON(enabled bool) {
	useJSON = enabled
}

// Debug выводит отладочное сообщение
func Debug(format string, args ...interface{}) {
	if currentLevel <= LevelDebug {
		write(LevelDebug, format, args...)
	}
}

// Info выводит информационное сообщение
func Info(format string, args ...interface{}) {
	if currentLevel <= LevelInfo {
		write(LevelInfo, format, args...)
	}
}

// Warn выводит предупреждение
func Warn(format string, args ...interface{}) {
	if currentLevel <= LevelWarn {
		write(LevelWarn, format, args...)
	}
}

// Error выводит сообщение об ошибке
func Error(format string, args ...interface{}) {
	if currentLevel <= LevelError {
		write(LevelError, format, args...)
	}
}

func write(level Level, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	levelStr := []string{"DEBUG", "INFO", "WARN", "ERROR"}[level]
	timestamp := time.Now().Format(time.RFC3339)

	if useJSON {
		entry := map[string]string{
			"timestamp": timestamp,
			"level":     levelStr,
			"message":   msg,
		}
		data, _ := json.Marshal(entry)
		fmt.Fprintln(os.Stderr, string(data))
	} else {
		fmt.Fprintf(os.Stderr, "[%s] %s: %s\n", timestamp, levelStr, msg)
	}
}
