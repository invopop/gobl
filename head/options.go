package head

import "github.com/invopop/gobl/schema"

// CorrectionOptions is used to define base correction options that can
// be shared between documents.
type CorrectionOptions struct {
	Head *Header `json:"-"` // copy of the original document header
}

func (co *CorrectionOptions) setHeader(header *Header) {
	co.Head = header
}

type correctionOptionsPtr interface {
	setHeader(*Header)
}

// WithHead ensures the original envelope's header is included in the set
// of correction options. If the head.CorrectionOptions is not defined
// in the options, this will be ignored.
func WithHead(header *Header) schema.Option {
	return func(o interface{}) {
		opts, ok := o.(correctionOptionsPtr)
		if ok {
			opts.setHeader(header)
		}
	}
}
