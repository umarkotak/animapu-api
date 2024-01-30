package app

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/config"
	"github.com/umarkotak/animapu-api/internal/controllers/anime_controller"
	"github.com/umarkotak/animapu-api/internal/controllers/animension_legacy_controller"
	"github.com/umarkotak/animapu-api/internal/controllers/health_controller"
	"github.com/umarkotak/animapu-api/internal/controllers/manga_controller"
	"github.com/umarkotak/animapu-api/internal/controllers/otakudesu_legacy_controller"
	"github.com/umarkotak/animapu-api/internal/controllers/proxy_controller"
	"github.com/umarkotak/animapu-api/internal/controllers/setting_controller"
	"github.com/umarkotak/animapu-api/internal/controllers/user_controller"
	"github.com/umarkotak/animapu-api/internal/local_db"
	"github.com/umarkotak/animapu-api/internal/repository"
	"github.com/umarkotak/animapu-api/internal/singelton"
	"github.com/umarkotak/animapu-api/internal/utils/logger"
)

func Initialize() {
	logger.Initialize()
	logrus.AddHook(&logger.AnimapuHook{})
	logrus.SetReportCaller(true)
	// logrus.SetFormatter(&logrus.TextFormatter{
	// 	CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
	// 		// frame.File: For exact path
	// 		return "", fmt.Sprintf("[%v:%v]:", frame.File, frame.Line)
	// 	},
	// })
	logrus.SetFormatter(&logger.Formatter{})
	// f, err := os.OpenFile("app.log", os.O_WRONLY|os.O_CREATE, 0755)
	// if err != nil {
	// 	logrus.Error(err)
	// }
	// logrus.SetOutput(f)

	config.Initialize()
	repository.Initialize()
	local_db.Initialize()
	singelton.Initialize()
}

func Start() {
	r := gin.New()
	r.Use(RequestID())
	r.Use(CORSMiddleware())
	r.Use(LogRequest())
	r.Use(LogVisitor())

	r.GET("/health", health_controller.GetHealth)
	r.GET("/logs", health_controller.GetLogs)
	r.GET("/visitor_logs", health_controller.GetVisitorLogs)

	r.GET("/mangas/sources", setting_controller.GetAvailableSource)

	r.GET("/mangas/:manga_source/latest", manga_controller.GetMangaLatest)
	r.GET("/mangas/:manga_source/detail/:manga_id", manga_controller.GetMangaDetail)
	r.GET("/mangas/:manga_source/read/:manga_id/:chapter_id", manga_controller.ReadManga)
	r.GET("/mangas/:manga_source/search", manga_controller.SearchManga)

	r.GET("/mangas/popular", manga_controller.GetMangaPopular)
	r.POST("/mangas/upvote", manga_controller.PostMangaUpvote)
	r.POST("/mangas/follow", manga_controller.PostMangaFollower)
	r.GET("/mangas/comments/disqus", manga_controller.GetMangaCommentsDisqus)

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
	r.GET("/mangas/komikindo/image_proxy/*url", proxy_controller.KomikindoImage)
	r.GET("/image_proxy", proxy_controller.GenericImage)

	r.GET("/animes/:anime_source/latest", anime_controller.GetLatest)
	r.GET("/animes/:anime_source/search", anime_controller.GetSearch)
	r.GET("/animes/:anime_source/random", anime_controller.GetRandom)
	r.GET("/animes/:anime_source/season/:release_year/:release_season", anime_controller.GetPerSeason)
	r.GET("/animes/:anime_source/detail/:anime_id", anime_controller.GetDetail)
	r.GET("/animes/:anime_source/watch/:anime_id/:episode_id", anime_controller.GetWatch)

	r.POST("/animension/quick_scrap/:anime_id", animension_legacy_controller.HandlerAnimensionQuickScrap)
	r.POST("/animension/sync_season", animension_legacy_controller.HandlerSyncSeason)

	r.POST("/otakudesu/scrap_otakudesu_all_animes", otakudesu_legacy_controller.HandlerAnimensionQuickScrap)

	r.Run(":" + config.Get().Port)
}
