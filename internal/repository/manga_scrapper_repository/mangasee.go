package manga_scrapper_repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/config"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/utils/utils"
)

type (
	MangaseeManga struct {
		SeriesID   string    `json:"SeriesID"`
		IndexName  string    `json:"IndexName"`
		SeriesName string    `json:"SeriesName"`
		ScanStatus string    `json:"ScanStatus"`
		Chapter    string    `json:"Chapter"`
		Genres     string    `json:"Genres"`
		Date       time.Time `json:"Date"`
		IsEdd      bool      `json:"IsEdd"`
	}

	MangaseeChapter struct {
		Chapter     string      `json:"Chapter"`
		Type        string      `json:"Type"`
		Date        string      `json:"Date"`
		ChapterName interface{} `json:"ChapterName"`
	}

	MangaseeSearchManga struct {
		IndexName  string   `json:"i"`
		SeriesName string   `json:"s"`
		AltNames   []string `json:"a"`
	}
)

type Mangasee struct {
	Host    string
	Source  string
	ImgHost string
}

func NewMangasee() Mangasee {
	return Mangasee{
		Host:    "https://www.mangasee123.com",
		Source:  "mangasee",
		ImgHost: "https://temp.compsci88.com",
	}
}

func (sc *Mangasee) GetHome(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	c := colly.NewCollector()
	// c.SetRequestTimeout(config.Get().CollyTimeout)
	c.SetRequestTimeout(config.Get().CollyTimeout)
	// t := &http.Transport{
	// 	Dial: (&net.Dialer{
	// 		Timeout:   60 * time.Second,
	// 		KeepAlive: 30 * time.Second,
	// 	}).Dial,
	// 	TLSHandshakeTimeout: 60 * time.Second,
	// }
	// c.WithTransport(t)

	mangas := []models.Manga{}

	c.OnHTML("body > script:nth-child(16)", func(e *colly.HTMLElement) {
		footerContent := e.Text
		splitted := strings.Split(footerContent, "vm.LatestJSON = ")
		if len(splitted) <= 0 {
			return
		}

		splitted = strings.Split(splitted[1], "vm.NewSeriesJSON")
		dataJson := splitted[0]
		dataJson = strings.ReplaceAll(dataJson, ";", "")
		dataJson = strings.TrimSpace(dataJson)

		mangaseeMangas := []MangaseeManga{}
		json.Unmarshal([]byte(dataJson), &mangaseeMangas)

		for _, oneMangaseeManga := range mangaseeMangas {
			chNumber := mangaseeDecodeCh(oneMangaseeManga.Chapter)

			mangas = append(mangas, models.Manga{
				ID:                  oneMangaseeManga.IndexName,
				SourceID:            oneMangaseeManga.IndexName,
				Source:              sc.Source,
				Title:               oneMangaseeManga.SeriesName,
				Genres:              []string{},
				LatestChapterID:     "",
				LatestChapterNumber: utils.ForceSanitizeStringToFloat(chNumber),
				LatestChapterTitle:  chNumber,
				Chapters:            []models.Chapter{},
				CoverImages: []models.CoverImage{
					{
						Index: 1,
						ImageUrls: []string{
							fmt.Sprintf("%v/cover/%v.jpg", sc.ImgHost, oneMangaseeManga.IndexName),
						},
					},
				},
			})
		}
	})

	err := c.Visit(fmt.Sprintf("%v", sc.Host))
	c.Wait()
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return mangas, err
	}

	return mangas, nil
}

