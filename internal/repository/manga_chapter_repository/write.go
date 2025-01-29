package manga_chapter_repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
)

func Insert(ctx context.Context, tx *sqlx.Tx, mangaChapter models.MangaChapter) (int64, error) {
	newID := int64(0)

	var row *sqlx.Row
	if tx == nil {
		row = stmtInsert.QueryRowContext(ctx, mangaChapter)

	} else {
		namedStmt, err := tx.PrepareNamedContext(ctx, queryInsert)
		if err != nil {
			logrus.WithContext(ctx).Error(err)
			return newID, err
		}

		row = namedStmt.QueryRowContext(ctx, mangaChapter)
	}

	err := row.Scan(&newID)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return newID, err
	}

	return newID, nil
}

func Update(ctx context.Context, tx *sqlx.Tx, mangaChapter models.MangaChapter) error {
	var err error
	var namedStmt *sqlx.NamedStmt

	if tx == nil {
		_, err = stmtUpdate.ExecContext(ctx, mangaChapter)

	} else {
		namedStmt, err = tx.PrepareNamedContext(ctx, queryUpdate)
		if err != nil {
			logrus.WithContext(ctx).Error(err)
			return err
		}

		_, err = namedStmt.ExecContext(ctx, mangaChapter)
	}
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return err
	}

	return nil
}

func UpdateByMangaIDAndSourceChapterID(ctx context.Context, tx *sqlx.Tx, mangaChapter models.MangaChapter) error {
	var err error
	var namedStmt *sqlx.NamedStmt

	if tx == nil {
		_, err = stmtUpdateByMangaIDAndSourceChapterID.ExecContext(ctx, mangaChapter)

	} else {
		namedStmt, err = tx.PrepareNamedContext(ctx, queryUpdateByMangaIDAndSourceChapterID)
		if err != nil {
			logrus.WithContext(ctx).Error(err)
			return err
		}

		_, err = namedStmt.ExecContext(ctx, mangaChapter)
	}
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return err
	}

	return nil
}
