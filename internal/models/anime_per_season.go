package models

type AnimePerSeason struct {
	ReleaseYear int64   `json:"release_year"`
	SeasonName  string  `json:"season_name"`
	SeasonIndex int64   `json:"season_index"`
	Animes      []Anime `json:"animes"`
}
