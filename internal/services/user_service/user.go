package user_service

import (
	"context"
	"database/sql"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository/user_repository"
)

func UpsertAndGetUser(ctx context.Context, user models.User) (models.User, error) {
	var err error

	if user.Guid.String == "" && user.VisitorId == "" {
		return models.User{}, err
	}

	if user.Guid.String == "" {
		existingUser, err := user_repository.GetByVisitorID(ctx, user.VisitorId)
		if err != nil && err != sql.ErrNoRows {
			logrus.WithContext(ctx).Error(err)
			return models.User{}, err
		}

		// user with visitor id exists
		if existingUser.ID != 0 {
			return existingUser, nil
		}

		user.Guid = sql.NullString{"", false}
		user.Email = sql.NullString{"", false}
		user.ID, err = user_repository.Insert(ctx, nil, user)
		if err != nil {
			logrus.WithContext(ctx).Error(err)
			return models.User{}, err
		}

		return user, nil
	}

	existingUser, err := user_repository.GetByGuid(ctx, user.Guid.String)
	if err != nil && err != sql.ErrNoRows {
		logrus.WithContext(ctx).Error(err)
		return models.User{}, err
	}

	// user with guid exists
	if existingUser.ID != 0 {
		return existingUser, nil
	}

	user.VisitorId = user.Guid.String
	user.ID, err = user_repository.Insert(ctx, nil, user)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return models.User{}, err
	}

	return user, nil
}
