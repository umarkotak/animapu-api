package manga_scrapper

import (
	"context"
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository"
	"github.com/umarkotak/animapu-api/internal/repository/manga_scrapper_repository"
)

func GetHome(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	mangas := []models.Manga{}
	var err error

	cachedMangas, found := repository.GoCache().Get(queryParams.ToKey("page_latest"))
	if found {
		cachedMangasByte, err := json.Marshal(cachedMangas)
		if err == nil {
			err = json.Unmarshal(cachedMangasByte, &mangas)
			if err == nil {
				return mangas, nil
			}
		}
	}

	mangaScrapper, err := mangaScrapperGenerator(queryParams.Source)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return mangas, err
	}

	mangas, err = mangaScrapper.GetHome(ctx, queryParams)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return mangas, err
	}

	if len(mangas) > 0 {
		go repository.GoCache().Set(queryParams.ToKey("page_latest"), mangas, 5*time.Minute)
	}

	return mangas, nil
}

func GetDetail(ctx context.Context, queryParams models.QueryParams) (models.Manga, error) {
	manga := models.Manga{}
	var err error

	cachedManga, found := repository.GoCache().Get(queryParams.ToKey("page_detail"))
	if found {
		cachedMangaByte, err := json.Marshal(cachedManga)
		if err == nil {
			err = json.Unmarshal(cachedMangaByte, &manga)
			if err == nil {
				return manga, nil
			}
		}
	}

	mangaScrapper, err := mangaScrapperGenerator(queryParams.Source)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return manga, err
	}

	manga, err = mangaScrapper.GetDetail(ctx, queryParams)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return manga, err
	}

	if len(manga.Chapters) > 0 {
		go cacheManga(context.Background(), queryParams.ToKey("page_detail"), manga)
	}

	repository.GoCache().Set(queryParams.ToKey("page_detail"), manga, 30*time.Minute)

	return manga, nil
}

func GetSearch(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, error) {
	mangas := []models.Manga{}
	var err error

	mangaScrapper, err := mangaScrapperGenerator(queryParams.Source)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return mangas, err
	}

	mangas, err = mangaScrapper.GetSearch(ctx, queryParams)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return mangas, err
	}

	return mangas, nil
}

func GetChapter(ctx context.Context, queryParams models.QueryParams) (models.Chapter, error) {
	var err error
	chapter := models.Chapter{}

	// _, chapterJsonByte, err := repository.FbGet(ctx, queryParams.ToFbKey("page_read"))
	// if err == nil {
	// 	err = json.Unmarshal(chapterJsonByte, &chapter)
	// 	if err == nil {
	// 		return chapter, nil
	// 	}
	// }

	mangaScrapper, err := mangaScrapperGenerator(queryParams.Source)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return chapter, err
	}

	chapter, err = mangaScrapper.GetChapter(ctx, queryParams)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return chapter, err
	}

	if err == nil && len(chapter.ChapterImages) > 5 {
		err = repository.FbSet(ctx, queryParams.ToFbKey("page_read"), chapter, time.Now().UTC().Add(30*24*time.Hour))
		if err != nil {
			logrus.WithContext(ctx).Error(err)
		}
	}

	return chapter, nil
}

func mangaScrapperGenerator(mangaSource string) (models.MangaScrapper, error) {
	var mangaScrapper models.MangaScrapper

	switch mangaSource {
	case "mangabat":
		mangabat := manga_scrapper_repository.NewMangabat()
		return &mangabat, nil
	}

	return mangaScrapper, models.ErrMangaSourceNotFound
}

func cacheManga(ctx context.Context, cacheKey string, manga models.Manga) {
	mangaByte, err := json.Marshal(manga)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return
	}
	_, err = repository.Redis().Set(ctx, cacheKey, string(mangaByte), 30*time.Minute).Result()
	if err != nil {
		logrus.WithContext(ctx).Error(err)
	}
}