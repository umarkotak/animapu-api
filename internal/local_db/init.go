package local_db

import (
	_ "embed"
	"encoding/json"

	"github.com/umarkotak/animapu-api/internal/models"
)

var (
	//go:embed db_animension_anime_index.json
	db_animension_anime_index []byte
	//go:embed db_animension_season_shorted_anime.json
	db_animension_season_shorted_anime []byte
)

var (
	AnimeLinkToDetailMap = map[string]models.Anime{}

	AnimensionAnimeIndex    = []models.AnimeDetail{}
	AnimensionSeasonShorted = []models.SeasonDetail{}
)

func Initialize() {
	json.Unmarshal([]byte(otakuRawDB), &AnimeLinkToDetailMap)

	json.Unmarshal(db_animension_anime_index, &AnimensionAnimeIndex)
	json.Unmarshal(db_animension_season_shorted_anime, &AnimensionSeasonShorted)
}
