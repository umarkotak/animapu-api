package manga_scrapper

import (
	"context"
	"fmt"
	"time"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
)

// ANIMAPU LITE API TEMPLATE GUIDELINE
// When adding new manga source, please follow this guideline in order to comply
// with the Front End Apps requirements.
// You can copy paste this template to a new file then modify it.
func GetMangasourceLatestManga(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
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
		Chapters:            []models.Chapter{},
		CoverImages: []models.CoverImage{
			{
				Index: 1,
				ImageUrls: []string{
					fmt.Sprintf("https://animapu-lite.vercel.app/images/manga/%v", "image_id"),
				},
			},
		},
	})

	err := c.Visit(fmt.Sprintf("https://animapu-lite.vercel.app/home/%v", queryParams.Page))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return mangas, err
	}

	return mangas, nil
}

func GetMangasourceDetailManga(ctx context.Context, queryParams models.QueryParams) (models.Manga, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	manga := models.Manga{
		ID:                  queryParams.SourceID,
		Source:              "source",
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

	err := c.Visit(fmt.Sprintf("https://animapu-lite.vercel.app/manga/%v", queryParams.SourceID))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return manga, err
	}

	return manga, nil
}

func GetMangasourceByQuery(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
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
		Chapters:            []models.Chapter{},
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

func GetMangasourceDetailChapter(ctx context.Context, queryParams models.QueryParams) (models.Chapter, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	chapter := models.Chapter{
		ID:            queryParams.ChapterID,
		SourceID:      queryParams.SourceID,
		Source:        "fizmanga",
		Number:        0,
		ChapterImages: []models.ChapterImage{},
	}

	err := c.Visit(fmt.Sprintf("https://animapu-lite.vercel.app/manga/%v/chapter/%v", queryParams.SourceID, queryParams.ChapterID))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return chapter, err
	}

	return chapter, nil
}
