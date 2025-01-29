package models

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type (
	MangaChapter struct {
		ID        int64        `json:"id" db:"id"`                 //
		CreatedAt time.Time    `json:"created_at" db:"created_at"` //
		UpdatedAt time.Time    `json:"updated_at" db:"updated_at"` //
		DeletedAt sql.NullTime `json:"deleted_at" db:"deleted_at"` //

		MangaID         int64          `json:"manga_id" db:"manga_id"`                   //
		SourceChapterID string         `json:"source_chapter_id" db:"source_chapter_id"` //
		ChapterNumber   float64        `json:"chapter_number" db:"chapter_number"`       //
		ImageUrls       pq.StringArray `json:"image_urls" db:"image_urls"`               //
	}
)
