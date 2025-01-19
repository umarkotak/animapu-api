package mal_api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/nstratos/go-myanimelist/mal"
	"github.com/sirupsen/logrus"
)

type clientIDTransport struct {
	Transport http.RoundTripper
	ClientID  string
}

func (c *clientIDTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if c.Transport == nil {
		c.Transport = http.DefaultTransport
	}
	req.Header.Add("X-MAL-CLIENT-ID", c.ClientID)
	return c.Transport.RoundTrip(req)
}

func GetSeasonalAnime(ctx context.Context, year int, season string) ([]mal.Anime, error) {
	// Validate season
	validSeasons := map[string]bool{
		"winter": true,
		"spring": true,
		"summer": true,
		"fall":   true,
	}
	if !validSeasons[season] {
		return nil, fmt.Errorf("invalid season: %s. Must be winter, spring, summer, or fall", season)
	}

	publicInfoClient := &http.Client{
		// Create client ID from https://myanimelist.net/apiconfig.
		Transport: &clientIDTransport{ClientID: "0082178c6dc4da4903b91c5d1ef1fdbe"},
	}

	// Create client
	c := mal.NewClient(publicInfoClient)

	animes, _, err := c.Anime.Seasonal(
		ctx, year, mal.AnimeSeason(season), mal.SortSeasonalByAnimeNumListUsers, mal.Limit(500),
		mal.Fields{"synopsis", "num_episodes", "alternative_titles"},
	)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return animes, err
	}

	return animes, nil
}
