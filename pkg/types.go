package pkg

// Error представляет ошибку валидации
type Error struct {
	Type          string `json:"type"`
	Message       string `json:"message"`
	Suggestion    string `json:"suggestion,omitempty"`
	Path          string `json:"path,omitempty"`
	Line          int    `json:"line,omitempty"`
	Column        int    `json:"column,omitempty"`
	Severity      string `json:"severity,omitempty"`
	DocumentIndex int    `json:"document_index,omitempty"` // для мультидокумента: индекс документа (1-based)
}
