package bill

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/data"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/jsonschema"
)

// CorrectionOptions defines a structure used to pass configuration options
// to correct a previous invoice. This is made available to make it easier to
// pass options between external services.
type CorrectionOptions struct {
	head.CorrectionOptions

	// The type of corrective invoice to produce.
	Type cbc.Key `json:"type" jsonschema:"title=Type"`
	// When the new corrective invoice's issue date should be set to.
	IssueDate *cal.Date `json:"issue_date,omitempty" jsonschema:"title=Issue Date"`
	// Stamps of the previous document to include in the preceding data.
	Stamps []*head.Stamp `json:"stamps,omitempty" jsonschema:"title=Stamps"`
	// Human readable reason for the corrective operation.
	Reason string `json:"reason,omitempty" jsonschema:"title=Reason"`
	// Changes keys that describe the specific changes according to the tax regime.
	Changes []cbc.Key `json:"changes,omitempty" jsonschema:"title=Changes"`

	// In case we want to use a raw json object as a source of the options.
	data json.RawMessage `json:"-"`
}

// WithOptions takes an already completed CorrectionOptions instance and
// uses this as a base instead of passing individual options. This is useful
// for passing options from an API, developers should use the regular option
// methods.
func WithOptions(opts *CorrectionOptions) schema.Option {
	return func(o interface{}) {
		o2 := o.(*CorrectionOptions)
		*o2 = *opts
	}
}

// WithData expects a raw JSON object that will be marshalled into a
// CorrectionOptions instance and used as the base for the correction.
func WithData(data json.RawMessage) schema.Option {
	return func(o interface{}) {
		opts := o.(*CorrectionOptions)
		opts.data = data
	}
}

// WithStamps provides a configuration option with stamp information
// usually included in the envelope header for a previously generated
// and processed invoice document.
func WithStamps(stamps []*head.Stamp) schema.Option {
	return func(o interface{}) {
		opts := o.(*CorrectionOptions)
		opts.Stamps = stamps
	}
}

// WithReason allows a reason to be provided for the corrective operation.
func WithReason(reason string) schema.Option {
	return func(o interface{}) {
		opts := o.(*CorrectionOptions)
		opts.Reason = reason
	}
}

// WithChanges adds the set of change keys to the invoice's preceding data,
// can be called multiple times.
func WithChanges(changes ...cbc.Key) schema.Option {
	return func(o interface{}) {
		opts := o.(*CorrectionOptions)
		opts.Changes = append(opts.Changes, changes...)
	}
}

// WithIssueDate can be used to override the issue date of the corrective invoice
// produced.
func WithIssueDate(date cal.Date) schema.Option {
	return func(o interface{}) {
		opts := o.(*CorrectionOptions)
		opts.IssueDate = &date
	}
}

// Corrective is used for creating corrective or rectified invoices
// that completely replace a previous document.
var Corrective schema.Option = func(o interface{}) {
	opts := o.(*CorrectionOptions)
	opts.Type = InvoiceTypeCorrective
}

// Credit indicates that the corrective operation requires a credit note
// or equivalent.
var Credit schema.Option = func(o interface{}) {
	opts := o.(*CorrectionOptions)
	opts.Type = InvoiceTypeCreditNote
}

// Debit indicates that the corrective operation is to append
// new items to the previous invoice, usually as a debit note.
var Debit schema.Option = func(o interface{}) {
	opts := o.(*CorrectionOptions)
	opts.Type = InvoiceTypeDebitNote
}

