package user_controller

import (
	"github.com/gin-gonic/gin"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/utils/render"
)

type (
	PostHistoriesParams struct {
		Manga models.Manga `json:"manga"`
	}
)

func GetHistories(c *gin.Context) {

}

func PostHistories(c *gin.Context) {
	var postHistoriesParams PostHistoriesParams
	c.BindJSON(&postHistoriesParams)

	render.Response(
		c.Request.Context(), c,
		map[string]string{
			"id": postHistoriesParams.Manga.ID,
		},
		nil, 200,
	)
}
