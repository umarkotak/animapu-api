package manga_scrapper_repository

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/config"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/utils/utils"
)

type Mangasee struct {
	Host    string
	Source  string
	ImgHost string
}

type MangaseeManga struct {
	SeriesID   string    `json:"SeriesID"`
	IndexName  string    `json:"IndexName"`
	SeriesName string    `json:"SeriesName"`
	ScanStatus string    `json:"ScanStatus"`
	Chapter    string    `json:"Chapter"`
	Genres     string    `json:"Genres"`
	Date       time.Time `json:"Date"`
	IsEdd      bool      `json:"IsEdd"`
}

type MangaseeChapter struct {
	Chapter     string      `json:"Chapter"`
	Type        string      `json:"Type"`
	Date        string      `json:"Date"`
	ChapterName interface{} `json:"ChapterName"`
}

func NewMangasee() Mangasee {
	return Mangasee{
		Host:    "https://www.mangasee123.com",
		Source:  "mangasee",
		ImgHost: "https://temp.compsci88.com",
	}
}

func (sc *Mangasee) GetHome(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	mangas := []models.Manga{}

	c.OnHTML("body > script:nth-child(16)", func(e *colly.HTMLElement) {
		footerContent := e.Text
		splitted := strings.Split(footerContent, "vm.LatestJSON = ")
		if len(splitted) <= 0 {
			return
		}

		splitted = strings.Split(splitted[1], "vm.NewSeriesJSON")
		dataJson := splitted[0]
		dataJson = strings.ReplaceAll(dataJson, ";", "")
		dataJson = strings.TrimSpace(dataJson)

		mangaseeMangas := []MangaseeManga{}
		json.Unmarshal([]byte(dataJson), &mangaseeMangas)

		for _, oneMangaseeManga := range mangaseeMangas {
			chNumber := mangaseeDecodeCh(oneMangaseeManga.Chapter)

			mangas = append(mangas, models.Manga{
				ID:                  oneMangaseeManga.IndexName,
				SourceID:            oneMangaseeManga.IndexName,
				Source:              sc.Source,
				Title:               oneMangaseeManga.SeriesName,
				Genres:              []string{},
				LatestChapterID:     "",
				LatestChapterNumber: utils.ForceSanitizeStringToFloat(chNumber),
				LatestChapterTitle:  chNumber,
				Chapters:            []models.Chapter{},
				CoverImages: []models.CoverImage{
					{
						Index: 1,
						ImageUrls: []string{
							fmt.Sprintf("%v/cover/%v.jpg", sc.ImgHost, oneMangaseeManga.IndexName),
						},
					},
				},
			})
		}
	})

	err := c.Visit(fmt.Sprintf("%v", sc.Host))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
	}

	return mangas, nil
}

