package mangamee_port

import (
	"context"

	"github.com/umarkotak/animapu-api/internal/models"
)

type Mangasee struct {
	AnimapuSourceID  string
	MangameeSourceID string
}

func NewMangasee() Mangasee {
	return Mangasee{
		AnimapuSourceID:  models.SOURCE_ASURA_COMIC,
		MangameeSourceID: "mangasee",
	}
}

func (t Mangasee) GetHome(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	return getHome(ctx, t.AnimapuSourceID, t.MangameeSourceID, queryParams.Page)
}

func (t Mangasee) GetDetail(ctx context.Context, queryParams models.QueryParams) (models.Manga, error) {
	return getDetail(ctx, t.AnimapuSourceID, t.MangameeSourceID, queryParams)
}

func (t Mangasee) GetSearch(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	return getSearch(ctx, t.AnimapuSourceID, t.MangameeSourceID, queryParams)
}

func (t Mangasee) GetChapter(ctx context.Context, queryParams models.QueryParams) (models.Chapter, error) {
	return getChapter(ctx, t.AnimapuSourceID, t.MangameeSourceID, queryParams)
}
