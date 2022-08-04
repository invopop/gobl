package bill

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/org"
	"github.com/invopop/jsonschema"
)

// InvoiceType defines the type of invoice document according to a subset of the UNTDID 1001
// standard list.
type InvoiceType org.Key

// Predefined list of the invoice type codes officially supported.
const (
	InvoiceTypeNone       InvoiceType = ""            // None specified
	InvoiceTypeProforma   InvoiceType = "proforma"    // Proforma invoice
	InvoiceTypeSimplified InvoiceType = "simplified"  // Simplified Invoice
	InvoiceTypePartial    InvoiceType = "partial"     // Partial Invoice
	InvoiceTypeCommercial InvoiceType = "commercial"  // Commercial Invoice
	InvoiceTypeCorrected  InvoiceType = "corrected"   // Corrected Invoice
	InvoiceTypeCreditNote InvoiceType = "credit-note" // Credit Note
	InvoiceTypeSelfBilled InvoiceType = "self-billed" // Self Billed Invoice
)

// InvoiceTypeDef is used to describe a type definition.
type InvoiceTypeDef struct {
	Key         InvoiceType `json:"key" jsonschema:"title=InvoiceType Key"`
	Description string      `json:"description" jsonschema:"title=Description"`
	UNTDID1001  org.Code    `json:"untdid1001" jsonschema:"title=UNTDID 1001 Code"`
}

// InvoiceTypeDefinitions describes each of the InvoiceTypes supported by
// GOBL invoices, and includes a reference to the matching
// UNTDID 1001 code.
var InvoiceTypeDefinitions = []InvoiceTypeDef{
	{InvoiceTypeProforma, "Proforma invoice, for a clients validation before sending a final invoice.", "325"},
	{InvoiceTypeSimplified, "Simplified invoice or receipt typically used for small transactions that don't require customer details.", "380"}, // same UNTDID as commercial
	{InvoiceTypePartial, "Partial invoice", "326"},
	{InvoiceTypeCommercial, "Commercial invoice, usually cross-border transactions requiring an invoice for customs.", "380"},
	{InvoiceTypeCorrected, "Corrected invoice", "384"},
	{InvoiceTypeCreditNote, "Credit note", "381"},
	{InvoiceTypeSelfBilled, "Self billed invoice", "389"},
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
func (k InvoiceType) UNTDID1001() org.Code {
	for _, d := range InvoiceTypeDefinitions {
		if d.Key == k {
			return d.UNTDID1001
		}
	}
	return org.CodeEmpty
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
			Const:       org.Key(v.Key).String(),
			Description: v.Description,
		}
	}
	return s
}
