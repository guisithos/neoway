package util

import (
	"bytes"
	"unicode"
)

func toInt(r rune) int {
	return int(r - '0')
}

// cleanNonDigits removes every rune that is not a digit.
func cleanNonDigits(doc *string) {
	buf := bytes.NewBufferString("")
	for _, r := range *doc {
		if unicode.IsDigit(r) {
			buf.WriteRune(r)
		}
	}

	*doc = buf.String()
}
