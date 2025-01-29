package manga_repository

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
)

func GetByID(ctx context.Context, mangaID int64) (models.Manga, error) {
	manga := models.Manga{}

	err := stmtGetByID.GetContext(ctx, &manga, map[string]any{
		"id": mangaID,
	})
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"manga_id": mangaID,
		}).Error(err)
		return manga, err
	}

	return manga, nil
}

func GetBySourceAndSourceID(ctx context.Context, source, sourceID string) (models.Manga, error) {
	manga := models.Manga{}

	err := stmtGetBySourceAndSourceID.GetContext(ctx, &manga, map[string]any{
		"source":    source,
		"source_id": sourceID,
	})
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"source":    source,
			"source_id": sourceID,
		}).Error(err)
		return manga, err
	}

	return manga, nil
}
