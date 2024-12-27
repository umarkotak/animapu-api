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
	SOURCE_MANGABAT   = "mangabat"
	SOURCE_ASURA_NACM = "asura_nacm"
	SOURCE_KOMIKINDO  = "komikindo"
	SOURCE_KOMIKU     = "komiku"
	SOURCE_MANGASEE   = "mangasee"

	// powered by mangamee
	SOURCE_M_MANGABAT  = "m_mangabat"
	SOURCE_MAIDMY      = "maidmy"
	SOURCE_MANGAREAD   = "mangaread"
	SOURCE_MANGATOWN   = "mangatown"
	SOURCE_ASURA_COMIC = "m_asura"
	SOURCE_MANGANATO   = "manganato"
	SOURCE_MANGANELO   = "manganelo"
	SOURCE_M_MANGASEE  = "m_mangasee"

	SOURCE_MANGAHUB   = "mangahub"
	SOURCE_KLIKMANGA  = "klikmanga"
	SOURCE_WEBTOONSID = "webtoonsid"
	SOURCE_MANGADEX   = "mangadex"

	ANIME_SOURCE_OTAKUDESU        = "otakudesu"
	ANIME_SOURCE_ANIMENSION       = "animension"
	ANIME_SOURCE_ANIMENSION_LOCAL = "animension_local"
)

var (
	MangaSources = []MangaSource{
		{
			ID:       SOURCE_KOMIKINDO,
			Language: "id",
			Title:    "KomikIndo",
			WebLink:  "https://komikindo.pw",
			Active:   true,
			Status:   "stable",
		},
		{
			ID:       SOURCE_KOMIKU,
			Language: "id",
			Title:    "Komiku",
			WebLink:  "https://komiku.id",
			Active:   true,
			Status:   "stable",
		},
		{
			ID:       SOURCE_MANGABAT,
			Language: "en",
			Title:    "Manga Bat",
			WebLink:  "https://m.mangabat.com/m",
			Active:   true,
			Status:   "stable",
		},
		{
			ID:       SOURCE_ASURA_NACM,
			Language: "en",
			Title:    "Asura",
			WebLink:  "https://asura.nacm.xyz",
			Active:   true,
			Status:   "stable",
		},
		{
			ID:       SOURCE_MANGASEE,
			Language: "en",
			Title:    "Mangasee",
			WebLink:  "https://www.mangasee123.com",
			Active:   true,
			Status:   "stable",
		},

		// mangamee
		{
			ID:       SOURCE_M_MANGABAT,
			Language: "en",
			Title:    "Mangabat",
			WebLink:  "https://m.mangabat.com/m",
			Active:   true,
			Status:   "stable - from mangamee",
		},
		{
			ID:       SOURCE_MAIDMY,
			Language: "id",
			Title:    "Maid My",
			WebLink:  "https://www.maid.my.id/",
			Active:   true,
			Status:   "stable - from mangamee",
		},
		{
			ID:       SOURCE_MANGAREAD,
			Language: "en",
			Title:    "Manga Read",
			WebLink:  "https://www.mangaread.org/",
			Active:   true,
			Status:   "stable - from mangamee",
		},
		{
			ID:       SOURCE_MANGATOWN,
			Language: "en",
			Title:    "Manga Town",
			WebLink:  "https://apkpure.com/id/manga-town-manga-reader/com.mangatown.app1/",
			Active:   true,
			Status:   "stable - from mangamee",
		},
		{
			ID:       SOURCE_ASURA_COMIC,
			Language: "en",
			Title:    "Asura",
			WebLink:  "https://asura.nacm.xyz",
			Active:   true,
			Status:   "stable - from mangamee",
		},
		{
			ID:       SOURCE_MANGANATO,
			Language: "en",
			Title:    "Manga Nato",
			WebLink:  "https://manganato.com/index.php",
			Active:   true,
			Status:   "stable - from mangamee",
		},
		{
			ID:       SOURCE_MANGANELO,
			Language: "en",
			Title:    "Manga Nelo",
			WebLink:  "https://ww6.manganelo.tv/home",
			Active:   true,
			Status:   "stable - from mangamee",
		},
		{
			ID:       SOURCE_M_MANGASEE,
			Language: "en",
			Title:    "Mangasee",
			WebLink:  "https://www.mangasee123.com",
			Active:   true,
			Status:   "stable - from mangamee",
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
	}
)

type Season struct {
	Index int64  `json:"index"`
	Name  string `json:"name"`
}

var (
	MONTH_TO_SEASON_MAP = map[string]Season{
		// winter
		"jan": {1, "winter"},
		"feb": {1, "winter"},
		"mar": {1, "winter"},
		"":    {1, "winter"},
		"10":  {1, "winter"},
		// spring
		"apr": {2, "spring"},
		"may": {2, "spring"},
		"jun": {2, "spring"},
		// summer
		"jul": {3, "summer"},
		"aug": {3, "summer"},
		"agu": {3, "summer"},
		"sep": {3, "summer"},
		// fall
		"oct": {4, "fall"},
		"okt": {4, "fall"},
		"nov": {4, "fall"},
		"des": {4, "fall"},
		"dec": {4, "fall"},
	}

	SEASON_TO_SEASON_INDEX = map[string]int64{
		"winter": 1,
		"spring": 2,
		"summer": 3,
		"fall":   4,
	}
)
