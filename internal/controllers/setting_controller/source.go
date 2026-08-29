package setting_controller

import (
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/utils/fiber_ctx"
	"github.com/umarkotak/animapu-api/internal/utils/render"
)

func GetAvailableSource(c *fiber_ctx.Context) {
	render.Response(
		c.Request.Context(),
		c,
		models.MangaSources,
		nil,
		200,
	)
}

func GetAnimeAvailableSource(c *fiber_ctx.Context) {
	render.Response(
		c.Request.Context(),
		c,
		models.AnimeSources,
		nil,
		200,
	)
}
