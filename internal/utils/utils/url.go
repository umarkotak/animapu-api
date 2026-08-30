package utils

import (
	"net/http"
	"strings"
)

func GetFinalURL(inputURL string) (string, error) {
	resp, err := http.Get(inputURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return strings.Trim(resp.Request.URL.String(), "/"), nil
}
