package tax

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/validation"
)

// ScenarioSet is a collection of tax scenarios for a given schema that can be used to
// determine special codes or notes that need to be included in the final document.
type ScenarioSet struct {
	// Partial or complete schema URL for the document type
	Schema string `json:"schema" jsonschema:"title=Schema"`
	// List of scenarios for the schema
	List []*Scenario `json:"list" jsonschema:"title=List"`
}

// Scenario is used to describe a tax scenario of a document based on the combination
// of document type and tag used.
type Scenario struct {
	// Type of document, if present.
	Types []cbc.Key `json:"type,omitempty" jsonschema:"title=Type"`

	// Tag that was applied to the document
	Tags []cbc.Key `json:"tags,omitempty" jsonschema:"title=Tag"`

	// A note to be added to the document if the scenario is applied.
	Note *cbc.Note `json:"note,omitempty" jsonschema:"title=Note"`

	// Any additional local meta data that may be useful in integrations.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Validate checks the scenario set for errors.
func (ss *ScenarioSet) Validate() error {
	err := validation.ValidateStruct(ss,
		validation.Field(&ss.Schema, validation.Required),
		validation.Field(&ss.List, validation.Required),
	)
	return err
}

// Validate checks the scenario for errors.
func (s *Scenario) Validate() error {
	err := validation.ValidateStruct(s,
		validation.Field(&s.Types),
		validation.Field(&s.Tags),
		validation.Field(&s.Note),
		validation.Field(&s.Meta),
	)
	return err
}
