package my_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/my"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeTaxIdentity(t *testing.T) {
	r := my.New()

	tests := []struct {
		Code     cbc.Code
		Expected cbc.Code
	}{
		{
			Code:     "sSt1234567890",
			Expected: "SST1234567890",
		},
		{
			Code:     " w10-12345678-123 ",
			Expected: "W10-12345678-123",
		},
		{
			Code:     "201901234567",
			Expected: "201901234567",
		},
	}
	for _, ts := range tests {
		tID := &tax.Identity{Country: "MY", Code: ts.Code}
		r.NormalizeObject(tID)
		assert.Equal(t, ts.Expected, tID.Code)
	}
}

func TestValidateTaxIdentity(t *testing.T) {
	r := my.New()

	tests := []struct {
		Code  cbc.Code
		Valid bool
	}{
		// Valid 12-digit business numbers
		{Code: "201901234567", Valid: true},
		{Code: "123456789012", Valid: true},

		// Valid SST numbers
		{Code: "SST1234567890", Valid: true},
		{Code: "W10-12345678-123", Valid: true},

		// Invalid ones
		{Code: "2019-123456", Valid: false},        // Wrong format
		{Code: "XYZ", Valid: false},                // Too short
		{Code: "SST-INVALID-NUMBER", Valid: false}, // Wrong pattern
	}

	for _, ts := range tests {
		t.Run(string(ts.Code), func(t *testing.T) {
			tID := &tax.Identity{Country: "MY", Code: ts.Code}
			err := r.ValidateObject(tID)
			if ts.Valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
