package animension_legacy_controller

import (
	"encoding/json"
	"os"
	"sort"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
)

type (
	SeasonDB []models.SeasonDetail
)

func syncSeason() {
	animeDB, err := fetchAnimeDB()
	if err != nil {
		logrus.Error(err)
		return
	}

	sort.Slice(animeDB, func(i, j int) bool {
		if animeDB[i].ReleaseYear == animeDB[j].ReleaseYear {
			if animeDB[i].ReleaseSeasonIndex == animeDB[j].ReleaseSeasonIndex {
				return animeDB[i].MalScore > animeDB[j].MalScore
			}
			return animeDB[i].ReleaseSeasonIndex > animeDB[j].ReleaseSeasonIndex
		}
		return animeDB[i].ReleaseYear > animeDB[j].ReleaseYear
	})

	seasonDB := SeasonDB{}

	for _, oneAnime := range animeDB {
		seasonIdx := -1

		for idx, oneSeason := range seasonDB {
			if oneAnime.ReleaseYear == oneSeason.Year && oneAnime.ReleaseSeasonIndex == oneSeason.SeasonIndex {
				seasonIdx = idx
			}
		}

		if seasonIdx == -1 {
			seasonDB = append(seasonDB, models.SeasonDetail{
				Year:        oneAnime.ReleaseYear,
				SeasonName:  oneAnime.ReleaseSeasonName,
				SeasonIndex: oneAnime.ReleaseSeasonIndex,
				AnimeData:   []models.AnimeSummary{},
			})

			for idx, oneSeason := range seasonDB {
				if oneAnime.ReleaseYear == oneSeason.Year && oneAnime.ReleaseSeasonIndex == oneSeason.SeasonIndex {
					seasonIdx = idx
				}
			}
		}

		tmpOneSeason := seasonDB[seasonIdx]
		tmpOneSeason.AnimeData = append(tmpOneSeason.AnimeData, models.AnimeSummary{
			AnimensionAnimeID:   oneAnime.AnimensionAnimeID,
			AnimensionEpisodeID: oneAnime.Episodes[0].AnimensionEpisodeID,
			CoverURL:            oneAnime.CoverURL,
			AvatarURL:           "/images/youtube.png",
			Title:               oneAnime.Title,
			MasterTitle:         oneAnime.MasterTitle,
			Status:              oneAnime.Status,
			LastEpisode:         oneAnime.Episodes[len(oneAnime.Episodes)-1].EpisodeNumber,
		})
		seasonDB[seasonIdx] = tmpOneSeason
	}

	resByte, err := json.MarshalIndent(seasonDB, "", "    ")
	if err != nil {
		logrus.Error(err)
		return
	}

	err = os.WriteFile("internal/local_db/db_animension_season_shorted_anime.json", resByte, 0644)
	if err != nil {
		logrus.Error(err)
		return
	}
}
