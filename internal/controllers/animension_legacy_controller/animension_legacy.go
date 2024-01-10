package animension_legacy_controller

import (
	"regexp"

	"github.com/gin-gonic/gin"
)

type (
	ReqBody struct {
		AnimeID string `json:"anime_id"`
	}
)

var (
	AnimensionHost = "https://animension.to"
	ExtractEpRegex = regexp.MustCompile(`Episode (\d+)`)
)

func HandlerAnimensionQuickScrap(c *gin.Context) {
	animeDetail, err := quickScrapAnimeDetail(c.Request.Context(), ReqBody{
		AnimeID: c.Param("anime_id"),
	})
	if err != nil {
		c.JSON(422, gin.H{"message": err.Error()})
		return
	}

	err = saveAnime(animeDetail)
	if err != nil {
		c.JSON(422, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "ok", "data": animeDetail})
}

func HandlerSyncSeason(c *gin.Context) {
	syncSeason()

	c.JSON(200, gin.H{"message": "ok", "data": "ok"})
}
