package models

type Episode struct {
	AnimeID      string  `json:"anime_id"`
	Source       string  `json:"source"`
	ID           string  `json:"id"`
	Number       float64 `json:"number"`
	Title        string  `json:"title"`
	OriginalLink string  `json:"original_link"`
}

type EpisodeWatch struct {
	StreamType   string   `json:"stream_type"` // Enum: hls, mp4, iframe
	RawStreamUrl string   `json:"raw_stream_url"`
	RawPageByte  []byte   `json:"raw_page_byte"`
	IframeUrl    string   `json:"iframe_url"`
	IframeUrls   []string `json:"iframe_urls"`
	OriginalUrl  string   `json:"original_url"`
}
