package xrechnung

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// BR-DE-17 - restricted subset of UNTDID document type codes
var validInvoiceUNTDIDDocumentTypeValues = []cbc.Code{
	"326", // Partial
	"380", // Commercial
	"384", // Corrected
	"389", // Self-billed
	"381", // Credit note
	"875", // Partial construction invoice
	"876", // Partial Final construction invoice
	"877", // Final construction invoice
}

// validateInvoice validates the invoice according to the XRechnung standard
func validateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		// BR-DE-17
		validation.Field(&inv.Tax,
			validation.By(validateInvoiceTax),
			validation.Skip,
		),
		validation.Field(&inv.Preceding,
			validation.When(
				inv.Type.In(
					bill.InvoiceTypeCorrective,
					bill.InvoiceTypeCreditNote,
				),
				validation.Required,
			),
			validation.Skip,
		),
		// BR-DE-6
		validation.Field(&inv.Customer,
			validation.Required,
			validation.By(validateParty),
			validation.Skip,
		),
	)
}

func validateInvoiceTax(value any) error {
	tx, ok := value.(*bill.Tax)
	if !ok || tx == nil {
		return nil
	}
	return validation.ValidateStruct(tx,
		validation.Field(&tx.Ext,
			tax.ExtensionsHasCodes(untdid.ExtKeyTaxCategory, validInvoiceUNTDIDDocumentTypeValues...),
			validation.Skip,
		),
	)
}

func validateParty(val any) error {
	party, ok := val.(*org.Party)
	if !ok {
		return nil
	}

	// Check if either a person has a telephone number or the party has a telephone number
	hasTelephone := len(party.Telephones) > 0

	// If party doesn't have telephones directly, check if any person has telephones
	if !hasTelephone {
		for _, person := range party.People {
			if len(person.Telephones) > 0 {
				hasTelephone = true
				break
			}
		}
	}

	if !hasTelephone {
		return validation.NewError(
			"xrechnung.party.telephone.required",
			"Either the party or at least one person must have a telephone number",
		)
	}

	return nil
}
