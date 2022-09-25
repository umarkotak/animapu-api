package mangamee_port

import (
	"context"

	"github.com/umarkotak/animapu-api/internal/models"
)

type Mangabat struct{}

func NewMangabat() Mangabat {
	return Mangabat{}
}

func (t Mangabat) GetHome(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	return getHome(ctx, models.SOURCE_MANGABAT, 3, queryParams.Page)
}

func (t Mangabat) GetDetail(ctx context.Context, queryParams models.QueryParams) (models.Manga, error) {
	return getDetail(ctx, models.SOURCE_MANGABAT, 3, queryParams)
}

func (t Mangabat) GetSearch(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	return getSearch(ctx, models.SOURCE_MANGABAT, 3, queryParams)
}

func (t Mangabat) GetChapter(ctx context.Context, queryParams models.QueryParams) (models.Chapter, error) {
	return getChapter(ctx, models.SOURCE_MANGABAT, 3, queryParams)
}
