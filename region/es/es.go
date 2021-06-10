package es

import (
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

// Spain holds everything related to spanish documents and taxes.
type Spain struct{}

// New provides the Spanish region definition
func New() *Spain {
	return new(Spain)
}

// Taxes provides all of this regions tax definitions.
func (Spain) Taxes() *tax.Region {
	return &taxRegion
}

func (Spain) Currency() *currency.Def {
	d, ok := currency.Get(currency.Code("EUR"))
	if !ok {
		return nil
	}
	return &d
}

// ValidateTaxID
func (Spain) ValidateTaxID(id *org.TaxID) error {
	code := id.Code
	return VerifyTaxCode(code)
}
