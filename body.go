package gobl

// Body represents the main payload of the document whose content's
// format is determined from the type defined in the header.
type Body interface {
	Type() BodyType // string representation of body type expected
}

// BodyType defines the accepted body types
type BodyType string

// Set of defined main body types. If not defined here,
// it cannot be used in the body of a GoBL document.
const (
	BodyTypeInvoice BodyType = "Invoice"
)
