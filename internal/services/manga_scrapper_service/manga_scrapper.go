package manga_scrapper_service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/datastore"
	"github.com/umarkotak/animapu-api/internal/contract"
	"github.com/umarkotak/animapu-api/internal/models"
)

func GetHome(ctx context.Context, queryParams models.QueryParams) ([]contract.Manga, models.Meta, error) {
	mangas := []contract.Manga{}
	var err error

	cachedMangas, found := datastore.Get().GoCache.Get(queryParams.ToKey("page_latest"))
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
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return mangas, models.Meta{}, err
	}

	if len(mangas) > 0 {
		go datastore.Get().GoCache.Set(queryParams.ToKey("page_latest"), mangas, 10*time.Minute)
		go MangaSync(context.Background(), mangas)
	}

	return mangas, models.Meta{}, nil
}

func GetDetail(ctx context.Context, queryParams models.QueryParams) (contract.Manga, models.Meta, error) {
	manga := contract.Manga{}
	var err error

	cachedManga, found := datastore.Get().GoCache.Get(queryParams.ToKey("page_detail"))
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
		go datastore.Get().GoCache.Set(queryParams.ToKey("page_detail"), manga, 15*time.Hour)
	}

	err = MangaSync(ctx, []contract.Manga{manga})
	if err != nil {
		logrus.WithContext(ctx).Error(err)
	}

	return manga, models.Meta{}, nil
}

func GetSearch(ctx context.Context, queryParams models.QueryParams) ([]contract.Manga, models.Meta, error) {
	mangas := []contract.Manga{}
	var err error

	cachedMangas, found := datastore.Get().GoCache.Get(queryParams.ToKey("page_search"))
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
		go datastore.Get().GoCache.Set(queryParams.ToKey("page_search"), mangas, 30*24*time.Hour)
		go MangaSync(context.Background(), mangas)
	}

	return mangas, models.Meta{}, nil
}

func GetChapter(ctx context.Context, queryParams models.QueryParams) (contract.Chapter, models.Meta, error) {
	var err error
	chapter := contract.Chapter{}

	cachedChapter, found := datastore.Get().GoCache.Get(queryParams.ToKey("page_read"))
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

	iterator := 0
	for i := 0; i <= iterator; i++ {
		chapter, err = mangaScrapper.GetChapter(ctx, queryParams)
		if err != nil {
			logrus.WithContext(ctx).Error(err)
			time.Sleep(200 * time.Millisecond)
			continue
		}

		break
	}
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return chapter, models.Meta{}, err
	}

	if len(chapter.ChapterImages) > 5 {
		go datastore.Get().GoCache.Set(queryParams.ToKey("page_read"), chapter, 30*24*time.Hour)
		go MangaChapterSync(context.Background(), queryParams, chapter)
	}

	return chapter, models.Meta{}, nil
}
