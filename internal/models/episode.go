package models

type Episode struct {
	ID           string  `json:"id"`
	Number       float64 `json:"number"`
	StreamType   string  `json:"stream_type"` // Enum: hls, mp4
	RawStreamUrl string  `json:"raw_stream_url"`
}
