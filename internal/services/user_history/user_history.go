package user_history

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository"
)

func RecordHistory(ctx context.Context, user models.User, manga models.Manga, chapter models.Chapter) (models.Manga, error) {
	err := repository.RecordUserReadHistory(ctx, user, manga)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return manga, err
	}
	return manga, nil
}

func GetReadHistories(ctx context.Context, user models.User) ([]models.Manga, map[string]models.Manga, error) {
	mangaHistories, mangaHistoriesMap, err := repository.GetUserReadHistories(ctx, user)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return mangaHistories, mangaHistoriesMap, err
	}
	return mangaHistories, mangaHistoriesMap, nil
}
