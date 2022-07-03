package manga_scrapper

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository"
	"github.com/umarkotak/animapu-api/internal/utils/utils"
)

func GetMangahubLatestManga(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	mangas := []models.Manga{}

	fbMangaHubHome, err := repository.FbGetHomeByMangaSource(ctx, models.SOURCE_MANGAHUB)
	cachedMangas := fbMangaHubHome.Mangas
	if err == nil &&
		time.Now().UTC().Before(fbMangaHubHome.ExpiredAt) &&
		len(cachedMangas) > 0 {

		return MangasPaginate(cachedMangas, queryParams.Page, 30), nil
	}

	scrapeNinjaResponse, err := repository.QuickScrape(ctx, "https://mangahub.io")
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return MangasPaginate(cachedMangas, queryParams.Page, 30), nil
	}
	if scrapeNinjaResponse.Info.StatusCode != 200 {
		logrus.WithContext(ctx).Error(fmt.Errorf("Scrape ninja non 200"))
		return MangasPaginate(cachedMangas, queryParams.Page, 30), nil
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(scrapeNinjaResponse.Body))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return MangasPaginate(cachedMangas, queryParams.Page, 30), nil
	}

	doc.Find(".iqzwK.list-group-item").Each(func(i int, s *goquery.Selection) {
		imageUrl, _ := s.Find("img").Attr("src")

		mangaLink, _ := s.Find(".media-left a").Attr("href")
		mangaLinkSplit := strings.Split(mangaLink, "/")
		if len(mangaLinkSplit) <= 0 {
			return
		}
		mangaID := mangaLinkSplit[len(mangaLinkSplit)-1]

		chapterLink, _ := s.Find("a._8Qtbo.text-secondary._2euQb").Attr("href")
		chapterLinkSplit := strings.Split(chapterLink, "/")
		if len(chapterLinkSplit) <= 0 {
			return
		}
		chapterID := chapterLinkSplit[len(chapterLinkSplit)-1]

		chapterNumberString := utils.RemoveNonNumeric(chapterID)
		chapterNumber, _ := strconv.ParseFloat(chapterNumberString, 64)

		mangas = append(mangas, models.Manga{
			ID:                  mangaID,
			SourceID:            mangaID,
			Source:              models.SOURCE_MANGAHUB,
			Title:               s.Find("a._31Z6T.text-secondary").Text(),
			Description:         "Description unavailable",
			Genres:              []string{},
			Status:              "Ongoing",
			Rating:              "10",
			LatestChapterID:     chapterID,
			LatestChapterNumber: chapterNumber,
			LatestChapterTitle:  chapterID,
			Chapters:            []models.Chapter{},
			CoverImages: []models.CoverImage{
				{
					Index: 1,
					ImageUrls: []string{
						imageUrl,
					},
				},
			},
		})
	})

	// Cache mangahub home to firebase
	go func() {
		if err != nil || len(mangas) <= 0 {
			return
		}
		_, err = repository.FbUpsertHomeByMangaSource(context.Background(), models.SOURCE_MANGAHUB, mangas)
		if err != nil {
			logrus.WithContext(context.Background()).Error(err)
		}
	}()

	return MangasPaginate(mangas, queryParams.Page, 30), nil
}

