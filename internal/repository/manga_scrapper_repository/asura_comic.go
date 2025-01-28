package manga_scrapper_repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-rod/rod"
	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/config"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/utils/utils"
)

type AsuraComic struct {
	Host string
}

func NewAsuraComic() AsuraComic {
	return AsuraComic{
		Host: "https://asuracomic.net",
	}
}

func (t *AsuraComic) GetHome(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(config.Get().CollyTimeout)

	mangas := []models.Manga{}

	c.OnHTML(`body > div:nth-child(4) > div > div > div > div.w-\[100\%\].float-left.min-\[882px\]\:w-\[68\.5\%\].min-\[1030px\]\:w-\[70\%\].max-\[600px\]\:w-\[100\%\] > div > div.grid.grid-cols-2.sm\:grid-cols-2.md\:grid-cols-5.gap-3.p-4 > a`, func(e *colly.HTMLElement) {
		latestChapterTitle := e.ChildText(`div > div > div.block.w-\[100\%\].h-auto.items-center > span.text-\[13px\].text-\[\#999\]`)
		latestChapterTitle = utils.RemoveNonNumeric(latestChapterTitle)
		latestChapter, _ := strconv.ParseFloat(latestChapterTitle, 64)

		mangaLink := e.Attr("href")
		// logrus.WithContext(ctx).Infof("manga link: %v", mangaLink)

		splitted := strings.Split(mangaLink, "/")
		mangaID := splitted[1]
		mangaID = strings.ReplaceAll(mangaID, "/", "")

		mangas = append(mangas, models.Manga{
			ID:                  mangaID,
			SourceID:            mangaID,
			Source:              "asura_nacm",
			Title:               e.ChildText(`div > div > div.block.w-\[100\%\].h-auto.items-center > span.block.text-\[13\.3px\].font-bold`),
			Genres:              []string{},
			LatestChapterID:     "",
			LatestChapterNumber: latestChapter,
			LatestChapterTitle:  latestChapterTitle,
			Chapters:            []models.Chapter{},
			CoverImages: []models.CoverImage{
				{
					Index: 1,
					ImageUrls: []string{
						e.ChildAttr(`div > div > div.flex.h-\[250px\].md\:h-\[200px\].overflow-hidden.relative.hover\:opacity-60 > img`, "src"),
					},
				},
			},
		})
	})

	err := c.Visit(fmt.Sprintf("%v/series?page=%v&order=update", t.Host, queryParams.Page))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
	}

	return mangas, nil
}

