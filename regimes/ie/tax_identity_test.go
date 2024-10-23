package ie_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/ie"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "valid ind old-style", code: "1234567T"},
		{name: "valid ind new-style", code: "1234567TW"},
		{name: "valid company", code: "1A23456T"},
		{
			name: "too many digits",
			code: "123456789",
			err:  "invalid format",
		},
		{
			name: "too few digits",
			code: "12345T",
			err:  "invalid format",
		},
		{
			name: "no digits",
			code: "ABCDEFGH",
			err:  "invalid format",
		},
		{
			name: "lower case",
			code: "1234567t",
			err:  "invalid format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "IE", Code: tt.code}
			err := ie.Validate(tID)
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
