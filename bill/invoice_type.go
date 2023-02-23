package bill

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
)

// InvoiceType defines the type of invoice document according to a subset of the UNTDID 1001
// standard list.
type InvoiceType cbc.Key

// Predefined list of the invoice type codes officially supported.
const (
	InvoiceTypeDefault    InvoiceType = ""
	InvoiceTypeProforma   InvoiceType = "proforma"
	InvoiceTypeSimplified InvoiceType = "simplified"
	InvoiceTypePartial    InvoiceType = "partial"
	InvoiceTypeCorrective InvoiceType = "corrective"
	InvoiceTypeCreditNote InvoiceType = "credit-note"
	InvoiceTypeSelfBilled InvoiceType = "self-billed"
)

// InvoiceTypeDef is used to describe a type definition.
type InvoiceTypeDef struct {
	Key         InvoiceType `json:"key" jsonschema:"title=InvoiceType Key"`
	Description string      `json:"description" jsonschema:"title=Description"`
	UNTDID1001  cbc.Code    `json:"untdid1001" jsonschema:"title=UNTDID 1001 Code"`
}

// InvoiceTypeDefinitions describes each of the InvoiceTypes supported by
// GOBL invoices, and includes a reference to the matching
// UNTDID 1001 code.
var InvoiceTypeDefinitions = []InvoiceTypeDef{
	{InvoiceTypeDefault, "A regular commercial invoice document between a supplier and customer.", "380"},
	{InvoiceTypeProforma, "For a clients validation before sending a final invoice.", "325"},
	{InvoiceTypeSimplified, "Typically used for small transactions that don't require customer details.", "380"}, // same UNTDID as commercial
	{InvoiceTypePartial, "Reflecting partial delivery of goods or services to be paid.", "326"},
	{InvoiceTypeSelfBilled, "Created by a customer on behalf of the supplier.", "389"},
	{InvoiceTypeCorrective, "Corrected invoice that completely replaces the preceding document.", "384"},
	{InvoiceTypeCreditNote, "Reflects a refund either partial or complete of the preceding document.", "381"},
}

var isValidInvoiceType = validation.In(validInvoiceTypes()...)

func validInvoiceTypes() []interface{} {
	list := make([]interface{}, len(InvoiceTypeDefinitions))
	for i, d := range InvoiceTypeDefinitions {
		list[i] = string(d.Key)
	}
	return list
}

// Validate is used to ensure the code provided is one of those we know
// about.
func (k InvoiceType) Validate() error {
	return validation.Validate(string(k), isValidInvoiceType)
}

// UNTDID1001 provides the official code number assigned to the type.
func (k InvoiceType) UNTDID1001() cbc.Code {
	for _, d := range InvoiceTypeDefinitions {
		if d.Key == k {
			return d.UNTDID1001
		}
	}
	return cbc.CodeEmpty
}

// In checks to see if the type key equals one of the
// provided set.
func (k InvoiceType) In(set ...InvoiceType) bool {
	for _, v := range set {
		if v == k {
			return true
		}
	}
	return false
}

// JSONSchema provides a representation of the struct for usage in Schema.
func (InvoiceType) JSONSchema() *jsonschema.Schema {
	s := &jsonschema.Schema{
		Title:       "Invoice Type",
		Type:        "string", // they're all strings
		OneOf:       make([]*jsonschema.Schema, len(InvoiceTypeDefinitions)),
		Description: "Defines an invoice type according to a subset of the UNTDID 1001 standard list.",
	}
	for i, v := range InvoiceTypeDefinitions {
		s.OneOf[i] = &jsonschema.Schema{
			Const:       cbc.Key(v.Key).String(),
			Description: v.Description,
		}
	}
	return s
}
