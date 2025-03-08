package manga_scrapper_repository

import (
	"context"
	"fmt"
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

	c.OnHTML("body > div.container > div.main-wrapper > div > div > div.list-truyen-item-wrap", func(e *colly.HTMLElement) {
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

	err := c.Visit(fmt.Sprintf("%v/genre/all?page=%v", m.Host, queryParams.Page))
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

	c.OnHTML("body > div.container > div.main-wrapper > div.leftCol > div.manga-info-top > ul > li:nth-child(1) > h1", func(e *colly.HTMLElement) {
		manga.Title = e.Text
	})

	c.OnHTML("body > div.container > div.main-wrapper > div.leftCol > div.manga-info-top > div > img", func(e *colly.HTMLElement) {
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

	idx := int64(1)
	c.OnHTML("#chapter > div > div.chapter-list > div", func(e *colly.HTMLElement) {
		chapterLink := e.ChildAttr("a", "href")
		splittedLink := strings.Split(chapterLink, "/")
		if len(splittedLink) == 0 {
			return
		}
		chapterID := splittedLink[len(splittedLink)-1]

		manga.Chapters = append(manga.Chapters, contract.Chapter{
			ID:       chapterID,
			Source:   "mangabat",
			SourceID: chapterID,
			Title:    e.ChildText("a"),
			Index:    idx,
			Number:   utils.ForceSanitizeStringToFloat(chapterID),
		})

		idx += 1
	})

	err := c.Visit(fmt.Sprintf("%s/manga/%s", m.Host, queryParams.SourceID))
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
	err := c.Visit(fmt.Sprintf("%s/search/story/%s", m.Host, query))
	c.Wait()
	if err != nil {
		logrus.WithContext(ctx).Error(err)
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
			logrus.WithContext(ctx).Error(err)
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
