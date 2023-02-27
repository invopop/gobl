package tax

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/validation"
)

// Scheme contains the definition of a scheme that belongs to a region and can be used
// to simplify validation processes for document contents.
type Scheme struct {
	/* Basic details of Scheme */

	// Key used to identify this scheme
	Key cbc.Key `json:"key" jsonschema:"title=Key"`
	// Name of this scheme.
	Name i18n.String `json:"name" jsonschema:"title=Name"`
	// Human details describing what this scheme is used for.
	Description i18n.String `json:"description,omitempty" jsonschema:"title=Description"`

	/* Conditions */

	// Set of invoice types supported by this scheme.
	InvoiceTypes []cbc.Key `json:"invoice_types,omitempty" jsonschema:"title=Invoice Types"`

	// Array of schemes that this scheme can appear under.
	Parents []cbc.Key `json:"parents,omitempty" jsonschema:"title=Parent Schemes"`

	/* Requirements */

	// List of tax category codes that can be used when this scheme is
	// applied.
	Categories []cbc.Code `json:"categories,omitempty" jsonschema:"title=Category Codes"`

	/* Outputs */

	// Note defines a message that should be added to a document
	// when this scheme is used.
	Note *cbc.Note `json:"note,omitempty" jsonschema:"title=Note"`

	// Any additional data that is relevant for generating documents in a
	// regime, but not meaningful inside GOBL.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Validate ensures the tax scheme details look valid.
func (s *Scheme) Validate() error {
	return validation.ValidateStruct(s,
		validation.Field(&s.Key, validation.Required),
		validation.Field(&s.Name, validation.Required),
		validation.Field(&s.Description),
		validation.Field(&s.Categories),
		validation.Field(&s.Note),
	)
}
