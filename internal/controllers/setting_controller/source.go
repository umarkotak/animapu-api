package setting_controller

import (
	"github.com/gin-gonic/gin"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/utils/render"
)

func GetAvailableSource(c *gin.Context) {
	render.Response(
		c.Request.Context(),
		c,
		models.MangaSources,
		nil,
		200,
	)
}

func GetAnimeAvailableSource(c *gin.Context) {
	render.Response(
		c.Request.Context(),
		c,
		models.AnimeSources,
		nil,
		200,
	)
}
