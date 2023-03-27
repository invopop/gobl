// Package bill provides models for dealing with Billing and specifically invoicing.
package bill

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/schema"
)

func init() {
	schema.Register(schema.GOBL.Add("bill"),
		// None of bill's sub-models are meant to be used outside an invoice.
		Invoice{},
	)
}

// Corrector defines the method used to build corrective invoices.
type Corrector interface {
	Correct(inv *Invoice, opts *Options) error
}

// Options defines a structure used to pass configuration options
// for certain method calls. This is only meant to be used internally.
type Options struct {
	Stamps           []*cbc.Stamp
	Refund           bool
	Append           bool
	Reason           string
	CorrectionMethod cbc.Key
	Corrections      []cbc.Key
	Previous         *Invoice
}

// WithStamps provides a configuration option with stamp information
// usually included in the envelope header for a previously generated
// and processed invoice document.
func WithStamps(stamps []*cbc.Stamp) cbc.Option {
	return func(o interface{}) {
		opts := o.(*Options)
		opts.Stamps = stamps
	}
}

// WithReason allows a reason to be provided for the corrective operation.
func WithReason(reason string) cbc.Option {
	return func(o interface{}) {
		opts := o.(*Options)
		opts.Reason = reason
	}
}

// WithCorrectionMethod defines the method used to correct the previous invoice.
func WithCorrectionMethod(method cbc.Key) cbc.Option {
	return func(o interface{}) {
		opts := o.(*Options)
		opts.CorrectionMethod = method
	}
}

// WithCorrection adds a single correction key to the invoice preceding data,
// use multiple times for multiple entries.
func WithCorrection(correction cbc.Key) cbc.Option {
	return func(o interface{}) {
		opts := o.(*Options)
		opts.Corrections = append(opts.Corrections, correction)
	}
}

// Refund indicates that the corrective operation is a refund.
var Refund cbc.Option = func(o interface{}) {
	opts := o.(*Options)
	opts.Refund = true
}

// Append indicates that the corrective operation is to append
// new items to the previous invoice.
var Append cbc.Option = func(o interface{}) {
	opts := o.(*Options)
	opts.Append = true
}
