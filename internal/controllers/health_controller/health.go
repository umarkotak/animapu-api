package health_controller

import (
	"github.com/umarkotak/animapu-api/internal/repository"
	"github.com/umarkotak/animapu-api/internal/utils/fiber_ctx"
	"github.com/umarkotak/animapu-api/internal/utils/logger"
	"github.com/umarkotak/animapu-api/internal/utils/render"
)

func GetHealth(c *fiber_ctx.Context) {
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

func GetLogs(c *fiber_ctx.Context) {
	render.Response(
		c.Request.Context(),
		c,
		logger.GlobalLog,
		nil,
		200,
	)
}

func GetVisitorLogs(c *fiber_ctx.Context) {
	render.Response(
		c.Request.Context(),
		c,
		repository.GetLogVisitor(),
		nil,
		200,
	)
}
