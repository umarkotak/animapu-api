package gogo_anime_new

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"slices"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/utils/utils"
)

type GogoAnimeNew struct {
	AnimapuSource string
	Source        string
	Host          string
}

func NewGogoAnimeNew() GogoAnimeNew {
	return GogoAnimeNew{
		AnimapuSource: models.ANIME_SOURCE_GOGO_ANIME,
		Source:        "gogo_anime_new",
		Host:          "https://gogoanime.by",
	}
}

func (s *GogoAnimeNew) GetLatest(ctx context.Context, queryParams models.AnimeQueryParams) ([]models.Anime, error) {
	animes := []models.Anime{}

	c := colly.NewCollector()

	c.OnHTML("#content > div > div.postbody > div.bixbox.bixboxarc.bbnofrm > div.mrgn > div.listupd > article", func(e *colly.HTMLElement) {
		coverUrl := e.ChildAttr("div > a > div.limit > img", "src")

		animeLink := e.ChildAttr("div > a", "href")
		animeID := strings.ReplaceAll(animeLink, s.Host, "")
		animeID = strings.ReplaceAll(animeID, "/series/", "")
		animeID = strings.TrimSuffix(animeID, "/")

		animes = append(animes, models.Anime{
			ID:            animeID,
			Source:        s.Source,
			Title:         e.ChildText("div > a > div.tt.tts > h2"),
			LatestEpisode: 0,
			CoverUrls:     []string{coverUrl},
			OriginalLink:  animeLink,
		})
	})

	targetUrl := fmt.Sprintf("%s/series/?page=%v&status=&type=&order=update", s.Host, queryParams.Page)
	err := c.Visit(targetUrl)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return animes, err
	}
	c.Wait()

	return animes, nil
}

func (s *GogoAnimeNew) GetSearch(ctx context.Context, queryParams models.AnimeQueryParams) ([]models.Anime, error) {
	animes := []models.Anime{}

	c := colly.NewCollector()

	c.OnHTML("#content > div > div.postbody > div > div.listupd > article", func(e *colly.HTMLElement) {
		coverUrl := e.ChildAttr("div > a > div.limit > img", "src")

		animeLink := e.ChildAttr("div > a", "href")
		animeID := strings.ReplaceAll(animeLink, s.Host, "")
		animeID = strings.ReplaceAll(animeID, "/series/", "")
		animeID = strings.TrimSuffix(animeID, "/")

		animes = append(animes, models.Anime{
			ID:            animeID,
			Source:        s.Source,
			Title:         e.ChildText("div > a > div.tt.tts > h2"),
			LatestEpisode: 0,
			CoverUrls:     []string{coverUrl},
			OriginalLink:  animeLink,
		})
	})

	targetUrl := fmt.Sprintf("%s/?s=%s", s.Host, url.QueryEscape(queryParams.Title))
	err := c.Visit(targetUrl)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return animes, err
	}
	c.Wait()

	return animes, nil
}

