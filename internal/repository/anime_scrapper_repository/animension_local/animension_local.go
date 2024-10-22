package anime_scrapper_animension_local

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/local_db"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/utils/anime_utils"
	"github.com/umarkotak/animapu-api/internal/utils/utils"
	"mvdan.cc/xurls"
)

type AnimensionLocal struct {
	AnimapuSource   string
	Source          string
	AnimensionHost  string
	DesusStreamHost string
}

var (
	ExtractEpRegex = regexp.MustCompile(`Episode (\d+)`)
)

func NewAnimensionLocal() AnimensionLocal {
	return AnimensionLocal{
		AnimapuSource:  models.ANIME_SOURCE_ANIMENSION_LOCAL,
		Source:         "animension",
		AnimensionHost: "https://animension.to",
	}
}

func (r *AnimensionLocal) GetLatest(ctx context.Context, queryParams models.AnimeQueryParams) ([]models.Anime, error) {
	if queryParams.Page <= 0 {
		queryParams.Page = 1
	}

	url := fmt.Sprintf("%s/public-api/index.php?page=%v&mode=sub", r.AnimensionHost, queryParams.Page)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return []models.Anime{}, err
	}

	req.Header.Add("authority", strings.ReplaceAll(r.AnimensionHost, "https://", ""))
	req.Header.Add("accept", "*/*")
	req.Header.Add("accept-language", "en-US,en;q=0.9,id;q=0.8")
	req.Header.Add("origin", r.AnimensionHost)
	req.Header.Add("sec-ch-ua", "\"Not_A Brand\";v=\"8\", \"Chromium\";v=\"120\", \"Google Chrome\";v=\"120\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"Windows\"")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return []models.Anime{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return []models.Anime{}, err
	}

	if res.StatusCode != 200 {
		err = fmt.Errorf("error animension %v", res.StatusCode)
		logrus.WithContext(ctx).Error(err)
		return []models.Anime{}, err
	}

	// fmt.Println(string(body))

	// Sample response:
	// [
	//   [
	//       "Jashin-chan Mame Anime",                                // title
	//       3624577909,                                              // id
	//       3975867423,                                              //
	//       21,                                                      // latest eps
	//       "https://gogocdn.net/cover/jashin-chan-mame-anime.png",  // cover image
	//       1704604031
	//   ]
	// ]
	data := []any{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return []models.Anime{}, err
	}

	animes := []models.Anime{}
	for _, oneElem := range data {
		arrAnime := []any{}
		objAnime := map[string]any{}

		tmpByte, _ := json.Marshal(oneElem)

		if strings.HasPrefix(string(tmpByte), "[") {
			err = json.Unmarshal(tmpByte, &arrAnime)
		} else {
			err = json.Unmarshal(tmpByte, &objAnime)
		}
		if err != nil {
			logrus.WithContext(ctx).Error(err)
			return []models.Anime{}, err
		}

		if strings.HasPrefix(string(tmpByte), "[") {
			latestEp, _ := strconv.ParseFloat(fmt.Sprint(arrAnime[3]), 64)

			animes = append(animes, models.Anime{
				ID:            fmt.Sprintf("%v", int64(arrAnime[1].(float64))),
				Source:        r.Source,
				Title:         fmt.Sprint(arrAnime[0]),
				LatestEpisode: latestEp,
				CoverUrls:     r.animensionImages(fmt.Sprintf("%s%v", r.AnimensionHost, arrAnime[4])),
				OriginalLink:  fmt.Sprintf("%s/%v", r.AnimensionHost, int64(arrAnime[1].(float64))),
			})
		} else {
			latestEp, _ := strconv.ParseFloat(fmt.Sprint(objAnime["3"]), 64)

			animes = append(animes, models.Anime{
				ID:            fmt.Sprintf("%v", int64(objAnime["1"].(float64))),
				Source:        r.Source,
				Title:         fmt.Sprint(objAnime["0"]),
				LatestEpisode: latestEp,
				CoverUrls:     r.animensionImages(fmt.Sprintf("%s%v", r.AnimensionHost, objAnime["4"])),
				OriginalLink:  fmt.Sprintf("%s/%v", r.AnimensionHost, int64(objAnime["1"].(float64))),
			})
		}
	}

	return animes, nil
}

