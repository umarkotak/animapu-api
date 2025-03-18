package manga_library_repository

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/datastore"
)

var (
	allColumns = strings.Join([]string{
		"ml.id",
		"ml.created_at",
		"ml.updated_at",
		"ml.deleted_at",
		"ml.user_id",
		"ml.manga_id",
	}, ", ")

	queryGetByID = fmt.Sprintf(`
		SELECT
			%s
		FROM manga_libraries ml
		WHERE
			ml.id = :id
			AND ml.deleted_at IS NULL
	`, allColumns)

	queryGetByMangaIDAndUserID = fmt.Sprintf(`
		SELECT
			%s
		FROM manga_libraries ml
		WHERE
			ml.user_id = :user_id
			AND ml.manga_id = :manga_id
			AND ml.deleted_at IS NULL
	`, allColumns)

	queryGetByMangaIDsAndUserID = fmt.Sprintf(`
		SELECT
			%s
		FROM manga_libraries ml
		WHERE
			ml.user_id = :user_id
			AND ml.manga_id = ANY(:manga_ids)
			AND ml.deleted_at IS NULL
	`, allColumns)

	queryGetByUserID = fmt.Sprintf(`
		SELECT
			%s
		FROM manga_libraries ml
		WHERE
			ml.user_id = :user_id
			AND ml.deleted_at IS NULL
	`, allColumns)

	queryInsert = `
		INSERT INTO manga_libraries (
			user_id,
			manga_id
		) VALUES (
			:user_id,
			:manga_id
		)
		ON CONFLICT (user_id, manga_id)
		DO UPDATE SET
			updated_at = NOW()
		RETURNING id
	`

	queryDelete = `
		DELETE FROM manga_libraries
		WHERE
			user_id = :user_id
			AND manga_id = :manga_id
	`

	queryGetByUserAndSourceDetail = `
		SELECT
			m.source_id
		FROM manga_libraries ml
		INNER JOIN mangas m ON m.id = ml.manga_id
		WHERE
			ml.user_id = :user_id
			AND m.source = ANY(:sources)
			AND m.source_id = ANY(:source_ids)
	`
)

var (
	stmtGetByID                  *sqlx.NamedStmt
	stmtGetByMangaIDAndUserID    *sqlx.NamedStmt
	stmtGetByMangaIDsAndUserID   *sqlx.NamedStmt
	stmtGetByUserID              *sqlx.NamedStmt
	stmtInsert                   *sqlx.NamedStmt
	stmtDelete                   *sqlx.NamedStmt
	stmtGetByUserAndSourceDetail *sqlx.NamedStmt
)

func Initialize() {
	var err error

	stmtGetByID, err = datastore.Get().Db.PrepareNamed(queryGetByID)
	if err != nil {
		logrus.Fatal(err)
	}

	stmtGetByMangaIDAndUserID, err = datastore.Get().Db.PrepareNamed(queryGetByMangaIDAndUserID)
	if err != nil {
		logrus.Fatal(err)
	}

	stmtGetByMangaIDsAndUserID, err = datastore.Get().Db.PrepareNamed(queryGetByMangaIDsAndUserID)
	if err != nil {
		logrus.Fatal(err)
	}

	stmtGetByUserID, err = datastore.Get().Db.PrepareNamed(queryGetByUserID)
	if err != nil {
		logrus.Fatal(err)
	}

	stmtInsert, err = datastore.Get().Db.PrepareNamed(queryInsert)
	if err != nil {
		logrus.Fatal(err)
	}

	stmtDelete, err = datastore.Get().Db.PrepareNamed(queryDelete)
	if err != nil {
		logrus.Fatal(err)
	}

	stmtGetByUserAndSourceDetail, err = datastore.Get().Db.PrepareNamed(queryGetByUserAndSourceDetail)
	if err != nil {
		logrus.Fatal(err)
	}
}
