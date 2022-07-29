package manga_controller

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository"
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

	cachedManga, found := repository.GoCache().Get(queryParams.ToKey("page_detail"))
	if found {
		render.Response(c.Request.Context(), c, cachedManga, nil, 200)
		return
	}

	cachedJson, err := repository.Redis().Get(c.Request.Context(), queryParams.ToKey("page_detail")).Result()
	if err != nil && err != redis.Nil {
		logrus.WithContext(c.Request.Context()).Error(err)
	}
	if err == nil {
		var cachedManga interface{}
		err = json.Unmarshal([]byte(cachedJson), &cachedManga)
		if err == nil {
			go repository.GoCache().Set(queryParams.ToKey("page_detail"), manga, 30*time.Minute)

			render.Response(c.Request.Context(), c, cachedManga, nil, 200)
			return
		}
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

	if len(manga.Chapters) > 0 {
		go cacheManga(context.Background(), queryParams.ToKey("page_detail"), manga)

		go repository.GoCache().Set(queryParams.ToKey("page_detail"), manga, 30*time.Minute)
	}

	render.Response(c.Request.Context(), c, manga, nil, 200)
}

func cacheManga(ctx context.Context, cacheKey string, manga models.Manga) {
	mangaByte, err := json.Marshal(manga)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return
	}
	_, err = repository.Redis().Set(ctx, cacheKey, string(mangaByte), 30*time.Minute).Result()
	if err != nil {
		logrus.WithContext(ctx).Error(err)
	}
}
