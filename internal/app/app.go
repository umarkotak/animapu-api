package app

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/config"
	"github.com/umarkotak/animapu-api/datastore"
	"github.com/umarkotak/animapu-api/internal/controllers/affiliate_link_controller"
	"github.com/umarkotak/animapu-api/internal/controllers/anime_controller"
	"github.com/umarkotak/animapu-api/internal/controllers/health_controller"
	"github.com/umarkotak/animapu-api/internal/controllers/manga_controller"
	"github.com/umarkotak/animapu-api/internal/controllers/migration_controller"
	"github.com/umarkotak/animapu-api/internal/controllers/proxy_controller"
	"github.com/umarkotak/animapu-api/internal/controllers/setting_controller"
	"github.com/umarkotak/animapu-api/internal/controllers/user_controller"
	"github.com/umarkotak/animapu-api/internal/repository"
	"github.com/umarkotak/animapu-api/internal/repository/affiliate_link_repository"
	"github.com/umarkotak/animapu-api/internal/repository/manga_chapter_repository"
	"github.com/umarkotak/animapu-api/internal/repository/manga_history_repository"
	"github.com/umarkotak/animapu-api/internal/repository/manga_library_repository"
	"github.com/umarkotak/animapu-api/internal/repository/manga_repository"
	"github.com/umarkotak/animapu-api/internal/repository/manga_scrapper_repository"
	"github.com/umarkotak/animapu-api/internal/repository/user_repository"
	"github.com/umarkotak/animapu-api/internal/utils/logger"
)

func Initialize() {
	logger.Initialize()
	logrus.AddHook(&logger.AnimapuHook{})
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logger.Formatter{})

	config.Initialize()
	repository.Initialize()
	datastore.Initialize()

	user_repository.Initialize()
	manga_repository.Initialize()
	manga_chapter_repository.Initialize()
	manga_history_repository.Initialize()
	manga_library_repository.Initialize()
	affiliate_link_repository.Initialize()
	manga_scrapper_repository.InitializeAsuraComic()
}

func Start() {
	r := gin.New()
	r.Use(RequestID())
	r.Use(CORSMiddleware())
	r.Use(LogRequest())
	r.Use(LogVisitor())
	r.Use(CommonCtx())
	r.Use(RegisterUser())

	r.GET("/migrate_up", migration_controller.MigrateUp)

	r.GET("/health", health_controller.GetHealth)
	r.GET("/logs", health_controller.GetLogs)
	r.GET("/visitor_logs", health_controller.GetVisitorLogs)

	r.GET("/mangas/sources", setting_controller.GetAvailableSource)
	r.GET("/animes/sources", setting_controller.GetAnimeAvailableSource)

	r.GET("/mangas/:manga_source/latest", manga_controller.GetMangaLatest)
	r.GET("/mangas/:manga_source/detail/:manga_id", manga_controller.GetMangaDetail)
	r.GET("/mangas/:manga_source/read/:manga_id/:chapter_id", manga_controller.ReadManga)
	r.GET("/mangas/:manga_source/read/:manga_id/:chapter_id/manga_chapter.pdf", manga_controller.DownloadMangaChapter)
	r.GET("/mangas/:manga_source/search", manga_controller.SearchManga)

	r.GET("/users/mangas/histories", user_controller.GetHistories)
	r.GET("/users/mangas/histories_v2", user_controller.GetHistoriesV2)
	r.POST("/users/mangas/histories", user_controller.FirebasePostHistories)
	r.GET("/users/mangas/activities", user_controller.GetUserMangaActivities)

	r.POST("/users/mangas/libraries/:source/:source_id/add", user_controller.AddLibrary)
	r.POST("/users/mangas/libraries/:source/:source_id/remove", user_controller.DeleteLibrary)
	r.GET("/users/mangas/libraries", user_controller.GetLibraries)

	r.GET("/mangas/mangabat/image_proxy/*url", proxy_controller.MangabatImage)
	r.GET("/mangas/klikmanga/image_proxy/*url", proxy_controller.KlikmangaImage)
	r.GET("/mangas/komikindo/image_proxy/*url", proxy_controller.KomikindoImage)
	r.GET("/image_proxy", proxy_controller.GenericImage)

	r.GET("/animes/:anime_source/latest", anime_controller.GetLatest)
	r.GET("/animes/:anime_source/search", anime_controller.GetSearch)
	r.GET("/animes/:anime_source/random", anime_controller.GetRandom)
	r.GET("/animes/:anime_source/season/:release_year/:release_season", anime_controller.GetPerSeason)
	r.GET("/animes/:anime_source/detail/:anime_id", anime_controller.GetDetail)
	r.GET("/animes/:anime_source/watch/:anime_id/:episode_id", anime_controller.GetWatch)

	r.POST("/affiliate_links/tokopedia/add", affiliate_link_controller.AddTokopediaAffiliateLink)
	r.GET("/affiliate_links/random", affiliate_link_controller.GetRandom)
	r.GET("/affiliate_links", affiliate_link_controller.GetList)

	defer manga_scrapper_repository.DeferAsuraComic()
	r.Run(":" + config.Get().Port)
}
