package anime_history_service

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/contract"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository/anime_history_repository"
	"github.com/umarkotak/animapu-api/internal/repository/user_repository"
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

func GetUserAnimeActivities(ctx context.Context, pagination models.Pagination) (contract.UserAnimeActivityData, error) {
	userAnimeActivityData := contract.UserAnimeActivityData{
		Users: []contract.UserAnimeActivity{},
	}

	animeHistories, err := anime_history_repository.GetRecentHistories(ctx, pagination)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return contract.UserAnimeActivityData{}, err
	}

	userIDs := []int64{}
	for _, animeHistory := range animeHistories {
		userIDs = append(userIDs, animeHistory.UserID)
	}

	users, err := user_repository.GetByIDs(ctx, userIDs)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return contract.UserAnimeActivityData{}, err
	}

	userAnimeActivityIdx := map[int64]int64{}
	for idx, user := range users {
		userAnimeActivityIdx[user.ID] = int64(idx)
		userAnimeActivityData.Users = append(userAnimeActivityData.Users, contract.UserAnimeActivity{
			VisitorID:      user.VisitorId,
			Email:          user.Email.String,
			AnimeHistories: []contract.AnimeHistory{},
		})
	}

	for _, animeHistory := range animeHistories {
		resAnimeHistory := contract.AnimeHistory{
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

		targetIdx := userAnimeActivityIdx[animeHistory.UserID]

		userAnimeActivityData.Users[targetIdx].AnimeHistories = append(userAnimeActivityData.Users[targetIdx].AnimeHistories, resAnimeHistory)
	}

	return userAnimeActivityData, nil
}
