package manga_scrapper

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/config"
	"github.com/umarkotak/animapu-api/internal/models"
)

func GetMangabatLatestManga(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	mangas := []models.Manga{}

	if queryParams.Page <= 0 {
		queryParams.Page = 1
	}

	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	c.OnHTML("body div.body-site div.container.container-main div.container-main-left div.panel-list-story .list-story-item", func(e *colly.HTMLElement) {
		sourceID := strings.Replace(e.ChildAttr("div > h3 > a", "href"), "https://read.mangabat.com/", "", -1)
		sourceID = strings.Replace(sourceID, "https://m.mangabat.com/", "", -1)
		sourceID = strings.Replace(sourceID, "https://readmangabat.com/", "", -1)
		sourceID = strings.Replace(sourceID, "https://h.mangabat.com/", "", -1)

		imageURL := e.ChildAttr("a > img", "src")

		latestChapterText := e.ChildText("div > a:nth-child(2)")
		latestChapterID := strings.Replace(e.ChildAttr("div > a:nth-child(2)", "href"), "https://read.mangabat.com/", "", -1)
		latestChapterID = strings.Replace(latestChapterID, "https://m.mangabat.com/", "", -1)
		latestChapterID = strings.Replace(latestChapterID, "https://readmangabat.com/", "", -1)

		latestChapterNumberString := strings.Replace(latestChapterText, "Chapter ", "", -1)
		latestChapterNumber, _ := strconv.ParseFloat(latestChapterNumberString, 64)

		if latestChapterNumber == 0 {
			latestChapterNumberSplitted := strings.Split(latestChapterID, "-")
			if len(latestChapterNumberSplitted) > 0 {
				latestChapterNumberString = latestChapterNumberSplitted[len(latestChapterNumberSplitted)-1]
				latestChapterNumber, _ = strconv.ParseFloat(latestChapterNumberString, 64)
			}
		}

		mangas = append(mangas, models.Manga{
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
			Chapters:            []models.Chapter{},
			CoverImages: []models.CoverImage{
				{
					Index: 1,
					ImageUrls: []string{
						imageURL,
					},
				},
			},
		})
	})

	err := c.Visit(fmt.Sprintf("https://m.mangabat.com/manga-list-all/%v", queryParams.Page))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return mangas, err
	}

	return mangas, nil
}

func GetMangabatDetailManga(ctx context.Context, queryParams models.QueryParams) (models.Manga, error) {
	manga := models.Manga{
		Source:      "mangabat",
		SourceID:    queryParams.SourceID,
		Status:      "Ongoing",
		Chapters:    []models.Chapter{},
		Description: "Description unavailable",
		CoverImages: []models.CoverImage{{ImageUrls: []string{""}}},
	}
	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	c.OnHTML("body > div.body-site > div.container.container-main > div.container-main-left > div.panel-story-info > div.story-info-right > h1", func(e *colly.HTMLElement) {
		manga.Title = e.Text
	})

	c.OnHTML("body > div.body-site > div.container.container-main > div.container-main-left > div.panel-story-info > div.story-info-left > span.info-image > img", func(e *colly.HTMLElement) {
		manga.CoverImages = []models.CoverImage{
			{
				Index:     1,
				ImageUrls: []string{e.Attr("src")},
			},
		}
	})

	c.OnHTML("#panel-story-info-description", func(e *colly.HTMLElement) {
		manga.Description = e.Text
	})

	idx := int64(1)
	c.OnHTML("body > div.body-site > div.container.container-main > div.container-main-left > div.panel-story-chapter-list > ul > li", func(e *colly.HTMLElement) {
		chapterLink := e.ChildAttr("a", "href")
		chapterLink = strings.Replace(chapterLink, "https://read.mangabat.com/", "", -1)
		chapterLink = strings.Replace(chapterLink, "https://m.mangabat.com/", "", -1)
		chapterLink = strings.Replace(chapterLink, "https://h.mangabat.com/", "", -1)

		splittedLink := strings.Split(chapterLink, "-")
		chapterNumber, _ := strconv.ParseFloat(splittedLink[len(splittedLink)-1], 64)
		id := fmt.Sprintf("chap-%v", splittedLink[len(splittedLink)-1])

		manga.Chapters = append(manga.Chapters, models.Chapter{
			ID:       id,
			Source:   "mangabat",
			SourceID: id,
			Title:    e.ChildText("a"),
			Index:    idx,
			Number:   chapterNumber,
		})

		idx += 1
	})

	err := c.Visit(fmt.Sprintf("https://m.mangabat.com/%v", queryParams.SourceID))

	if manga.Title == "" {
		err = c.Visit(fmt.Sprintf("https://read.mangabat.com/%v", queryParams.SourceID))
	}

	if manga.Title == "" {
		err = c.Visit(fmt.Sprintf("https://readmangabat.com/%v", queryParams.SourceID))
	}

	if manga.Title == "" {
		err = c.Visit(fmt.Sprintf("https://h.mangabat.com/%v", queryParams.SourceID))
	}

	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return manga, err
	}

	manga.GenerateLatestChapter()

	return manga, nil
}

