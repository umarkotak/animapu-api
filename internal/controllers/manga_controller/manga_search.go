package manga_controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/services/manga_scrapper"
	"github.com/umarkotak/animapu-api/internal/utils/render"
)

func SearchManga(c *gin.Context) {
	page, _ := strconv.ParseInt(c.Request.URL.Query().Get("page"), 10, 64)
	queryParams := models.QueryParams{
		Source: c.Param("manga_source"),
		Page:   page,
		Title:  c.Request.URL.Query().Get("title"),
	}

	mangas := []models.Manga{}
	var err error

	switch queryParams.Source {
	case "mangaupdates":
		mangas, err = manga_scrapper.GetMangaupdatesByQuery(c.Request.Context(), queryParams)
	case "mangabat":
		mangas, err = manga_scrapper.GetMangabatByQuery(c.Request.Context(), queryParams)
	case "klikmanga":
		mangas, err = manga_scrapper.GetKlikmangaByQuery(c.Request.Context(), queryParams)
	case "webtoonsid":
		mangas, err = manga_scrapper.GetWebtoonsidByQuery(c.Request.Context(), queryParams)
	case "fizmanga":
		mangas, err = manga_scrapper.GetFizmangaByQuery(c.Request.Context(), queryParams)
	case "mangadex":
		render.ErrorResponse(c.Request.Context(), c, models.ErrMangaSourceNotImplemented, false)
		return
	case "maidmy":
		render.ErrorResponse(c.Request.Context(), c, models.ErrMangaSourceNotImplemented, false)
		return
	case "mangareadorg":
		render.ErrorResponse(c.Request.Context(), c, models.ErrMangaSourceNotImplemented, false)
		return
	default:
		render.ErrorResponse(c.Request.Context(), c, models.ErrMangaSourceNotFound, false)
		return
	}

	if err != nil {
		logrus.WithContext(c.Request.Context()).Error(err)
		render.ErrorResponse(c.Request.Context(), c, err, false)
		return
	}

	render.Response(c.Request.Context(), c, mangas, nil, 200)
}
