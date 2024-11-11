// Package ae_test provides tests for the UAE TRN (Tax Registration Number) validation.
package ae_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/ae"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "good 1", code: "123456789012345"},
		{name: "good 2", code: "187654321098765"},
		{name: "good 3", code: "100111222333444"},

		// Invalid formats
		{name: "too short", code: "12345678901234", err: "invalid format: TRN must be a 15-digit number"},
		{name: "too long", code: "1234567890123456", err: "invalid format: TRN must be a 15-digit number"},
		{name: "non-numeric", code: "12345678ABCD345", err: "invalid format: TRN must be a 15-digit number"},
		{name: "not normalized", code: "1234-5678-9012-345", err: "invalid format: TRN must be a 15-digit number"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "AE", Code: tt.code}
			err := ae.Validate(tID)
			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), tt.err)
				}
			}
		})
	}
}
