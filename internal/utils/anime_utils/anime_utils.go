package anime_utils

import "strings"

func SeasonToIndex(seasonName string) int {
	seasonName = strings.ToLower(seasonName)

	if seasonName == "winter" {
		return 1
	}
	if seasonName == "spring" {
		return 2
	}
	if seasonName == "summer" {
		return 3
	}
	if seasonName == "fall" {
		return 4
	}
	return 0
}

func SeasonIndexToSeason(index int) string {
	return ""
}
