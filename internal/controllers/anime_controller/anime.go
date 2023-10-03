package anime_controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/services/anime_scrapper_service"
	"github.com/umarkotak/animapu-api/internal/utils/render"
)

func GetWatch(c *gin.Context) {
	ctx := c.Request.Context()

	queryParams := models.AnimeQueryParams{
		Source:    c.Param("anime_source"),
		SourceID:  c.Param("anime_id"),
		EpisodeID: c.Param("episode_id"),
	}

	episode, meta, err := anime_scrapper_service.GetEpisode(ctx, queryParams)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		render.ErrorResponse(ctx, c, err, false)
		return
	}

	c.Writer.Header().Set("Res-From-Cache", fmt.Sprintf("%v", meta.FromCache))
	render.Response(ctx, c, episode, nil, 200)
}