func (sc *Mangasee) GetDetail(ctx context.Context, queryParams models.QueryParams) (models.Manga, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(config.Get().CollyTimeout)
	c.AllowURLRevisit = true

	manga := models.Manga{
		ID:          queryParams.SourceID,
		Source:      sc.Source,
		SourceID:    queryParams.SourceID,
		Title:       strings.ReplaceAll(queryParams.SourceID, "-", " "),
		Description: "Description unavailable",
		Genres:      []string{},
		Status:      "Ongoing",
		CoverImages: []models.CoverImage{{ImageUrls: []string{
			fmt.Sprintf("%v/cover/%v.jpg", sc.ImgHost, queryParams.SourceID),
		}}},
		Chapters: []models.Chapter{},
	}

	c.OnHTML("div.top-5.Content", func(e *colly.HTMLElement) {
		manga.Description = e.Text
	})

	c.OnHTML("body > script:nth-child(16)", func(e *colly.HTMLElement) {
		footerContent := e.Text
		splitted := strings.Split(footerContent, "vm.Chapters = ")
		if len(splitted) <= 0 {
			return
		}

		splitted = strings.Split(splitted[1], "vm.NumSubs")
		dataJson := splitted[0]
		dataJson = strings.ReplaceAll(dataJson, ";", "")
		dataJson = strings.TrimSpace(dataJson)

		mangaseeChapters := []MangaseeChapter{}
		json.Unmarshal([]byte(dataJson), &mangaseeChapters)

		for i, oneMangaseeChapter := range mangaseeChapters {
			firstNum := oneMangaseeChapter.Chapter[0:1]
			lastNum := oneMangaseeChapter.Chapter[len(oneMangaseeChapter.Chapter)-1:]
			chNumberS := mangaseeDecodeCh(oneMangaseeChapter.Chapter)
			if lastNum != "0" {
				chNumberS = fmt.Sprintf("%v.%v", chNumberS, lastNum)
			}

			chNumer := utils.ForceSanitizeStringToFloat(chNumberS)

			manga.Chapters = append(manga.Chapters, models.Chapter{
				ID:                fmt.Sprint(chNumer),
				Source:            sc.Source,
				SourceID:          fmt.Sprint(chNumer),
				SecondarySourceID: fmt.Sprint(firstNum),
				Title:             fmt.Sprintf("%v %v", oneMangaseeChapter.Type, chNumer),
				Index:             int64(i),
				Number:            chNumer,
			})
		}
	})

	targetUrl := fmt.Sprintf("%v/manga/%v", sc.Host, queryParams.SourceID)
	err := c.Visit(targetUrl)
	c.Wait()
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"target_url": targetUrl,
		}).Error(err)
		return manga, err
	}

	manga.GenerateLatestChapter()

	return manga, nil
}

func (sc *Mangasee) GetSearch(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/_search.php", sc.Host), nil)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return []models.Manga{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return []models.Manga{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return []models.Manga{}, err
	}

	mangaseeSearchDatas := []MangaseeSearchManga{}
	err = json.Unmarshal(body, &mangaseeSearchDatas)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return []models.Manga{}, err
	}

	mangas := []models.Manga{}
	for _, oneMangaseeSearch := range mangaseeSearchDatas {
		if !strings.Contains(strings.ToLower(oneMangaseeSearch.SeriesName), strings.ToLower(queryParams.Title)) {
			continue
		}

		mangas = append(mangas, models.Manga{
			ID:                  oneMangaseeSearch.IndexName,
			SourceID:            oneMangaseeSearch.IndexName,
			Source:              sc.Source,
			Title:               oneMangaseeSearch.SeriesName,
			Genres:              []string{},
			LatestChapterID:     "",
			LatestChapterNumber: 0,
			LatestChapterTitle:  "",
			Chapters:            []models.Chapter{},
			CoverImages: []models.CoverImage{
				{
					Index: 1,
					ImageUrls: []string{
						fmt.Sprintf("%v/cover/%v.jpg", sc.ImgHost, oneMangaseeSearch.IndexName),
					},
				},
			},
		})
	}

	return mangas, nil
}

