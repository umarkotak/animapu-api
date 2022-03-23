package manga_scrapper

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
)

// Used at home page
func GetMangaupdatesLatestManga(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	mangas := []models.Manga{}

	if queryParams.Page <= 0 {
		queryParams.Page = 1
	}

	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	var sourceDetailLink, title, sourceID, secondarySourceID string
	c.OnHTML("div.alt.p-1 div.row.no-gutters div", func(e *colly.HTMLElement) {
		if e.Attr("class") == "col-6 pbreak" {
			sourceDetailLink = e.ChildAttr("a", "href")
			sourceDetailLinkSplitted := strings.Split(sourceDetailLink, "series.html?id=")
			if len(sourceDetailLinkSplitted) >= 2 {
				sourceID = sourceDetailLinkSplitted[1]
			}
			title = e.ChildText("a")
			secondarySourceID = convertTitleToMangahubTitle(title)
		}

		if e.Attr("class") == "col-2 pl-1 pbreak" && sourceID != "" {
			latestChapterSplitted := strings.Split(e.Text, "c.")
			var latestChapterRaw string
			if len(latestChapterSplitted) > 0 {
				latestChapterRaw = latestChapterSplitted[len(latestChapterSplitted)-1]
			} else {
				latestChapterRaw = "0"
			}
			latestChapterSplitted = strings.Split(latestChapterRaw, "-")
			if len(latestChapterSplitted) > 0 {
				latestChapterRaw = latestChapterSplitted[0]
			}
			latestChapter, _ := strconv.ParseFloat(latestChapterRaw, 64)

			mangas = append(mangas, models.Manga{
				ID:                  sourceID,
				SourceID:            sourceID,
				SecondarySourceID:   secondarySourceID,
				Source:              "mangaupdates",
				SecondarySource:     "mangahub",
				Title:               title,
				Description:         "",
				Genres:              []string{},
				Status:              "",
				Rating:              "",
				LatestChapterID:     fmt.Sprintf("%v", latestChapter),
				LatestChapterNumber: latestChapter,
				LatestChapterTitle:  fmt.Sprintf("%v", latestChapter),
				Chapters:            []models.Chapter{},
				CoverImages: []models.CoverImage{
					{
						Index: 1,
						ImageUrls: []string{
							fmt.Sprintf("https://thumb.mghubcdn.com/mn/%s.jpg", secondarySourceID),
							fmt.Sprintf("https://thumb.mghubcdn.com/md/%s.jpg`", secondarySourceID),
							fmt.Sprintf("https://thumb.mghubcdn.com/m4l/%s.jpg`", secondarySourceID),
						},
					},
				},
			})
		}
	})

	err := c.Visit(fmt.Sprintf("https://www.mangaupdates.com/releases.html?page=%v", queryParams.Page))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return mangas, err
	}

	return mangas, nil
}

// Used at manga detail page
func GetMangaupdatesDetailManga(ctx context.Context, queryParams models.QueryParams) (models.Manga, error) {
	manga := models.Manga{
		Source:            "mangaupdates",
		SourceID:          queryParams.SourceID,
		SecondarySource:   "mangahub",
		SecondarySourceID: queryParams.SecondarySourceID,
		Status:            "Ongoing",
		Chapters:          []models.Chapter{},
	}
	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	c.OnHTML("#main_content > div:nth-child(2) > div.row.no-gutters > div.col-12.p-2 > span.releasestitle.tabletitle", func(e *colly.HTMLElement) {
		manga.Title = e.Text
	})
	c.OnHTML("#div_desc_link", func(e *colly.HTMLElement) {
		manga.Description = e.Text
	})
	c.OnHTML("#main_content > div:nth-child(2) > div.row.no-gutters > div:nth-child(4) > div:nth-child(5) > a", func(e *colly.HTMLElement) {
		manga.Genres = append(manga.Genres, e.Text)
	})
	c.OnHTML("#main_content > div:nth-child(2) > div.row.no-gutters > div:nth-child(3) > div:nth-child(20)", func(e *colly.HTMLElement) {
		if strings.Contains(e.Text, "Complete") {
			manga.Status = "Complete"
		}
	})
	c.OnHTML("#main_content > div:nth-child(2) > div.row.no-gutters > div:nth-child(4) > div:nth-child(2) > center > img", func(e *colly.HTMLElement) {
		manga.CoverImages = []models.CoverImage{
			{
				Index:     1,
				ImageUrls: []string{e.Attr("src")},
			},
		}
	})

	err := c.Visit(fmt.Sprintf("https://www.mangaupdates.com/series.html?id=%v", manga.SourceID))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return manga, err
	}

	cc := colly.NewCollector()
	cc.SetRequestTimeout(60 * time.Second)

	cc.OnHTML("body > div.body-site > div.container.container-main > div.container-main-left > div.panel-list-story > div:nth-child(1) > div > a:nth-child(2)", func(e *colly.HTMLElement) {
		manga.LatestChapterTitle = e.Text
		latestChapterSplitted := strings.Split(manga.LatestChapterTitle, " ")
		if len(latestChapterSplitted) > 0 {
			manga.LatestChapterID = strings.Replace(latestChapterSplitted[1], ":", "", -1)
			manga.LatestChapterNumber, _ = strconv.ParseFloat(manga.LatestChapterID, 64)
		}
	})

	err = cc.Visit(fmt.Sprintf("https://m.mangabat.com/search/manga/%v", manga.SecondarySourceID))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return manga, err
	}

	idx := int64(1)
	for i := int64(manga.LatestChapterNumber); i > 0; i-- {
		manga.Chapters = append(manga.Chapters, models.Chapter{
			ID:                fmt.Sprintf("%v", i),
			SourceID:          manga.SourceID,
			Source:            "mangaupdates",
			SecondarySourceID: manga.SecondarySourceID,
			SecondarySource:   "mangahub",
			Title:             fmt.Sprintf("Chapter %v", i),
			Index:             idx,
			Number:            float64(i),
		})
		idx++
	}

	return manga, nil
}

