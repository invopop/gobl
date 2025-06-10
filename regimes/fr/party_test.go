package fr_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeTaxIdentity(t *testing.T) {
	r := fr.New()
	tests := []struct {
		Code     cbc.Code
		Expected cbc.Code
	}{
		{
			Code:     "356000000", // SIREN to VAT
			Expected: "39356000000",
		},
		{
			Code:     "44 73282 9320 ",
			Expected: "44732829320",
		},
		{
			Code:     "391-838-042",
			Expected: "44391838042",
		},
		{
			Code:     "FR391838042",
			Expected: "44391838042",
		},
		{
			Code:     "FR44391838042",
			Expected: "44391838042",
		},
		{
			Code:     "341200068",
			Expected: "341200068",
		},
		{
			Code:     "39183804212345",
			Expected: "44391838042",
		},
		{
			Code:     " INV-ALID ",
			Expected: "INVALID",
		},
	}
	for _, ts := range tests {
		party := &org.Party{
			TaxID: &tax.Identity{Country: "FR", Code: ts.Code},
		}
		r.NormalizeObject(party)
		assert.Equal(t, ts.Expected, party.TaxID.Code)
	}
}
