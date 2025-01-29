package gogo_anime

import (
	"context"
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/contract"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/utils/utils"
)

type GogoAnime struct {
	AnimapuSource string
	Source        string
	Host          string
}

func NewGogoAnime() GogoAnime {
	return GogoAnime{
		AnimapuSource: models.ANIME_SOURCE_GOGO_ANIME,
		Source:        "gogo_anime",
		Host:          "https://ww10.gogoanimes.org",
	}
}

func (s *GogoAnime) GetLatest(ctx context.Context, queryParams models.AnimeQueryParams) ([]contract.Anime, error) {
	animes := []contract.Anime{}

	c := colly.NewCollector()

	c.OnHTML("body > div.last_episodes.loaddub > ul", func(e *colly.HTMLElement) {
		e.ForEach("li", func(i int, h *colly.HTMLElement) {
			coverUrl := h.ChildAttr("div > a > img", "src")

			episodeLink := h.ChildAttr("div > a", "href")
			animeID := strings.ReplaceAll(episodeLink, "/watch/", "")
			animeID = strings.TrimSuffix(animeID, "/")
			splitted := strings.Split(animeID, "-episode-")
			if len(splitted) != 2 {
				return
			}
			animeID = splitted[0]

			animes = append(animes, contract.Anime{
				ID:            animeID,
				Source:        s.Source,
				Title:         h.ChildText("p.name > a"),
				LatestEpisode: utils.ForceSanitizeStringToFloat(h.ChildText("p.episode")),
				CoverUrls:     []string{coverUrl},
				OriginalLink:  fmt.Sprintf("%s/category/%s", s.Host, animeID),
			})

		})
	})

	targetUrl := fmt.Sprintf("%v/ajax/page-recent-release?page=%v&type=1", s.Host, queryParams.Page)
	err := c.Visit(targetUrl)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return animes, err
	}
	c.Wait()

	return animes, nil
}

func (s *GogoAnime) GetSearch(ctx context.Context, queryParams models.AnimeQueryParams) ([]contract.Anime, error) {
	animes := []contract.Anime{}

	return animes, nil
}

func (s *GogoAnime) GetDetail(ctx context.Context, queryParams models.AnimeQueryParams) (contract.Anime, error) {
	targetUrl := fmt.Sprintf("%v/category/%v", s.Host, queryParams.SourceID)

	anime := contract.Anime{
		ID:             queryParams.SourceID,
		Source:         s.AnimapuSource,
		Title:          "",
		LatestEpisode:  0,          // done
		CoverUrls:      []string{}, // done
		Episodes:       []contract.Episode{},
		OriginalLink:   targetUrl,
		MultipleServer: true,
	}

	c := colly.NewCollector()

	c.OnHTML("#wrapper_bg > section > section.content_left > div.main_body > div.anime_info_body > div.anime_info_body_bg > img", func(e *colly.HTMLElement) {
		anime.CoverUrls = append(anime.CoverUrls, e.Attr("src"))
	})

	c.OnHTML("#wrapper_bg > section > section.content_left > div.main_body > div.anime_info_body > div.anime_info_body_bg > h1", func(e *colly.HTMLElement) {
		anime.Title = e.Text
	})

	c.OnHTML("#wrapper_bg > section > section.content_left > div.main_body > div.anime_info_body > div.anime_info_body_bg > p:nth-child(5)", func(e *colly.HTMLElement) {
		anime.Description = e.Text
	})

	c.OnHTML("#episode_related > li", func(e *colly.HTMLElement) {
		episodeLink := e.ChildAttr("a", "href")
		episodeID := strings.ReplaceAll(episodeLink, "/watch/", "")

		episode := contract.Episode{
			AnimeID:      queryParams.SourceID,
			Source:       s.Source,
			ID:           episodeID,
			Number:       utils.ForceSanitizeStringToFloat(e.ChildText("a > div.name")),
			Title:        e.ChildText("a > div.name"),
			OriginalLink: fmt.Sprintf("%v%v", s.Host, episodeLink),
			// UseTitle:     true,
		}
		anime.Episodes = append(anime.Episodes, episode)
	})

	err := c.Visit(targetUrl)
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"target_url": targetUrl,
		}).Error(err)
		return anime, err
	}
	c.Wait()

	epListUrl := fmt.Sprintf("%s/ajaxajax/load-list-episode?ep_start=0&ep_end=&id=0&default_ep=&alias=/category/%s", s.Host, queryParams.SourceID)

	err = c.Visit(epListUrl)
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"target_url": targetUrl,
		}).Error(err)
		return anime, err
	}
	c.Wait()

	return anime, nil
}

func (s *GogoAnime) Watch(ctx context.Context, queryParams models.AnimeQueryParams) (contract.EpisodeWatch, error) {
	episode := contract.EpisodeWatch{
		StreamType:  "iframe",
		IframeUrl:   "",
		IframeUrls:  map[string]string{},
		OriginalUrl: "",
	}

	targetUrl := fmt.Sprintf("%v/watch/%v", s.Host, queryParams.EpisodeID)

	c := colly.NewCollector()

	c.OnHTML("#load_anime > div > div > iframe", func(e *colly.HTMLElement) {
		episode.IframeUrl = s.cleanUpUrl(e.Attr("src"))
	})

	err := c.Visit(targetUrl)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return episode, err
	}
	c.Wait()

	return episode, nil
}

func (s *GogoAnime) GetPerSeason(ctx context.Context, queryParams models.AnimeQueryParams) (contract.AnimePerSeason, error) {
	animePerSeason := contract.AnimePerSeason{}

	return animePerSeason, nil
}

func (s *GogoAnime) GetRandom(ctx context.Context, queryParams models.AnimeQueryParams) ([]contract.Anime, error) {
	animes := []contract.Anime{}

	return animes, nil
}

func (r *GogoAnime) cleanUpUrl(str string) string {
	return strings.ReplaceAll(str, "\u0026", "&")
}
