package anime_scrapper_service

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
)

func Watch(ctx context.Context, queryParams models.AnimeQueryParams) (models.EpisodeWatch, models.Meta, error) {
	episodeWatch := models.EpisodeWatch{}

	animeScrapper, err := animeScrapperGenerator(queryParams.Source)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return episodeWatch, models.Meta{}, err
	}

	episodeWatch, err = animeScrapper.Watch(ctx, queryParams)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return episodeWatch, models.Meta{}, err
	}

	return episodeWatch, models.Meta{}, nil
}
