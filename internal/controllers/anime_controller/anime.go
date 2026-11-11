package anime_controller

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
	app_config "github.com/umarkotak/animapu-api/config"
	"github.com/umarkotak/animapu-api/internal/contract"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository/anime_repository"
	"github.com/umarkotak/animapu-api/internal/services/anime_scrapper_service"
	"github.com/umarkotak/animapu-api/internal/utils/common_ctx"
	"github.com/umarkotak/animapu-api/internal/utils/fiber_ctx"
	"github.com/umarkotak/animapu-api/internal/utils/render"
	"github.com/umarkotak/animapu-api/internal/utils/tag_utils"
)

func GetLatest(c *fiber_ctx.Context) {
	ctx := c.Request.Context()

	commonCtx := common_ctx.GetFromFiberCtx(c)

	page, _ := strconv.ParseInt(c.Query("page"), 10, 64)
	queryParams := models.AnimeQueryParams{
		Source:        c.Param("anime_source"),
		Page:          page,
		SpecialFilter: c.Query("special_filter"),
	}

	animes, meta, err := anime_scrapper_service.GetList(ctx, queryParams)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		render.ErrorResponse(ctx, c, err, true)
		return
	}

	animes = anime_scrapper_service.MultiInjectLibraryAndHistoryForLatest(c.Request.Context(), commonCtx.User, animes)

	c.Writer.Header().Set("Res-From-Cache", fmt.Sprintf("%v", meta.FromCache))
	render.Response(ctx, c, animes, nil, 200)
}

func GetSearch(c *fiber_ctx.Context) {
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

func GetPerSeason(c *fiber_ctx.Context) {
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

func GetDetail(c *fiber_ctx.Context) {
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

	stored, err := anime_repository.GetBySourceAndSourceID(ctx, queryParams.Source, queryParams.SourceID)
	if err != nil && err != sql.ErrNoRows {
		render.ErrorResponse(ctx, c, err, true)
		return
	}
	anime.Tags = stored.Tags

	c.Writer.Header().Set("Res-From-Cache", fmt.Sprintf("%v", meta.FromCache))
	render.Response(ctx, c, anime, nil, 200)
}

func GetWatch(c *fiber_ctx.Context) {
	ctx := c.Request.Context()

	queryParams := models.AnimeQueryParams{
		Source:          c.Param("anime_source"),
		SourceID:        c.Param("anime_id"),
		EpisodeID:       c.Param("episode_id"),
		WatchVersion:    c.Request.URL.Query().Get("watch_version"),
		Resolution:      c.Request.URL.Query().Get("resolution"),
		StreamIdx:       c.Request.URL.Query().Get("stream_idx"),
		ManualServerOpt: c.Request.URL.Query().Get("manual_server_opt"),
		User:            common_ctx.GetFromFiberCtx(c).User,
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

func GetRandom(c *fiber_ctx.Context) {
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

func GetAnimeKids(c *fiber_ctx.Context) {
	ctx := c.Request.Context()
	animes, err := anime_repository.GetByTag(ctx, "for_kids:true")
	if err != nil {
		render.ErrorResponse(ctx, c, err, false)
		return
	}
	result := make([]contract.Anime, 0, len(animes))
	for _, anime := range animes {
		result = append(result, contract.Anime{
			ID: anime.SourceID, Source: anime.Source, Title: anime.Title,
			CoverUrls: anime.CoverUrls, LatestEpisode: anime.LatestEpisode, Tags: anime.Tags,
		})
	}
	render.Response(ctx, c, result, nil, 200)
}

func UpdateAnimeTagsBySource(c *fiber_ctx.Context) {
	ctx := c.Request.Context()
	if !app_config.IsAdminEmail(common_ctx.GetFromFiberCtx(c).User.Email.String) {
		render.ErrorResponse(ctx, c, models.ErrUnauthorized, true)
		return
	}
	params := struct {
		Tags []string `json:"tags"`
		Mode string   `json:"mode"`
	}{}
	if err := c.BindJSON(&params); err != nil {
		render.ErrorResponse(ctx, c, models.ErrInvalidFormat, true)
		return
	}
	if _, ok := tag_utils.Update(nil, params.Tags, params.Mode); !ok {
		render.ErrorResponse(ctx, c, models.ErrInvalidFormat, true)
		return
	}
	source, sourceID := c.Param("anime_source"), c.Param("source_id")
	anime, err := anime_repository.GetBySourceAndSourceID(ctx, source, sourceID)
	if err == sql.ErrNoRows {
		// Search and season results may not have been saved yet.
		detail, _, detailErr := anime_scrapper_service.GetDetail(ctx, models.AnimeQueryParams{Source: source, SourceID: sourceID})
		if detailErr != nil {
			render.ErrorResponse(ctx, c, detailErr, true)
			return
		}
		_, err = anime_repository.Insert(ctx, nil, models.Anime{
			Source: source, SourceID: sourceID, Title: detail.Title,
			CoverUrls: detail.CoverUrls, LatestEpisode: detail.LatestEpisode,
		})
		if err == nil {
			anime, err = anime_repository.GetBySourceAndSourceID(ctx, source, sourceID)
		}
	}
	if err != nil {
		render.ErrorResponse(ctx, c, err, true)
		return
	}
	tags, _ := tag_utils.Update(anime.Tags, params.Tags, params.Mode)
	anime.Tags = tags
	if anime.Tags == nil {
		anime.Tags = []string{}
	}
	if err := anime_repository.UpdateTags(ctx, anime); err != nil {
		render.ErrorResponse(ctx, c, err, true)
		return
	}
	render.Response(ctx, c, anime, nil, 200)
}

func GetCapabilities(c *fiber_ctx.Context) {
	ctx := c.Request.Context()
	capabilities, err := anime_scrapper_service.GetCapabilities(ctx, models.AnimeQueryParams{Source: c.Param("anime_source")})
	if err != nil {
		render.ErrorResponse(ctx, c, err, true)
		return
	}
	render.Response(ctx, c, capabilities, nil, http.StatusOK)
}
