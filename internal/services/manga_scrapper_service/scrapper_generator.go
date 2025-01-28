package manga_scrapper_service

import (
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository/manga_scrapper_repository"
	"github.com/umarkotak/animapu-api/internal/repository/mangamee_port"
)

func mangaScrapperGenerator(mangaSource string) (models.MangaScrapper, error) {
	var mangaScrapper models.MangaScrapper

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

	case models.SOURCE_M_MANGABAT:
		mangaScrapper := mangamee_port.NewMangabat()
		return &mangaScrapper, nil
	case models.SOURCE_MANGAREAD:
		mangaScrapper := mangamee_port.NewMangaread()
		return &mangaScrapper, nil
	case models.SOURCE_MANGATOWN:
		mangaScrapper := mangamee_port.NewMangatown()
		return &mangaScrapper, nil
	case models.SOURCE_MAIDMY:
		mangaScrapper := mangamee_port.NewMaidmy()
		return &mangaScrapper, nil
	case models.SOURCE_ASURA_COMIC:
		mangaScrapper := mangamee_port.NewAsuraComic()
		return &mangaScrapper, nil
	case models.SOURCE_MANGANATO:
		mangaScrapper := mangamee_port.NewMangaNato()
		return &mangaScrapper, nil
	case models.SOURCE_MANGANELO:
		mangaScrapper := mangamee_port.NewMangaNelo()
		return &mangaScrapper, nil
	case models.SOURCE_M_MANGASEE:
		mangaScrapper := mangamee_port.NewMangasee()
		return &mangaScrapper, nil

	case models.SOURCE_KLIKMANGA:
		mangaScrapper := manga_scrapper_repository.NewKlikmanga()
		return &mangaScrapper, nil
	case models.SOURCE_WEBTOONSID:
		mangaScrapper := manga_scrapper_repository.NewWebtoonsid()
		return &mangaScrapper, nil
	}

	return mangaScrapper, models.ErrMangaSourceNotFound
}