func (r *AnimensionLocal) GetSearchLegacy(ctx context.Context, queryParams models.AnimeQueryParams) ([]models.Anime, error) {
	q := strings.ToLower(queryParams.Title)
	animes := []models.Anime{}

	for _, oneAnime := range local_db.AnimensionAnimeIndex {
		if strings.Contains(strings.ToLower(oneAnime.Title), q) || strings.Contains(strings.ToLower(oneAnime.AltTitle), q) {
			animes = append(animes, models.Anime{
				ID:            fmt.Sprintf("%v", oneAnime.AnimensionAnimeID),
				Source:        r.Source,
				Title:         fmt.Sprint(oneAnime.Title),
				LatestEpisode: oneAnime.Episodes[len(oneAnime.Episodes)-1].EpisodeNumber,
				CoverUrls:     []string{oneAnime.CoverURL},
				OriginalLink:  fmt.Sprintf("%s/%v", r.AnimensionHost, oneAnime.AnimensionAnimeID),
			})
		}
	}

	return animes, nil
}

func (r *AnimensionLocal) GetSearch(ctx context.Context, queryParams models.AnimeQueryParams) ([]models.Anime, error) {
	animes := []models.Anime{}

	url := fmt.Sprintf("%v/public-api/search.php?search_text=%v&sort=popular-week&page=1", r.AnimensionHost, queryParams.Title)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return animes, err
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
		return animes, err
	}
	defer resp.Body.Close()

	// [
	// 	[
	// 		"Solo Leveling",
	// 		3028690795,
	// 		"https:\/\/s4.anilist.co\/file\/anilistcdn\/media\/anime\/cover\/medium\/bx151807-m1gX3iwfIsLu.png",
	// 		0,
	// 	]
	// ]

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return animes, err
	}

	res := []any{}
	d := json.NewDecoder(strings.NewReader(string(body)))
	d.UseNumber()
	err = d.Decode(&res)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return animes, err
	}

	for _, oneAnimeRes := range res {
		arrAnime := []any{}
		objAnime := map[string]any{}

		tmpByte, _ := json.Marshal(oneAnimeRes)
		d := json.NewDecoder(strings.NewReader(string(tmpByte)))
		d.UseNumber()

		if strings.HasPrefix(string(tmpByte), "[") {
			err = d.Decode(&arrAnime)
		} else {
			err = d.Decode(&objAnime)
		}
		if err != nil {
			logrus.WithContext(ctx).Error(err)
			return []models.Anime{}, err
		}

		if strings.HasPrefix(string(tmpByte), "[") {
			animes = append(animes, models.Anime{
				ID:            fmt.Sprintf("%v", arrAnime[1]),
				Source:        r.Source,
				Title:         fmt.Sprint(arrAnime[0]),
				LatestEpisode: 0,
				CoverUrls:     r.animensionImages(fmt.Sprintf("%s%v", r.AnimensionHost, arrAnime[2])),
				OriginalLink:  fmt.Sprintf("%s/%v", r.AnimensionHost, fmt.Sprintf("%v", arrAnime[1])),
			})
		} else {
			animes = append(animes, models.Anime{
				ID:            fmt.Sprintf("%v", objAnime["1"]),
				Source:        r.Source,
				Title:         fmt.Sprint(objAnime["0"]),
				LatestEpisode: 0,
				CoverUrls:     r.animensionImages(fmt.Sprintf("%s%v", r.AnimensionHost, objAnime["2"])),
				OriginalLink:  fmt.Sprintf("%s/%v", r.AnimensionHost, fmt.Sprintf("%v", objAnime["1"])),
			})
		}
	}

	return animes, nil
}

