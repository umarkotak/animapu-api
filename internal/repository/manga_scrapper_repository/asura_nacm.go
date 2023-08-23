package manga_scrapper_repository

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
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

	c.OnHTML("#content > div.wrapper > div.postbody > div.bixbox > div.listupd > div.bs", func(e *colly.HTMLElement) {
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

	var err error

	pageCount := 3
	for i := 1; i <= pageCount; i++ {
		query := strings.Replace(queryParams.Title, " ", "+", -1)
		err = c.Visit(fmt.Sprintf("https://asura.nacm.xyz/page/%v/?s=%v", i, query))
		if err != nil {
			logrus.WithContext(ctx).Error(err)
		}
	}
	if err != nil {
		return mangas, err
	}

	return mangas, nil
}

func (t *AsuraNacm) GetChapter(ctx context.Context, queryParams models.QueryParams) (models.Chapter, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	chapterNumber := float64(0)

	splitted := strings.Split(queryParams.ChapterID, "chapter-")
	if len(splitted) > 0 {
		newSplitted := strings.Split(splitted[len(splitted)-1], "-")
		if len(newSplitted) > 0 {
			chapterNumber = utils.ForceSanitizeStringToFloat(newSplitted[0])
		}
	}

	chapter := models.Chapter{
		ID:            queryParams.ChapterID,
		SourceID:      queryParams.SourceID,
		Source:        "asura_nacm",
		Number:        chapterNumber,
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

	err := c.Visit(fmt.Sprintf("https://asura.nacm.xyz/%v", queryParams.ChapterID))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return chapter, err
	}

	return chapter, nil
}
