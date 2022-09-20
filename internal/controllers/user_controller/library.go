package user_controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/utils/render"
	"github.com/umarkotak/animapu-api/internal/utils/request"
)

type (
	SyncLibrariesParams struct {
		Mangas []models.Manga `json:"mangas"`
	}

	PostLibraryParams struct {
		Manga models.Manga `json:"manga"`
	}
)

func GetLibraries(c *gin.Context) {
	user, err := request.ReqToUser(c.Request)
	if err != nil {
		logrus.WithContext(c.Request.Context()).Error(err)
		render.ErrorResponse(c.Request.Context(), c, err, false)
		return
	}

	if user.Uid == "" {
		err = models.ErrUnauthorized
		render.ErrorResponse(c.Request.Context(), c, err, false)
		return
	}

	// TODO: Implement logic
	mangas := []models.Manga{}

	render.Response(c.Request.Context(), c, mangas, nil, 200)
}

func PostLibrary(c *gin.Context) {
	var postLibraryParams PostLibraryParams
	c.BindJSON(&postLibraryParams)

	user, err := request.ReqToUser(c.Request)
	if err != nil {
		logrus.WithContext(c.Request.Context()).Error(err)
		render.ErrorResponse(c.Request.Context(), c, err, false)
		return
	}

	if user.Uid == "" {
		err = models.ErrUnauthorized
		render.ErrorResponse(c.Request.Context(), c, err, false)
		return
	}

	// TODO: Implement logic

	render.Response(
		c.Request.Context(), c,
		map[string]string{
			"library_saved": "success",
		},
		nil, 200,
	)
}

func DeleteLibrary(c *gin.Context) {
	user, err := request.ReqToUser(c.Request)
	if err != nil {
		logrus.WithContext(c.Request.Context()).Error(err)
		render.ErrorResponse(c.Request.Context(), c, err, false)
		return
	}

	if user.Uid == "" {
		err = models.ErrUnauthorized
		render.ErrorResponse(c.Request.Context(), c, err, false)
		return
	}

	// TODO: Implement logic

	render.Response(
		c.Request.Context(), c,
		map[string]string{
			"library_deleted": "success",
		},
		nil, 200,
	)
}

func SyncLibraries(c *gin.Context) {
	var syncLibrariesParams SyncLibrariesParams
	c.BindJSON(&syncLibrariesParams)

	user, err := request.ReqToUser(c.Request)
	if err != nil {
		logrus.WithContext(c.Request.Context()).Error(err)
		render.ErrorResponse(c.Request.Context(), c, err, false)
		return
	}

	if user.Uid == "" {
		err = models.ErrUnauthorized
		render.ErrorResponse(c.Request.Context(), c, err, false)
		return
	}

	// TODO: Implement sync library services
	// err = user_library_service.Sync(c.Request.Context(), user, syncLibrariesParams.Mangas)
	if err != nil {
		logrus.WithContext(c.Request.Context()).Error(err)
		render.ErrorResponse(c.Request.Context(), c, err, false)
		return
	}

	render.Response(
		c.Request.Context(), c,
		map[string]string{
			"synced": "success",
		},
		nil, 200,
	)
}
