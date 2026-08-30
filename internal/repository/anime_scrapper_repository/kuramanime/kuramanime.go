package kuramanime

import (
	"context"

	"github.com/umarkotak/animapu-api/internal/contract"
	"github.com/umarkotak/animapu-api/internal/models"
)

type Kuramanime struct {
	KuramanimeAggregator string
}

func New() Kuramanime {
	return Kuramanime{
		KuramanimeAggregator: "https://kuramalink.app",
	}
}

func (s *Kuramanime) GetLatest(ctx context.Context, queryParams models.AnimeQueryParams) ([]contract.Anime, error) {
	return []contract.Anime{}, nil
}

func (s *Kuramanime) GetSearch(ctx context.Context, queryParams models.AnimeQueryParams) ([]contract.Anime, error) {
	return []contract.Anime{}, nil
}

func (s *Kuramanime) GetDetail(ctx context.Context, queryParams models.AnimeQueryParams) (contract.Anime, error) {
	return contract.Anime{}, nil
}

func (s *Kuramanime) Watch(ctx context.Context, queryParams models.AnimeQueryParams) (contract.EpisodeWatch, error) {
	return contract.EpisodeWatch{}, nil
}

func (s *Kuramanime) GetPerSeason(ctx context.Context, queryParams models.AnimeQueryParams) (contract.AnimePerSeason, error) {
	return contract.AnimePerSeason{}, nil
}

func (s *Kuramanime) GetRandom(ctx context.Context, queryParams models.AnimeQueryParams) ([]contract.Anime, error) {
	return []contract.Anime{}, nil
}
