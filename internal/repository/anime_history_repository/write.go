package anime_history_repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
)

func Insert(ctx context.Context, tx *sqlx.Tx, obj models.AnimeHistory) (int64, error) {
	newID := int64(0)

	var row *sqlx.Row
	if tx == nil {
		row = stmtInsert.QueryRowContext(ctx, obj)

	} else {
		namedStmt, err := tx.PrepareNamedContext(ctx, queryInsert)
		if err != nil {
			logrus.WithContext(ctx).Error(err)
			return newID, err
		}

		row = namedStmt.QueryRowContext(ctx, obj)
	}

	err := row.Scan(&newID)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return newID, err
	}

	return newID, nil
}

func Update(ctx context.Context, tx *sqlx.Tx, obj models.AnimeHistory) error {
	var err error
	var namedStmt *sqlx.NamedStmt

	if tx == nil {
		_, err = stmtUpdate.ExecContext(ctx, obj)

	} else {
		namedStmt, err = tx.PrepareNamedContext(ctx, queryUpdate)
		if err != nil {
			logrus.WithContext(ctx).Error(err)
			return err
		}

		_, err = namedStmt.ExecContext(ctx, obj)
	}
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return err
	}

	return nil
}

func UpdateByAnimeIDAndUserID(ctx context.Context, tx *sqlx.Tx, obj models.AnimeHistory) error {
	var err error
	var namedStmt *sqlx.NamedStmt

	if tx == nil {
		_, err = stmtUpdateByAnimeIDAndUserID.ExecContext(ctx, obj)

	} else {
		namedStmt, err = tx.PrepareNamedContext(ctx, queryUpdateByAnimeIDAndUserID)
		if err != nil {
			logrus.WithContext(ctx).Error(err)
			return err
		}

		_, err = namedStmt.ExecContext(ctx, obj)
	}
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return err
	}

	return nil
}
