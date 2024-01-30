package anime_scrapper_otakudesu

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"slices"
	"sort"
	"strings"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/local_db"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/utils/anime_utils"
	"github.com/umarkotak/animapu-api/internal/utils/utils"
)

type Otakudesu struct {
	AnimapuSource      string
	Source             string
	OtakudesuHost      string
	OtakudesuAuthority string
	DesusStreamHost    string
}

func NewOtakudesu() Otakudesu {
	return Otakudesu{
		AnimapuSource:      models.ANIME_SOURCE_OTAKUDESU,
		Source:             "otakudesu",
		OtakudesuHost:      "https://otakudesu.media",
		OtakudesuAuthority: "otakudesu.media",
		DesusStreamHost:    "https://desustream.me",
	}
}

func (s *Otakudesu) GetLatest(ctx context.Context, queryParams models.AnimeQueryParams) ([]models.Anime, error) {
	animes := []models.Anime{}

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

		animes = append(animes, models.Anime{
			ID:            id,
			Source:        s.Source,
			Title:         e.ChildText("div > div.thumb > a > div > h2"),
			LatestEpisode: utils.ForceSanitizeStringToFloat(e.ChildText("div > div.epz")),
			CoverUrls:     []string{coverUrl},
			OriginalLink:  animeLink,
		})
	})

	maxPage := 3
	for i := 1; i <= maxPage; i++ {
		targetUrl := fmt.Sprintf("%v/ongoing-anime/page/%v", s.OtakudesuHost, i)
		err := c.Visit(targetUrl)
		if err != nil {
			logrus.WithContext(ctx).Error(err)
			return animes, err
		}
		c.Wait()
	}

	return animes, nil
}

func (s *Otakudesu) GetSearch(ctx context.Context, queryParams models.AnimeQueryParams) ([]models.Anime, error) {
	return []models.Anime{}, nil
}

func (s *Otakudesu) GetDetail(ctx context.Context, queryParams models.AnimeQueryParams) (models.Anime, error) {
	targetUrl := fmt.Sprintf("%v/anime/%v", s.OtakudesuHost, queryParams.SourceID)
	anime := models.Anime{
		ID:            queryParams.SourceID,
		Source:        s.AnimapuSource,
		Title:         "",
		LatestEpisode: 0,          // done
		CoverUrls:     []string{}, // done
		Episodes:      []models.Episode{},
		OriginalLink:  targetUrl,
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

	c.OnHTML("div.episodelist > ul > li", func(e *colly.HTMLElement) {
		episodeLink := e.ChildAttr("span > a", "href")
		splitted := strings.Split(episodeLink, "/episode/")
		id := ""
		if len(splitted) > 0 {
			id = strings.ReplaceAll(splitted[len(splitted)-1], "/", "")
		}

		epNo := utils.ForceSanitizeStringToFloat(e.ChildText("span > a"))
		if epNo > 1500 {
			return
		}

		episode := models.Episode{
			AnimeID:      queryParams.SourceID,
			Source:       s.Source,
			ID:           id,
			Number:       epNo,
			Title:        e.ChildText("span > a"),
			OriginalLink: episodeLink,
		}
		anime.Episodes = append(anime.Episodes, episode)
	})

	err := c.Visit(targetUrl)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return anime, err
	}
	c.Wait()

	slices.Reverse(anime.Episodes)

	return anime, nil
}

