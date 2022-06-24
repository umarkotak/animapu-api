package repository

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
)

func FbMangahubGetHome(ctx context.Context) (models.FbMangahubHome, error) {
	homeRef := mangahubRef.Child("home")

	fbMangahubHome := models.FbMangahubHome{}
	err := homeRef.Get(ctx, &fbMangahubHome)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return fbMangahubHome, err
	}

	return fbMangahubHome, nil
}

func FbMangahubUpsertHome(ctx context.Context, mangas []models.Manga) ([]models.Manga, error) {
	homeRef := mangahubRef.Child("home")

	now := time.Now().UTC()
	err := homeRef.Set(ctx, models.FbMangahubHome{
		UpdatedAt:     now,
		UpdatedAtUnix: now.Unix(),
		ExpiredAt:     now.Add(4 * time.Hour), // 4 hours
		Mangas:        mangas,
	})
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return mangas, err
	}

	return mangas, nil
}

func FbMangahubGetDetail(ctx context.Context, manga models.Manga) (models.FbMangahubDetail, error) {
	var err error
	fbMangahubDetail := models.FbMangahubDetail{}
	detailRef := mangahubRef.Child("detail")

	if detailRef == nil {
		err = detailRef.Set(ctx, map[string]interface{}{})
		if err != nil {
			logrus.WithContext(ctx).Error(err)
			return fbMangahubDetail, err
		}
	}

	fbMangahubDetailSelected := detailRef.Child(manga.SourceID)
	if fbMangahubDetailSelected == nil {
		err = models.ErrCacheNotFound
		return fbMangahubDetail, err
	}

	err = fbMangahubDetailSelected.Get(ctx, &fbMangahubDetail)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return fbMangahubDetail, err
	}

	return fbMangahubDetail, nil
}

func FbMangahubUpsertDetail(ctx context.Context, manga models.Manga) (models.Manga, error) {
	var err error
	detailRef := mangahubRef.Child("detail")

	if detailRef == nil {
		err = detailRef.Set(ctx, map[string]interface{}{})
		if err != nil {
			logrus.WithContext(ctx).Error(err)
			return manga, err
		}
	}

	fbMangahubDetailSelected := detailRef.Child(manga.SourceID)

	now := time.Now().UTC()
	err = fbMangahubDetailSelected.Set(ctx, models.FbMangahubDetail{
		UpdatedAt:     now,
		UpdatedAtUnix: now.Unix(),
		ExpiredAt:     now.Add(24 * 3 * time.Hour), // 3 days
		Manga:         manga,
	})
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return manga, err
	}

	return manga, nil
}
