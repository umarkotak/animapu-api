package mangamee_port

import (
	"context"

	"github.com/umarkotak/animapu-api/internal/models"
)

type AsuraComic struct {
	AnimapuSourceID  string
	MangameeSourceID string
}

func NewAsuraComic() AsuraComic {
	return AsuraComic{
		AnimapuSourceID:  models.SOURCE_ASURA_COMIC,
		MangameeSourceID: "asuracomic",
	}
}

func (t AsuraComic) GetHome(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	return getHome(ctx, t.AnimapuSourceID, t.MangameeSourceID, queryParams.Page)
}

func (t AsuraComic) GetDetail(ctx context.Context, queryParams models.QueryParams) (models.Manga, error) {
	return getDetail(ctx, t.AnimapuSourceID, t.MangameeSourceID, queryParams)
}

func (t AsuraComic) GetSearch(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	return getSearch(ctx, t.AnimapuSourceID, t.MangameeSourceID, queryParams)
}

func (t AsuraComic) GetChapter(ctx context.Context, queryParams models.QueryParams) (models.Chapter, error) {
	return getChapter(ctx, t.AnimapuSourceID, t.MangameeSourceID, queryParams)
}
