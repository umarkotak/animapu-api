package manga_scrapper_service

import (
	"context"

	"github.com/umarkotak/animapu-api/internal/contract"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository/manga_scrapper_repository/komikindo"
	"github.com/umarkotak/animapu-api/internal/repository/manga_scrapper_repository/komiku"
	"github.com/umarkotak/animapu-api/internal/repository/manga_scrapper_repository/weebcentral"
)

type (
	MangaScrapper interface {
		GetHome(ctx context.Context, queryParams models.QueryParams) ([]contract.Manga, error)
		GetDetail(ctx context.Context, queryParams models.QueryParams) (contract.Manga, error)
		GetSearch(ctx context.Context, queryParams models.QueryParams) ([]contract.Manga, error)
		GetChapter(ctx context.Context, queryParams models.QueryParams) (contract.Chapter, error)
	}
)

func mangaScrapperGenerator(mangaSource string) (MangaScrapper, error) {
	var mangaScrapper MangaScrapper

	switch mangaSource {
	case models.SOURCE_KOMIKINDO:
		mangaScrapper := komikindo.New()
		return &mangaScrapper, nil
	case models.SOURCE_KOMIKU:
		mangaScrapper := komiku.New()
		return &mangaScrapper, nil
	case models.SOURCE_WEEB_CENTRAL:
		mangaScrapper := weebcentral.New()
		return &mangaScrapper, nil
	}

	return mangaScrapper, models.ErrMangaSourceNotFound
}
