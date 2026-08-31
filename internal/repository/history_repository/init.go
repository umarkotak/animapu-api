package history_repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/datastore"
)

const queryGetByUserID = `
	SELECT * FROM (
		SELECT
			'manga' AS media_type,
			m.source,
			m.source_id,
			m.title,
			m.cover_urls,
			m.latest_chapter AS latest_number,
			mh.chapter_number AS progress,
			mh.frontend_path AS last_link,
			mh.updated_at
		FROM manga_histories mh
		INNER JOIN mangas m ON m.id = mh.manga_id
		WHERE mh.user_id = :user_id AND mh.deleted_at IS NULL

		UNION ALL

		SELECT
			'anime' AS media_type,
			a.source,
			a.source_id,
			a.title,
			a.cover_urls,
			a.latest_episode AS latest_number,
			ah.episode_number AS progress,
			ah.frontend_path AS last_link,
			ah.updated_at
		FROM anime_histories ah
		INNER JOIN animes a ON a.id = ah.anime_id
		WHERE ah.user_id = :user_id AND ah.deleted_at IS NULL
	) histories
	ORDER BY updated_at DESC
	LIMIT :limit OFFSET :offset
`

var stmtGetByUserID *sqlx.NamedStmt

func Initialize() {
	var err error
	stmtGetByUserID, err = datastore.Get().Db.PrepareNamed(queryGetByUserID)
	if err != nil {
		logrus.Fatal(err)
	}
}
