package local_db

import (
	"encoding/json"

	"github.com/umarkotak/animapu-api/internal/models"
)

var AnimeLinkToDetailMap = map[string]models.Anime{}

func Initialize() {
	json.Unmarshal([]byte(otakuRawDB), &AnimeLinkToDetailMap)
}
