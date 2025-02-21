package manga_scrapper_repository

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/config"
	"github.com/umarkotak/animapu-api/internal/contract"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/utils/utils"
)

type (
	WeebCentralManga struct {
		SeriesID   string    `json:"SeriesID"`
		IndexName  string    `json:"IndexName"`
		SeriesName string    `json:"SeriesName"`
		ScanStatus string    `json:"ScanStatus"`
		Chapter    string    `json:"Chapter"`
		Genres     string    `json:"Genres"`
		Date       time.Time `json:"Date"`
		IsEdd      bool      `json:"IsEdd"`
	}

	WeebCentralChapter struct {
		Chapter     string      `json:"Chapter"`
		Type        string      `json:"Type"`
		Date        string      `json:"Date"`
		ChapterName interface{} `json:"ChapterName"`
	}

	WeebCentralSearchManga struct {
		IndexName  string   `json:"i"`
		SeriesName string   `json:"s"`
		AltNames   []string `json:"a"`
	}
)

type WeebCentral struct {
	Host    string
	Source  string
	ImgHost string
}

func NewWeebCentral() WeebCentral {
	return WeebCentral{
		Source:  "weeb_central",
		Host:    "https://weebcentral.com",
		ImgHost: "https://temp.compsci88.com",
	}
}

func (sc *WeebCentral) GetHome(ctx context.Context, queryParams models.QueryParams) ([]contract.Manga, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(config.Get().CollyTimeout)

	mangas := []contract.Manga{}

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	})

	c.OnHTML("body > abbr > article", func(e *colly.HTMLElement) {
		mangaLink := e.ChildAttr("a.aspect-square.overflow-hidden", "href")
		mangaID := strings.ReplaceAll(mangaLink, sc.Host, "")
		mangaID = strings.TrimPrefix(mangaID, "/series/")
		mangaID = strings.ReplaceAll(mangaID, "/", "---")

		mangas = append(mangas, contract.Manga{
			ID:                  mangaID,
			Source:              sc.Source,
			SourceID:            mangaID,
			Title:               e.ChildText("a.min-w-0.flex.flex-col.justify-center.pe-4 > div:nth-child(1) > div"),
			Genres:              []string{},
			LatestChapterID:     "",
			LatestChapterNumber: utils.ForceSanitizeStringToFloat(e.ChildText("a.min-w-0.flex.flex-col.justify-center.pe-4 > div:nth-child(2) > span")),
			LatestChapterTitle:  e.ChildText("a.min-w-0.flex.flex-col.justify-center.pe-4 > div:nth-child(2) > span"),
			Chapters:            []contract.Chapter{},
			CoverImages: []contract.CoverImage{
				{
					Index: 1,
					ImageUrls: []string{
						e.ChildAttr("a.aspect-square.overflow-hidden > picture > source", "srcset"),
						e.ChildAttr("a.aspect-square.overflow-hidden > picture > img", "src"),
					},
				},
			},
		})
	})

	targetLinks := []string{
		fmt.Sprintf("%v/latest-updates/%v", sc.Host, queryParams.Page),
	}
	for _, targetLink := range targetLinks {
		err := c.Visit(targetLink)
		c.Wait()
		if err != nil || len(mangas) <= 0 {
			logrus.WithContext(ctx).Error(err)
			continue
		}
		break
	}

	return mangas, nil
}

func (sc *WeebCentral) GetDetail(ctx context.Context, queryParams models.QueryParams) (contract.Manga, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(config.Get().CollyTimeout)
	c.AllowURLRevisit = true

	newMangaIDSplit := strings.Split(queryParams.SourceID, "---")
	newMangaID := newMangaIDSplit[0]

	manga := contract.Manga{
		ID:          queryParams.SourceID,
		Source:      sc.Source,
		SourceID:    queryParams.SourceID,
		Title:       "",
		Description: "Description unavailable",
		Genres:      []string{},
		Status:      "Ongoing",
		CoverImages: []contract.CoverImage{{ImageUrls: []string{
			fmt.Sprintf("https://temp.compsci88.com/cover/normal/%s.webp", newMangaID),
		}}},
		Chapters: []contract.Chapter{},
	}

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	})

	c.OnHTML("#top > section.flex.flex-col > section.flex.flex-col.gap-4 > h1", func(e *colly.HTMLElement) {
		manga.Title = e.Text
	})

	c.OnHTML("#top > section.flex.flex-col > section.flex.flex-col.gap-4 > section:nth-child(3) > ul > li > p", func(e *colly.HTMLElement) {
		manga.Description = e.Text
	})

	c.OnHTML("#top > section.flex.flex-col.md:flex-row.gap-4.md:gap-8 > section.md:w-4/12.flex.flex-col.gap-4 > section:nth-child(3) > picture > img", func(e *colly.HTMLElement) {
		manga.CoverImages = []contract.CoverImage{{ImageUrls: []string{
			e.Attr("src"),
		}}}
	})

	c.OnHTML("body > div > a", func(e *colly.HTMLElement) {
		chapterUrl := e.Attr("href")
		chapterID := strings.ReplaceAll(chapterUrl, sc.Host, "")
		chapterID = strings.TrimPrefix(chapterID, "/chapters/")

		chText := e.ChildText("span.grow.flex.items-center.gap-2 > span:nth-child(1)")
		chNumer := utils.ForceSanitizeStringToFloat(chText)

		manga.Chapters = append(manga.Chapters, contract.Chapter{
			ID:                chapterID,
			Source:            sc.Source,
			SourceID:          chapterID,
			SecondarySourceID: "",
			Title:             chText,
			Index:             int64(chNumer),
			Number:            chNumer,
		})
	})

	queryParams.SourceID = strings.ReplaceAll(queryParams.SourceID, "---", "/")
	targetLinks := []string{
		fmt.Sprintf("%v/series/%v", sc.Host, queryParams.SourceID),
		fmt.Sprintf("%v/series/%v/full-chapter-list", sc.Host, newMangaID),
	}
	for _, targetLink := range targetLinks {
		c.Visit(targetLink)
	}

	manga.GenerateLatestChapter()

	return manga, nil
}