func GetMangahubDetailManga(ctx context.Context, queryParams models.QueryParams) (models.Manga, error) {
	manga := models.Manga{
		ID:          queryParams.SourceID,
		Source:      models.SOURCE_MANGAHUB,
		SourceID:    queryParams.SourceID,
		Description: "Description unavailable",
		Genres:      []string{},
		Status:      "Ongoing",
		CoverImages: []models.CoverImage{{ImageUrls: []string{}}},
		Chapters:    []models.Chapter{},
	}

	fbMangaHubDetail, err := repository.FbGetMangaDetailByMangaSource(ctx, models.SOURCE_MANGAHUB, manga)
	cachedManga := fbMangaHubDetail.Manga
	if err != nil &&
		time.Now().UTC().Before(fbMangaHubDetail.ExpiredAt) &&
		len(manga.Chapters) > 0 {

		return cachedManga, nil
	}

	scrapeNinjaResponse, err := repository.QuickScrape(ctx, fmt.Sprintf("https://mangahub.io/manga/%v", manga.SourceID))
	if err != nil {
		if len(manga.Chapters) > 0 {
			return cachedManga, nil
		}
		logrus.WithContext(ctx).Error(err)
		return manga, nil
	}
	if scrapeNinjaResponse.Info.StatusCode != 200 {
		if len(manga.Chapters) > 0 {
			return cachedManga, nil
		}
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"scrape_ninja_response": scrapeNinjaResponse,
		}).Error(fmt.Errorf("Scrape ninja non 200"))
		return manga, nil
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(scrapeNinjaResponse.Body))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return manga, nil
	}

	manga.Title = doc.Find("#mangadetail h1._3xnDj").Text()
	additionalTitle := doc.Find("#mangadetail h1._3xnDj small").Text()
	manga.Title = strings.Replace(manga.Title, additionalTitle, "", -1)
	newTag := doc.Find("#mangadetail h1._3xnDj a").Text()
	manga.Title = strings.Replace(manga.Title, newTag, "", -1)

	manga.Description = additionalTitle
	manga.Description = fmt.Sprintf("%v\n%v", manga.Description, doc.Find("#noanim-content-tab-pane-99 div p.ZyMp7").Text())

	imageUrl, _ := doc.Find("#mangadetail img.img-responsive.manga-thumb").Attr("src")
	manga.CoverImages = []models.CoverImage{{ImageUrls: []string{
		imageUrl,
	}}}

	idx := int64(1)
	doc.Find("li._287KE.list-group-item").Each(func(i int, s *goquery.Selection) {
		chapterLink, _ := s.Find("a").Attr("href")
		chapterLink = strings.Replace(chapterLink, "https://mangahub.io/chapter/", "", -1)
		chapterLink = strings.Replace(chapterLink, queryParams.SourceID, "", -1)
		chapterID := strings.Replace(chapterLink, "/", "", -1)

		chapterNumberString := utils.RemoveNonNumeric(chapterID)
		chapterNumber, _ := strconv.ParseFloat(chapterNumberString, 64)

		manga.Chapters = append(manga.Chapters, models.Chapter{
			ID:       chapterID,
			Source:   models.SOURCE_MANGAHUB,
			SourceID: chapterID,
			Title:    fmt.Sprintf("%v %v", chapterNumber, s.Find("a span._2IG5P").Text()),
			Index:    idx,
			Number:   chapterNumber,
		})

		idx += 1
	})

	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return manga, err
	}

	// Cache mangahub manga detail to firebase
	if err == nil && len(manga.Chapters) > 0 {
		_, err = repository.FbUpsertMangaDetailByMangaSource(context.Background(), models.SOURCE_MANGAHUB, manga)
		if err != nil {
			logrus.WithContext(context.Background()).Error(err)
		}
		logrus.WithContext(context.Background()).Infof("detail caching finished")
	}

	return manga, nil
}

func GetMangahubByQuery(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	mangas := []models.Manga{}

	mangas = append(mangas, models.Manga{
		ID:                  "",
		SourceID:            "",
		Source:              "source",
		SecondarySourceID:   "",
		SecondarySource:     "secondary_source",
		Title:               "Untitled",
		Description:         "Description unavailable",
		Genres:              []string{},
		Status:              "Ongoing",
		Rating:              "10",
		LatestChapterID:     "chapter_id",
		LatestChapterNumber: 0,
		LatestChapterTitle:  "Chapter 0",
		CoverImages: []models.CoverImage{
			{
				Index: 1,
				ImageUrls: []string{
					fmt.Sprintf("https://animapu-lite.vercel.app/images/manga/%v", "image_id"),
				},
			},
		},
	})

	err := c.Visit(fmt.Sprintf("https://animapu-lite.vercel.app/search/%v", queryParams.Page))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return mangas, err
	}

	return mangas, nil
}

func GetMangahubDetailChapter(ctx context.Context, queryParams models.QueryParams) (models.Chapter, error) {
	pageCountConfig := int64(100)

	chapterNumber, _ := strconv.ParseFloat(utils.RemoveNonNumeric(queryParams.ChapterID), 64)

	chapter := models.Chapter{
		ID:            queryParams.ChapterID,
		SourceID:      queryParams.SourceID,
		Source:        models.SOURCE_MANGAHUB,
		Title:         "",
		Index:         0,
		Number:        chapterNumber,
		ChapterImages: []models.ChapterImage{},
	}

	for i := int64(1); i <= pageCountConfig; i++ {
		chapter.ChapterImages = append(chapter.ChapterImages, models.ChapterImage{
			Index: i,
			ImageUrls: []string{
				fmt.Sprintf("https://img.mghubcdn.com/file/imghub/%v/%v/%v.jpg", queryParams.SourceID, chapterNumber, i),
				fmt.Sprintf("https://img.mghubcdn.com/file/imghub/%v/%v/%v.png", queryParams.SourceID, chapterNumber, i),
				fmt.Sprintf("https://img.mghubcdn.com/file/imghub/%v/%v/%v.jpeg", queryParams.SourceID, chapterNumber, i),
				fmt.Sprintf("https://img.mghubcdn.com/file/imghub/%v/%v/%v.webp", queryParams.SourceID, chapterNumber, i),
			},
		})
	}

	chapter.SourceLink = "#"

	return chapter, nil
}

func MangasPaginate(mangas []models.Manga, page, perpage int64) []models.Manga {
	startIndex := (page - 1) * perpage
	endIndex := startIndex + perpage
	if startIndex >= int64(len(mangas)) {
		return []models.Manga{}
	}
	if endIndex >= int64(len(mangas)) {
		endIndex = int64(len(mangas) - 1)
	}
	return mangas[startIndex:endIndex]
}
