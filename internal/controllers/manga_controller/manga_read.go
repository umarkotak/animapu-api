package manga_controller

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository"
	"github.com/umarkotak/animapu-api/internal/services/manga_scrapper"
	"github.com/umarkotak/animapu-api/internal/utils/render"
)

func ReadManga(c *gin.Context) {
	ctx := c.Request.Context()
	var err error
	chapter := models.Chapter{}

	queryParams := models.QueryParams{
		Source:            c.Param("manga_source"),
		SourceID:          c.Param("manga_id"),
		SecondarySourceID: c.Request.URL.Query().Get("secondary_source_id"),
		ChapterID:         c.Param("chapter_id"),
	}

	_, chapterJsonByte, err := repository.FbGet(ctx, queryParams.ToKey("page_read"))
	if err == nil {
		json.Unmarshal(chapterJsonByte, &chapter)
		render.Response(ctx, c, chapter, nil, 200)
		return
	}

	switch queryParams.Source {
	case "mangaupdates":
		chapter, err = manga_scrapper.GetMangaupdatesDetailChapter(ctx, queryParams)
	case "mangabat":
		chapter, err = manga_scrapper.GetMangabatDetailChapter(ctx, queryParams)
	case "klikmanga":
		chapter, err = manga_scrapper.GetKlikmangaDetailChapter(ctx, queryParams)
	case "webtoonsid":
		chapter, err = manga_scrapper.GetWebtoonsidDetailChapter(ctx, queryParams)
	case "fizmanga":
		chapter, err = manga_scrapper.GetFizmangaDetailChapter(ctx, queryParams)
	case "mangahub":
		chapter, err = manga_scrapper.GetMangahubDetailChapter(ctx, queryParams)
	case "mangadex":
		render.ErrorResponse(ctx, c, models.ErrMangaSourceNotImplemented, false)
		return
	case "maidmy":
		render.ErrorResponse(ctx, c, models.ErrMangaSourceNotImplemented, false)
		return
	case "mangareadorg":
		render.ErrorResponse(ctx, c, models.ErrMangaSourceNotImplemented, false)
		return
	default:
		render.ErrorResponse(ctx, c, models.ErrMangaSourceNotFound, false)
		return
	}

	if err != nil {
		logrus.WithContext(ctx).Error(err)
		render.ErrorResponse(ctx, c, err, false)
		return
	}

	if err == nil && len(chapter.ChapterImages) > 5 {
		err = repository.FbSet(ctx, queryParams.ToKey("page_read"), chapter, time.Now().UTC().Add(30*24*time.Hour))
		if err != nil {
			logrus.WithContext(ctx).Error(err)
		}
	}

	render.Response(ctx, c, chapter, nil, 200)
}
