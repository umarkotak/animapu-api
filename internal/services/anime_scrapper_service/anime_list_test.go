package anime_scrapper_service

import (
	"context"
	"errors"
	"testing"

	"github.com/umarkotak/animapu-api/internal/contract"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository/anime_scrapper_repository"
)

type listScrapper struct {
	anime_scrapper_repository.Template
}

func (*listScrapper) GetLatest(_ context.Context, qp models.AnimeQueryParams) ([]contract.Anime, error) {
	return []contract.Anime{{ID: "latest", LatestEpisode: float64(qp.Page)}}, nil
}
func (*listScrapper) GetCompleted(_ context.Context, qp models.AnimeQueryParams) ([]contract.Anime, error) {
	return []contract.Anime{{ID: "completed", LatestEpisode: float64(qp.Page)}}, nil
}
func (*listScrapper) GetAll(_ context.Context, qp models.AnimeQueryParams) ([]contract.Anime, error) {
	return []contract.Anime{{ID: "all", LatestEpisode: float64(qp.Page)}}, nil
}

func TestAnimeListMethod(t *testing.T) {
	for _, tc := range []struct{ filter, want string }{{"", "latest"}, {"completed", "completed"}, {"all", "all"}} {
		t.Run(tc.want, func(t *testing.T) {
			method, err := animeListMethod(&listScrapper{}, tc.filter)
			if err != nil {
				t.Fatal(err)
			}
			result, err := method(context.Background(), models.AnimeQueryParams{Page: 2})
			if err != nil || len(result) != 1 || result[0].ID != tc.want || result[0].LatestEpisode != 2 {
				t.Fatalf("result = %#v, error = %v", result, err)
			}
		})
	}
	if _, err := animeListMethod(&listScrapper{}, "invalid"); !errors.Is(err, models.ErrInvalidFormat) {
		t.Fatalf("invalid filter error = %v", err)
	}
	if _, _, err := GetList(context.Background(), models.AnimeQueryParams{SpecialFilter: "invalid"}); !errors.Is(err, models.ErrInvalidFormat) {
		t.Fatalf("invalid list request error = %v", err)
	}
}

func TestAnimeListCacheKeys(t *testing.T) {
	keys := map[string]bool{}
	for _, filter := range []string{"", "completed", "all"} {
		for _, page := range []int64{1, 2} {
			qp := models.AnimeQueryParams{Source: "kuramanime", SpecialFilter: filter, Page: page}
			key := qp.ToKey("GetList")
			if keys[key] {
				t.Fatalf("duplicate cache key: %s", key)
			}
			keys[key] = true
		}
	}
}

func TestAnimeCapabilities(t *testing.T) {
	for _, source := range []string{models.ANIME_SOURCE_KURAMANIME, models.ANIME_SOURCE_OTAKUDESU} {
		capabilities, err := GetCapabilities(context.Background(), models.AnimeQueryParams{Source: source})
		if err != nil || len(capabilities) != 1 {
			t.Fatalf("%s capabilities = %#v, error = %v", source, capabilities, err)
		}
		capability := capabilities[0]
		if !capability.GetLatest || !capability.GetSearch || !capability.GetDetail || !capability.Watch {
			t.Fatalf("%s missing implemented capabilities: %#v", source, capability)
		}
		if capability.GetCompleted != (source == models.ANIME_SOURCE_KURAMANIME) || capability.GetAll != (source == models.ANIME_SOURCE_KURAMANIME) {
			t.Fatalf("%s has unexpected list capabilities", source)
		}
	}
	if _, err := GetCapabilities(context.Background(), models.AnimeQueryParams{Source: "invalid"}); !errors.Is(err, models.ErrAnimeSourceNotFound) {
		t.Fatalf("invalid source error = %v", err)
	}
}
