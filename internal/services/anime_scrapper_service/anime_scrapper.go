package anime_scrapper_service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository"
)

func GetLatest(ctx context.Context, queryParams models.AnimeQueryParams) ([]models.Anime, models.Meta, error) {
	animes := []models.Anime{}

	cachedAnimes, found := repository.GoCache().Get(queryParams.ToKey("GetLatest"))
	if found {
		cachedAnimesByte, err := json.Marshal(cachedAnimes)
		if err == nil {
			err = json.Unmarshal(cachedAnimesByte, &animes)
			if err == nil {
				return animes, models.Meta{FromCache: true}, nil
			}
		}
	}

	animeScrapper, err := animeScrapperGenerator(queryParams.Source)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return animes, models.Meta{}, err
	}

	animes, err = animeScrapper.GetLatest(ctx, queryParams)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return animes, models.Meta{}, err
	}

	if len(animes) > 0 {
		repository.GoCache().Set(queryParams.ToKey("GetLatest"), animes, 24*time.Hour)
	}

	return animes, models.Meta{}, nil
}

func GetPerSeason(ctx context.Context, queryParams models.AnimeQueryParams) (models.AnimePerSeason, models.Meta, error) {
	animePerSeason := models.AnimePerSeason{}

	animeScrapper, err := animeScrapperGenerator(queryParams.Source)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return animePerSeason, models.Meta{}, err
	}

	animePerSeason, err = animeScrapper.GetPerSeason(ctx, queryParams)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return animePerSeason, models.Meta{}, err
	}

	return animePerSeason, models.Meta{}, nil
}

func GetDetail(ctx context.Context, queryParams models.AnimeQueryParams) (models.Anime, models.Meta, error) {
	anime := models.Anime{}

	cachedAnime, found := repository.GoCache().Get(queryParams.ToKey("GetDetail"))
	if found {
		cachedAnimeByte, err := json.Marshal(cachedAnime)
		if err == nil {
			err = json.Unmarshal(cachedAnimeByte, &anime)
			if err == nil {
				return anime, models.Meta{FromCache: true}, nil
			}
		}
	}

	animeScrapper, err := animeScrapperGenerator(queryParams.Source)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return anime, models.Meta{}, err
	}

	anime, err = animeScrapper.GetDetail(ctx, queryParams)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return anime, models.Meta{}, err
	}

	if anime.ID != "" {
		repository.GoCache().Set(queryParams.ToKey("GetDetail"), anime, 24*time.Hour)
	}

	return anime, models.Meta{}, nil
}

func Watch(ctx context.Context, queryParams models.AnimeQueryParams) (models.EpisodeWatch, models.Meta, error) {
	episodeWatch := models.EpisodeWatch{}

	cachedEpisodeWatch, found := repository.GoCache().Get(queryParams.ToKey("Watch"))
	if found {
		cachedEpisodeWatchByte, err := json.Marshal(cachedEpisodeWatch)
		if err == nil {
			err = json.Unmarshal(cachedEpisodeWatchByte, &episodeWatch)
			if err == nil {
				return episodeWatch, models.Meta{FromCache: true}, nil
			}
		}
	}

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

	if episodeWatch.RawStreamUrl != "" || episodeWatch.IframeUrl != "" {
		repository.GoCache().Set(queryParams.ToKey("Watch"), episodeWatch, 24*time.Hour)
	}

	return episodeWatch, models.Meta{}, nil
}
