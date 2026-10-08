// Package il_test provides tests for the Israeli Osek Murshe tax ID validation.
package il_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/il"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		// Valid: 9 numeric digits
		{name: "valid 1", code: "123456789"},
		{name: "valid 2", code: "000000001"},
		{name: "valid 3", code: "516179157"}, // real-world example

		// Invalid: wrong length
		{name: "too short", code: "12345678", err: "must be a 9-digit number"},
		{name: "too long", code: "1234567890", err: "must be a 9-digit number"},

		// Invalid: non-numeric
		{name: "contains letters", code: "12345678A", err: "must be a 9-digit number"},
		{name: "contains hyphens", code: "123-456-7", err: "must be a 9-digit number"},
		{name: "contains spaces", code: "123 456 7", err: "must be a 9-digit number"},

		// Empty: skip validation
		{name: "empty code", code: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "IL", Code: tt.code}
			err := il.Validate(tID)
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
