package contract

import "time"

type History struct {
	MediaType    string    `json:"media_type"`
	Source       string    `json:"source"`
	SourceID     string    `json:"source_id"`
	Title        string    `json:"title"`
	CoverURLs    []string  `json:"cover_urls"`
	LatestNumber float64   `json:"latest_number"`
	Progress     float64   `json:"progress"`
	LastLink     string    `json:"last_link"`
	UpdatedAt    time.Time `json:"updated_at"`
}
