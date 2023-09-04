package base

// Option is a generic single function intended to be used
// for handling options to method calls.
type Option func(o interface{})

// CorrectionOptions is used to define base correction options that can
// be shared between documents.
type CorrectionOptions struct {
	Head *Header `json:"-"` // copy of the original document header
}

// WithHead ensures the original envelope's header is included in the set
// of correction options.
func WithHead(header *Header) Option {
	return func(o interface{}) {
		opts := o.(*CorrectionOptions)
		opts.Head = header
	}
}
