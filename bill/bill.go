package bill

import (
	"errors"

	"github.com/invopop/gobl/schema"
)

func init() {
	// None of TypeKey's sub-models are meant to be used outside an invoice.
	schema.Register(schema.GOBL.Add("bill"), Invoice{})
}

// TypeKey defines the type of invoice document according to a subset of the UNTDID 1001
// standard list.
type TypeKey string

// Predefined list of the invoice type codes officially supported.
const (
	TypeKeyCommercial TypeKey = ""            // Commercial Invoice, default
	TypeKeyProforma   TypeKey = "proforma"    // Proforma invoice
	TypeKeySimplified TypeKey = "simplified"  // Simplified Invoice
	TypeKeyPartial    TypeKey = "partial"     // Partial Invoice
	TypeKeyCorrected  TypeKey = "corrected"   // Corrected Invoice
	TypeKeyCreditNote TypeKey = "credit-note" // Credit Note
	TypeKeySelfBilled TypeKey = "self-billed" // Self Billed Invoice
)

// UNTDID1001TypeKeyMap offers a way to convert the GOBL invoice type code into
// one supported by our subset of the UNTDID 1001 official list.
var UNTDID1001TypeKeyMap = map[TypeKey]string{
	TypeKeyProforma:   "325",
	TypeKeyPartial:    "326",
	TypeKeyCommercial: "380",
	TypeKeySimplified: "380", // same as commercial
	TypeKeyCorrected:  "384",
	TypeKeyCreditNote: "381",
	TypeKeySelfBilled: "389",
}

// Validate is used to ensure the code provided is one of those we know
// about.
func (c TypeKey) Validate() error {
	_, ok := UNTDID1001TypeKeyMap[c]
	if !ok {
		return errors.New("not found")
	}
	return nil
}

// UNTDID1001 provides the official code number assigned to the type.
func (c TypeKey) UNTDID1001() string {
	s, ok := UNTDID1001TypeKeyMap[c]
	if !ok {
		return "na"
	}
	return s
}
