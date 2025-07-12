package anime_scrapper_otakudesu

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"slices"
	"strings"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/contract"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/utils/anime_utils"
	"github.com/umarkotak/animapu-api/internal/utils/utils"
)

type Otakudesu struct {
	AnimapuSource        string
	Source               string
	OtakudesuHost        string
	AllowedStreamServers []string
}

func NewOtakudesu() Otakudesu {
	return Otakudesu{
		AnimapuSource: models.ANIME_SOURCE_OTAKUDESU,
		Source:        "otakudesu",
		OtakudesuHost: "https://otakudesu.cloud",
		AllowedStreamServers: []string{
			"filelions",
			"ondesuhd",
			"otakustream",
			"odstream",
			"pdrain",
			"", // whitelist all
		},
	}
}

func (s *Otakudesu) GetLatest(ctx context.Context, queryParams models.AnimeQueryParams) ([]contract.Anime, error) {
	animes := []contract.Anime{}

	c := colly.NewCollector()

	c.OnHTML("#venkonten > div > div.venser > div.venutama > div.rseries > div > div.venz > ul > li", func(e *colly.HTMLElement) {
		coverUrl := e.ChildAttr("div > div.thumb > a > div > img", "src")

		animeLink := e.ChildAttr("div > div.thumb > a", "href")
		splitted := strings.Split(animeLink, "/anime/")
		id := ""
		if len(splitted) > 0 {
			id = strings.ReplaceAll(splitted[len(splitted)-1], "/", "")
		}

		if id == "" {
			return
		}

		animes = append(animes, contract.Anime{
			ID:            id,
			Source:        s.Source,
			Title:         e.ChildText("div > div.thumb > a > div > h2"),
			LatestEpisode: utils.ForceSanitizeStringToFloat(e.ChildText("div > div.epz")),
			CoverUrls:     []string{coverUrl},
			OriginalLink:  animeLink,
		})
	})

	targetUrl := fmt.Sprintf("%v/ongoing-anime/page/%v", s.OtakudesuHost, queryParams.Page)
	err := c.Visit(targetUrl)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return animes, err
	}
	c.Wait()

	if queryParams.Page >= 4 {
		targetUrl := fmt.Sprintf("%v/complete-anime/page/%v", s.OtakudesuHost, queryParams.Page-3)
		err := c.Visit(targetUrl)
		if err != nil {
			logrus.WithContext(ctx).Error(err)
			return animes, err
		}
		c.Wait()
	}

	return animes, nil
}

func (s *Otakudesu) GetSearch(ctx context.Context, queryParams models.AnimeQueryParams) ([]contract.Anime, error) {
	queryParams.Title = strings.ReplaceAll(queryParams.Title, " ", "+")

	animes := []contract.Anime{}

	c := colly.NewCollector()

	c.OnHTML("#venkonten > div > div.venser > div > div > ul > li", func(e *colly.HTMLElement) {
		coverUrl := e.ChildAttr("img", "src")

		animeLink := e.ChildAttr("h2 > a", "href")
		splitted := strings.Split(animeLink, "/anime/")
		id := ""
		if len(splitted) > 0 {
			id = strings.ReplaceAll(splitted[len(splitted)-1], "/", "")
		}

		if id == "" {
			return
		}

		animes = append(animes, contract.Anime{
			ID:            id,
			Source:        s.Source,
			Title:         e.ChildText("h2 > a"),
			LatestEpisode: 0,
			CoverUrls:     []string{coverUrl},
			OriginalLink:  animeLink,
		})
	})

	targetUrl := fmt.Sprintf("%s/?s=%s&post_type=anime", s.OtakudesuHost, queryParams.Title)
	err := c.Visit(targetUrl)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return animes, err
	}
	c.Wait()

	return animes, nil
}

