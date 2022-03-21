package models

type (
	QueryParams struct {
		MangaSource string `json:"manga_source"`
		MangaID     string `json:"manga_id"`
		Page        int64  `json:"page"`
	}
)
