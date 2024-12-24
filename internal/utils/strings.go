package utils

import "strings"

func ExtractString(stringToExtract string) []string {
	return strings.FieldsFunc(stringToExtract, func(r rune) bool {
		return r == ',' || r == '\n'
	})
}
