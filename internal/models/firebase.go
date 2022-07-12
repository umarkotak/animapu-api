package models

import "time"

type FbMangaHomeCache struct {
	UpdatedAt     time.Time `json:"updated_at"`
	UpdatedAtUnix int64     `json:"updated_at_unix"`
	ExpiredAt     time.Time `json:"expired_at"`
	Mangas        []Manga   `json:"mangas"`
}

type FbMangaDetailCache struct {
	UpdatedAt     time.Time `json:"updated_at"`
	UpdatedAtUnix int64     `json:"updated_at_unix"`
	ExpiredAt     time.Time `json:"expired_at"`
	Manga         Manga     `json:"manga"`
}

type FbGenericCache struct {
}
