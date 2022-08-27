package manga_scrapper

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/config"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/utils/utils"
)

func GetFizmangaLatestManga(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	var err error
	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	mangas := []models.Manga{}

	selector := "div#loop-content > div.page-listing-item > div.row > div.badge-pos-1"
	if queryParams.Page > 1 {
		selector = "div.page-listing-item > div.row > div.badge-pos-1"
	}

	c.OnHTML(selector, func(e *colly.HTMLElement) {
		sourceID := e.ChildAttr("div.item-thumb a", "href")
		sourceID = strings.Replace(sourceID, "https://fizmanga.com/manga/", "", -1)
		sourceID = strings.Replace(sourceID, "/", "", -1)

		latestChapterID := e.ChildAttr("div.list-chapter > div:nth-child(1) > span.chapter.font-meta > a", "href")
		latestChapterID = strings.Replace(latestChapterID, "https://fizmanga.com/manga/", "", -1)
		latestChapterID = strings.Replace(latestChapterID, sourceID, "", -1)
		latestChapterID = strings.Replace(latestChapterID, "/", "", -1)

		reg, _ := regexp.Compile("[^0-9]+")
		latestChapterNumberString := reg.ReplaceAllString(latestChapterID, "")
		number, _ := strconv.ParseFloat(latestChapterNumberString, 64)

		imageUrl := fmt.Sprintf("%v/mangas/fizmanga/image_proxy/%v", config.Get().AnimapuOnlineHost, e.ChildAttr("a img", "data-lazy-src"))
		if queryParams.Page > 1 {
			imageUrl = fmt.Sprintf("%v/mangas/fizmanga/image_proxy/%v", config.Get().AnimapuOnlineHost, e.ChildAttr("a img", "src"))
		}

		mangas = append(mangas, models.Manga{
			ID:                  sourceID,
			Source:              "fizmanga",
			SourceID:            sourceID,
			Title:               e.ChildText("h3.h5 a"),
			LatestChapterID:     latestChapterID,
			LatestChapterNumber: number,
			LatestChapterTitle:  e.ChildText("div.list-chapter > div:nth-child(1) > span.chapter.font-meta > a"),
			CoverImages: []models.CoverImage{
				{
					Index:     1,
					ImageUrls: []string{imageUrl},
				},
			},
		})
	})

	if queryParams.Page <= 1 {
		err = c.Visit("https://fizmanga.com/")
	} else {
		requestData := strings.NewReader(
			strings.Join([]string{
				"action=madara_load_more",
				fmt.Sprintf("page=%v", queryParams.Page-1),
				"template=madara-core/content/content-archive",
				"vars[orderby]=meta_value_num",
				"vars[paged]=1",
				"vars[posts_per_page]=20",
				"vars[tax_query][relation]=OR",
				"vars[meta_query][0][relation]=AND",
				"vars[meta_query][relation]=OR",
				"vars[post_type]=wp-manga",
				"vars[post_status]=publish",
				"vars[meta_key]=_latest_update",
				"vars[order]=desc",
				"vars[sidebar]=right",
				"vars[manga_archives_item_layout]=big_thumbnail",
			}, "&"),
		)
		err = c.Request(
			"POST",
			"https://fizmanga.com/wp-admin/admin-ajax.php",
			requestData,
			colly.NewContext(),
			http.Header{
				"content-type":       []string{"application/x-www-form-urlencoded; charset=UTF-8"},
				"Authority":          []string{"fizmanga.com"},
				"Sec-Ch-Ua":          []string{"\"Google Chrome\";v=\"93\", \" Not;A Brand\";v=\"99\", \"Chromium\";v=\"93\""},
				"Accept-Language":    []string{"en-US,en;q=0.9,id;q=0.8"},
				"Sec-Ch-Ua-Mobile":   []string{"?0"},
				"User-Agent":         []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.82 Safari/537.36"},
				"Content-Type":       []string{"application/x-www-form-urlencoded; charset=UTF-8"},
				"Accept":             []string{"*/*"},
				"X-Requested-With":   []string{"XMLHttpRequest"},
				"Sec-Ch-Ua-Platform": []string{"\"macOS\""},
				"Origin":             []string{"https://fizmanga.com"},
				"Sec-Fetch-Site":     []string{"same-origin"},
				"Sec-Fetch-Mode":     []string{"cors"},
				"Sec-Fetch-Dest":     []string{"empty"},
				"Referer":            []string{"https://fizmanga.com/"},
			},
		)
	}

	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return mangas, err
	}

	return mangas, nil
}

