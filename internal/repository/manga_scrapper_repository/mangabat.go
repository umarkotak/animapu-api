package manga_scrapper_repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/config"
	"github.com/umarkotak/animapu-api/internal/contract"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/utils/utils"
)

type Mangabat struct {
	Host     string
	ReadHost string
}

func NewMangabat() Mangabat {
	return Mangabat{
		Host:     "https://www.mangabats.com",
		ReadHost: "https://readmangabat.com",
	}
}

func (m *Mangabat) GetHome(ctx context.Context, queryParams models.QueryParams) ([]contract.Manga, error) {
	mangas := []contract.Manga{}

	if queryParams.Page <= 0 {
		queryParams.Page = 1
	}

	c := colly.NewCollector()
	c.SetRequestTimeout(config.Get().CollyTimeout)

	c.OnHTML("body > div.container > div.main-wrapper > div > div > div.list-comic-item-wrap", func(e *colly.HTMLElement) {
		sourceID := ""
		mangaLink := e.ChildAttr("a.cover", "href")
		mangaLinkSplitted := strings.Split(mangaLink, "/")
		if len(mangaLinkSplitted) > 0 {
			sourceID = mangaLinkSplitted[len(mangaLinkSplitted)-1]
		}

		latestChapterID := ""
		latestChapterLink := e.ChildAttr("a.list-story-item-wrap-chapter", "href")
		mangaLinkSplittedSplitted := strings.Split(latestChapterLink, "/")
		if len(mangaLinkSplittedSplitted) > 0 {
			latestChapterID = mangaLinkSplittedSplitted[len(mangaLinkSplittedSplitted)-1]
		}

		latestChapterText := latestChapterID
		latestChapterNumber := utils.ForceSanitizeStringToFloat(latestChapterText)

		imageURL := e.ChildAttr("a.cover > img", "src")
		imageURL = fmt.Sprintf("%v/mangas/mangabat/image_proxy/%v", config.Get().AnimapuOnlineHost, imageURL)

		mangas = append(mangas, contract.Manga{
			ID:                  sourceID,
			Source:              "mangabat",
			SourceID:            sourceID,
			Title:               e.ChildText("div > h3 > a"),
			Description:         "",
			Genres:              []string{},
			Status:              "",
			Rating:              "",
			LatestChapterID:     latestChapterID,
			LatestChapterNumber: latestChapterNumber,
			LatestChapterTitle:  latestChapterText,
			Chapters:            []contract.Chapter{},
			CoverImages: []contract.CoverImage{
				{
					Index: 1,
					ImageUrls: []string{
						imageURL,
					},
				},
			},
		})
	})

	err := c.Visit(fmt.Sprintf("%v/manga-list/latest-manga?page=%v", m.Host, queryParams.Page))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return mangas, err
	}
	c.Wait()

	return mangas, nil
}

func (m *Mangabat) GetDetail(ctx context.Context, queryParams models.QueryParams) (contract.Manga, error) {
	manga := contract.Manga{
		Source:      "mangabat",
		SourceID:    queryParams.SourceID,
		Status:      "Ongoing",
		Chapters:    []contract.Chapter{},
		Description: "Description unavailable",
		CoverImages: []contract.CoverImage{{ImageUrls: []string{""}}},
	}
	c := colly.NewCollector()
	c.SetRequestTimeout(config.Get().CollyTimeout)

	c.OnHTML("body > div.container > div.main-wrapper > div.leftCol > div.manga-info-top > div.manga-info-content > ul > li:nth-child(1) > h1", func(e *colly.HTMLElement) {
		manga.Title = e.Text
	})

	c.OnHTML("body > div.container > div.main-wrapper > div.leftCol > div.manga-info-top > div.manga-info-pic > img", func(e *colly.HTMLElement) {
		imageURL := e.Attr("src")
		imageURL = fmt.Sprintf("%v/mangas/mangabat/image_proxy/%v", config.Get().AnimapuOnlineHost, imageURL)
		manga.CoverImages = []contract.CoverImage{
			{
				Index:     1,
				ImageUrls: []string{imageURL},
			},
		}
	})

	c.OnHTML("#contentBox", func(e *colly.HTMLElement) {
		manga.Description = e.Text
	})

	// Fetch chapters via API
	chapters, err := m.fetchChapters(ctx, queryParams.SourceID)
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"source_id": queryParams.SourceID,
		}).Warn("Failed to fetch chapters via API, continuing without chapters: ", err)
	} else {
		manga.Chapters = chapters
	}

	err = c.Visit(fmt.Sprintf("%s/manga/%s", m.Host, queryParams.SourceID))
	c.Wait()

	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return manga, err
	}

	manga.GenerateLatestChapter()

	return manga, nil
}

