package bill

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
)

// CorrectionOptions defines a structure used to pass configuration options
// to correct a previous invoice. This is made available to make it easier to
// pass options between external services.
type CorrectionOptions struct {
	// When the new corrective invoice's issue date should be set to.
	IssueDate *cal.Date `json:"issue_date,omitempty" jsonschema:"title=Issue Date"`
	// Stamps of the previous document to include in the preceding data.
	Stamps []*cbc.Stamp `json:"stamps,omitempty" jsonschema:"title=Stamps"`
	// Credit when true indicates that the corrective document should cancel the previous document.
	Credit bool `json:"credit,omitempty" jsonschema:"title=Credit"`
	// Debit when true indicates that the corrective document should add new items to the previous document.
	Debit bool `json:"debit,omitempty" jsonschema:"title=Debit"`
	// Human readable reason for the corrective operation.
	Reason string `json:"reason,omitempty" jsonschema:"title=Reason"`
	// Correction method as defined by the tax regime.
	CorrectionMethod cbc.Key `json:"correction_method,omitempty" jsonschema:"title=Correction Method"`
	// Correction keys that describe the specific changes according to the tax regime.
	Corrections []cbc.Key `json:"corrections,omitempty" jsonschema:"title=Corrections"`

	// In case we want to use a raw json object as a source of the options.
	data json.RawMessage `json:"-"`
}

// WithOptions takes an already completed CorrectionOptions instance and
// uses this as a base instead of passing individual options. This is useful
// for passing options from an API, developers should use the regular option
// methods.
func WithOptions(opts *CorrectionOptions) cbc.Option {
	return func(o interface{}) {
		o2 := o.(*CorrectionOptions)
		*o2 = *opts
	}
}

// WithData expects a raw JSON object that will be marshalled into a
// CorrectionOptions instance and used as the base for the correction.
func WithData(data json.RawMessage) cbc.Option {
	return func(o interface{}) {
		opts := o.(*CorrectionOptions)
		opts.data = data
	}
}

// WithStamps provides a configuration option with stamp information
// usually included in the envelope header for a previously generated
// and processed invoice document.
func WithStamps(stamps []*cbc.Stamp) cbc.Option {
	return func(o interface{}) {
		opts := o.(*CorrectionOptions)
		opts.Stamps = stamps
	}
}

// WithReason allows a reason to be provided for the corrective operation.
func WithReason(reason string) cbc.Option {
	return func(o interface{}) {
		opts := o.(*CorrectionOptions)
		opts.Reason = reason
	}
}

// WithCorrectionMethod defines the method used to correct the previous invoice.
func WithCorrectionMethod(method cbc.Key) cbc.Option {
	return func(o interface{}) {
		opts := o.(*CorrectionOptions)
		opts.CorrectionMethod = method
	}
}

// WithCorrection adds a single correction key to the invoice preceding data,
// use multiple times for multiple entries.
func WithCorrection(correction cbc.Key) cbc.Option {
	return func(o interface{}) {
		opts := o.(*CorrectionOptions)
		opts.Corrections = append(opts.Corrections, correction)
	}
}

// WithIssueDate can be used to override the issue date of the corrective invoice
// produced.
func WithIssueDate(date cal.Date) cbc.Option {
	return func(o interface{}) {
		opts := o.(*CorrectionOptions)
		opts.IssueDate = &date
	}
}

// Credit indicates that the corrective operation requires a credit note
// or equivalent.
var Credit cbc.Option = func(o interface{}) {
	opts := o.(*CorrectionOptions)
	opts.Credit = true
}

// Debit indicates that the corrective operation is to append
// new items to the previous invoice, usually as a debit note.
var Debit cbc.Option = func(o interface{}) {
	opts := o.(*CorrectionOptions)
	opts.Debit = true
}

// Correct moves key fields of the current invoice to the preceding
// structure and performs any regime specific actions defined by the
// regime's configuration.
// If the existing document doesn't have a code, we'll raise an error, for
// most use cases this will prevent looping over the same invoice.
func (inv *Invoice) Correct(opts ...cbc.Option) error {
	o := new(CorrectionOptions)
	for _, row := range opts {
		row(o)
	}

	// If we have a raw json object, this will override any of the other options
	if len(o.data) > 0 {
		if err := json.Unmarshal(o.data, o); err != nil {
			return fmt.Errorf("failed to unmarshal correction options: %w", err)
		}
	}

	if o.Credit && o.Debit {
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
	pre := &Preceding{
		UUID:             inv.UUID,
		Series:           inv.Series,
		Code:             inv.Code,
		IssueDate:        inv.IssueDate.Clone(),
		Reason:           o.Reason,
		Corrections:      o.Corrections,
		CorrectionMethod: o.CorrectionMethod,
	}
	inv.UUID = nil
	inv.Series = ""
	inv.Code = ""
	if o.IssueDate != nil {
		inv.IssueDate = *o.IssueDate
	} else {
		inv.IssueDate = cal.Today()
	}

	// Take the regime def to figure out what needs to be copied
	if o.Credit {
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
	} else if o.Debit {
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
			for _, row := range o.Stamps {
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