func GetFizmangaDetailManga(ctx context.Context, queryParams models.QueryParams) (models.Manga, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	manga := models.Manga{
		ID:                  queryParams.SourceID,
		Source:              "fizmanga",
		SourceID:            queryParams.SourceID,
		Title:               "Untitled",
		Description:         "Description unavailable",
		Genres:              []string{},
		Status:              "Ongoing",
		CoverImages:         []models.CoverImage{{ImageUrls: []string{}}},
		Chapters:            []models.Chapter{},
		LatestChapterID:     "",
		LatestChapterNumber: 0,
		LatestChapterTitle:  "",
	}

	c.OnHTML("div.post-title > h1", func(e *colly.HTMLElement) {
		manga.Title = e.Text
	})

	c.OnHTML("div.description-summary > div.summary__content", func(e *colly.HTMLElement) {
		manga.Description = e.Text
	})

	c.OnHTML("div.tab-summary > div.summary_image > a > img", func(e *colly.HTMLElement) {
		manga.CoverImages = []models.CoverImage{
			{Index: 1, ImageUrls: []string{
				fmt.Sprintf("%v/mangas/fizmanga/image_proxy/%v", config.Get().AnimapuOnlineHost, e.Attr("data-lazy-src")),
			}},
		}
	})

	var mangaShortID string
	c.OnHTML("input.rating-post-id", func(e *colly.HTMLElement) {
		mangaShortID = e.Attr("value")
	})

	err := c.Visit(fmt.Sprintf("https://fizmanga.com/manga/%v", queryParams.SourceID))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return manga, err
	}

	idx := int64(1)
	c.OnHTML("li.wp-manga-chapter", func(e *colly.HTMLElement) {
		rawChapterLink := e.ChildAttr("a", "href")
		chapterLinkId := strings.Replace(rawChapterLink, fmt.Sprintf("https://fizmanga.com/manga/%v", queryParams.SourceID), "", -1)
		chapterLinkId = strings.Replace(chapterLinkId, "/", "", -1)
		chapterNumber, _ := strconv.ParseFloat(utils.RemoveNonNumeric(chapterLinkId), 64)

		manga.Chapters = append(manga.Chapters, models.Chapter{
			SourceID:      queryParams.SourceID,
			Source:        "fizmanga",
			Index:         idx,
			ID:            chapterLinkId,
			Title:         e.ChildText("a"),
			Number:        chapterNumber,
			ChapterImages: []models.ChapterImage{},
		})
		idx += 1
	})

	requestData := strings.NewReader(
		strings.Join([]string{
			"action=manga_get_chapters",
			fmt.Sprintf("manga=%v", mangaShortID),
		}, "&"),
	)
	err = c.Request(
		"POST",
		"https://fizmanga.com/wp-admin/admin-ajax.php",
		requestData,
		colly.NewContext(),
		http.Header{
			"content-type":     []string{"application/x-www-form-urlencoded; charset=UTF-8"},
			"Authority":        []string{"fizmanga.com"},
			"Sec-Ch-Ua":        []string{"\"Google Chrome\";v=\"93\", \" Not;A Brand\";v=\"99\", \"Chromium\";v=\"93\""},
			"X-Requested-With": []string{"XMLHttpRequest"},
			"Origin":           []string{"https://fizmanga.com"},
			"Referer":          []string{"https://fizmanga.com/"},
		},
	)

	manga.GenerateLatestChapter()

	return manga, nil
}

