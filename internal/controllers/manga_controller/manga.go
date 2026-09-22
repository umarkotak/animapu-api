package manga_controller

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/page"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/extension"
	"github.com/johnfercher/maroto/v2/pkg/props"
	"github.com/umarkotak/animapu-api/internal/utils/fiber_ctx"

	"github.com/sirupsen/logrus"
	app_config "github.com/umarkotak/animapu-api/config"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository/manga_repository"
	"github.com/umarkotak/animapu-api/internal/services/manga_scrapper_service"
	"github.com/umarkotak/animapu-api/internal/utils/common_ctx"
	"github.com/umarkotak/animapu-api/internal/utils/render"
	"github.com/umarkotak/animapu-api/internal/utils/utils"
)

type UpdateMangaTagsParams struct {
	Tags []string `json:"tags"`
	Mode string   `json:"mode"`
}

func GetMangaKids(c *fiber_ctx.Context) {
	mangas, err := manga_repository.GetByTag(c.Request.Context(), "for_kids:true")
	if err != nil {
		render.ErrorResponse(c.Request.Context(), c, err, false)
		return
	}

	render.Response(c.Request.Context(), c, mangas, nil, 200)
}

func UpdateMangaTags(c *fiber_ctx.Context) {
	ctx := c.Request.Context()
	if !app_config.IsAdminEmail(common_ctx.GetFromFiberCtx(c).User.Email.String) {
		render.ErrorResponse(ctx, c, models.ErrUnauthorized, true)
		return
	}

	mangaID, err := strconv.ParseInt(c.Param("manga_id"), 10, 64)
	if err != nil {
		render.ErrorResponse(ctx, c, models.ErrInvalidFormat, true)
		return
	}

	manga, err := manga_repository.GetByID(ctx, mangaID)
	if err != nil {
		render.ErrorResponse(ctx, c, models.ErrNotFound, true)
		return
	}
	updateMangaTags(c, manga)
}

func UpdateMangaTagsBySource(c *fiber_ctx.Context) {
	ctx := c.Request.Context()
	if !app_config.IsAdminEmail(common_ctx.GetFromFiberCtx(c).User.Email.String) {
		render.ErrorResponse(ctx, c, models.ErrUnauthorized, true)
		return
	}

	manga, err := manga_repository.GetBySourceAndSourceID(ctx, c.Param("manga_source"), c.Param("source_id"))
	if err != nil {
		render.ErrorResponse(ctx, c, models.ErrNotFound, true)
		return
	}
	updateMangaTags(c, manga)
}

func updateMangaTags(c *fiber_ctx.Context, manga models.Manga) {
	ctx := c.Request.Context()
	params := UpdateMangaTagsParams{}
	if err := c.BindJSON(&params); err != nil {
		render.ErrorResponse(ctx, c, models.ErrInvalidFormat, true)
		return
	}
	tags, ok := updateTags(manga.Tags, params.Tags, params.Mode)
	if !ok {
		render.ErrorResponse(ctx, c, models.ErrInvalidFormat, true)
		return
	}
	manga.Tags = tags
	if err := manga_repository.UpdateTags(ctx, manga); err != nil {
		render.ErrorResponse(ctx, c, err, true)
		return
	}

	render.Response(ctx, c, manga, nil, 200)
}

func updateTags(current, incoming []string, mode string) ([]string, bool) {
	switch mode {
	case "replace":
		return incoming, true
	case "add":
		seen := make(map[string]bool, len(current)+len(incoming))
		for _, tag := range current {
			seen[tag] = true
		}
		for _, tag := range incoming {
			if !seen[tag] {
				current = append(current, tag)
				seen[tag] = true
			}
		}
		return current, true
	case "remove":
		remove := make(map[string]bool, len(incoming))
		for _, tag := range incoming {
			remove[tag] = true
		}
		tags := current[:0]
		for _, tag := range current {
			if !remove[tag] {
				tags = append(tags, tag)
			}
		}
		return tags, true
	default:
		return nil, false
	}
}

func GetMangaLatest(c *fiber_ctx.Context) {
	commonCtx := common_ctx.GetFromFiberCtx(c)

	page, _ := strconv.ParseInt(c.Request.URL.Query().Get("page"), 10, 64)
	queryParams := models.QueryParams{
		Source: c.Param("manga_source"),
		Page:   page,
	}

	mangas, meta, err := manga_scrapper_service.GetHome(c.Request.Context(), queryParams)
	if err != nil {
		logrus.WithContext(c.Request.Context()).Error(err)
		render.ErrorResponse(c.Request.Context(), c, err, false)
		return
	}

	mangas = manga_scrapper_service.MultiInjectLibraryAndHistory(c.Request.Context(), commonCtx.User, mangas)

	c.Writer.Header().Set("Res-From-Cache", fmt.Sprintf("%v", meta.FromCache))
	render.Response(c.Request.Context(), c, mangas, nil, 200)
}

func GetMangaDetail(c *fiber_ctx.Context) {
	queryParams := models.QueryParams{
		Source:   c.Param("manga_source"),
		SourceID: c.Param("manga_id"),
	}

	manga, meta, err := manga_scrapper_service.GetDetail(c.Request.Context(), queryParams)
	if err != nil {
		logrus.WithContext(c.Request.Context()).Error(err)
		render.ErrorResponse(c.Request.Context(), c, err, false)
		return
	}

	c.Writer.Header().Set("Res-From-Cache", fmt.Sprintf("%v", meta.FromCache))
	render.Response(c.Request.Context(), c, manga, nil, 200)
}