// Used at search manga page
func GetMangaSearchByQuery(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	mangas := []models.Manga{}

	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	c.OnHTML("#main_content > div.p-2.pt-2.pb-2.text > div:nth-child(2) > div", func(e *colly.HTMLElement) {
		title := e.ChildText("div > div.col.text.p-1.pl-3 > div > div:nth-child(1) > a > u > b")
		if title == "" {
			return
		}
		mangaupdatesDetailLink := e.ChildAttr("div.col-auto.align-self-center.series_thumb.p-0 > a", "href")
		mangaupdatesDetailLinkSplitted := strings.Split(mangaupdatesDetailLink, "series.html?id=")
		var mangaupdatesID string
		if len(mangaupdatesDetailLinkSplitted) >= 2 {
			mangaupdatesID = mangaupdatesDetailLinkSplitted[1]
		}

		mangas = append(mangas, models.Manga{
			ID:                  mangaupdatesID,
			SourceID:            mangaupdatesID,
			SecondarySourceID:   convertTitleToMangahubTitle(title),
			Source:              "mangaupdates",
			SecondarySource:     "mangahub",
			Title:               title,
			Description:         "",
			Genres:              []string{},
			Status:              "",
			Rating:              "",
			LatestChapterID:     "0",
			LatestChapterNumber: 0,
			LatestChapterTitle:  "Chapter 0",
			Chapters:            []models.Chapter{},
			CoverImages: []models.CoverImage{
				{
					Index: 1,
					ImageUrls: []string{
						e.ChildAttr("div.col-auto.align-self-center.series_thumb.p-0 > a > img", "src"),
					},
				},
			},
		})
	})

	err := c.Visit(fmt.Sprintf("https://www.mangaupdates.com/series.html?search=%v", url.QueryEscape(queryParams.Title)))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return mangas, err
	}

	return mangas, nil
}

// Used at reading manga chapter
func GetMangaupdatesDetailChapter(ctx context.Context, queryParams models.QueryParams) (models.Chapter, error) {
	pageCountConfig := int64(150)

	chapterNumber, _ := strconv.ParseFloat(queryParams.ChapterID, 64)

	chapter := models.Chapter{
		ID:                queryParams.ChapterID,
		SourceID:          queryParams.SourceID,
		Source:            "mangaupdates",
		SecondarySourceID: queryParams.SecondarySourceID,
		SecondarySource:   "mangahub",
		Title:             "",
		Index:             0,
		Number:            chapterNumber,
	}

	for i := int64(1); i <= pageCountConfig; i++ {
		chapter.ChapterImages = append(chapter.ChapterImages, models.ChapterImage{
			Index: i,
			ImageUrls: []string{
				fmt.Sprintf("https://img.mghubcdn.com/file/imghub/%v/%v/%v.jpg", queryParams.SecondarySourceID, chapterNumber, i),
				fmt.Sprintf("https://img.mghubcdn.com/file/imghub/%v/%v/%v.png", queryParams.SecondarySourceID, chapterNumber, i),
				fmt.Sprintf("https://img.mghubcdn.com/file/imghub/%v/%v/%v.jpeg", queryParams.SecondarySourceID, chapterNumber, i),
				fmt.Sprintf("https://img.mghubcdn.com/file/imghub/%v/%v/%v.webp", queryParams.SecondarySourceID, chapterNumber, i),
			},
		})
	}

	return chapter, nil
}

func convertTitleToMangahubTitle(initialTitle string) string {
	result := strings.ToLower(initialTitle)
	result = strings.Replace(result, "%", "", -1)
	result = strings.Replace(result, "'", "-", -1)
	result = strings.Replace(result, "!", "", -1)
	result = strings.Replace(result, "?", "", -1)
	result = strings.Replace(result, ".", "", -1)
	result = strings.Replace(result, "&", "", -1)
	result = strings.Replace(result, ":", "", -1)
	result = strings.Replace(result, ",", "", -1)
	result = strings.Replace(result, "(", "", -1)
	result = strings.Replace(result, ")", "", -1)
	result = strings.Replace(result, "-", "", -1)
	result = strings.Replace(result, "\"", "", -1)
	result = strings.Replace(result, "  ", "-", -1)
	result = strings.Replace(result, " ", "-", -1)
	return result
}
