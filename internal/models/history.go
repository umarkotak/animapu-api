package models

import (
	"time"

	"github.com/lib/pq"
)

type History struct {
	MediaType    string         `json:"media_type" db:"media_type"`
	Source       string         `json:"source" db:"source"`
	SourceID     string         `json:"source_id" db:"source_id"`
	Title        string         `json:"title" db:"title"`
	CoverURLs    pq.StringArray `json:"cover_urls" db:"cover_urls"`
	LatestNumber float64        `json:"latest_number" db:"latest_number"`
	Progress     float64        `json:"progress" db:"progress"`
	LastLink     string         `json:"last_link" db:"last_link"`
	UpdatedAt    time.Time      `json:"updated_at" db:"updated_at"`
}
