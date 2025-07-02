//nolint:revive // "common" package name is acceptable here for shared regime utilities
package common

import "strconv"

// ComputeLuhnCheckDigit expects as argument a number string excluding the check
// digit. The returned integer should be checked against the check digit by the
// caller.
// Luhn Algorithm definition: https://en.wikipedia.org/wiki/Luhn_algorithm
func ComputeLuhnCheckDigit(number string) string {
	sum := 0
	pos := 0

	for i := len(number) - 1; i >= 0; i-- {
		digit := int(number[i] - '0')

		if pos%2 == 0 {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		pos++
	}

	return strconv.FormatInt(int64((10-(sum%10))%10), 10)
}
