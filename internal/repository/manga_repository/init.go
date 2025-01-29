package manga_repository

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/datastore"
)

var (
	allColumns = strings.Join([]string{
		"m.id",
		"m.created_at",
		"m.updated_at",
		"m.deleted_at",
		"m.source",
		"m.source_id",
		"m.title",
		"m.cover_urls",
		"m.latest_chapter",
	}, ", ")

	queryGetByID = fmt.Sprintf(`
		SELECT
			%s
		FROM mangas m
		WHERE
			m.id = :id
			AND m.deleted_at IS NULL
	`, allColumns)

	queryGetBySourceAndSourceID = fmt.Sprintf(`
		SELECT
			%s
		FROM mangas m
		WHERE
			m.source = :source
			AND m.source_id = :source_id
			AND m.deleted_at IS NULL
	`, allColumns)

	queryInsert = `
		INSERT INTO mangas (
			source,
			source_id,
			title,
			cover_urls,
			latest_chapter
		) VALUES (
			:source,
			:source_id,
			:title,
			:cover_urls,
			:latest_chapter
		)
		ON CONFLICT (source, source_id)
		DO UPDATE SET
			latest_chapter = :latest_chapter
		RETURNING id
	`

	queryUpdate = `
		UPDATE mangas
		SET
			source = :source,
			source_id = :source_id,
			title = :title,
			cover_urls = :cover_urls,
			latest_chapter = :latest_chapter
		WHERE
			id = :id
	`

	queryUpdateBySourceAndSourceID = `
		UPDATE mangas
		SET
			title = :title,
			cover_urls = :cover_urls,
			latest_chapter = :latest_chapter
		WHERE
			source = :source
			AND source_id = :source_id
	`
)

var (
	stmtGetByID                   *sqlx.NamedStmt
	stmtGetBySourceAndSourceID    *sqlx.NamedStmt
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
