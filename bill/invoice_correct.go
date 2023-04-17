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

// Correct moves key fields of the current invoice to the preceding
// structure and performs any regime specific actions defined by the
// regime's configuration.
// If the existing document doesn't have a code, we'll raise an error, for
// most use cases this will prevent looping over the same invoice.
func (inv *Invoice) Correct(opts ...cbc.Option) error {
	o := new(correctionOptions)
	for _, row := range opts {
		row(o)
	}
	if o.credit && o.debit {
		return errors.New("cannot use both credit and debit options")
	}
	if inv.Code == "" {
		return errors.New("cannot correct an invoice without a code")
	}

	r := taxRegimeFor(inv.Supplier)
	if r == nil {
		return errors.New("failed to load supplier regime")
	}

	// Copy and prepare the basic fields
	iDate := inv.IssueDate
	pre := &Preceding{
		UUID:             inv.UUID,
		Series:           inv.Series,
		Code:             inv.Code,
		IssueDate:        &iDate,
		Reason:           o.reason,
		Corrections:      o.corrections,
		CorrectionMethod: o.correctionMethod,
	}
	inv.UUID = nil
	inv.Series = ""
	inv.Code = ""
	inv.IssueDate = cal.Today()

	// Take the regime def to figure out what needs to be copied
	if o.credit {
		if r.Preceding.HasType(InvoiceTypeCreditNote) {
			// regular credit note
			inv.Type = InvoiceTypeCreditNote
		} else if r.Preceding.HasType(InvoiceTypeCorrective) {
			// corrective invoice with negative values
			inv.Type = InvoiceTypeCorrective
			inv.Invert()
		} else {
			return errors.New("credit not supported by regime")
		}
		inv.Payment.ResetAdvances()
	} else if o.debit {
		if r.Preceding.HasType(InvoiceTypeDebitNote) {
			// regular debit note, implies no rows as new ones
			// will be added
			inv.Type = InvoiceTypeDebitNote
			inv.Empty()
		} else {
			return errors.New("debit not supported by regime")
		}
	} else {
		if r.Preceding.HasType(InvoiceTypeCorrective) {
			inv.Type = InvoiceTypeCorrective
		} else {
			return errors.New("correction not supported by regime")
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
				return fmt.Errorf("missing stamp: %v", k)
			}
			pre.Stamps = append(pre.Stamps, s)
		}
	}

	// Replace all previous preceding data
	inv.Preceding = []*Preceding{pre}

	// Running a Calculate feels a bit out of place, but not performing
	// this operation on the corrected invoice results in potentially
	// conflicting or incomplete data.
	return inv.Calculate()
}
