package app

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/controllers/health_controller"
	"github.com/umarkotak/animapu-api/internal/controllers/manga_controller"
	"github.com/umarkotak/animapu-api/internal/utils/logger"
)

func Initialize() {
	logger.Initialize()
	logrus.AddHook(&logger.AnimapuHook{})
}

func Start() {
	r := gin.New()
	r.Use(RequestID())
	r.Use(CORSMiddleware())

	r.GET("/health", health_controller.GetHealth)
	r.GET("/logs", health_controller.GetLogs)

	r.GET("/mangas/:source/home", manga_controller.GetMangaHome)

	port := os.Getenv("PORT")
	if port == "" {
		port = "6000"
	}
	r.Run(":" + port)
}
