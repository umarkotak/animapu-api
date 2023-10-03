package anime_scrapper_service

import (
	"context"

	"github.com/umarkotak/animapu-api/internal/models"
)

func GetEpisode(ctx context.Context, queryParams models.AnimeQueryParams) (models.Episode, models.Meta, error) {
	episode := models.Episode{}

	return episode, models.Meta{}, nil
}
