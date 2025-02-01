package manga_history_repository

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/datastore"
)

var (
	allColumns = strings.Join([]string{
		"mh.id",
		"mh.created_at",
		"mh.updated_at",
		"mh.deleted_at",
		"mh.user_id",
		"mh.manga_id",
		"mh.chapter_number",
		"mh.source_chapter_id",
		"mh.frontend_path",
	}, ", ")

	queryGetByID = fmt.Sprintf(`
		SELECT
			%s
		FROM manga_histories mh
		WHERE
			mh.id = :id
			AND mh.deleted_at IS NULL
	`, allColumns)

	queryGetByMangaIDAndUserID = fmt.Sprintf(`
		SELECT
			%s
		FROM manga_histories mh
		WHERE
			mh.user_id = :user_id
			AND mh.manga_id = :manga_id
			AND mh.deleted_at IS NULL
	`, allColumns)

	queryGetByUserID = fmt.Sprintf(`
		SELECT
			%s
		FROM manga_histories mh
		WHERE
			mh.user_id = :user_id
			AND mh.deleted_at IS NULL
	`, allColumns)

	queryGetByUserIDDetailed = fmt.Sprintf(`
		SELECT
			%s,
			m.source AS manga_source,
			m.source_id AS manga_source_id,
			m.title AS manga_title,
			m.cover_urls AS manga_cover_urls,
			m.latest_chapter AS manga_latest_chapter
		FROM manga_histories mh
		INNER JOIN mangas m ON m.id = mh.manga_id
		WHERE
			mh.user_id = :user_id
			AND mh.deleted_at IS NULL
		ORDER BY mh.updated_at DESC
		LIMIT :limit OFFSET :offset
	`, allColumns)

	queryGetByUserAndSourceDetail = fmt.Sprintf(`
		SELECT
			%s,
			m.source AS manga_source,
			m.source_id AS manga_source_id,
			m.title AS manga_title,
			m.cover_urls AS manga_cover_urls,
			m.latest_chapter AS manga_latest_chapter
		FROM manga_histories mh
		INNER JOIN mangas m ON m.id = mh.manga_id
		WHERE
			mh.user_id = :user_id
			AND m.source = ANY(:sources)
			AND m.source_id = ANY(:source_ids)
	`, allColumns)

	queryInsert = `
		INSERT INTO manga_histories (
			user_id,
			manga_id,
			chapter_number,
			source_chapter_id,
			frontend_path
		) VALUES (
			:user_id,
			:manga_id,
			:chapter_number,
			:source_chapter_id,
			:frontend_path
		)
		ON CONFLICT (user_id, manga_id)
		DO UPDATE SET
			chapter_number = :chapter_number,
			source_chapter_id = :source_chapter_id,
			frontend_path = :frontend_path,
			updated_at = NOW()
		RETURNING id
	`

	queryUpdate = `
		UPDATE manga_histories
		SET
			user_id = :user_id,
			manga_id = :manga_id,
			chapter_number = :chapter_number,
			source_chapter_id = :source_chapter_id,
			frontend_path = :frontend_path
		WHERE
			id = :id
	`

	queryUpdateByMangaIDAndUserID = `
		UPDATE manga_histories
		SET
			chapter_number = :chapter_number,
			source_chapter_id = :source_chapter_id,
			frontend_path = :frontend_path
		WHERE
			user_id = :user_id
			AND manga_id = :manga_id
	`

	queryGetRecentHistories = fmt.Sprintf(`
		SELECT
			%s,
			m.source AS manga_source,
			m.source_id AS manga_source_id,
			m.title AS manga_title,
			m.cover_urls AS manga_cover_urls,
			m.latest_chapter AS manga_latest_chapter
		FROM manga_histories mh
		INNER JOIN mangas m ON m.id = mh.manga_id
		ORDER BY mh.updated_at DESC
		LIMIT :limit OFFSET :offset
	`, allColumns)
)

var (
	stmtGetByID                  *sqlx.NamedStmt
	stmtGetByMangaIDAndUserID    *sqlx.NamedStmt
	stmtGetByUserID              *sqlx.NamedStmt
	stmtGetByUserIDDetailed      *sqlx.NamedStmt
	stmtGetByUserAndSourceDetail *sqlx.NamedStmt
	stmtInsert                   *sqlx.NamedStmt
	stmtUpdate                   *sqlx.NamedStmt
	stmtUpdateByMangaIDAndUserID *sqlx.NamedStmt
	stmtGetRecentHistories       *sqlx.NamedStmt
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

	stmtGetByUserID, err = datastore.Get().Db.PrepareNamed(queryGetByUserID)
	if err != nil {
		logrus.Fatal(err)
	}

	stmtGetByUserIDDetailed, err = datastore.Get().Db.PrepareNamed(queryGetByUserIDDetailed)
	if err != nil {
		logrus.Fatal(err)
	}

	stmtGetByUserAndSourceDetail, err = datastore.Get().Db.PrepareNamed(queryGetByUserAndSourceDetail)
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

	stmtUpdateByMangaIDAndUserID, err = datastore.Get().Db.PrepareNamed(queryUpdateByMangaIDAndUserID)
	if err != nil {
		logrus.Fatal(err)
	}

	stmtGetRecentHistories, err = datastore.Get().Db.PrepareNamed(queryGetRecentHistories)
	if err != nil {
		logrus.Fatal(err)
	}
}
