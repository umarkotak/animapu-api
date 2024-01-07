package anime_scrapper_animension_local

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

// sample resp:
// [[
//
//	208726,
//	2201723854, -- episode id
//	5,          -- episode number
//	1699015324
//
// ]]
func (r *AnimensionLocal) GetEpisodes(ctx context.Context, animeID string) ([][]interface{}, error) {
	res := [][]interface{}{}

	url := fmt.Sprintf("%v/public-api/episodes.php?id=%v", r.AnimensionHost, animeID)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return res, err
	}
	req.Header.Add("authority", strings.ReplaceAll(r.AnimensionHost, "https://", ""))
	req.Header.Add("accept", "*/*")
	req.Header.Add("accept-language", "en-US,en;q=0.9,id;q=0.8")
	req.Header.Add("origin", r.AnimensionHost)
	req.Header.Add("sec-ch-ua", "\"Chromium\";v=\"118\", \"Google Chrome\";v=\"118\", \"Not=A?Brand\";v=\"99\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return res, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return res, err
	}

	d := json.NewDecoder(strings.NewReader(string(body)))
	d.UseNumber()
	err = d.Decode(&res)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return res, err
	}

	return res, nil
}
