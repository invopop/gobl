package nz_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/nz"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidIRDNumbers(t *testing.T) {
	tests := []string{
		"49091850", // Spec ex.1: primary weights, remainder 0, check digit 0
		"49-091-850",
		"35901981",  // Spec ex.2: primary weights, remainder 10, 11-10=1
		"49098576",  // Spec ex.3: primary gives 10, secondary gives 6
		"136410132", // Spec ex.4: 9-digit, secondary weights
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
		name   string
		ird    string
		errMsg string
	}{
		{"TooShort", "1234567", "must be 8-9 digit"},
		{"TooLong", "1234567890", "must be 8-9 digit"},
		{"OutOfRangeLow", "09999999", "out of valid range"},
		{"OutOfRangeHigh", "200000001", "out of valid range"},
		{"BadCheckDigit", "49091851", "check digit mismatch"},
		{"BadCheckDigitSecondary", "136410133", "check digit mismatch"}, // Spec ex.5: secondary check digit mismatch
		{"BelowRange", "9125568", "must be 8-9 digit"},                  // Spec ex.6: below valid range
		{"AllZeros", "00000000", "out of valid range"},
		{"Letters", "4909185A", "must be 8-9 digit"},
		// Both primary and secondary weights give checkDigit == 10
		{"BothWeightsGive10", "10000050", "check digit calculation failed"},
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
			assert.Contains(t, err.Error(), tt.errMsg)
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

func TestNZBNAsTaxIdentityFails(t *testing.T) {
	r := nz.New()
	tID := &tax.Identity{
		Country: "NZ",
		Code:    "9429041234563",
	}
	r.Normalizer(tID)
	err := r.Validator(tID)
	assert.Error(t, err, "NZBN should not be accepted as a tax identity")
}

func TestIdentityKeys(t *testing.T) {
	r := nz.New()
	require.NotNil(t, r.Identities)

	var foundIRD, foundNZBN bool
	for _, id := range r.Identities {
		if id.Key == nz.IdentityKeyIRD {
			foundIRD = true
		}
		if id.Key == org.IdentityKeyGLN {
			foundNZBN = true
		}
	}

	assert.True(t, foundIRD, "Should have IRD identity")
	assert.True(t, foundNZBN, "Should have NZBN identity")
}

func TestNormalizeTaxIdentityNil(_ *testing.T) {
	r := nz.New()
	var tID *tax.Identity
	r.Normalizer(tID)
}

func TestNormalizeTaxIdentityEmptyCode(t *testing.T) {
	r := nz.New()
	tID := &tax.Identity{Country: "NZ", Code: ""}
	r.Normalizer(tID)
	assert.Equal(t, "", tID.Code.String())
}

func TestValidateTaxIdentityEmptyCode(t *testing.T) {
	r := nz.New()
	tID := &tax.Identity{Country: "NZ", Code: ""}
	err := r.Validator(tID)
	assert.NoError(t, err, "Empty code should pass validation")
}
