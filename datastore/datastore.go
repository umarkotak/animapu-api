package datastore

import (
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/config"
)

type DataStore struct {
	Db      *sqlx.DB // required
	GoCache *cache.Cache
}

var dataStore DataStore

func Initialize() error {
	db, err := sqlx.Connect("postgres", config.Get().DbUrl)
	if err != nil {
		logrus.Error(err)
		return err
	}

	goCache := cache.New(5*time.Minute, 10*time.Minute)

	dataStore = DataStore{
		Db:      db,
		GoCache: goCache,
	}

	return nil
}

func Get() DataStore {
	return dataStore
}
