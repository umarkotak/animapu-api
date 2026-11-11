package anime_scrapper_repository

import (
	"context"

	"github.com/umarkotak/animapu-api/internal/contract"
	"github.com/umarkotak/animapu-api/internal/models"
)

type (
	AnimeScrapper interface {
		GetCapabilities(ctx context.Context, queryParams models.AnimeQueryParams) ([]contract.AnimeCapability, error)
		GetLatest(ctx context.Context, queryParams models.AnimeQueryParams) ([]contract.Anime, error)
		GetSearch(ctx context.Context, queryParams models.AnimeQueryParams) ([]contract.Anime, error)
		GetRandom(ctx context.Context, queryParams models.AnimeQueryParams) ([]contract.Anime, error)
		GetDetail(ctx context.Context, queryParams models.AnimeQueryParams) (contract.Anime, error)
		Watch(ctx context.Context, queryParams models.AnimeQueryParams) (contract.EpisodeWatch, error)
		GetPerSeason(ctx context.Context, queryParams models.AnimeQueryParams) (contract.AnimePerSeason, error)
		GetCompleted(ctx context.Context, queryParams models.AnimeQueryParams) ([]contract.Anime, error)
		GetAll(ctx context.Context, queryParams models.AnimeQueryParams) ([]contract.Anime, error)
	}
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

func (s *Template) GetCapabilities(ctx context.Context, queryParams models.AnimeQueryParams) ([]contract.AnimeCapability, error) {
	return []contract.AnimeCapability{{}}, nil
}

func (s *Template) GetCompleted(ctx context.Context, queryParams models.AnimeQueryParams) ([]contract.Anime, error) {
	return []contract.Anime{}, nil
}

func (s *Template) GetAll(ctx context.Context, queryParams models.AnimeQueryParams) ([]contract.Anime, error) {
	return []contract.Anime{}, nil
}
