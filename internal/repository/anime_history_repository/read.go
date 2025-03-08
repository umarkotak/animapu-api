package anime_history_repository

import (
	"context"

	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
)

func GetByID(ctx context.Context, historyID int64) (models.AnimeHistory, error) {
	obj := models.AnimeHistory{}

	err := stmtGetByID.GetContext(ctx, &obj, map[string]any{
		"id": historyID,
	})
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"history_id": historyID,
		}).Error(err)
		return obj, err
	}

	return obj, nil
}

func GetByAnimeIDAndUserID(ctx context.Context, userID, animeID int64) (models.AnimeHistory, error) {
	obj := models.AnimeHistory{}

	err := stmtGetByAnimeIDAndUserID.GetContext(ctx, &obj, map[string]any{
		"user_id":  userID,
		"anime_id": animeID,
	})
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"user_id":  userID,
			"anime_id": animeID,
		}).Error(err)
		return obj, err
	}

	return obj, nil
}

func GetByUserID(ctx context.Context, userID int64) ([]models.AnimeHistory, error) {
	objs := []models.AnimeHistory{}

	err := stmtGetByUserID.SelectContext(ctx, &objs, map[string]any{
		"user_id": userID,
	})
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"user_id": userID,
		}).Error(err)
		return objs, err
	}

	return objs, nil
}

func GetByUserIDDetailed(ctx context.Context, userID int64, pagination models.Pagination) ([]models.AnimeHistoryDetailed, error) {
	objs := []models.AnimeHistoryDetailed{}

	err := stmtGetByUserIDDetailed.SelectContext(ctx, &objs, map[string]any{
		"user_id": userID,
		"limit":   pagination.Limit,
		"offset":  pagination.Offset,
	})
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"user_id": userID,
		}).Error(err)
		return objs, err
	}

	return objs, nil
}

func GetByUserAndSourceDetail(ctx context.Context, userID int64, sources, sourceIDs pq.StringArray) ([]models.AnimeHistoryDetailed, error) {
	objs := []models.AnimeHistoryDetailed{}

	err := stmtGetByUserAndSourceDetail.SelectContext(ctx, &objs, map[string]any{
		"user_id":    userID,
		"sources":    sources,
		"source_ids": sourceIDs,
	})
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"user_id": userID,
		}).Error(err)
		return objs, err
	}

	return objs, nil
}

func GetRecentHistories(ctx context.Context, pagination models.Pagination) ([]models.AnimeHistoryDetailed, error) {
	objs := []models.AnimeHistoryDetailed{}

	err := stmtGetRecentHistories.SelectContext(ctx, &objs, map[string]any{
		"limit":  pagination.Limit,
		"offset": pagination.Offset,
	})
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return objs, err
	}

	return objs, nil
}