func (r *AnimensionLocal) GetDetail(ctx context.Context, queryParams models.AnimeQueryParams) (models.Anime, error) {
	cl := colly.NewCollector()

	animeDetail := models.AnimeDetail{
		AnimensionAnimeID:  queryParams.SourceID,     // done
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
		animeDetail.CoverURL = fmt.Sprintf("%s%v", r.AnimensionHost, e.Attr("src"))
		animeDetail.CoverURLs = r.animensionImages(fmt.Sprintf("%s%v", r.AnimensionHost, e.Attr("src")))
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
			AnimensionAnimeID:   queryParams.SourceID, // done
			AnimensionEpisodeID: "",                   // done
			EpisodeTitle:        "",                   // done
			EpisodeNumber:       0,                    // done
			RawHlsPlaybackURL:   "",                   // tba
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
				CoverUrl:     fmt.Sprintf("%s%v", r.AnimensionHost, h.ChildAttr("div > a > div.limit > img", "src")),
				CoverUrls:    r.animensionImages(fmt.Sprintf("%s%v", r.AnimensionHost, h.ChildAttr("div > a > div.limit > img", "src"))),
			}
			animeDetail.Relations = append(animeDetail.Relations, animeRelation)
		})
	})

	cl.OnRequest(func(r *colly.Request) {
		r.Headers.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		r.Headers.Set("accept-language", "en-US,en;q=0.9,id;q=0.8")
		r.Headers.Set("cache-control", "max-age=0")
		r.Headers.Set("cookie", "MicrosoftApplicationsTelemetryDeviceId=6ee9fc8b-8ce4-43d2-af98-ac4f7a5e97da; MicrosoftApplicationsTelemetryFirstLaunchTime=2024-04-26T00:27:41.344Z")
		r.Headers.Set("priority", "u=0, i")
		r.Headers.Set("sec-ch-ua", "\"Chromium\";v=\"124\", \"Google Chrome\";v=\"124\", \"Not-A.Brand\";v=\"99\"")
		r.Headers.Set("sec-ch-ua-mobile", "?0")
		r.Headers.Set("sec-ch-ua-platform", "\"Windows\"")
		r.Headers.Set("sec-fetch-dest", "document")
		r.Headers.Set("sec-fetch-mode", "navigate")
		r.Headers.Set("sec-fetch-site", "same-origin")
		r.Headers.Set("sec-fetch-user", "?1")
		r.Headers.Set("upgrade-insecure-requests", "1")
		r.Headers.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	})

	targetUrl := fmt.Sprintf("%v/%v", r.AnimensionHost, queryParams.SourceID)
	err := cl.Visit(targetUrl)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return models.Anime{}, err
	}

	cl.Wait()

	if !relationsExist {
		animeDetail.Relations = []models.AnimeRelation{}
	}

	episodesArr, err := r.GetEpisodes(ctx, queryParams.SourceID)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return models.Anime{}, err
	}

	for _, rawEpisode := range episodesArr {
		if len(rawEpisode) < 3 {
			continue
		}

		animeEpisode := models.AnimeEpisode{
			AnimensionAnimeID:   queryParams.SourceID,
			AnimensionEpisodeID: fmt.Sprintf("%v", rawEpisode[1]),
			EpisodeTitle:        fmt.Sprintf("Episode %v", rawEpisode[2]),
			EpisodeNumber:       utils.StringMustFloat64(fmt.Sprint(rawEpisode[2])),
			RawHlsPlaybackURL:   "",
		}

		animeDetail.Episodes = append(animeDetail.Episodes, animeEpisode)
	}

	sort.Slice(animeDetail.Episodes, func(i, j int) bool {
		return animeDetail.Episodes[i].EpisodeNumber < animeDetail.Episodes[j].EpisodeNumber
	})

	return r.animeDetailToAnime(animeDetail), nil
}