func (sc *Mangasee) GetDetail(ctx context.Context, queryParams models.QueryParams) (models.Manga, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)
	c.AllowURLRevisit = true

	manga := models.Manga{
		ID:          queryParams.SourceID,
		Source:      sc.Source,
		SourceID:    queryParams.SourceID,
		Title:       strings.ReplaceAll(queryParams.SourceID, "-", " "),
		Description: "Description unavailable",
		Genres:      []string{},
		Status:      "Ongoing",
		CoverImages: []models.CoverImage{{ImageUrls: []string{
			fmt.Sprintf("%v/cover/%v.jpg", sc.ImgHost, queryParams.SourceID),
		}}},
		Chapters: []models.Chapter{},
	}

	c.OnHTML("div.top-5.Content", func(e *colly.HTMLElement) {
		manga.Description = e.Text
	})

	c.OnHTML("body > script:nth-child(16)", func(e *colly.HTMLElement) {
		footerContent := e.Text
		splitted := strings.Split(footerContent, "vm.Chapters = ")
		if len(splitted) <= 0 {
			return
		}

		splitted = strings.Split(splitted[1], "vm.NumSubs")
		dataJson := splitted[0]
		dataJson = strings.ReplaceAll(dataJson, ";", "")
		dataJson = strings.TrimSpace(dataJson)

		mangaseeChapters := []MangaseeChapter{}
		json.Unmarshal([]byte(dataJson), &mangaseeChapters)

		for i, oneMangaseeChapter := range mangaseeChapters {
			firstNum := oneMangaseeChapter.Chapter[0:1]
			lastNum := oneMangaseeChapter.Chapter[len(oneMangaseeChapter.Chapter)-1:]
			chNumberS := mangaseeDecodeCh(oneMangaseeChapter.Chapter)
			if lastNum != "0" {
				chNumberS = fmt.Sprintf("%v.%v", chNumberS, lastNum)
			}

			chNumer := utils.ForceSanitizeStringToFloat(chNumberS)

			manga.Chapters = append(manga.Chapters, models.Chapter{
				ID:                fmt.Sprint(chNumer),
				Source:            sc.Source,
				SourceID:          fmt.Sprint(chNumer),
				SecondarySourceID: fmt.Sprint(firstNum),
				Title:             fmt.Sprintf("%v %v", oneMangaseeChapter.Type, chNumer),
				Index:             int64(i),
				Number:            chNumer,
			})
		}
	})

	targetUrl := fmt.Sprintf("%v/manga/%v", sc.Host, queryParams.SourceID)
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

func (sc *Mangasee) GetSearch(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
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

func (sc *Mangasee) GetChapter(ctx context.Context, queryParams models.QueryParams) (models.Chapter, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	targetLink := fmt.Sprintf("%v/read-online/%v-chapter-%v.html", sc.Host, queryParams.SourceID, queryParams.ChapterID)
	if queryParams.SecondarySourceID == "2" {
		targetLink = fmt.Sprintf("%v/read-online/%v-chapter-%v-index-2.html", sc.Host, queryParams.SourceID, queryParams.ChapterID)
	}

	chapter := models.Chapter{
		ID:            queryParams.ChapterID,
		SourceID:      queryParams.SourceID,
		Source:        sc.Source,
		SourceLink:    targetLink,
		Number:        utils.ForceSanitizeStringToFloat(queryParams.ChapterID),
		ChapterImages: []models.ChapterImage{},
	}

	c.OnHTML("body > script:nth-child(19)", func(e *colly.HTMLElement) {
		re := regexp.MustCompile(`vm\.CurPathName\s*=\s*("[^"]+")`)

		// Find the first match of the pattern
		imageHost := re.FindStringSubmatch(e.Text)

		chFloat := utils.ForceSanitizeStringToFloat(queryParams.ChapterID)
		chInt := int(chFloat)
		modifier := ""
		if (chFloat - float64(chInt)) > 0 {
			splitted := strings.Split(queryParams.ChapterID, ".")
			modifier = fmt.Sprintf(".%v", splitted[1])
		}

		for i := 1; i <= 100; i++ {
			chapter.ChapterImages = append(chapter.ChapterImages, models.ChapterImage{
				Index: 0,
				ImageUrls: []string{
					fmt.Sprintf(
						"https://%v/manga/%v/%04d%v-%03d.png",
						strings.ReplaceAll(imageHost[1], `"`, ""), queryParams.SourceID, chInt, modifier, i,
					),
					fmt.Sprintf(
						"https://%v/manga/%v/Mag-Official/%04d%v-%03d.png",
						strings.ReplaceAll(imageHost[1], `"`, ""), queryParams.SourceID, chInt, modifier, i,
					),
				},
			})
		}
	})

	err := c.Visit(targetLink)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return chapter, err
	}

	return chapter, nil
}

func mangaseeDecodeCh(s string) string {
	return trimFirstRune(trimLastChar(s))
}

func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}

func trimLastChar(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return s[:len(s)-size]
}