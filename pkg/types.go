package pkg

// Error представляет ошибку валидации
type Error struct {
	Type       string `json:"type"`
	Message    string `json:"message"`
	Suggestion string `json:"suggestion,omitempty"` // подсказка по исправлению (как gofmt)
	Path       string `json:"path,omitempty"`
	Line       int    `json:"line,omitempty"`
	Column     int    `json:"column,omitempty"`
	Severity   string `json:"severity,omitempty"` // "error" (default) или "warning" (5.7)
}