func GetMangabatByQuery(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	mangas := []models.Manga{}

	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	c.OnHTML("body > div.body-site > div.container.container-main > div.container-main-left > div.panel-list-story > div", func(e *colly.HTMLElement) {
		detailUrl := e.ChildAttr("a.item-img", "href")
		sourceID := strings.Replace(detailUrl, "https://read.mangabat.com/", "", -1)
		sourceID = strings.Replace(sourceID, "https://m.mangabat.com/", "", -1)
		sourceID = strings.Replace(sourceID, "https://readmangabat.com/", "", -1)
		sourceID = strings.Replace(sourceID, "https://h.mangabat.com/", "", -1)

		title := e.ChildText("div > h3 > a")

		latestChapterTitle := e.ChildAttr("div > a:nth-child(2)", "href")
		latestChapterTitleSplitted := strings.Split(latestChapterTitle, "-")
		latestChapterID := latestChapterTitleSplitted[len(latestChapterTitleSplitted)-1]
		latestChapterNumber, _ := strconv.ParseFloat(latestChapterID, 64)

		mangas = append(mangas, models.Manga{
			ID:                  sourceID,
			SourceID:            sourceID,
			Source:              "mangabat",
			Title:               title,
			Description:         "",
			Genres:              []string{},
			Status:              "",
			Rating:              "",
			LatestChapterID:     latestChapterID,
			LatestChapterNumber: latestChapterNumber,
			LatestChapterTitle:  latestChapterTitle,
			Chapters:            []models.Chapter{},
			CoverImages: []models.CoverImage{
				{
					Index: 1,
					ImageUrls: []string{
						e.ChildAttr("a > img", "src"),
					},
				},
			},
		})
	})

	query := strings.Replace(queryParams.Title, " ", "_", -1)
	err := c.Visit(fmt.Sprintf("https://m.mangabat.com/search/manga/%v", query))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return mangas, err
	}

	return mangas, nil
}

func GetMangabatDetailChapter(ctx context.Context, queryParams models.QueryParams) (models.Chapter, error) {
	c := colly.NewCollector()

	chapterNumberSplitted := strings.Split(queryParams.ChapterID, "-")
	chapterNumber, _ := strconv.ParseFloat(chapterNumberSplitted[1], 64)

	chapter := models.Chapter{
		ID:            queryParams.ChapterID,
		SourceID:      queryParams.SourceID,
		Source:        "mangabat",
		Number:        chapterNumber,
		ChapterImages: []models.ChapterImage{},
	}

	c.OnHTML("body > div.body-site > div.container-chapter-reader > img", func(e *colly.HTMLElement) {
		chapter.ChapterImages = append(chapter.ChapterImages, models.ChapterImage{
			Index: 0,
			ImageUrls: []string{
				fmt.Sprintf("%v/image_proxy?referer=%v&target=%v", config.Get().AnimapuOnlineHost, "https://m.mangabat.com/", e.Attr("src")),
			},
		})
	})

	var err error
	targets := []string{
		fmt.Sprintf("https://m.mangabat.com/%v-%v", queryParams.SourceID, queryParams.ChapterID),
		fmt.Sprintf("https://read.mangabat.com/%v-%v", queryParams.SourceID, queryParams.ChapterID),
		fmt.Sprintf("https://readmangabat.com/%v-%v", queryParams.SourceID, queryParams.ChapterID),
		fmt.Sprintf("https://h.mangabat.com/%v-%v", queryParams.SourceID, queryParams.ChapterID),
	}

	for _, targetLink := range targets {
		err = c.Request(
			"GET",
			targetLink,
			strings.NewReader("{}"),
			colly.NewContext(),
			http.Header{
				"Authority":                 []string{"readmangabat.com"},
				"Accept":                    []string{"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
				"Accept-Language":           []string{"en-US,en;q=0.9,id;q=0.8"},
				"Cache-Control":             []string{"max-age=0"},
				"Sec-Fetch-Site":            []string{"same-origin"},
				"Sec-Fetch-Mode":            []string{"navigate"},
				"Sec-Fetch-Dest":            []string{"document"},
				"Sec-Fetch-User":            []string{"?1"},
				"Upgrade-Insecure-Requests": []string{"1"},
				"User-Agent":                []string{"Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1"},
				"Referer":                   []string{targetLink},
			},
		)
		if err != nil {
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