func (r *AnimensionLocal) WatchLegacy(ctx context.Context, queryParams models.AnimeQueryParams) (models.EpisodeWatch, error) {
	hlsUrl, err := r.GetHlsUrl(ctx, queryParams.EpisodeID)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return models.EpisodeWatch{}, err
	}

	return models.EpisodeWatch{
		StreamType:   "hls",
		RawStreamUrl: hlsUrl,
		OriginalUrl:  fmt.Sprintf("%s/%v#%v", r.AnimensionHost, queryParams.SourceID, queryParams.EpisodeID),
	}, nil
}

func (r *AnimensionLocal) Watch(ctx context.Context, queryParams models.AnimeQueryParams) (models.EpisodeWatch, error) {
	targetUrl := fmt.Sprintf("%s/public-api/episode.php?id=%s", r.AnimensionHost, queryParams.EpisodeID)
	req, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return models.EpisodeWatch{}, err
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
		return models.EpisodeWatch{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return models.EpisodeWatch{}, err
	}

	data := []any{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return models.EpisodeWatch{}, err
	}

	type EmbedUrls struct {
		VidCDN     string `json:"VidCDN-embed"`
		Streamwish string `json:"Streamwish-embed"`
		Mp4upload  string `json:"Mp4upload-embed"`
		Doodstream string `json:"Doodstream-embed"`
		Vidhide    string `json:"Vidhide-embed"`
	}
	embedUrls := EmbedUrls{}

	for _, val := range data {
		stringVal := fmt.Sprint(val)

		if !strings.Contains(stringVal, "VidCDN-embed") {
			continue
		}

		err = json.Unmarshal([]byte(stringVal), &embedUrls)
		if err != nil {
			logrus.WithContext(ctx).Error(err)
			return models.EpisodeWatch{}, err
		}
	}

	return models.EpisodeWatch{
		StreamType: "iframe",
		IframeUrl:  r.cleanUpUrl(embedUrls.VidCDN),
		IframeUrls: map[string]string{
			"VidCDN":     r.cleanUpUrl(embedUrls.VidCDN),
			"Streamwish": r.cleanUpUrl(embedUrls.Streamwish),
			"Mp4upload":  r.cleanUpUrl(embedUrls.Mp4upload),
			"Doodstream": r.cleanUpUrl(embedUrls.Doodstream),
			"Vidhide":    r.cleanUpUrl(embedUrls.Vidhide),
		},
		OriginalUrl: fmt.Sprintf("%s/%v#%v", r.AnimensionHost, queryParams.SourceID, queryParams.EpisodeID),
	}, nil
}

func (r *AnimensionLocal) GetPerSeason(ctx context.Context, queryParams models.AnimeQueryParams) (models.AnimePerSeason, error) {
	animes := []models.Anime{}
	animesSummary := []models.AnimeSummary{}

	for _, oneSeason := range local_db.AnimensionSeasonShorted {
		if oneSeason.Year == int(queryParams.ReleaseYear) && oneSeason.SeasonName == queryParams.ReleaseSeason {
			animesSummary = oneSeason.AnimeData
			break
		}
	}

	for _, oneAnime := range animesSummary {
		animes = append(animes, models.Anime{
			ID:            fmt.Sprintf("%v", oneAnime.AnimensionAnimeID),
			Source:        r.Source,
			Title:         oneAnime.Title,
			LatestEpisode: 0,
			CoverUrls:     r.animensionImages(oneAnime.CoverURL),
			OriginalLink:  fmt.Sprintf("%s/%v", r.AnimensionHost, oneAnime.AnimensionAnimeID),
		})
	}

	return models.AnimePerSeason{
		ReleaseYear: queryParams.ReleaseYear,
		SeasonName:  queryParams.ReleaseSeason,
		SeasonIndex: int64(anime_utils.SeasonToIndex(queryParams.ReleaseSeason)),
		Animes:      animes,
	}, nil
}

