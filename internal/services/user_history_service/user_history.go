package user_history_service

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/contract"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository"
)

func FirebaseRecordHistory(ctx context.Context, user models.UserFirebase, manga contract.Manga, chapter contract.Chapter) (contract.Manga, error) {
	err := repository.FirebaseRecordUserReadHistory(ctx, user, manga)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return manga, err
	}
	return manga, nil
}

func FirebaseGetReadHistories(ctx context.Context, user models.UserFirebase) ([]contract.Manga, map[string]contract.Manga, error) {
	mangaHistories, mangaHistoriesMap, err := repository.FirebaseGetUserReadHistories(ctx, user)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return mangaHistories, mangaHistoriesMap, err
	}
	return mangaHistories, mangaHistoriesMap, nil
}
