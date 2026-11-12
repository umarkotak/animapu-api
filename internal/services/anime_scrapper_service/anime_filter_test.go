package anime_scrapper_service

import (
	"context"
	"reflect"
	"testing"

	"github.com/umarkotak/animapu-api/internal/contract"
	"github.com/umarkotak/animapu-api/internal/models"
)

func TestExcludeNSFW(t *testing.T) {
	animes := []contract.Anime{
		{Source: "one", ID: "marked"},
		{Source: "one", ID: "safe"},
		{Source: "one", ID: "unknown"},
		{Source: "two", ID: "marked"},
	}
	original := append([]contract.Anime(nil), animes...)
	dbAnimes := []models.Anime{
		{Source: "one", SourceID: "marked", Tags: []string{"for_kids:true", "nsfw:true"}},
		{Source: "one", SourceID: "safe", Tags: []string{"nsfw:false"}},
	}

	got := excludeNSFW(animes, dbAnimes)
	if !reflect.DeepEqual(got, original[1:]) {
		t.Fatalf("unexpected filtered animes: %+v", got)
	}
	if !reflect.DeepEqual(animes, original) {
		t.Fatal("filter modified the cached source results")
	}
	if got := excludeNSFW(animes[:1], dbAnimes); got == nil || len(got) != 0 {
		t.Fatalf("expected an empty list when all results are NSFW, got %+v", got)
	}
	if got := excludeNSFW(animes, nil); !reflect.DeepEqual(got, animes) {
		t.Fatalf("unknown animes should remain visible: %+v", got)
	}
}

func TestFilterNSFWEmpty(t *testing.T) {
	got, err := FilterNSFW(context.Background(), "one", []contract.Anime{})
	if err != nil || got == nil || len(got) != 0 {
		t.Fatalf("expected an empty list without a database lookup, got %+v, %v", got, err)
	}
}
