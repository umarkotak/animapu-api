package manga_scrapper_repository

import (
	"strings"

	"github.com/umarkotak/animapu-api/internal/models"
)

func mangasPaginate(mangas []models.Manga, page, perpage int64) []models.Manga {
	startIndex := (page - 1) * perpage
	endIndex := startIndex + perpage
	if startIndex >= int64(len(mangas)) {
		return []models.Manga{}
	}
	if endIndex >= int64(len(mangas)) {
		endIndex = int64(len(mangas) - 1)
	}
	return mangas[startIndex:endIndex]
}

func convertTitleToMangahubTitle(initialTitle string) string {
	result := strings.ToLower(initialTitle)
	result = strings.Replace(result, "%", "", -1)
	result = strings.Replace(result, "'", "-", -1)
	result = strings.Replace(result, "!", "", -1)
	result = strings.Replace(result, "?", "", -1)
	result = strings.Replace(result, ".", "", -1)
	result = strings.Replace(result, "&", "", -1)
	result = strings.Replace(result, ":", "", -1)
	result = strings.Replace(result, ",", "", -1)
	result = strings.Replace(result, "(", "", -1)
	result = strings.Replace(result, ")", "", -1)
	result = strings.Replace(result, "-", "", -1)
	result = strings.Replace(result, "\"", "", -1)
	result = strings.Replace(result, "  ", "-", -1)
	result = strings.Replace(result, " ", "-", -1)
	return result
}
