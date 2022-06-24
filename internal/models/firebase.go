package models

import "time"

type FbMangahubHome struct {
	UpdatedAt     time.Time `json:"updated_at"`
	UpdatedAtUnix int64     `json:"updated_at_unix"`
	ExpiredAt     time.Time `json:"expired_at"`
	Mangas        []Manga   `json:"mangas"`
}

type FbMangahubDetail struct {
	UpdatedAt     time.Time `json:"updated_at"`
	UpdatedAtUnix int64     `json:"updated_at_unix"`
	ExpiredAt     time.Time `json:"expired_at"`
	Manga         Manga     `json:"manga"`
}
