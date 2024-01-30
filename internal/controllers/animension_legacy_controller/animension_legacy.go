package animension_legacy_controller

import (
	"fmt"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
