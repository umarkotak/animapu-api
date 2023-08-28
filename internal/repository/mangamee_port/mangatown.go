package mangamee_port

import (
	"context"

	"github.com/umarkotak/animapu-api/internal/models"
)

type Mangatown struct{}

func NewMangatown() Mangatown {
	return Mangatown{}
}

func (t Mangatown) GetHome(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	return getHome(ctx, models.SOURCE_MANGATOWN, "mangatown", queryParams.Page)
}

func (t Mangatown) GetDetail(ctx context.Context, queryParams models.QueryParams) (models.Manga, error) {
	return getDetail(ctx, models.SOURCE_MANGATOWN, "mangatown", queryParams)
}

func (t Mangatown) GetSearch(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	return getSearch(ctx, models.SOURCE_MANGATOWN, "mangatown", queryParams)
}

func (t Mangatown) GetChapter(ctx context.Context, queryParams models.QueryParams) (models.Chapter, error) {
	return getChapter(ctx, models.SOURCE_MANGATOWN, "mangatown", queryParams)
}
