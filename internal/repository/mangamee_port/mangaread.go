package mangamee_port

import (
	"context"

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
	return getSearch(ctx, models.SOURCE_MANGAREAD, 1, queryParams)
}

func (t *Mangaread) GetChapter(ctx context.Context, queryParams models.QueryParams) (models.Chapter, error) {
	return getChapter(ctx, models.SOURCE_MANGAREAD, 1, queryParams)
}