func (s *Otakudesu) Watch(ctx context.Context, queryParams models.AnimeQueryParams) (models.EpisodeWatch, error) {
	episodeWatch := models.EpisodeWatch{}

	c := colly.NewCollector()

	shortLink := ""
	c.OnHTML("link", func(e *colly.HTMLElement) {
		if e.Attr("rel") == "shortlink" {
			shortLink = e.Attr("href")
		}
	})

	baypassedStreams := []string{
		// "ondesuhd",
		// "otakustream",
		// "odstream",
		"pdrain",
	}
	streams := map[string][]map[string]string{
		"360p": {},
		"480p": {},
		"720p": {},
	}
	c.OnHTML("#embed_holder > div.mirrorstream > ul.m720p", func(e *colly.HTMLElement) {
		e.ForEach("a", func(i int, h *colly.HTMLElement) {
			for _, oneBaypassedStream := range baypassedStreams {
				if strings.Contains(h.Text, oneBaypassedStream) {
					streams["720p"] = append(streams["720p"], map[string]string{
						"name": h.Text,
						"idx":  fmt.Sprint(i),
					})
				}
			}
		})
	})
	// c.OnHTML("#embed_holder > div.mirrorstream > ul.m360p", func(e *colly.HTMLElement) {
	// 	e.ForEach("a", func(i int, h *colly.HTMLElement) {
	// 		for _, oneBaypassedStream := range baypassedStreams {
	// 			if strings.Contains(h.Text, oneBaypassedStream) {
	// 				streams["720p"] = append(streams["720p"], map[string]string{
	// 					"name": h.Text,
	// 					"idx":  fmt.Sprint(i),
	// 				})
	// 			}
	// 		}
	// 	})
	// })

	targetUrl := fmt.Sprintf("%v/episode/%v", s.OtakudesuHost, queryParams.EpisodeID)
	err := c.Visit(targetUrl)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return episodeWatch, err
	}
	c.Wait()

	if len(streams["720p"]) <= 0 {
		err = fmt.Errorf("720 stream server not found")
		logrus.WithContext(ctx).Error(err)
		return episodeWatch, err
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
	logrus.Infof("NONCE BODY: %+v", string(nonceBody))

	nonceData := map[string]string{}
	json.Unmarshal(nonceBody, &nonceData)
	nonce := nonceData["data"]
	if nonce == "" {
		err = fmt.Errorf("missing nonce p")
		logrus.WithContext(ctx).Error(err)
		return episodeWatch, err
	}

	iframeFinalUrl := ""
	iframeFinalUrls := []string{}
	for _, ondesuIdx := range streams["720p"] {
		iframeBody, err := s.AdminAjaxCaller("2a3505c93b0035d3f455df82bf976b84", []string{
			fmt.Sprintf("id=%v", p),
			fmt.Sprintf("i=%v", ondesuIdx["idx"]),
			"q=720p",
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

		re := regexp.MustCompile(`src="(https:\/\/desustream\.me[^"]+)"`)
		re2 := regexp.MustCompile(`src="(https:\/\/www\.pixeldrain\.com[^"]+)"`)

		// Find the matches
		matches := re.FindStringSubmatch(string(iframeBase64Decoded))
		matches = append(matches, re2.FindStringSubmatch(string(iframeBase64Decoded))...)

		if len(matches) >= 2 {
			iframeFinalUrl = matches[1]
		}
		if iframeFinalUrl == "" {
			err = fmt.Errorf("missing final iframe url")
			logrus.WithContext(ctx).WithFields(logrus.Fields{
				"iframe_element": string(iframeBase64Decoded),
			}).Error(err)
			continue
		}

		iframeFinalUrls = append(iframeFinalUrls, iframeFinalUrl)
	}

	episodeWatch = models.EpisodeWatch{
		StreamType:  "iframe",
		IframeUrl:   iframeFinalUrl,
		IframeUrls:  iframeFinalUrls,
		OriginalUrl: targetUrl,
	}

	return episodeWatch, nil
}

func (s *Otakudesu) GetPerSeason(ctx context.Context, queryParams models.AnimeQueryParams) (models.AnimePerSeason, error) {
	animePerSeason := models.AnimePerSeason{
		ReleaseYear: queryParams.ReleaseYear,
		SeasonName:  queryParams.ReleaseSeason,
		SeasonIndex: models.SEASON_TO_SEASON_INDEX[queryParams.ReleaseSeason],
		Animes:      []models.Anime{},
	}

	otakudesuDB := local_db.AnimeLinkToDetailMap

	for _, oneAnime := range otakudesuDB {
		if oneAnime.ReleaseYear != queryParams.ReleaseYear {
			continue
		}

		if oneAnime.ReleaseSeason != queryParams.ReleaseSeason {
			continue
		}

		oneAnime.LatestEpisode = 0

		animePerSeason.Animes = append(animePerSeason.Animes, oneAnime)
	}

	sort.Slice(animePerSeason.Animes, func(i, j int) bool {
		return animePerSeason.Animes[i].Score < animePerSeason.Animes[j].Score
	})

	return animePerSeason, nil
}

func (s *Otakudesu) GetRandom(ctx context.Context, queryParams models.AnimeQueryParams) ([]models.Anime, error) {
	return []models.Anime{}, nil
}
