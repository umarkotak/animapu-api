package app

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/config"
	"github.com/umarkotak/animapu-api/internal/controllers/health_controller"
	"github.com/umarkotak/animapu-api/internal/controllers/manga_controller"
	"github.com/umarkotak/animapu-api/internal/utils/logger"
)

func Initialize() {
	logger.Initialize()
	logrus.AddHook(&logger.AnimapuHook{})
	config.Initialize()
}

func Start() {
	r := gin.New()
	r.Use(RequestID())
	r.Use(CORSMiddleware())

	r.GET("/health", health_controller.GetHealth)
	r.GET("/logs", health_controller.GetLogs)

	r.GET("/mangas/:manga_source/latest", manga_controller.GetMangaLatest)
	r.GET("/mangas/:manga_source/detail/:manga_id", manga_controller.GetMangaDetail)
	r.GET("/mangas/:manga_source/read/:manga_id/:chapter_id", manga_controller.ReadManga)
	r.GET("/mangas/:manga_source/search", manga_controller.SearchManga)

	port := os.Getenv("PORT")
	if port == "" {
		port = "6000"
	}
	r.Run(":" + port)
}
