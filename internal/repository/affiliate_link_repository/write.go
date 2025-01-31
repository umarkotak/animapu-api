package affiliate_link_repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
)

func Insert(ctx context.Context, tx *sqlx.Tx, obj models.AffiliateLink) (int64, error) {
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
