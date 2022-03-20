package models

type (
	Manga struct {
		ID                 string    `json:"id"`
		SourceID           string    `json:"source_id"`
		SecondarySourceID  string    `json:"secondary_source_id"`
		Source             string    `json:"source"`
		Title              string    `json:"title"`
		Description        string    `json:"description"`
		Genres             []string  `json:"genres"`
		Status             string    `json:"status"`
		Rating             string    `json:"rating"`
		LatestChapterID    string    `json:"latest_chapter_id"`
		LatestChapterTitle string    `json:"latest_chapter_title"`
		Chapters           []Chapter `json:"chapters"`
	}
)
