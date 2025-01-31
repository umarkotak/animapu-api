package models

import (
	"database/sql"
	"time"
)

type (
	AffiliateLink struct {
		ID        int64        `json:"id" db:"id"`                 //
		CreatedAt time.Time    `json:"created_at" db:"created_at"` //
		UpdatedAt time.Time    `json:"updated_at" db:"updated_at"` //
		DeletedAt sql.NullTime `json:"-" db:"deleted_at"`          //

		ShortLink string `json:"short_link" db:"short_link"` //
		LongLink  string `json:"long_link" db:"long_link"`   //
		ImageUrl  string `json:"image_url" db:"image_url"`   //
		Name      string `json:"name" db:"name"`             //
	}
)
