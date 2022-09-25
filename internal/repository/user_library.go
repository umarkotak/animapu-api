package repository

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
)

func SyncUserLibrary(ctx context.Context, user models.User, mangas []models.Manga) error {
	var err error

	oneUser := usersRef.Child(user.Uid)

	if oneUser == nil {
		user.Libraries = mangas

		err = oneUser.Set(ctx, user)
		if err != nil {
			logrus.WithContext(ctx).Error(err)
			return err
		}

		return nil
	}

	return nil
}
