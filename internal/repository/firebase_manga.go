package repository

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
)

func FbGetHomeByMangaSource(ctx context.Context, mangaSource string) (models.FbMangaHomeCache, error) {
	mangaSourceRef := animapuLiteApiRef.Child(mangaSource)
	if mangaSourceRef == nil {
		return models.FbMangaHomeCache{}, models.ErrCacheNotFound
	}

	homeRef := mangaSourceRef.Child("home")

	fbMangaHome := models.FbMangaHomeCache{}
	err := homeRef.Get(ctx, &fbMangaHome)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return fbMangaHome, err
	}

	return fbMangaHome, nil
}

func FbUpsertHomeByMangaSource(ctx context.Context, mangaSource string, mangas []models.Manga) ([]models.Manga, error) {
	mangaSourceRef := animapuLiteApiRef.Child(mangaSource)
	if mangaSourceRef == nil {
		return []models.Manga{}, models.ErrCacheNotFound
	}

	homeRef := mangaSourceRef.Child("home")

	now := time.Now().UTC()
	err := homeRef.Set(ctx, models.FbMangaHomeCache{
		UpdatedAt:     now,
		UpdatedAtUnix: now.Unix(),
		ExpiredAt:     now.Add(4 * time.Hour), // 4 hours
		Mangas:        mangas,
	})
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return mangas, err
	}

	return mangas, nil
}

func FbGetMangaDetailByMangaSource(ctx context.Context, mangaSource string, manga models.Manga) (models.FbMangaDetailCache, error) {
	var err error

	mangaSourceRef := animapuLiteApiRef.Child(mangaSource)
	if mangaSourceRef == nil {
		return models.FbMangaDetailCache{}, models.ErrCacheNotFound
	}

	fbMangaDetail := models.FbMangaDetailCache{}
	detailRef := mangaSourceRef.Child("detail")

	if detailRef == nil {
		err = detailRef.Set(ctx, map[string]interface{}{})
		if err != nil {
			logrus.WithContext(ctx).Error(err)
			return fbMangaDetail, err
		}
	}

	fbMangaDetailSelected := detailRef.Child(manga.SourceID)
	if fbMangaDetailSelected == nil {
		err = models.ErrCacheNotFound
		return fbMangaDetail, err
	}

	err = fbMangaDetailSelected.Get(ctx, &fbMangaDetail)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return fbMangaDetail, err
	}

	return fbMangaDetail, nil
}

func FbUpsertMangaDetailByMangaSource(ctx context.Context, mangaSource string, manga models.Manga) (models.Manga, error) {
	var err error

	mangaSourceRef := animapuLiteApiRef.Child(mangaSource)
	if mangaSourceRef == nil {
		return models.Manga{}, models.ErrCacheNotFound
	}

	detailRef := mangaSourceRef.Child("detail")

	if detailRef == nil {
		err = detailRef.Set(ctx, map[string]interface{}{})
		if err != nil {
			logrus.WithContext(ctx).Error(err)
			return manga, err
		}
	}

	fbMangaDetailSelected := detailRef.Child(manga.SourceID)

	now := time.Now().UTC()
	err = fbMangaDetailSelected.Set(ctx, models.FbMangaDetailCache{
		UpdatedAt:     now,
		UpdatedAtUnix: now.Unix(),
		ExpiredAt:     now.Add(24 * 3 * time.Hour), // 3 days
		Manga:         manga,
	})
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return manga, err
	}

	return manga, nil
}

func FbUpvoteManga(ctx context.Context, manga models.Manga) (models.Manga, error) {
	onePopularRef := popularMangaRef.Child(manga.GetUniqueKey())

	var err error
	if onePopularRef == nil {
		err = onePopularRef.Set(ctx, manga)
	} else {
		err = onePopularRef.Get(ctx, &manga)
		if err != nil {
			logrus.WithContext(ctx).Error(err)
			return manga, err
		}
		manga.PopularityPoint += 1
		err = onePopularRef.Set(ctx, manga)
	}
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return manga, err
	}

	return manga, nil
}

func FbGetPopularManga(ctx context.Context) ([]models.Manga, error) {
	mangaMap := map[string]models.Manga{}
	mangas := []models.Manga{}

	err := popularMangaRef.Get(ctx, &mangaMap)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return mangas, err
	}

	for _, oneManga := range mangaMap {
		mangas = append(mangas, oneManga)
	}

	// TODO: sort by popularity point

	return mangas, nil
}
