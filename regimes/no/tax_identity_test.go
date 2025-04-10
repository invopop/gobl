// Package no_test provides tests for the Norwegian TRN (Tax Registration Number) validation.
package no_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/no"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "valid TRN 1", code: "123456785"}, // Valid checksum
		{name: "valid TRN 2", code: "290883970"}, // Corrected valid TRN (check digit computed as 0)
		{name: "valid TRN 3", code: "974760673"}, // Valid checksum

		// Invalid formats
		{name: "too short", code: "12345678", err: "must be a 9-digit number"},
		{name: "too long", code: "1234567890", err: "must be a 9-digit number"},
		{name: "non-numeric", code: "12345ABCD", err: "must be a 9-digit number"},
		{name: "invalid checksum", code: "123456789", err: "invalid checksum for TRN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "NO", Code: tt.code}
			err := no.Validate(tID)
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
