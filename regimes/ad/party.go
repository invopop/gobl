package ad

import (
	"errors"

	"github.com/invopop/gobl/org"
	"github.com/invopop/validation"
)

// validateParty checks that non-resident parties supply a tax ID or
// passport number, as required by Forms 980-A and 980-B.
func validateParty(party *org.Party) error {
	if party == nil || party.TaxID == nil {
		return nil
	}
	if party.TaxID.Country.In("AD") {
		return nil // domestic parties are fine
	}
	// For non-resident parties a tax ID code is required.
	// Full fiscal-representative enforcement (name + NRT) is handled
	// at the invoice addon layer via the "non-resident-b2b" tag.
	return validation.ValidateStruct(party,
		validation.Field(&party.TaxID,
			validation.By(func(v interface{}) error {
				if party.TaxID.Code == "" {
					return errors.New("non-resident party must provide a tax ID or passport number")
				}
				return nil
			}),
		),
	)
}