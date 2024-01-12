package anime_controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/services/anime_scrapper_service"
	"github.com/umarkotak/animapu-api/internal/utils/render"
)

func GetLatest(c *gin.Context) {
	ctx := c.Request.Context()

	page, _ := strconv.ParseInt(c.Query("page"), 10, 64)
	queryParams := models.AnimeQueryParams{
		Source: c.Param("anime_source"),
		Page:   page,
	}

	animes, meta, err := anime_scrapper_service.GetLatest(ctx, queryParams)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		render.ErrorResponse(ctx, c, err, true)
		return
	}

	c.Writer.Header().Set("Res-From-Cache", fmt.Sprintf("%v", meta.FromCache))
	render.Response(ctx, c, animes, nil, 200)
}

func GetSearch(c *gin.Context) {
	ctx := c.Request.Context()

	queryParams := models.AnimeQueryParams{
		Source: c.Param("anime_source"),
		Title:  c.Query("title"),
	}

	animes, meta, err := anime_scrapper_service.GetSearch(ctx, queryParams)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		render.ErrorResponse(ctx, c, err, true)
		return
	}

	c.Writer.Header().Set("Res-From-Cache", fmt.Sprintf("%v", meta.FromCache))
	render.Response(ctx, c, animes, nil, 200)
}

func GetPerSeason(c *gin.Context) {
	ctx := c.Request.Context()

	releaseYear, _ := strconv.ParseInt(c.Param("release_year"), 10, 64)
	queryParams := models.AnimeQueryParams{
		Source:        c.Param("anime_source"),
		ReleaseYear:   releaseYear,
		ReleaseSeason: c.Param("release_season"),
	}

	animes, meta, err := anime_scrapper_service.GetPerSeason(ctx, queryParams)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		render.ErrorResponse(ctx, c, err, true)
		return
	}

	c.Writer.Header().Set("Res-From-Cache", fmt.Sprintf("%v", meta.FromCache))
	render.Response(ctx, c, animes, nil, 200)
}

func GetDetail(c *gin.Context) {
	ctx := c.Request.Context()

	queryParams := models.AnimeQueryParams{
		Source:   c.Param("anime_source"),
		SourceID: c.Param("anime_id"),
	}

	anime, meta, err := anime_scrapper_service.GetDetail(ctx, queryParams)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		render.ErrorResponse(ctx, c, err, true)
		return
	}

	c.Writer.Header().Set("Res-From-Cache", fmt.Sprintf("%v", meta.FromCache))
	render.Response(ctx, c, anime, nil, 200)
}

func GetWatch(c *gin.Context) {
	ctx := c.Request.Context()

	queryParams := models.AnimeQueryParams{
		Source:       c.Param("anime_source"),
		SourceID:     c.Param("anime_id"),
		EpisodeID:    c.Param("episode_id"),
		WatchVersion: c.Request.URL.Query().Get("watch_version"),
	}

	episodeWatch, meta, err := anime_scrapper_service.Watch(ctx, queryParams)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		render.ErrorResponse(ctx, c, err, true)
		return
	}

	if episodeWatch.RawPageByte != nil {
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write(episodeWatch.RawPageByte)
		return
	}

	c.Writer.Header().Set("Res-From-Cache", fmt.Sprintf("%v", meta.FromCache))
	render.Response(ctx, c, episodeWatch, nil, 200)
}

func GetRandom(c *gin.Context) {
	ctx := c.Request.Context()

	queryParams := models.AnimeQueryParams{
		Source: c.Param("anime_source"),
	}

	animes, meta, err := anime_scrapper_service.GetRandom(ctx, queryParams)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		render.ErrorResponse(ctx, c, err, true)
		return
	}

	c.Writer.Header().Set("Res-From-Cache", fmt.Sprintf("%v", meta.FromCache))
	render.Response(ctx, c, animes, nil, 200)
}
