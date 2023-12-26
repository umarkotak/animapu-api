package manga_scrapper_repository

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/config"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/utils/utils"
)

type AsuraNacm struct {
	Host string
}

func NewAsuraNacm() AsuraNacm {
	return AsuraNacm{
		Host: "https://asuratoon.com",
	}
}

func (t *AsuraNacm) GetHome(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(config.Get().CollyTimeout)

	mangas := []models.Manga{}

	c.OnHTML("#content > div.wrapper > div.postbody > div.bixbox > div.mrgn > div.listupd > div.bs", func(e *colly.HTMLElement) {
		latestChapterTitle := e.ChildText("div.bsx > a > div.bigor > div.adds > div.epxs")
		latestChapterTitle = utils.RemoveNonNumeric(latestChapterTitle)
		latestChapter, _ := strconv.ParseFloat(latestChapterTitle, 64)

		mangaLink := e.ChildAttr("div.bsx > a", "href")

		splitted := strings.Split(mangaLink, "/manga/")
		mangaID := splitted[1]
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

	err := c.Visit(fmt.Sprintf("%v/manga/?page=%v&order=update", t.Host, queryParams.Page))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
	}

	return mangas, nil
}

func (t *AsuraNacm) GetDetail(ctx context.Context, queryParams models.QueryParams) (models.Manga, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(config.Get().CollyTimeout)
	c.AllowURLRevisit = true

	manga := models.Manga{
		ID:          queryParams.SourceID,
		Source:      "asura_nacm",
		SourceID:    queryParams.SourceID,
		Title:       "Untitled",
		Description: "Description unavailable",
		Genres:      []string{},
		Status:      "Ongoing",
		CoverImages: []models.CoverImage{{Index: 0, ImageUrls: []string{""}}},
		Chapters:    []models.Chapter{},
	}

	c.OnHTML("div.bixbox.animefull > div.bigcontent.nobigcover > div.infox > h1", func(e *colly.HTMLElement) {
		if e.Text == "" {
			return
		}
		manga.Title = e.Text
	})

	c.OnHTML("div.bixbox.animefull > div.bigcontent > div.infox > h1", func(e *colly.HTMLElement) {
		if e.Text == "" {
			return
		}
		manga.Title = e.Text
	})

	c.OnHTML("div.entry-content.entry-content-single", func(e *colly.HTMLElement) {
		manga.Description = e.Text
	})

	c.OnHTML("div.bixbox.animefull > div.bigcontent.nobigcover > div.thumbook > div.thumb > img", func(e *colly.HTMLElement) {
		if e.Attr("src") == "" {
			return
		}
		manga.CoverImages = []models.CoverImage{{ImageUrls: []string{
			e.Attr("src"),
		}}}
	})

	c.OnHTML("div.bixbox.animefull > div.bigcontent > div.thumbook > div.thumb > img", func(e *colly.HTMLElement) {
		if e.Attr("src") == "" {
			return
		}
		manga.CoverImages = []models.CoverImage{{ImageUrls: []string{
			e.Attr("src"),
		}}}
	})

	c.OnHTML("#chapterlist > ul", func(e *colly.HTMLElement) {
		e.ForEach("li", func(i int, h *colly.HTMLElement) {
			chapterLink := h.ChildAttr("div > div > a", "href")
			chapterUrl, _ := url.Parse(chapterLink)
			chapterID := chapterUrl.Path
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

	targetUrl := fmt.Sprintf("%v/manga/%v", t.Host, queryParams.SourceID)
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

func (t *AsuraNacm) GetSearch(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(config.Get().CollyTimeout)

	mangas := []models.Manga{}

	c.OnHTML("#content > div.wrapper > div.postbody > div.bixbox > div.listupd > div.bs", func(e *colly.HTMLElement) {
		latestChapterTitle := e.ChildText("div.bsx > a > div.bigor > div.adds > div.epxs")
		latestChapterTitle = utils.RemoveNonNumeric(latestChapterTitle)
		latestChapter, _ := strconv.ParseFloat(latestChapterTitle, 64)

		mangaLink := e.ChildAttr("div.bsx > a", "href")

		splitted := strings.Split(mangaLink, "/manga/")
		mangaID := splitted[1]
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

	var err error

	pageCount := 5
	for i := 1; i <= pageCount; i++ {
		query := strings.Replace(queryParams.Title, " ", "+", -1)
		err = c.Visit(fmt.Sprintf("%v/page/%v/?s=%v", t.Host, i, query))
		if err != nil {
			logrus.WithContext(ctx).Error(err)
		}
	}
	if err != nil {
		return mangas, nil
	}

	return mangas, nil
}

func (t *AsuraNacm) GetChapter(ctx context.Context, queryParams models.QueryParams) (models.Chapter, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(config.Get().CollyTimeout)

	chapterNumber := float64(0)

	splitted := strings.Split(queryParams.ChapterID, "chapter-")
	if len(splitted) > 0 {
		newSplitted := strings.Split(splitted[len(splitted)-1], "-")
		if len(newSplitted) > 0 {
			chapterNumber = utils.ForceSanitizeStringToFloat(newSplitted[0])
		}
	}

	targetLink := fmt.Sprintf("%v/%v", t.Host, queryParams.ChapterID)

	chapter := models.Chapter{
		ID:            queryParams.ChapterID,
		SourceID:      queryParams.SourceID,
		Source:        "asura_nacm",
		SourceLink:    targetLink,
		Number:        chapterNumber,
		ChapterImages: []models.ChapterImage{},
	}

	c.OnHTML("img.ts-main-image", func(e *colly.HTMLElement) {
		chapter.ChapterImages = append(chapter.ChapterImages, models.ChapterImage{
			Index: 0,
			ImageUrls: []string{
				e.Attr("src"),
			},
		})
	})

	// Sample target: https://www.asurascans.com/?p=225072
	c.OnHTML("#comments > script", func(e *colly.HTMLElement) {
		pattern := `https:\/\/www\.asurascans\.com\/\?p=\d+`
		re := regexp.MustCompile(pattern)

		matches := re.FindAllString(e.Text, -1)

		asuraDisqusID := ""

		for _, match := range matches {
			asuraDisqusID = match
			break
		}

		asuraDisqusUrl, _ := url.Parse(asuraDisqusID)

		oneAsuraDisqusID := asuraDisqusUrl.Query().Get("p")

		// Disqus format: 242992 https://asura.nacm.xyz/?p=242992
		chapter.GenericDiscussion.DisqusID = fmt.Sprintf("%v %v", oneAsuraDisqusID, asuraDisqusID)
	})

	err := c.Visit(targetLink)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return chapter, err
	}

	return chapter, nil
}
