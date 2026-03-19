package at_test

import (
	"testing"

	_ "github.com/invopop/gobl/regimes/at"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	_ "github.com/invopop/gobl/regimes/at"
)

func TestTaxIdentityRules(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "good 1", code: "U00000033"},
		{name: "good 2", code: "U38516405"},
		{name: "good 3", code: "U64727905"},
		{
			name: "zeros",
			code: "U00000000",
			err:  "IDENTITY-01",
		},
		{
			name: "bad mid length",
			code: "U12345678910",
			err:  "IDENTITY-01",
		},
		{
			name: "too long",
			code: "U1234567890123",
			err:  "IDENTITY-01",
		},
		{
			name: "too short",
			code: "U123456",
			err:  "IDENTITY-01",
		},
		{
			name: "not normalized",
			code: "U-385.16.405",
			err:  "IDENTITY-01",
		},
		{
			name: "bad checksum",
			code: "U10223001",
			err:  "IDENTITY-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "AT", Code: tt.code}
			err := rules.Validate(tID)
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
