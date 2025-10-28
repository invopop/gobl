package luhn_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pkg/luhn"
	"github.com/stretchr/testify/assert"
)

func TestCheck(t *testing.T) {
	t.Parallel()
	validNumbers := []cbc.Code{
		"0",  // Single digit 0 is valid
		"18", // 1×2 = 2, 2+8 = 10, 10 mod 10 = 0
		"26", // 2×2 = 4, 4+6 = 10, 10 mod 10 = 0
		"34", // 3×2 = 6, 6+4 = 10, 10 mod 10 = 0
		"42", // 4×2 = A, 8+2 = 10, 10 mod 10 = 0
		"59", // 5×2 = 10, 1+0+9 = 10, 10 mod 10 = 0
		"67", // 6×2 = 12, 1+2+7 = 10, 10 mod 10 = 0
		"75", // 7×2 = 14, 1+4+5 = 10, 10 mod 10 = 0
		"83", // 8×2 = 16, 1+6+3 = 10, 10 mod 10 = 0
		"91", // 9×2 = 18, 1+8+1 = 10, 10 mod 10 = 0
		"79927398713",
		"4111111111111111", // Valid credit card number
	}

	invalidNumbers := []cbc.Code{
		"1",                // 1 is not valid
		"19",               // 1×2 = 2, 2+9 = 11, 11 mod 10 = 1
		"123",              // 1×2 = 2, 2+2+3 = 7, 7 mod 10 = 7
		"79927398714",      // Invalid by one
		"4111111111111112", // Invalid credit card
		"",                 // Empty string
		"123a",             // Non-numeric character
	}

	for _, number := range validNumbers {
		t.Run("Valid: "+number.String(), func(t *testing.T) {
			t.Parallel()
			result := luhn.Check(number)
			assert.True(t, result, "Should be valid: %s", number)
		})
	}

	for _, number := range invalidNumbers {
		t.Run("Invalid: "+number.String(), func(t *testing.T) {
			t.Parallel()
			result := luhn.Check(number)
			assert.False(t, result, "Should be invalid: %s", number)
		})
	}
}
