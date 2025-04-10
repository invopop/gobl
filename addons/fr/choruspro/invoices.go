package choruspro

import (
	"errors"
	"fmt"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/validation"
)

func normalizeInvoice(inv *bill.Invoice) {
	normalizeParty(inv.Supplier)
	normalizeParty(inv.Customer)
}

func normalizeParty(party *org.Party) {
	if party == nil || party.TaxID == nil {
		return
	}

	for _, identity := range party.Identities {
		if identity.Type == fr.IdentityTypeSIREN || identity.Type == fr.IdentityTypeSIRET {
			return
		}
	}

	fmt.Println("party.TaxID.Code.String()", party.TaxID.Code.String())
	// These checks ensure that we only try to add an identity if we can create a valid SIREN.
	// This is because this occurs before regime normalization.
	if len(party.TaxID.Code.String()) == 11 {
		if party.Identities == nil {
			party.Identities = make([]*org.Identity, 0)
		}
		party.Identities = append(party.Identities, &org.Identity{
			Type: fr.IdentityTypeSIREN,
			Code: cbc.Code(party.TaxID.Code.String()[2:]),
		})
	}

	if len(party.TaxID.Code.String()) == 9 {
		if party.Identities == nil {
			party.Identities = make([]*org.Identity, 0)
		}
		party.Identities = append(party.Identities, &org.Identity{
			Type: fr.IdentityTypeSIRET,
			Code: cbc.Code(party.TaxID.Code.String()),
		})
	}

	if len(party.TaxID.Code.String()) == 14 {
		if party.Identities == nil {
			party.Identities = make([]*org.Identity, 0)
		}
		party.Identities = append(party.Identities, &org.Identity{
			Type: fr.IdentityTypeSIRET,
			Code: cbc.Code(party.TaxID.Code.String()),
		})
	}

}

func validateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Payment,
			validation.Required,
			validation.By(validatePayment),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.Required,
			validation.By(validateParty),
			validation.Skip,
		),
		validation.Field(&inv.Supplier,
			validation.Required,
			validation.By(validateParty),
			validation.Skip,
		),
	)
}

func validatePayment(val any) error {
	payment, ok := val.(*bill.PaymentDetails)
	if !ok {
		return nil
	}
	// Rest of the validation is handled by en16931 addon
	return validation.ValidateStruct(payment,
		validation.Field(&payment.Instructions,
			validation.Required,
			validation.Skip,
		),
	)
}

func validateParty(val any) error {
	party, ok := val.(*org.Party)
	if !ok {
		return nil
	}

	return validation.ValidateStruct(party,
		validation.Field(&party.Identities,
			validation.Required,
			validation.By(validateIdentities),
			validation.Skip,
		),
	)
}

func validateIdentities(val any) error {
	identities, ok := val.([]*org.Identity)
	if !ok {
		return nil
	}
	for _, identity := range identities {
		if identity.Type == fr.IdentityTypeSIREN || identity.Type == fr.IdentityTypeSIRET {
			return nil
		}
	}

	return errors.New("at least one identity must be SIREN or SIRET")
}
