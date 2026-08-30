package kuramanime

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/bytedance/sonic"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/datastore"
	"github.com/umarkotak/animapu-api/internal/contract"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/utils/anime_utils"
	"github.com/umarkotak/animapu-api/internal/utils/utils"
)

type Kuramanime struct {
	KuramanimeAggregator string
}

var mp4URLPattern = regexp.MustCompile(`(?i)https?://[^\s"'\\]+\.mp4[^\s"'\\]*`)

func mp4URLs(value string) []string {
	return mp4URLPattern.FindAllString(strings.ReplaceAll(value, `\/`, "/"), -1)
}

func New() Kuramanime {
	return Kuramanime{
		KuramanimeAggregator: "https://kuramalink.app",
	}
}

var (
	WAIT_DURATION    = 15 * time.Second
	TARGET_MP4_COUNT = 6
)

func (s *Kuramanime) GetLatest(ctx context.Context, queryParams models.AnimeQueryParams) ([]contract.Anime, error) {
	page := queryParams.Page
	if page < 1 {
		page = 1
	}
	host := s.getHost()
	targetURL := fmt.Sprintf("%s/quick/ongoing?order_by=updated&page=%d&need_json=true", host, page)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("kuramanime returned status %s", resp.Status)
	}

	var payload struct {
		Animes struct {
			Data []struct {
				ID               int64   `json:"id"`
				Title            string  `json:"title"`
				Synopsis         string  `json:"synopsis"`
				Slug             string  `json:"slug"`
				LatestEpisode    float64 `json:"latest_episode"`
				Score            float64 `json:"score"`
				ImagePortraitURL string  `json:"image_portrait_url"`
				FullAltTitles    string  `json:"full_alt_titles"`
			} `json:"data"`
		} `json:"animes"`
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := sonic.Unmarshal(body, &payload); err != nil {
		return nil, err
	}

	animes := make([]contract.Anime, 0, len(payload.Animes.Data))
	for _, item := range payload.Animes.Data {
		altTitles := []string{}
		if item.FullAltTitles != "" {
			altTitles = strings.Split(item.FullAltTitles, ", ")
		}

		animes = append(animes, contract.Anime{
			ID:            fmt.Sprint(item.ID),
			Source:        models.ANIME_SOURCE_KURAMANIME,
			Title:         item.Title,
			AltTitles:     altTitles,
			Description:   item.Synopsis,
			LatestEpisode: item.LatestEpisode,
			CoverUrls:     []string{item.ImagePortraitURL},
			OriginalLink:  fmt.Sprintf("%s/anime/%d/%s", host, item.ID, item.Slug),
			Score:         item.Score,
		})
	}

	return animes, nil
}

func (s *Kuramanime) GetSearch(ctx context.Context, queryParams models.AnimeQueryParams) ([]contract.Anime, error) {
	host := s.getHost()
	targetURL := fmt.Sprintf("%s/anime?search=%s&order_by=oldest&need_json=true", s.getHost(), strings.ReplaceAll(queryParams.Title, " ", "+"))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("kuramanime returned status %s", resp.Status)
	}

	var payload struct {
		Animes struct {
			Data []struct {
				ID               int64   `json:"id"`
				Title            string  `json:"title"`
				Synopsis         string  `json:"synopsis"`
				Slug             string  `json:"slug"`
				LatestEpisode    float64 `json:"latest_episode"`
				Score            float64 `json:"score"`
				ImagePortraitURL string  `json:"image_portrait_url"`
				FullAltTitles    string  `json:"full_alt_titles"`
			} `json:"data"`
		} `json:"animes"`
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := sonic.Unmarshal(body, &payload); err != nil {
		return nil, err
	}

	animes := make([]contract.Anime, 0, len(payload.Animes.Data))
	for _, item := range payload.Animes.Data {
		altTitles := []string{}
		if item.FullAltTitles != "" {
			altTitles = strings.Split(item.FullAltTitles, ", ")
		}

		animes = append(animes, contract.Anime{
			ID:            fmt.Sprint(item.ID),
			Source:        models.ANIME_SOURCE_KURAMANIME,
			Title:         item.Title,
			AltTitles:     altTitles,
			Description:   item.Synopsis,
			LatestEpisode: item.LatestEpisode,
			CoverUrls:     []string{item.ImagePortraitURL},
			OriginalLink:  fmt.Sprintf("%s/anime/%d/%s", host, item.ID, item.Slug),
			Score:         item.Score,
		})
	}

	return animes, nil
}

