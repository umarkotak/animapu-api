package mangamee_port

import (
	"context"
	"fmt"
	"time"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
)

type Mangaread struct{}

func NewMangaread() Mangaread {
	return Mangaread{}
}

func (t *Mangaread) GetHome(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	return getHome(ctx, models.SOURCE_MANGAREAD, 1, queryParams.Page)
}

func (t *Mangaread) GetDetail(ctx context.Context, queryParams models.QueryParams) (models.Manga, error) {
	return getDetail(ctx, models.SOURCE_MANGAREAD, 1, queryParams)
}

func (t *Mangaread) GetSearch(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
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

func (t *Mangaread) GetChapter(ctx context.Context, queryParams models.QueryParams) (models.Chapter, error) {
	return getChapter(ctx, models.SOURCE_MANGAREAD, 1, queryParams)
}
