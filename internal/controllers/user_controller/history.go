package user_controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/services/user_history"
	"github.com/umarkotak/animapu-api/internal/utils/render"
	"github.com/umarkotak/animapu-api/internal/utils/request"
)

type (
	PostHistoriesParams struct {
		Manga           models.Manga   `json:"manga"`
		LastReadChapter models.Chapter `json:"last_read_chapter"`
	}
)

func GetHistories(c *gin.Context) {
	user, err := request.ReqToUser(c.Request)
	if err != nil {
		logrus.WithContext(c.Request.Context()).Error(err)
		render.ErrorResponse(c.Request.Context(), c, err, false)
		return
	}

	mangaHistories, mangaHistoriesMap, err := user_history.GetReadHistories(c.Request.Context(), user)
	if err != nil {
		logrus.WithContext(c.Request.Context()).Error(err)
		render.ErrorResponse(c.Request.Context(), c, err, false)
		return
	}

	render.Response(
		c.Request.Context(), c,
		map[string]interface{}{
			"manga_histories":     mangaHistories,
			"manga_histories_map": mangaHistoriesMap,
		},
		nil, 200,
	)
}

func PostHistories(c *gin.Context) {
	var postHistoriesParams PostHistoriesParams
	c.BindJSON(&postHistoriesParams)

	user, err := request.ReqToUser(c.Request)
	if err != nil {
		logrus.WithContext(c.Request.Context()).Error(err)
		render.ErrorResponse(c.Request.Context(), c, err, false)
		return
	}

	manga, err := user_history.RecordHistory(c.Request.Context(), user, postHistoriesParams.Manga, postHistoriesParams.LastReadChapter)
	if err != nil {
		logrus.WithContext(c.Request.Context()).Error(err)
		render.ErrorResponse(c.Request.Context(), c, err, false)
		return
	}

	render.Response(
		c.Request.Context(), c,
		map[string]string{
			"id":                  manga.ID,
			"source":              manga.Source,
			"source_id":           manga.SourceID,
			"secondary_source":    manga.SecondarySource,
			"secondary_source_id": manga.SecondarySourceID,
		},
		nil, 200,
	)
}
