package manga_controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository"
	"github.com/umarkotak/animapu-api/internal/utils/render"
)

func PostMangaFollower(c *gin.Context) {
	manga := models.Manga{}
	err := c.BindJSON(&manga)
	if err != nil {
		logrus.WithContext(c.Request.Context()).Error(err)
		render.ErrorResponse(c.Request.Context(), c, err, false)
		return
	}
	manga.Chapters = []models.Chapter{}
	manga, err = repository.FbAddFollowManga(c.Request.Context(), manga)
	if err != nil {
		logrus.WithContext(c.Request.Context()).Error(err)
		render.ErrorResponse(c.Request.Context(), c, err, false)
		return
	}
	render.Response(c.Request.Context(), c, manga, nil, 200)
}
