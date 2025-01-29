package anime_scrapper_service

import (
	"context"

	"github.com/umarkotak/animapu-api/internal/contract"
	"github.com/umarkotak/animapu-api/internal/models"
	anime_scrapper_animeindo "github.com/umarkotak/animapu-api/internal/repository/anime_scrapper_repository/animeindo"
	anime_scrapper_gogo_anime "github.com/umarkotak/animapu-api/internal/repository/anime_scrapper_repository/gogo_anime"
	anime_scrapper_gogo_anime_new "github.com/umarkotak/animapu-api/internal/repository/anime_scrapper_repository/gogo_anime_new"
	anime_scrapper_otakudesu "github.com/umarkotak/animapu-api/internal/repository/anime_scrapper_repository/otakudesu"
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
		animeScrapper := anime_scrapper_otakudesu.NewOtakudesu()
		return &animeScrapper, nil
	case models.ANIME_SOURCE_ANIMEINDO:
		animeScrapper := anime_scrapper_animeindo.NewAnimeindo()
		return &animeScrapper, nil
	case models.ANIME_SOURCE_GOGO_ANIME:
		animeScrapper := anime_scrapper_gogo_anime.NewGogoAnime()
		return &animeScrapper, nil
	case models.ANIME_SOURCE_GOGO_ANIME_NEW:
		animeScrapper := anime_scrapper_gogo_anime_new.NewGogoAnimeNew()
		return &animeScrapper, nil
	}

	return animeScrapper, models.ErrAnimeSourceNotFound
}
