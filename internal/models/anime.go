package models

type Anime struct {
	ID                 string    `json:"id"`
	Source             string    `json:"source"`
	Title              string    `json:"title"`
	LatestEpisode      float64   `json:"latest_episode"`
	CoverUrls          []string  `json:"cover_urls"`
	Episodes           []Episode `json:"episodes"`
	OriginalLink       string    `json:"original_link"`
	ReleaseMonth       string    `json:"release_month"`
	ReleaseSeason      string    `json:"release_season"`
	ReleaseSeasonIndex int64     `json:"release_season_index"`
	ReleaseYear        int64     `json:"release_year"`
	ReleaseDate        string    `json:"release_date"`
	Score              float64   `json:"score"`
}
