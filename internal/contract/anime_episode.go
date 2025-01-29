package contract

type (
	Episode struct {
		AnimeID      string   `json:"anime_id"`
		Source       string   `json:"source"`
		ID           string   `json:"id"`
		Number       float64  `json:"number"`
		Title        string   `json:"title"`
		OriginalLink string   `json:"original_link"`
		CoverUrl     string   `json:"cover_url"`
		CoverUrls    []string `json:"cover_urls"`
		UseTitle     bool     `json:"use_title"`
	}

	EpisodeWatch struct {
		StreamType    string            `json:"stream_type"`    // Enum: hls, mp4, iframe
		RawStreamUrl  string            `json:"raw_stream_url"` //
		RawPageByte   []byte            `json:"raw_page_byte"`  //
		IframeUrl     string            `json:"iframe_url"`     //
		IframeUrls    map[string]string `json:"iframe_urls"`    //
		OriginalUrl   string            `json:"original_url"`   //
		StreamOptions []StreamOption    `json:"stream_options"` //
		Resolution    string            `json:"resolution"`     //
		StreamIdx     string            `json:"stream_idx"`     //
	}

	StreamOption struct {
		Resolution string `json:"resolution"` // Enum: 720p, 480p, 360p
		Index      string `json:"index"`
		Name       string `json:"name"`
		Used       bool   `json:"used"`
	}
)
