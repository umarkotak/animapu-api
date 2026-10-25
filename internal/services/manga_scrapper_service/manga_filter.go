package manga_scrapper_service

import (
	"context"
	"slices"

	"github.com/umarkotak/animapu-api/internal/contract"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository/manga_repository"
)

func FilterNSFW(ctx context.Context, source string, mangas []contract.Manga) ([]contract.Manga, error) {
	if len(mangas) == 0 {
		return mangas, nil
	}

	sourceIDs := make([]string, 0, len(mangas))
	for _, manga := range mangas {
		sourceIDs = append(sourceIDs, manga.SourceID)
	}
	dbMangas, err := manga_repository.GetBySourceAndSourceIDs(ctx, source, sourceIDs)
	if err != nil {
		return nil, err
	}

	return excludeNSFW(mangas, dbMangas), nil
}

func excludeNSFW(mangas []contract.Manga, dbMangas []models.Manga) []contract.Manga {
	type mangaKey struct{ source, sourceID string }
	nsfw := make(map[mangaKey]bool)
	for _, manga := range dbMangas {
		if slices.Contains(manga.Tags, "nsfw:true") {
			nsfw[mangaKey{manga.Source, manga.SourceID}] = true
		}
	}

	filtered := make([]contract.Manga, 0, len(mangas))
	for _, manga := range mangas {
		if !nsfw[mangaKey{manga.Source, manga.SourceID}] {
			filtered = append(filtered, manga)
		}
	}
	return filtered
}