func (s *Otakudesu) GetDetail(ctx context.Context, queryParams models.AnimeQueryParams) (contract.Anime, error) {
	targetUrl := fmt.Sprintf("%v/anime/%v", s.OtakudesuHost, queryParams.SourceID)
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

	c.OnHTML("#venkonten > div.venser > div.fotoanime > div.infozin > div > p:nth-child(1) > span", func(e *colly.HTMLElement) {
		anime.Title = strings.ReplaceAll(e.Text, "Judul: ", "")
	})

	c.OnHTML("#venkonten > div.venser > div.fotoanime > img", func(e *colly.HTMLElement) {
		anime.CoverUrls = append(anime.CoverUrls, e.Attr("src"))
	})

	c.OnHTML("#venkonten > div.venser > div.fotoanime > div.infozin > div > p > span", func(e *colly.HTMLElement) {
		if strings.Contains(e.Text, "Tanggal Rilis") {
			releaseDateRaw := strings.ReplaceAll(e.Text, "Tanggal Rilis: ", "")
			anime.ReleaseDate = releaseDateRaw
			splitted := strings.Split(releaseDateRaw, " ")
			if len(splitted) != 3 {
				return
			}
			anime.ReleaseMonth = splitted[0]
			anime.ReleaseYear = utils.StringMustInt64(utils.RemoveNonNumeric(splitted[2]))
			anime.ReleaseSeason = anime_utils.OtakudesuMonthToSeason(anime.ReleaseMonth)
		}
	})

	c.OnHTML("#venkonten > div.venser > div.fotoanime > div.infozin > div > p > span", func(e *colly.HTMLElement) {
		if strings.Contains(e.Text, "Genre") {
			genreRaw := strings.ReplaceAll(e.Text, "Genre: ", "")
			anime.Genres = strings.Split(strings.ToLower(genreRaw), ",")
		}
	})

	c.OnHTML("#venkonten > div.venser > div.fotoanime > div.infozin > div > p > span", func(e *colly.HTMLElement) {
		if strings.Contains(e.Text, "Skor") {
			scoreRaw := strings.ReplaceAll(e.Text, "Skor: ", "")
			anime.Score = utils.ForceSanitizeStringToFloat(scoreRaw)
		}
	})

	c.OnHTML("div.episodelist", func(e *colly.HTMLElement) {
		checkText := e.ChildText("div > span > span")
		if !strings.Contains(strings.ToLower(checkText), strings.ToLower("Link Download Episode + Streaming")) {
			return
		}

		e.ForEach("ul > li", func(i int, h *colly.HTMLElement) {
			episodeLink := h.ChildAttr("span > a", "href")
			splitted := strings.Split(episodeLink, "/episode/")
			id := ""
			if len(splitted) > 0 {
				id = strings.ReplaceAll(splitted[len(splitted)-1], "/", "")
			}

			epTitle := h.ChildText("span > a")
			epTitleSplitted := strings.Split(epTitle, " ")

			epNo := float64(0)
			for _, content := range epTitleSplitted {
				if utils.ForceSanitizeStringToFloat(content) > 0 {
					epNo = utils.ForceSanitizeStringToFloat(content)
				}
			}

			episode := contract.Episode{
				AnimeID:      queryParams.SourceID,
				Source:       s.Source,
				ID:           id,
				Number:       epNo,
				Title:        epTitle,
				OriginalLink: episodeLink,
				UseTitle:     true,
			}
			anime.Episodes = append(anime.Episodes, episode)
		})
	})

	err := c.Visit(targetUrl)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return anime, err
	}
	c.Wait()

	for idx, _ := range anime.Episodes {
		// anime.Episodes[idx].Number = float64(len(anime.Episodes) - idx)
		anime.Episodes[idx].CoverUrl = anime.CoverUrls[0]
		anime.Episodes[idx].CoverUrls = anime.CoverUrls
	}

	slices.Reverse(anime.Episodes)

	return anime, nil
}

