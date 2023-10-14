package anime_scrapper_service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/local_db"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository/anime_scrapper_repository"
)

func GetLatest(ctx context.Context, queryParams models.AnimeQueryParams) ([]models.Anime, models.Meta, error) {
	animes := []models.Anime{}

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

	return anime, models.Meta{}, nil
}

func Watch(ctx context.Context, queryParams models.AnimeQueryParams) (models.EpisodeWatch, models.Meta, error) {
	episodeWatch := models.EpisodeWatch{}

	animeScrapper, err := animeScrapperGenerator(queryParams.Source)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return episodeWatch, models.Meta{}, err
	}

	if queryParams.WatchVersion == "2" {
		episodeWatch, err = animeScrapper.WatchV2(ctx, queryParams)
	} else {
		episodeWatch, err = animeScrapper.Watch(ctx, queryParams)
	}
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return episodeWatch, models.Meta{}, err
	}

	return episodeWatch, models.Meta{}, nil
}

func ScrapOtakudesuAllAnimes(ctx context.Context) error {
	c := colly.NewCollector()

	targets := []string{}
	// Search page scrap
	c.OnHTML("a.hodebgst", func(e *colly.HTMLElement) {
		targets = append(targets, e.Attr("href"))
	})

	targetUrl := "https://otakudesu.wiki/anime-list"
	err := c.Visit(targetUrl)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return err
	}
	c.Wait()

	otakudesuDB := local_db.AnimeLinkToDetailMap

	searchMap := map[string]models.Anime{}
	for _, oneTarget := range targets {
		splitted := strings.Split(oneTarget, "/anime/")
		id := ""
		if len(splitted) > 0 {
			id = strings.ReplaceAll(splitted[len(splitted)-1], "/", "")
		}
		if id == "" {
			err = fmt.Errorf("id not found")
			logrus.WithContext(ctx).WithFields(logrus.Fields{
				"url": oneTarget,
				"id":  id,
			}).Error(err)
			continue
		}

		anime, found := otakudesuDB[id]
		if !found {
			otakudesuScrapper := anime_scrapper_repository.NewOtakudesu()
			anime, err = otakudesuScrapper.GetDetail(ctx, models.AnimeQueryParams{
				SourceID: id,
			})
			if err != nil {
				logrus.WithContext(ctx).Error(err)
				continue
			}
		}

		seasonObj, found := models.MONTH_TO_SEASON_MAP[strings.ToLower(anime.ReleaseMonth)]
		if !found {
			seasonObj = models.Season{1, "winter"}
		}
		anime.ReleaseSeason = seasonObj.Name
		anime.ReleaseSeasonIndex = seasonObj.Index
		searchMap[id] = anime

		result, _ := json.MarshalIndent(searchMap, " ", "  ")
		err = os.WriteFile("otakudesu_search_map", result, 0644)
		if err != nil {
			logrus.WithContext(ctx).Error(err)
			continue
		}
	}

	return nil
}
