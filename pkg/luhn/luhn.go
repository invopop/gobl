// Package luhn provides utilities for working with the Luhn algorithm for
// validating identification numbers.
package luhn

import (
	"regexp"

	"github.com/invopop/gobl/cbc"
)

var nonNumberRegExp = regexp.MustCompile(`[^\d]+`)

// Check if the code containing digits passes the Luhn algorithm validation.
//
// `code` must contain only numeric characters.
//
// [Source]
//
// [Source]: https://github.com/luhnmod10/go
func Check(code cbc.Code) bool {
	if code == cbc.CodeEmpty || nonNumberRegExp.MatchString(string(code)) {
		return false
	}

	var checksum int

	numberLen := len(code)
	for i := numberLen - 1; i >= 0; i -= 2 {
		n := code[i] - '0'
		checksum += int(n)
	}
	for i := numberLen - 2; i >= 0; i -= 2 {
		n := code[i] - '0'
		n *= 2
		if n > 9 {
			n -= 9
		}
		checksum += int(n)
	}

	return checksum%10 == 0
}
