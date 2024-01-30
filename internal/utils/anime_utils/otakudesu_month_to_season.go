package anime_utils

import "strings"

func OtakudesuMonthToSeason(month string) string {
	m := strings.ToLower(month)
	switch m {
	case "jan":
		return "winter"
	case "feb":
		return "winter"
	case "mar":
		return "winter"
	case "apr":
		return "spring"
	case "may":
		return "spring"
	case "jun":
		return "spring"
	case "jul":
		return "summer"
	case "agu":
		return "summer"
	case "sep":
		return "summer"
	case "okt":
		return "fall"
	case "nov":
		return "fall"
	case "des":
		return "fall"
	}
	return ""
}
