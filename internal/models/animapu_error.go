package models

import (
	"errors"
)

type AnimapuError struct {
	StatusCode int    `json:"status_code"`
	ErrorCode  string `json:"error_code"`
	Message    string `json:"message"`
	RawError   string `json:"raw_error"`
}

var (
	ErrInternal                  = errors.New("unknown error")
	ErrMangaSourceNotFound       = errors.New("manga source not found")
	ErrMangaSourceNotImplemented = errors.New("manga source not implemented yet")
	ErrInvalidFormat             = errors.New("invalid format")
	ErrInvalidTargetURL          = errors.New("invalid target url")
	ErrCacheNotFound             = errors.New("cache not found")
	ErrNotFound                  = errors.New("not found")
	ErrUnauthorized              = errors.New("unauthorized")
	ErrMangamee                  = errors.New("mangamee error")

	ErrAnimeSourceNotFound          = errors.New("anime source not found")
	ErrOtakudesuFrameSourceNotFound = errors.New("otakudesu frame source not found")
)

var (
	ERROR_MAP = map[error]AnimapuError{
		ErrInternal: {
			StatusCode: 500,
			ErrorCode:  "Internal server error",
			Message:    "Terjadi kesalahan internal pada system, silahkan coba lagi.",
		},
		ErrMangaSourceNotFound: {
			StatusCode: 404,
			ErrorCode:  "Manga source not found",
			Message:    "Sumber manga yang anda cari tidak dapat ditemukan.",
		},
		ErrNotFound: {
			StatusCode: 404,
			ErrorCode:  "Not found",
			Message:    "Tidak dapat ditemukan.",
		},
		ErrMangaSourceNotImplemented: {
			StatusCode: 422,
			ErrorCode:  "Manga source not implemented",
			Message:    "Sumber yang anda pilih masih dalam proses pengerjaan.",
		},
		ErrInvalidFormat: {
			StatusCode: 400,
			ErrorCode:  "Invalid format",
			Message:    "Format anda salah.",
		},
		ErrInvalidTargetURL: {
			StatusCode: 400,
			ErrorCode:  "Invalid target url",
			Message:    "Target URL anda salah.",
		},
		ErrUnauthorized: {
			StatusCode: 401,
			ErrorCode:  "Unauthorized",
			Message:    "Tidak ada akses.",
		},
		ErrMangamee: {
			StatusCode: 422,
			ErrorCode:  "mangamee_error",
			Message:    "Terjadi kesalahan di server mangamee",
		},

		ErrAnimeSourceNotFound: {
			StatusCode: 404,
			ErrorCode:  "Anime source not found",
			Message:    "Sumber anime yang anda cari tidak dapat ditemukan.",
		},
		ErrOtakudesuFrameSourceNotFound: {
			StatusCode: 404,
			ErrorCode:  "Otakudesu frame source not found",
			Message:    "Sumber stream otakudesu tidak dapat ditemukan.",
		},
	}
)
