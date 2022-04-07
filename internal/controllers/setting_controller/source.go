package setting_controller

import (
	"github.com/gin-gonic/gin"
	"github.com/umarkotak/animapu-api/internal/utils/render"
)

func GetAvailableSource(c *gin.Context) {
	render.Response(
		c.Request.Context(),
		c,
		[]string{
			"mangaupdates",
			"mangadex",
			"maidmy",
			"klikmanga",
			"mangareadorg",
			"mangabat",
		},
		nil,
		200,
	)
}