func (s *GogoAnimeNew) GetDetail(ctx context.Context, queryParams models.AnimeQueryParams) (models.Anime, error) {
	targetUrl := fmt.Sprintf("%v/series/%v/", s.Host, queryParams.SourceID)

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

	c.OnHTML("article > div.bixbox.animefull > div > div.thumbook > div.thumb > img", func(e *colly.HTMLElement) {
		anime.CoverUrls = append(anime.CoverUrls, e.Attr("src"))
	})

	c.OnHTML("article > div.bixbox.animefull > div > div.infox > h1", func(e *colly.HTMLElement) {
		anime.Title = e.Text
	})

	c.OnHTML("article > div.bixbox.animefull > div > div.infox > div > span", func(e *colly.HTMLElement) {
		anime.Description = e.Text
	})

	c.OnHTML("article > div.episodes-container > div.episode-item", func(e *colly.HTMLElement) {
		episodeLink := e.ChildAttr("a", "href")
		episodeID := strings.ReplaceAll(episodeLink, s.Host, "")
		episodeID = strings.Trim(episodeID, "/")

		episode := models.Episode{
			AnimeID:      queryParams.SourceID,
			Source:       s.Source,
			ID:           episodeID,
			Number:       utils.ForceSanitizeStringToFloat(e.ChildText("a")),
			Title:        e.ChildText("a"),
			OriginalLink: fmt.Sprintf("%v%v", s.Host, episodeLink),
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

	slices.Reverse(anime.Episodes)

	return anime, nil
}

func (s *GogoAnimeNew) Watch(ctx context.Context, queryParams models.AnimeQueryParams) (models.EpisodeWatch, error) {
	episode := models.EpisodeWatch{
		StreamType:    "iframe",
		IframeUrl:     "",
		IframeUrls:    map[string]string{},
		OriginalUrl:   "",
		StreamOptions: []models.StreamOption{},
	}

	serverParams := []GogoAnimeServerParams{}

	targetUrl := fmt.Sprintf("%v/%v", s.Host, queryParams.EpisodeID)

	c := colly.NewCollector()

	featureImage := ""
	c.OnHTML("article > div.megavid > div > div.item.meta > div.tb > img", func(e *colly.HTMLElement) {
		featureImage = e.Attr("src")
	})

	c.OnHTML("#w-servers > div.servers > div > ul", func(e *colly.HTMLElement) {
		e.ForEach("li", func(i int, h *colly.HTMLElement) {
			gogoAnimeServerParams := GogoAnimeServerParams{
				DataType:          h.Attr("data-type"),
				DataEncryptedUrl1: h.Attr("data-encrypted-url1"),
				DataEncryptedUrl2: h.Attr("data-encrypted-url2"),
				DataEncryptedUrl3: h.Attr("data-encrypted-url3"),
			}
			serverParams = append(serverParams, gogoAnimeServerParams)

			episode.StreamOptions = append(episode.StreamOptions, models.StreamOption{
				Resolution: "",
				Index:      gogoAnimeServerParams.DataType,
				Name:       gogoAnimeServerParams.DataType,
				Used:       false,
			})
		})
	})

	err := c.Visit(targetUrl)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return episode, err
	}
	c.Wait()

	if len(serverParams) == 0 {
		err = fmt.Errorf("server unavailable")
		logrus.WithContext(ctx).Error(err)
		return episode, err
	}

	usedGogoAnimeServerParams := serverParams[len(serverParams)-1]

	if queryParams.StreamIdx != "" {
		for _, oneServer := range serverParams {
			if oneServer.DataType == queryParams.StreamIdx {
				usedGogoAnimeServerParams = oneServer
				break
			}
		}
	}

	usedGogoAnimeServerParams.FeatureImage = featureImage
	usedGogoAnimeServerParams.EpisodeUrl = targetUrl

	gogoAnimeStream, err := s.getStreamUrl(ctx, usedGogoAnimeServerParams)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return episode, err
	}

	if gogoAnimeStream.Mode == "iframe" {
		episode.StreamType = "iframe"
		episode.IframeUrl = gogoAnimeStream.Src
	} else if gogoAnimeStream.Mode == "mp4" {
		episode.StreamType = "mp4"
		episode.RawStreamUrl = s.cleanUpUrl(gogoAnimeStream.Src)
	}

	for idx, _ := range episode.StreamOptions {
		if episode.StreamOptions[idx].Index == usedGogoAnimeServerParams.DataType {
			episode.StreamOptions[idx].Used = true
			break
		}
	}

	return episode, nil
}

func (s *GogoAnimeNew) GetPerSeason(ctx context.Context, queryParams models.AnimeQueryParams) (models.AnimePerSeason, error) {
	animePerSeason := models.AnimePerSeason{}

	return animePerSeason, nil
}

func (s *GogoAnimeNew) GetRandom(ctx context.Context, queryParams models.AnimeQueryParams) ([]models.Anime, error) {
	animes := []models.Anime{}

	return animes, nil
}

func (s *GogoAnimeNew) getStreamUrl(ctx context.Context, params GogoAnimeServerParams) (GogoAnimeStream, error) {
	gogoAnimeStream := GogoAnimeStream{}

	apiParams := []string{
		fmt.Sprintf("%s=%s", params.DataType, params.DataEncryptedUrl1),
		fmt.Sprintf("url2=%s", params.DataEncryptedUrl2),
		fmt.Sprintf("url3=%s", params.DataEncryptedUrl3),
		fmt.Sprintf("feature_image=%s", url.QueryEscape(params.FeatureImage)),
		fmt.Sprintf("user_agent=Mozilla%%2F5.0+(Linux%%3B+Android+10%%3B+K)+AppleWebKit%%2F537.36+(KHTML%%2C+like+Gecko)+Chrome%%2F131.0.0.0+Mobile+Safari%%2F537.36"),
	}

	targetUrl := fmt.Sprintf("%s/wp-content/plugins/video-player/includes/player/player.php?%s", s.Host, strings.Join(apiParams, "&"))

	c := colly.NewCollector()

	iframeSrc := ""
	c.OnHTML("iframe", func(e *colly.HTMLElement) {
		iframeSrc = e.Attr("src")
	})

	videoURL := ""
	c.OnHTML("script", func(e *colly.HTMLElement) {
		scriptContent := e.Text

		// Use regex to extract the file URL
		re := regexp.MustCompile(`file:\s*"(https?://[^"]+)"`)
		matches := re.FindStringSubmatch(scriptContent)

		if len(matches) < 2 {
			err := fmt.Errorf("video media not found")
			logrus.WithContext(ctx).WithFields(logrus.Fields{
				"target_url": targetUrl,
				"mode":       params.DataType,
			}).Error(err)
		}
		videoURL = matches[1]
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("accept", "text/html, */*; q=0.01")
		r.Headers.Set("accept-language", "en-US,en;q=0.9,id;q=0.8")
		r.Headers.Set("cookie", "_ga=GA1.1.1050574509.1735699415; _ga_8KW6LYG84H=GS1.1.1735729345.2.1.1735730376.0.0.0")
		r.Headers.Set("priority", "u=1, i")
		r.Headers.Set("referer", params.EpisodeUrl)
		r.Headers.Set("sec-ch-ua", "\"Google Chrome\";v=\"131\", \"Chromium\";v=\"131\", \"Not_A Brand\";v=\"24\"")
		r.Headers.Set("sec-ch-ua-mobile", "?0")
		r.Headers.Set("sec-ch-ua-platform", "\"macOS\"")
		r.Headers.Set("sec-fetch-dest", "empty")
		r.Headers.Set("sec-fetch-mode", "cors")
		r.Headers.Set("sec-fetch-site", "same-origin")
		r.Headers.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
		r.Headers.Set("x-requested-with", "XMLHttpRequest")
	})

	err := c.Visit(targetUrl)
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"target_url": targetUrl,
		}).Error(err)
		return gogoAnimeStream, err
	}
	c.Wait()

	if params.DataType == "hianime" {
		gogoAnimeStream = GogoAnimeStream{
			Mode: "iframe",
			Src:  fmt.Sprintf("%s/%s", s.Host, iframeSrc),
		}
	} else if params.DataType == "double_player" {
		gogoAnimeStream = GogoAnimeStream{
			Mode: "mp4",
			Src:  videoURL,
		}
	} else {
		err := fmt.Errorf("stream mode not supported")
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"target_url": targetUrl,
			"mode":       params.DataType,
		}).Error(err)
		return gogoAnimeStream, err
	}

	if gogoAnimeStream.Mode == "" || gogoAnimeStream.Src == "" {
		err := fmt.Errorf("stream not found")
		logrus.WithContext(ctx).Error(err)
		return gogoAnimeStream, err
	}

	return gogoAnimeStream, nil
}

func (s *GogoAnimeNew) cleanUpUrl(str string) string {
	return strings.ReplaceAll(str, "\u0026", "&")
}
