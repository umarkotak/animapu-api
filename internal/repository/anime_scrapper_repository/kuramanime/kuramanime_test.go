package kuramanime

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/umarkotak/animapu-api/internal/models"
)

func TestGetLatest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/quick/ongoing" {
			return
		}
		if r.URL.Query().Get("page") != "2" {
			t.Fatalf("page = %q, want 2", r.URL.Query().Get("page"))
		}
		_, _ = w.Write([]byte(`{"animes":{"data":[{"id":5044,"title":"Saijo no Osewa","slug":"saijo-no-osewa","synopsis":"Synopsis","latest_episode":9,"score":7.3,"image_portrait_url":"https://example.com/cover.jpg","full_alt_titles":"Title A, Title B"}]}}`))
	}))
	defer server.Close()

	animes, err := (&Kuramanime{KuramanimeAggregator: server.URL}).GetLatest(context.Background(), models.AnimeQueryParams{Page: 2})
	if err != nil {
		t.Fatal(err)
	}
	if len(animes) != 1 {
		t.Fatalf("len(animes) = %d, want 1", len(animes))
	}

	anime := animes[0]
	if anime.ID != "5044" || anime.LatestEpisode != 9 || anime.OriginalLink != server.URL+"/anime/5044/saijo-no-osewa" {
		t.Fatalf("unexpected anime: %#v", anime)
	}
}

func TestGetDetail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/anime/5044" {
			return
		}
		_, _ = w.Write([]byte(`<!doctype html><html><body>
			<div class="anime__details__pic set-bg" data-setbg="https://example.com/cover.jpg"></div>
			<div class="anime__details__title"><h3>Test Anime</h3><span>Alt One, Alt Two</span></div>
			<p id="synopsisField">Test synopsis</p>
			<div class="anime__details__widget"><li><div><div class="col-3"><span>Tayang:</span></div><div class="col-9">05 Jul 2026 s/d ?</div></div></li><li><div><div class="col-3"><span>Genre:</span></div><div class="col-9"><a>Comedy,</a><a>Romance</a></div></div></li><li><div><div class="col-3"><span>Skor:</span></div><div class="col-9">7.30 / 10.00</div></div></li></div>
			<a id="episodeLists" data-content="&lt;a href='/anime/5044/test-anime/episode/1'&gt;Ep 1&lt;/a&gt;&lt;a href='/anime/5044/test-anime/episode/2'&gt;Ep 2&lt;/a&gt;"></a>
		</body></html>`))
	}))
	defer server.Close()

	anime, err := (&Kuramanime{KuramanimeAggregator: server.URL}).GetDetail(context.Background(), models.AnimeQueryParams{SourceID: "5044"})
	if err != nil {
		t.Fatal(err)
	}
	if anime.Title != "Test Anime" || anime.ReleaseSeason != "summer" || len(anime.Episodes) != 2 {
		t.Fatalf("unexpected anime: %#v", anime)
	}
	if anime.Episodes[1].Number != 2 || anime.Episodes[1].CoverUrl != "https://example.com/cover.jpg" || anime.Episodes[1].CoverUrls[0] != "https://example.com/cover.jpg" {
		t.Fatalf("unexpected episode: %#v", anime.Episodes[1])
	}
}

func TestMP4URLs(t *testing.T) {
	urls := mp4URLs(`{"first":"https:\/\/cdn.example.com\/first.mp4?token=one","second":"https://cdn.example.com/second.MP4"}`)
	if len(urls) != 2 || urls[0] != "https://cdn.example.com/first.mp4?token=one" || urls[1] != "https://cdn.example.com/second.MP4" {
		t.Fatalf("mp4URLs() = %#v", urls)
	}
}
