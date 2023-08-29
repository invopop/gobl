package tax

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
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

	// Name of the scenario for further information.
	Name i18n.String `json:"name,omitempty" jsonschema:"title=Name"`

	// A note to be added to the document if the scenario is applied.
	Note *cbc.Note `json:"note,omitempty" jsonschema:"title=Note"`

	// Codes is used to define additional codes for regime specific
	// situations.
	Codes cbc.CodeMap `json:"codes,omitempty" jsonschema:"title=Codes"`

	// Any additional local meta data that may be useful in integrations.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// ScenarioSummary is the result after running through a set of
// scenarios and determining which combinations of Notes and Meta
// are viable.
type ScenarioSummary struct {
	Notes []*cbc.Note
	Codes cbc.CodeMap
	Meta  cbc.Meta
}

// Validate checks the scenario set for errors.
func (ss *ScenarioSet) Validate() error {
	err := validation.ValidateStruct(ss,
		validation.Field(&ss.Schema, validation.Required),
		validation.Field(&ss.List, validation.Required),
	)
	return err
}

// SummaryFor returns a summary by applying the scenarios to the
// supplied document.
func (ss *ScenarioSet) SummaryFor(docType cbc.Key, docTags []cbc.Key) *ScenarioSummary {
	summary := &ScenarioSummary{
		Notes: make([]*cbc.Note, 0),
		Codes: make(cbc.CodeMap),
		Meta:  make(cbc.Meta),
	}
	for _, s := range ss.List {
		if s.match(docType, docTags) {
			if s.Note != nil {
				summary.Notes = append(summary.Notes, s.Note)
			}
			for k, v := range s.Codes {
				summary.Codes[k] = v
			}
			for k, v := range s.Meta {
				summary.Meta[k] = v
			}
		}
	}
	return summary
}

// match checks if the scenario has a matching doc type or set of tags.
// Empty types or tags in the scenario implies that all values are valid.
func (s *Scenario) match(docType cbc.Key, docTags []cbc.Key) bool {
	if len(s.Types) > 0 {
		if !s.hasType(docType) {
			return false
		}
	}
	if len(s.Tags) > 0 {
		if !s.hasTags(docTags) {
			return false
		}
	}
	return true
}

// hasType returns true if the scenario has the specified document type.
func (s *Scenario) hasType(docType cbc.Key) bool {
	return docType.In(s.Types...)
}

// hasTags returns true if the the provided document tags is a subset of the
// scenarios tags.
func (s *Scenario) hasTags(docTags []cbc.Key) bool {
	if len(s.Tags) > 0 {
		for _, t := range s.Tags {
			if !t.In(docTags...) {
				return false
			}
		}
		return true
	}
	return false
}

// Validate checks the scenario for errors.
func (s *Scenario) Validate() error {
	err := validation.ValidateStruct(s,
		validation.Field(&s.Types),
		validation.Field(&s.Tags),
		validation.Field(&s.Name),
		validation.Field(&s.Note),
		validation.Field(&s.Codes),
		validation.Field(&s.Meta),
	)
	return err
}
