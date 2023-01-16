package health_controller

import (
	"github.com/gin-gonic/gin"
	"github.com/umarkotak/animapu-api/internal/repository"
	"github.com/umarkotak/animapu-api/internal/utils/logger"
	"github.com/umarkotak/animapu-api/internal/utils/render"
)

func GetHealth(c *gin.Context) {
	render.Response(
		c.Request.Context(),
		c,
		map[string]string{
			"version": "1.0.0",
			"health":  "ok",
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

func GetVisitorLogs(c *gin.Context) {
	render.Response(
		c.Request.Context(),
		c,
		repository.GetLogVisitor(),
		nil,
		200,
	)
}
