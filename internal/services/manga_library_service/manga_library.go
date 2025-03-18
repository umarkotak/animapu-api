package manga_library_service

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/contract"
	"github.com/umarkotak/animapu-api/internal/repository/manga_library_repository"
)

func GetLibraries(ctx context.Context, params contract.MangaLibraryParams) ([]contract.MangaLibrary, error) {
	mangaLibraries, err := manga_library_repository.GetByUserIDDetailed(ctx, params)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return []contract.MangaLibrary{}, err
	}

	resMangaLibraries := []contract.MangaLibrary{}
	for _, mangaLibrary := range mangaLibraries {
		coverImage := contract.CoverImage{
			Index:     0,
			ImageUrls: mangaLibrary.MangaCoverUrls,
		}

		resMangaLibrary := contract.MangaLibrary{
			Source:              mangaLibrary.MangaSource,
			SourceID:            mangaLibrary.MangaSourceID,
			Title:               mangaLibrary.MangaTitle,
			LatestChapterNumber: mangaLibrary.MangaLatestChapter,
			CoverImages:         []contract.CoverImage{coverImage},
			LastChapterRead:     mangaLibrary.ChapterNumber,
			LastLink:            mangaLibrary.FrontendPath,
			UpdatedAt:           mangaLibrary.MangaUpdatedAt,
		}

		resMangaLibraries = append(resMangaLibraries, resMangaLibrary)
	}

	return resMangaLibraries, nil
}
