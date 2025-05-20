package bill

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

// Predefined list of the invoice type codes officially supported.
const (
	InvoiceTypeStandard   cbc.Key = "standard"
	InvoiceTypeProforma   cbc.Key = "proforma"
	InvoiceTypeCorrective cbc.Key = "corrective"
	InvoiceTypeCreditNote cbc.Key = "credit-note"
	InvoiceTypeDebitNote  cbc.Key = "debit-note"
	InvoiceTypeOther      cbc.Key = "other"
)

const (
	// UNTDID1001Key is the key used to identify the UNTDID 1001 code
	// associated with an invoice type.
	UNTDID1001Key cbc.Key = "untdid1001"
)

// InvoiceTypes describes each of the InvoiceTypes supported.
var InvoiceTypes = []*cbc.Definition{
	{
		Key: InvoiceTypeStandard,
		Name: i18n.String{
			i18n.EN: "Standard",
		},
		Desc: i18n.String{
			i18n.EN: "A regular commercial invoice document between a supplier and customer.",
		},
		Map: cbc.CodeMap{
			UNTDID1001Key: "380",
		},
	},
	{
		Key: InvoiceTypeProforma,
		Name: i18n.String{
			i18n.EN: "Proforma",
		},
		Desc: i18n.String{
			i18n.EN: "For a clients validation before sending a final invoice.",
		},
		Map: cbc.CodeMap{
			UNTDID1001Key: "325",
		},
	},
	{
		Key: InvoiceTypeCorrective,
		Name: i18n.String{
			i18n.EN: "Corrective",
		},
		Desc: i18n.String{
			i18n.EN: "Corrected invoice that completely *replaces* the preceding document.",
		},
		Map: cbc.CodeMap{
			UNTDID1001Key: "384",
		},
	},
	{
		Key: InvoiceTypeCreditNote,
		Name: i18n.String{
			i18n.EN: "Credit Note",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Reflects a refund either partial or complete of the preceding document. A 
				credit note effectively *extends* the previous document.
			`),
		},
		Map: cbc.CodeMap{
			UNTDID1001Key: "381",
		},
	},
	{
		Key: InvoiceTypeDebitNote,
		Name: i18n.String{
			i18n.EN: "Debit Note",
		},
		Desc: i18n.String{
			i18n.EN: "An additional set of charges to be added to the preceding document.",
		},
		Map: cbc.CodeMap{
			UNTDID1001Key: "383",
		},
	},
	{
		Key: InvoiceTypeOther,
		Name: i18n.String{
			i18n.EN: "Other",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Any other type of invoice that does not fit into the standard categories and implies
				that any scenarios defined in tax regimes or addons will not be applied.

				This is useful for being able to create invoices with custom types in extensions,
				but is not recommend for general use.
			`),
		},
	},
}

var isValidInvoiceType = cbc.InKeyDefs(InvoiceTypes)

// UNTDID1001 provides the official code number assigned with the Invoice type.
func (inv *Invoice) UNTDID1001() cbc.Code {
	for _, d := range InvoiceTypes {
		if d.Key == inv.Type {
			return d.Map[UNTDID1001Key]
		}
	}
	return cbc.CodeEmpty
}
