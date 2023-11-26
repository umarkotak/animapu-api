package utils

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	reg1, _ = regexp.Compile("[^0-9.]+")
	reg2, _ = regexp.Compile("[^a-zA-Z]+")
	reg3, _ = regexp.Compile("[^a-zA-Z ]+")
)

func RemoveNonNumeric(str string) string {
	return reg1.ReplaceAllString(str, "")
}

func RemoveNonAlphabet(str string) string {
	return reg2.ReplaceAllString(str, "")
}

func RemoveNonAlphabetExceptSpace(str string) string {
	return reg3.ReplaceAllString(str, "")
}

func StringContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func StringMustInt64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func StringMustFloat64(s string) float64 {
	i, _ := strconv.ParseFloat(s, 64)
	return i
}

func ForceSanitizeStringToFloat(s string) float64 {
	s = RemoveNonNumeric(s)
	i, _ := strconv.ParseFloat(s, 64)
	return i
}

func SeasonToIndex(seasonName string) int {
	seasonName = strings.ToLower(seasonName)
	switch seasonName {
	case "winter":
		return 1
	case "spring":
		return 2
	case "summer":
		return 3
	case "fall":
		return 4
	}
	return 0
}