func GetFizmangaByQuery(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)
	var err error

	mangas := []models.Manga{}

	selector := "div.c-tabs-item > div.row.c-tabs-item__content"
	if queryParams.Page > 1 {
		selector = "div.row.c-tabs-item__content"
	}

	c.OnHTML(selector, func(e *colly.HTMLElement) {
		sourceID := e.ChildAttr("div.tab-thumb.c-image-hover a", "href")
		sourceID = strings.Replace(sourceID, "https://fizmanga.com/manga/", "", -1)
		sourceID = strings.Replace(sourceID, "/", "", -1)

		latestChapterID := e.ChildAttr("div.font-meta.chapter a", "href")
		latestChapterID = strings.Replace(latestChapterID, "https://fizmanga.com/manga/", "", -1)
		latestChapterID = strings.Replace(latestChapterID, sourceID, "", -1)
		latestChapterID = strings.Replace(latestChapterID, "/", "", -1)

		reg, _ := regexp.Compile("[^0-9]+")
		latestChapterNumberString := reg.ReplaceAllString(latestChapterID, "")
		number, _ := strconv.ParseFloat(latestChapterNumberString, 64)

		mangas = append(mangas, models.Manga{
			ID:                  sourceID,
			Source:              "fizmanga",
			SourceID:            sourceID,
			Title:               e.ChildText("h3.h4 a"),
			LatestChapterID:     latestChapterID,
			LatestChapterNumber: number,
			LatestChapterTitle:  e.ChildText("div.font-meta.chapter a"),
			CoverImages: []models.CoverImage{
				{
					Index: 1,
					ImageUrls: []string{
						fmt.Sprintf("%v/mangas/fizmanga/image_proxy/%v", config.Get().AnimapuOnlineHost, e.ChildAttr("a img", "src")),
					},
				},
			},
		})
	})

	if queryParams.Page <= 1 {
		err = c.Visit(fmt.Sprintf("https://fizmanga.com/?s=%v&post_type=wp-manga", queryParams.Title))
	} else {
		requestData := strings.NewReader(
			strings.Join([]string{
				"action=madara_load_more",
				fmt.Sprintf("page=%v", queryParams.Page-1),
				"template=madara-core/content/content-search",
				fmt.Sprintf("vars[s]:%v", queryParams.Title),
				"vars[orderby]:",
				"vars[paged]:1",
				"vars[template]:search",
				"vars[meta_query][0][relation]:AND",
				"vars[meta_query][relation]:OR",
				"vars[post_type]:wp-manga",
				"vars[post_status]:publish",
				"vars[manga_archives_item_layout]:default",
			}, "&"),
		)
		err = c.Request(
			"POST",
			"Request URL: https://fizmanga.com/wp-admin/admin-ajax.php",
			requestData,
			colly.NewContext(),
			http.Header{
				"Authority":      []string{"fizmanga.com"},
				"Origin":         []string{"https://fizmanga.com"},
				"Sec-Fetch-Site": []string{"same-origin"},
				"Sec-Fetch-Mode": []string{"cors"},
				"Sec-Fetch-Dest": []string{"empty"},
				"Referer":        []string{"https://fizmanga.com/"},
			},
		)
	}
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return mangas, err
	}

	return mangas, nil
}

func GetFizmangaDetailChapter(ctx context.Context, queryParams models.QueryParams) (models.Chapter, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	chapterNumberSplitted := strings.Split(queryParams.ChapterID, "-")
	chapterNumber, _ := strconv.ParseFloat(chapterNumberSplitted[1], 64)

	chapter := models.Chapter{
		ID:            queryParams.ChapterID,
		SourceID:      queryParams.SourceID,
		Source:        "fizmanga",
		Number:        chapterNumber,
		ChapterImages: []models.ChapterImage{},
		SourceLink:    fmt.Sprintf("https://fizmanga.com/manga/%v/%v", queryParams.SourceID, queryParams.ChapterID),
	}

	idx := int64(1)
	c.OnHTML("img.wp-manga-chapter-img", func(e *colly.HTMLElement) {
		chapter.ChapterImages = append(chapter.ChapterImages, models.ChapterImage{
			Index: idx,
			ImageUrls: []string{
				e.Attr("data-lazy-src"),
				fmt.Sprintf("%v/mangas/fizmanga/image_proxy/%v", config.Get().AnimapuLocalHost, e.Attr("data-lazy-src")),
			},
		})
		idx += 1
	})

	err := c.Visit(fmt.Sprintf("https://fizmanga.com/manga/%v/%v", queryParams.SourceID, queryParams.ChapterID))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return chapter, err
	}

	return chapter, nil
}
