package health_controller

import (
	"github.com/gin-gonic/gin"
	"github.com/umarkotak/animapu-api/internal/utils/logger"
	"github.com/umarkotak/animapu-api/internal/utils/render"
)

func GetHealth(c *gin.Context) {
	render.Response(
		c.Request.Context(),
		c,
		map[string]string{
			"health": "ok",
		},
		nil,
		200,
	)
}

func GetLogs(c *gin.Context) {
	render.Response(
		c.Request.Context(),
		c,
		logger.GlobalLog,
		nil,
		200,
	)
}
