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
	stamps           []*cbc.Stamp
	credit           bool
	debit            bool
	reason           string
	correctionMethod cbc.Key
	corrections      []cbc.Key
}

// WithStamps provides a configuration option with stamp information
// usually included in the envelope header for a previously generated
// and processed invoice document.
func WithStamps(stamps []*cbc.Stamp) cbc.Option {
	return func(o interface{}) {
		opts := o.(*correctionOptions)
		opts.stamps = stamps
	}
}

// WithReason allows a reason to be provided for the corrective operation.
func WithReason(reason string) cbc.Option {
	return func(o interface{}) {
		opts := o.(*correctionOptions)
		opts.reason = reason
	}
}

// WithCorrectionMethod defines the method used to correct the previous invoice.
func WithCorrectionMethod(method cbc.Key) cbc.Option {
	return func(o interface{}) {
		opts := o.(*correctionOptions)
		opts.correctionMethod = method
	}
}

// WithCorrection adds a single correction key to the invoice preceding data,
// use multiple times for multiple entries.
func WithCorrection(correction cbc.Key) cbc.Option {
	return func(o interface{}) {
		opts := o.(*correctionOptions)
		opts.corrections = append(opts.corrections, correction)
	}
}

// Credit indicates that the corrective operation requires a credit note
// or equivalent.
var Credit cbc.Option = func(o interface{}) {
	opts := o.(*correctionOptions)
	opts.credit = true
}

// Debit indicates that the corrective operation is to append
// new items to the previous invoice, usually as a debit note.
var Debit cbc.Option = func(o interface{}) {
	opts := o.(*correctionOptions)
	opts.debit = true
}

// Correct builds a new invoice using this one as a base for the preceding
// data.
func (inv *Invoice) Correct(opts ...cbc.Option) (*Invoice, error) {
	o := new(correctionOptions)
	for _, row := range opts {
		row(o)
	}
	if o.credit && o.debit {
		return nil, errors.New("cannot use both credit and debit options")
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
		Reason:           o.reason,
		Corrections:      o.corrections,
		CorrectionMethod: o.correctionMethod,
	}

	// Take the regime def to figure out what needs to be copied
	if o.credit {
		if r.Preceding.HasType(InvoiceTypeCreditNote) {
			// regular credit note
			i2.Type = InvoiceTypeCreditNote
		} else if r.Preceding.HasType(InvoiceTypeCorrective) {
			// corrective invoice with negative values
			i2.Type = InvoiceTypeCorrective
			i2.Invert()
		} else {
			return nil, errors.New("credit not supported by regime")
		}
	} else if o.debit {
		if r.Preceding.HasType(InvoiceTypeDebitNote) {
			// regular debit note, implies no rows as new ones
			// will be added
			i2.Type = InvoiceTypeDebitNote
			i2.Empty()
		} else {
			return nil, errors.New("debit not supported by regime")
		}
	} else {
		if r.Preceding.HasType(InvoiceTypeCorrective) {
			i2.Type = InvoiceTypeCorrective
		} else {
			return nil, errors.New("correction not supported by regime")
		}
	}

	// Make sure the stamps are there too
	if r.Preceding != nil {
		for _, k := range r.Preceding.Stamps {
			var s *cbc.Stamp
			for _, row := range o.stamps {
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
	}

	// Replace all previous preceding data
	i2.Preceding = []*Preceding{pre}
	return i2, nil
}
