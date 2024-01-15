package manga_scrapper_service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository"
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

	iterator := 5
	for i := 0; i <= iterator; i++ {
		mangas, err = mangaScrapper.GetHome(ctx, queryParams)
		if err != nil {
			logrus.WithContext(ctx).Error(err)
			time.Sleep(200 * time.Millisecond)
			continue
		}
		break
	}

	if len(mangas) > 0 {
		go repository.GoCache().Set(queryParams.ToKey("page_latest"), mangas, 30*time.Minute)
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

	iterator := 5
	for i := 0; i <= iterator; i++ {
		manga, err = mangaScrapper.GetDetail(ctx, queryParams)
		if err != nil {
			logrus.WithContext(ctx).Error(err)
			time.Sleep(200 * time.Millisecond)
			continue
		}
		break
	}

	if len(manga.Chapters) > 0 {
		go repository.GoCache().Set(queryParams.ToKey("page_detail"), manga, 12*time.Hour)
		// go cacheManga(context.Background(), queryParams.ToKey("page_detail"), manga)
	}

	return manga, models.Meta{}, nil
}

func GetSearch(ctx context.Context, queryParams models.QueryParams) ([]models.Manga, models.Meta, error) {
	mangas := []models.Manga{}
	var err error

	cachedMangas, found := repository.GoCache().Get(queryParams.ToKey("page_search"))
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

	mangas, err = mangaScrapper.GetSearch(ctx, queryParams)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return mangas, models.Meta{}, err
	}

	if len(mangas) > 0 {
		go repository.GoCache().Set(queryParams.ToKey("page_search"), mangas, 30*24*time.Hour)
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
		go repository.GoCache().Set(queryParams.ToKey("page_read"), chapter, 30*24*time.Hour)
	}

	return chapter, models.Meta{}, nil
}
