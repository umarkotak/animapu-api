package anime_repository

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
)

func GetByID(ctx context.Context, id int64) (models.Anime, error) {
	obj := models.Anime{}

	err := stmtGetByID.GetContext(ctx, &obj, map[string]any{
		"id": id,
	})
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"id": id,
		}).Error(err)
		return obj, err
	}

	return obj, nil
}

func GetBySourceAndSourceID(ctx context.Context, source, sourceID string) (models.Anime, error) {
	obj := models.Anime{}

	err := stmtGetBySourceAndSourceID.GetContext(ctx, &obj, map[string]any{
		"source":    source,
		"source_id": sourceID,
	})
	if err != nil {
		if err != sql.ErrNoRows {
			logrus.WithContext(ctx).WithFields(logrus.Fields{
				"source":    source,
				"source_id": sourceID,
			}).Error(err)
		}
		return obj, err
	}

	return obj, nil
}

func GetBySourceAndSourceIDs(ctx context.Context, source, sourceIDs pq.StringArray) ([]models.Anime, error) {
	objs := []models.Anime{}

	err := stmtGetBySourceAndSourceID.SelectContext(ctx, &objs, map[string]any{
		"source":     source,
		"source_ids": sourceIDs,
	})
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"source":     source,
			"source_ids": sourceIDs,
		}).Error(err)
		return objs, err
	}

	return objs, nil
}
