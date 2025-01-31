package affiliate_link_repository

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
)

func GetByID(ctx context.Context, id int64) (models.AffiliateLink, error) {
	obj := models.AffiliateLink{}

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

func GetByShortLink(ctx context.Context, shortLink string) (models.AffiliateLink, error) {
	obj := models.AffiliateLink{}

	err := stmtGetByShortLink.GetContext(ctx, &obj, map[string]any{
		"short_link": shortLink,
	})
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"short_link": shortLink,
		}).Error(err)
		return obj, err
	}

	return obj, nil
}

func GetRandom(ctx context.Context, pagination models.Pagination) ([]models.AffiliateLink, error) {
	objs := []models.AffiliateLink{}

	err := stmtGetRandom.SelectContext(ctx, &objs, pagination)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return objs, err
	}

	return objs, nil
}

func GetList(ctx context.Context, pagination models.Pagination) ([]models.AffiliateLink, error) {
	objs := []models.AffiliateLink{}

	err := stmtGetList.SelectContext(ctx, &objs, pagination)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return objs, err
	}

	return objs, nil
}
