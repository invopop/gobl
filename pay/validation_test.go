package pay_test

import (
	"testing"

	"github.com/invopop/gobl/pay"
	"github.com/stretchr/testify/assert"
)

func TestIsIBAN(t *testing.T) {
	tests := []struct {
		name string
		iban string
		want bool
	}{
		{"Finnish IBAN", "FI2112345600000785", true},
		{"German IBAN", "DE89370400440532013000", true},
		{"Norwegian IBAN at the minimum length", "NO9386011117947", true},
		{"check digits do not match", "FI2112345600000786", false},
		{"check digits match but too short for ISO 13616", "FI2100021", false},
		{"lower case", "fi2112345600000785", false},
		{"grouping spaces", "FI21 1234 5600 0007 85", false},
		{"not an account number", "NOT-AN-IBAN", false},
		{"empty", "", false},
		{"check digits only", "FI21", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, pay.IsIBAN.Check(tt.iban))
		})
	}

	t.Run("not a string", func(t *testing.T) {
		assert.False(t, pay.IsIBAN.Check(42))
	})

	t.Run("description lists both conditions", func(t *testing.T) {
		assert.Equal(t,
			`matches ^[A-Z]{2}\d{2}[A-Z0-9]{11,30}$, and has matching check digits`,
			pay.IsIBAN.String())
	})
}

func TestIsBIC(t *testing.T) {
	tests := []struct {
		name string
		bic  string
		want bool
	}{
		{"eight characters", "NDEAFIHH", true},
		{"eleven characters with branch code", "DEUTDEFF500", true},
		{"too short", "NDEA", false},
		{"nine characters", "NDEAFIHHX", false},
		{"lower case", "ndeafihh", false},
		{"digit in the institution code", "N1EAFIHH", false},
		// The pattern tests in rules/is skip absent values so that presence
		// is asserted separately.
		{"empty", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, pay.IsBIC.Check(tt.bic))
		})
	}
}
