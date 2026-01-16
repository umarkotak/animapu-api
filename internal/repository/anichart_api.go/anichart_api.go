package anichart_api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// AniList API endpoint
const (
	aniListAPI = "https://graphql.anilist.co"

	// GraphQL query to get anime by year and season
	animeByYearAndSeasonQuery = `query ($year: Int, $season: MediaSeason) {
		Page(page: 1, perPage: 50) { # Adjust perPage as needed
			media(type: ANIME, season: $season, seasonYear: $year, sort: POPULARITY_DESC) {
				id
				title {
					romaji
					english
					native
				}
				genres
				averageScore
				episodes
				status
				season
				seasonYear
				coverImage {
					large
				}
				description
				externalLinks {
					site
					url
				}
			}
		}
	}`
)

// AniList response structure
type (
	AniListResponse struct {
		Data struct {
			Page struct {
				Media []Anime `json:"media"`
			} `json:"Page"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	// Anime structure
	Anime struct {
		ID    int `json:"id"`
		Title struct {
			Romaji  string `json:"romaji"`
			English string `json:"english"`
			Native  string `json:"native"`
		} `json:"title"`
		Genres       []string `json:"genres"`
		AverageScore int      `json:"averageScore"`
		Episodes     int      `json:"episodes"`
		Status       string   `json:"status"`
		Season       string   `json:"season"`
		SeasonYear   int      `json:"seasonYear"`
		CoverImage   struct {
			Large string `json:"large"`
		} `json:"coverImage"`
		Description   string `json:"description"`
		ExternalLinks []struct {
			Site string `json:"site"`
			Url  string `json:"url"`
		} `json:"externalLinks"`
	}
)

func GetSeasonalAnime(ctx context.Context, year int, season string) ([]Anime, error) {
	variables := map[string]any{
		"year":   year,
		"season": strings.ToUpper(season),
	}

	requestBody, err := json.Marshal(map[string]any{
		"query":     animeByYearAndSeasonQuery,
		"variables": variables,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, aniListAPI, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var aniListResponse AniListResponse
	err = json.Unmarshal(body, &aniListResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w, Body: %s", err, string(body))
	}

	if len(aniListResponse.Errors) > 0 {
		return nil, fmt.Errorf("AniList API returned errors: %v", aniListResponse.Errors)
	}

	return aniListResponse.Data.Page.Media, nil
}
