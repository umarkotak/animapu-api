package manga_scrapper

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
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
		chapterID := e.ChildAttr("a", "href")
		chapterID = strings.Replace(chapterID, "https://www.webtoons.com/id/", "", -1)
		chapterID = strings.Replace(chapterID, "/", "Z2F", -1)
		chapterID = strings.Replace(chapterID, "?", "Z3F", -1)

		numberString := e.ChildText("span.subj")
		reg, _ := regexp.Compile("[^0-9]+")
		numberString = reg.ReplaceAllString(numberString, "")
		number, _ := strconv.ParseFloat(numberString, 64)

		manga.Chapters = append(manga.Chapters, models.Chapter{
			ID:       chapterID,
			Source:   "webtoonsid",
			SourceID: queryParams.SourceID,
			Title:    strings.Replace(e.ChildText("span.subj"), "UP", "", -1),
			Index:    idx,
			Number:   number,
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

	checkNextPage := len(manga.Chapters) > 0 && manga.Chapters[len(manga.Chapters)-1].Number/10 >= 1
	nextPage := int64(2)
	for checkNextPage {
		err := c.Visit(fmt.Sprintf("https://www.webtoons.com/id/%v&page=%v", formattedID, nextPage))
		if err != nil {
			logrus.WithContext(ctx).Error(err)
			return manga, err
		}
		nextPage += 1

		if manga.Chapters[len(manga.Chapters)-1].Number <= 1 {
			checkNextPage = false
		}

		if nextPage == 20 {
			checkNextPage = false
		}
	}

	return manga, nil
}

func GetWebtoonsidByQuery(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	mangas := []models.Manga{}

	err := c.Visit(fmt.Sprintf("https://www.webtoons.com/id/search?keyword=%v", queryParams.Title))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return mangas, err
	}

	return mangas, nil
}

func GetWebtoonsidDetailChapter(ctx context.Context, queryParams models.QueryParams) (models.Chapter, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	chapter := models.Chapter{
		ID:            queryParams.ChapterID,
		SourceID:      queryParams.SourceID,
		Source:        "webtoonsid",
		Title:         "",
		Index:         0,
		Number:        0,
		ChapterImages: []models.ChapterImage{},
	}

	idx := int64(1)
	c.OnHTML("div#container > div#content > div.cont_box > div.viewer_lst > div.viewer_img._img_viewer_area > img", func(e *colly.HTMLElement) {
		chapter.ChapterImages = append(chapter.ChapterImages, models.ChapterImage{
			Index: idx,
			ImageUrls: []string{
				fmt.Sprintf("http://localhost:6001/mangas/webtoons/image_proxy/%v", e.Attr("data-url")),
				fmt.Sprintf("https://animapu-api.herokuapp.com/mangas/webtoons/image_proxy/%v", e.Attr("data-url")),
			},
		})
		idx += 1
	})

	c.OnHTML("div.paginate.v2 > span.tx", func(e *colly.HTMLElement) {
		numberString := e.Text
		reg, _ := regexp.Compile("[^0-9]+")
		numberString = reg.ReplaceAllString(numberString, "")
		number, _ := strconv.ParseFloat(numberString, 64)
		chapter.Number = number
	})

	formattedID := queryParams.ChapterID
	formattedID = strings.Replace(formattedID, "Z2F", "/", -1)
	formattedID = strings.Replace(formattedID, "Z3F", "?", -1)
	err := c.Visit(fmt.Sprintf("https://www.webtoons.com/id/%v", formattedID))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return chapter, err
	}

	return chapter, nil
}
