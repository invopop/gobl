package es

// BillInvoiceLineMeta defines additional fields that may be added and used
// in an invoice line.
type BillInvoiceLineMeta struct {
	// When true, this line should be considered as being sourced from a provider
	// under a "Equivalence Surcharge VAT" regime.
	Supplied bool `json:"supplied,omitempty" jsonschema:"title=Supplied"`

	// Message that explains why this line is exempt of taxes.
	Exempt string `json:"exempt,omitempty" jsonschema:"title=Exempt"`
}
