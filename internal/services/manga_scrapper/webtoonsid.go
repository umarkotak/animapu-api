package manga_scrapper

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
)

// This is only a template file to be easily copy-pasted

func GetWebtoonsidLatestManga(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	mangas := []models.Manga{}

	c.OnHTML("html body div#wrap div#container div#content div div#dailyList li", func(e *colly.HTMLElement) {
		sourceID := e.ChildAttr("a", "href")
		sourceID = strings.Replace(sourceID, "https://www.webtoons.com/id/", "", -1)
		sourceID = strings.Replace(sourceID, "/", "Z2F", -1)
		sourceID = strings.Replace(sourceID, "?", "Z3F", -1)

		mangas = append(mangas, models.Manga{
			ID:                  sourceID,
			Source:              "webtoonsid",
			SourceID:            sourceID,
			Title:               e.ChildText("a > div > p.subj"),
			Description:         "",
			Genres:              []string{},
			Status:              "",
			Rating:              "",
			LatestChapterID:     "",
			LatestChapterNumber: 0,
			LatestChapterTitle:  "",
			Chapters:            []models.Chapter{},
			CoverImages: []models.CoverImage{
				{
					Index: 1,
					ImageUrls: []string{
						fmt.Sprintf("http://localhost:6001/mangas/webtoons/image_proxy/%v", e.ChildAttr("img", "src")),
						fmt.Sprintf("https://animapu-api.herokuapp.com/mangas/webtoons/image_proxy/%v", e.ChildAttr("img", "src")),
					},
				},
			},
		})
	})

	err := c.Visit(fmt.Sprintf("https://www.webtoons.com/id/dailySchedule"))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return mangas, err
	}

	return mangas, nil
}

func GetWebtoonsidDetailManga(ctx context.Context, queryParams models.QueryParams) (models.Manga, error) {
	fmt.Println("hello")
	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	manga := models.Manga{
		Source:           "webtoonsid",
		SourceID:         queryParams.SourceID,
		Status:           "Ongoing",
		Chapters:         []models.Chapter{},
		Description:      "Description unavailable",
		CoverImages:      []models.CoverImage{{ImageUrls: []string{""}}},
		ChapterPaginated: true,
	}

	c.OnHTML("span.thmb > img", func(e *colly.HTMLElement) {
		manga.CoverImages = []models.CoverImage{
			{
				Index:     1,
				ImageUrls: []string{fmt.Sprintf("http://localhost:6001/mangas/webtoons/image_proxy/%v", e.Attr("src"))},
			},
		}
	})

	c.OnHTML("#content > div.cont_box > div.detail_header.type_white > div.info > h1", func(e *colly.HTMLElement) {
		manga.Title = e.Text
	})

	idx := int64(1)
	c.OnHTML("div.detail_body > div.detail_lst > ul#_listUl > li", func(e *colly.HTMLElement) {
		manga.Chapters = append(manga.Chapters, models.Chapter{
			ID:       "id",
			Source:   "webtoonsid",
			SourceID: "id",
			Title:    e.ChildText("span.subj"),
			Index:    idx,
			Number:   0,
		})

		idx += 1
	})

	formattedID := queryParams.SourceID
	formattedID = strings.Replace(formattedID, "Z2F", "/", -1)
	formattedID = strings.Replace(formattedID, "Z3F", "?", -1)
	err := c.Visit(fmt.Sprintf("https://www.webtoons.com/id/%v", formattedID))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return manga, err
	}

	return manga, nil
}

func GetWebtoonsidByQuery(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	mangas := []models.Manga{}

	err := c.Visit(fmt.Sprintf("https://m.mangabat.com/manga-list-all/%v", queryParams.Page))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return mangas, err
	}

	return mangas, nil
}

func GetWebtoonsidDetailChapter(ctx context.Context, queryParams models.QueryParams) (models.Chapter, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	chapter := models.Chapter{}

	err := c.Visit(fmt.Sprintf("https://m.mangabat.com/manga-list-all/%v", queryParams.Page))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return chapter, err
	}

	return chapter, nil
}
