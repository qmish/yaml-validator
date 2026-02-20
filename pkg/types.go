package pkg

// Error представляет ошибку валидации
type Error struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Path    string `json:"path,omitempty"`
	Line    int    `json:"line,omitempty"`
	Column  int    `json:"column,omitempty"`
}
