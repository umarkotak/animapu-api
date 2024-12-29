package animeindo

import (
	"context"
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/utils/utils"
)

type Animeindo struct {
	AnimapuSource string
	Source        string
	Host          string
}

func NewAnimeindo() Animeindo {
	return Animeindo{
		AnimapuSource: models.ANIME_SOURCE_ANIMEINDO,
		Source:        "animeindo",
		Host:          "https://anime-indo.lol",
	}
}

func (s *Animeindo) GetLatest(ctx context.Context, queryParams models.AnimeQueryParams) ([]models.Anime, error) {
	animes := []models.Anime{}

	c := colly.NewCollector()

	c.OnHTML("#content-wrap > div.ngiri > div.menu > a", func(e *colly.HTMLElement) {
		coverUrl := e.ChildAttr("div > img", "data-original")

		episodeLink := e.Attr("href")
		animeID := strings.ReplaceAll(episodeLink, "/", "")
		splitted := strings.Split(animeID, "-episode-")
		if len(splitted) != 2 {
			return
		}
		animeID = splitted[0]

		animes = append(animes, models.Anime{
			ID:            animeID,
			Source:        s.Source,
			Title:         e.ChildText("div > p"),
			LatestEpisode: utils.ForceSanitizeStringToFloat(e.ChildText("div > span")),
			CoverUrls:     []string{coverUrl},
			OriginalLink:  fmt.Sprintf("%s/anime/%s", s.Host, animeID),
		})
	})

	targetUrl := fmt.Sprintf("%v/page/%v/", s.Host, queryParams.Page)
	err := c.Visit(targetUrl)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return animes, err
	}
	c.Wait()

	return animes, nil
}

func (s *Animeindo) GetSearch(ctx context.Context, queryParams models.AnimeQueryParams) ([]models.Anime, error) {
	animes := []models.Anime{}

	return animes, nil
}

func (s *Animeindo) GetDetail(ctx context.Context, queryParams models.AnimeQueryParams) (models.Anime, error) {
	targetUrl := fmt.Sprintf("%v/anime/%v/", s.Host, queryParams.SourceID)

	anime := models.Anime{
		ID:             queryParams.SourceID,
		Source:         s.AnimapuSource,
		Title:          "",
		LatestEpisode:  0,          // done
		CoverUrls:      []string{}, // done
		Episodes:       []models.Episode{},
		OriginalLink:   targetUrl,
		MultipleServer: true,
	}

	c := colly.NewCollector()

	c.OnHTML("div.detail > img", func(e *colly.HTMLElement) {
		anime.CoverUrls = append(anime.CoverUrls, fmt.Sprintf("%s/%s", s.Host, e.Attr("src")))
	})

	err := c.Visit(targetUrl)
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"target_url": targetUrl,
		}).Error(err)
		return anime, err
	}
	c.Wait()

	return anime, nil
}

func (s *Animeindo) Watch(ctx context.Context, queryParams models.AnimeQueryParams) (models.EpisodeWatch, error) {
	episode := models.EpisodeWatch{}

	return episode, nil
}

func (s *Animeindo) GetPerSeason(ctx context.Context, queryParams models.AnimeQueryParams) (models.AnimePerSeason, error) {
	animePerSeason := models.AnimePerSeason{}

	return animePerSeason, nil
}

func (s *Animeindo) GetRandom(ctx context.Context, queryParams models.AnimeQueryParams) ([]models.Anime, error) {
	animes := []models.Anime{}

	return animes, nil
}
