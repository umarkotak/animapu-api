package anime_scrapper_service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/contract"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository/anime_history_repository"
	"github.com/umarkotak/animapu-api/internal/repository/anime_repository"
)

func AnimeSync(ctx context.Context, animes []contract.Anime) error {
	for _, anime := range animes {
		existingAnime, err := anime_repository.GetBySourceAndSourceID(ctx, anime.Source, anime.ID)
		if err != nil && err != sql.ErrNoRows {
			logrus.WithContext(ctx).Error(err)
			continue
		}

		if existingAnime.ID != 0 {
			if existingAnime.LatestEpisode < existingAnime.LatestEpisode && existingAnime.ImageURLsEqual(existingAnime.CoverUrls) {
				continue
			}

			existingAnime.CoverUrls = anime.CoverUrls
			existingAnime.LatestEpisode = anime.LatestEpisode
			err = anime_repository.Update(ctx, nil, existingAnime)

		} else {
			newAnime := models.Anime{
				Source:        anime.Source,
				SourceID:      anime.ID,
				Title:         anime.Title,
				CoverUrls:     anime.CoverUrls,
				LatestEpisode: anime.LatestEpisode,
			}
			newAnime.ID, err = anime_repository.Insert(ctx, nil, newAnime)
		}
		if err != nil {
			logrus.WithContext(ctx).Error(err)
			continue
		}
	}

	return nil
}

func AnimeEpisodeSync(ctx context.Context, queryParams models.AnimeQueryParams, episode contract.Episode) error {
	anime, _, err := GetDetail(ctx, models.AnimeQueryParams{Source: episode.Source, SourceID: episode.AnimeID})
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return err
	}

	existingAnime, err := anime_repository.GetBySourceAndSourceID(ctx, episode.Source, episode.AnimeID)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return err
	}

	episodeNumber := float64(0)
	for _, animeEpisode := range anime.Episodes {
		if animeEpisode.ID == episode.ID {
			episodeNumber = animeEpisode.Number
			break
		}
	}

	animeHistory := models.AnimeHistory{
		UserID:          queryParams.User.ID,
		AnimeID:         existingAnime.ID,
		EpisodeNumber:   episodeNumber,
		SourceEpisodeID: episode.ID,
		FrontendPath:    fmt.Sprintf("/anime/%s/detail/%s/watch/%s", episode.Source, episode.AnimeID, episode.ID),
	}
	_, err = anime_history_repository.Insert(ctx, nil, animeHistory)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return err
	}

	return nil
}
