package utils

import (
	"strings"
	"unicode"
)

// Slugify converts a string into a URL-friendly slug.
func Slugify(s string) string {

	var slug strings.Builder

	for _, r := range s {

		if unicode.IsLetter(r) || unicode.IsDigit(r) {

			slug.WriteRune(unicode.ToLower(r))

		} else if unicode.IsSpace(r) {

			slug.WriteRune('-')

		}

	}

	return slug.String()

}
