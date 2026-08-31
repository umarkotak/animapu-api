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

var r2URLPattern = regexp.MustCompile(`(?i)https?://[^\s"'\\]*r2\.cloudflarestorage\.com[^\s"'\\]*`)

func r2URLs(value string) []string {
	return r2URLPattern.FindAllString(strings.ReplaceAll(value, `\/`, "/"), -1)
}

func New() Kuramanime {
	return Kuramanime{
		KuramanimeAggregator: "https://kuramalink.app",
	}
}

var WAIT_DURATION = 60 * time.Second

func logError(ctx context.Context, operation string, err error, fields logrus.Fields) {
	logrus.WithContext(ctx).WithError(err).WithFields(fields).Error("Kuramanime " + operation)
}

func (s *Kuramanime) GetLatest(ctx context.Context, queryParams models.AnimeQueryParams) ([]contract.Anime, error) {
	page := queryParams.Page
	if page < 1 {
		page = 1
	}
	host := s.getHost()
	targetURL := fmt.Sprintf("%s/quick/ongoing?order_by=updated&page=%d&need_json=true", host, page)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		logError(ctx, "create latest request", err, logrus.Fields{"url": targetURL, "page": page})
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logError(ctx, "request latest anime", err, logrus.Fields{"url": targetURL, "page": page})
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		err := fmt.Errorf("kuramanime returned status %s", resp.Status)
		logError(ctx, "request latest anime", err, logrus.Fields{"url": targetURL, "page": page, "status_code": resp.StatusCode})
		return nil, err
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
		logError(ctx, "read latest response", err, logrus.Fields{"url": targetURL})
		return nil, err
	}
	if err := sonic.Unmarshal(body, &payload); err != nil {
		logError(ctx, "parse latest response", err, logrus.Fields{"url": targetURL, "body_length": len(body)})
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
		logError(ctx, "create search request", err, logrus.Fields{"url": targetURL, "title": queryParams.Title})
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logError(ctx, "request search anime", err, logrus.Fields{"url": targetURL, "title": queryParams.Title})
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		err := fmt.Errorf("kuramanime returned status %s", resp.Status)
		logError(ctx, "request search anime", err, logrus.Fields{"url": targetURL, "title": queryParams.Title, "status_code": resp.StatusCode})
		return nil, err
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
		logError(ctx, "read search response", err, logrus.Fields{"url": targetURL, "title": queryParams.Title})
		return nil, err
	}
	if err := sonic.Unmarshal(body, &payload); err != nil {
		logError(ctx, "parse search response", err, logrus.Fields{"url": targetURL, "title": queryParams.Title, "body_length": len(body)})
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
	if finalURL, err := utils.GetFinalURL(targetURL); err != nil {
		logError(ctx, "resolve detail URL", err, logrus.Fields{"url": targetURL, "anime_id": queryParams.SourceID})
	} else {
		targetURL = finalURL
	}

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
			logError(ctx, "parse episode list", err, logrus.Fields{"url": targetURL, "anime_id": queryParams.SourceID})
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
		logError(ctx, "request anime detail", err, logrus.Fields{"url": targetURL, "anime_id": queryParams.SourceID})
		return anime, err
	}
	c.Wait()

	return anime, nil
}

func (s *Kuramanime) Watch(ctx context.Context, queryParams models.AnimeQueryParams) (contract.EpisodeWatch, error) {
	for attempt := 1; attempt <= 2; attempt++ {
		episodeWatch, err := s.watchOnce(ctx, queryParams)
		if err != nil || episodeWatch.StreamType != "" || attempt == 2 {
			return episodeWatch, err
		}

		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"anime_id":   queryParams.SourceID,
			"episode_id": queryParams.EpisodeID,
			"attempt":    attempt + 1,
		}).Warn("Kuramanime stream discovery timed out; retrying")
	}

	return contract.EpisodeWatch{}, nil
}

