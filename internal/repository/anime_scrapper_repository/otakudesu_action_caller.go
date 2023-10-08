package anime_scrapper_repository

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (s *Otakudesu) AdminAjaxCaller(action string, additionals []string) ([]byte, error) {
	url := fmt.Sprintf("%v/wp-admin/admin-ajax.php", s.OtakudesuHost)
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(
		`action=%v&%v`, action, strings.Join(additionals, "&"),
	))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return nil, err
	}
	req.Header.Add("authority", "otakudesu.wiki")
	req.Header.Add("accept", "*/*")
	req.Header.Add("accept-language", "en-US,en;q=0.9,id;q=0.8")
	req.Header.Add("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("cookie", "_ga=GA1.2.1861696737.1695554920; _gid=GA1.2.526398125.1696081168; _gat=1; _ga_025LZFQCB2=GS1.2.1696081167.2.1.1696082456.0.0.0")
	req.Header.Add("origin", fmt.Sprintf("%v", s.OtakudesuHost))
	req.Header.Add("referer", fmt.Sprintf("%v/episode/mt-ithd-s2-episode-1-sub-indo/", s.OtakudesuHost))
	req.Header.Add("sec-ch-ua", "\"Google Chrome\";v=\"117\", \"Not;A=Brand\";v=\"8\", \"Chromium\";v=\"117\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"Windows\"")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")
	req.Header.Add("x-requested-with", "XMLHttpRequest")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
