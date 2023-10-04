package anime_scrapper_service

import (
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository/anime_scrapper_repository"
)

func animeScrapperGenerator(animeSource string) (models.AnimeScrapper, error) {
	var animeScrapper models.AnimeScrapper

	switch animeSource {
	case models.ANIME_SOURCE_OTAKUDESU:
		animeScrapper := anime_scrapper_repository.NewOtakudesu()
		return &animeScrapper, nil
	}

	return animeScrapper, models.ErrAnimeSourceNotFound
}
