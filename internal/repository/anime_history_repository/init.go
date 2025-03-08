package anime_history_repository

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/datastore"
)

var (
	allColumns = strings.Join([]string{
		"ah.id",
		"ah.created_at",
		"ah.updated_at",
		"ah.deleted_at",
		"ah.user_id",
		"ah.anime_id",
		"ah.episode_number",
		"ah.source_episode_id",
		"ah.frontend_path",
	}, ", ")

	queryGetByID = fmt.Sprintf(`
		SELECT
			%s
		FROM anime_histories ah
		WHERE
			ah.id = :id
			AND ah.deleted_at IS NULL
	`, allColumns)

	queryGetByAnimeIDAndUserID = fmt.Sprintf(`
		SELECT
			%s
		FROM anime_histories ah
		WHERE
			ah.user_id = :user_id
			AND ah.anime_id = :anime_id
			AND ah.deleted_at IS NULL
	`, allColumns)

	queryGetByUserID = fmt.Sprintf(`
		SELECT
			%s
		FROM anime_histories ah
		WHERE
			ah.user_id = :user_id
			AND ah.deleted_at IS NULL
	`, allColumns)

	queryGetByUserIDDetailed = fmt.Sprintf(`
		SELECT
			%s,
			a.source AS anime_source,
			a.source_id AS anime_source_id,
			a.title AS anime_title,
			a.cover_urls AS anime_cover_urls,
			a.latest_episode AS anime_latest_episode
		FROM anime_histories ah
		INNER JOIN animes a ON a.id = ah.anime_id
		WHERE
			ah.user_id = :user_id
			AND ah.deleted_at IS NULL
		ORDER BY ah.updated_at DESC
		LIMIT :limit OFFSET :offset
	`, allColumns)

	queryGetByUserAndSourceDetail = fmt.Sprintf(`
		SELECT
			%s,
			a.source AS anime_source,
			a.source_id AS anime_source_id,
			a.title AS anime_title,
			a.cover_urls AS anime_cover_urls,
			a.latest_episode AS anime_latest_episode
		FROM anime_histories ah
		INNER JOIN animes a ON a.id = ah.anime_id
		WHERE
			ah.user_id = :user_id
			AND a.source = ANY(:sources)
			AND a.source_id = ANY(:source_ids)
	`, allColumns)

	queryInsert = `
		INSERT INTO anime_histories (
			user_id,
			anime_id,
			episode_number,
			source_episode_id,
			frontend_path
		) VALUES (
			:user_id,
			:anime_id,
			:episode_number,
			:source_episode_id,
			:frontend_path
		)
		ON CONFLICT (user_id, anime_id)
		DO UPDATE SET
			episode_number = :episode_number,
			source_episode_id = :source_episode_id,
			frontend_path = :frontend_path,
			updated_at = NOW()
		RETURNING id
	`

	queryUpdate = `
		UPDATE anime_histories
		SET
			user_id = :user_id,
			anime_id = :anime_id,
			episode_number = :episode_number,
			source_episode_id = :source_episode_id,
			frontend_path = :frontend_path
		WHERE
			id = :id
	`

	queryUpdateByAnimeIDAndUserID = `
		UPDATE anime_histories
		SET
			episode_number = :episode_number,
			source_episode_id = :source_episode_id,
			frontend_path = :frontend_path
		WHERE
			user_id = :user_id
			AND anime_id = :anime_id
	`

	queryGetRecentHistories = fmt.Sprintf(`
		SELECT
			%s,
			a.source AS anime_source,
			a.source_id AS anime_source_id,
			a.title AS anime_title,
			a.cover_urls AS anime_cover_urls,
			a.latest_episode AS anime_latest_episode
		FROM anime_histories ah
		INNER JOIN animes a ON a.id = ah.anime_id
		ORDER BY ah.updated_at DESC
		LIMIT :limit OFFSET :offset
	`, allColumns)
)

var (
	stmtGetByID                  *sqlx.NamedStmt
	stmtGetByAnimeIDAndUserID    *sqlx.NamedStmt
	stmtGetByUserID              *sqlx.NamedStmt
	stmtGetByUserIDDetailed      *sqlx.NamedStmt
	stmtGetByUserAndSourceDetail *sqlx.NamedStmt
	stmtInsert                   *sqlx.NamedStmt
	stmtUpdate                   *sqlx.NamedStmt
	stmtUpdateByAnimeIDAndUserID *sqlx.NamedStmt
	stmtGetRecentHistories       *sqlx.NamedStmt
)

func Initialize() {
	var err error

	stmtGetByID, err = datastore.Get().Db.PrepareNamed(queryGetByID)
	if err != nil {
		logrus.Fatal(err)
	}

	stmtGetByAnimeIDAndUserID, err = datastore.Get().Db.PrepareNamed(queryGetByAnimeIDAndUserID)
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

	stmtUpdateByAnimeIDAndUserID, err = datastore.Get().Db.PrepareNamed(queryUpdateByAnimeIDAndUserID)
	if err != nil {
		logrus.Fatal(err)
	}

	stmtGetRecentHistories, err = datastore.Get().Db.PrepareNamed(queryGetRecentHistories)
	if err != nil {
		logrus.Fatal(err)
	}
}
