package head

import "github.com/invopop/gobl/schema"

// CorrectionOptions is used to define base correction options that can
// be shared between documents.
type CorrectionOptions struct {
	Head *Header `json:"-"` // copy of the original document header
}

// WithHead ensures the original envelope's header is included in the set
// of correction options.
func WithHead(header *Header) schema.Option {
	return func(o interface{}) {
		opts := o.(*CorrectionOptions)
		opts.Head = header
	}
}
