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

	// indonesia
	ANIME_SOURCE_OTAKUDESU     = "otakudesu"
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
			ID:       SOURCE_KOMIKCAST,
			Language: "id",
			Title:    "Komik Cast",
			WebLink:  "https://komikcast.bz",
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
		{
			ID:       SOURCE_WEEB_CENTRAL,
			Language: "en",
			Title:    "Weeb Central",
			WebLink:  "https://weebcentral.com",
			Active:   true,
			Status:   "stable",
		},
		{
			ID:       SOURCE_MANGADEX,
			Language: "en",
			Title:    "Mangadex",
			WebLink:  "https://mangadex.org",
			Active:   false,
			Status:   "wip",
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
		{
			ID:       ANIME_SOURCE_GOGO_ANIME,
			Language: "en",
			Title:    "Gogo Anime Old",
			WebLink:  "https://ww10.gogoanimes.org",
			Active:   false,
			Status:   "stable",
		},
		{
			ID:       ANIME_SOURCE_GOGO_ANIME_NEW,
			Language: "en",
			Title:    "Gogo Anime New",
			WebLink:  "https://gogoanime.by", // https://gogoanime.by,
			Active:   false,
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
