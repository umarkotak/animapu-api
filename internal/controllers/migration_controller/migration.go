package migration_controller

import (
	"github.com/umarkotak/animapu-api/datastore"
	"github.com/umarkotak/animapu-api/internal/utils/fiber_ctx"
	"github.com/umarkotak/animapu-api/internal/utils/render"
)

func MigrateUp(c *fiber_ctx.Context) {
	datastore.MigrateUp()

	render.Response(
		c.Request.Context(), c,
		map[string]any{},
		nil, 200,
	)
}
