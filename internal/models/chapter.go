package models

type (
	Chapter struct {
		ID       string   `json:"id"`
		SourceID string   `json:"source_id"`
		Source   string   `json:"source"`
		Images   []string `json:"images"`
	}
)
