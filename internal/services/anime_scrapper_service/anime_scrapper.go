package anime_scrapper_service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository"
	"github.com/umarkotak/animapu-api/internal/repository/mal_api"
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

	for i := 1; i <= 3; i++ {
		animes, err = animeScrapper.GetLatest(ctx, queryParams)
		if err == nil {
			break
		} else {
			time.Sleep(2 * time.Second)
		}
	}
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
	animePerSeason := models.AnimePerSeason{
		ReleaseYear: queryParams.ReleaseYear,
		SeasonName:  queryParams.ReleaseSeason,
		SeasonIndex: models.SEASON_TO_SEASON_INDEX[queryParams.ReleaseSeason],
	}

	malAnimes, err := mal_api.GetSeasonalAnime(ctx, int(queryParams.ReleaseYear), queryParams.ReleaseSeason)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return animePerSeason, models.Meta{}, err
	}

	animes := []models.Anime{}
	for _, malAnime := range malAnimes {
		altTitles := []string{}

		if malAnime.AlternativeTitles.En != "" {
			altTitles = append(altTitles, malAnime.AlternativeTitles.En)
		}

		if malAnime.AlternativeTitles.Ja != "" {
			altTitles = append(altTitles, malAnime.AlternativeTitles.Ja)
		}

		if malAnime.AlternativeTitles.Synonyms != nil {
			altTitles = append(altTitles, malAnime.AlternativeTitles.Synonyms...)
		}

		anime := models.Anime{
			ID:                 fmt.Sprint(malAnime.ID),
			Source:             "mal",
			Title:              malAnime.Title,
			AltTitles:          altTitles,
			Description:        malAnime.Synopsis,
			LatestEpisode:      float64(malAnime.NumEpisodes),
			CoverUrls:          []string{malAnime.MainPicture.Medium},
			Genres:             []string{},
			Episodes:           []models.Episode{},
			OriginalLink:       "",
			ReleaseMonth:       "",
			ReleaseSeason:      queryParams.ReleaseSeason,
			ReleaseSeasonIndex: models.SEASON_TO_SEASON_INDEX[queryParams.ReleaseSeason],
			ReleaseYear:        queryParams.ReleaseYear,
			ReleaseDate:        "",
			Score:              float64(malAnime.MyListStatus.Score),
			Relations:          []models.Anime{},
			Relationship:       "",
			MultipleServer:     false,
			SearchTitle:        "",
		}
		animes = append(animes, anime)
	}
	animePerSeason.Animes = animes

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

	for i := 1; i <= 3; i++ {
		anime, err = animeScrapper.GetDetail(ctx, queryParams)
		if err == nil {
			break
		} else {
			time.Sleep(2 * time.Second)
		}
	}
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

	for i := 1; i <= 3; i++ {
		episodeWatch, err = animeScrapper.Watch(ctx, queryParams)
		if err == nil {
			break
		} else {
			time.Sleep(2 * time.Second)
		}
	}
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return episodeWatch, models.Meta{}, err
	}

	// if episodeWatch.RawStreamUrl != "" || episodeWatch.IframeUrl != "" {
	// 	repository.GoCache().Set(queryParams.ToKey("Watch"), episodeWatch, 24*time.Hour)
	// }

	return episodeWatch, models.Meta{}, nil
}

func GetSearch(ctx context.Context, queryParams models.AnimeQueryParams) ([]models.Anime, models.Meta, error) {
	animes := []models.Anime{}

	cachedAnimes, found := repository.GoCache().Get(queryParams.ToKey("GetSearch"))
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

	for i := 1; i <= 3; i++ {
		animes, err = animeScrapper.GetSearch(ctx, queryParams)
		if err == nil {
			break
		} else {
			time.Sleep(2 * time.Second)
		}
	}
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return animes, models.Meta{}, err
	}

	if len(animes) > 0 {
		repository.GoCache().Set(queryParams.ToKey("GetSearch"), animes, 24*time.Hour)
	}

	return animes, models.Meta{}, nil
}

func GetRandom(ctx context.Context, queryParams models.AnimeQueryParams) ([]models.Anime, models.Meta, error) {
	animes := []models.Anime{}

	animeScrapper, err := animeScrapperGenerator(queryParams.Source)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return animes, models.Meta{}, err
	}

	for i := 1; i <= 3; i++ {
		animes, err = animeScrapper.GetRandom(ctx, queryParams)
		if err == nil {
			break
		} else {
			time.Sleep(2 * time.Second)
		}
	}
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return animes, models.Meta{}, err
	}

	return animes, models.Meta{}, nil
}
