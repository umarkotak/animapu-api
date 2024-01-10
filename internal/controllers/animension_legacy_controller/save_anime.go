package animension_legacy_controller

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
)

type (
	AnimeDBObj struct {
		AnimeDB AnimeDB
	}

	AnimeDB []models.AnimeDetail
)

func saveAnime(animeDetail models.AnimeDetail) error {
	animeDB, err := fetchAnimeDB()
	if err != nil {
		logrus.Error(err)
		return err
	}
	animeDBObj := AnimeDBObj{
		AnimeDB: animeDB,
	}

	existingAnime, existingIdx, err := animeDBObj.FindAnimeAndIndex(animeDetail.AnimensionAnimeID)
	// Create new anime
	if err != nil {
		animeDB = append(animeDB, animeDetail)
		err = writeAnimeDB(animeDB)
		if err != nil {
			logrus.Error(err)
			return err
		}
		return nil
	}

	// Update existing anime
	animeDetail.MasterTitleTag = existingAnime.MasterTitleTag
	animeDetail.SeasonIndex = existingAnime.SeasonIndex
	animeDetail.SeasonLabel = existingAnime.SeasonLabel
	animeDB[existingIdx] = animeDetail
	err = writeAnimeDB(animeDB)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func fetchAnimeDB() (AnimeDB, error) {
	animeDbFile, err := os.Open("internal/local_db/db_animension_anime_index.json")
	if err != nil {
		logrus.Error(err)
		return AnimeDB{}, err
	}
	defer animeDbFile.Close()

	animeDBData, err := io.ReadAll(animeDbFile)
	if err != nil {
		logrus.Error(err)
		return AnimeDB{}, err
	}

	animeDB := AnimeDB{}
	err = json.Unmarshal(animeDBData, &animeDB)
	if err != nil {
		logrus.Error(err)
		return AnimeDB{}, err
	}

	return animeDB, nil
}

func (adbo *AnimeDBObj) FindAnimeAndIndex(animeID string) (models.AnimeDetail, int, error) {
	for idx, oneAnimeDetail := range adbo.AnimeDB {
		if oneAnimeDetail.AnimensionAnimeID == animeID {
			return oneAnimeDetail, idx, nil
		}
	}
	return models.AnimeDetail{}, 0, fmt.Errorf("anime not found")
}

func writeAnimeDB(animeDB AnimeDB) error {
	resByte, err := json.MarshalIndent(animeDB, "", "    ")
	if err != nil {
		logrus.Error(err)
		return err
	}

	err = os.WriteFile("internal/local_db/db_animension_anime_index.json", resByte, 0644)
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}
