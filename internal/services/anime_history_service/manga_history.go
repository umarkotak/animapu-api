package anime_history_service

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/contract"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository/anime_history_repository"
)

func GetHistories(ctx context.Context, user models.User, pagination models.Pagination) ([]contract.AnimeHistory, error) {
	animeHistories, err := anime_history_repository.GetByUserIDDetailed(ctx, user.ID, pagination)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return []contract.AnimeHistory{}, err
	}

	resAnimeHistories := []contract.AnimeHistory{}
	for _, animeHistory := range animeHistories {
		resMangaHistory := contract.AnimeHistory{
			ID:               animeHistory.AnimeSourceID,
			Source:           animeHistory.AnimeSource,
			Title:            animeHistory.AnimeTitle,
			LatestEpisode:    animeHistory.AnimeLatestEpisode,
			CoverUrls:        animeHistory.AnimeCoverUrls,
			LastEpisodeWatch: animeHistory.EpisodeNumber,
			LastLink:         animeHistory.FrontendPath,
			IsInLibrary:      false,
			LastWatchAt:      animeHistory.UpdatedAt,
		}

		resAnimeHistories = append(resAnimeHistories, resMangaHistory)
	}

	return resAnimeHistories, nil
}
