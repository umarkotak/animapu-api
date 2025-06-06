package manga_scrapper_repository

import (
	"context"
	"fmt"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/config"
	"github.com/umarkotak/animapu-api/internal/contract"
	"github.com/umarkotak/animapu-api/internal/models"
)

type Template struct{}

func NewTemplate() Template {
	return Template{}
}

func (t *Template) GetHome(ctx context.Context, queryParams models.QueryParams) ([]contract.Manga, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(config.Get().CollyTimeout)

	mangas := []contract.Manga{}

	mangas = append(mangas, contract.Manga{
		ID:                  "",
		SourceID:            "",
		Source:              "source",
		Title:               "Untitled",
		Description:         "Description unavailable",
		Genres:              []string{},
		Status:              "Ongoing",
		Rating:              "10",
		LatestChapterID:     "chapter_id",
		LatestChapterNumber: 0,
		LatestChapterTitle:  "Chapter 0",
		Chapters:            []contract.Chapter{},
		CoverImages: []contract.CoverImage{
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

func (t *Template) GetDetail(ctx context.Context, queryParams models.QueryParams) (contract.Manga, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(config.Get().CollyTimeout)

	manga := contract.Manga{
		ID:                  queryParams.SourceID,
		Source:              "source",
		SourceID:            queryParams.SourceID,
		Title:               "Untitled",
		Description:         "Description unavailable",
		Genres:              []string{},
		Status:              "Ongoing",
		CoverImages:         []contract.CoverImage{{ImageUrls: []string{}}},
		Chapters:            []contract.Chapter{},
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

func (t *Template) GetSearch(ctx context.Context, queryParams models.QueryParams) ([]contract.Manga, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(config.Get().CollyTimeout)

	mangas := []contract.Manga{}

	mangas = append(mangas, contract.Manga{
		ID:                  "",
		SourceID:            "",
		Source:              "source",
		Title:               "Untitled",
		Description:         "Description unavailable",
		Genres:              []string{},
		Status:              "Ongoing",
		Rating:              "10",
		LatestChapterID:     "chapter_id",
		LatestChapterNumber: 0,
		LatestChapterTitle:  "Chapter 0",
		Chapters:            []contract.Chapter{},
		CoverImages: []contract.CoverImage{
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

func (t *Template) GetChapter(ctx context.Context, queryParams models.QueryParams) (contract.Chapter, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(config.Get().CollyTimeout)

	chapter := contract.Chapter{
		ID:            queryParams.ChapterID,
		SourceID:      queryParams.SourceID,
		Source:        "fizmanga",
		Number:        0,
		ChapterImages: []contract.ChapterImage{},
	}

	err := c.Visit(fmt.Sprintf("https://animapu-lite.vercel.app/manga/%v/chapter/%v", queryParams.SourceID, queryParams.ChapterID))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return chapter, err
	}

	return chapter, nil
}
