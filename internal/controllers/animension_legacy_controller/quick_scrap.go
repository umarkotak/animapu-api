package animension_legacy_controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/utils/anime_utils"
	"github.com/umarkotak/animapu-api/internal/utils/utils"
	"mvdan.cc/xurls"
)

func quickScrapAnimeDetail(ctx context.Context, params ReqBody) (models.AnimeDetail, error) {
	cl := colly.NewCollector()

	animeDetail := models.AnimeDetail{
		AnimensionAnimeID:  params.AnimeID,           // done
		MasterTitle:        "",                       // done
		MasterTitleTag:     "",                       // done
		Title:              "",                       // done
		AltTitle:           "",                       // done
		Description:        "",                       // done
		Status:             "",                       // done
		Genres:             []string{},               // done
		SeasonLabel:        "",                       // tba
		SeasonIndex:        1,                        // tba
		ReleaseYear:        0,                        // done
		ReleaseSeasonName:  "",                       // done
		ReleaseSeasonIndex: 0,                        // done
		CoverURL:           "",                       // done
		HeaderCoverURL:     "",                       // done
		MalScore:           0,                        // done
		TotalEpisode:       0,                        // done
		Type:               "",                       // done
		Episodes:           []models.AnimeEpisode{},  // done
		Relations:          []models.AnimeRelation{}, // tba
		VideoSources:       []models.VideoSource{},   // tba
		LastSyncAt:         time.Now(),
	}

	cl.OnHTML("#content > div > div.postbody > div > div.bixbox.animefull.animefull-bixbox > div.bigcontent > div.infox > h1", func(e *colly.HTMLElement) {
		animeDetail.Title = e.Text
	})
	cl.OnHTML("#content > div > div.postbody > div > div.bixbox.animefull.animefull-bixbox > div.bigcontent > div.infox > h2", func(e *colly.HTMLElement) {
		animeDetail.AltTitle = e.Text
	})
	cl.OnHTML("#content > div > div.postbody > div > div.bixbox.animefull.animefull-bixbox > div.bigcontent > div.infox > div > div > div.desc", func(e *colly.HTMLElement) {
		animeDetail.Description = e.Text
	})
	cl.OnHTML("#content > div > div.postbody > div > div.bixbox.animefull.animefull-bixbox > div.bigcontent > div.infox > div > div > div.spe > span:nth-child(1)", func(e *colly.HTMLElement) {
		animeDetail.Status = strings.TrimSpace(strings.ReplaceAll(strings.ToLower(e.Text), "status:", ""))
	})
	cl.OnHTML("#content > div > div.postbody > div > div.bixbox.animefull.animefull-bixbox > div.bigcontent > div.infox > div > div > div.genxed > span > a", func(e *colly.HTMLElement) {
		animeDetail.Genres = append(animeDetail.Genres, strings.ToLower(e.Text))
	})
	cl.OnHTML("#content > div > div.postbody > div > div.bixbox.animefull.animefull-bixbox > div.bigcontent > div.infox > div > div > div.spe > span:nth-child(2) > a", func(e *colly.HTMLElement) {
		splitted := strings.Split(strings.ToLower(e.Text), " ")
		if len(splitted) < 2 {
			return
		}
		animeDetail.ReleaseSeasonName = splitted[0]
		releaseYear, _ := strconv.ParseInt(splitted[1], 10, 64)
		animeDetail.ReleaseYear = int(releaseYear)
		animeDetail.ReleaseSeasonIndex = anime_utils.SeasonToIndex(animeDetail.ReleaseSeasonName)
	})
	cl.OnHTML("#thumbook > div.thumb > img", func(e *colly.HTMLElement) {
		animeDetail.CoverURL = e.Attr("src")

		u, err := url.Parse(animeDetail.CoverURL)
		if err != nil || u.Host == "" {
			animeDetail.CoverURL = fmt.Sprintf("%s/%v", AnimensionHost, animeDetail.CoverURL)
		}
	})
	cl.OnHTML("#bigcover > div", func(e *colly.HTMLElement) {
		animeDetail.HeaderCoverURL = xurls.Relaxed.FindString(e.Attr("style"))
	})
	cl.OnHTML("#thumbook > div.rt > div.rating > strong", func(e *colly.HTMLElement) {
		animeDetail.MalScore = utils.ForceSanitizeStringToFloat(e.Text)
	})
	cl.OnHTML("#content > div > div.postbody > div > div.bixbox.animefull.animefull-bixbox > div.bigcontent > div.infox > div > div > div.spe > span:nth-child(3)", func(e *colly.HTMLElement) {
		animeDetail.TotalEpisode = utils.StringMustInt64(utils.RemoveNonNumeric(e.Text))
	})
	cl.OnHTML("#anime_episodes > ul > li", func(e *colly.HTMLElement) {
		animeEpisode := models.AnimeEpisode{
			AnimensionAnimeID:   params.AnimeID, // done
			AnimensionEpisodeID: "",             // done
			EpisodeTitle:        "",             // done
			EpisodeNumber:       0,              // done
			RawHlsPlaybackURL:   "",             // tba
		}

		match := ExtractEpRegex.FindStringSubmatch(e.ChildAttr("div.sli-name > a", "href"))
		if len(match) > 1 {
			animeEpisode.EpisodeNumber = utils.StringMustFloat64(match[1])
			animeEpisode.EpisodeTitle = fmt.Sprintf("Episode %v", animeEpisode.EpisodeNumber)
		}

		animeEpisode.AnimensionEpisodeID = strings.ReplaceAll(e.ChildAttr("div.sli-btn > a", "id"), "episode-", "")

		animeDetail.Episodes = append(animeDetail.Episodes, animeEpisode)
	})

	relationsExist := false
	cl.OnHTML("#content > div > div.postbody > div > div:nth-child(4) > div.releases > h3 > span", func(e *colly.HTMLElement) {
		if strings.ToLower(e.Text) == "relations" {
			relationsExist = true
		}
	})

	cl.OnHTML("#content > div > div.postbody > div > div:nth-child(4) > div.listupd", func(e *colly.HTMLElement) {
		e.ForEach("article", func(i int, h *colly.HTMLElement) {
			animeRelation := models.AnimeRelation{
				AnimeID:      strings.ReplaceAll(h.ChildAttr("div > a", "href"), "/", ""),
				Relationship: strings.ToLower(h.ChildText("div > a > div.limit > div.bt > span")),
				Title:        strings.TrimSpace(h.ChildText("div > a > div.tt")),
				CoverUrl:     h.ChildAttr("div > a > div.limit > img", "src"),
			}
			animeDetail.Relations = append(animeDetail.Relations, animeRelation)
		})
	})

	targetUrl := fmt.Sprintf("%v/%v", AnimensionHost, params.AnimeID)
	err := cl.Visit(targetUrl)
	if err != nil {
		return models.AnimeDetail{}, err
	}

	cl.Wait()

	if !relationsExist {
		animeDetail.Relations = []models.AnimeRelation{}
	}

	episodesArr, err := getAnimensionEpisodes(params.AnimeID)
	if err != nil {
		return models.AnimeDetail{}, err
	}
	for _, rawEpisode := range episodesArr {
		if len(rawEpisode) < 3 {
			continue
		}

		animeEpisode := models.AnimeEpisode{
			AnimensionAnimeID:   params.AnimeID,                                     // done
			AnimensionEpisodeID: fmt.Sprintf("%v", rawEpisode[1]),                   // done
			EpisodeTitle:        fmt.Sprintf("Episode %v", rawEpisode[2]),           // done
			EpisodeNumber:       utils.StringMustFloat64(fmt.Sprint(rawEpisode[2])), // done
			RawHlsPlaybackURL:   "",                                                 // done
		}

		hlsUrl, err := getHlsUrl(animeEpisode.AnimensionEpisodeID)
		if err != nil {
			logrus.Errorf("Episode error: %v - %v", animeEpisode.AnimensionEpisodeID, err)
		}
		animeEpisode.RawHlsPlaybackURL = hlsUrl

		animeDetail.Episodes = append(animeDetail.Episodes, animeEpisode)
	}

	animeDetail.GenerateDefault()

	sort.Slice(animeDetail.Episodes, func(i, j int) bool {
		return animeDetail.Episodes[i].EpisodeNumber < animeDetail.Episodes[j].EpisodeNumber
	})

	return animeDetail, nil
}

