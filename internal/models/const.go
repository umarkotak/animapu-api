package models

type MangaSource struct {
	ID       string `json:"id"`
	Language string `json:"language"`
	Title    string `json:"title"`
	WebLink  string `json:"web_link"`
	Active   bool   `json:"active"`
	Status   string `json:"status"`
}

const (
	SOURCE_MANGAHUB     = "mangahub"
	SOURCE_MANGAUPDATES = "mangaupdates"
	SOURCE_MANGABAT     = "mangabat"
	SOURCE_M_MANGABAT   = "m_mangabat"
	SOURCE_KLIKMANGA    = "klikmanga"
	SOURCE_WEBTOONSID   = "webtoonsid"
	SOURCE_FIZMANGA     = "fizmanga"
	SOURCE_MANGADEX     = "mangadex"
	SOURCE_MAIDMY       = "maidmy"
	SOURCE_MANGAREAD    = "mangaread"
	SOURCE_MANGATOWN    = "mangatown"
)

var (
	MangaSources = []MangaSource{
		{
			ID:       SOURCE_MANGABAT,
			Language: "en",
			Title:    "Manga Bat",
			WebLink:  "https://m.mangabat.com/m",
			Active:   true,
			Status:   "stable",
		},
		{
			ID:       SOURCE_FIZMANGA,
			Language: "en",
			Title:    "Fizmanga",
			WebLink:  "https://fizmanga.com/",
			Active:   false,
			Status:   "not-stable",
		},
		{
			ID:       SOURCE_MANGAUPDATES,
			Language: "en",
			Title:    "Manga Updates",
			WebLink:  "https://www.mangaupdates.com/",
			Active:   false,
			Status:   "not-stable",
		},
		{
			ID:       SOURCE_MANGAHUB,
			Language: "en",
			Title:    "Manga Hub",
			WebLink:  "https://www.mangahub.io/",
			Active:   false,
			Status:   "not-stable",
		},
		{
			ID:       SOURCE_KLIKMANGA,
			Language: "id",
			Title:    "Klik Manga",
			WebLink:  "https://klikmanga.id/",
			Active:   false,
			Status:   "stable",
		},
		{
			ID:       SOURCE_WEBTOONSID,
			Language: "id",
			Title:    "WebToon ID",
			WebLink:  "https://www.webtoons.com/id/",
			Active:   false,
			Status:   "stable",
		},
		{
			ID:       SOURCE_MANGADEX,
			Language: "mix",
			Title:    "Manga Dex",
			WebLink:  "https://mangadex.org/",
			Active:   false,
			Status:   "unavailable",
		},
		{
			ID:       SOURCE_MAIDMY,
			Language: "id",
			Title:    "Maid My",
			WebLink:  "https://www.maid.my.id/",
			Active:   false,
			Status:   "stable - from mangamee",
		},
		{
			ID:       SOURCE_MANGAREAD,
			Language: "en",
			Title:    "Manga Read",
			WebLink:  "https://www.mangaread.org/",
			Active:   false,
			Status:   "stable - from mangamee",
		},
		{
			ID:       SOURCE_MANGATOWN,
			Language: "en",
			Title:    "Manga Town",
			WebLink:  "https://apkpure.com/id/manga-town-manga-reader/com.mangatown.app1/",
			Active:   false,
			Status:   "stable - from mangamee",
		},
		{
			ID:       SOURCE_M_MANGABAT,
			Language: "en",
			Title:    "Mangabat",
			WebLink:  "https://m.mangabat.com/m",
			Active:   false,
			Status:   "stable - from mangamee",
		},
	}
)