func (r *AnimensionLocal) GetRandom(ctx context.Context, queryParams models.AnimeQueryParams) ([]models.Anime, error) {
	animes := []models.Anime{}

	randomAnimes := r.getRandomElements(local_db.AnimensionAnimeIndex, 40)
	for _, oneAnime := range randomAnimes {
		animes = append(animes, models.Anime{
			ID:            fmt.Sprintf("%v", oneAnime.AnimensionAnimeID),
			Source:        r.Source,
			Title:         fmt.Sprint(oneAnime.Title),
			LatestEpisode: 0,
			CoverUrls:     r.animensionImages(oneAnime.CoverURL),
			OriginalLink:  fmt.Sprintf("%s/%v", r.AnimensionHost, oneAnime.AnimensionAnimeID),
		})
	}

	return animes, nil
}

func (r *AnimensionLocal) animeDetailToAnime(animeDetail models.AnimeDetail) models.Anime {
	episodes := []models.Episode{}
	for _, oneEp := range animeDetail.Episodes {
		episodes = append(episodes, models.Episode{
			AnimeID:      animeDetail.AnimensionAnimeID,
			Source:       r.AnimensionHost,
			ID:           oneEp.AnimensionEpisodeID,
			Number:       oneEp.EpisodeNumber,
			Title:        fmt.Sprintf("Episode %v", oneEp.EpisodeNumber),
			CoverUrl:     animeDetail.CoverURL,
			CoverUrls:    r.animensionImages(animeDetail.CoverURL),
			OriginalLink: fmt.Sprintf("%s/%v#%v", r.AnimensionHost, animeDetail.AnimensionAnimeID, oneEp.AnimensionEpisodeID),
		})
	}

	relations := []models.Anime{}
	for _, oneRelation := range animeDetail.Relations {
		relations = append(relations, models.Anime{
			ID:           oneRelation.AnimeID,
			Source:       r.AnimensionHost,
			Title:        oneRelation.Title,
			CoverUrls:    r.animensionImages(oneRelation.CoverUrl),
			OriginalLink: fmt.Sprintf("%s/%v", r.AnimensionHost, oneRelation.AnimeID),
			Relationship: oneRelation.Relationship,
		})
	}

	return models.Anime{
		ID:                 animeDetail.AnimensionAnimeID,
		Source:             r.AnimensionHost,
		Title:              animeDetail.Title,
		LatestEpisode:      animeDetail.Episodes[len(animeDetail.Episodes)-1].EpisodeNumber,
		Description:        animeDetail.Description,
		Genres:             animeDetail.Genres,
		CoverUrls:          r.animensionImages(animeDetail.CoverURL),
		Episodes:           episodes,
		OriginalLink:       fmt.Sprintf("%s/%v", r.AnimensionHost, animeDetail.AnimensionAnimeID),
		ReleaseSeason:      animeDetail.ReleaseSeasonName,
		ReleaseSeasonIndex: int64(animeDetail.ReleaseSeasonIndex),
		ReleaseYear:        int64(animeDetail.ReleaseYear),
		ReleaseMonth:       "",
		ReleaseDate:        "",
		Score:              animeDetail.MalScore,
		Relations:          relations,
	}
}

func (r *AnimensionLocal) getRandomElements(arr []models.AnimeDetail, count int) []models.AnimeDetail {
	// Check if the count is greater than the array length
	if count > len(arr) {
		count = len(arr)
	}

	// Shuffle the array using Fisher-Yates algorithm
	shuffled := make([]models.AnimeDetail, len(arr))
	copy(shuffled, arr)
	for i := len(shuffled) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	// Return the first 'count' elements
	return shuffled[:count]
}

func (r *AnimensionLocal) animensionImages(url string) []string {
	return []string{
		url,
		utils.AnimensionImgProxy(url),
		"/images/animehub_cover.jpeg",
	}
}

func (r *AnimensionLocal) cleanUpUrl(str string) string {
	return strings.ReplaceAll(str, "\u0026", "&")
}
