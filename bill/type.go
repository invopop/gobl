package bill

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/org"
	"github.com/invopop/jsonschema"
)

// Type defines the type of invoice document according to a subset of the UNTDID 1001
// standard list.
type Type org.Key

// Predefined list of the invoice type codes officially supported.
const (
	TypeNone       Type = ""            // None specified
	TypeProforma   Type = "proforma"    // Proforma invoice
	TypeSimplified Type = "simplified"  // Simplified Invoice
	TypePartial    Type = "partial"     // Partial Invoice
	TypeCommercial Type = "commercial"  // Commercial Invoice
	TypeCorrected  Type = "corrected"   // Corrected Invoice
	TypeCreditNote Type = "credit-note" // Credit Note
	TypeSelfBilled Type = "self-billed" // Self Billed Invoice
)

// TypeDef is used to describe a type definition.
type TypeDef struct {
	Key         Type     `json:"key" jsonschema:"title=Type Key"`
	Description string   `json:"description" jsonschema:"title=Description"`
	UNTDID1001  org.Code `json:"untdid1001" jsonschema:"title=UNTDID 1001 Code"`
}

// TypeDefinitions describes each of the Types supported by
// GOBL invoices, and includes a reference to the matching
// UNTDID 1001 code.
var TypeDefinitions = []TypeDef{
	{TypeProforma, "Proforma invoice, for a clients validation before sending a final invoice.", "325"},
	{TypeSimplified, "Simplified invoice or receipt typically used for small transactions that don't require customer details.", "380"}, // same UNTDID as commercial
	{TypePartial, "Partial invoice", "326"},
	{TypeCommercial, "Commercial invoice, usually cross-border transactions requiring a customs invoice.", "380"},
	{TypeCorrected, "Corrected invoice", "384"},
	{TypeCreditNote, "Credit note", "381"},
	{TypeSelfBilled, "Self billed invoice", "389"},
}

var isValidType = validation.In(validTypes()...)

func validTypes() []interface{} {
	list := make([]interface{}, len(TypeDefinitions))
	for i, d := range TypeDefinitions {
		list[i] = string(d.Key)
	}
	return list
}

// Validate is used to ensure the code provided is one of those we know
// about.
func (k Type) Validate() error {
	return validation.Validate(string(k), isValidType)
}

// UNTDID1001 provides the official code number assigned to the type.
func (k Type) UNTDID1001() org.Code {
	for _, d := range TypeDefinitions {
		if d.Key == k {
			return d.UNTDID1001
		}
	}
	return org.CodeEmpty
}

// In checks to see if the type key equals one of the
// provided set.
func (k Type) In(set ...Type) bool {
	for _, v := range set {
		if v == k {
			return true
		}
	}
	return false
}

// JSONSchema provides a representation of the struct for usage in Schema.
func (Type) JSONSchema() *jsonschema.Schema {
	s := &jsonschema.Schema{
		Title:       "Type",
		Type:        "string", // they're all strings
		OneOf:       make([]*jsonschema.Schema, len(TypeDefinitions)),
		Description: "Defines the type of invoice document according to a subset of the UNTDID 1001 standard list.",
	}
	for i, v := range TypeDefinitions {
		s.OneOf[i] = &jsonschema.Schema{
			Const:       org.Key(v.Key).String(),
			Description: v.Description,
		}
	}
	return s
}
