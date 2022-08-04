package bill

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/org"
	"github.com/invopop/jsonschema"
)

// TypeKey defines the type of invoice document according to a subset of the UNTDID 1001
// standard list.
type TypeKey org.Key

// Predefined list of the invoice type codes officially supported.
const (
	TypeKeyNone       TypeKey = ""            // None specified
	TypeKeyCommercial TypeKey = "commercial"  // Commercial Invoice
	TypeKeyProforma   TypeKey = "proforma"    // Proforma invoice
	TypeKeySimplified TypeKey = "simplified"  // Simplified Invoice
	TypeKeyPartial    TypeKey = "partial"     // Partial Invoice
	TypeKeyCorrected  TypeKey = "corrected"   // Corrected Invoice
	TypeKeyCreditNote TypeKey = "credit-note" // Credit Note
	TypeKeySelfBilled TypeKey = "self-billed" // Self Billed Invoice
)

// TypeKeyDef is used to describe a type definition.
type TypeKeyDef struct {
	Key         TypeKey  `json:"key" jsonschema:"title=Type Key"`
	Description string   `json:"description" jsonschema:"title=Description"`
	UNTDID1001  org.Code `json:"untdid1001" jsonschema:"title=UNTDID 1001 Code"`
}

// TypeKeyDefinitions describes each of the TypeKeys supported by
// GOBL invoices, and includes a reference to the matching
// UNTDID 1001 code.
var TypeKeyDefinitions = []TypeKeyDef{
	{TypeKeyCommercial, "Commercial invoice", "380"},
	{TypeKeyProforma, "Proforma invoice", "325"},
	{TypeKeySimplified, "Simplified invoice or receipt", "380"}, // same UNTDID as commercial
	{TypeKeyPartial, "Partial invoice", "326"},
	{TypeKeyCorrected, "Corrected invoice", "384"},
	{TypeKeyCreditNote, "Credit note", "381"},
	{TypeKeySelfBilled, "Self billed invoice", "389"},
}

var isValidTypeKey = validation.In(validTypeKeys()...)

func validTypeKeys() []interface{} {
	list := make([]interface{}, len(TypeKeyDefinitions))
	for i, d := range TypeKeyDefinitions {
		list[i] = string(d.Key)
	}
	return list
}

// Validate is used to ensure the code provided is one of those we know
// about.
func (k TypeKey) Validate() error {
	return validation.Validate(string(k), isValidTypeKey)
}

// UNTDID1001 provides the official code number assigned to the type.
func (k TypeKey) UNTDID1001() org.Code {
	for _, d := range TypeKeyDefinitions {
		if d.Key == k {
			return d.UNTDID1001
		}
	}
	return org.CodeEmpty
}

// In checks to see if the type key equals one of the
// provided set.
func (k TypeKey) In(set ...TypeKey) bool {
	for _, v := range set {
		if v == k {
			return true
		}
	}
	return false
}

// JSONSchema provides a representation of the struct for usage in Schema.
func (TypeKey) JSONSchema() *jsonschema.Schema {
	s := &jsonschema.Schema{
		Title:       "Type Key",
		Type:        "string", // they're all strings
		OneOf:       make([]*jsonschema.Schema, len(TypeKeyDefinitions)),
		Description: "TypeKey defines the type of invoice document according to a subset of the UNTDID 1001 standard list.",
	}
	for i, v := range TypeKeyDefinitions {
		s.OneOf[i] = &jsonschema.Schema{
			Const:       org.Key(v.Key).String(),
			Description: v.Description,
		}
	}
	return s
}
