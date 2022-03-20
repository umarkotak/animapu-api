package models

import (
	"errors"
)

type AnimapuError struct {
	StatusCode int    `json:"status_code"`
	ErrorCode  string `json:"error_code"`
	Message    string `json:"message"`
}

var (
	ErrInternal = errors.New("unknown error")
)

var (
	ERROR_MAP = map[error]AnimapuError{
		ErrInternal: {
			StatusCode: 500,
			ErrorCode:  "Internal server error",
			Message:    "Terjadi kesalahan internal pada system, silahkan coba lagi.",
		},
	}
)
