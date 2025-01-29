package anime_scrapper_repository

import (
	"context"

	"github.com/umarkotak/animapu-api/internal/contract"
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

func (s *Template) GetLatest(ctx context.Context, queryParams models.AnimeQueryParams) ([]contract.Anime, error) {
	return []contract.Anime{}, nil
}

func (s *Template) GetSearch(ctx context.Context, queryParams models.AnimeQueryParams) ([]contract.Anime, error) {
	return []contract.Anime{}, nil
}

func (s *Template) GetDetail(ctx context.Context, queryParams models.AnimeQueryParams) (contract.Anime, error) {
	return contract.Anime{}, nil
}

func (s *Template) Watch(ctx context.Context, queryParams models.AnimeQueryParams) (contract.EpisodeWatch, error) {
	return contract.EpisodeWatch{}, nil
}

func (s *Template) GetPerSeason(ctx context.Context, queryParams models.AnimeQueryParams) (contract.AnimePerSeason, error) {
	return contract.AnimePerSeason{}, nil
}

func (s *Template) GetRandom(ctx context.Context, queryParams models.AnimeQueryParams) ([]contract.Anime, error) {
	return []contract.Anime{}, nil
}
