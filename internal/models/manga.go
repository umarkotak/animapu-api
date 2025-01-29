package models

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type (
	Manga struct {
		ID        int64        `json:"id" db:"id"`                 //
		CreatedAt time.Time    `json:"created_at" db:"created_at"` //
		UpdatedAt time.Time    `json:"updated_at" db:"updated_at"` //
		DeletedAt sql.NullTime `json:"deleted_at" db:"deleted_at"` //

		Source        string         `json:"source" db:"source"`                 //
		SourceID      string         `json:"source_id" db:"source_id"`           //
		Title         string         `json:"title" db:"title"`                   //
		CoverUrls     pq.StringArray `json:"cover_urls" db:"cover_urls"`         //
		LatestChapter float64        `json:"latest_chapter" db:"latest_chapter"` //
	}
)
