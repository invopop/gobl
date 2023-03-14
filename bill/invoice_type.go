package bill

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/validation"
)

// Predefined list of the invoice type codes officially supported.
const (
	InvoiceTypeStandard   cbc.Key = "standard"
	InvoiceTypeProforma   cbc.Key = "proforma"
	InvoiceTypeCorrective cbc.Key = "corrective"
	InvoiceTypeCreditNote cbc.Key = "credit-note"
	InvoiceTypeDebitNote  cbc.Key = "debit-note"
)

type invoiceTypeDefs []InvoiceTypeDef

// InvoiceTypeDef is used to describe a type definition.
type InvoiceTypeDef struct {
	Key         cbc.Key  `json:"key" jsonschema:"title=InvoiceType Key"`
	Description string   `json:"description" jsonschema:"title=Description"`
	UNTDID1001  cbc.Code `json:"untdid1001" jsonschema:"title=UNTDID 1001 Code"`
}

// InvoiceTypes describes each of the InvoiceTypes supported by
// GOBL invoices, and includes a reference to the matching
// UNTDID 1001 code.
var InvoiceTypes = invoiceTypeDefs{
	{InvoiceTypeStandard, "A regular commercial invoice document between a supplier and customer.", "380"},
	{InvoiceTypeProforma, "For a clients validation before sending a final invoice.", "325"},
	{InvoiceTypeCorrective, "Corrected invoice that completely replaces the preceding document.", "384"},
	{InvoiceTypeCreditNote, "Reflects a refund either partial or complete of the preceding document.", "381"},
	{InvoiceTypeDebitNote, "An additional set of charges to be added to the preceding document.", "383"},
}

var isValidInvoiceType = validation.In(validInvoiceTypes()...)

// UNTDID1001 provides the official code number assigned with the Invoice type.
func (l invoiceTypeDefs) UNTDID1001(key cbc.Key) cbc.Code {
	for _, d := range l {
		if d.Key == key {
			return d.UNTDID1001
		}
	}
	return cbc.CodeEmpty
}

func validInvoiceTypes() []interface{} {
	list := make([]interface{}, len(InvoiceTypes))
	for i, d := range InvoiceTypes {
		list[i] = d.Key
	}
	return list
}

// UNTDID1001 provides the official code number assigned with the Invoice type.
func (i *Invoice) UNTDID1001() cbc.Code {
	for _, d := range InvoiceTypes {
		if d.Key == i.Type {
			return d.UNTDID1001
		}
	}
	return cbc.CodeEmpty
}
