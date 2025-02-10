package anime_scrapper_service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/datastore"
	"github.com/umarkotak/animapu-api/internal/contract"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository/anichart_api.go"
)

func GetLatest(ctx context.Context, queryParams models.AnimeQueryParams) ([]contract.Anime, models.Meta, error) {
	animes := []contract.Anime{}

	cachedAnimes, found := datastore.Get().GoCache.Get(queryParams.ToKey("GetLatest"))
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
		datastore.Get().GoCache.Set(queryParams.ToKey("GetLatest"), animes, 24*time.Hour)
	}

	return animes, models.Meta{}, nil
}

func GetPerSeason(ctx context.Context, queryParams models.AnimeQueryParams) (contract.AnimePerSeason, models.Meta, error) {
	animePerSeason := contract.AnimePerSeason{
		ReleaseYear: queryParams.ReleaseYear,
		SeasonName:  queryParams.ReleaseSeason,
		SeasonIndex: models.SEASON_TO_SEASON_INDEX[queryParams.ReleaseSeason],
	}

	cachedAnimePerSeason, found := datastore.Get().GoCache.Get(queryParams.ToKey("GetPerSeason"))
	if found {
		cachedAnimesByte, err := json.Marshal(cachedAnimePerSeason)
		if err == nil {
			err = json.Unmarshal(cachedAnimesByte, &animePerSeason)
			if err == nil {
				return animePerSeason, models.Meta{FromCache: true}, nil
			}
		}
	}

	anichartAnimes, err := anichart_api.GetSeasonalAnime(ctx, int(queryParams.ReleaseYear), queryParams.ReleaseSeason)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return animePerSeason, models.Meta{}, err
	}

	animes := []contract.Anime{}
	for _, anichartAnime := range anichartAnimes {
		altTitles := []string{}

		if anichartAnime.Title.English != "" {
			altTitles = append(altTitles, anichartAnime.Title.English)
		}

		if anichartAnime.Title.Native != "" {
			altTitles = append(altTitles, anichartAnime.Title.Native)
		}

		if anichartAnime.Title.Romaji != "" {
			altTitles = append(altTitles, anichartAnime.Title.Romaji)
		}

		anime := contract.Anime{
			ID:                 fmt.Sprint(anichartAnime.ID),
			Source:             "mal",
			Title:              anichartAnime.Title.Romaji,
			AltTitles:          altTitles,
			Description:        anichartAnime.Description,
			LatestEpisode:      float64(anichartAnime.Episodes),
			CoverUrls:          []string{anichartAnime.CoverImage.Large},
			Genres:             []string{},
			Episodes:           []contract.Episode{},
			OriginalLink:       "",
			ReleaseMonth:       "",
			ReleaseSeason:      queryParams.ReleaseSeason,
			ReleaseSeasonIndex: models.SEASON_TO_SEASON_INDEX[queryParams.ReleaseSeason],
			ReleaseYear:        queryParams.ReleaseYear,
			ReleaseDate:        "",
			Score:              float64(anichartAnime.AverageScore),
			Relations:          []contract.Anime{},
			Relationship:       "",
			MultipleServer:     false,
			SearchTitle:        "",
		}
		animes = append(animes, anime)
	}
	animePerSeason.Animes = animes

	if len(animePerSeason.Animes) > 0 {
		datastore.Get().GoCache.Set(queryParams.ToKey("GetPerSeason"), animes, 7*24*time.Hour)
	}

	return animePerSeason, models.Meta{}, nil
}

func GetDetail(ctx context.Context, queryParams models.AnimeQueryParams) (contract.Anime, models.Meta, error) {
	anime := contract.Anime{}

	cachedAnime, found := datastore.Get().GoCache.Get(queryParams.ToKey("GetDetail"))
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
		datastore.Get().GoCache.Set(queryParams.ToKey("GetDetail"), anime, 24*time.Hour)
	}

	return anime, models.Meta{}, nil
}

func Watch(ctx context.Context, queryParams models.AnimeQueryParams) (contract.EpisodeWatch, models.Meta, error) {
	episodeWatch := contract.EpisodeWatch{}

	cachedEpisodeWatch, found := datastore.Get().GoCache.Get(queryParams.ToKey("Watch"))
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

	iter := 1
	for i := 1; i <= iter; i++ {
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
	// 	datastore.Get().GoCache.Set(queryParams.ToKey("Watch"), episodeWatch, 24*time.Hour)
	// }

	return episodeWatch, models.Meta{}, nil
}

func GetSearch(ctx context.Context, queryParams models.AnimeQueryParams) ([]contract.Anime, models.Meta, error) {
	animes := []contract.Anime{}

	cachedAnimes, found := datastore.Get().GoCache.Get(queryParams.ToKey("GetSearch"))
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
		datastore.Get().GoCache.Set(queryParams.ToKey("GetSearch"), animes, 24*time.Hour)
	}

	return animes, models.Meta{}, nil
}

func GetRandom(ctx context.Context, queryParams models.AnimeQueryParams) ([]contract.Anime, models.Meta, error) {
	animes := []contract.Anime{}

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
