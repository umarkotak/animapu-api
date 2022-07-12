package utils

import "regexp"

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
