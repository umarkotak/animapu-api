package manga_controller

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/page"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/extension"
	"github.com/johnfercher/maroto/v2/pkg/props"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/services/manga_scrapper_service"
	"github.com/umarkotak/animapu-api/internal/utils/common_ctx"
	"github.com/umarkotak/animapu-api/internal/utils/render"
	"github.com/umarkotak/animapu-api/internal/utils/utils"
)

func GetMangaLatest(c *gin.Context) {
	commonCtx := common_ctx.GetFromGinCtx(c)

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

func GetMangaDetail(c *gin.Context) {
	queryParams := models.QueryParams{
		Source:            c.Param("manga_source"),
		SourceID:          c.Param("manga_id"),
		SecondarySourceID: c.Request.URL.Query().Get("secondary_source_id"),
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

func ReadManga(c *gin.Context) {
	ctx := c.Request.Context()

	queryParams := models.QueryParams{
		Source:            c.Param("manga_source"),
		SourceID:          c.Param("manga_id"),
		SecondarySourceID: c.Request.URL.Query().Get("secondary_source_id"),
		ChapterID:         c.Param("chapter_id"),
		User:              common_ctx.GetFromGinCtx(c).User,
	}

	chapter, meta, err := manga_scrapper_service.GetChapter(ctx, queryParams)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		render.ErrorResponse(ctx, c, err, false)
		return
	}

	c.Writer.Header().Set("Res-From-Cache", fmt.Sprintf("%v", meta.FromCache))
	render.Response(ctx, c, chapter, nil, 200)
}

func SearchManga(c *gin.Context) {
	commonCtx := common_ctx.GetFromGinCtx(c)

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

func DownloadMangaChapter(c *gin.Context) {
	ctx := c.Request.Context()

	queryParams := models.QueryParams{
		Source:            c.Param("manga_source"),
		SourceID:          c.Param("manga_id"),
		SecondarySourceID: c.Request.URL.Query().Get("secondary_source_id"),
		ChapterID:         c.Param("chapter_id"),
		User:              common_ctx.GetFromGinCtx(c).User,
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

	c.Writer.Header().Set("Content-Type", "application/pdf")
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.pdf", fmt.Sprintf("%v - chapter %v", manga.Title, chapter.Number))) // Customize filename
	c.Writer.Header().Set("Content-Length", fmt.Sprint(len(document.GetBytes())))

	c.Writer.Write(document.GetBytes())
}
