package anime_scrapper_service

import (
	"github.com/umarkotak/animapu-api/internal/models"
	anime_scrapper_animension_local "github.com/umarkotak/animapu-api/internal/repository/anime_scrapper_repository/animension_local"
	anime_scrapper_otakudesu "github.com/umarkotak/animapu-api/internal/repository/anime_scrapper_repository/otakudesu"
)

func animeScrapperGenerator(animeSource string) (models.AnimeScrapper, error) {
	var animeScrapper models.AnimeScrapper

	switch animeSource {
	case models.ANIME_SOURCE_OTAKUDESU:
		animeScrapper := anime_scrapper_otakudesu.NewOtakudesu()
		return &animeScrapper, nil
	case models.ANIME_SOURCE_ANIMENSION_LOCAL:
		animeScrapper := anime_scrapper_animension_local.NewAnimensionLocal()
		return &animeScrapper, nil
	case models.ANIME_SOURCE_ANIMENSION:
		animeScrapper := anime_scrapper_animension_local.NewAnimensionLocal()
		return &animeScrapper, nil
	}

	return animeScrapper, models.ErrAnimeSourceNotFound
}
