package models

type (
	QueryParams struct {
		Source            string `json:"source"`
		SourceID          string `json:"source_id"`
		SecondarySourceID string `json:"secondary_source_id"`
		Page              int64  `json:"page"`
		ChapterID         string `json:"chapter_id"`
		Title             string `json:"title"`
	}
)
