package manga_controller

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository"
	"github.com/umarkotak/animapu-api/internal/services/manga_scrapper"
	"github.com/umarkotak/animapu-api/internal/utils/render"
)

func GetMangaLatest(c *gin.Context) {
	logrus.WithContext(c.Request.Context()).Error(models.ErrMangaSourceNotFound)

	page, _ := strconv.ParseInt(c.Request.URL.Query().Get("page"), 10, 64)
	queryParams := models.QueryParams{
		Source: c.Param("manga_source"),
		Page:   page,
	}

	mangas := []models.Manga{}
	var err error

	cachedMangas, found := repository.GoCache().Get(queryParams.ToKey("page_latest"))
	if found {
		render.Response(c.Request.Context(), c, cachedMangas, nil, 200)
		return
	}

	switch queryParams.Source {
	case "mangaupdates":
		mangas, err = manga_scrapper.GetMangaupdatesLatestManga(c.Request.Context(), queryParams)
	case "mangabat":
		mangas, err = manga_scrapper.GetMangabatLatestManga(c.Request.Context(), queryParams)
	case "klikmanga":
		mangas, err = manga_scrapper.GetKlikmangaLatestManga(c.Request.Context(), queryParams)
	case "webtoonsid":
		mangas, err = manga_scrapper.GetWebtoonsidLatestManga(c.Request.Context(), queryParams)
	case "fizmanga":
		mangas, err = manga_scrapper.GetFizmangaLatestManga(c.Request.Context(), queryParams)
	case "mangahub":
		mangas, err = manga_scrapper.GetMangahubLatestManga(c.Request.Context(), queryParams)
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

	if len(mangas) > 0 {
		go repository.GoCache().Set(queryParams.ToKey("page_latest"), mangas, 5*time.Minute)
	}

	render.Response(c.Request.Context(), c, mangas, nil, 200)
}
