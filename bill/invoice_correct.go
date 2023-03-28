package bill

import (
	"errors"
	"fmt"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
)

// correctionOptions defines a structure used to pass configuration options
// for certain method calls.
type correctionOptions struct {
	Stamps           []*cbc.Stamp
	Refund           bool
	Append           bool
	Reason           string
	CorrectionMethod cbc.Key
	Corrections      []cbc.Key
}

// WithStamps provides a configuration option with stamp information
// usually included in the envelope header for a previously generated
// and processed invoice document.
func WithStamps(stamps []*cbc.Stamp) cbc.Option {
	return func(o interface{}) {
		opts := o.(*correctionOptions)
		opts.Stamps = stamps
	}
}

// WithReason allows a reason to be provided for the corrective operation.
func WithReason(reason string) cbc.Option {
	return func(o interface{}) {
		opts := o.(*correctionOptions)
		opts.Reason = reason
	}
}

// WithCorrectionMethod defines the method used to correct the previous invoice.
func WithCorrectionMethod(method cbc.Key) cbc.Option {
	return func(o interface{}) {
		opts := o.(*correctionOptions)
		opts.CorrectionMethod = method
	}
}

// WithCorrection adds a single correction key to the invoice preceding data,
// use multiple times for multiple entries.
func WithCorrection(correction cbc.Key) cbc.Option {
	return func(o interface{}) {
		opts := o.(*correctionOptions)
		opts.Corrections = append(opts.Corrections, correction)
	}
}

// Refund indicates that the corrective operation is a refund.
var Refund cbc.Option = func(o interface{}) {
	opts := o.(*correctionOptions)
	opts.Refund = true
}

// Append indicates that the corrective operation is to append
// new items to the previous invoice.
var Append cbc.Option = func(o interface{}) {
	opts := o.(*correctionOptions)
	opts.Append = true
}

// Correct builds a new invoice using this one as a base for the preceding
// data.
func (inv *Invoice) Correct(opts ...cbc.Option) (*Invoice, error) {
	o := new(correctionOptions)
	for _, row := range opts {
		row(o)
	}

	r := taxRegimeFor(inv.Supplier)
	if r == nil {
		return nil, errors.New("failed to load supplier regime")
	}

	i2, err := inv.Clone()
	if err != nil {
		return nil, err
	}

	// Copy and prepare the basic fields
	i2.UUID = nil
	i2.Series = ""
	i2.Code = ""
	i2.IssueDate = cal.Today()
	pre := &Preceding{
		UUID:             inv.UUID,
		Series:           inv.Series,
		Code:             inv.Code,
		IssueDate:        &inv.IssueDate,
		Reason:           o.Reason,
		Corrections:      o.Corrections,
		CorrectionMethod: o.CorrectionMethod,
	}

	// Take the regime def to figure out what needs to be copied
	if o.Refund {
		if InvoiceTypeCreditNote.In(r.Preceding.Types...) {
			// regular credit note
			i2.Type = InvoiceTypeCreditNote
		} else if InvoiceTypeCorrective.In(r.Preceding.Types...) {
			// corrective invoice with negative values
			i2.Type = InvoiceTypeCorrective
			i2.Invert()
		} else {
			return nil, errors.New("refund not supported by regime")
		}
	} else if o.Append {
		if InvoiceTypeDebitNote.In(r.Preceding.Types...) {
			// regular debit note
			i2.Type = InvoiceTypeDebitNote
		} else {
			return nil, errors.New("append not supported by regime")
		}
	} else {
		if InvoiceTypeCorrective.In(r.Preceding.Types...) {
			i2.Type = InvoiceTypeCorrective
		} else {
			return nil, errors.New("correction not supported by regime")
		}
	}

	// Make sure the stamps are there too
	for _, k := range r.Preceding.Stamps {
		var s *cbc.Stamp
		for _, row := range o.Stamps {
			if row.Provider == k {
				s = row
				break
			}
		}
		if s == nil {
			return nil, fmt.Errorf("missing stamp: %v", k)
		}
		pre.Stamps = append(pre.Stamps, s)
	}

	i2.Preceding = append(i2.Preceding, pre)
	return i2, nil
}
