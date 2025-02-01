package user_repository

import (
	"context"

	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
)

func GetByID(ctx context.Context, userID int64) (models.User, error) {
	user := models.User{}

	err := stmtGetByID.GetContext(ctx, &user, map[string]any{
		"id": userID,
	})
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"user_id": userID,
		}).Error(err)
		return user, err
	}

	return user, nil
}

func GetByIDs(ctx context.Context, userIDs pq.Int64Array) ([]models.User, error) {
	users := []models.User{}

	err := stmtGetByIDs.SelectContext(ctx, &users, map[string]any{
		"ids": userIDs,
	})
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return users, err
	}

	return users, nil
}

func GetByGuid(ctx context.Context, userGuid string) (models.User, error) {
	user := models.User{}

	err := stmtGetByGuid.GetContext(ctx, &user, map[string]any{
		"guid": userGuid,
	})
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"user_guid": userGuid,
		}).Error(err)
		return user, err
	}

	return user, nil
}

func GetByEmail(ctx context.Context, email string) (models.User, error) {
	user := models.User{}

	err := stmtGetByEmail.GetContext(ctx, &user, map[string]any{
		"email": email,
	})
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"email": email,
		}).Error(err)
		return user, err
	}

	return user, nil
}

func GetByVisitorID(ctx context.Context, visitorID string) (models.User, error) {
	user := models.User{}

	err := stmtGetByVisitorID.GetContext(ctx, &user, map[string]any{
		"visitor_id": visitorID,
	})
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"visitor_id": visitorID,
		}).Error(err)
		return user, err
	}

	return user, nil
}
