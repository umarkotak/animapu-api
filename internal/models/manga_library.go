package models

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type (
	MangaLibrary struct {
		ID        int64        `json:"id" db:"id"`                 //
		CreatedAt time.Time    `json:"created_at" db:"created_at"` //
		UpdatedAt time.Time    `json:"updated_at" db:"updated_at"` //
		DeletedAt sql.NullTime `json:"deleted_at" db:"deleted_at"` //

		UserID  int64 `json:"user_id" db:"user_id"`   //
		MangaID int64 `json:"manga_id" db:"manga_id"` //
	}

	MangaLibraryDetailed struct {
		ID        int64        `json:"id" db:"id"`                 //
		CreatedAt time.Time    `json:"created_at" db:"created_at"` //
		UpdatedAt time.Time    `json:"updated_at" db:"updated_at"` //
		DeletedAt sql.NullTime `json:"deleted_at" db:"deleted_at"` //

		UserID  int64 `json:"user_id" db:"user_id"`   //
		MangaID int64 `json:"manga_id" db:"manga_id"` //

		HistoryChapterNumber   float64 `json:"chapter_number" db:"history_chapter_number"`       //
		HistorySourceChapterID string  `json:"source_chapter_id" db:"history_source_chapter_id"` //
		HistoryFrontendPath    string  `json:"frontend_path" db:"history_frontend_path"`         //

		MangaSource        string         `json:"manga_source" db:"manga_source"`                 //
		MangaSourceID      string         `json:"manga_source_id" db:"manga_source_id"`           //
		MangaTitle         string         `json:"manga_title" db:"manga_title"`                   //
		MangaCoverUrls     pq.StringArray `json:"manga_cover_urls" db:"manga_cover_urls"`         //
		MangaLatestChapter float64        `json:"manga_latest_chapter" db:"manga_latest_chapter"` //
	}
)
