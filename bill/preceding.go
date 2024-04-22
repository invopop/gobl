package bill

import (
	"context"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
)

// Preceding allows for information to be provided about a previous invoice that this one
// will replace, subtract from, or add to. If this is used, the invoice type code will most likely need
// to be set to `corrective`, `credit-note`, or similar.
type Preceding struct {
	// Preceding document's UUID.
	UUID uuid.UUID `json:"uuid,omitempty" jsonschema:"title=UUID"`
	// Type of the preceding document
	Type cbc.Key `json:"type,omitempty" jsonschema:"title=Type"`
	// Series identification code
	Series string `json:"series,omitempty" jsonschema:"title=Series"`
	// Code of the previous document.
	Code string `json:"code" jsonschema:"title=Code"`
	// The issue date of the previous document.
	IssueDate *cal.Date `json:"issue_date,omitempty" jsonschema:"title=Issue Date"`
	// Human readable description on why the preceding invoice is being replaced.
	Reason string `json:"reason,omitempty" jsonschema:"title=Reason"`
	// Seals of approval from other organisations that may need to be listed.
	Stamps []*head.Stamp `json:"stamps,omitempty" jsonschema:"title=Stamps"`
	// Tax period in which the previous invoice had an effect required by some tax regimes and formats.
	Period *cal.Period `json:"period,omitempty" jsonschema:"title=Period"`
	// Extensions for region specific requirements.
	Ext tax.Extensions `json:"ext,omitempty" jsonschema:"title=Ext"`
	// Additional semi-structured data that may be useful in specific regions
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Validate ensures the preceding details look okay
func (p *Preceding) Validate() error {
	return p.ValidateWithContext(context.Background())
}

// Calculate tries to normalize the preceding data
func (p *Preceding) Calculate() error {
	if p == nil {
		return nil
	}
	p.Stamps = head.NormalizeStamps(p.Stamps)
	p.Ext = tax.NormalizeExtensions(p.Ext)
	return nil
}

// ValidateWithContext ensures the preceding details look okay
func (p *Preceding) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, p,
		validation.Field(&p.UUID),
		validation.Field(&p.Type),
		validation.Field(&p.Series),
		validation.Field(&p.Code, validation.Required),
		validation.Field(&p.IssueDate, cal.DateNotZero()),
		validation.Field(&p.Stamps),
		validation.Field(&p.Period),
		validation.Field(&p.Ext),
		validation.Field(&p.Meta),
	)
}

// JSONSchemaExtend extends the schema with additional property details
func (Preceding) JSONSchemaExtend(schema *jsonschema.Schema) {
	props := schema.Properties
	if prop, ok := props.Get("series"); ok {
		prop.Pattern = InvoiceCodePattern
	}
	if prop, ok := props.Get("code"); ok {
		prop.Pattern = InvoiceCodePattern
	}
	// Extend type list
	if its, ok := props.Get("type"); ok {
		its.OneOf = make([]*jsonschema.Schema, len(InvoiceTypes))
		for i, kd := range InvoiceTypes {
			its.OneOf[i] = &jsonschema.Schema{
				Const:       kd.Key.String(),
				Title:       kd.Name.String(),
				Description: kd.Desc.String(),
			}
		}
	}
}
