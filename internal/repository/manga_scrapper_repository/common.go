package manga_scrapper_repository

import (
	"strings"

	"github.com/umarkotak/animapu-api/internal/contract"
)

func mangasPaginate(mangas []contract.Manga, page, perpage int64) []contract.Manga {
	startIndex := (page - 1) * perpage
	endIndex := startIndex + perpage
	if startIndex >= int64(len(mangas)) {
		return []contract.Manga{}
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
