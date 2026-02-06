// Package gs1 provides utilities for validating GS1 identification numbers
// such as the Global Location Number (GLN).
package gs1

import (
	"regexp"

	"github.com/invopop/gobl/cbc"
)

var glnPattern = regexp.MustCompile(`^\d{13}$`)

// CheckGLN validates a GS1 Global Location Number (GLN).
// Returns true if the code is exactly 13 digits with a valid GS1 Modulo-10
// check digit.
func CheckGLN(code cbc.Code) bool {
	val := code.String()
	if !glnPattern.MatchString(val) {
		return false
	}

	sum := 0
	for i := 0; i < 12; i++ {
		d := int(val[i] - '0')
		if i%2 == 0 {
			sum += d
		} else {
			sum += d * 3
		}
	}

	check := (10 - sum%10) % 10
	return check == int(val[12]-'0')
}

// HasPrefix checks whether the code starts with the given GS1 prefix.
// Country prefixes are assigned by GS1 member organizations (e.g. "94" for
// New Zealand).
func HasPrefix(code cbc.Code, prefix string) bool {
	val := code.String()
	if len(val) < len(prefix) {
		return false
	}
	return val[:len(prefix)] == prefix
}
