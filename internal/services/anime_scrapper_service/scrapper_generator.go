package anime_scrapper_service

import (
	"context"

	"github.com/umarkotak/animapu-api/internal/contract"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository/anime_scrapper_repository/kuramanime"
	"github.com/umarkotak/animapu-api/internal/repository/anime_scrapper_repository/otakudesu"
)

type (
	AnimeScrapper interface {
		GetLatest(ctx context.Context, queryParams models.AnimeQueryParams) ([]contract.Anime, error)
		GetSearch(ctx context.Context, queryParams models.AnimeQueryParams) ([]contract.Anime, error)
		GetRandom(ctx context.Context, queryParams models.AnimeQueryParams) ([]contract.Anime, error)
		GetDetail(ctx context.Context, queryParams models.AnimeQueryParams) (contract.Anime, error)
		Watch(ctx context.Context, queryParams models.AnimeQueryParams) (contract.EpisodeWatch, error)
		GetPerSeason(ctx context.Context, queryParams models.AnimeQueryParams) (contract.AnimePerSeason, error)
	}
)

func animeScrapperGenerator(animeSource string) (AnimeScrapper, error) {
	var animeScrapper AnimeScrapper

	switch animeSource {
	case models.ANIME_SOURCE_OTAKUDESU:
		animeScrapper := otakudesu.New()
		return &animeScrapper, nil
	case models.ANIME_SOURCE_KURAMANIME:
		animeScrapper := kuramanime.New()
		return &animeScrapper, nil
	}

	return animeScrapper, models.ErrAnimeSourceNotFound
}
