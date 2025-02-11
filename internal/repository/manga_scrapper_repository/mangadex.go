package manga_scrapper_repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/config"
	"github.com/umarkotak/animapu-api/internal/contract"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/utils/utils"
)

type (
	Mangadex struct {
		Host string
	}

	Response struct {
		Data []struct {
			ID         string `json:"id"`
			Attributes struct {
				Title         map[string]string `json:"title"`
				LatestChapter string            `json:"latestChapter"`
			} `json:"attributes"`
			Relationships []struct {
				Type       string         `json:"type"`
				ID         string         `json:"id"`
				Attributes map[string]any `json:"attributes"`
			} `json:"relationships"`
		} `json:"data"`
	}
)

func NewMangadex() Mangadex {
	return Mangadex{
		Host: "https://api.mangadex.org",
	}
}

func (t *Mangadex) GetHome(ctx context.Context, queryParams models.QueryParams) ([]contract.Manga, error) {
	limit := 20
	offset := (int(queryParams.Page) - 1) * limit

	queries := []string{
		fmt.Sprintf("limit=%v", limit),
		fmt.Sprintf("offset=%v", offset),
		"includes[]=cover_art",
		"contentRating[]=safe",
		"contentRating[]=suggestive",
		"contentRating[]=erotica",
		"originalLanguage[]=ja",
		"originalLanguage[]=ko",
		"availableTranslatedLanguage[]=en",
		"availableTranslatedLanguage[]=id",
		"order[createdAt]=desc",
		"includedTagsMode=AND",
		"excludedTagsMode=OR",
	}
	url := fmt.Sprintf("%v/manga?%v", t.Host, strings.Join(queries, "&"))
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	mangas := []contract.Manga{}

	for _, item := range result.Data {
		title, exists := item.Attributes.Title["en"]

		if exists {
			var coverImage, latestChapter string
			for _, rel := range item.Relationships {
				if rel.Type == "cover_art" {
					coverImage = fmt.Sprintf("https://uploads.mangadex.org/covers/%v/%v.256.jpg", item.ID, rel.Attributes["fileName"])
				}
			}
			latestChapter = item.Attributes.LatestChapter

			mangas = append(mangas, contract.Manga{
				ID:                  "",
				SourceID:            "",
				Source:              "source",
				Title:               title,
				Description:         "Description unavailable",
				Genres:              []string{},
				Status:              "Ongoing",
				Rating:              "",
				LatestChapterID:     latestChapter,
				LatestChapterNumber: utils.StringMustFloat64(latestChapter),
				LatestChapterTitle:  latestChapter,
				Chapters:            []contract.Chapter{},
				CoverImages: []contract.CoverImage{
					{
						Index: 1,
						ImageUrls: []string{
							coverImage,
						},
					},
				},
			})
		}
	}

	return mangas, nil
}

func (t *Mangadex) GetDetail(ctx context.Context, queryParams models.QueryParams) (contract.Manga, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(config.Get().CollyTimeout)

	manga := contract.Manga{
		ID:                  queryParams.SourceID,
		Source:              "source",
		SourceID:            queryParams.SourceID,
		Title:               "Untitled",
		Description:         "Description unavailable",
		Genres:              []string{},
		Status:              "Ongoing",
		CoverImages:         []contract.CoverImage{{ImageUrls: []string{}}},
		Chapters:            []contract.Chapter{},
		LatestChapterID:     "",
		LatestChapterNumber: 0,
		LatestChapterTitle:  "",
	}

	err := c.Visit(fmt.Sprintf("https://animapu-lite.vercel.app/manga/%v", queryParams.SourceID))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return manga, err
	}

	return manga, nil
}

func (t *Mangadex) GetSearch(ctx context.Context, queryParams models.QueryParams) ([]contract.Manga, error) {
	limit := 64
	offset := (int(queryParams.Page) - 1) * limit

	queries := []string{
		fmt.Sprintf("limit=%v", limit),
		fmt.Sprintf("offset=%v", offset),
		fmt.Sprintf("title=%v", queryParams.Title),
		"includes[]=cover_art",
		"contentRating[]=safe",
		"contentRating[]=suggestive",
		"contentRating[]=erotica",
		"originalLanguage[]=ja",
		"originalLanguage[]=ko",
		"availableTranslatedLanguage[]=en",
		"availableTranslatedLanguage[]=id",
		"order[createdAt]=desc",
		"includedTagsMode=AND",
		"excludedTagsMode=OR",
	}
	url := fmt.Sprintf("%v/manga?%v", t.Host, strings.Join(queries, "&"))
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	mangas := []contract.Manga{}

	for _, item := range result.Data {
		title, exists := item.Attributes.Title["en"]

		if exists {
			var coverImage, latestChapter string
			for _, rel := range item.Relationships {
				if rel.Type == "cover_art" {
					coverImage = fmt.Sprintf("https://uploads.mangadex.org/covers/%v/%v.256.jpg", item.ID, rel.Attributes["fileName"])
				}
			}
			latestChapter = item.Attributes.LatestChapter

			mangas = append(mangas, contract.Manga{
				ID:                  "",
				SourceID:            "",
				Source:              "source",
				Title:               title,
				Description:         "Description unavailable",
				Genres:              []string{},
				Status:              "Ongoing",
				Rating:              "",
				LatestChapterID:     latestChapter,
				LatestChapterNumber: utils.StringMustFloat64(latestChapter),
				LatestChapterTitle:  latestChapter,
				Chapters:            []contract.Chapter{},
				CoverImages: []contract.CoverImage{
					{
						Index: 1,
						ImageUrls: []string{
							coverImage,
						},
					},
				},
			})
		}
	}

	return mangas, nil
}

func (t *Mangadex) GetChapter(ctx context.Context, queryParams models.QueryParams) (contract.Chapter, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(config.Get().CollyTimeout)

	chapter := contract.Chapter{
		ID:            queryParams.ChapterID,
		SourceID:      queryParams.SourceID,
		Source:        "fizmanga",
		Number:        0,
		ChapterImages: []contract.ChapterImage{},
	}

	err := c.Visit(fmt.Sprintf("https://animapu-lite.vercel.app/manga/%v/chapter/%v", queryParams.SourceID, queryParams.ChapterID))
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return chapter, err
	}

	return chapter, nil
}
