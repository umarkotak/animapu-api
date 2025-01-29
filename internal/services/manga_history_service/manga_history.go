package manga_history_service

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/contract"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository/manga_history_repository"
)

func GetHistories(ctx context.Context, user models.User, limit, page int64) ([]contract.MangaHistory, error) {
	mangaHistories, err := manga_history_repository.GetByUserIDDetailed(ctx, user.ID)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return []contract.MangaHistory{}, err
	}

	resMangaHistories := []contract.MangaHistory{}
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
		}

		resMangaHistories = append(resMangaHistories, resMangaHistory)
	}

	return resMangaHistories, nil
}
