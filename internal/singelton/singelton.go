package singelton

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
)

type Singelton struct {
	Inited       bool
	DbAnimension models.DbAnimension
}

var singleton Singelton

func Initialize() error {
	var err error

	if singleton.Inited {
		return fmt.Errorf("singleton has been initialized")
	}
	dbAnimension := models.DbAnimension{}

	dbAnimension.AnimeList, err = animensionLoadAnimeIndex()
	if err != nil {
		logrus.Error(err)
		return err
	}
	dbAnimension.AnimeMap = animensionLoadAnimeListToMap(dbAnimension.AnimeList)
	dbAnimension.SeasonList, err = animensionLoadSeason()
	if err != nil {
		logrus.Error(err)
		return err
	}
	dbAnimension.SeasonMap = animensionLoadSeasonListToMap(dbAnimension.SeasonList)

	singleton = Singelton{
		Inited:       true,
		DbAnimension: dbAnimension,
	}

	return nil
}

func Get() Singelton {
	return singleton
}

func animensionLoadAnimeIndex() ([]models.AnimeDetail, error) {
	animeDbFile, err := os.Open("internal/local_db/db_animension_anime_index.json")
	if err != nil {
		logrus.Error(err)
		return []models.AnimeDetail{}, err
	}
	defer animeDbFile.Close()

	animeDBData, err := io.ReadAll(animeDbFile)
	if err != nil {
		logrus.Error(err)
		return []models.AnimeDetail{}, err
	}

	animeList := []models.AnimeDetail{}
	err = json.Unmarshal(animeDBData, &animeList)
	if err != nil {
		logrus.Error(err)
		return []models.AnimeDetail{}, err
	}

	return animeList, nil
}

func animensionLoadAnimeListToMap(animeList []models.AnimeDetail) map[string]models.AnimeDetail {
	animeMap := map[string]models.AnimeDetail{}

	for _, oneAnime := range animeList {
		animeMap[oneAnime.AnimensionAnimeID] = oneAnime
	}

	return animeMap
}

func animensionLoadSeason() ([]models.SeasonDetail, error) {
	seasonDbFile, err := os.Open("internal/local_db/db_animension_season_shorted_anime.json")
	if err != nil {
		logrus.Error(err)
		return []models.SeasonDetail{}, err
	}
	defer seasonDbFile.Close()

	seasonDBData, err := io.ReadAll(seasonDbFile)
	if err != nil {
		logrus.Error(err)
		return []models.SeasonDetail{}, err
	}

	seasonList := []models.SeasonDetail{}
	err = json.Unmarshal(seasonDBData, &seasonList)
	if err != nil {
		logrus.Error(err)
		return []models.SeasonDetail{}, err
	}

	return seasonList, nil
}

func animensionLoadSeasonListToMap(seasonList []models.SeasonDetail) map[string]models.SeasonDetail {
	seasonMap := map[string]models.SeasonDetail{}

	for _, oneSeason := range seasonList {
		seasonMap[""] = oneSeason
	}

	return seasonMap
}
