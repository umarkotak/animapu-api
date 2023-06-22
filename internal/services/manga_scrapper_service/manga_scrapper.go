package manga_scrapper_service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository"
	"github.com/umarkotak/animapu-api/internal/repository/manga_scrapper_repository"
	"github.com/umarkotak/animapu-api/internal/repository/mangamee_port"
)

func GetHome(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, models.Meta, error) {
	mangas := []models.Manga{}
	var err error

	cachedMangas, found := repository.GoCache().Get(queryParams.ToKey("page_latest"))
	if found {
		cachedMangasByte, err := json.Marshal(cachedMangas)
		if err == nil {
			err = json.Unmarshal(cachedMangasByte, &mangas)
			if err == nil {
				return mangas, models.Meta{FromCache: true}, nil
			}
		}
	}

	mangaScrapper, err := mangaScrapperGenerator(queryParams.Source)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return mangas, models.Meta{}, err
	}

	mangas, err = mangaScrapper.GetHome(ctx, queryParams)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return mangas, models.Meta{}, err
	}

	if len(mangas) > 0 {
		go repository.GoCache().Set(queryParams.ToKey("page_latest"), mangas, 5*time.Minute)
	}

	return mangas, models.Meta{}, nil
}

func GetDetail(ctx context.Context, queryParams models.QueryParams) (models.Manga, models.Meta, error) {
	manga := models.Manga{}
	var err error

	cachedManga, found := repository.GoCache().Get(queryParams.ToKey("page_detail"))
	if found {
		cachedMangaByte, err := json.Marshal(cachedManga)
		if err == nil {
			err = json.Unmarshal(cachedMangaByte, &manga)
			if err == nil {
				return manga, models.Meta{FromCache: true}, nil
			}
		}
	}

	mangaScrapper, err := mangaScrapperGenerator(queryParams.Source)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return manga, models.Meta{}, err
	}

	manga, err = mangaScrapper.GetDetail(ctx, queryParams)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return manga, models.Meta{}, err
	}

	if len(manga.Chapters) > 0 {
		go repository.GoCache().Set(queryParams.ToKey("page_detail"), manga, 5*time.Minute)
		// go cacheManga(context.Background(), queryParams.ToKey("page_detail"), manga)
	}

	return manga, models.Meta{}, nil
}

func GetSearch(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, models.Meta, error) {
	mangas := []models.Manga{}
	var err error

	mangaScrapper, err := mangaScrapperGenerator(queryParams.Source)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return mangas, models.Meta{}, err
	}

	mangas, err = mangaScrapper.GetSearch(ctx, queryParams)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return mangas, models.Meta{}, err
	}

	return mangas, models.Meta{}, nil
}

func GetChapter(ctx context.Context, queryParams models.QueryParams) (models.Chapter, models.Meta, error) {
	var err error
	chapter := models.Chapter{}

	// _, chapterJsonByte, err := repository.FbGet(ctx, queryParams.ToFbKey("page_read"))
	// if err == nil {
	// 	err = json.Unmarshal(chapterJsonByte, &chapter)
	// 	if err == nil {
	// 		return chapter, nil
	// 	}
	// }

	cachedChapter, found := repository.GoCache().Get(queryParams.ToKey("page_read"))
	if found {
		cachedChapterByte, err := json.Marshal(cachedChapter)
		if err == nil {
			err = json.Unmarshal(cachedChapterByte, &chapter)
			if err == nil {
				return chapter, models.Meta{FromCache: true}, nil
			}
		}
	}

	mangaScrapper, err := mangaScrapperGenerator(queryParams.Source)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return chapter, models.Meta{}, err
	}

	chapter, err = mangaScrapper.GetChapter(ctx, queryParams)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return chapter, models.Meta{}, err
	}

	if err == nil && len(chapter.ChapterImages) > 5 {
		// err = repository.FbSet(ctx, queryParams.ToFbKey("page_read"), chapter, time.Now().UTC().Add(30*24*time.Hour))
		// if err != nil {
		// 	logrus.WithContext(ctx).Error(err)
		// }
		go repository.GoCache().Set(queryParams.ToKey("page_read"), chapter, 5*time.Minute)
	}

	return chapter, models.Meta{}, nil
}

func mangaScrapperGenerator(mangaSource string) (models.MangaScrapper, error) {
	var mangaScrapper models.MangaScrapper

	switch mangaSource {
	case models.SOURCE_MANGABAT:
		mangaScrapper := manga_scrapper_repository.NewMangabat()
		return &mangaScrapper, nil
	case models.SOURCE_KLIKMANGA:
		mangaScrapper := manga_scrapper_repository.NewKlikmanga()
		return &mangaScrapper, nil
	case models.SOURCE_WEBTOONSID:
		mangaScrapper := manga_scrapper_repository.NewWebtoonsid()
		return &mangaScrapper, nil
	case models.SOURCE_MANGAREAD:
		mangaScrapper := mangamee_port.NewMangaread()
		return &mangaScrapper, nil
	case models.SOURCE_MANGATOWN:
		mangaScrapper := mangamee_port.NewMangatown()
		return &mangaScrapper, nil
	case models.SOURCE_MAIDMY:
		mangaScrapper := mangamee_port.NewMaidmy()
		return &mangaScrapper, nil
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