func (s *Kuramanime) watchOnce(ctx context.Context, queryParams models.AnimeQueryParams) (contract.EpisodeWatch, error) {
	host := s.getHost()
	targetURL := fmt.Sprintf("%s/anime/%s", host, queryParams.SourceID)
	if finalURL, err := utils.GetFinalURL(targetURL); err != nil {
		logError(ctx, "resolve watch URL", err, logrus.Fields{"url": targetURL, "anime_id": queryParams.SourceID, "episode_id": queryParams.EpisodeID})
	} else {
		targetURL = finalURL
	}
	targetURL = fmt.Sprintf("%s/episode/%s", targetURL, queryParams.EpisodeID)

	episodeWatch := contract.EpisodeWatch{
		OriginalUrl: targetURL,
		GdriveConf:  contract.GdriveConf{},
	}

	lease, err := datastore.NewBrowser(ctx)
	if err != nil {
		logError(ctx, "create isolated browser", err, logrus.Fields{"url": targetURL, "anime_id": queryParams.SourceID, "episode_id": queryParams.EpisodeID})
		return episodeWatch, err
	}
	defer lease.Release()
	browser := lease.Browser

	page, err := browser.Page(proto.TargetCreateTarget{})
	if err != nil {
		logError(ctx, "create watch page", err, logrus.Fields{"url": targetURL, "anime_id": queryParams.SourceID, "episode_id": queryParams.EpisodeID})
		return episodeWatch, err
	}
	defer page.Close()

	driveTokenFound := make(chan contract.GdriveConf, 1)
	driveTokenURL := fmt.Sprintf("%s/misc/token/drive-token", host)
	driveBrowser, stopDriveWait := browser.Context(ctx).WithCancel()
	driveEvents := driveBrowser.Event()
	driveDone := make(chan struct{})
	if err := (&proto.FetchEnable{Patterns: []*proto.FetchRequestPattern{{
		URLPattern:   driveTokenURL + "*",
		RequestStage: proto.FetchRequestStageResponse,
	}}}).Call(browser); err != nil {
		logError(ctx, "enable Drive token listener", err, logrus.Fields{"url": driveTokenURL})
		return episodeWatch, err
	}
	go func() {
		defer close(driveDone)
		for event := range driveEvents {
			if event.Method != "Fetch.requestPaused" {
				continue
			}
			paused := &proto.FetchRequestPaused{}
			if !event.Load(paused) || paused.Request == nil || paused.ResponseStatusCode == nil {
				continue
			}

			var foundDriveToken *contract.GdriveConf
			body, err := (proto.FetchGetResponseBody{RequestID: paused.RequestID}).Call(browser)
			if err != nil {
				logError(ctx, "read Drive token response", err, logrus.Fields{"url": paused.Request.URL, "status_code": *paused.ResponseStatusCode})
			} else {
				responseBody := body.Body
				if body.Base64Encoded {
					decoded, err := base64.StdEncoding.DecodeString(responseBody)
					if err != nil {
						logError(ctx, "decode Drive token response", err, logrus.Fields{"url": paused.Request.URL, "response_length": len(responseBody)})
						responseBody = ""
					} else {
						responseBody = string(decoded)
					}
				}

				var driveToken struct {
					Token string `json:"access_token"`
					ID    string `json:"gid"`
				}
				if err := sonic.Unmarshal([]byte(responseBody), &driveToken); err != nil {
					logError(ctx, "parse Drive token response", err, logrus.Fields{"url": paused.Request.URL, "status_code": *paused.ResponseStatusCode, "response_length": len(responseBody)})
				} else if driveToken.Token == "" || driveToken.ID == "" {
					logrus.WithContext(ctx).WithFields(logrus.Fields{"url": paused.Request.URL, "status_code": *paused.ResponseStatusCode, "has_token": driveToken.Token != "", "has_gid": driveToken.ID != ""}).Warn("Kuramanime Drive token response is incomplete")
				} else {
					foundDriveToken = &contract.GdriveConf{AccessToken: driveToken.Token, Gid: driveToken.ID}
				}
			}

			if err := (proto.FetchContinueRequest{RequestID: paused.RequestID}).Call(browser); err != nil {
				logError(ctx, "continue Drive token response", err, logrus.Fields{"url": paused.Request.URL})
				continue
			}
			if foundDriveToken != nil {
				select {
				case driveTokenFound <- *foundDriveToken:
				default:
				}
			}
		}
	}()
	defer func() {
		stopDriveWait()
		<-driveDone
		if err := (proto.FetchDisable{}).Call(browser); err != nil {
			logError(ctx, "disable Drive token listener", err, logrus.Fields{"url": driveTokenURL})
		}
	}()

	page = page.Context(ctx)
	if err := page.Navigate(targetURL); err != nil {
		logError(ctx, "navigate watch page", err, logrus.Fields{"url": targetURL, "anime_id": queryParams.SourceID, "episode_id": queryParams.EpisodeID})
		return episodeWatch, err
	}
	if err := page.WaitLoad(); err != nil {
		logError(ctx, "wait for watch page", err, logrus.Fields{"url": targetURL, "anime_id": queryParams.SourceID, "episode_id": queryParams.EpisodeID})
		return episodeWatch, err
	}

	restoreNetwork := page.EnableDomain(&proto.NetworkEnable{})
	defer restoreNetwork()

	mediaBrowser, stopMediaWait := browser.Context(ctx).WithCancel()
	mediaEvents := mediaBrowser.Event()
	mediaDone := make(chan struct{})
	r2StreamFound := make(chan string, 1)
	go func() {
		defer close(mediaDone)
		apiResponses := map[proto.NetworkRequestID]struct{}{}
		foundR2URLs := map[string]struct{}{}
		recordR2URLs := func(value string) {
			for _, url := range r2URLs(value) {
				if _, found := foundR2URLs[url]; found {
					continue
				}
				foundR2URLs[url] = struct{}{}
				logrus.WithField("url", url).Info("Kuramanime Cloudflare R2 media")
				select {
				case r2StreamFound <- url:
				default:
				}
			}
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
				recordR2URLs(response.Response.URL)
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
					logError(ctx, "read XHR response", err, logrus.Fields{"request_id": finished.RequestID})
					continue
				}
				responseBody := body.Body
				if body.Base64Encoded {
					decoded, err := base64.StdEncoding.DecodeString(responseBody)
					if err != nil {
						logError(ctx, "decode XHR response", err, logrus.Fields{"request_id": finished.RequestID, "response_length": len(responseBody)})
						continue
					}
					responseBody = string(decoded)
				}
				recordR2URLs(responseBody)
			}
		}
	}()
	defer func() {
		stopMediaWait()
		<-mediaDone
	}()

	playerButton, err := page.Element("#animeVideoPlayer > div.mb-3 > div > div.plyr.plyr--full-ui.plyr--video")
	if err != nil {
		logError(ctx, "find video player", err, logrus.Fields{"url": targetURL, "anime_id": queryParams.SourceID, "episode_id": queryParams.EpisodeID})
		return episodeWatch, err
	}
	time.Sleep(time.Second)
	if err := playerButton.Click(proto.InputMouseButtonLeft, 1); err != nil {
		logError(ctx, "click video player", err, logrus.Fields{"url": targetURL, "anime_id": queryParams.SourceID, "episode_id": queryParams.EpisodeID})
		return episodeWatch, err
	}

	select {
	case <-ctx.Done():
		err := ctx.Err()
		logError(ctx, "wait for stream", err, logrus.Fields{"url": targetURL, "anime_id": queryParams.SourceID, "episode_id": queryParams.EpisodeID})
		return episodeWatch, err
	case gdriveConf := <-driveTokenFound:
		episodeWatch.StreamType = "gdrive"
		episodeWatch.GdriveConf = gdriveConf
	case r2StreamURL := <-r2StreamFound:
		episodeWatch.StreamType = "mp4"
		episodeWatch.RawStreamUrl = strings.ReplaceAll(r2StreamURL, "&amp;", "&")
	case <-time.After(WAIT_DURATION):
		logrus.WithContext(ctx).WithFields(logrus.Fields{"url": targetURL, "anime_id": queryParams.SourceID, "episode_id": queryParams.EpisodeID, "wait_duration": WAIT_DURATION.String()}).Warn("Kuramanime stream discovery timed out")
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
