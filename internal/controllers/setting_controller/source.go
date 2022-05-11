package setting_controller

import (
	"github.com/gin-gonic/gin"
	"github.com/umarkotak/animapu-api/internal/utils/render"
)

type MangaSource struct {
	ID       string `json:"id"`
	Language string `json:"language"`
	Title    string `json:"title"`
	WebLink  string `json:"web_link"`
	Active   bool   `json:"active"`
}

func GetAvailableSource(c *gin.Context) {
	render.Response(
		c.Request.Context(),
		c,
		[]MangaSource{
			{
				ID:       "mangabat",
				Language: "en",
				Title:    "Manga Bat",
				WebLink:  "https://m.mangabat.com/m",
				Active:   true,
			},
			{
				ID:       "mangaupdates",
				Language: "en",
				Title:    "Manga Updates",
				WebLink:  "https://www.mangaupdates.com/",
				Active:   true,
			},
			{
				ID:       "klikmanga",
				Language: "id",
				Title:    "Klik Manga",
				WebLink:  "https://klikmanga.id/",
				Active:   true,
			},
			{
				ID:       "fizmanga",
				Language: "en",
				Title:    "Fizmanga",
				WebLink:  "https://fizmanga.com/",
				Active:   true,
			},
			{
				ID:       "webtoonsid",
				Language: "id",
				Title:    "WebToon ID",
				WebLink:  "https://www.webtoons.com/id/",
				Active:   true,
			},
			{
				ID:       "mangadex",
				Language: "mix",
				Title:    "Manga Dex",
				WebLink:  "https://mangadex.org/",
				Active:   false,
			},
			{
				ID:       "maidmy",
				Language: "id",
				Title:    "Maid My",
				WebLink:  "https://www.maid.my.id/",
				Active:   false,
			},
			{
				ID:       "mangaread",
				Language: "en",
				Title:    "Manga Read",
				WebLink:  "https://www.mangaread.org/",
				Active:   false,
			},
		},
		nil,
		200,
	)
}
