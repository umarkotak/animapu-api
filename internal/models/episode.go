package models

type Episode struct {
	ID           string  `json:"id"`
	Number       float64 `json:"number"`
	StreamType   string  `json:"stream_type"` // Enum: hls, mp4
	RawStreamUrl string  `json:"raw_stream_url"`
	OriginalUrl  string  `json:"original_url"`
	RawPageByte  []byte  `json:"raw_page_byte"`
}

type EpisodeWatch struct {
	RawStreamUrl string `json:"raw_stream_url"`
	RawPageByte  []byte `json:"raw_page_byte"`
	IframeUrl    string `json:"iframe_url"`
}
