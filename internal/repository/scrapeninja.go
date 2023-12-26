package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/config"
	"github.com/umarkotak/animapu-api/internal/models"
)

func QuickScrape(ctx context.Context, targetUrl string) (models.ScrapeNinjaResponse, error) {
	var err error
	var res *http.Response
	scrapeNinjaResponse := models.ScrapeNinjaResponse{}

	scrapeUrl := fmt.Sprintf("%v/scrape", config.Get().ScrapeNinjaConfig.Host)

	for _, rapidApiKey := range config.Get().ScrapeNinjaConfig.RapidApiKeys {
		params := map[string]interface{}{
			"url": targetUrl,
		}
		payload, _ := json.Marshal(params)
		req, _ := http.NewRequest(
			"POST", scrapeUrl, strings.NewReader(string(payload)),
		)
		req.Header.Add("content-type", "application/json")
		req.Header.Add("X-RapidAPI-Host", "scrapeninja.p.rapidapi.com")

		req.Header.Del("X-RapidAPI-Key")
		req.Header.Add("X-RapidAPI-Key", rapidApiKey)

		res, err = http.DefaultClient.Do(req)

		if err != nil || res.StatusCode != 200 {
			continue
		}
		break
	}
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return scrapeNinjaResponse, err
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	json.Unmarshal(body, &scrapeNinjaResponse)
	return scrapeNinjaResponse, nil
}
