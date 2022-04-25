package manga_scrapper

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
)

func GetKlikmangaLatestManga(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	mangas := []models.Manga{}

	selector := "#loop-content > div > div > div > div"
	if queryParams.Page > 1 {
		selector = "div.page-item-detail"
	}

	c.OnHTML(selector, func(e *colly.HTMLElement) {
		sourceID := e.ChildAttr("div.item-summary > div.post-title.font-title > h3 > a", "href")
		sourceID = strings.Replace(sourceID, "https://klikmanga.id/manga/", "", -1)
		sourceID = strings.Replace(sourceID, "/", "", -1)

		lastestChapterID := e.ChildAttr("div.item-summary > div.list-chapter > div:nth-child(1) > span.chapter.font-meta > a", "href")
		lastestChapterIDSplitted := strings.Split(lastestChapterID, "/")
		if len(lastestChapterIDSplitted) > 0 {
			lastestChapterID = lastestChapterIDSplitted[len(lastestChapterIDSplitted)-2]
		}

		latestChapterTitle := strings.Replace(lastestChapterID, "-", " ", -1)
		lastestChapterNumberString := strings.Replace(lastestChapterID, "chapter-", "", -1)
		lastestChapterNumberString = strings.Replace(lastestChapterNumberString, "-", ".", -1)
		latestChapterNumber, _ := strconv.ParseFloat(lastestChapterNumberString, 64)

		mangas = append(mangas, models.Manga{
			ID:                  sourceID,
			Source:              "klikmanga",
			SourceID:            sourceID,
			Title:               e.ChildText("div.item-summary > div.post-title.font-title > h3 > a"),
			LatestChapterID:     lastestChapterID,
			LatestChapterNumber: latestChapterNumber,
			LatestChapterTitle:  latestChapterTitle,
			CoverImages: []models.CoverImage{
				{
					Index: 1,
					ImageUrls: []string{
						fmt.Sprintf("https://thumb.mghubcdn.com/mn/%s.jpg", sourceID),
						fmt.Sprintf("https://thumb.mghubcdn.com/md/%s.jpg", sourceID),
						fmt.Sprintf("https://thumb.mghubcdn.com/m4l/%s.jpg", sourceID),
					},
				},
			},
		})
	})

	var err error
	if queryParams.Page <= 1 {
		err = c.Visit(fmt.Sprintf("https://klikmanga.id/"))
	} else {
		requestData := strings.NewReader(fmt.Sprintf(`action=madara_load_more&page=%v&template=madara-core/content/content-archive&vars[orderby]=meta_value_num&vars[paged]=1&vars[posts_per_page]=40&vars[tax_query][relation]=OR&vars[meta_query][0][relation]=AND&vars[meta_query][relation]=OR&vars[post_type]=wp-manga&vars[post_status]=publish&vars[meta_key]=_latest_update&vars[order]=desc&vars[sidebar]=right&vars[manga_archives_item_layout]=big_thumbnail`, queryParams.Page-1))
		err = c.Request(
			"POST", "https://klikmanga.id/wp-admin/admin-ajax.php", requestData, colly.NewContext(),
			http.Header{
				"content-type":       []string{"application/x-www-form-urlencoded; charset=UTF-8"},
				"Authority":          []string{"klikmanga.id"},
				"Sec-Ch-Ua":          []string{"\"Google Chrome\";v=\"93\", \" Not;A Brand\";v=\"99\", \"Chromium\";v=\"93\""},
				"Accept-Language":    []string{"en-US,en;q=0.9,id;q=0.8"},
				"Sec-Ch-Ua-Mobile":   []string{"?0"},
				"User-Agent":         []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.82 Safari/537.36"},
				"Content-Type":       []string{"application/x-www-form-urlencoded; charset=UTF-8"},
				"Accept":             []string{"*/*"},
				"X-Requested-With":   []string{"XMLHttpRequest"},
				"Sec-Ch-Ua-Platform": []string{"\"macOS\""},
				"Origin":             []string{"https://klikmanga.id"},
				"Sec-Fetch-Site":     []string{"same-origin"},
				"Sec-Fetch-Mode":     []string{"cors"},
				"Sec-Fetch-Dest":     []string{"empty"},
				"Referer":            []string{"https://klikmanga.id/"},
			},
		)
	}

	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return mangas, err
	}

	return mangas, nil
}

func GetKlikmangaDetailManga(ctx context.Context, queryParams models.QueryParams) (models.Manga, error) {
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

func GetKlikmangaByQuery(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
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

func GetKlikmangaDetailChapter(ctx context.Context, queryParams models.QueryParams) (models.Chapter, error) {
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
