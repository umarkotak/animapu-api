package history_repository

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
)

func GetByUserID(ctx context.Context, userID int64, pagination models.Pagination) ([]models.History, error) {
	histories := []models.History{}
	if err := stmtGetByUserID.SelectContext(ctx, &histories, map[string]any{
		"user_id": userID,
		"limit":   pagination.Limit,
		"offset":  pagination.Offset,
	}); err != nil {
		logrus.WithContext(ctx).WithField("user_id", userID).Error(err)
		return histories, err
	}

	return histories, nil
}
