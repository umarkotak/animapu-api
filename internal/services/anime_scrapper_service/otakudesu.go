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
	anime_scrapper_otakudesu "github.com/umarkotak/animapu-api/internal/repository/anime_scrapper_repository/otakudesu"
)

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
	for idx, oneTarget := range targets {
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
			otakudesuScrapper := anime_scrapper_otakudesu.NewOtakudesu()
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

		logrus.Infof("[%v/%v] Adding: %v", idx, len(targets), anime.Title)
		searchMap[id] = anime
	}

	result, _ := json.MarshalIndent(searchMap, " ", "  ")
	err = os.WriteFile("otakudesu_search_map", result, 0644)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return err
	}

	return nil
}
