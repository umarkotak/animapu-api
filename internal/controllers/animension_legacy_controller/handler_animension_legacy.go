package animension_legacy_controller

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/schollz/progressbar/v3"
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

	animeIDs := []int64{}

	for i := 1; i <= int(page); i++ {
		seasonAnimeIDs, err := getAnimensionAnimesBySeason(
			c.Request.Context(), seasonID, page,
		)
		if err != nil {
			c.JSON(422, gin.H{"message": err.Error()})
			return
		}
		animeIDs = append(animeIDs, seasonAnimeIDs...)

		if len(seasonAnimeIDs) < 25 {
			break
		}
	}

	c.JSON(200, gin.H{"message": "ok", "data": map[string]any{
		"page":      page,
		"total":     len(animeIDs),
		"anime_ids": animeIDs,
	}})
}

func HandlerAnimensionSyncMultiSeasons(c *gin.Context) {
	seasonIDs := strings.Split(c.Query("season_ids"), ",")
	page := utils.StringMustInt64(c.Query("page"))

	for idx, seasonID := range seasonIDs {
		animeIDs := []int64{}

		for i := 1; i <= int(page); i++ {
			seasonAnimeIDs, err := getAnimensionAnimesBySeason(
				c.Request.Context(), seasonID, page,
			)
			if err != nil {
				c.JSON(422, gin.H{"message": err.Error()})
				return
			}
			animeIDs = append(animeIDs, seasonAnimeIDs...)

			if len(seasonAnimeIDs) < 25 {
				break
			}
			time.Sleep(1 * time.Second)
		}

		bar := progressbar.Default(int64(len(animeIDs)))
		logrus.WithContext(c).Infof("Start Processing %v (%v/%v)", seasonID, idx+1, len(seasonIDs))
		for _, oneAnimeID := range animeIDs {
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

			bar.Add(1)
			time.Sleep(2 * time.Second)
		}

		syncSeason()
		logrus.WithContext(c).Infof("Finish Processing %v", seasonID)
		time.Sleep(5 * time.Second)
	}

	c.JSON(200, gin.H{"message": "ok", "data": map[string]any{}})
}
