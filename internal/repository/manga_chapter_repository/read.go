package manga_chapter_repository

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
)

func GetByID(ctx context.Context, chapterID int64) (models.MangaChapter, error) {
	mangaChapter := models.MangaChapter{}

	err := stmtGetByID.GetContext(ctx, &mangaChapter, map[string]any{
		"id": chapterID,
	})
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"chapter_id": chapterID,
		}).Error(err)
		return mangaChapter, err
	}

	return mangaChapter, nil
}

func GetByMangaIDAndSourceChapterID(ctx context.Context, mangaID int64, sourceChapterID string) (models.MangaChapter, error) {
	mangaChapter := models.MangaChapter{}

	err := stmtGetByMangaIDAndSourceChapterID.GetContext(ctx, &mangaChapter, map[string]any{
		"manga_id":          mangaID,
		"source_chapter_id": sourceChapterID,
	})
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"manga_id":          mangaID,
			"source_chapter_id": sourceChapterID,
		}).Error(err)
		return mangaChapter, err
	}

	return mangaChapter, nil
}
