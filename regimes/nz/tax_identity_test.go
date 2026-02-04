package nz_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/nz"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidIRDNumbers(t *testing.T) {
	tests := []string{
		"49091850",
		"49-091-850",
		"136410132",
		"136-410-132",
	}

	r := nz.New()
	for _, ird := range tests {
		t.Run(ird, func(t *testing.T) {
			tID := &tax.Identity{
				Country: "NZ",
				Code:    cbc.Code(ird),
			}
			r.Normalizer(tID)
			err := r.Validator(tID)
			assert.NoError(t, err, "IRD %s should be valid", ird)
		})
	}
}

func TestInvalidIRDNumbers(t *testing.T) {
	tests := []struct {
		name string
		ird  string
	}{
		{"TooShort", "1234567"},
		{"TooLong", "1234567890"},
		{"OutOfRangeLow", "09999999"},
		{"OutOfRangeHigh", "150000001"},
		{"BadCheckDigit", "49091851"},
		{"AllZeros", "00000000"},
		{"Letters", "4909185A"},
	}

	r := nz.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{
				Country: "NZ",
				Code:    cbc.Code(tt.ird),
			}
			r.Normalizer(tID)
			err := r.Validator(tID)
			assert.Error(t, err, "IRD %s should be invalid", tt.ird)
		})
	}
}

func TestNormalizeTaxIdentity(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"IRDWithDashes", "49-091-850", "49091850"},
		{"IRDWithSpaces", "49 091 850", "49091850"},
		{"IRDClean", "49091850", "49091850"},
	}

	r := nz.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{
				Country: "NZ",
				Code:    cbc.Code(tt.input),
			}
			r.Normalizer(tID)
			assert.Equal(t, tt.expected, tID.Code.String())
		})
	}
}

func TestIdentityKeys(t *testing.T) {
	r := nz.New()
	require.NotNil(t, r.Identities)

	var foundIRD, foundGLN bool
	for _, id := range r.Identities {
		if id.Key == nz.IdentityKeyIRD {
			foundIRD = true
		}
		if id.Key == "gln" {
			foundGLN = true
		}
	}

	assert.True(t, foundIRD, "Should have nz-ird identity")
	assert.True(t, foundGLN, "Should have gln identity")
}
