package models

type (
	Manga struct {
		ID                  string       `json:"id"`
		SourceID            string       `json:"source_id"`
		SecondarySourceID   string       `json:"secondary_source_id"`
		Source              string       `json:"source"`
		SecondarySource     string       `json:"secondary_source"`
		Title               string       `json:"title"`
		Description         string       `json:"description"`
		Genres              []string     `json:"genres"`
		Status              string       `json:"status"`
		Rating              string       `json:"rating"`
		LatestChapterID     string       `json:"latest_chapter_id"`
		LatestChapterNumber float64      `json:"latest_chapter_number"`
		LatestChapterTitle  string       `json:"latest_chapter_title"`
		Chapters            []Chapter    `json:"chapters"`
		CoverImages         []CoverImage `json:"cover_image"`
	}

	CoverImage struct {
		Index     int64    `json:"index"`
		ImageUrls []string `json:"image_urls"`
	}
)
