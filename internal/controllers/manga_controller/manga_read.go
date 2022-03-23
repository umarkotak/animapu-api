package manga_controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/services/manga_scrapper"
	"github.com/umarkotak/animapu-api/internal/utils/render"
)

func ReadManga(c *gin.Context) {
	queryParams := models.QueryParams{
		Source:            c.Param("manga_source"),
		SourceID:          c.Param("manga_id"),
		SecondarySourceID: c.Request.URL.Query().Get("secondary_source_id"),
		ChapterID:         c.Param("chapter_id"),
	}

	chapter := models.Chapter{}
	var err error

	switch queryParams.Source {
	case "mangaupdates":
		chapter, err = manga_scrapper.GetMangaupdatesDetailChapter(c.Request.Context(), queryParams)
	case "mangadex":
		render.ErrorResponse(c.Request.Context(), c, models.ErrMangaSourceNotImplemented)
		return
	case "maidmy":
		render.ErrorResponse(c.Request.Context(), c, models.ErrMangaSourceNotImplemented)
		return
	case "klikmanga":
		render.ErrorResponse(c.Request.Context(), c, models.ErrMangaSourceNotImplemented)
		return
	case "mangareadorg":
		render.ErrorResponse(c.Request.Context(), c, models.ErrMangaSourceNotImplemented)
		return
	case "mangabat":
		render.ErrorResponse(c.Request.Context(), c, models.ErrMangaSourceNotImplemented)
		return
	default:
		render.ErrorResponse(c.Request.Context(), c, models.ErrMangaSourceNotFound)
		return
	}

	if err != nil {
		logrus.WithContext(c.Request.Context()).Error(err)
		render.ErrorResponse(c.Request.Context(), c, err)
		return
	}

	render.Response(c.Request.Context(), c, chapter, nil, 200)
}
