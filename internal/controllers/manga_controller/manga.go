package manga_controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/services/manga_scrapper_service"
	"github.com/umarkotak/animapu-api/internal/utils/render"
)

func GetMangaLatest(c *gin.Context) {
	page, _ := strconv.ParseInt(c.Request.URL.Query().Get("page"), 10, 64)
	queryParams := models.QueryParams{
		Source: c.Param("manga_source"),
		Page:   page,
	}

	mangas, err := manga_scrapper_service.GetHome(c.Request.Context(), queryParams)
	if err != nil {
		logrus.WithContext(c.Request.Context()).Error(err)
		render.ErrorResponse(c.Request.Context(), c, err, false)
		return
	}

	render.Response(c.Request.Context(), c, mangas, nil, 200)
}

func GetMangaDetail(c *gin.Context) {
	queryParams := models.QueryParams{
		Source:            c.Param("manga_source"),
		SourceID:          c.Param("manga_id"),
		SecondarySourceID: c.Request.URL.Query().Get("secondary_source_id"),
	}

	manga, err := manga_scrapper_service.GetDetail(c.Request.Context(), queryParams)
	if err != nil {
		logrus.WithContext(c.Request.Context()).Error(err)
		render.ErrorResponse(c.Request.Context(), c, err, false)
		return
	}

	render.Response(c.Request.Context(), c, manga, nil, 200)
}

func ReadManga(c *gin.Context) {
	ctx := c.Request.Context()

	queryParams := models.QueryParams{
		Source:            c.Param("manga_source"),
		SourceID:          c.Param("manga_id"),
		SecondarySourceID: c.Request.URL.Query().Get("secondary_source_id"),
		ChapterID:         c.Param("chapter_id"),
	}

	chapter, err := manga_scrapper_service.GetChapter(ctx, queryParams)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		render.ErrorResponse(ctx, c, err, false)
		return
	}

	render.Response(ctx, c, chapter, nil, 200)
}

func SearchManga(c *gin.Context) {
	page, _ := strconv.ParseInt(c.Request.URL.Query().Get("page"), 10, 64)
	queryParams := models.QueryParams{
		Source: c.Param("manga_source"),
		Page:   page,
		Title:  c.Request.URL.Query().Get("title"),
	}

	mangas, err := manga_scrapper_service.GetSearch(c.Request.Context(), queryParams)
	if err != nil {
		logrus.WithContext(c.Request.Context()).Error(err)
		render.ErrorResponse(c.Request.Context(), c, err, false)
		return
	}

	render.Response(c.Request.Context(), c, mangas, nil, 200)
}
