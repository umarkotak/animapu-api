package manga_chapter_repository

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/datastore"
)

var (
	allColumns = strings.Join([]string{
		"mc.id",
		"mc.created_at",
		"mc.updated_at",
		"mc.deleted_at",
		"mc.manga_id",
		"mc.source_chapter_id",
		"mc.chapter_number",
		"mc.image_urls",
	}, ", ")

	queryGetByID = fmt.Sprintf(`
		SELECT
			%s
		FROM manga_chapters mc
		WHERE
			mc.id = :id
			AND mc.deleted_at IS NULL
	`, allColumns)

	queryGetByMangaIDAndSourceChapterID = fmt.Sprintf(`
		SELECT
			%s
		FROM manga_chapters mc
		WHERE
			mc.manga_id = :manga_id
			AND mc.source_chapter_id = :source_chapter_id
			AND mc.deleted_at IS NULL
	`, allColumns)

	queryInsert = `
		INSERT INTO manga_chapters (
			manga_id,
			source_chapter_id,
			chapter_number,
			image_urls
		) VALUES (
			:manga_id,
			:source_chapter_id,
			:chapter_number,
			:image_urls
		)
		ON CONFLICT (manga_id, source_chapter_id)
		DO UPDATE SET
			chapter_number = :chapter_number,
			image_urls = :image_urls
		RETURNING id
	`

	queryUpdate = `
		UPDATE manga_chapters
		SET
			manga_id = :manga_id,
			source_chapter_id = :source_chapter_id,
			chapter_number = :chapter_number,
			image_urls = :image_urls
		WHERE
			id = :id
	`

	queryUpdateByMangaIDAndSourceChapterID = `
		UPDATE manga_chapters
		SET
			chapter_number = :chapter_number,
			image_urls = :image_urls
		WHERE
			manga_id = :manga_id
			AND source_chapter_id = :source_chapter_id
	`
)

var (
	stmtGetByID                           *sqlx.NamedStmt
	stmtGetByMangaIDAndSourceChapterID    *sqlx.NamedStmt
	stmtInsert                            *sqlx.NamedStmt
	stmtUpdate                            *sqlx.NamedStmt
	stmtUpdateByMangaIDAndSourceChapterID *sqlx.NamedStmt
)

func Initialize() {
	var err error

	stmtGetByID, err = datastore.Get().Db.PrepareNamed(queryGetByID)
	if err != nil {
		logrus.Fatal(err)
	}

	stmtGetByMangaIDAndSourceChapterID, err = datastore.Get().Db.PrepareNamed(queryGetByMangaIDAndSourceChapterID)
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

	stmtUpdateByMangaIDAndSourceChapterID, err = datastore.Get().Db.PrepareNamed(queryUpdateByMangaIDAndSourceChapterID)
	if err != nil {
		logrus.Fatal(err)
	}
}
