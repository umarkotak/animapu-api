package user_controller

import (
	"github.com/gin-gonic/gin"
	"github.com/umarkotak/animapu-api/internal/contract"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository/manga_library_repository"
	"github.com/umarkotak/animapu-api/internal/repository/manga_repository"
	"github.com/umarkotak/animapu-api/internal/services/manga_library_service"
	"github.com/umarkotak/animapu-api/internal/services/manga_scrapper_service"
	"github.com/umarkotak/animapu-api/internal/utils/common_ctx"
	"github.com/umarkotak/animapu-api/internal/utils/render"
)

type (
	SyncLibrariesParams struct {
		Mangas []contract.Manga `json:"mangas"`
	}

	PostLibraryParams struct {
		User     models.User `json:"-"`
		Source   string      `json:"-"`
		SourceID string      `json:"-"`
	}
)

func GetLibraries(c *gin.Context) {
	ctx := c.Request.Context()

	user := common_ctx.GetFromGinCtx(c).User

	mangas, err := manga_library_service.GetLibraries(ctx, user, 1000, 1)
	if err != nil {
		render.ErrorResponse(ctx, c, err, true)
		return
	}

	mangas = manga_scrapper_service.MultiInjectLibraryAndHistoryForLibrary(ctx, user, mangas)

	render.Response(ctx, c, mangas, nil, 200)
}

func AddLibrary(c *gin.Context) {
	ctx := c.Request.Context()

	postLibraryParams := PostLibraryParams{
		User:     common_ctx.GetFromGinCtx(c).User,
		Source:   c.Param("source"),
		SourceID: c.Param("source_id"),
	}

	manga_scrapper_service.GetDetail(ctx, models.QueryParams{
		Source:   postLibraryParams.Source,
		SourceID: postLibraryParams.SourceID,
	})

	manga, err := manga_repository.GetBySourceAndSourceID(ctx, postLibraryParams.Source, postLibraryParams.SourceID)
	if err != nil {
		render.ErrorResponse(ctx, c, err, true)
		return
	}

	_, err = manga_library_repository.Insert(ctx, nil, models.MangaHistory{
		UserID:  postLibraryParams.User.ID,
		MangaID: manga.ID,
	})
	if err != nil {
		render.ErrorResponse(ctx, c, err, true)
		return
	}

	render.Response(
		ctx, c,
		map[string]string{
			"library_saved": "success",
		},
		nil, 200,
	)
}

func DeleteLibrary(c *gin.Context) {
	ctx := c.Request.Context()

	postLibraryParams := PostLibraryParams{
		User:     common_ctx.GetFromGinCtx(c).User,
		Source:   c.Param("source"),
		SourceID: c.Param("source_id"),
	}

	manga, err := manga_repository.GetBySourceAndSourceID(ctx, postLibraryParams.Source, postLibraryParams.SourceID)
	if err != nil {
		render.ErrorResponse(ctx, c, err, true)
		return
	}

	err = manga_library_repository.Delete(ctx, nil, postLibraryParams.User.ID, manga.ID)
	if err != nil {
		render.ErrorResponse(ctx, c, err, true)
		return
	}

	render.Response(
		c.Request.Context(), c,
		map[string]string{
			"library_deleted": "success",
		},
		nil, 200,
	)
}