func (sc *Mangasee) GetChapter(ctx context.Context, queryParams models.QueryParams) (models.Chapter, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(10 * time.Minute)
	// t := &http.Transport{
	// 	Dial: (&net.Dialer{
	// 		Timeout:   60 * time.Second,
	// 		KeepAlive: 30 * time.Second,
	// 	}).Dial,
	// 	TLSHandshakeTimeout: 60 * time.Second,
	// }
	// c.WithTransport(t)

	targetLink := fmt.Sprintf("%v/read-online/%v-chapter-%v.html", sc.Host, queryParams.SourceID, queryParams.ChapterID)
	if queryParams.SecondarySourceID == "2" {
		targetLink = fmt.Sprintf("%v/read-online/%v-chapter-%v-index-2.html", sc.Host, queryParams.SourceID, queryParams.ChapterID)
	}

	chapter := models.Chapter{
		ID:            queryParams.ChapterID,
		SourceID:      queryParams.SourceID,
		Source:        sc.Source,
		SourceLink:    targetLink,
		Number:        utils.ForceSanitizeStringToFloat(queryParams.ChapterID),
		ChapterImages: []models.ChapterImage{},
	}

	c.OnHTML("body > script:nth-child(19)", func(e *colly.HTMLElement) {
		re := regexp.MustCompile(`vm\.CurPathName\s*=\s*("[^"]+")`)

		// Find the first match of the pattern
		imageHost := re.FindStringSubmatch(e.Text)

		chFloat := utils.ForceSanitizeStringToFloat(queryParams.ChapterID)
		chInt := int(chFloat)
		modifier := ""
		if (chFloat - float64(chInt)) > 0 {
			splitted := strings.Split(queryParams.ChapterID, ".")
			modifier = fmt.Sprintf(".%v", splitted[1])
		}

		currChString := ""
		splittedCurrCh := strings.Split(e.Text, "vm.CurChapter = ")
		if len(splittedCurrCh) >= 2 {
			splittedCurrCh = strings.Split(splittedCurrCh[1], ";")
			if len(splittedCurrCh) > 0 {
				currChString = splittedCurrCh[0]
			}
		}
		type MangaseeChapter struct {
			Chapter     string  `json:"Chapter"`
			Type        string  `json:"Type"`
			Page        string  `json:"Page"`
			Directory   string  `json:"Directory"`
			Date        string  `json:"Date"`
			ChapterName *string `json:"ChapterName"`
		}
		mangaseeChapter := MangaseeChapter{}
		json.Unmarshal([]byte(currChString), &mangaseeChapter)
		pageInt, _ := strconv.ParseInt(mangaseeChapter.Page, 10, 54)
		if pageInt == 0 {
			pageInt = 150
		}

		dir := ""
		if mangaseeChapter.Directory != "" {
			dir = fmt.Sprintf("%s/", dir)
		}

		// https://{{vm.CurPathName}}/manga/Dandadan/{{vm.CurChapter.Directory == '' ? '' : vm.CurChapter.Directory+'/'}}{{vm.ChapterImage(vm.CurChapter.Chapter)}}-{{vm.PageImage(Page)}}.png

		for i := 1; i <= int(pageInt); i++ {
			chapter.ChapterImages = append(chapter.ChapterImages, models.ChapterImage{
				Index: 0,
				ImageUrls: []string{
					fmt.Sprintf(
						"https://%v/manga/%v/%v%04d%v-%03d.png",
						strings.ReplaceAll(imageHost[1], `"`, ""), queryParams.SourceID, dir, chInt, modifier, i,
					),
					// fmt.Sprintf(
					// 	"https://%v/manga/%v/%04d%v-%03d.png",
					// 	strings.ReplaceAll(imageHost[1], `"`, ""), queryParams.SourceID, chInt, modifier, i,
					// ),
					// fmt.Sprintf(
					// 	"https://%v/manga/%v/Mag-Official/%04d%v-%03d.png",
					// 	strings.ReplaceAll(imageHost[1], `"`, ""), queryParams.SourceID, chInt, modifier, i,
					// ),
				},
			})
		}
	})

	var err error
	for i := 0; i < 5; i++ {
		err = c.Visit(targetLink)
		c.Wait()
		if err != nil {
			logrus.WithContext(ctx).Error(err)
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}
	if err != nil {
		return chapter, err
	}

	return chapter, nil
}

func mangaseeDecodeCh(s string) string {
	return trimFirstRune(trimLastChar(s))
}

func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}

func trimLastChar(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return s[:len(s)-size]
}
