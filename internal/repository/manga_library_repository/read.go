package manga_library_repository

import (
	"context"
	"fmt"

	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/datastore"
	"github.com/umarkotak/animapu-api/internal/contract"
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

func GetByUserIDDetailed(ctx context.Context, params contract.MangaLibraryParams) ([]models.MangaHistoryDetailed, error) {
	mangaHistories := []models.MangaHistoryDetailed{}

	sort := "m.updated_at DESC"
	switch params.Sort {
	case "latest_update":
		sort = "m.updated_at DESC"
	case "recent_added":
		sort = "ml.id DESC"
	}

	queryGetByUserIDDetailed := fmt.Sprintf(`
		SELECT
			%s,
			m.source AS manga_source,
			m.source_id AS manga_source_id,
			m.title AS manga_title,
			m.cover_urls AS manga_cover_urls,
			m.latest_chapter AS manga_latest_chapter,
			m.updated_at AS manga_updated_at
		FROM manga_libraries ml
		INNER JOIN mangas m ON m.id = ml.manga_id
		WHERE
			ml.user_id = :user_id
			AND ml.deleted_at IS NULL
		ORDER BY %s
	`, allColumns, sort)

	stmtGetByUserIDDetailed, err := datastore.Get().Db.PrepareNamed(queryGetByUserIDDetailed)
	if err != nil {
		logrus.Fatal(err)
	}

	err = stmtGetByUserIDDetailed.SelectContext(ctx, &mangaHistories, map[string]any{
		"user_id": params.UserID,
	})
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"user_id": params.UserID,
		}).Error(err)
		return mangaHistories, err
	}

	return mangaHistories, nil
}

func GetByUserAndSourceDetail(ctx context.Context, userID int64, sources, sourceIDs pq.StringArray) ([]string, error) {
	mangaSourceIDs := []string{}

	err := stmtGetByUserAndSourceDetail.SelectContext(ctx, &mangaSourceIDs, map[string]any{
		"user_id":    userID,
		"sources":    sources,
		"source_ids": sourceIDs,
	})
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"user_id": userID,
		}).Error(err)
		return mangaSourceIDs, err
	}

	return mangaSourceIDs, nil
}
