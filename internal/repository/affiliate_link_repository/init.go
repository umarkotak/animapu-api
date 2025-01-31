package affiliate_link_repository

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/datastore"
)

var (
	allColumns = strings.Join([]string{
		"al.id",
		"al.created_at",
		"al.updated_at",
		"al.deleted_at",
		"al.short_link",
		"al.long_link",
		"al.image_url",
		"al.name",
	}, ", ")

	queryGetByID = fmt.Sprintf(`
		SELECT
			%s
		FROM affiliate_links al
		WHERE
			al.id = :id
			AND al.deleted_at IS NULL
	`, allColumns)

	queryGetByShortLink = fmt.Sprintf(`
		SELECT
			%s
		FROM affiliate_links al
		WHERE
			al.short_link = :short_link
			AND al.deleted_at IS NULL
	`, allColumns)

	queryGetRandom = fmt.Sprintf(`
		SELECT
			%s
		FROM affiliate_links al
		WHERE
			al.deleted_at IS NULL
		ORDER BY RANDOM()
		LIMIT :limit
	`, allColumns)

	queryInsert = `
		INSERT INTO affiliate_links (
			short_link,
			long_link,
			image_url,
			name
		) VALUES (
			:short_link,
			:long_link,
			:image_url,
			:name
		)
		RETURNING id
	`
)

var (
	stmtGetByID        *sqlx.NamedStmt
	stmtGetByShortLink *sqlx.NamedStmt
	stmtGetRandom      *sqlx.NamedStmt
	stmtInsert         *sqlx.NamedStmt
)

func Initialize() {
	var err error

	stmtGetByID, err = datastore.Get().Db.PrepareNamed(queryGetByID)
	if err != nil {
		logrus.Fatal(err)
	}

	stmtGetByShortLink, err = datastore.Get().Db.PrepareNamed(queryGetByShortLink)
	if err != nil {
		logrus.Fatal(err)
	}

	stmtGetRandom, err = datastore.Get().Db.PrepareNamed(queryGetRandom)
	if err != nil {
		logrus.Fatal(err)
	}

	stmtInsert, err = datastore.Get().Db.PrepareNamed(queryInsert)
	if err != nil {
		logrus.Fatal(err)
	}

}
