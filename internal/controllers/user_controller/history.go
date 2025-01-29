package user_controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/contract"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/services/manga_history_service"
	"github.com/umarkotak/animapu-api/internal/services/user_history_service"
	"github.com/umarkotak/animapu-api/internal/utils/common_ctx"
	"github.com/umarkotak/animapu-api/internal/utils/render"
	"github.com/umarkotak/animapu-api/internal/utils/request"
)

type (
	PostHistoriesParams struct {
		Manga           contract.Manga   `json:"manga"`
		LastReadChapter contract.Chapter `json:"last_read_chapter"`
	}
)

func GetHistories(c *gin.Context) {
	user, err := request.ReqToUser(c.Request)
	if err != nil {
		logrus.WithContext(c.Request.Context()).Error(err)
		render.ErrorResponse(c.Request.Context(), c, err, false)
		return
	}

	if user.Uid == "" {
		err = models.ErrUnauthorized
		render.ErrorResponse(c.Request.Context(), c, err, false)
		return
	}

	mangaHistories, mangaHistoriesMap, err := user_history_service.FirebaseGetReadHistories(c.Request.Context(), user)
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

func GetHistoriesV2(c *gin.Context) {
	commonCtx := common_ctx.GetFromGinCtx(c)

	mangaHistories, err := manga_history_service.GetHistories(c.Request.Context(), commonCtx.User, 1000, 0)
	if err != nil {
		logrus.WithContext(c.Request.Context()).Error(err)
		render.ErrorResponse(c.Request.Context(), c, err, false)
		return
	}

	render.Response(
		c.Request.Context(), c,
		map[string]any{
			"manga_histories": mangaHistories,
		},
		nil, 200,
	)
}

func FirebasePostHistories(c *gin.Context) {
	var postHistoriesParams PostHistoriesParams
	c.BindJSON(&postHistoriesParams)

	user, err := request.ReqToUser(c.Request)
	if err != nil {
		logrus.WithContext(c.Request.Context()).Error(err)
		render.ErrorResponse(c.Request.Context(), c, err, false)
		return
	}

	if user.Uid == "" {
		err = models.ErrUnauthorized
		render.ErrorResponse(c.Request.Context(), c, err, false)
		return
	}

	postHistoriesParams.Manga.Chapters = []contract.Chapter{}
	postHistoriesParams.Manga.Description = ""

	manga, err := user_history_service.FirebaseRecordHistory(c.Request.Context(), user, postHistoriesParams.Manga, postHistoriesParams.LastReadChapter)
	if err != nil {
		logrus.WithContext(c.Request.Context()).Error(err)
		render.ErrorResponse(c.Request.Context(), c, err, false)
		return
	}

	render.Response(
		c.Request.Context(), c,
		map[string]string{
			"id":        manga.ID,
			"source":    manga.Source,
			"source_id": manga.SourceID,
		},
		nil, 200,
	)
}
