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
						fmt.Sprintf("http://localhost:60001/mangas/klikmanga/image_proxy/%s", e.ChildAttr("img.img-responsive", "src")),
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

	manga := models.Manga{
		ID:          queryParams.SourceID,
		Source:      "klikmanga",
		SourceID:    queryParams.SourceID,
		Chapters:    []models.Chapter{},
		CoverImages: []models.CoverImage{},
	}

	c.OnHTML("body > div.wrap > div > div.site-content > div > div.profile-manga > div > div > div > div.tab-summary > div.summary_image > a > img", func(e *colly.HTMLElement) {
		// manga.CoverImages = []models.CoverImage{
		// 	{Index: 1, ImageUrls: []string{e.Attr("src")}},
		// }

		manga.CoverImages = []models.CoverImage{{
			Index: 1,
			ImageUrls: []string{
				fmt.Sprintf("http://localhost:60001/mangas/klikmanga/image_proxy/%s", e.Attr("src")),
				fmt.Sprintf("https://thumb.mghubcdn.com/mn/%s.jpg", queryParams.SourceID),
				fmt.Sprintf("https://thumb.mghubcdn.com/md/%s.jpg", queryParams.SourceID),
				fmt.Sprintf("https://thumb.mghubcdn.com/m4l/%s.jpg", queryParams.SourceID),
			},
		}}
	})

	c.OnHTML("body > div.wrap > div > div.site-content > div > div.c-page-content.style-1 > div > div > div > div.main-col.col-md-8.col-sm-8 > div > div.c-page > div > div.description-summary > div > p", func(e *colly.HTMLElement) {
		manga.Description = e.Text
	})

	c.OnHTML("body > div.wrap > div > div.site-content > div > div.profile-manga > div > div > div > div.tab-summary > div.summary_content_wrap > div > div.post-content > div:nth-child(8) > div.summary-content > div", func(e *colly.HTMLElement) {
		manga.Genres = strings.Split(e.Text, ", ")
	})

	c.OnHTML("body > div.wrap > div > div.site-content > div > div.profile-manga > div > div > div > div.post-title > h1", func(e *colly.HTMLElement) {
		manga.Title = e.Text
	})

	idx := int64(1)
	c.OnHTML("body > div.wrap > div > div.site-content > div > div.c-page-content.style-1 > div > div > div > div.main-col.col-md-8.col-sm-8 > div > div.c-page > div > div.page-content-listing.single-page > div > ul > li", func(e *colly.HTMLElement) {
		rawChapterLink := e.ChildAttr("a", "href")
		chapterLinkId := strings.Replace(rawChapterLink, fmt.Sprintf("https://klikmanga.id/manga/%v", queryParams.SourceID), "", -1)
		chapterLinkId = strings.Replace(chapterLinkId, "/", "", -1)

		chapterString := e.ChildText("a")
		chapterString = strings.Replace(chapterString, "Chapter ", "", -1)
		chapterNumber, _ := strconv.ParseFloat(chapterString, 64)

		manga.Chapters = append(manga.Chapters, models.Chapter{
			ID:            chapterLinkId,
			SourceID:      queryParams.SourceID,
			Source:        "klikmanga",
			Title:         chapterString,
			Index:         idx,
			Number:        chapterNumber,
			ChapterImages: []models.ChapterImage{},
		})
		idx += 1
	})

	err := c.Visit(fmt.Sprintf("https://klikmanga.id/manga/%v", queryParams.SourceID))
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

	selector := "body > div.wrap > div > div.site-content > div.c-page-content > div > div > div > div > div.main-col-inner > div > div.tab-content-wrap > div > div"
	if queryParams.Page > 1 {
		selector = "div.row.c-tabs-item__content"
	}

	c.OnHTML(selector, func(e *colly.HTMLElement) {
		rawLink := e.ChildAttr("div div a", "href")
		sourceID := strings.Replace(rawLink, "https://klikmanga.id/manga/", "", -1)
		sourceID = strings.Replace(sourceID, "/", "", -1)

		lastChapterLink := e.ChildAttr("div.col-8.col-12.col-md-10 > div.tab-meta > div.meta-item.latest-chap > span.font-meta.chapter > a", "href")
		prefix := fmt.Sprintf("https://klikmanga.id/manga/%v", sourceID)
		lastChapterID := strings.Replace(lastChapterLink, prefix, "", -1)
		lastChapterID = strings.Replace(lastChapterID, "/", "", -1)

		lastChapterString := strings.Replace(lastChapterID, "chapter-", "", -1)
		lastChapterNumber, _ := strconv.ParseFloat(lastChapterString, 64)

		mangas = append(mangas, models.Manga{
			ID:                  sourceID,
			SourceID:            sourceID,
			Source:              "klikmanga",
			Title:               e.ChildText("div div h3 a"),
			LatestChapterID:     lastChapterID,
			LatestChapterNumber: lastChapterNumber,
			LatestChapterTitle:  e.ChildText("div.col-8.col-12.col-md-10 > div.tab-meta > div.meta-item.latest-chap > span.font-meta.chapter > a"),
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
		err = c.Visit(fmt.Sprintf("https://klikmanga.id/?s=%v&post_type=wp-manga&op=&author=&artist=&release=&adult=&m_orderby=latest", queryParams.Title))
	} else {
		requestData := strings.NewReader(fmt.Sprintf(`action=madara_load_more&page=%v&template=madara-core/content/content-search&vars[s]=%v&vars[orderby]=meta_value_num&vars[paged]=1&vars[template]=search&vars[meta_query][0][0][key]=_wp_manga_status&vars[meta_query][0][0][value][]=end&vars[meta_query][0][0][compare]=IN&vars[meta_query][0][relation]=AND&vars[meta_query][relation]=OR&vars[post_type]=wp-manga&vars[post_status]=publish&vars[meta_key]=_latest_update&vars[order]=desc&vars[manga_archives_item_layout]=big_thumbnail`, queryParams.Page, queryParams.Title))
		err = c.Request(
			"POST",
			"https://klikmanga.id/wp-admin/admin-ajax.php",
			requestData,
			colly.NewContext(),
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

func GetKlikmangaDetailChapter(ctx context.Context, queryParams models.QueryParams) (models.Chapter, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	chapter := models.Chapter{
		ID:            queryParams.ChapterID,
		SourceID:      queryParams.ChapterID,
		Source:        "klikmanga",
		ChapterImages: []models.ChapterImage{},
	}

	idx := int64(1)
	c.OnHTML("body > div.wrap > div > div.site-content > div > div > div > div > div > div > div.c-blog-post > div.entry-content > div > div > div.reading-content > div", func(e *colly.HTMLElement) {
		chapter.ChapterImages = append(chapter.ChapterImages, models.ChapterImage{
			Index:     idx,
			ImageUrls: []string{e.ChildAttr("img.wp-manga-chapter-img", "src")},
		})

		idx += 1
	})

	err := c.Visit(fmt.Sprintf("https://klikmanga.id/manga/%v/%v/?style=list", queryParams.SourceID, queryParams.ChapterID))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return chapter, err
	}

	chapter.SourceLink = fmt.Sprintf("https://klikmanga.id/manga/%v/%v/?style=list", queryParams.SourceID, queryParams.ChapterID)

	return chapter, nil
}
