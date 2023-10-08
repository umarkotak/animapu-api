package models

type Anime struct {
	ID            string    `json:"id"`
	Source        string    `json:"source"`
	Title         string    `json:"title"`
	LatestEpisode float64   `json:"latest_episode"`
	CoverUrls     []string  `json:"cover_urls"`
	Episodes      []Episode `json:"episodes"`
	OriginalLink  string    `json:"original_link"`
}
