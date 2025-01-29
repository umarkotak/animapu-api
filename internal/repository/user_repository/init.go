package user_repository

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/datastore"
)

var (
	allColumns = strings.Join([]string{
		"u.id",
		"u.created_at",
		"u.updated_at",
		"u.deleted_at",
		"u.visitor_id",
		"u.guid",
		"u.email",
	}, ", ")

	queryGetByEmail = fmt.Sprintf(`
		SELECT
			%s
		FROM users u
		WHERE
			u.email = :email
			AND u.deleted_at IS NULL
	`, allColumns)

	queryGetByID = fmt.Sprintf(`
		SELECT
			%s
		FROM users u
		WHERE
			u.id = :id
			AND u.deleted_at IS NULL
	`, allColumns)

	queryGetByGuid = fmt.Sprintf(`
		SELECT
			%s
		FROM users u
		WHERE
			u.guid = :guid
			AND u.deleted_at IS NULL
	`, allColumns)

	queryGetByVisitorID = fmt.Sprintf(`
		SELECT
			%s
		FROM users u
		WHERE
			u.visitor_id = :visitor_id
			AND u.deleted_at IS NULL
	`, allColumns)

	queryInsert = `
		INSERT INTO users (
			visitor_id,
			guid,
			email
		) VALUES (
			:visitor_id,
			:guid,
			:email
		)
		ON CONFLICT (visitor_id)
		DO UPDATE SET
			visitor_id = :visitor_id
		RETURNING id
	`

	queryUpdate = `
		UPDATE users
		SET
			visitor_id = :visitor_id,
			guid = :guid,
			email = :email,
			updated_at = NOW()
		WHERE
			id = :id
	`
)

var (
	stmtGetByEmail     *sqlx.NamedStmt
	stmtGetByID        *sqlx.NamedStmt
	stmtGetByGuid      *sqlx.NamedStmt
	stmtGetByVisitorID *sqlx.NamedStmt
	stmtInsert         *sqlx.NamedStmt
	stmtUpdate         *sqlx.NamedStmt
)

func Initialize() {
	var err error

	stmtGetByEmail, err = datastore.Get().Db.PrepareNamed(queryGetByEmail)
	if err != nil {
		logrus.Fatal(err)
	}

	stmtGetByID, err = datastore.Get().Db.PrepareNamed(queryGetByID)
	if err != nil {
		logrus.Fatal(err)
	}

	stmtGetByGuid, err = datastore.Get().Db.PrepareNamed(queryGetByGuid)
	if err != nil {
		logrus.Fatal(err)
	}

	stmtGetByVisitorID, err = datastore.Get().Db.PrepareNamed(queryGetByVisitorID)
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
}
