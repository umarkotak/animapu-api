package anime_repository

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/datastore"
)

var (
	allColumns = strings.Join([]string{
		"a.id",
		"a.created_at",
		"a.updated_at",
		"a.deleted_at",
		"a.source",
		"a.source_id",
		"a.title",
		"a.cover_urls",
		"a.latest_episode",
	}, ", ")

	queryGetByID = fmt.Sprintf(`
		SELECT
			%s
		FROM animes a
		WHERE
			a.id = :id
			AND a.deleted_at IS NULL
	`, allColumns)

	queryGetBySourceAndSourceID = fmt.Sprintf(`
		SELECT
			%s
		FROM animes a
		WHERE
			a.source = :source
			AND a.source_id = :source_id
			AND a.deleted_at IS NULL
	`, allColumns)

	queryGetBySourceAndSourceIDs = fmt.Sprintf(`
		SELECT
			%s
		FROM animes a
		WHERE
			a.source = :source
			AND a.source_id = ANY(:source_ids)
			AND a.deleted_at IS NULL
	`, allColumns)

	queryInsert = `
		INSERT INTO animes (
			source,
			source_id,
			title,
			cover_urls,
			latest_episode
		) VALUES (
			:source,
			:source_id,
			:title,
			:cover_urls,
			:latest_episode
		)
		ON CONFLICT (source, source_id)
		DO UPDATE SET
			latest_episode = :latest_episode
		RETURNING id
	`

	queryUpdate = `
		UPDATE animes
		SET
			source = :source,
			source_id = :source_id,
			title = :title,
			cover_urls = :cover_urls,
			latest_episode = :latest_episode
		WHERE
			id = :id
	`

	queryUpdateBySourceAndSourceID = `
		UPDATE animes
		SET
			title = :title,
			cover_urls = :cover_urls,
			latest_episode = :latest_episode
		WHERE
			source = :source
			AND source_id = :source_id
	`
)

var (
	stmtGetByID                   *sqlx.NamedStmt
	stmtGetBySourceAndSourceID    *sqlx.NamedStmt
	stmtGetBySourceAndSourceIDs   *sqlx.NamedStmt
	stmtInsert                    *sqlx.NamedStmt
	stmtUpdate                    *sqlx.NamedStmt
	stmtUpdateBySourceAndSourceID *sqlx.NamedStmt
)

func Initialize() {
	var err error

	stmtGetByID, err = datastore.Get().Db.PrepareNamed(queryGetByID)
	if err != nil {
		logrus.Fatal(err)
	}

	stmtGetBySourceAndSourceID, err = datastore.Get().Db.PrepareNamed(queryGetBySourceAndSourceID)
	if err != nil {
		logrus.Fatal(err)
	}

	stmtGetBySourceAndSourceIDs, err = datastore.Get().Db.PrepareNamed(queryGetBySourceAndSourceIDs)
	if err != nil {
		logrus.Fatal(err)
	}

	stmtInsert, err = datastore.Get().Db.PrepareNamed(queryInsert)
	if err != nil {
		logrus.Fatal(err)
	}

	stmtUpdate, err = datastore.Get().Db.PrepareNamed(queryUpdate)
	if err != nil {
		logrus.Fatal(err)
	}

	stmtUpdateBySourceAndSourceID, err = datastore.Get().Db.PrepareNamed(queryUpdateBySourceAndSourceID)
	if err != nil {
		logrus.Fatal(err)
	}
}
