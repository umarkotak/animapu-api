package manga_scrapper_repository

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/config"
	"github.com/umarkotak/animapu-api/internal/contract"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/utils/utils"
)

type Komiku struct {
	Host    string
	ApiHost string
	Source  string
}

func NewKomiku() Komiku {
	return Komiku{
		Host:    "https://komiku.org",
		ApiHost: "https://api.komiku.org",
		Source:  "komiku",
	}
}

func (sc *Komiku) GetHome(ctx context.Context, queryParams models.QueryParams) ([]contract.Manga, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(config.Get().CollyTimeout)

	mangas := []contract.Manga{}

	c.OnHTML("body > div.bge", func(e *colly.HTMLElement) {
		mangaID := e.ChildAttr("div.bgei > a", "href")
		mangaID = strings.ReplaceAll(mangaID, sc.Host, "")
		mangaID = strings.ReplaceAll(mangaID, "/manga/", "")
		mangaID = strings.TrimSuffix(mangaID, "/")

		latestChapterNumber := float64(0)
		e.ForEach("div.new1", func(i int, h *colly.HTMLElement) {
			if !strings.Contains(h.Text, "Terbaru") {
				return
			}
			latestChapterNumber = utils.ForceSanitizeStringToFloat(h.Text)
		})

		mangas = append(mangas, contract.Manga{
			ID:                  mangaID,
			SourceID:            mangaID,
			Source:              sc.Source,
			Title:               e.ChildText("body > div > div.kan > a > h3"),
			Genres:              []string{},
			LatestChapterID:     "",
			LatestChapterNumber: latestChapterNumber,
			LatestChapterTitle:  "",
			Chapters:            []contract.Chapter{},
			CoverImages: []contract.CoverImage{
				{
					Index: 1,
					ImageUrls: []string{
						// strings.ReplaceAll(e.ChildAttr("div.bgei > a > img", "src"), "450,235", "270,450"),
						strings.ReplaceAll(e.ChildAttr("div.bgei > a > img", "src"), "resize", "resizee"),
					},
				},
			},
		})
	})

	err := c.Visit(fmt.Sprintf("%s/manga/page/%v/?orderby=modified&category_name&genre&genre2&status", sc.ApiHost, queryParams.Page))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
	}

	return mangas, nil
}

func (sc *Komiku) GetDetail(ctx context.Context, queryParams models.QueryParams) (contract.Manga, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(config.Get().CollyTimeout)
	c.AllowURLRevisit = true

	manga := contract.Manga{
		ID:          queryParams.SourceID,
		Source:      sc.Source,
		SourceID:    queryParams.SourceID,
		Title:       "Untitled",
		Description: "Description unavailable",
		Genres:      []string{},
		Status:      "Ongoing",
		CoverImages: []contract.CoverImage{{ImageUrls: []string{}}},
		Chapters:    []contract.Chapter{},
	}

	c.OnHTML("#Judul > p.j2", func(e *colly.HTMLElement) {
		if e.Text == "" {
			return
		}
		manga.Title = e.Text
	})

	c.OnHTML("#Judul > p.desc", func(e *colly.HTMLElement) {
		if e.Text == "" {
			return
		}
		manga.Description = e.Text
	})

	c.OnHTML("#Informasi > div > img", func(e *colly.HTMLElement) {
		if e.Attr("src") == "" {
			return
		}
		manga.CoverImages = []contract.CoverImage{{ImageUrls: []string{
			e.Attr("src"),
		}}}
	})

	c.OnHTML("#daftarChapter > tr", func(e *colly.HTMLElement) {
		chapterLink := e.ChildAttr("td.judulseries > a", "href")
		chapterID := chapterLink
		chapterID = strings.ReplaceAll(chapterID, "/", "")

		if chapterLink == "" {
			return
		}

		manga.Chapters = append(manga.Chapters, contract.Chapter{
			ID:       chapterID,
			Source:   sc.Source,
			SourceID: chapterID,
			Title:    e.ChildText("td.judulseries > a"),
			Index:    utils.ForceSanitizeStringToInt64(chapterLink),
			Number:   utils.ForceSanitizeStringToFloat(chapterLink),
		})
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

func (sc *Komiku) GetSearch(ctx context.Context, queryParams models.QueryParams) ([]contract.Manga, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(config.Get().CollyTimeout)

	mangas := []contract.Manga{}

	c.OnHTML("body > div.bge", func(e *colly.HTMLElement) {
		mangaID := e.ChildAttr("div.bgei > a", "href")
		mangaID = strings.ReplaceAll(mangaID, sc.Host, "")
		mangaID = strings.ReplaceAll(mangaID, "/manga/", "")
		mangaID = strings.TrimSuffix(mangaID, "/")

		mangas = append(mangas, contract.Manga{
			ID:                  mangaID,
			SourceID:            mangaID,
			Source:              sc.Source,
			Title:               e.ChildText("body > div > div.kan > a > h3"),
			Genres:              []string{},
			LatestChapterID:     "",
			LatestChapterNumber: 0,
			LatestChapterTitle:  "",
			Chapters:            []contract.Chapter{},
			CoverImages: []contract.CoverImage{
				{
					Index: 1,
					ImageUrls: []string{
						// strings.ReplaceAll(e.ChildAttr("div.bgei > a > img", "src"), "450,235", "270,450"),
						strings.ReplaceAll(e.ChildAttr("div.bgei > a > img", "src"), "resize", "resizee"),
					},
				},
			},
		})
	})

	q := strings.ReplaceAll(queryParams.Title, " ", "+")
	err := c.Visit(fmt.Sprintf("%s/?post_type=manga&s=%s", sc.ApiHost, q))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
	}
	if err != nil {
		return mangas, err
	}

	return mangas, nil
}

func (sc *Komiku) GetChapter(ctx context.Context, queryParams models.QueryParams) (contract.Chapter, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(config.Get().CollyTimeout)

	chapterNumber := float64(0)

	splitted := strings.Split(queryParams.ChapterID, "chapter-")
	if len(splitted) > 0 {
		chapterNumber = utils.ForceSanitizeStringToFloat(splitted[len(splitted)-1])
	}

	targetLink := fmt.Sprintf("%v/%v", sc.Host, queryParams.ChapterID)

	chapter := contract.Chapter{
		ID:            queryParams.ChapterID,
		SourceID:      queryParams.SourceID,
		Source:        sc.Source,
		SourceLink:    targetLink,
		Number:        chapterNumber,
		ChapterImages: []contract.ChapterImage{},
	}

	c.OnHTML("#Baca_Komik", func(e *colly.HTMLElement) {
		e.ForEach("img", func(i int, h *colly.HTMLElement) {
			chapterImage := contract.ChapterImage{
				Index: 0,
				ImageUrls: []string{
					h.Attr("src"),
				},
			}

			altImage := sc.ExtractOnErrorImg(h.Attr("onerror"))
			if altImage != "" {
				chapterImage.ImageUrls = append(chapterImage.ImageUrls, altImage)
			}

			chapter.ChapterImages = append(chapter.ChapterImages, chapterImage)
		})
	})

	err := c.Visit(targetLink)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return chapter, err
	}

	return chapter, nil
}

func (sc *Komiku) ExtractOnErrorImg(str string) string {
	re := regexp.MustCompile(`(?i)this\.src='([^']+)'`)

	matches := re.FindStringSubmatch(str)

	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}
