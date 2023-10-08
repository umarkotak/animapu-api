package manga_scrapper_repository

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/config"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/utils/utils"
)

type Komikindo struct {
	Host   string
	Source string
}

func NewKomikindo() Komikindo {
	return Komikindo{
		Host:   "https://komikindo.tv",
		Source: "komikindo",
	}
}

func (sc *Komikindo) GetHome(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	mangas := []models.Manga{}

	c.OnHTML("#content > div.postbody > section.whites > div.widget-body > div > div > div.listupd > div.animepost", func(e *colly.HTMLElement) {
		latestChapterTitle := e.ChildText("div.animposx > div.bigor > div > div > a")
		latestChapterTitle = utils.RemoveNonNumeric(strings.ReplaceAll(latestChapterTitle, "Ch.", ""))
		latestChapter, _ := strconv.ParseFloat(latestChapterTitle, 64)

		mangaLink := e.ChildAttr("div.animposx > div.bigor > a", "href")

		splitted := strings.Split(mangaLink, "/komik/")
		mangaID := splitted[1]
		mangaID = strings.ReplaceAll(mangaID, "/", "")

		mangas = append(mangas, models.Manga{
			ID:                  mangaID,
			SourceID:            mangaID,
			Source:              sc.Source,
			Title:               e.ChildText("div.animposx > div.bigor > a > div > h4"),
			Genres:              []string{},
			LatestChapterID:     "",
			LatestChapterNumber: latestChapter,
			LatestChapterTitle:  latestChapterTitle,
			Chapters:            []models.Chapter{},
			CoverImages: []models.CoverImage{
				{
					Index: 1,
					ImageUrls: []string{
						fmt.Sprintf("%v/mangas/komikindo/image_proxy/%v", config.Get().AnimapuOnlineHost, e.ChildAttr("a > div > img", "src")),
					},
				},
			},
		})
	})

	err := c.Visit(fmt.Sprintf("%v/komik-terbaru/page/%v", sc.Host, queryParams.Page))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
	}

	return mangas, nil
}

func (sc *Komikindo) GetDetail(ctx context.Context, queryParams models.QueryParams) (models.Manga, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)
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

	c.OnHTML("div.infoanime > h1.entry-title", func(e *colly.HTMLElement) {
		if e.Text == "" {
			return
		}
		manga.Title = e.Text
	})

	c.OnHTML("div.entry-content.entry-content-single", func(e *colly.HTMLElement) {
		manga.Description = e.Text
	})

	c.OnHTML("div.thumb > img", func(e *colly.HTMLElement) {
		if e.Attr("src") == "" {
			return
		}
		manga.CoverImages = []models.CoverImage{{ImageUrls: []string{
			fmt.Sprintf("%v/mangas/komikindo/image_proxy/%v", config.Get().AnimapuOnlineHost, e.Attr("src")),
		}}}
	})

	c.OnHTML("#chapter_list > ul", func(e *colly.HTMLElement) {
		e.ForEach("li", func(i int, h *colly.HTMLElement) {
			chapterLink := h.ChildAttr("span.lchx > a", "href")
			chapterUrl, _ := url.Parse(chapterLink)
			chapterID := chapterUrl.Path
			chapterID = strings.ReplaceAll(chapterID, "/", "")

			manga.Chapters = append(manga.Chapters, models.Chapter{
				ID:       chapterID,
				Source:   sc.Source,
				SourceID: chapterID,
				Title:    h.ChildText("span.lchx > a"),
				Index:    int64(i),
				Number:   utils.ForceSanitizeStringToFloat(h.ChildText("span.lchx > a")),
			})
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

func (sc *Komikindo) GetSearch(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	mangas := []models.Manga{}

	c.OnHTML("#content > div.postbody > section > div.film-list > div.animepost", func(e *colly.HTMLElement) {
		mangaLink := e.ChildAttr("div.animposx > a", "href")

		splitted := strings.Split(mangaLink, "/komik/")
		mangaID := splitted[1]
		mangaID = strings.ReplaceAll(mangaID, "/", "")

		mangas = append(mangas, models.Manga{
			ID:                  mangaID,
			SourceID:            mangaID,
			Source:              sc.Source,
			Title:               e.ChildText("div.animposx > div.bigors > a > div > h4"),
			Genres:              []string{},
			LatestChapterID:     "",
			LatestChapterNumber: 0,
			LatestChapterTitle:  "",
			Chapters:            []models.Chapter{},
			CoverImages: []models.CoverImage{
				{
					Index: 1,
					ImageUrls: []string{
						fmt.Sprintf("%v/mangas/komikindo/image_proxy/%v", config.Get().AnimapuOnlineHost, e.ChildAttr("a > div > img", "src")),
					},
				},
			},
		})
	})

	err := c.Visit(fmt.Sprintf("%v/?s=%v", sc.Host, queryParams.Title))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
	}
	if err != nil {
		return mangas, err
	}

	return mangas, nil
}

func (sc *Komikindo) GetChapter(ctx context.Context, queryParams models.QueryParams) (models.Chapter, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	chapterNumber := float64(0)

	splitted := strings.Split(queryParams.ChapterID, "chapter-")
	if len(splitted) > 0 {
		chapterNumber = utils.ForceSanitizeStringToFloat(splitted[len(splitted)-1])
	}

	targetLink := fmt.Sprintf("%v/%v", sc.Host, queryParams.ChapterID)

	chapter := models.Chapter{
		ID:            queryParams.ChapterID,
		SourceID:      queryParams.SourceID,
		Source:        sc.Source,
		SourceLink:    targetLink,
		Number:        chapterNumber,
		ChapterImages: []models.ChapterImage{},
	}

	c.OnHTML("#chimg-auh", func(e *colly.HTMLElement) {
		e.ForEach("img", func(i int, h *colly.HTMLElement) {
			chapter.ChapterImages = append(chapter.ChapterImages, models.ChapterImage{
				Index: 0,
				ImageUrls: []string{
					fmt.Sprintf("%v/mangas/komikindo/image_proxy/%v", config.Get().AnimapuOnlineHost, h.Attr("src")),
				},
			})
		})
	})

	// TODO: Adjust discus
	// // Sample target: https://www.asurascans.com/?p=225072
	// c.OnHTML("#disqus_embed-js-extra", func(e *colly.HTMLElement) {
	// 	pattern := `https:\/\/www\.komikindo\.one\/\?p=\d+`
	// 	re := regexp.MustCompile(pattern)

	// 	matches := re.FindAllString(e.Text, -1)

	// 	asuraDisqusID := ""

	// 	for _, match := range matches {
	// 		asuraDisqusID = match
	// 		break
	// 	}

	// 	asuraDisqusUrl, _ := url.Parse(asuraDisqusID)

	// 	oneAsuraDisqusID := asuraDisqusUrl.Query().Get("p")

	// 	// Disqus format: 242992 https://asura.nacm.xyz/?p=242992
	// 	chapter.GenericDiscussion.DisqusID = fmt.Sprintf("%v %v", oneAsuraDisqusID, asuraDisqusID)
	// })

	err := c.Visit(targetLink)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return chapter, err
	}

	return chapter, nil
}
