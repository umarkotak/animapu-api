package mangamee_port

import (
	"context"

	"github.com/umarkotak/animapu-api/internal/models"
)

type MangaNato struct {
	AnimapuSourceID  string
	MangameeSourceID string
}

func NewMangaNato() MangaNato {
	return MangaNato{
		AnimapuSourceID:  models.SOURCE_MANGANATO,
		MangameeSourceID: "manganato",
	}
}

func (t MangaNato) GetHome(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	return getHome(ctx, t.AnimapuSourceID, t.MangameeSourceID, queryParams.Page)
}

func (t MangaNato) GetDetail(ctx context.Context, queryParams models.QueryParams) (models.Manga, error) {
	return getDetail(ctx, t.AnimapuSourceID, t.MangameeSourceID, queryParams)
}

func (t MangaNato) GetSearch(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	return getSearch(ctx, t.AnimapuSourceID, t.MangameeSourceID, queryParams)
}

func (t MangaNato) GetChapter(ctx context.Context, queryParams models.QueryParams) (models.Chapter, error) {
	return getChapter(ctx, t.AnimapuSourceID, t.MangameeSourceID, queryParams)
}
