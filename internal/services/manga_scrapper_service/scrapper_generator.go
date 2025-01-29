package manga_scrapper_service

import (
	"context"

	"github.com/umarkotak/animapu-api/internal/contract"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository/manga_scrapper_repository"
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
	case models.SOURCE_MANGABAT:
		mangaScrapper := manga_scrapper_repository.NewMangabat()
		return &mangaScrapper, nil
	case models.SOURCE_ASURA_NACM:
		mangaScrapper := manga_scrapper_repository.NewAsuraComic()
		return &mangaScrapper, nil
	case models.SOURCE_KOMIKINDO:
		mangaScrapper := manga_scrapper_repository.NewKomikindo()
		return &mangaScrapper, nil
	case models.SOURCE_KOMIKU:
		mangaScrapper := manga_scrapper_repository.NewKomiku()
		return &mangaScrapper, nil
	case models.SOURCE_KOMIKCAST:
		mangaScrapper := manga_scrapper_repository.NewKomikcast()
		return &mangaScrapper, nil
	case models.SOURCE_MANGASEE:
		mangaScrapper := manga_scrapper_repository.NewMangasee()
		return &mangaScrapper, nil
	case models.SOURCE_WEEB_CENTRAL:
		mangaScrapper := manga_scrapper_repository.NewWeebCentral()
		return &mangaScrapper, nil
	}

	return mangaScrapper, models.ErrMangaSourceNotFound
}
