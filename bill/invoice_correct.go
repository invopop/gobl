package bill

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/iancoleman/orderedmap"
	"github.com/invopop/gobl/build"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
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

	// When the new corrective invoice's issue date should be set to.
	IssueDate *cal.Date `json:"issue_date,omitempty" jsonschema:"title=Issue Date"`
	// Stamps of the previous document to include in the preceding data.
	Stamps []*head.Stamp `json:"stamps,omitempty" jsonschema:"title=Stamps"`
	// Credit when true indicates that the corrective document should cancel the previous document.
	Credit bool `json:"credit,omitempty" jsonschema:"title=Credit"`
	// Debit when true indicates that the corrective document should add new items to the previous document.
	Debit bool `json:"debit,omitempty" jsonschema:"title=Debit"`
	// Human readable reason for the corrective operation.
	Reason string `json:"reason,omitempty" jsonschema:"title=Reason"`
	// Correction method as defined by the tax regime.
	Method cbc.Key `json:"method,omitempty" jsonschema:"title=Method"`
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

// WithMethod defines the method used to correct the previous invoice.
func WithMethod(method cbc.Key) schema.Option {
	return func(o interface{}) {
		opts := o.(*CorrectionOptions)
		opts.Method = method
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

// Credit indicates that the corrective operation requires a credit note
// or equivalent.
var Credit schema.Option = func(o interface{}) {
	opts := o.(*CorrectionOptions)
	opts.Credit = true
}

// Debit indicates that the corrective operation is to append
// new items to the previous invoice, usually as a debit note.
var Debit schema.Option = func(o interface{}) {
	opts := o.(*CorrectionOptions)
	opts.Debit = true
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
	data, err := build.Content.ReadFile("schemas/bill/correction-options.json")
	if err != nil {
		return nil, fmt.Errorf("loading schema option data: %w", err)
	}
	if err := json.Unmarshal(data, schema); err != nil {
		return nil, fmt.Errorf("unmarshalling options schema: %w", err)
	}
	schema = schema.Definitions["CorrectionOptions"]

	// Improve the quality of the schema
	schema.Required = append(schema.Required, "credit")

	cd := r.CorrectionDefinitionFor(ShortSchemaInvoice)
	if cd == nil {
		return schema, nil
	}

	if cd.ReasonRequired {
		schema.Required = append(schema.Required, "reason")
	}

	// These methods are quite ugly as the jsonschema was not designed
	// for being able to load documents.
	if len(cd.Methods) > 0 {
		schema.Required = append(schema.Required, "method")
		if prop, ok := schema.Properties.Get("method"); ok {
			ps := prop.(orderedmap.OrderedMap)
			oneOf := make([]*jsonschema.Schema, len(cd.Methods))
			for i, v := range cd.Methods {
				oneOf[i] = &jsonschema.Schema{
					Const: v.Key.String(),
					Title: v.Name.String(),
				}
				if !v.Desc.IsEmpty() {
					oneOf[i].Description = v.Desc.String()
				}
			}
			ps.Set("oneOf", oneOf)
			schema.Properties.Set("method", ps)
		}
	}

	if len(cd.Keys) > 0 {
		schema.Required = append(schema.Required, "changes")
		if prop, ok := schema.Properties.Get("changes"); ok {
			ps := prop.(orderedmap.OrderedMap)
			items, _ := ps.Get("items")
			pi := items.(orderedmap.OrderedMap)

			oneOf := make([]*jsonschema.Schema, len(cd.Keys))
			for i, v := range cd.Keys {
				oneOf[i] = &jsonschema.Schema{
					Const: v.Key.String(),
					Title: v.Name.String(),
				}
				if !v.Desc.IsEmpty() {
					oneOf[i].Description = v.Desc.String()
				}
			}
			pi.Set("oneOf", oneOf)
			ps.Set("items", pi)
			schema.Properties.Set("changes", ps)
		}
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
		UUID:             inv.UUID,
		Series:           inv.Series,
		Code:             inv.Code,
		IssueDate:        inv.IssueDate.Clone(),
		Reason:           o.Reason,
		CorrectionMethod: o.Method,
		Changes:          o.Changes,
	}
	inv.UUID = nil
	inv.Series = ""
	inv.Code = ""
	if o.IssueDate != nil {
		inv.IssueDate = *o.IssueDate
	} else {
		inv.IssueDate = cal.Today()
	}

	cd := r.CorrectionDefinitionFor(ShortSchemaInvoice)

	if err := inv.prepareCorrectionType(o, cd); err != nil {
		return err
	}

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

	if o.Credit && o.Debit {
		return errors.New("cannot use both credit and debit options")
	}
	return nil
}

func (inv *Invoice) prepareCorrectionType(o *CorrectionOptions, cd *tax.CorrectionDefinition) error {
	// Take the regime def to figure out what needs to be copied
	if o.Credit {
		if cd.HasType(InvoiceTypeCreditNote) {
			// regular credit note
			inv.Type = InvoiceTypeCreditNote
		} else if cd.HasType(InvoiceTypeCorrective) {
			// corrective invoice with negative values
			inv.Type = InvoiceTypeCorrective
			inv.Invert()
		} else {
			return errors.New("credit note not supported by regime")
		}
		inv.Payment.ResetAdvances()
	} else if o.Debit {
		if cd.HasType(InvoiceTypeDebitNote) {
			// regular debit note, implies no rows as new ones
			// will be added
			inv.Type = InvoiceTypeDebitNote
			inv.Empty()
		} else {
			return errors.New("debit note not supported by regime")
		}
	} else {
		if cd.HasType(InvoiceTypeCorrective) {
			inv.Type = InvoiceTypeCorrective
		} else {
			return fmt.Errorf("corrective invoice type not supported by regime, try credit or debit")
		}
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

	if len(cd.Methods) > 0 {
		if pre.CorrectionMethod == cbc.KeyEmpty {
			return errors.New("missing correction method")
		}
		if !cd.HasMethod(pre.CorrectionMethod) {
			return fmt.Errorf("invalid correction method: %v", pre.CorrectionMethod)
		}
	}

	if len(cd.Keys) > 0 {
		if len(pre.Changes) == 0 {
			return errors.New("missing changes")
		}
		for _, k := range pre.Changes {
			if !cd.HasKey(k) {
				return fmt.Errorf("invalid change key: '%v'", k)
			}
		}
	}

	if cd.ReasonRequired && pre.Reason == "" {
		return errors.New("missing corrective reason")
	}

	return nil
}
