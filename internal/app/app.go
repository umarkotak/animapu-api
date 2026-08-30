package app

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
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
	"github.com/umarkotak/animapu-api/internal/repository/affiliate_link_repository"
	"github.com/umarkotak/animapu-api/internal/repository/anime_history_repository"
	"github.com/umarkotak/animapu-api/internal/repository/anime_repository"
	"github.com/umarkotak/animapu-api/internal/repository/manga_chapter_repository"
	"github.com/umarkotak/animapu-api/internal/repository/manga_history_repository"
	"github.com/umarkotak/animapu-api/internal/repository/manga_library_repository"
	"github.com/umarkotak/animapu-api/internal/repository/manga_repository"
	"github.com/umarkotak/animapu-api/internal/repository/user_repository"
	"github.com/umarkotak/animapu-api/internal/utils/fiber_ctx"
	"github.com/umarkotak/animapu-api/internal/utils/logger"
)

func Initialize() {
	logger.Initialize()
	logrus.AddHook(&logger.AnimapuHook{})
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logger.Formatter{})

	config.Initialize()
	datastore.Initialize()

	user_repository.Initialize()
	manga_repository.Initialize()
	manga_chapter_repository.Initialize()
	manga_history_repository.Initialize()
	manga_library_repository.Initialize()
	affiliate_link_repository.Initialize()
	anime_repository.Initialize()
	anime_history_repository.Initialize()
}

func Start() error {
	r := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})
	r.Use(RequestID())
	r.Use(CORSMiddleware())
	r.Use(LogRequest())
	r.Use(LogVisitor())
	r.Use(CommonCtx())
	r.Use(RegisterUser())

	r.Get("/migrate_up", fiber_ctx.Wrap(migration_controller.MigrateUp))

	r.Get("/health", fiber_ctx.Wrap(health_controller.GetHealth))
	r.Get("/logs", fiber_ctx.Wrap(health_controller.GetLogs))
	r.Get("/visitor_logs", fiber_ctx.Wrap(health_controller.GetVisitorLogs))

	r.Get("/mangas/sources", fiber_ctx.Wrap(setting_controller.GetAvailableSource))
	r.Get("/animes/sources", fiber_ctx.Wrap(setting_controller.GetAnimeAvailableSource))

	r.Get("/mangas/:manga_source/latest", fiber_ctx.Wrap(manga_controller.GetMangaLatest))
	r.Get("/mangas/:manga_source/detail/:manga_id", fiber_ctx.Wrap(manga_controller.GetMangaDetail))
	r.Get("/mangas/:manga_source/read/:manga_id/:chapter_id", fiber_ctx.Wrap(manga_controller.ReadManga))
	r.Get("/mangas/:manga_source/read/:manga_id/:chapter_id/manga_chapter.pdf", fiber_ctx.Wrap(manga_controller.DownloadMangaChapter))
	r.Get("/mangas/:manga_source/search", fiber_ctx.Wrap(manga_controller.SearchManga))

	r.Get("/users/mangas/histories_v2", fiber_ctx.Wrap(user_controller.GetHistoriesV2))
	r.Get("/users/mangas/activities", fiber_ctx.Wrap(user_controller.GetUserMangaActivities))
	r.Get("/users/animes/histories", fiber_ctx.Wrap(user_controller.GetAnimeHistories))
	r.Get("/users/animes/activities", fiber_ctx.Wrap(user_controller.GetUserAnimeActivities))

	r.Post("/users/mangas/libraries/:source/:source_id/add", fiber_ctx.Wrap(user_controller.AddLibrary))
	r.Post("/users/mangas/libraries/:source/:source_id/remove", fiber_ctx.Wrap(user_controller.DeleteLibrary))
	r.Get("/users/mangas/libraries", fiber_ctx.Wrap(user_controller.GetLibraries))

	r.Get("/mangas/mangabat/image_proxy/*url", fiber_ctx.Wrap(proxy_controller.MangabatImage))
	r.Get("/mangas/weeb_central/image_proxy/*url", fiber_ctx.Wrap(proxy_controller.WeebCentralImage))
	r.Get("/mangas/klikmanga/image_proxy/*url", fiber_ctx.Wrap(proxy_controller.KlikmangaImage))
	r.Get("/mangas/komikindo/image_proxy/*url", fiber_ctx.Wrap(proxy_controller.KomikindoImage))
	r.Get("/image_proxy", fiber_ctx.Wrap(proxy_controller.GenericImage))

	r.Get("/animes/:anime_source/latest", fiber_ctx.Wrap(anime_controller.GetLatest))
	r.Get("/animes/:anime_source/search", fiber_ctx.Wrap(anime_controller.GetSearch))
	r.Get("/animes/:anime_source/random", fiber_ctx.Wrap(anime_controller.GetRandom))
	r.Get("/animes/:anime_source/season/:release_year/:release_season", fiber_ctx.Wrap(anime_controller.GetPerSeason))
	r.Get("/animes/:anime_source/detail/:anime_id", fiber_ctx.Wrap(anime_controller.GetDetail))
	r.Get("/animes/:anime_source/watch/:anime_id/:episode_id", fiber_ctx.Wrap(anime_controller.GetWatch))

	r.Post("/affiliate_links/tokopedia/add", fiber_ctx.Wrap(affiliate_link_controller.AddTokopediaAffiliateLink))
	r.Get("/affiliate_links/random", fiber_ctx.Wrap(affiliate_link_controller.GetRandom))
	r.Get("/affiliate_links", fiber_ctx.Wrap(affiliate_link_controller.GetList))

	return r.Listen(":" + config.Get().Port)
}
