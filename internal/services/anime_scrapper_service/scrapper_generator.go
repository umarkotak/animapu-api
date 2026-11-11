package anime_scrapper_service

import (
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository/anime_scrapper_repository"
	"github.com/umarkotak/animapu-api/internal/repository/anime_scrapper_repository/kuramanime"
	"github.com/umarkotak/animapu-api/internal/repository/anime_scrapper_repository/otakudesu"
)

func animeScrapperGenerator(animeSource string) (anime_scrapper_repository.AnimeScrapper, error) {
	var animeScrapper anime_scrapper_repository.AnimeScrapper

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
