package animension_legacy_controller

import (
	"fmt"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/utils/utils"
)

type (
	ReqBody struct {
		AnimeID string `json:"anime_id"`
	}
)

var (
	AnimensionBase = "animension.to"
	AnimensionHost = fmt.Sprintf("https://%s", AnimensionBase)
	ExtractEpRegex = regexp.MustCompile(`Episode (\d+)`)
)

func HandlerAnimensionQuickScrap(c *gin.Context) {
	reqBody := struct {
		AnimeIDs []int64 `json:"anime_ids"`
	}{}
	c.BindJSON(&reqBody)

	for _, oneAnimeID := range reqBody.AnimeIDs {
		animeDetail, err := quickScrapAnimeDetail(c.Request.Context(), ReqBody{
			AnimeID: fmt.Sprint(oneAnimeID),
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

		logrus.WithContext(c).Infof("%v DONE", oneAnimeID)
		time.Sleep(2 * time.Second)
	}

	c.JSON(200, gin.H{"message": "ok", "data": reqBody.AnimeIDs})
}

func HandlerSyncSeason(c *gin.Context) {
	syncSeason()

	c.JSON(200, gin.H{"message": "ok", "data": "ok"})
}

func HandlerAnimensionQuickScrapSeason(c *gin.Context) {
	seasonID := c.Param("season_id") // Eg: fall-2024, [summer, fall, winter, spring]-{year}
	page := utils.StringMustInt64(c.Query("page"))

	animeIDs, err := getAnimensionAnimesBySeason(
		c.Request.Context(), seasonID, page,
	)
	if err != nil {
		c.JSON(422, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "ok", "data": map[string]any{
		"page":      page,
		"total":     len(animeIDs),
		"anime_ids": animeIDs,
	}})
}
