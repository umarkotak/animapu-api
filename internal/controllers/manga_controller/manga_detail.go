package manga_controller

import (
	"context"
	"encoding/json"
	"fmt"
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

	cacheKey := fmt.Sprintf("CACHED_MANGA:%v:%v:%v", queryParams.Source, queryParams.SourceID, queryParams.SecondarySourceID)

	cachedJson, err := repository.Redis().Get(c.Request.Context(), cacheKey).Result()
	if err != nil && err != redis.Nil {
		logrus.WithContext(c.Request.Context()).Error(err)
	}
	if err == nil {
		var cachedManga interface{}
		err = json.Unmarshal([]byte(cachedJson), &cachedManga)
		if err == nil {
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

	go cacheManga(context.Background(), cacheKey, manga)

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