func (m *Mangabat) GetSearch(ctx context.Context, queryParams models.QueryParams) ([]contract.Manga, error) {
	mangas := []contract.Manga{}

	c := colly.NewCollector()
	c.SetRequestTimeout(config.Get().CollyTimeout)

	c.OnRequest(func(r *colly.Request) {
		// r.Headers.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		// r.Headers.Set("accept-language", "en-US,en;q=0.9,id;q=0.8")
		// r.Headers.Set("cache-control", "max-age=0")
		// r.Headers.Set("if-modified-since", "Sun, 01 Jun 2025 06:15:42 GMT")
		// r.Headers.Set("priority", "u=0, i")
		r.Headers.Set("referer", "https://www.mangabats.com/")
		// r.Headers.Set("sec-ch-ua", "\"Chromium\";v=\"136\", \"Google Chrome\";v=\"136\", \"Not.A/Brand\";v=\"99\"")
		// r.Headers.Set("sec-ch-ua-mobile", "?0")
		// r.Headers.Set("sec-ch-ua-platform", "\"macOS\"")
		// r.Headers.Set("sec-fetch-dest", "document")
		// r.Headers.Set("sec-fetch-mode", "navigate")
		// r.Headers.Set("sec-fetch-site", "same-origin")
		// r.Headers.Set("sec-fetch-user", "?1")
		// r.Headers.Set("upgrade-insecure-requests", "1")
		// r.Headers.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36")
	})

	c.OnHTML("body > div.container > div.main-wrapper > div.leftCol > div.daily-update > div > div", func(e *colly.HTMLElement) {
		sourceID := ""
		mangaLink := e.ChildAttr("a", "href")
		mangaLinkSplitted := strings.Split(mangaLink, "/")
		if len(mangaLinkSplitted) > 0 {
			sourceID = mangaLinkSplitted[len(mangaLinkSplitted)-1]
		}

		title := e.ChildText("div > h3 > a")

		latestChapterLink := e.ChildAttr("div > em:nth-child(2) > a", "href")
		latestChapterLinkSplitted := strings.Split(latestChapterLink, "/")
		latestChapterID := ""
		if len(latestChapterLinkSplitted) > 0 {
			latestChapterID = latestChapterLinkSplitted[len(latestChapterLinkSplitted)-1]
		}

		imageURL := e.ChildAttr("a > img", "src")
		imageURL = fmt.Sprintf("%v/mangas/mangabat/image_proxy/%v", config.Get().AnimapuOnlineHost, imageURL)

		mangas = append(mangas, contract.Manga{
			ID:                  sourceID,
			SourceID:            sourceID,
			Source:              "mangabat",
			Title:               title,
			Description:         "",
			Genres:              []string{},
			Status:              "",
			Rating:              "",
			LatestChapterID:     latestChapterID,
			LatestChapterNumber: utils.ForceSanitizeStringToFloat(latestChapterID),
			LatestChapterTitle:  latestChapterID,
			Chapters:            []contract.Chapter{},
			CoverImages: []contract.CoverImage{
				{Index: 1, ImageUrls: []string{imageURL}},
			},
		})
	})

	query := strings.Replace(queryParams.Title, " ", "_", -1)
	url := fmt.Sprintf("%s/search/story/%s", m.Host, query)
	err := c.Visit(url)
	c.Wait()
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"url": url,
		}).Error(err)
		return mangas, err
	}

	return mangas, nil
}

