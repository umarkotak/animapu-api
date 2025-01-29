package user_repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
)

func Insert(ctx context.Context, tx *sqlx.Tx, user models.User) (int64, error) {
	newID := int64(0)

	var row *sqlx.Row
	if tx == nil {
		row = stmtInsert.QueryRowContext(ctx, user)

	} else {
		namedStmt, err := tx.PrepareNamedContext(ctx, queryInsert)
		if err != nil {
			logrus.WithContext(ctx).Error(err)
			return newID, err
		}

		row = namedStmt.QueryRowContext(ctx, user)
	}

	err := row.Scan(&newID)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return newID, err
	}

	return newID, nil
}

func Update(ctx context.Context, tx *sqlx.Tx, user models.User) error {
	var err error
	var namedStmt *sqlx.NamedStmt

	if tx == nil {
		_, err = stmtUpdate.ExecContext(ctx, user)

	} else {
		namedStmt, err = tx.PrepareNamedContext(ctx, queryUpdate)
		if err != nil {
			logrus.WithContext(ctx).Error(err)
			return err
		}

		_, err = namedStmt.ExecContext(ctx, user)
	}
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return err
	}

	return nil
}