func (sc *WeebCentral) GetSearch(ctx context.Context, queryParams models.QueryParams) ([]contract.Manga, error) {
	mangas := []contract.Manga{}

	// Create a new collector
	c := colly.NewCollector()

	// Set up the headers
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("HX-Trigger", "quick-search-input")
		r.Headers.Set("HX-Trigger-Name", "text")
		r.Headers.Set("sec-ch-ua-platform", "macOS")
		r.Headers.Set("Referer", "https://weebcentral.com/")
		r.Headers.Set("HX-Target", "quick-search-result")
		r.Headers.Set("HX-Current-URL", "https://weebcentral.com/")
		r.Headers.Set("sec-ch-ua", `"Not A(Brand";v="8", "Chromium";v="132", "Google Chrome";v="132"`)
		r.Headers.Set("sec-ch-ua-mobile", "?0")
		r.Headers.Set("HX-Request", "true")
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36")
		r.Headers.Set("Content-Type", "application/x-www-form-urlencoded")
	})

	c.OnHTML("#quick-search-result > div.w-full.join.join-vertical.flex.absolute.inset-x-0.z-10.mt-4.rounded-none > a", func(e *colly.HTMLElement) {
		mangaLink := e.Attr("href")
		mangaID := strings.ReplaceAll(mangaLink, sc.Host, "")
		mangaID = strings.TrimPrefix(mangaID, "/series/")
		mangaID = strings.ReplaceAll(mangaID, "/", "---")

		mangas = append(mangas, contract.Manga{
			ID:                  mangaID,
			Source:              sc.Source,
			SourceID:            mangaID,
			Title:               e.ChildText("div.flex-1.overflow-hidden.text-left.text-ellipsis.leading-normal.line-clamp-2"),
			Genres:              []string{},
			LatestChapterID:     "",
			LatestChapterNumber: 0,
			LatestChapterTitle:  "",
			Chapters:            []contract.Chapter{},
			CoverImages: []contract.CoverImage{
				{
					Index: 1,
					ImageUrls: []string{
						e.ChildAttr("div > picture > source", "srcset"),
						e.ChildAttr("div > picture > img", "src"),
					},
				},
			},
		})
	})

	// Prepare the form data
	formData := map[string]string{
		"text": queryParams.Title,
	}

	// Send the POST request
	err := c.Post(fmt.Sprintf("%v/search/simple?location=main", sc.Host), formData)
	if err != nil {
		log.Fatal(err)
	}

	return mangas, nil
}

func (sc *WeebCentral) GetChapter(ctx context.Context, queryParams models.QueryParams) (contract.Chapter, error) {
	var err error
	c := colly.NewCollector()
	c.SetRequestTimeout(10 * time.Minute)

	chapter := contract.Chapter{
		ID:            queryParams.ChapterID,
		SourceID:      queryParams.SourceID,
		Source:        sc.Source,
		Number:        0,
		ChapterImages: []contract.ChapterImage{},
		SourceLink:    "",
	}

	// sample link: https://scans-hot.planeptune.us/manga/Kingdom/0825-001.png
	c.OnHTML("head > link:nth-child(18)", func(e *colly.HTMLElement) {
		firstImageLink := e.Attr("href")
		imageLinkSplit := strings.Split(firstImageLink, "/")

		if len(imageLinkSplit) == 0 {
			return
		}

		chapterAndImageIdx := imageLinkSplit[len(imageLinkSplit)-1]
		chapterAndImageIdxSplit := strings.Split(chapterAndImageIdx, "-")
		if len(chapterAndImageIdxSplit) != 2 {
			return
		}

		chapterNoStr := chapterAndImageIdxSplit[0]
		chapter.Number = utils.ForceSanitizeStringToFloat(chapterNoStr)

		imageNoAndExtensionSplit := strings.Split(chapterAndImageIdxSplit[1], ".")
		if len(imageNoAndExtensionSplit) != 2 {
			return
		}

		// imageNoStr := imageNoAndExtensionSplit[0]
		extension := imageNoAndExtensionSplit[1]

		// https://scans-hot.planeptune.us/manga/Kingdom
		imageLinkPrefix := strings.TrimSuffix(firstImageLink, chapterAndImageIdx)

		for i := 1; i <= 150; i++ {
			imageUrl := fmt.Sprintf("%s%s-%03d.%s", imageLinkPrefix, chapterNoStr, i, extension)

			// _, err = http.Head(imageUrl)
			// if err != nil {
			// 	break
			// }

			chapter.ChapterImages = append(chapter.ChapterImages, contract.ChapterImage{
				Index:     int64(i),
				ImageUrls: []string{imageUrl},
			})
		}
	})

	targetLinks := []string{
		fmt.Sprintf("%v/chapters/%v", sc.Host, queryParams.ChapterID),
	}
	for _, targetLink := range targetLinks {
		c.Visit(targetLink)
	}

	chapter.SourceLink = targetLinks[0]

	return chapter, err
}
