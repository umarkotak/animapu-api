package manga_scrapper_repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/config"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/utils/utils"
)

type Komikcast struct {
	Host    string
	ApiHost string
	Source  string
}

func NewKomikcast() Komikcast {
	return Komikcast{
		Host:    "https://komikcast.bz",
		ApiHost: "https://komikcast.bz",
		Source:  "komikcast",
	}
}

func (sc *Komikcast) GetHome(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(config.Get().CollyTimeout)

	mangas := []models.Manga{}

	c.OnHTML("#content > div > div > div.komiklist_filter > div > div > div.list-update_items-wrapper > div.list-update_item", func(e *colly.HTMLElement) {
		mangaID := e.ChildAttr("a", "href")
		mangaID = strings.ReplaceAll(mangaID, sc.Host, "")
		mangaID = strings.ReplaceAll(mangaID, "/komik/", "")
		mangaID = strings.TrimSuffix(mangaID, "/")

		mangas = append(mangas, models.Manga{
			ID:                  mangaID,
			SourceID:            mangaID,
			Source:              sc.Source,
			Title:               e.ChildText("a > div.list-update_item-info > h3"),
			Genres:              []string{},
			LatestChapterID:     "",
			LatestChapterNumber: utils.ForceSanitizeStringToFloat(strings.ReplaceAll(e.ChildText("a > div.list-update_item-info > div > div.chapter"), "Ch.", "")),
			LatestChapterTitle:  "",
			Chapters:            []models.Chapter{},
			CoverImages: []models.CoverImage{
				{
					Index: 1,
					ImageUrls: []string{
						e.ChildAttr("a > div.list-update_item-image > img", "src"),
					},
				},
			},
		})
	})

	err := c.Visit(fmt.Sprintf("%s/daftar-komik/page/%v/?orderby=update", sc.ApiHost, queryParams.Page))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
	}

	return mangas, nil
}

func (sc *Komikcast) GetDetail(ctx context.Context, queryParams models.QueryParams) (models.Manga, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(config.Get().CollyTimeout)
	c.AllowURLRevisit = true

	manga := models.Manga{
		ID:          queryParams.SourceID,
		Source:      sc.Source,
		SourceID:    queryParams.SourceID,
		Title:       "Untitled",
		Description: "Description unavailable",
		Genres:      []string{},
		Status:      "Ongoing",
		CoverImages: []models.CoverImage{{ImageUrls: []string{}}},
		Chapters:    []models.Chapter{},
	}

	c.OnHTML("h1.komik_info-content-body-title", func(e *colly.HTMLElement) {
		if e.Text == "" {
			return
		}
		manga.Title = e.Text
	})

	c.OnHTML("div.komik_info-description > div.komik_info-description-sinopsis > p", func(e *colly.HTMLElement) {
		if e.Text == "" {
			return
		}
		manga.Description = e.Text
	})

	c.OnHTML("div.komik_info-content > div.komik_info-content-thumbnail > img", func(e *colly.HTMLElement) {
		if e.Attr("src") == "" {
			return
		}
		manga.CoverImages = []models.CoverImage{{ImageUrls: []string{
			e.Attr("src"),
		}}}
	})

	c.OnHTML("#chapter-wrapper > li", func(e *colly.HTMLElement) {
		chapterLink := e.ChildAttr("a", "href")
		chapterID := chapterLink
		chapterID = strings.ReplaceAll(chapterID, sc.Host, "")
		chapterID = strings.ReplaceAll(chapterID, "/chapter/", "")
		chapterID = strings.TrimSuffix(chapterID, "/")

		if chapterLink == "" {
			return
		}

		manga.Chapters = append(manga.Chapters, models.Chapter{
			ID:       chapterID,
			Source:   sc.Source,
			SourceID: chapterID,
			Title:    e.ChildText("a"),
			Index:    utils.ForceSanitizeStringToInt64(e.ChildText("a")),
			Number:   utils.ForceSanitizeStringToFloat(e.ChildText("a")),
		})
	})

	targetUrl := fmt.Sprintf("%v/komik/%v", sc.Host, queryParams.SourceID)
	err := c.Visit(targetUrl)
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"target_url": targetUrl,
		}).Error(err)
		return manga, err
	}

	manga.GenerateLatestChapter()

	return manga, nil
}

func (sc *Komikcast) GetSearch(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(config.Get().CollyTimeout)

	mangas := []models.Manga{}

	c.OnHTML("div.list-update_item", func(e *colly.HTMLElement) {
		mangaID := e.ChildAttr("a", "href")
		mangaID = strings.ReplaceAll(mangaID, sc.Host, "")
		mangaID = strings.ReplaceAll(mangaID, "/komik/", "")
		mangaID = strings.TrimSuffix(mangaID, "/")

		mangas = append(mangas, models.Manga{
			ID:                  mangaID,
			SourceID:            mangaID,
			Source:              sc.Source,
			Title:               e.ChildText("a > div.list-update_item-info > h3"),
			Genres:              []string{},
			LatestChapterID:     "",
			LatestChapterNumber: utils.ForceSanitizeStringToFloat(strings.ReplaceAll(e.ChildText("a > div.list-update_item-info > div > div.chapter"), "Ch.", "")),
			LatestChapterTitle:  "",
			Chapters:            []models.Chapter{},
			CoverImages: []models.CoverImage{
				{
					Index: 1,
					ImageUrls: []string{
						e.ChildAttr("a > div.list-update_item-image > img", "src"),
					},
				},
			},
		})
	})

	q := strings.ReplaceAll(queryParams.Title, " ", "+")
	err := c.Visit(fmt.Sprintf("%s/?s=%s", sc.ApiHost, q))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
	}
	if err != nil {
		return mangas, err
	}

	return mangas, nil
}

func (sc *Komikcast) GetChapter(ctx context.Context, queryParams models.QueryParams) (models.Chapter, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(config.Get().CollyTimeout)

	targetLink := fmt.Sprintf("%v/chapter/%v", sc.Host, queryParams.ChapterID)

	chapter := models.Chapter{
		ID:            queryParams.ChapterID,
		SourceID:      queryParams.SourceID,
		Source:        sc.Source,
		SourceLink:    targetLink,
		Number:        0,
		ChapterImages: []models.ChapterImage{},
	}

	c.OnHTML("#content > div > div > div.chapter_headpost > h1", func(e *colly.HTMLElement) {
		chapter.Number = utils.ForceSanitizeStringToFloat(e.Text)
	})

	c.OnHTML("#chapter_body > div.main-reading-area > img", func(e *colly.HTMLElement) {
		chapterImage := models.ChapterImage{
			Index: 0,
			ImageUrls: []string{
				e.Attr("src"),
			},
		}

		chapter.ChapterImages = append(chapter.ChapterImages, chapterImage)
	})

	err := c.Visit(targetLink)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return chapter, err
	}

	return chapter, nil
}
