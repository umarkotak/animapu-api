package manga_scrapper_service

import (
	"context"
	"reflect"
	"testing"

	"github.com/umarkotak/animapu-api/internal/contract"
	"github.com/umarkotak/animapu-api/internal/models"
)

func TestExcludeNSFW(t *testing.T) {
	mangas := []contract.Manga{
		{Source: "one", SourceID: "marked"},
		{Source: "one", SourceID: "safe"},
		{Source: "one", SourceID: "unknown"},
		{Source: "two", SourceID: "marked"},
	}
	original := append([]contract.Manga(nil), mangas...)
	dbMangas := []models.Manga{
		{Source: "one", SourceID: "marked", Tags: []string{"for_kids:true", "nsfw:true"}},
		{Source: "one", SourceID: "safe", Tags: []string{"nsfw:false"}},
	}

	got := excludeNSFW(mangas, dbMangas)
	if !reflect.DeepEqual(got, original[1:]) {
		t.Fatalf("unexpected filtered mangas: %+v", got)
	}
	if !reflect.DeepEqual(mangas, original) {
		t.Fatal("filter modified the cached source results")
	}
	if got := excludeNSFW(mangas[:1], dbMangas); got == nil || len(got) != 0 {
		t.Fatalf("expected an empty list when all results are NSFW, got %+v", got)
	}
	if got := excludeNSFW(mangas, nil); !reflect.DeepEqual(got, mangas) {
		t.Fatalf("unknown mangas should remain visible: %+v", got)
	}
}

func TestFilterNSFWEmpty(t *testing.T) {
	got, err := FilterNSFW(context.Background(), "one", []contract.Manga{})
	if err != nil || got == nil || len(got) != 0 {
		t.Fatalf("expected an empty list without a database lookup, got %+v, %v", got, err)
	}
}
