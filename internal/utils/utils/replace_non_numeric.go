package utils

import "regexp"

var (
	reg, _ = regexp.Compile("[^0-9.]+")
)

func RemoveNonNumeric(str string) string {
	return reg.ReplaceAllString(str, "")
}
