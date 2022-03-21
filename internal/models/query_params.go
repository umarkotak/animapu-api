package models

type (
	QueryParams struct {
		Source            string `json:"source"`
		SourceID          string `json:"source_id"`
		SecondarySourceID string `json:"secondary_source_id"`
		Page              int64  `json:"page"`
	}
)
