package manga_history_repository

import (
	"context"

	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
)

func GetByID(ctx context.Context, historyID int64) (models.MangaHistory, error) {
	mangaHistory := models.MangaHistory{}

	err := stmtGetByID.GetContext(ctx, &mangaHistory, map[string]any{
		"id": historyID,
	})
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"history_id": historyID,
		}).Error(err)
		return mangaHistory, err
	}

	return mangaHistory, nil
}

func GetByMangaIDAndUserID(ctx context.Context, userID, mangaID int64) (models.MangaHistory, error) {
	mangaHistory := models.MangaHistory{}

	err := stmtGetByMangaIDAndUserID.GetContext(ctx, &mangaHistory, map[string]any{
		"user_id":  userID,
		"manga_id": mangaID,
	})
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"user_id":  userID,
			"manga_id": mangaID,
		}).Error(err)
		return mangaHistory, err
	}

	return mangaHistory, nil
}

func GetByUserID(ctx context.Context, userID int64) ([]models.MangaHistory, error) {
	mangaHistories := []models.MangaHistory{}

	err := stmtGetByUserID.SelectContext(ctx, &mangaHistories, map[string]any{
		"user_id": userID,
	})
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"user_id": userID,
		}).Error(err)
		return mangaHistories, err
	}

	return mangaHistories, nil
}

func GetByUserIDDetailed(ctx context.Context, userID int64, pagination models.Pagination) ([]models.MangaHistoryDetailed, error) {
	mangaHistories := []models.MangaHistoryDetailed{}

	err := stmtGetByUserIDDetailed.SelectContext(ctx, &mangaHistories, map[string]any{
		"user_id": userID,
		"limit":   pagination.Limit,
		"offset":  pagination.Offset,
	})
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"user_id": userID,
		}).Error(err)
		return mangaHistories, err
	}

	return mangaHistories, nil
}

func GetByUserAndSourceDetail(ctx context.Context, userID int64, sources, sourceIDs pq.StringArray) ([]models.MangaHistoryDetailed, error) {
	mangaHistories := []models.MangaHistoryDetailed{}

	err := stmtGetByUserAndSourceDetail.SelectContext(ctx, &mangaHistories, map[string]any{
		"user_id":    userID,
		"sources":    sources,
		"source_ids": sourceIDs,
	})
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"user_id": userID,
		}).Error(err)
		return mangaHistories, err
	}

	return mangaHistories, nil
}
