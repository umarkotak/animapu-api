package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/config"
	"github.com/umarkotak/animapu-api/internal/models"
)

func QuickScrape(ctx context.Context, targetUrl string) (models.ScrapeNinjaResponse, error) {
	var err error
	var res *http.Response
	scrapeNinjaResponse := models.ScrapeNinjaResponse{}

	scrapeUrl := fmt.Sprintf("%v/scrape", config.Get().ScrapeNinjaConfig.Host)
	params := map[string]interface{}{
		"url": targetUrl,
	}
	payload, _ := json.Marshal(params)
	req, _ := http.NewRequest(
		"POST", scrapeUrl, strings.NewReader(string(payload)),
	)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("X-RapidAPI-Host", "scrapeninja.p.rapidapi.com")

	for _, rapidApiKey := range config.Get().ScrapeNinjaConfig.RapidApiKeys {
		req.Header.Del("X-RapidAPI-Key")
		req.Header.Add("X-RapidAPI-Key", rapidApiKey)

		res, err = http.DefaultClient.Do(req)

		if err != nil {
			continue
		}
		break
	}
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return scrapeNinjaResponse, err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	json.Unmarshal(body, &scrapeNinjaResponse)
	return scrapeNinjaResponse, nil
}
