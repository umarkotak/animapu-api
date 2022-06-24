package manga_controller

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/config"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/services/manga_scrapper"
	"github.com/umarkotak/animapu-api/internal/utils/render"
)

func GetMangaDetail(c *gin.Context) {
	queryParams := models.QueryParams{
		Source:            c.Param("manga_source"),
		SourceID:          c.Param("manga_id"),
		SecondarySourceID: c.Request.URL.Query().Get("secondary_source_id"),
	}

	manga := models.Manga{}
	var err error

	cacheKey := fmt.Sprintf("CACHED_MANGA:%v:%v:%v", queryParams.Source, queryParams.SourceID, queryParams.SecondarySourceID)
	cachedManga, found := config.Get().CacheObj.Get(cacheKey)
	if found {
		render.Response(c.Request.Context(), c, cachedManga, nil, 200)
		return
	}

	switch queryParams.Source {
	case "mangaupdates":
		manga, err = manga_scrapper.GetMangaupdatesDetailManga(c.Request.Context(), queryParams)
	case "mangabat":
		manga, err = manga_scrapper.GetMangabatDetailManga(c.Request.Context(), queryParams)
	case "klikmanga":
		manga, err = manga_scrapper.GetKlikmangaDetailManga(c.Request.Context(), queryParams)
	case "webtoonsid":
		manga, err = manga_scrapper.GetWebtoonsidDetailManga(c.Request.Context(), queryParams)
	case "fizmanga":
		manga, err = manga_scrapper.GetFizmangaDetailManga(c.Request.Context(), queryParams)
	case "mangahub":
		manga, err = manga_scrapper.GetMangahubDetailManga(c.Request.Context(), queryParams)
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

	config.Get().CacheObj.Set(cacheKey, manga, 30*time.Minute)

	render.Response(c.Request.Context(), c, manga, nil, 200)
}
