package anime_scrapper_animension_local

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
)

func (r *AnimensionLocal) GetHlsUrl(ctx context.Context, epid string) (string, error) {
	targetUrl := fmt.Sprintf("%s/public-api/episode.php?id=%s", r.AnimensionHost, epid)
	req, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return "", err
	}

	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "en-US,en;q=0.9,id;q=0.8")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"116\", \"Not)A;Brand\";v=\"24\", \"Google Chrome\";v=\"116\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return "", err
	}

	var data []interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return "", err
	}

	hlsUrl := ""

	for _, val := range data {
		tryParse := func(s string) (bool, string) {
			if strings.Contains(s, ".m3u8") {
				var insideData map[string]interface{}
				err := json.Unmarshal([]byte(s), &insideData)
				if err != nil {
					logrus.WithContext(ctx).Error(err)
					return false, fmt.Sprintf("%s %v", s, err)
				}

				tmpUrl, err := url.Parse(insideData["VidCDN-embed"].(string))
				if err != nil {
					logrus.WithContext(ctx).Error(err)
					return false, fmt.Sprintf("%s %v", s, err)
				}

				hlsUrl = tmpUrl.Path

				hlsUrl = strings.TrimPrefix(hlsUrl, "/")
				hlsUrl = strings.TrimSuffix(hlsUrl, ".php")
				hlsUrl = fmt.Sprintf("https://%s", hlsUrl)

				return true, ""
			}

			return false, ""
		}

		if ok, errMsg := tryParse(fmt.Sprint(val)); ok {
			return hlsUrl, nil
		} else if errMsg != "" {
			logrus.WithContext(ctx).Error(fmt.Errorf(errMsg))
		}
	}

	return "", fmt.Errorf("HLS URL not found")
}
