package app

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/config"
	"github.com/umarkotak/animapu-api/internal/controllers/health_controller"
	"github.com/umarkotak/animapu-api/internal/controllers/manga_controller"
	"github.com/umarkotak/animapu-api/internal/controllers/proxy_controller"
	"github.com/umarkotak/animapu-api/internal/controllers/setting_controller"
	"github.com/umarkotak/animapu-api/internal/controllers/user_controller"
	"github.com/umarkotak/animapu-api/internal/repository"
	"github.com/umarkotak/animapu-api/internal/utils/logger"
)

func Initialize() {
	logger.Initialize()
	logrus.AddHook(&logger.AnimapuHook{})
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			// frame.File: For exact path
			return "", fmt.Sprintf("[%v:%v]:", path.Base(frame.File), frame.Line)
		},
	})
	// f, err := os.OpenFile("app.log", os.O_WRONLY|os.O_CREATE, 0755)
	// if err != nil {
	// 	logrus.Error(err)
	// }
	// logrus.SetOutput(f)

	config.Initialize()
	repository.Initialize()
}

func Start() {
	r := gin.New()
	r.Use(RequestID())
	r.Use(CORSMiddleware())
	r.Use(LogRequest())

	r.GET("/health", health_controller.GetHealth)
	r.GET("/logs", health_controller.GetLogs)

	r.GET("/mangas/sources", setting_controller.GetAvailableSource)

	r.GET("/mangas/:manga_source/latest", manga_controller.GetMangaLatest)
	r.GET("/mangas/:manga_source/detail/:manga_id", manga_controller.GetMangaDetail)
	r.GET("/mangas/:manga_source/read/:manga_id/:chapter_id", manga_controller.ReadManga)
	r.GET("/mangas/:manga_source/search", manga_controller.SearchManga)

	r.GET("/mangas/popular", manga_controller.GetMangaPopular)
	r.POST("/mangas/upvote", manga_controller.PostMangaUpvote)

	r.GET("/users/mangas/histories", user_controller.GetHistories)
	r.POST("/users/mangas/histories", user_controller.PostHistories)

	r.POST("/users/mangas/libraries", user_controller.PostLibrary)
	r.POST("/users/mangas/libraries/:source/:source_id/remove", user_controller.DeleteLibrary)
	r.GET("/users/mangas/libraries", user_controller.GetLibraries)
	r.POST("/users/mangas/libraries/sync", user_controller.SyncLibraries)

	r.GET("/mangas/mangabat/image_proxy/*url", proxy_controller.MangabatImage)
	r.GET("/mangas/webtoons/image_proxy/*url", proxy_controller.WebtoonsImage)
	r.GET("/mangas/fizmanga/image_proxy/*url", proxy_controller.FizmangaImage)
	r.GET("/mangas/klikmanga/image_proxy/*url", proxy_controller.KlikmangaImage)
	r.GET("/image_proxy", proxy_controller.GenericImage)

	port := os.Getenv("PORT")
	if port == "" {
		port = "6001"
	}
	r.Run(":" + port)
}
