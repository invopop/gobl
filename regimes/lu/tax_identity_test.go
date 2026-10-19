package lu_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	_ "github.com/invopop/gobl/regimes/lu"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestTaxIdentityRules(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "good 1", code: "15027442"},
		{name: "good 2", code: "13669580"},
		{name: "good 3", code: "10000053"},
		{
			name: "bad checksum",
			code: "15027443",
			err:  "IDENTITY-01",
		},
		{
			name: "too short",
			code: "1502744",
			err:  "IDENTITY-01",
		},
		{
			name: "too long",
			code: "150274421",
			err:  "IDENTITY-01",
		},
		{
			name: "not normalized",
			code: "150-274.42",
			err:  "IDENTITY-01",
		},
		{
			name: "letters",
			code: "1502744A",
			err:  "IDENTITY-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "LU", Code: tt.code}
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
