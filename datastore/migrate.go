package datastore

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/config"
)

func MigrateUp() error {
	m, err := migrate.New("file://db/migrations", config.Get().DbUrl)
	if err != nil {
		logrus.Error(err)
		return err
	}

	err = m.Up()
	if err != nil && err.Error() != "no change" {
		logrus.Error(err)
		return err
	}

	return nil
}
