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
	}
)