func (s *Otakudesu) Watch(ctx context.Context, queryParams models.AnimeQueryParams) (contract.EpisodeWatch, error) {
	if queryParams.Resolution == "" {
		queryParams.Resolution = "720p"
	}

	streamOptions := []contract.StreamOption{}

	episodeWatch := contract.EpisodeWatch{}

	c := colly.NewCollector()

	shortLink := ""
	c.OnHTML("link", func(e *colly.HTMLElement) {
		if e.Attr("rel") == "shortlink" {
			shortLink = e.Attr("href")
		}
	})

	type serverOpt struct {
		Name string
		Idx  string
	}
	streams := map[string][]serverOpt{
		"720p": {},
		"480p": {},
		"360p": {},
	}
	c.OnHTML("#venkonten > div.venser > div.venutama > div.mirrorstream > ul.m720p", func(e *colly.HTMLElement) {
		e.ForEach("a", func(i int, h *colly.HTMLElement) {
			streams["720p"] = append(streams["720p"], serverOpt{
				Name: h.Text,
				Idx:  fmt.Sprint(i),
			})
			streamOptions = append(streamOptions, contract.StreamOption{
				Resolution: "720p",
				Index:      fmt.Sprint(i),
				Name:       h.Text,
			})
		})
	})
	c.OnHTML("#venkonten > div.venser > div.venutama > div.mirrorstream > ul.m480p", func(e *colly.HTMLElement) {
		e.ForEach("a", func(i int, h *colly.HTMLElement) {
			streams["480p"] = append(streams["480p"], serverOpt{
				Name: h.Text,
				Idx:  fmt.Sprint(i),
			})
			streamOptions = append(streamOptions, contract.StreamOption{
				Resolution: "480p",
				Index:      fmt.Sprint(i),
				Name:       h.Text,
			})
		})
	})
	c.OnHTML("#venkonten > div.venser > div.venutama > div.mirrorstream > ul.m360p", func(e *colly.HTMLElement) {
		e.ForEach("a", func(i int, h *colly.HTMLElement) {
			streams["360p"] = append(streams["360p"], serverOpt{
				Name: h.Text,
				Idx:  fmt.Sprint(i),
			})
			streamOptions = append(streamOptions, contract.StreamOption{
				Resolution: "360p",
				Index:      fmt.Sprint(i),
				Name:       h.Text,
			})
		})
	})

	targetUrl := fmt.Sprintf("%v/episode/%v", s.OtakudesuHost, queryParams.EpisodeID)
	err := c.Visit(targetUrl)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return episodeWatch, err
	}
	c.Wait()

	// logrus.Infof("STREAM SERVER: %+v", streams)

	if len(streams[queryParams.Resolution]) <= 0 {
		backupFound := false
		for k, oneStream := range streams {
			if len(oneStream) > 0 {
				queryParams.Resolution = k
				backupFound = true
				break
			}
		}

		if !backupFound {
			err = fmt.Errorf(fmt.Sprintf("%s stream server not found", queryParams.Resolution))
			logrus.WithContext(ctx).Error(err)
			return episodeWatch, err
		}
	}

	shortLinkUrl, err := url.Parse(shortLink)
	if err != nil {
		err = fmt.Errorf("invalid short link url")
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"short_link": shortLink,
		}).Error(err)
		return episodeWatch, err
	}

	p := shortLinkUrl.Query().Get("p")
	if p == "" {
		err = fmt.Errorf("missing short link p")
		logrus.WithContext(ctx).Error(err)
		return episodeWatch, err
	}

	nonceBody, err := s.AdminAjaxCaller("aa1208d27f29ca340c92c66d1926f13f", []string{})
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return episodeWatch, err
	}
	// logrus.Infof("NONCE BODY: %+v", string(nonceBody))

	nonceData := map[string]string{}
	json.Unmarshal(nonceBody, &nonceData)
	nonce := nonceData["data"]
	if nonce == "" {
		err = fmt.Errorf("missing nonce p")
		logrus.WithContext(ctx).Error(err)
		return episodeWatch, err
	}

	iframeFinalUrl := ""
	logrus.Infof("QUERY PARAMS: %+v", queryParams)

	if queryParams.StreamIdx != "" {
		iframeBody, err := s.AdminAjaxCaller("2a3505c93b0035d3f455df82bf976b84", []string{
			fmt.Sprintf("id=%v", p),
			fmt.Sprintf("i=%v", queryParams.StreamIdx),
			fmt.Sprintf("q=%v", queryParams.Resolution),
			fmt.Sprintf("nonce=%v", nonce),
		})
		if err != nil {
			logrus.WithContext(ctx).Error(err)
			return contract.EpisodeWatch{}, err
		}

		iframeBase64Data := map[string]string{}
		json.Unmarshal(iframeBody, &iframeBase64Data)
		iframeBase64 := iframeBase64Data["data"]
		if iframeBase64 == "" {
			err = fmt.Errorf("missing iframe data")
			logrus.WithContext(ctx).Error(err)
			return contract.EpisodeWatch{}, err
		}

		iframeBase64Decoded, err := base64.StdEncoding.DecodeString(iframeBase64)
		if err != nil {
			logrus.WithContext(ctx).Error(err)
			return contract.EpisodeWatch{}, err
		}

		re := regexp.MustCompile(`src="([^"]+)"`)

		// Find the matches
		match := re.FindStringSubmatch(string(iframeBase64Decoded))

		if len(match) > 1 {
			iframeFinalUrl = match[1]
		}

		if iframeFinalUrl != "" {
			episodeWatch = contract.EpisodeWatch{
				StreamType:    "iframe",
				IframeUrl:     iframeFinalUrl,
				OriginalUrl:   targetUrl,
				StreamOptions: streamOptions,
				Resolution:    queryParams.Resolution,
				StreamIdx:     queryParams.StreamIdx,
			}

			return episodeWatch, nil
		}
	}

	selectedResolution := ""
	selectedStreamIdx := ""
	if iframeFinalUrl == "" {
		for _, ondesuIdx := range streams[queryParams.Resolution] {
			iframeBody, err := s.AdminAjaxCaller("2a3505c93b0035d3f455df82bf976b84", []string{
				fmt.Sprintf("id=%v", p),
				fmt.Sprintf("i=%v", ondesuIdx.Idx),
				fmt.Sprintf("q=%v", queryParams.Resolution),
				fmt.Sprintf("nonce=%v", nonce),
			})
			if err != nil {
				logrus.WithContext(ctx).Error(err)
				continue
			}

			iframeBase64Data := map[string]string{}
			json.Unmarshal(iframeBody, &iframeBase64Data)
			iframeBase64 := iframeBase64Data["data"]
			if iframeBase64 == "" {
				err = fmt.Errorf("missing iframe data")
				logrus.WithContext(ctx).Error(err)
				continue
			}

			iframeBase64Decoded, err := base64.StdEncoding.DecodeString(iframeBase64)
			if err != nil {
				logrus.WithContext(ctx).Error(err)
				continue
			}

			re := regexp.MustCompile(`src="([^"]+)"`)

			// Find the matches
			match := re.FindStringSubmatch(string(iframeBase64Decoded))

			if len(match) > 1 {
				iframeFinalUrl = match[1]
				selectedResolution = queryParams.Resolution
				selectedStreamIdx = ondesuIdx.Idx
				break
			}

			if iframeFinalUrl == "" {
				err = fmt.Errorf("missing final iframe url")
				logrus.WithContext(ctx).WithFields(logrus.Fields{
					"iframe_element": string(iframeBase64Decoded),
				}).Error(err)
				continue
			}
		}
	}

	if iframeFinalUrl == "" {
		err = fmt.Errorf("final iframe url not found at all")
		return contract.EpisodeWatch{}, err
	}

	episodeWatch = contract.EpisodeWatch{
		StreamType:    "iframe",
		IframeUrl:     iframeFinalUrl,
		OriginalUrl:   targetUrl,
		StreamOptions: streamOptions,
		Resolution:    selectedResolution,
		StreamIdx:     selectedStreamIdx,
	}

	return episodeWatch, nil
}

func (s *Otakudesu) GetPerSeason(ctx context.Context, queryParams models.AnimeQueryParams) (contract.AnimePerSeason, error) {
	animePerSeason := contract.AnimePerSeason{
		ReleaseYear: queryParams.ReleaseYear,
		SeasonName:  queryParams.ReleaseSeason,
		SeasonIndex: models.SEASON_TO_SEASON_INDEX[queryParams.ReleaseSeason],
		Animes:      []contract.Anime{},
	}

	return animePerSeason, nil
}

func (s *Otakudesu) GetRandom(ctx context.Context, queryParams models.AnimeQueryParams) ([]contract.Anime, error) {
	return []contract.Anime{}, nil
}
