package kuramanime

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/umarkotak/animapu-api/internal/contract"
	"github.com/umarkotak/animapu-api/internal/models"
)

func TestGetCompleted(t *testing.T) {
	for _, tc := range []struct {
		name      string
		page      int64
		wantPage  string
		status    int
		body      string
		wantError bool
		wantCount int
	}{
		{"page two", 2, "2", 200, `{"animes":{"data":[{"id":5044,"title":"Finished Anime","slug":"finished-anime","synopsis":"Synopsis","latest_episode":12,"score":7.3,"image_portrait_url":"https://example.com/cover.jpg","full_alt_titles":"Title A, Title B"}]}}`, false, 1},
		{"default page", 0, "1", 200, `{"animes":{"data":[]}}`, false, 0},
		{"upstream error", 2, "2", 503, `unavailable`, true, 0},
		{"invalid json", 2, "2", 200, `invalid`, true, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// The aggregator host lookup visits the root before the list request.
				if r.URL.Path == "/" {
					return
				}
				if r.URL.Path != "/quick/finished" || r.URL.Query().Get("page") != tc.wantPage || r.URL.Query().Get("order_by") != "latest" || r.URL.Query().Get("need_json") != "true" {
					t.Errorf("unexpected request: %s", r.URL)
				}
				w.WriteHeader(tc.status)
				_, _ = fmt.Fprint(w, tc.body)
			}))
			defer server.Close()
			result, err := (&Kuramanime{KuramanimeAggregator: server.URL}).GetCompleted(context.Background(), models.AnimeQueryParams{Page: tc.page})
			if (err != nil) != tc.wantError {
				t.Fatalf("error = %v, want error = %v", err, tc.wantError)
			}
			if len(result) != tc.wantCount {
				t.Fatalf("results = %#v", result)
			}
			if tc.wantCount == 1 {
				want := contract.Anime{
					ID: "5044", Source: models.ANIME_SOURCE_KURAMANIME, Title: "Finished Anime",
					Description: "Synopsis", LatestEpisode: 12, Score: 7.3,
					AltTitles: []string{"Title A", "Title B"}, CoverUrls: []string{"https://example.com/cover.jpg"},
					OriginalLink: server.URL + "/anime/5044/finished-anime",
				}
				if !reflect.DeepEqual(result[0], want) {
					t.Fatalf("anime = %#v, want %#v", result[0], want)
				}
			}
		})
	}
}
