package models

import "context"

type (
	AnimeScrapper interface {
		GetLatest(ctx context.Context, queryParams AnimeQueryParams) ([]Anime, error)
		GetSearch(ctx context.Context, queryParams AnimeQueryParams) ([]Anime, error)
		GetRandom(ctx context.Context, queryParams AnimeQueryParams) ([]Anime, error)
		GetDetail(ctx context.Context, queryParams AnimeQueryParams) (Anime, error)
		Watch(ctx context.Context, queryParams AnimeQueryParams) (EpisodeWatch, error)
		GetPerSeason(ctx context.Context, queryParams AnimeQueryParams) (AnimePerSeason, error)
	}
)
