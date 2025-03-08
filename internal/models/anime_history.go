package models

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type (
	AnimeHistory struct {
		ID        int64        `json:"id" db:"id"`                 //
		CreatedAt time.Time    `json:"created_at" db:"created_at"` //
		UpdatedAt time.Time    `json:"updated_at" db:"updated_at"` //
		DeletedAt sql.NullTime `json:"deleted_at" db:"deleted_at"` //

		UserID          int64   `json:"user_id" db:"user_id"`                     //
		AnimeID         int64   `json:"anime_id" db:"anime_id"`                   //
		EpisodeNumber   float64 `json:"episode_number" db:"episode_number"`       //
		SourceEpisodeID string  `json:"source_episode_id" db:"source_episode_id"` //
		FrontendPath    string  `json:"frontend_path" db:"frontend_path"`         //
	}

	AnimeHistoryDetailed struct {
		ID        int64        `json:"id" db:"id"`                 //
		CreatedAt time.Time    `json:"created_at" db:"created_at"` //
		UpdatedAt time.Time    `json:"updated_at" db:"updated_at"` //
		DeletedAt sql.NullTime `json:"deleted_at" db:"deleted_at"` //

		UserID          int64   `json:"user_id" db:"user_id"`                     //
		AnimeID         int64   `json:"anime_id" db:"anime_id"`                   //
		EpisodeNumber   float64 `json:"episode_number" db:"episode_number"`       //
		SourceEpisodeID string  `json:"source_episode_id" db:"source_episode_id"` //
		FrontendPath    string  `json:"frontend_path" db:"frontend_path"`         //

		AnimeSource        string         `json:"anime_source" db:"anime_source"`                 //
		AnimeSourceID      string         `json:"anime_source_id" db:"anime_source_id"`           //
		AnimeTitle         string         `json:"anime_title" db:"anime_title"`                   //
		AnimeCoverUrls     pq.StringArray `json:"anime_cover_urls" db:"anime_cover_urls"`         //
		AnimeLatestEpisode float64        `json:"anime_latest_episode" db:"anime_latest_episode"` //
	}
)
