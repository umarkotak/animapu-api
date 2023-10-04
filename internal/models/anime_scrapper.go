package models

import "context"

type (
	AnimeScrapper interface {
		Watch(ctx context.Context, queryParams AnimeQueryParams) (EpisodeWatch, error)
	}
)
