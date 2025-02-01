package manga_scrapper_service

import (
	"context"
	"fmt"
	"slices"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/contract"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository/manga_history_repository"
	"github.com/umarkotak/animapu-api/internal/repository/manga_library_repository"
)

func MultiInjectLibraryAndHistory(ctx context.Context, user models.User, mangas []contract.Manga) []contract.Manga {
	if len(mangas) == 0 {
		return mangas
	}

	sources := []string{}
	sourceIDs := []string{}
	for _, manga := range mangas {
		sourceIDs = append(sourceIDs, manga.SourceID)
		sources = append(sources, manga.Source)
	}

	librarySourceIDs, err := manga_library_repository.GetByUserAndSourceDetail(ctx, user.ID, sources, sourceIDs)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
	}

	mangaHistories, err := manga_history_repository.GetByUserAndSourceDetail(ctx, user.ID, sources, sourceIDs)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
	}

	mangaHistoryMap := map[string]models.MangaHistoryDetailed{}
	for _, mangaHistory := range mangaHistories {
		key := fmt.Sprintf("%s::%s", mangaHistory.MangaSource, mangaHistory.MangaSourceID)
		mangaHistoryMap[key] = mangaHistory
	}

	for idx, manga := range mangas {
		if slices.Contains(librarySourceIDs, manga.SourceID) {
			mangas[idx].IsInLibrary = true
		}

		key := fmt.Sprintf("%s::%s", manga.Source, manga.SourceID)
		mangaHistory, found := mangaHistoryMap[key]
		if found {
			mangas[idx].LastChapterRead = mangaHistory.ChapterNumber
			mangas[idx].LastLink = mangaHistory.FrontendPath
			mangas[idx].LastReadAt = mangaHistory.UpdatedAt
		}
	}

	return mangas
}

func MultiInjectLibraryAndHistoryForLibrary(ctx context.Context, user models.User, mangas []contract.MangaLibrary) []contract.MangaLibrary {
	if len(mangas) == 0 {
		return mangas
	}

	sources := []string{}
	sourceIDs := []string{}
	for _, manga := range mangas {
		sourceIDs = append(sourceIDs, manga.SourceID)
		sources = append(sources, manga.Source)
	}

	librarySourceIDs, err := manga_library_repository.GetByUserAndSourceDetail(ctx, user.ID, sources, sourceIDs)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
	}

	mangaHistories, err := manga_history_repository.GetByUserAndSourceDetail(ctx, user.ID, sources, sourceIDs)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
	}

	mangaHistoryMap := map[string]models.MangaHistoryDetailed{}
	for _, mangaHistory := range mangaHistories {
		key := fmt.Sprintf("%s::%s", mangaHistory.MangaSource, mangaHistory.MangaSourceID)
		mangaHistoryMap[key] = mangaHistory
	}

	for idx, manga := range mangas {
		if slices.Contains(librarySourceIDs, manga.SourceID) {
			mangas[idx].IsInLibrary = true
		}

		key := fmt.Sprintf("%s::%s", manga.Source, manga.SourceID)
		mangaHistory, found := mangaHistoryMap[key]
		if found {
			mangas[idx].LastChapterRead = mangaHistory.ChapterNumber
			mangas[idx].LastLink = mangaHistory.FrontendPath
			mangas[idx].LastReadAt = mangaHistory.UpdatedAt
		}
	}

	return mangas
}

func MultiInjectLibraryAndHistoryForHistory(ctx context.Context, user models.User, mangas []contract.MangaHistory) []contract.MangaHistory {
	if len(mangas) == 0 {
		return mangas
	}

	sources := []string{}
	sourceIDs := []string{}
	for _, manga := range mangas {
		sourceIDs = append(sourceIDs, manga.SourceID)
		sources = append(sources, manga.Source)
	}

	librarySourceIDs, err := manga_library_repository.GetByUserAndSourceDetail(ctx, user.ID, sources, sourceIDs)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
	}

	mangaHistories, err := manga_history_repository.GetByUserAndSourceDetail(ctx, user.ID, sources, sourceIDs)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
	}

	mangaHistoryMap := map[string]models.MangaHistoryDetailed{}
	for _, mangaHistory := range mangaHistories {
		key := fmt.Sprintf("%s::%s", mangaHistory.MangaSource, mangaHistory.MangaSourceID)
		mangaHistoryMap[key] = mangaHistory
	}

	for idx, manga := range mangas {
		if slices.Contains(librarySourceIDs, manga.SourceID) {
			mangas[idx].IsInLibrary = true
		}

		key := fmt.Sprintf("%s::%s", manga.Source, manga.SourceID)
		mangaHistory, found := mangaHistoryMap[key]
		if found {
			mangas[idx].LastChapterRead = mangaHistory.ChapterNumber
			mangas[idx].LastLink = mangaHistory.FrontendPath
			mangas[idx].LastReadAt = mangaHistory.UpdatedAt
		}
	}

	return mangas
}
