package manga_scrapper_service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/contract"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository/manga_chapter_repository"
	"github.com/umarkotak/animapu-api/internal/repository/manga_history_repository"
	"github.com/umarkotak/animapu-api/internal/repository/manga_repository"
)

func MangaSync(ctx context.Context, mangas []contract.Manga) error {
	for _, manga := range mangas {
		existingManga, err := manga_repository.GetBySourceAndSourceID(ctx, manga.Source, manga.SourceID)
		if err != nil && err != sql.ErrNoRows {
			logrus.WithContext(ctx).Error(err)
			continue
		}

		if existingManga.ID != 0 {
			if manga.LatestChapterNumber < existingManga.LatestChapter && manga.ImageURLsEqual(existingManga.CoverUrls) {
				continue
			}

			existingManga.CoverUrls = manga.ImageURLs()
			existingManga.LatestChapter = manga.LatestChapterNumber
			err = manga_repository.Update(ctx, nil, existingManga)

		} else {
			newManga := models.Manga{
				Source:        manga.Source,
				SourceID:      manga.SourceID,
				Title:         manga.Title,
				CoverUrls:     manga.ImageURLs(),
				LatestChapter: manga.LatestChapterNumber,
			}
			newManga.ID, err = manga_repository.Insert(ctx, nil, newManga)
		}
		if err != nil {
			logrus.WithContext(ctx).Error(err)
			continue
		}
	}

	return nil
}

func MangaChapterSync(ctx context.Context, queryParams models.QueryParams, chapter contract.Chapter) error {
	existingManga, err := manga_repository.GetBySourceAndSourceID(ctx, chapter.Source, chapter.SourceID)
	if err != nil && err != sql.ErrNoRows {
		logrus.WithContext(ctx).Error(err)
		return nil
	}

	if existingManga.ID == 0 {
		GetDetail(ctx, models.QueryParams{Source: chapter.Source, SourceID: chapter.SourceID})

		existingManga, err = manga_repository.GetBySourceAndSourceID(ctx, chapter.Source, chapter.SourceID)
		if err != nil {
			logrus.WithContext(ctx).Error(err)
			return err
		}
	}

	mangaChapter := models.MangaChapter{
		MangaID:         existingManga.ID,
		SourceChapterID: chapter.ID,
		ChapterNumber:   chapter.Number,
		ImageUrls:       chapter.ImageURLs(),
	}
	mangaChapter.ID, err = manga_chapter_repository.Insert(ctx, nil, mangaChapter)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return err
	}

	mangaHistory := models.MangaHistory{
		UserID:          queryParams.User.ID,
		MangaID:         mangaChapter.MangaID,
		ChapterNumber:   mangaChapter.ChapterNumber,
		SourceChapterID: mangaChapter.SourceChapterID,
		FrontendPath:    fmt.Sprintf("/mangas/%s/%s/read/%s", existingManga.Source, existingManga.SourceID, mangaChapter.SourceChapterID),
	}
	_, err = manga_history_repository.Insert(ctx, nil, mangaHistory)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return err
	}

	return nil
}