func ReadManga(c *fiber_ctx.Context) {
	ctx := c.Request.Context()

	queryParams := models.QueryParams{
		Source:    c.Param("manga_source"),
		SourceID:  c.Param("manga_id"),
		ChapterID: c.Param("chapter_id"),
		User:      common_ctx.GetFromFiberCtx(c).User,
	}

	chapter, meta, err := manga_scrapper_service.GetChapter(ctx, queryParams)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		render.ErrorResponse(ctx, c, err, false)
		return
	}
	if err := manga_scrapper_service.MangaChapterSync(ctx, queryParams, chapter); err != nil {
		logrus.WithContext(ctx).Error(err)
	}

	c.Writer.Header().Set("Res-From-Cache", fmt.Sprintf("%v", meta.FromCache))
	render.Response(ctx, c, chapter, nil, 200)
}

func SearchManga(c *fiber_ctx.Context) {
	commonCtx := common_ctx.GetFromFiberCtx(c)

	page, _ := strconv.ParseInt(c.Request.URL.Query().Get("page"), 10, 64)
	queryParams := models.QueryParams{
		Source: c.Param("manga_source"),
		Page:   page,
		Title:  c.Request.URL.Query().Get("title"),
	}

	mangas, meta, err := manga_scrapper_service.GetSearch(c.Request.Context(), queryParams)
	if err != nil {
		logrus.WithContext(c.Request.Context()).Error(err)
		render.ErrorResponse(c.Request.Context(), c, err, false)
		return
	}

	mangas = manga_scrapper_service.MultiInjectLibraryAndHistory(c.Request.Context(), commonCtx.User, mangas)

	c.Writer.Header().Set("Res-From-Cache", fmt.Sprintf("%v", meta.FromCache))
	render.Response(c.Request.Context(), c, mangas, nil, 200)
}

func DownloadMangaChapter(c *fiber_ctx.Context) {
	ctx := c.Request.Context()

	queryParams := models.QueryParams{
		Source:    c.Param("manga_source"),
		SourceID:  c.Param("manga_id"),
		ChapterID: c.Param("chapter_id"),
		User:      common_ctx.GetFromFiberCtx(c).User,
	}

	manga, _, err := manga_scrapper_service.GetDetail(c.Request.Context(), queryParams)
	if err != nil {
		logrus.WithContext(c.Request.Context()).Error(err)
		render.ErrorResponse(c.Request.Context(), c, err, false)
		return
	}

	chapter, _, err := manga_scrapper_service.GetChapter(ctx, queryParams)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		render.ErrorResponse(ctx, c, err, false)
		return
	}

	cfg := config.NewBuilder().
		WithDimensions(800, 1200).
		Build()
	m := maroto.New(cfg)

	type ImagePacket struct {
		Idx       int
		ImageByte []byte
	}

	type ImageCompiled struct {
		mu           sync.Mutex
		ImagePackets []ImagePacket
	}

	imageCompiled := ImageCompiled{
		ImagePackets: []ImagePacket{},
	}

	wg := sync.WaitGroup{}

	logrus.WithContext(ctx).Infof("Start fetch: %+v", time.Now())
	for idx, chapterImage := range chapter.ChapterImages {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, imageUrl := range chapterImage.ImageUrls {
				logrus.WithContext(ctx).Infof("Start fetch image %v: %+v", idx, time.Now())
				imageByte, err := utils.ImageUrlToJpegByte(imageUrl)
				if err != nil {
					logrus.WithContext(ctx).Error(err)
					continue
				}

				imageCompiled.mu.Lock()
				defer imageCompiled.mu.Unlock()

				imageCompiled.ImagePackets = append(imageCompiled.ImagePackets, ImagePacket{
					Idx:       idx,
					ImageByte: imageByte,
				})
				logrus.WithContext(ctx).Infof("Finish fetch image %v: %+v", idx, time.Now())
				break
			}
		}()
	}
	wg.Wait()
	logrus.WithContext(ctx).Infof("Finish fetch: %+v", time.Now())

	sort.Slice(imageCompiled.ImagePackets, func(i, j int) bool {
		return imageCompiled.ImagePackets[i].Idx < imageCompiled.ImagePackets[j].Idx
	})

	for _, imagePacket := range imageCompiled.ImagePackets {
		m.AddPages(
			page.New().Add(
				image.NewAutoFromBytesRow(imagePacket.ImageByte, extension.Jpeg, props.Rect{}),
			),
		)
	}

	logrus.WithContext(ctx).Infof("Start gen: %+v", time.Now())
	document, err := m.Generate()
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		render.ErrorResponse(ctx, c, err, false)
		return
	}
	logrus.WithContext(ctx).Infof("Finish gen: %+v", time.Now())

	chapterTitle := fmt.Sprintf("%s - chapter %v", manga.Title, chapter.Number)

	c.Writer.Header().Set("Content-Type", "application/pdf")
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.pdf", chapterTitle)) // Customize filename
	c.Writer.Header().Set("Content-Length", fmt.Sprint(len(document.GetBytes())))

	c.Writer.Write(document.GetBytes())
}
