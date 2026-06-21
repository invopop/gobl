package ee_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"

	_ "github.com/invopop/gobl/regimes/ee"
)

func TestTaxIdentityRules(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		// Valid KMKR numbers (format + checksum verified against EMTA/VIES)
		{name: "valid 1", code: "100207415"},
		{name: "valid 2", code: "100931558"},
		{name: "valid 3", code: "100594102"},
		// Format errors
		{
			name: "too short",
			code: "10020741",
			err:  "IDENTITY-01",
		},
		{
			name: "too long",
			code: "1002074159",
			err:  "IDENTITY-01",
		},
		{
			name: "not normalized, EE prefix present",
			code: "EE100207415",
			err:  "IDENTITY-01",
		},
		{
			name: "contains non-digits",
			code: "10020741X",
			err:  "IDENTITY-01",
		},
		// Checksum errors
		{
			name: "bad checksum, last digit off by one",
			code: "100207410",
			err:  "IDENTITY-01",
		},
		{
			name: "bad checksum, sequential digits",
			code: "123456789",
			err:  "IDENTITY-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "EE", Code: tt.code}
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
