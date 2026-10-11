package weebcentral

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/umarkotak/animapu-api/config"
	"github.com/umarkotak/animapu-api/internal/models"
)

func TestGetChapter(t *testing.T) {
	body, err := os.ReadFile("testdata/chapter-images.html")
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("ANIMAPU_API_HOST", "https://api.example.com")
	if err := config.Initialize(); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		body   string
		status int
		count  int
		number float64
	}{
		{name: "sample", body: string(body), status: http.StatusOK, count: 18, number: 104},
		{name: "arbitrary filenames", body: `<img src="https://example.com/cover.png"><section id="chapter-images"><img><img src="https://example.com/page-b.webp"><img src="https://example.com/page-a.jpg"></section>`, status: http.StatusOK, count: 2},
		{name: "empty response", body: `<section id="chapter-images"></section>`, status: http.StatusOK},
		{name: "request failure", status: http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/chapters/01M43Q7CFX4XXN7WBVZH1MTEFS/images" || r.URL.RawQuery != "reading_style=long_strip" {
					t.Errorf("unexpected request: %s", r.URL)
				}
				for name, want := range map[string]string{
					"Upgrade-Insecure-Requests": "1",
					"User-Agent":                "Mozilla/5.0 (Linux; Android 16; Pixel 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/154.0.0.0 Mobile Safari/537.36",
					"sec-ch-ua":                 `"Chromium";v="154", "Google Chrome";v="154", "Not A(Brand";v="99"`,
					"sec-ch-ua-mobile":          "?1",
					"sec-ch-ua-platform":        `"Android"`,
				} {
					if got := r.Header.Get(name); got != want {
						t.Errorf("header %s = %q, want %q", name, got, want)
					}
				}
				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer server.Close()
			sc := New()
			sc.Host = server.URL
			params := models.QueryParams{ChapterID: "01M43Q7CFX4XXN7WBVZH1MTEFS", SourceID: "series-id"}
			chapter, err := sc.GetChapter(context.Background(), params)
			if tt.count == 0 {
				if err == nil {
					t.Fatal("expected an error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if len(chapter.ChapterImages) != tt.count || chapter.Number != tt.number {
				t.Fatalf("got %d images, chapter %v; want %d images, chapter %v", len(chapter.ChapterImages), chapter.Number, tt.count, tt.number)
			}
			if chapter.ID != params.ChapterID || chapter.SourceID != params.SourceID || chapter.Source != sc.Source || chapter.SourceLink != server.URL+"/chapters/"+params.ChapterID {
				t.Fatalf("unexpected chapter metadata: %#v", chapter)
			}
			for i, img := range chapter.ChapterImages {
				src := fmt.Sprintf("https://scans-hot.planeptune.us/manga/Centuria/0104-%03d.png", i+1)
				if tt.name == "arbitrary filenames" {
					src = []string{"https://example.com/page-b.webp", "https://example.com/page-a.jpg"}[i]
				}
				want := "https://api.example.com/mangas/weeb_central/image_proxy/" + src
				if img.Index != int64(i+1) || len(img.ImageUrls) != 1 || img.ImageUrls[0] != want {
					t.Errorf("image %d = %#v, want index %d and URL %s", i, img, i+1, want)
				}
			}
		})
	}
}
