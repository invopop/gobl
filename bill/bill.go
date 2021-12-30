package bill

import "errors"

// Standard billing model types that can be incorporated into an
// envelope.
const (
	InvoiceType = "bill.Invoice"
)

// TypeCode defines the "Invoice Type Code" according to a subset of the UNTDID 1001
// standard list.
type TypeCode string

// Predefined list of the invoice type codes officially supported.
const (
	CommercialTypeCode TypeCode = ""            // Commercial Invoice, default
	ProformaTypeCode   TypeCode = "proforma"    // Proforma invoice
	SimplifiedTypeCode TypeCode = "simplified"  // Simplified Invoice
	PartialTypeCode    TypeCode = "partial"     // Partial Invoice
	CorrectedTypeCode  TypeCode = "corrected"   // Corrected Invoice
	CreditNoteTypeCode TypeCode = "credit-note" // Credit Note
	SelfBilledTypeCode TypeCode = "self-billed" // Self Billed Invoice
)

// UNTDID1001TypeCodeMap offers a way to convert the GOBL invoice type code into
// one supported by our subset of the UNTDID 1001 official list.
var UNTDID1001TypeCodeMap = map[TypeCode]string{
	ProformaTypeCode:   "325",
	PartialTypeCode:    "326",
	CommercialTypeCode: "380",
	SimplifiedTypeCode: "380", // same as commercial
	CorrectedTypeCode:  "384",
	CreditNoteTypeCode: "381",
	SelfBilledTypeCode: "389",
}

// Validate is used to ensure the code provided is one of those we know
// about.
func (c TypeCode) Validate() error {
	_, ok := UNTDID1001TypeCodeMap[c]
	if !ok {
		return errors.New("not found")
	}
	return nil
}

// UNTDID1001 provides the official code number assigned to the type.
func (c TypeCode) UNTDID1001() string {
	s, ok := UNTDID1001TypeCodeMap[c]
	if !ok {
		return "na"
	}
	return s
}
