package manga_scrapper

import (
	"context"
	"fmt"
	"time"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
)

// This is only a template file to be easily copy-pasted

func GetFizmangaLatestManga(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
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
