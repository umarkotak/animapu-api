package models

type (
	Chapter struct {
		ID            string         `json:"id"`
		SourceID      string         `json:"source_id"`
		Source        string         `json:"source"`
		Title         string         `json:"title"`
		Index         int64          `json:"index"`
		Number        float64        `json:"number"`
		ChapterImages []ChapterImage `json:"chapter_images"`
	}

	ChapterImage struct {
		Index     int64    `json:"index"`
		ImageUrls []string `json:"image_urls"`
	}
)
