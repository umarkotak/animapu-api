package models

type ScrapeNinjaResponse struct {
	Info ScrapeNinjaInfo `json:"info"`
	Body string          `json:"body"`
}

type ScrapeNinjaInfo struct {
	Version       string            `json:"version"`
	StatusCode    int64             `json:"statusCode"`
	StatusMessage string            `json:"statusMessage"`
	Headers       map[string]string `json:"headers"`
}
