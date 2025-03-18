package models

import (
	"database/sql"
	"time"
)

type (
	User struct {
		ID        int64          `json:"id" db:"id"`                 //
		CreatedAt time.Time      `json:"created_at" db:"created_at"` //
		UpdatedAt time.Time      `json:"updated_at" db:"updated_at"` //
		DeletedAt sql.NullTime   `json:"deleted_at" db:"deleted_at"` //
		VisitorId string         `json:"visitor_id" db:"visitor_id"` //
		Guid      sql.NullString `json:"guid" db:"guid"`             // sha generated from front end
		Email     sql.NullString `json:"email" db:"email"`           //
	}
)
