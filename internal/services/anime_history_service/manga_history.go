package anime_history_service

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/contract"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository/anime_history_repository"
	"github.com/umarkotak/animapu-api/internal/repository/manga_history_repository"
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

func GetUserAnimeActivities(ctx context.Context, pagination models.Pagination) (contract.UserMangaActivityData, error) {
	userMangaActivityData := contract.UserMangaActivityData{
		Users: []contract.UserMangaActivity{},
	}

	mangaHistories, err := manga_history_repository.GetRecentHistories(ctx, pagination)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return contract.UserMangaActivityData{}, err
	}

	userIDs := []int64{}
	for _, mangaHistory := range mangaHistories {
		userIDs = append(userIDs, mangaHistory.UserID)
	}

	users, err := user_repository.GetByIDs(ctx, userIDs)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return contract.UserMangaActivityData{}, err
	}

	userMangaActivityIdx := map[int64]int64{}
	for idx, user := range users {
		userMangaActivityIdx[user.ID] = int64(idx)
		userMangaActivityData.Users = append(userMangaActivityData.Users, contract.UserMangaActivity{
			VisitorID:      user.VisitorId,
			Email:          user.Email.String,
			MangaHistories: []contract.MangaHistory{},
		})
	}

	for _, mangaHistory := range mangaHistories {
		coverImage := contract.CoverImage{
			Index:     0,
			ImageUrls: mangaHistory.MangaCoverUrls,
		}

		resMangaHistory := contract.MangaHistory{
			Source:              mangaHistory.MangaSource,
			SourceID:            mangaHistory.MangaSourceID,
			Title:               mangaHistory.MangaTitle,
			LatestChapterNumber: mangaHistory.MangaLatestChapter,
			CoverImages:         []contract.CoverImage{coverImage},
			LastChapterRead:     mangaHistory.ChapterNumber,
			LastLink:            mangaHistory.FrontendPath,
			LastReadAt:          mangaHistory.UpdatedAt,
		}

		targetIdx := userMangaActivityIdx[mangaHistory.UserID]

		userMangaActivityData.Users[targetIdx].MangaHistories = append(userMangaActivityData.Users[targetIdx].MangaHistories, resMangaHistory)
	}

	return userMangaActivityData, nil
}
