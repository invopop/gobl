package ad_test

import (
	"testing"

	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/ad"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestValidateParty(t *testing.T) {
	tests := []struct {
		name  string
		party *org.Party
		err   string
	}{
		{
			name:  "nil party",
			party: nil,
		},
		{
			name:  "no tax ID",
			party: &org.Party{Name: "No Tax"},
		},
		{
			name: "domestic AD party",
			party: &org.Party{
				Name:  "Fusta i Disseny S.L.",
				TaxID: &tax.Identity{Country: "AD", Code: "L123456A"},
			},
		},
		{
			name: "non-resident with tax ID",
			party: &org.Party{
				Name:  "Acme S.A.",
				TaxID: &tax.Identity{Country: "ES", Code: "B12345678"},
			},
		},
		{
			name: "non-resident missing tax ID",
			party: &org.Party{
				Name:  "Ghost Corp",
				TaxID: &tax.Identity{Country: "FR"},
			},
			err: "non-resident party must provide a tax ID or passport number",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ad.Validate(tt.party)
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