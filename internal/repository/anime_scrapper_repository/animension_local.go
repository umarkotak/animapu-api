package anime_scrapper_repository

import (
	"context"

	"github.com/umarkotak/animapu-api/internal/models"
)

type AnimensionLocal struct {
	AnimapuSource   string
	Source          string
	AnimensionHost  string
	DesusStreamHost string
}

func NewAnimensionLocal() AnimensionLocal {
	return AnimensionLocal{
		AnimapuSource:  models.ANIME_SOURCE_ANIMENSION_LOCAL,
		Source:         "animension",
		AnimensionHost: "https://animension.to",
	}
}

func (s *AnimensionLocal) GetLatest(ctx context.Context, queryParams models.AnimeQueryParams) ([]models.Anime, error) {
	return []models.Anime{}, nil
}

func (s *AnimensionLocal) GetDetail(ctx context.Context, queryParams models.AnimeQueryParams) (models.Anime, error) {
	return models.Anime{}, nil
}

func (s *AnimensionLocal) Watch(ctx context.Context, queryParams models.AnimeQueryParams) (models.EpisodeWatch, error) {
	return models.EpisodeWatch{}, nil
}

func (s *AnimensionLocal) GetPerSeason(ctx context.Context, queryParams models.AnimeQueryParams) (models.AnimePerSeason, error) {
	return models.AnimePerSeason{}, nil
}
