package models

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type (
	MangaHistory struct {
		ID        int64        `json:"id" db:"id"`                 //
		CreatedAt time.Time    `json:"created_at" db:"created_at"` //
		UpdatedAt time.Time    `json:"updated_at" db:"updated_at"` //
		DeletedAt sql.NullTime `json:"deleted_at" db:"deleted_at"` //

		UserID          int64   `json:"user_id" db:"user_id"`                     //
		MangaID         int64   `json:"manga_id" db:"manga_id"`                   //
		ChapterNumber   float64 `json:"chapter_number" db:"chapter_number"`       //
		SourceChapterID string  `json:"source_chapter_id" db:"source_chapter_id"` //
		FrontendPath    string  `json:"frontend_path" db:"frontend_path"`         //
	}

	MangaHistoryDetailed struct {
		ID        int64        `json:"id" db:"id"`                 //
		CreatedAt time.Time    `json:"created_at" db:"created_at"` //
		UpdatedAt time.Time    `json:"updated_at" db:"updated_at"` //
		DeletedAt sql.NullTime `json:"deleted_at" db:"deleted_at"` //

		UserID          int64   `json:"user_id" db:"user_id"`                     //
		MangaID         int64   `json:"manga_id" db:"manga_id"`                   //
		ChapterNumber   float64 `json:"chapter_number" db:"chapter_number"`       //
		SourceChapterID string  `json:"source_chapter_id" db:"source_chapter_id"` //
		FrontendPath    string  `json:"frontend_path" db:"frontend_path"`         //

		MangaSource        string         `json:"manga_source" db:"manga_source"`                 //
		MangaSourceID      string         `json:"manga_source_id" db:"manga_source_id"`           //
		MangaTitle         string         `json:"manga_title" db:"manga_title"`                   //
		MangaCoverUrls     pq.StringArray `json:"manga_cover_urls" db:"manga_cover_urls"`         //
		MangaLatestChapter float64        `json:"manga_latest_chapter" db:"manga_latest_chapter"` //
		MangaUpdatedAt     time.Time      `json:"manga_updated_at" db:"manga_updated_at"`         //
	}
)
