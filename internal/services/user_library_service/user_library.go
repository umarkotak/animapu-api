package user_library_service

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository"
)

func Sync(ctx context.Context, user models.User, mangas []models.Manga) error {
	err := repository.SyncUserLibrary(ctx, user, mangas)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return err
	}
	return nil
}