// CorrectionOptionsSchema provides a dynamic JSON schema of the options
// that can be used on the invoice in order to correct it. Data is
// extracted from the tax regime associated with the supplier.
func (inv *Invoice) CorrectionOptionsSchema() (interface{}, error) {
	r := taxRegimeFor(inv.Supplier)
	if r == nil {
		return nil, nil
	}

	schema := new(jsonschema.Schema)

	// try to load the pre-generated schema, this is just way more efficient
	// than trying to generate the configuration options manually.
	data, err := data.Content.ReadFile("schemas/bill/correction-options.json")
	if err != nil {
		return nil, fmt.Errorf("loading schema option data: %w", err)
	}
	if err := json.Unmarshal(data, schema); err != nil {
		return nil, fmt.Errorf("unmarshalling options schema: %w", err)
	}

	// Add our regime to the schema ID
	code := strings.ToLower(r.Code().String())
	id := fmt.Sprintf("%s?tax_regime=%s", schema.ID.String(), code)
	schema.ID = jsonschema.ID(id)
	schema.Comments = fmt.Sprintf("Generated dynamically for %s", code)

	cos := schema.Definitions["CorrectionOptions"]

	cd := r.CorrectionDefinitionFor(ShortSchemaInvoice)
	if cd == nil {
		return schema, nil
	}

	if len(cd.Types) > 0 {
		if ps, ok := cos.Properties.Get("type"); ok {
			ps.OneOf = make([]*jsonschema.Schema, len(cd.Types))
			for i, v := range cd.Types {
				ps.OneOf[i] = &jsonschema.Schema{
					Const: v.String(),
					Title: v.String(),
				}
			}
		}
	}

	if len(cd.Changes) > 0 {
		cos.Required = append(cos.Required, "changes")
		if ps, ok := cos.Properties.Get("changes"); ok {
			items := ps.Items
			items.OneOf = make([]*jsonschema.Schema, len(cd.Changes))
			for i, v := range cd.Changes {
				items.OneOf[i] = &jsonschema.Schema{
					Const: v.Key.String(),
					Title: v.Name.String(),
				}
				if !v.Desc.IsEmpty() {
					items.OneOf[i].Description = v.Desc.String()
				}
			}
		}
	}

	if cd.ReasonRequired {
		cos.Required = append(cos.Required, "reason")
	}

	return schema, nil
}

// Correct moves key fields of the current invoice to the preceding
// structure and performs any regime specific actions defined by the
// regime's configuration.
// If the existing document doesn't have a code, we'll raise an error, for
// most use cases this will prevent looping over the same invoice.
func (inv *Invoice) Correct(opts ...schema.Option) error {
	o := new(CorrectionOptions)
	if err := prepareCorrectionOptions(o, opts...); err != nil {
		return err
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
		UUID:      inv.UUID,
		Type:      inv.Type,
		Series:    inv.Series,
		Code:      inv.Code,
		IssueDate: inv.IssueDate.Clone(),
		Reason:    o.Reason,
		Changes:   o.Changes,
	}
	inv.UUID = nil
	inv.Type = o.Type
	inv.Series = ""
	inv.Code = ""
	if o.IssueDate != nil {
		inv.IssueDate = *o.IssueDate
	} else {
		inv.IssueDate = cal.Today()
	}

	cd := r.CorrectionDefinitionFor(ShortSchemaInvoice)

	if err := inv.validatePrecedingData(o, cd, pre); err != nil {
		return err
	}

	// Replace all previous preceding data
	inv.Preceding = []*Preceding{pre}

	// Running a Calculate feels a bit out of place, but not performing
	// this operation on the corrected invoice results in potentially
	// conflicting or incomplete data.
	return inv.Calculate()
}

func prepareCorrectionOptions(o *CorrectionOptions, opts ...schema.Option) error {
	for _, row := range opts {
		row(o)
	}

	// Copy over the stamps from the previous header
	if o.Head != nil && len(o.Head.Stamps) > 0 {
		o.Stamps = append(o.Stamps, o.Head.Stamps...)
	}

	// If we have a raw json object, this will override any of the other options
	if len(o.data) > 0 {
		if err := json.Unmarshal(o.data, o); err != nil {
			return fmt.Errorf("failed to unmarshal correction options: %w", err)
		}
	}

	if o.Type == cbc.KeyEmpty {
		return errors.New("missing correction type")
	}

	return nil
}

func (inv *Invoice) validatePrecedingData(o *CorrectionOptions, cd *tax.CorrectionDefinition, pre *Preceding) error {
	if cd == nil {
		return nil
	}
	for _, k := range cd.Stamps {
		var s *head.Stamp
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

	if !o.Type.In(cd.Types...) {
		return fmt.Errorf("invalid correction type: %v", o.Type.String())
	}

	if len(cd.Changes) > 0 {
		if len(pre.Changes) == 0 {
			return errors.New("missing correction changes")
		}
		for _, k := range pre.Changes {
			if !cd.HasChange(k) {
				return fmt.Errorf("invalid correction change key: '%v'", k)
			}
		}
	}

	if cd.ReasonRequired && pre.Reason == "" {
		return errors.New("missing corrective reason")
	}

	return nil
}
