package anime_scrapper_service

import (
	"context"
	"fmt"
	"slices"

	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/contract"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository/anime_history_repository"
	"github.com/umarkotak/animapu-api/internal/repository/manga_history_repository"
	"github.com/umarkotak/animapu-api/internal/repository/manga_library_repository"
)

func MultiInjectLibraryAndHistoryForLatest(ctx context.Context, user models.User, animes []contract.Anime) []contract.Anime {
	if len(animes) == 0 {
		return animes
	}

	sources := []string{}
	sourceIDs := []string{}
	for _, anime := range animes {
		sourceIDs = append(sourceIDs, anime.ID)
		sources = append(sources, anime.Source)
	}

	// librarySourceIDs, err := manga_library_repository.GetByUserAndSourceDetail(ctx, user.ID, sources, sourceIDs)
	// if err != nil {
	// 	logrus.WithContext(ctx).Error(err)
	// }

	animeHistories, err := anime_history_repository.GetByUserAndSourceDetail(ctx, user.ID, sources, sourceIDs)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
	}

	animeHistoryMap := map[string]models.AnimeHistoryDetailed{}
	for _, animeHistory := range animeHistories {
		key := fmt.Sprintf("%s::%s", animeHistory.AnimeSource, animeHistory.AnimeSourceID)
		animeHistoryMap[key] = animeHistory
	}

	for idx, anime := range animes {
		// if slices.Contains(librarySourceIDs, anime.ID) {
		// 	animes[idx].IsInLibrary = true
		// }

		key := fmt.Sprintf("%s::%s", anime.Source, anime.ID)
		animeHistory, found := animeHistoryMap[key]
		if found {
			animes[idx].LastEpisodeWatch = animeHistory.EpisodeNumber
			animes[idx].LastLink = animeHistory.FrontendPath
			animes[idx].LastWatchAt = animeHistory.UpdatedAt
		}
	}

	return animes
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

func MultiInjectLibraryAndHistoryForHistory(ctx context.Context, user models.User, animes []contract.AnimeHistory) []contract.AnimeHistory {
	if len(animes) == 0 {
		return animes
	}

	sources := []string{}
	sourceIDs := []string{}
	for _, anime := range animes {
		sourceIDs = append(sourceIDs, anime.ID)
		sources = append(sources, anime.Source)
	}

	// librarySourceIDs, err := manga_library_repository.GetByUserAndSourceDetail(ctx, user.ID, sources, sourceIDs)
	// if err != nil {
	// 	logrus.WithContext(ctx).Error(err)
	// }

	animeHistories, err := anime_history_repository.GetByUserAndSourceDetail(ctx, user.ID, sources, sourceIDs)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
	}

	animeHistoryMap := map[string]models.AnimeHistoryDetailed{}
	for _, animeHistory := range animeHistories {
		key := fmt.Sprintf("%s::%s", animeHistory.AnimeSource, animeHistory.AnimeSourceID)
		animeHistoryMap[key] = animeHistory
	}

	for idx, anime := range animes {
		// if slices.Contains(librarySourceIDs, anime.ID) {
		// 	animes[idx].IsInLibrary = true
		// }

		key := fmt.Sprintf("%s::%s", anime.Source, anime.ID)
		animeHistory, found := animeHistoryMap[key]
		if found {
			animes[idx].LastEpisodeWatch = animeHistory.EpisodeNumber
			animes[idx].LastLink = animeHistory.FrontendPath
			animes[idx].LastWatchAt = animeHistory.UpdatedAt
		}
	}

	return animes
}
