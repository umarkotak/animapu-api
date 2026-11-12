package anime_scrapper_service

import (
	"context"
	"slices"

	"github.com/umarkotak/animapu-api/internal/contract"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository/anime_repository"
)

func FilterNSFW(ctx context.Context, source string, animes []contract.Anime) ([]contract.Anime, error) {
	if len(animes) == 0 {
		return animes, nil
	}

	sourceIDs := make([]string, 0, len(animes))
	for _, anime := range animes {
		sourceIDs = append(sourceIDs, anime.ID)
	}
	dbAnimes, err := anime_repository.GetBySourceAndSourceIDs(ctx, source, sourceIDs)
	if err != nil {
		return nil, err
	}

	return excludeNSFW(animes, dbAnimes), nil
}

func excludeNSFW(animes []contract.Anime, dbAnimes []models.Anime) []contract.Anime {
	type animeKey struct{ source, sourceID string }
	nsfw := make(map[animeKey]bool)
	for _, anime := range dbAnimes {
		if slices.Contains(anime.Tags, "nsfw:true") {
			nsfw[animeKey{anime.Source, anime.SourceID}] = true
		}
	}

	filtered := make([]contract.Anime, 0, len(animes))
	for _, anime := range animes {
		if !nsfw[animeKey{anime.Source, anime.ID}] {
			filtered = append(filtered, anime)
		}
	}
	return filtered
}
