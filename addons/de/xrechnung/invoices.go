package xrechnung

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
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
