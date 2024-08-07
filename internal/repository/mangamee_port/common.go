package mangamee_port

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
	"github.com/umarkotak/animapu-api/internal/utils/utils"
)

type (
	MangameeHomeResponse struct {
		Code    int64           `json:"code"`
		Message string          `json:"message"`
		Data    []MangameeManga `json:"data"`
	}
	MangameeDetailResponse struct {
		Code    int64         `json:"code"`
		Message string        `json:"message"`
		Data    MangameeManga `json:"data"`
	}

	MangameeManga struct {
		Id             string            `json:"id"`
		Cover          string            `json:"cover"`
		Title          string            `json:"title"`
		LastChapter    string            `json:"last_chapter"`
		LastRead       string            `json:"last_read"`
		Status         string            `json:"status"`
		Summary        string            `json:"summary"`
		Chapters       []MangameeChapter `json:"chapters"`
		DataImages     MangameeDataImage `json:"data_images"`
		ChapterName    string            `json:"chapter_name"`
		Images         []MangameeImage   `json:"images"`
		OriginalServer string            `json:"original_server"`
	}
	MangameeDataImage struct {
		ChapterName string          `json:"chapter_name"`
		Images      []MangameeImage `json:"images"`
	}
	MangameeChapter struct {
		Id          string `json:"id"`
		Name        string `json:"name"`
		ChapterName string `json:"chapter_name"`
	}
	MangameeImage struct {
		Image string `json:"image"`
	}
)

func getHome(ctx context.Context, animapuSource string, mangameeSource string, page int64) ([]models.Manga, error) {
	url := fmt.Sprintf("%v/api/manga/index?source=%v&page=%v", config.Get().MangameeApiHost, mangameeSource, page)

	req, _ := http.NewRequest(
		"GET", url, strings.NewReader("{}"),
	)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return []models.Manga{}, err
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	if res.StatusCode != 200 {
		err = models.ErrMangamee
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"response_body": string(body),
		}).Error(err)
		return []models.Manga{}, err
	}

	var mangameeHomeResponse MangameeHomeResponse
	json.Unmarshal(body, &mangameeHomeResponse)

	// PORT mangamee home to animapu home
	mangas := []models.Manga{}

	for _, mangameeData := range mangameeHomeResponse.Data {
		chapterNumberString := utils.RemoveNonNumeric(mangameeData.LastChapter)

		mangas = append(mangas, models.Manga{
			ID:                  mangameeData.Id,
			SourceID:            mangameeData.Id,
			Source:              animapuSource,
			Title:               mangameeData.Title,
			Status:              "Ongoing",
			Rating:              "0",
			LatestChapterID:     "chapter_id",
			LatestChapterNumber: utils.StringMustFloat64(chapterNumberString),
			LatestChapterTitle:  mangameeData.LastChapter,
			Chapters:            []models.Chapter{},
			CoverImages: []models.CoverImage{
				{
					Index: 1,
					ImageUrls: []string{
						fmt.Sprintf(mangameeData.Cover),
					},
				},
			},
		})
	}

	return mangas, nil
}

func getDetail(ctx context.Context, animapuSource string, mangameeSource string, queryParams models.QueryParams) (models.Manga, error) {
	url := fmt.Sprintf("%v/api/manga/detail?source=%v&mangaid=%v", config.Get().MangameeApiHost, mangameeSource, queryParams.SourceID)

	req, _ := http.NewRequest(
		"GET", url, strings.NewReader("{}"),
	)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return models.Manga{}, err
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	if res.StatusCode != 200 {
		err = models.ErrMangamee
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"response_body": string(body),
		}).Error(err)
		return models.Manga{}, err
	}

	var mangameeDetailResponse MangameeDetailResponse
	json.Unmarshal(body, &mangameeDetailResponse)

	// logrus.Infof("URL: %+v", url)
	// logrus.Infof("MANGAMEE DETAIL: %+v", mangameeDetailResponse)

	chapterNumber := utils.StringMustFloat64(utils.RemoveNonNumeric(mangameeDetailResponse.Data.Chapters[0].Name))

	manga := models.Manga{
		ID:                  queryParams.SourceID,
		Source:              animapuSource,
		SourceID:            queryParams.SourceID,
		Title:               mangameeDetailResponse.Data.Title,
		Description:         mangameeDetailResponse.Data.Summary,
		Genres:              []string{},
		Status:              "Ongoing",
		CoverImages:         []models.CoverImage{{ImageUrls: []string{mangameeDetailResponse.Data.Cover}}},
		Chapters:            []models.Chapter{},
		LatestChapterID:     mangameeDetailResponse.Data.Chapters[0].Id,
		LatestChapterNumber: chapterNumber,
		LatestChapterTitle:  mangameeDetailResponse.Data.Chapters[0].Name,
	}

	idx := int64(1)
	for _, mangameeChapter := range mangameeDetailResponse.Data.Chapters {
		manga.Chapters = append(manga.Chapters, models.Chapter{
			ID:       mangameeChapter.Id,
			Source:   animapuSource,
			SourceID: mangameeChapter.Id,
			Title:    mangameeChapter.Name,
			Index:    idx,
			Number:   utils.StringMustFloat64(utils.RemoveNonNumeric(mangameeChapter.Name)),
		})
		idx += 1
	}

	return manga, nil
}

