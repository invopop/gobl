package nz_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/nz"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "good 8 digits", code: "12345678"},
		{name: "good 9 digits", code: "123456789"},

		{name: "too short", code: "1234567", err: "must be 8 or 9 digits"},
		{name: "too long", code: "1234567890", err: "must be 8 or 9 digits"},
		{name: "non-numeric", code: "12A45678", err: "must be 8 or 9 digits"},
		{name: "not normalized", code: "12-345-678", err: "must be 8 or 9 digits"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "NZ", Code: tt.code}
			err := nz.Validate(tID)
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

func TestValidateTaxIdentityEmpty(t *testing.T) {
	// empty identity
	tID := &tax.Identity{}
	err := nz.Validate(tID)
	assert.NoError(t, err)
}