func (m *Mangabat) GetChapter(ctx context.Context, queryParams models.QueryParams) (contract.Chapter, error) {
	c := colly.NewCollector()

	chapter := contract.Chapter{
		ID:            queryParams.ChapterID,
		SourceID:      queryParams.SourceID,
		Source:        "mangabat",
		Number:        utils.ForceSanitizeStringToFloat(queryParams.ChapterID),
		ChapterImages: []contract.ChapterImage{},
	}

	c.OnRequest(func(r *colly.Request) {
		// r.Headers.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		// r.Headers.Set("accept-language", "en-US,en;q=0.9,id;q=0.8")
		// r.Headers.Set("cache-control", "max-age=0")
		// r.Headers.Set("if-modified-since", "Sun, 01 Jun 2025 06:15:42 GMT")
		// r.Headers.Set("priority", "u=0, i")
		r.Headers.Set("referer", "https://www.mangabats.com/")
		// r.Headers.Set("sec-ch-ua", "\"Chromium\";v=\"136\", \"Google Chrome\";v=\"136\", \"Not.A/Brand\";v=\"99\"")
		// r.Headers.Set("sec-ch-ua-mobile", "?0")
		// r.Headers.Set("sec-ch-ua-platform", "\"macOS\"")
		// r.Headers.Set("sec-fetch-dest", "document")
		// r.Headers.Set("sec-fetch-mode", "navigate")
		// r.Headers.Set("sec-fetch-site", "same-origin")
		// r.Headers.Set("sec-fetch-user", "?1")
		// r.Headers.Set("upgrade-insecure-requests", "1")
		// r.Headers.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36")
	})

	c.OnHTML("body > div.container-chapter-reader > img", func(e *colly.HTMLElement) {
		imageURL := e.Attr("src")
		imageURL = fmt.Sprintf("%v/mangas/mangabat/image_proxy/%v", config.Get().AnimapuOnlineHost, imageURL)

		chapter.ChapterImages = append(chapter.ChapterImages, contract.ChapterImage{
			Index:     0,
			ImageUrls: []string{imageURL},
		})
	})

	var err error
	targets := []string{
		fmt.Sprintf("%s/manga/%s/%s", m.Host, queryParams.SourceID, queryParams.ChapterID),
	}

	for _, targetLink := range targets {
		err = c.Request(
			"GET",
			targetLink,
			strings.NewReader("{}"),
			colly.NewContext(),
			http.Header{
				// "Authority":                 []string{strings.ReplaceAll(m.ReadHost, "https://", "")},
				// "Accept":                    []string{"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
				// "Accept-Language":           []string{"en-US,en;q=0.9,id;q=0.8"},
				// "Cache-Control":             []string{"max-age=0"},
				// "Sec-Fetch-Site":            []string{"same-origin"},
				// "Sec-Fetch-Mode":            []string{"navigate"},
				// "Sec-Fetch-Dest":            []string{"document"},
				// "Sec-Fetch-User":            []string{"?1"},
				// "Upgrade-Insecure-Requests": []string{"1"},
				// "User-Agent":                []string{"Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1"},
				// "Referer":                   []string{targetLink},
			},
		)
		if err != nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{
				"target": targetLink,
			}).Error(err)
			continue
		}
		if len(chapter.ChapterImages) > 0 {
			chapter.SourceLink = targetLink
			break
		}
	}

	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return chapter, err
	}

	return chapter, nil
}

// ChapterAPIResponse represents the response from the chapters API
type ChapterAPIResponse struct {
	Success bool                   `json:"success"`
	Data    ChapterAPIResponseData `json:"data"`
}

type ChapterAPIResponseData struct {
	Chapters   []ChapterAPIItem     `json:"chapters"`
	Pagination ChapterAPIPagination `json:"pagination"`
}

type ChapterAPIItem struct {
	ChapterName string  `json:"chapter_name"`
	ChapterSlug string  `json:"chapter_slug"`
	ChapterNum  float64 `json:"chapter_num"`
	UpdatedAt   string  `json:"updated_at"`
	View        int     `json:"view"`
}

type ChapterAPIPagination struct {
	Total   int  `json:"total"`
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
	HasMore bool `json:"has_more"`
}

// fetchChapters fetches chapters from the mangabat API with pagination
func (m *Mangabat) fetchChapters(ctx context.Context, sourceID string) ([]contract.Chapter, error) {
	chapters := []contract.Chapter{}
	limit := 50
	offset := 0

	for {
		apiURL := fmt.Sprintf("%s/api/manga/%s/chapters?limit=%d&offset=%d", m.Host, sourceID, limit, offset)

		req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
		if err != nil {
			return chapters, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("accept", "*/*")
		req.Header.Set("accept-language", "en-US,en;q=0.9")
		req.Header.Set("referer", fmt.Sprintf("%s/manga/%s", m.Host, sourceID))
		req.Header.Set("sec-fetch-dest", "empty")
		req.Header.Set("sec-fetch-mode", "cors")
		req.Header.Set("sec-fetch-site", "same-origin")
		req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36")

		client := &http.Client{Timeout: config.Get().CollyTimeout}
		resp, err := client.Do(req)
		if err != nil {
			return chapters, fmt.Errorf("failed to fetch chapters: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return chapters, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return chapters, fmt.Errorf("failed to read response body: %w", err)
		}

		var apiResp ChapterAPIResponse
		if err := json.Unmarshal(body, &apiResp); err != nil {
			return chapters, fmt.Errorf("failed to parse response: %w", err)
		}

		for _, ch := range apiResp.Data.Chapters {
			chapters = append(chapters, contract.Chapter{
				ID:       ch.ChapterSlug,
				Source:   "mangabat",
				SourceID: ch.ChapterSlug,
				Title:    ch.ChapterName,
				Index:    int64(len(chapters) + 1),
				Number:   ch.ChapterNum,
			})
		}

		// Check if there are more chapters to fetch
		if !apiResp.Data.Pagination.HasMore {
			break
		}

		offset += limit
	}

	return chapters, nil
}