func getSearch(ctx context.Context, animapuSource string, mangameeSource string, queryParams models.QueryParams) ([]models.Manga, error) {
	url := fmt.Sprintf(
		"%v/api/manga/search?source=%v&title=%v",
		config.Get().MangameeApiHost, mangameeSource, strings.Replace(queryParams.Title, " ", "%20", -1),
	)

	req, _ := http.NewRequest(
		"GET", url, strings.NewReader("{}"),
	)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return []models.Manga{}, err
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	if res.StatusCode != 200 {
		err = models.ErrMangamee
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"response_body": string(body),
		}).Error(err)
		return []models.Manga{}, err
	}

	var mangameeSearchResponse MangameeHomeResponse
	json.Unmarshal(body, &mangameeSearchResponse)

	mangas := []models.Manga{}

	for _, mangameeData := range mangameeSearchResponse.Data {
		chapterNumberString := utils.RemoveNonNumeric(mangameeData.LastChapter)

		mangas = append(mangas, models.Manga{
			ID:                  mangameeData.Id,
			SourceID:            mangameeData.Id,
			Source:              animapuSource,
			Title:               mangameeData.Title,
			Status:              "Ongoing",
			Rating:              "0",
			LatestChapterID:     "chapter_id",
			LatestChapterNumber: utils.StringMustFloat64(chapterNumberString),
			LatestChapterTitle:  mangameeData.LastChapter,
			Chapters:            []models.Chapter{},
			CoverImages: []models.CoverImage{
				{
					Index: 1,
					ImageUrls: []string{
						fmt.Sprintf(mangameeData.Cover),
					},
				},
			},
		})
	}

	return mangas, nil
}

func getChapter(ctx context.Context, animapuSource string, mangameeSource string, queryParams models.QueryParams) (models.Chapter, error) {
	url := fmt.Sprintf("%v/api/manga/image?source=%v&mangaid=%v&chapterid=%v", config.Get().MangameeApiHost, mangameeSource, queryParams.SourceID, queryParams.ChapterID)

	req, _ := http.NewRequest(
		"GET", url, strings.NewReader("{}"),
	)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return models.Chapter{}, err
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	if res.StatusCode != 200 {
		err = models.ErrMangamee
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"response_body": string(body),
		}).Error(err)
		return models.Chapter{}, err
	}

	var mangameeDetailResponse MangameeDetailResponse
	json.Unmarshal(body, &mangameeDetailResponse)

	chapter := models.Chapter{
		ID:            queryParams.ChapterID,
		SourceID:      queryParams.SourceID,
		Source:        animapuSource,
		Number:        utils.StringMustFloat64(utils.RemoveNonNumeric(mangameeDetailResponse.Data.ChapterName)),
		ChapterImages: []models.ChapterImage{},
		SourceLink:    mangameeDetailResponse.Data.OriginalServer,
	}

	idx := int64(1)
	for _, image := range mangameeDetailResponse.Data.Images {
		chapter.ChapterImages = append(chapter.ChapterImages, models.ChapterImage{
			Index:     idx,
			ImageUrls: []string{image.Image},
		})
		idx += 1
	}

	return chapter, nil
}
