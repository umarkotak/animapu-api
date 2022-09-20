package mangamee_port

import (
	"context"

	"github.com/umarkotak/animapu-api/internal/models"
)

type Maidmy struct{}

func NewMaidmy() Maidmy {
	return Maidmy{}
}

func (t Maidmy) GetHome(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	return getHome(ctx, models.SOURCE_MAIDMY, 4, queryParams.Page)
}

func (t Maidmy) GetDetail(ctx context.Context, queryParams models.QueryParams) (models.Manga, error) {
	return getDetail(ctx, models.SOURCE_MAIDMY, 4, queryParams)
}

func (t Maidmy) GetSearch(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	mangas := []models.Manga{}

	return mangas, nil
}

func (t Maidmy) GetChapter(ctx context.Context, queryParams models.QueryParams) (models.Chapter, error) {
	return getChapter(ctx, models.SOURCE_MAIDMY, 4, queryParams)
}
