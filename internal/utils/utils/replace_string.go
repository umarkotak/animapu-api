package utils

import (
	"regexp"
	"strconv"
)

var (
	reg1, _ = regexp.Compile("[^0-9.]+")
	reg2, _ = regexp.Compile("[^a-zA-Z]+")
)

func RemoveNonNumeric(str string) string {
	return reg1.ReplaceAllString(str, "")
}

func RemoveNonAlphabet(str string) string {
	return reg2.ReplaceAllString(str, "")
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
	i, _ := strconv.ParseFloat(s, 10)
	return i
}

func ForceSanitizeStringToFloat(s string) float64 {
	s = RemoveNonNumeric(s)
	i, _ := strconv.ParseFloat(s, 10)
	return i
}
