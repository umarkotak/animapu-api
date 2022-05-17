package manga_scrapper

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
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

		fmt.Println(e.ChildAttr("div.item-thumb > a > img", "data-lazy-src"))

		mangas = append(mangas, models.Manga{
			ID:                  sourceID,
			Source:              "fizmanga",
			SourceID:            sourceID,
			Title:               e.ChildText("h3.h5 a"),
			LatestChapterID:     "lastestChapterID",
			LatestChapterNumber: 0,
			LatestChapterTitle:  "latestChapterTitle",
			CoverImages: []models.CoverImage{
				{
					Index: 1,
					ImageUrls: []string{
						fmt.Sprintf("http://localhost:6001/mangas/fizmannga/image_proxy/%v", e.ChildAttr("a img", "data-lazy-src")),
						fmt.Sprintf("http://localhost:6001/mangas/fizmannga/image_proxy/%v", e.ChildAttr("a img", "src")),
					},
				},
			},
		})
	})

	if queryParams.Page <= 1 {
		err = c.Visit("https://fizmanga.com/")
	} else {
		requestData := strings.NewReader(fmt.Sprintf(`action=madara_load_more&page=%v&template=madara-core/content/content-archive&vars[orderby]=meta_value_num&vars[paged]=1&vars[posts_per_page]=40&vars[tax_query][relation]=OR&vars[meta_query][0][relation]=AND&vars[meta_query][relation]=OR&vars[post_type]=wp-manga&vars[post_status]=publish&vars[meta_key]=_latest_update&vars[order]=desc&vars[sidebar]=right&vars[manga_archives_item_layout]=big_thumbnail`, queryParams.Page-1))
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

	manga := models.Manga{}

	err := c.Visit(fmt.Sprintf("https://m.mangabat.com/manga-list-all/%v", queryParams.Page))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return manga, err
	}

	return manga, nil
}

func GetFizmangaByQuery(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	mangas := []models.Manga{}

	err := c.Visit(fmt.Sprintf("https://m.mangabat.com/manga-list-all/%v", queryParams.Page))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return mangas, err
	}

	return mangas, nil
}

func GetFizmangaDetailChapter(ctx context.Context, queryParams models.QueryParams) (models.Chapter, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	chapter := models.Chapter{}

	err := c.Visit(fmt.Sprintf("https://m.mangabat.com/manga-list-all/%v", queryParams.Page))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return chapter, err
	}

	return chapter, nil
}
