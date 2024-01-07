package anime_scrapper_repository

import (
	"context"

	"github.com/umarkotak/animapu-api/internal/models"
)

type Template struct {
	AnimapuSource string
}

func NewTemplate() Template {
	return Template{
		AnimapuSource: "",
	}
}

func (s *Template) GetLatest(ctx context.Context, queryParams models.AnimeQueryParams) ([]models.Anime, error) {
	return []models.Anime{}, nil
}

func (s *Template) GetSearch(ctx context.Context, queryParams models.AnimeQueryParams) ([]models.Anime, error) {
	return []models.Anime{}, nil
}

func (s *Template) GetDetail(ctx context.Context, queryParams models.AnimeQueryParams) (models.Anime, error) {
	return models.Anime{}, nil
}

func (s *Template) Watch(ctx context.Context, queryParams models.AnimeQueryParams) (models.EpisodeWatch, error) {
	return models.EpisodeWatch{}, nil
}

func (s *Template) GetPerSeason(ctx context.Context, queryParams models.AnimeQueryParams) (models.AnimePerSeason, error) {
	return models.AnimePerSeason{}, nil
}

func (s *Template) GetRandom(ctx context.Context, queryParams models.AnimeQueryParams) ([]models.Anime, error) {
	return []models.Anime{}, nil
}
