package datastore

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func Transaction(ctx context.Context, db *sqlx.DB, p func(tx *sqlx.Tx) error) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		if tx != nil {
			tx.Commit()
		}
		logrus.WithContext(ctx).Error(err)
		return err
	}
	defer func() {
		err := tx.Rollback()
		if err != nil && err != sql.ErrTxDone {
			logrus.WithContext(ctx).Error(err)
		}
	}()

	err = p(tx)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return err
	}

	return nil
}
