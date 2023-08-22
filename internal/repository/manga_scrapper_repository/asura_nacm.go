package manga_scrapper_repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/utils/utils"
)

type AsuraNacm struct{}

func NewAsuraNacm() AsuraNacm {
	return AsuraNacm{}
}

func (t *AsuraNacm) GetHome(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	mangas := []models.Manga{}

	c.OnHTML("#content > div.wrapper > div.postbody > div.bixbox > div.mrgn > div.listupd > div.bs", func(e *colly.HTMLElement) {
		latestChapterTitle := e.ChildText("div.bsx > a > div.bigor > div.adds > div.epxs")
		latestChapterTitle = utils.RemoveNonNumeric(latestChapterTitle)
		latestChapter, _ := strconv.ParseFloat(latestChapterTitle, 64)

		mangaLink := e.ChildAttr("div.bsx > a", "href")

		mangaID := strings.ReplaceAll(mangaLink, "https://asura.nacm.xyz/manga/", "")
		mangaID = strings.ReplaceAll(mangaID, "/", "")

		mangas = append(mangas, models.Manga{
			ID:                  mangaID,
			SourceID:            mangaID,
			Source:              "asura_nacm",
			Title:               e.ChildText("div.bsx > a > div.bigor > div.tt"),
			Genres:              []string{},
			LatestChapterID:     "",
			LatestChapterNumber: latestChapter,
			LatestChapterTitle:  latestChapterTitle,
			Chapters:            []models.Chapter{},
			CoverImages: []models.CoverImage{
				{
					Index: 1,
					ImageUrls: []string{
						e.ChildAttr("div.bsx > a > div.limit > img.ts-post-image", "src"),
					},
				},
			},
		})
	})

	err := c.Visit(fmt.Sprintf("https://asura.nacm.xyz/manga/?page=%v&order=update", queryParams.Page))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
	}

	return mangas, nil
}

func (t *AsuraNacm) GetDetail(ctx context.Context, queryParams models.QueryParams) (models.Manga, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)
	c.AllowURLRevisit = true

	manga := models.Manga{
		ID:          queryParams.SourceID,
		Source:      "asura_nacm",
		SourceID:    queryParams.SourceID,
		Title:       "Untitled",
		Description: "Description unavailable",
		Genres:      []string{},
		Status:      "Ongoing",
		CoverImages: []models.CoverImage{{ImageUrls: []string{}}},
		Chapters:    []models.Chapter{},
	}

	c.OnHTML("article.hentry > div.bixbox.animefull > div.bigcontent.nobigcover > div.infox > h1", func(e *colly.HTMLElement) {
		manga.Title = e.Text
	})

	c.OnHTML("div.entry-content.entry-content-single", func(e *colly.HTMLElement) {
		manga.Description = e.Text
	})

	c.OnHTML("article.hentry > div.bixbox.animefull > div.bigcontent.nobigcover > div.thumbook > div.thumb > img", func(e *colly.HTMLElement) {
		manga.CoverImages = []models.CoverImage{{ImageUrls: []string{
			e.Attr("src"),
		}}}
	})

	c.OnHTML("#chapterlist > ul", func(e *colly.HTMLElement) {
		e.ForEach("li", func(i int, h *colly.HTMLElement) {
			chapterLink := h.ChildAttr("div > div > a", "href")
			chapterID := strings.ReplaceAll(chapterLink, "https://asura.nacm.xyz/", "")
			chapterID = strings.ReplaceAll(chapterID, "/", "")

			manga.Chapters = append(manga.Chapters, models.Chapter{
				ID:       chapterID,
				Source:   "asura_nacm",
				SourceID: chapterID,
				Title:    h.ChildText("div > div > a > span.chapternum"),
				Index:    int64(i),
				Number:   utils.ForceSanitizeStringToFloat(h.ChildText("div > div > a > span.chapternum")),
			})
		})
	})

	targetUrl := fmt.Sprintf("https://asura.nacm.xyz/manga/%v", queryParams.SourceID)
	err := c.Visit(targetUrl)
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"target_url": targetUrl,
		}).Error(err)
		return manga, err
	}

	return manga, nil
}

func (t *AsuraNacm) GetSearch(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	mangas := []models.Manga{}

	mangas = append(mangas, models.Manga{
		ID:                  "",
		SourceID:            "",
		Source:              "source",
		SecondarySourceID:   "",
		SecondarySource:     "secondary_source",
		Title:               "Untitled",
		Description:         "Description unavailable",
		Genres:              []string{},
		Status:              "Ongoing",
		Rating:              "10",
		LatestChapterID:     "chapter_id",
		LatestChapterNumber: 0,
		LatestChapterTitle:  "Chapter 0",
		Chapters:            []models.Chapter{},
		CoverImages: []models.CoverImage{
			{
				Index: 1,
				ImageUrls: []string{
					fmt.Sprintf("https://animapu-lite.vercel.app/images/manga/%v", "image_id"),
				},
			},
		},
	})

	err := c.Visit(fmt.Sprintf("https://animapu-lite.vercel.app/search/%v", queryParams.Page))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return mangas, err
	}

	return mangas, nil
}

func (t *AsuraNacm) GetChapter(ctx context.Context, queryParams models.QueryParams) (models.Chapter, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	chapter := models.Chapter{
		ID:            queryParams.ChapterID,
		SourceID:      queryParams.SourceID,
		Source:        "asura_nacm",
		Number:        0,
		ChapterImages: []models.ChapterImage{},
	}

	c.OnHTML("#readerarea > p", func(e *colly.HTMLElement) {
		chapter.ChapterImages = append(chapter.ChapterImages, models.ChapterImage{
			Index: 0,
			ImageUrls: []string{
				e.ChildAttr("img", "src"),
			},
		})
	})

	err := c.Visit(fmt.Sprintf("https://asura.nacm.xyz/%v", queryParams.ChapterID))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return chapter, err
	}

	return chapter, nil
}
