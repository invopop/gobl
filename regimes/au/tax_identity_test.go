package au_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/au"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "ABN good 1", code: "53004085616"},
		{name: "ABN good 2", code: "12004044937"},
		{name: "ABN good 3", code: "84085334037"},
		{
			name: "zeros",
			code: "00000000000",
			err:  "zeros",
		},
		{
			name: "too long",
			code: "123456789012",
			err:  "invalid format",
		},
		{
			name: "too short",
			code: "12345678",
			err:  "invalid format",
		},
		{
			name: "not normalized",
			code: "12.449.95-4",
			err:  "invalid format",
		},
		{
			name: "letter",
			code: "12a44939544",
			err:  "invalid format",
		},
		{
			name: "bad ABN checksum 1",
			code: "99999999123",
			err:  "checksum mismatch",
		},
		{
			name: "bad ABN checksum 2",
			code: "73827573823",
			err:  "checksum mismatch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "AU", Code: tt.code}
			err := au.Validate(tID)
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
