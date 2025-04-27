package cbc

import (
	"regexp"
)

var nonNumberRegExp = regexp.MustCompile(`[^\d]+`)

// ValidateLuhn checks if a string of digits passes the Luhn algorithm validation.
//
// `number` must contain only numeric characters.
//
// [Source]
//
// [Source]: https://github.com/luhnmod10/go
func ValidateLuhn(number string) bool {
	if number == "" || nonNumberRegExp.MatchString(number) {
		return false
	}

	var checksum int

	numberLen := len(number)
	for i := numberLen - 1; i >= 0; i -= 2 {
		n := number[i] - '0'
		checksum += int(n)
	}
	for i := numberLen - 2; i >= 0; i -= 2 {
		n := number[i] - '0'
		n *= 2
		if n > 9 {
			n -= 9
		}
		checksum += int(n)
	}

	return checksum%10 == 0
}