func getAnimensionEpisodes(animeID string) ([][]interface{}, error) {
	res := [][]interface{}{}

	url := fmt.Sprintf("%v/public-api/episodes.php?id=%v", AnimensionHost, animeID)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return res, err
	}
	req.Header.Add("authority", AnimensionBase)
	req.Header.Add("accept", "*/*")
	req.Header.Add("accept-language", "en-US,en;q=0.9,id;q=0.8")
	req.Header.Add("cookie", "token=2296020162393900497; username=umarkotak; id=39682; loggedin=1")
	req.Header.Add("origin", AnimensionHost)
	req.Header.Add("sec-ch-ua", "\"Chromium\";v=\"118\", \"Google Chrome\";v=\"118\", \"Not=A?Brand\";v=\"99\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	d := json.NewDecoder(strings.NewReader(string(body)))
	d.UseNumber()
	err = d.Decode(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func getHlsUrl(epid string) (string, error) {
	targetUrl := fmt.Sprintf("%s/public-api/episode.php?id=%s", AnimensionHost, epid)
	req, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		logrus.Error(err)
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
		logrus.WithFields(logrus.Fields{"ep_id": epid}).Error(err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(logrus.Fields{"ep_id": epid}).Error(err)
		return "", err
	}

	var data []interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		logrus.WithFields(logrus.Fields{"ep_id": epid}).Error(err)
		return "", err
	}

	hlsUrl := ""

	for _, val := range data {
		tryParse := func(s string) (bool, string) {
			if strings.Contains(s, ".m3u8") {
				var insideData map[string]interface{}
				err := json.Unmarshal([]byte(s), &insideData)
				if err != nil {
					return false, fmt.Sprintf("%s %v", s, err)
				}

				tmpUrl, err := url.Parse(insideData["VidCDN-embed"].(string))
				if err != nil {
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
			fmt.Println(errMsg)
		}
	}

	return "", fmt.Errorf("HLS URL not found")
}

func getAnimensionAnimesBySeason(ctx context.Context, seasonID string, page int64) ([]int64, error) {
	animeIDs := []int64{}

	url := fmt.Sprintf("%s/public-api/search.php?season=%s&sort=popular-week&page=%v", AnimensionHost, seasonID, page)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return animeIDs, err
	}
	req.Header.Add("authority", AnimensionBase)
	req.Header.Add("accept", "*/*")
	req.Header.Add("accept-language", "en-US,en;q=0.9,id;q=0.8")
	req.Header.Add("cookie", "token=2296020162393900497; username=umarkotak; id=39682; loggedin=1")
	req.Header.Add("origin", AnimensionHost)
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
		return animeIDs, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return animeIDs, err
	}

	rawAnimes := []any{}
	err = json.Unmarshal(body, &rawAnimes)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return animeIDs, err
	}

	// logrus.WithContext(ctx).Infof("RESULT STRUCT: %+v\n", rawAnimes)

	for _, rawAnime := range rawAnimes {
		switch animeMap := rawAnime.(type) {
		case map[string]any:
			tmpByte, _ := json.Marshal(animeMap)

			tmpAnime := struct {
				ID int64 `json:"1"`
			}{}
			json.Unmarshal(tmpByte, &tmpAnime)

			animeIDs = append(animeIDs, tmpAnime.ID)

		case []any:
			tmpByte, _ := json.Marshal(animeMap)

			tmpAnime := []any{}
			d := json.NewDecoder(strings.NewReader(string(tmpByte)))
			d.UseNumber()
			d.Decode(&tmpAnime)

			animeIDs = append(animeIDs, utils.StringMustInt64(fmt.Sprint(tmpAnime[1])))
		}
	}

	return animeIDs, nil
}