func (t *AsuraComic) GetDetail(ctx context.Context, queryParams models.QueryParams) (models.Manga, error) {
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

	c.OnHTML(`body > div:nth-child(4) > div > div > div > div.w-\[100\%\].float-left.min-\[882px\]\:w-\[68\.5\%\].min-\[1030px\]\:w-\[70\%\].max-\[600px\]\:w-\[100\%\] > div > div.space-y-7 > div.space-y-4 > div:nth-child(2) > div.relative.z-10.grid.grid-cols-12.gap-4.pt-4.pl-4.pr-4.pb-12 > div.col-span-12.sm\:col-span-9 > div.text-center.sm\:text-left > span`, func(e *colly.HTMLElement) {
		if e.Text == "" {
			return
		}
		manga.Title = e.Text
	})

	c.OnHTML(`body > div:nth-child(4) > div > div > div > div.w-\[100\%\].float-left.min-\[882px\]\:w-\[68\.5\%\].min-\[1030px\]\:w-\[70\%\].max-\[600px\]\:w-\[100\%\] > div > div.space-y-7 > div.space-y-4 > div:nth-child(2) > div.relative.z-10.grid.grid-cols-12.gap-4.pt-4.pl-4.pr-4.pb-12 > div.col-span-12.sm\:col-span-9 > span`, func(e *colly.HTMLElement) {
		manga.Description = e.Text
	})

	c.OnHTML(`body > div:nth-child(4) > div > div > div > div.w-\[100\%\].float-left.min-\[882px\]\:w-\[68\.5\%\].min-\[1030px\]\:w-\[70\%\].max-\[600px\]\:w-\[100\%\] > div > div.space-y-7 > div.space-y-4 > div:nth-child(2) > div.relative.z-10.grid.grid-cols-12.gap-4.pt-4.pl-4.pr-4.pb-12 > div.relative.col-span-12.sm\:col-span-3.space-y-3.px-6.sm\:px-0 > img`, func(e *colly.HTMLElement) {
		if e.Attr("src") == "" {
			return
		}
		manga.CoverImages = []models.CoverImage{{ImageUrls: []string{
			e.Attr("src"),
		}}}
	})

	c.OnHTML(`div.bg-\[\#222222\] > div > div > div`, func(e *colly.HTMLElement) {
		chapterLink := e.ChildAttr(`a`, "href")

		if chapterLink == "" {
			return
		}

		chapterLinkSplitted := strings.Split(chapterLink, "/")
		chapterID := chapterLinkSplitted[len(chapterLinkSplitted)-1]

		chapter := models.Chapter{
			ID:       chapterID,
			Source:   "asura_nacm",
			SourceID: chapterID,
			Title:    e.ChildText("a > h3.text-sm.text-white.font-medium.flex.flex-row"),
			Index:    int64(0),
			Number:   utils.ForceSanitizeStringToFloat(e.ChildText("a > h3.text-sm.text-white.font-medium.flex.flex-row")),
		}

		if chapter.Title == "" {
			return
		}

		manga.Chapters = append(manga.Chapters, chapter)
	})

	targetUrl := fmt.Sprintf("%v/series/%v", t.Host, queryParams.SourceID)
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

func (t *AsuraComic) GetSearch(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(config.Get().CollyTimeout)

	mangas := []models.Manga{}

	c.OnHTML(`body > div:nth-child(4) > div > div > div > div.w-\[100\%\].float-left.min-\[882px\]\:w-\[68\.5\%\].min-\[1030px\]\:w-\[70\%\].max-\[600px\]\:w-\[100\%\] > div > div.grid.grid-cols-2.sm\:grid-cols-2.md\:grid-cols-5.gap-3.p-4 > a`, func(e *colly.HTMLElement) {
		latestChapterTitle := e.ChildText(`div > div > div.block.w-\[100\%\].h-auto.items-center > span.text-\[13px\].text-\[\#999\]`)
		latestChapterTitle = utils.RemoveNonNumeric(latestChapterTitle)
		latestChapter, _ := strconv.ParseFloat(latestChapterTitle, 64)

		mangaLink := e.Attr("href")
		// logrus.WithContext(ctx).Infof("manga link: %v", mangaLink)

		splitted := strings.Split(mangaLink, "/")
		mangaID := splitted[1]
		mangaID = strings.ReplaceAll(mangaID, "/", "")

		mangas = append(mangas, models.Manga{
			ID:                  mangaID,
			SourceID:            mangaID,
			Source:              "asura_nacm",
			Title:               e.ChildText(`div > div > div.block.w-\[100\%\].h-auto.items-center > span.block.text-\[13\.3px\].font-bold`),
			Genres:              []string{},
			LatestChapterID:     "",
			LatestChapterNumber: latestChapter,
			LatestChapterTitle:  latestChapterTitle,
			Chapters:            []models.Chapter{},
			CoverImages: []models.CoverImage{
				{
					Index: 1,
					ImageUrls: []string{
						e.ChildAttr(`div > div > div.flex.h-\[250px\].md\:h-\[200px\].overflow-hidden.relative.hover\:opacity-60 > img`, "src"),
					},
				},
			},
		})
	})

	pageCount := 5
	for i := 1; i <= pageCount; i++ {
		err := c.Visit(fmt.Sprintf("%v/series?page=%v&name=%v&order=update", t.Host, i, strings.ReplaceAll(queryParams.Title, " ", "%20")))
		if err != nil {
			logrus.WithContext(ctx).Error(err)
		}
		c.Wait()

		logrus.Infof("LEN MANGAS: %v", len(mangas))
		if len(mangas)%15 != 0 {
			break
		}
	}

	return mangas, nil
}

func (t *AsuraComic) GetChapter(ctx context.Context, queryParams models.QueryParams) (models.Chapter, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(config.Get().CollyTimeout)

	chapterNumber := utils.ForceSanitizeStringToFloat(queryParams.ChapterID)

	targetLink := fmt.Sprintf("%v/series/%v/chapter/%v", t.Host, queryParams.SourceID, queryParams.ChapterID)

	chapter := models.Chapter{
		ID:            queryParams.ChapterID,
		SourceID:      queryParams.SourceID,
		Source:        "asura_nacm",
		SourceLink:    targetLink,
		Number:        chapterNumber,
		ChapterImages: []models.ChapterImage{},
	}

	c.OnHTML(`body > div > div > div > div > div > div > div`, func(e *colly.HTMLElement) {
		chapter.ChapterImages = append(chapter.ChapterImages, models.ChapterImage{
			Index: 0,
			ImageUrls: []string{
				"https://gg.asuracomic.net/storage/media/267218/conversions/00-kopya-optimized.webp",
			},
		})
	})

	// err := c.Visit(targetLink)
	// if err != nil {
	// 	logrus.WithContext(ctx).Error(err)
	// 	return chapter, err
	// }

	page := rod.New().MustConnect().MustPage(targetLink)

	el := page.MustElement("div.w-full.mx-auto.center > img")
	el.Attribute("src")
	// rod_utils.OutputFile(fmt.Sprintf("hello-%v.png", 9000), el.MustResource())

	els := page.MustElements("div.w-full.mx-auto.center > img")

	for _, el := range els {
		// rod_utils.OutputFile(fmt.Sprintf("hello-%v.png", i), el.MustResource())

		res, err := el.Attribute("src")

		if err != nil {
			logrus.WithContext(ctx).Error(err)
			continue
		}

		empty := ""
		if res == nil || res == &empty {
			continue
		}

		chapter.ChapterImages = append(chapter.ChapterImages, models.ChapterImage{
			Index: 0,
			ImageUrls: []string{
				*res,
			},
		})
	}

	return chapter, nil
}
