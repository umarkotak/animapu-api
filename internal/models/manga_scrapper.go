package models

import "context"

type (
	MangaScrapper interface {
		GetHome(ctx context.Context, queryParams QueryParams) ([]Manga, error)
		GetDetail(ctx context.Context, queryParams QueryParams) (Manga, error)
		GetSearch(ctx context.Context, queryParams QueryParams) ([]Manga, error)
		GetChapter(ctx context.Context, queryParams QueryParams) (Chapter, error)
	}
)
