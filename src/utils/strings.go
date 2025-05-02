package utils

import (
	"unicode"
)

func StringStartsWithUpper(s string) bool {
	if s == "" {
		return false
	}

	for _, r := range s[:1] {
		return unicode.IsUpper(r)
	}

	return false
}
