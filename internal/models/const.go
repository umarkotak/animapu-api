package models

type (
	MangaSource struct {
		ID       string `json:"id"`
		Language string `json:"language"`
		Title    string `json:"title"`
		WebLink  string `json:"web_link"`
		Active   bool   `json:"active"`
		Status   string `json:"status"`
	}

	AnimeSource struct {
		ID       string `json:"id"`
		Language string `json:"language"`
		Title    string `json:"title"`
		WebLink  string `json:"web_link"`
		Active   bool   `json:"active"`
		Status   string `json:"status"`
	}
)

const (
	// indonesia
	SOURCE_KOMIKINDO = "komikindo"
	SOURCE_KOMIKU    = "komiku"
	SOURCE_KOMIKCAST = "komikcast"

	// english
	SOURCE_MANGABAT     = "mangabat"
	SOURCE_ASURA_NACM   = "asura_nacm"
	SOURCE_MANGASEE     = "mangasee"
	SOURCE_WEEB_CENTRAL = "weeb_central"
	SOURCE_MANGADEX     = "mangadex"
	SOURCE_MANGAKAKALOT = "manga_kakalot" // TODO
	SOURCE_MANGANATO    = "manga_nato"    // TODO

	// indonesia
	ANIME_SOURCE_OTAKUDESU     = "otakudesu"
	ANIME_SOURCE_KURAMANIME    = "kuramanime"
	ANIME_SOURCE_ANIMEINDO     = "animeindo"     // TODO: https://anime-indo.lol
	ANIME_SOURCE_SAMEHADAKU_AC = "samehadaku_ac" // TODO: https://samehadaku.ac

	// english
	ANIME_SOURCE_GOGO_ANIME     = "gogo_anime"
	ANIME_SOURCE_GOGO_ANIME_NEW = "gogo_anime_new"
	ANIME_SOURCE_GOGO_ANIME_VC  = "gogo_anime_vc" // TODO: https://gogoanime.org.vc
	ANIME_SOURCE_GOGO_ANIME_CZ  = "gogo_anime_cz" // TODO: https://gogoanime.co.cz
	// TODO: https://hianime.to, https://theindex.moe/collection/self-hosted-streaming-sites
)

var (
	MangaSources = []MangaSource{
		{
			ID:       SOURCE_KOMIKINDO,
			Language: "id",
			Title:    "Komik Indo",
			WebLink:  "https://komikindo.ch",
			Active:   true,
			Status:   "stable",
		},
		{
			ID:       SOURCE_KOMIKU,
			Language: "id",
			Title:    "Komiku",
			WebLink:  "https://komiku.org",
			Active:   true,
			Status:   "stable",
		},
		{
			ID:       SOURCE_WEEB_CENTRAL,
			Language: "en",
			Title:    "Weeb Central",
			WebLink:  "https://weebcentral.com",
			Active:   true,
			Status:   "stable",
		},
	}

	AnimeSources = []AnimeSource{
		{
			ID:       ANIME_SOURCE_OTAKUDESU,
			Language: "id",
			Title:    "Otakudesu",
			WebLink:  "https://otakudesu.cloud",
			Active:   true,
			Status:   "stable",
		},
	}

	AdminEmails = []string{
		"umarkotak@gmail.com",
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