func (s *Kuramanime) GetDetail(ctx context.Context, queryParams models.AnimeQueryParams) (contract.Anime, error) {
	host := s.getHost()
	targetURL := fmt.Sprintf("%s/anime/%s", host, queryParams.SourceID)
	targetURL, _ = utils.GetFinalURL(targetURL)

	anime := contract.Anime{
		ID:           queryParams.SourceID,
		Source:       models.ANIME_SOURCE_KURAMANIME,
		CoverUrls:    []string{},
		AltTitles:    []string{},
		Genres:       []string{},
		Episodes:     []contract.Episode{},
		OriginalLink: targetURL,
	}

	c := colly.NewCollector()
	c.OnHTML("body > section.anime-details.spad.pb-0 > div > div > div > div.col-lg-3 > div.set-bg", func(e *colly.HTMLElement) {
		if coverURL := e.Attr("data-setbg"); coverURL != "" {
			anime.CoverUrls = []string{coverURL}
		}
	})
	c.OnHTML(".anime__details__title", func(e *colly.HTMLElement) {
		anime.Title = strings.TrimSpace(e.ChildText("h3"))
		if altTitles := strings.TrimSpace(e.ChildText("h3 + span")); altTitles != "" {
			anime.AltTitles = strings.Split(altTitles, ", ")
		}
	})
	c.OnHTML("#synopsisField", func(e *colly.HTMLElement) {
		anime.Description = strings.TrimSpace(e.Text)
	})
	c.OnHTML(".anime__details__widget li", func(e *colly.HTMLElement) {
		label := strings.TrimSpace(e.ChildText("div > div.col-3 > span"))
		value := strings.TrimSpace(e.ChildText("div > div.col-9"))
		switch label {
		case "Tayang:":
			anime.ReleaseDate = strings.TrimSpace(strings.Split(value, "s/d")[0])
			parts := strings.Fields(anime.ReleaseDate)
			if len(parts) == 3 {
				anime.ReleaseMonth = parts[1]
				anime.ReleaseYear = utils.StringMustInt64(parts[2])
				anime.ReleaseSeason = anime_utils.OtakudesuMonthToSeason(parts[1])
				anime.ReleaseSeasonIndex = int64(anime_utils.SeasonToIndex(anime.ReleaseSeason))
			}
		case "Genre:":
			e.ForEach("div > div.col-9 > a", func(_ int, genre *colly.HTMLElement) {
				anime.Genres = append(anime.Genres, strings.ToLower(strings.Trim(strings.TrimSpace(genre.Text), ",")))
			})
		case "Skor:":
			anime.Score = utils.ForceSanitizeStringToFloat(value)
		}
	})
	c.OnHTML("#episodeLists", func(e *colly.HTMLElement) {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(e.Attr("data-content")))
		if err != nil {
			return
		}
		doc.Find("a[href]").Each(func(_ int, episode *goquery.Selection) {
			link, _ := episode.Attr("href")
			title := strings.TrimSpace(episode.Text())
			number := utils.ForceSanitizeStringToFloat(title)
			if link == "" || number == 0 {
				return
			}
			anime.Episodes = append(anime.Episodes, contract.Episode{
				AnimeID:      anime.ID,
				Source:       models.ANIME_SOURCE_KURAMANIME,
				ID:           fmt.Sprint(number),
				Number:       number,
				Title:        title,
				OriginalLink: link,
				CoverUrl:     strings.Join(anime.CoverUrls, ""),
				CoverUrls:    anime.CoverUrls,
				UseTitle:     true,
			})
		})
	})

	if err := c.Visit(targetURL); err != nil {
		return anime, err
	}
	c.Wait()

	return anime, nil
}

