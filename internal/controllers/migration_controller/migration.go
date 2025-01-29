package migration_controller

import (
	"github.com/gin-gonic/gin"
	"github.com/umarkotak/animapu-api/datastore"
	"github.com/umarkotak/animapu-api/internal/utils/render"
)

func MigrateUp(c *gin.Context) {
	datastore.MigrateUp()

	render.Response(
		c.Request.Context(), c,
		map[string]any{},
		nil, 200,
	)
}
