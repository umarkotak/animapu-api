package contract

import (
	"github.com/umarkotak/animapu-api/internal/models"
)

type (
	MangaLibraryParams struct {
		UserID            int64  //
		Sort              string // Enum: latest_update, recent_added
		models.Pagination        //
	}
)
