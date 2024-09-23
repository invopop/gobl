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
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
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
	// Series to assign to the new corrective invoice.
	Series cbc.Code `json:"series,omitempty" jsonschema:"title=Series"`
	// Stamps of the previous document to include in the preceding data.
	Stamps []*head.Stamp `json:"stamps,omitempty" jsonschema:"title=Stamps"`
	// Human readable reason for the corrective operation.
	Reason string `json:"reason,omitempty" jsonschema:"title=Reason"`
	// Extensions for region specific requirements.
	Ext tax.Extensions `json:"ext,omitempty" jsonschema:"title=Extensions"`

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

// WithSeries assigns a new series to the corrective document.
func WithSeries(value cbc.Code) schema.Option {
	return func(o interface{}) {
		opts := o.(*CorrectionOptions)
		opts.Series = value
	}
}

// WithReason allows a reason to be provided for the corrective operation.
func WithReason(reason string) schema.Option {
	return func(o interface{}) {
		opts := o.(*CorrectionOptions)
		opts.Reason = reason
	}
}

// WithExtension adds a specific extension combination to the invoice's preceding data,
// can be called multiple times.
func WithExtension(key cbc.Key, code tax.ExtValue) schema.Option {
	return func(o interface{}) {
		opts := o.(*CorrectionOptions)
		if opts.Ext == nil {
			opts.Ext = make(tax.Extensions)
		}
		opts.Ext[key] = code
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
	js := new(jsonschema.Schema)

	// try to load the pre-generated schema, this is just way more efficient
	// than trying to generate the configuration options manually.
	data, err := data.Content.ReadFile("schemas/bill/correction-options.json")
	if err != nil {
		return nil, fmt.Errorf("loading schema option data: %w", err)
	}
	if err := json.Unmarshal(data, js); err != nil {
		return nil, fmt.Errorf("unmarshalling options schema: %w", err)
	}

	// Add our tax country code to the schema ID
	code := strings.ToLower(inv.GetRegime().String())
	id := fmt.Sprintf("%s?tax_regime=%s", js.ID.String(), code)
	js.ID = jsonschema.ID(id)
	js.Comments = fmt.Sprintf("Generated dynamically for %s", code)

	cos := js.Definitions["CorrectionOptions"]

	// Always recommend the series
	recommended := []string{"series"}

	// Try to load the regime and its correction definition for the document
	// type if there is one defined.
	r := inv.RegimeDef()
	if r == nil {
		return js, nil
	}
	cd := inv.correctionDef()
	if cd == nil {
		return js, nil
	}

	if len(cd.Types) > 0 {
		if ps, ok := cos.Properties.Get("type"); ok {
			ps.Default = cd.Types[0].String() // pick first one
			ps.OneOf = make([]*jsonschema.Schema, len(cd.Types))
			for i, v := range cd.Types {
				kd := cbc.GetKeyDefinition(v, InvoiceTypes)
				ps.OneOf[i] = &jsonschema.Schema{
					Const:       v.String(),
					Title:       kd.Name.String(),
					Description: kd.Desc.String(),
				}
			}
		}
	}

	if inv.Series != "" {
		if ps, ok := cos.Properties.Get("series"); ok {
			ps.Default = inv.Series // copy series from invoice
		}
	}

	// Try to add all the specific options for the extensions
	if len(cd.Extensions) > 0 {
		if ext, ok := cos.Properties.Get("ext"); ok {
			ext.Ref = "" // remove the ref
			ext.Type = "object"
			ext.Properties = jsonschema.NewProperties()
			rcmd := make([]string, 0)
			for _, pk := range cd.Extensions {
				re := tax.ExtensionForKey(pk)
				if re == nil {
					continue
				}
				prop := &jsonschema.Schema{
					Title: re.Name.String(),
					Type:  "string",
				}
				if !re.Desc.IsEmpty() {
					prop.Description = re.Desc.String()
				}
				var oneOf []*jsonschema.Schema
				if len(re.Values) > 0 {
					oneOf = make([]*jsonschema.Schema, 0, len(re.Values))
					for _, c := range re.Values {
						ci := &jsonschema.Schema{
							Const: c.Value,
							Title: c.Name.String(),
						}
						if len(c.Desc) > 0 {
							ci.Description = c.Desc.String()
						}
						oneOf = append(oneOf, ci)
					}
				}
				if oneOf != nil {
					prop.OneOf = oneOf
				}
				ext.Properties.Set(pk.String(), prop)
				rcmd = append(rcmd, pk.String())
			}
			if len(rcmd) > 0 {
				// Add the "recommended" extensions as extras
				ext.Extras = map[string]any{
					schema.Recommended: rcmd,
				}
			}
		}
		recommended = append(recommended, "ext")
	} else {
		// Remove extensions, they're not needed if not defined
		cos.Properties.Delete("ext")
	}

	if cd.ReasonRequired {
		cos.Required = append(cos.Required, "reason")
	}

	if len(recommended) > 0 {
		cos.Extras = map[string]any{
			schema.Recommended: recommended,
		}
	}

	return js, nil
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

	// Copy and prepare the basic fields
	pre := &org.DocumentRef{
		Identify:  uuid.Identify{UUID: inv.UUID},
		Type:      inv.Type,
		Series:    cbc.Code(inv.Series),
		Code:      cbc.Code(inv.Code),
		IssueDate: inv.IssueDate.Clone(),
		Reason:    o.Reason,
		Ext:       o.Ext,
	}
	inv.UUID = ""
	inv.Type = o.Type
	if o.Series != "" {
		inv.Series = o.Series
	}
	inv.Code = ""
	if o.IssueDate != nil {
		inv.IssueDate = *o.IssueDate
	} else {
		inv.IssueDate = cal.Today()
	}

	cd := inv.correctionDef()
	if err := inv.validatePrecedingData(o, cd, pre); err != nil {
		return err
	}

	// Replace all previous preceding data
	inv.Preceding = []*org.DocumentRef{pre}

	// Running a Calculate feels a bit out of place, but not performing
	// this operation on the corrected invoice results in potentially
	// conflicting or incomplete data.
	return inv.Calculate()
}

// correctionDef tries to determine a final correction definition
// by merge potentially multiple sources. The results include
// a key that can be used to identify the definition.
func (inv *Invoice) correctionDef() *tax.CorrectionDefinition {
	cd := &tax.CorrectionDefinition{
		Schema: ShortSchemaInvoice,
	}

	r := inv.RegimeDef()
	if r != nil {
		cd = cd.Merge(r.Corrections.Def(ShortSchemaInvoice))
	}
	for _, a := range inv.GetAddonDefs() {
		cd = cd.Merge(a.Corrections.Def(ShortSchemaInvoice))
	}

	return cd
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

func (inv *Invoice) validatePrecedingData(o *CorrectionOptions, cd *tax.CorrectionDefinition, pre *org.DocumentRef) error {
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

	if cd.ReasonRequired && pre.Reason == "" {
		return errors.New("missing corrective reason")
	}

	return nil
}
