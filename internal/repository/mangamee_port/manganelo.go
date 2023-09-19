package mangamee_port

import (
	"context"

	"github.com/umarkotak/animapu-api/internal/models"
)

type MangaNelo struct {
	AnimapuSourceID  string
	MangameeSourceID string
}

func NewMangaNelo() MangaNelo {
	return MangaNelo{
		AnimapuSourceID:  models.SOURCE_MANGANELO,
		MangameeSourceID: "manganelo",
	}
}

func (t MangaNelo) GetHome(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	return getHome(ctx, t.AnimapuSourceID, t.MangameeSourceID, queryParams.Page)
}

func (t MangaNelo) GetDetail(ctx context.Context, queryParams models.QueryParams) (models.Manga, error) {
	return getDetail(ctx, t.AnimapuSourceID, t.MangameeSourceID, queryParams)
}

func (t MangaNelo) GetSearch(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	return getSearch(ctx, t.AnimapuSourceID, t.MangameeSourceID, queryParams)
}

func (t MangaNelo) GetChapter(ctx context.Context, queryParams models.QueryParams) (models.Chapter, error) {
	return getChapter(ctx, t.AnimapuSourceID, t.MangameeSourceID, queryParams)
}
