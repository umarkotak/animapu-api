package manga_scrapper

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
)

func GetMangaupdatesLatestManga(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	mangas := []models.Manga{}

	if queryParams.Page <= 0 {
		queryParams.Page = 1
	}

	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	idIncrement := 1
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
				ID:                  fmt.Sprintf("%v-%v", queryParams.Page, idIncrement),
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
				LatestCahpterNumber: latestChapter,
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
			idIncrement++
		}
	})

	err := c.Visit(fmt.Sprintf("https://www.mangaupdates.com/releases.html?page=%v", queryParams.Page))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return mangas, err
	}

	return mangas, nil
}

func GetMangaupdatesDetailManga(ctx context.Context, queryParams models.QueryParams) (models.Manga, error) {
	return models.Manga{}, nil
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