func (s *Kuramanime) Watch(ctx context.Context, queryParams models.AnimeQueryParams) (contract.EpisodeWatch, error) {
	host := s.getHost()
	targetURL := fmt.Sprintf("%s/anime/%s", host, queryParams.SourceID)
	targetURL, _ = utils.GetFinalURL(targetURL)
	targetURL = fmt.Sprintf("%s/episode/%s", targetURL, queryParams.EpisodeID)

	episodeWatch := contract.EpisodeWatch{
		OriginalUrl: targetURL,
		StreamType:  "gdrive",
		GdriveConf:  contract.GdriveConf{},
	}

	browser := datastore.Get().Browser

	page, err := browser.Page(proto.TargetCreateTarget{})
	if err != nil {
		return episodeWatch, err
	}
	defer page.Close()

	driveTokenFound := make(chan contract.GdriveConf, 1)
	router := browser.HijackRequests()
	if err := router.Add("*", "", func(h *rod.Hijack) {
		isMiscRequest := strings.HasPrefix(h.Request.URL().String(), fmt.Sprintf("%s/misc/token/drive-token", host))
		isAPIRequest := h.Request.Type() == proto.NetworkResourceTypeXHR || h.Request.Type() == proto.NetworkResourceTypeFetch
		if !isMiscRequest && !isAPIRequest {
			h.ContinueRequest(&proto.FetchContinueRequest{})
			return
		}
		if err := h.LoadResponse(http.DefaultClient, true); err != nil {
			logrus.WithError(err).WithField("url", h.Request.URL()).Warn("failed to read Kuramanime API response")
			h.ContinueRequest(&proto.FetchContinueRequest{})
			return
		}
		if isMiscRequest {
			var driveToken struct {
				Token string `json:"access_token"`
				ID    string `json:"gid"`
			}
			if err := sonic.Unmarshal([]byte(h.Response.Body()), &driveToken); err != nil {
				logrus.WithError(err).Warn("parse Kuramanime Drive token")
			} else {
				// logrus.WithFields(logrus.Fields{"token": driveToken.Token, "id": driveToken.ID}).Info("Kuramanime Drive token")
				if driveToken.Token != "" {
					select {
					case driveTokenFound <- contract.GdriveConf{AccessToken: driveToken.Token, Gid: driveToken.ID}:
					default:
					}
				}
			}
		}

		logrus.WithFields(logrus.Fields{
			"method":   h.Request.Method(),
			"type":     h.Request.Type(),
			"url":      h.Request.URL(),
			"request":  h.Request.Body(),
			"response": h.Response.Body(),
		}).Info("Kuramanime API response")
	}); err != nil {
		return episodeWatch, err
	}
	go router.Run()
	defer router.Stop()

	page = page.Context(ctx)
	if err := page.Navigate(targetURL); err != nil {
		return episodeWatch, err
	}
	if err := page.WaitLoad(); err != nil {
		return episodeWatch, err
	}

	restoreNetwork := page.EnableDomain(&proto.NetworkEnable{})
	defer restoreNetwork()

	mediaBrowser, stopMediaWait := browser.Context(ctx).WithCancel()
	mediaEvents := mediaBrowser.Event()
	mediaDone := make(chan struct{})
	go func() {
		defer close(mediaDone)
		apiResponses := map[proto.NetworkRequestID]struct{}{}
		foundMP4URLs := map[string]struct{}{}
		recordMP4URLs := func(value string) bool {
			for _, url := range mp4URLs(value) {
				if _, found := foundMP4URLs[url]; found {
					continue
				}
				foundMP4URLs[url] = struct{}{}
				logrus.WithField("url", url).Info("Kuramanime MP4 media")
				if len(foundMP4URLs) >= TARGET_MP4_COUNT {
					return true
				}
			}
			return false
		}
		for event := range mediaEvents {
			if event.SessionID != page.SessionID {
				continue
			}
			switch event.Method {
			case "Network.responseReceived":
				response := &proto.NetworkResponseReceived{}
				if !event.Load(response) || response.Response == nil {
					continue
				}
				if recordMP4URLs(response.Response.URL) {
					return
				}
				if response.Type == proto.NetworkResourceTypeXHR || response.Type == proto.NetworkResourceTypeFetch {
					apiResponses[response.RequestID] = struct{}{}
				}
			case "Network.loadingFinished":
				finished := &proto.NetworkLoadingFinished{}
				if !event.Load(finished) {
					continue
				}
				if _, found := apiResponses[finished.RequestID]; !found {
					continue
				}
				delete(apiResponses, finished.RequestID)
				body, err := (proto.NetworkGetResponseBody{RequestID: finished.RequestID}).Call(page)
				if err != nil {
					logrus.WithError(err).Debug("read Kuramanime API response")
					continue
				}
				responseBody := body.Body
				if body.Base64Encoded {
					if decoded, err := base64.StdEncoding.DecodeString(responseBody); err == nil {
						responseBody = string(decoded)
					}
				}
				if recordMP4URLs(responseBody) {
					return
				}
			}
		}
	}()
	defer func() {
		stopMediaWait()
		<-mediaDone
	}()

	playerButton, err := page.Element("#animeVideoPlayer > div.mb-3 > div > div.plyr.plyr--full-ui.plyr--video")
	if err != nil {
		return episodeWatch, err
	}
	if err := playerButton.Click(proto.InputMouseButtonLeft, 1); err != nil {
		return episodeWatch, err
	}

	select {
	case <-ctx.Done():
		return episodeWatch, ctx.Err()
	case gdriveConf := <-driveTokenFound:
		episodeWatch.GdriveConf = gdriveConf
	case <-mediaDone:
	case <-time.After(WAIT_DURATION):
	}

	return episodeWatch, nil
}

func (s *Kuramanime) GetPerSeason(ctx context.Context, queryParams models.AnimeQueryParams) (contract.AnimePerSeason, error) {
	return contract.AnimePerSeason{}, nil
}

func (s *Kuramanime) GetRandom(ctx context.Context, queryParams models.AnimeQueryParams) ([]contract.Anime, error) {
	return []contract.Anime{}, nil
}

func (s *Kuramanime) getHost() string {
	finalUrl, _ := utils.GetFinalURL(s.KuramanimeAggregator)

	return finalUrl
}
