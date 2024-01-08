package anime_scrapper_animension_local

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/sirupsen/logrus"
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

	client := &http.Client{}
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
	data := [][]interface{}{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return []models.Anime{}, err
	}

	animes := []models.Anime{}
	for _, oneElem := range data {
		latestEp, _ := strconv.ParseFloat(fmt.Sprint(oneElem[3]), 64)

		animes = append(animes, models.Anime{
			ID:            fmt.Sprintf("%v", int64(oneElem[1].(float64))),
			Source:        r.Source,
			Title:         fmt.Sprint(oneElem[0]),
			LatestEpisode: latestEp,
			CoverUrls:     []string{fmt.Sprint(oneElem[4])},
			OriginalLink:  fmt.Sprintf("%s/%v", r.AnimensionHost, int64(oneElem[1].(float64))),
		})
	}

	return animes, nil
}

func (r *AnimensionLocal) GetSearch(ctx context.Context, queryParams models.AnimeQueryParams) ([]models.Anime, error) {
	return []models.Anime{}, nil
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
		animeDetail.CoverURL = e.Attr("src")
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
				CoverUrl:     h.ChildAttr("div > a > div.limit > img", "src"),
			}
			animeDetail.Relations = append(animeDetail.Relations, animeRelation)
		})
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

	return r.AnimeDetailToAnime(animeDetail), nil
}

func (r *AnimensionLocal) Watch(ctx context.Context, queryParams models.AnimeQueryParams) (models.EpisodeWatch, error) {
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

func (r *AnimensionLocal) GetPerSeason(ctx context.Context, queryParams models.AnimeQueryParams) (models.AnimePerSeason, error) {
	return models.AnimePerSeason{}, nil
}

func (r *AnimensionLocal) GetRandom(ctx context.Context, queryParams models.AnimeQueryParams) ([]models.Anime, error) {
	return []models.Anime{}, nil
}

func (r *AnimensionLocal) AnimeDetailToAnime(animeDetail models.AnimeDetail) models.Anime {
	episodes := []models.Episode{}
	for _, oneEp := range animeDetail.Episodes {
		episodes = append(episodes, models.Episode{
			AnimeID:      animeDetail.AnimensionAnimeID,
			Source:       r.AnimensionHost,
			ID:           oneEp.AnimensionEpisodeID,
			Number:       oneEp.EpisodeNumber,
			Title:        fmt.Sprintf("Episode %v", oneEp.EpisodeNumber),
			CoverUrl:     animeDetail.CoverURL,
			OriginalLink: fmt.Sprintf("%s/%v#%v", r.AnimensionHost, animeDetail.AnimensionAnimeID, oneEp.AnimensionEpisodeID),
		})
	}

	relations := []models.Anime{}
	for _, oneRelation := range animeDetail.Relations {
		relations = append(relations, models.Anime{
			ID:           oneRelation.AnimeID,
			Source:       r.AnimensionHost,
			Title:        oneRelation.Title,
			CoverUrls:    []string{oneRelation.CoverUrl},
			OriginalLink: fmt.Sprintf("%s/%v", r.AnimensionHost, oneRelation.AnimeID),
			Relationship: oneRelation.Relationship,
		})
	}

	return models.Anime{
		ID:                 animeDetail.AnimensionAnimeID,
		Source:             r.AnimensionHost,
		Title:              animeDetail.Title,
		LatestEpisode:      animeDetail.Episodes[len(animeDetail.Episodes)-1].EpisodeNumber,
		CoverUrls:          []string{animeDetail.CoverURL},
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
