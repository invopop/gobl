package tax

import (
	"context"

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
//
// There are effectively two parts to a scenario, the filters that are used to determine
// if the scenario is applicable to a document and the output that is applied or data to
// be used by conversion processes.
type Scenario struct {
	// Name of the scenario for further information.
	Name i18n.String `json:"name,omitempty" jsonschema:"title=Name"`

	/* Filters */

	// Type of document, if present.
	Types []cbc.Key `json:"type,omitempty" jsonschema:"title=Type"`

	// Tag that was applied to the document
	Tags []cbc.Key `json:"tags,omitempty" jsonschema:"title=Tag"`

	// Extension key that must be present in the document.
	ExtKey cbc.Key `json:"ext_key,omitempty" jsonschema:"title=Extension Key"`

	// Extension value that along side the key must be present for a match
	// to happen. This cannot be used without an `ExtKey`. The value will
	// be copied to the note code if needed.
	ExtValue ExtValue `json:"ext_value,omitempty" jsonschema:"title=Extension Value"`

	/* Outputs */

	// A note to be added to the document if the scenario is applied.
	Note *cbc.Note `json:"note,omitempty" jsonschema:"title=Note"`

	// Codes is used to define additional codes for regime specific
	// situations.
	Codes cbc.CodeMap `json:"codes,omitempty" jsonschema:"title=Codes"`

	// Ext represents a set of tax extensions that should be applied to
	// the document in the appropriate "tax" context.
	Ext Extensions `json:"ext,omitempty" jsonschema:"title=Extensions"`
}

// ScenarioSummary is the result after running through a set of
// scenarios and determining which combinations of Notes, Codes, Meta,
// and extensions are viable.
type ScenarioSummary struct {
	Notes []*cbc.Note
	Codes cbc.CodeMap
	Ext   Extensions
}

// ValidateWithContext checks the scenario set for errors.
func (ss *ScenarioSet) ValidateWithContext(ctx context.Context) error {
	err := validation.ValidateStructWithContext(ctx, ss,
		validation.Field(&ss.Schema, validation.Required),
		validation.Field(&ss.List, validation.Required),
	)
	return err
}

// SummaryFor returns a summary by applying the scenarios to the
// supplied document.
func (ss *ScenarioSet) SummaryFor(docType cbc.Key, docTags []cbc.Key, docExt []Extensions) *ScenarioSummary {
	summary := &ScenarioSummary{
		Notes: make([]*cbc.Note, 0),
		Codes: make(cbc.CodeMap),
		Ext:   make(Extensions),
	}
	for _, s := range ss.List {
		if s.match(docType, docTags, docExt) {
			if s.Note != nil {
				summary.addNote(s.Note.WithCode(s.ExtValue.String()))
			}
			for k, v := range s.Codes {
				summary.Codes[k] = v
			}
			for k, v := range s.Ext {
				summary.Ext[k] = v
			}
		}
	}
	return summary
}

func (ss *ScenarioSummary) addNote(note *cbc.Note) {
	for i, n := range ss.Notes {
		if n.SameAs(note) {
			// replace
			ss.Notes[i] = note
			return
		}
	}
	ss.Notes = append(ss.Notes, note)
}

// match checks if the scenario has a matching doc type or set of tags.
// Empty types or tags in the scenario implies that all values are valid.
// The list of extensions can contain duplicate extension maps to make recompilation
// of the array easier.
func (s *Scenario) match(docType cbc.Key, docTags []cbc.Key, docExt []Extensions) bool {
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
	if s.ExtKey != cbc.KeyEmpty {
		// For extensions we need to find a complete match
		// and reject if none found. We intentionally don't try
		// to combine extensions from the document.
		for _, ext := range docExt {
			v, ok := ext[s.ExtKey]
			if !ok {
				continue // try next extension
			}
			if s.ExtValue != "" {
				if v == s.ExtValue {
					return true
				}
			} else {
				return true
			}
		}
		return false
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

// ValidateWithContext checks the scenario for errors, using the regime in the context
// to validate the list of tags.
func (s *Scenario) ValidateWithContext(ctx context.Context) error {
	r := ctx.Value(KeyRegime).(*Regime)
	err := validation.ValidateStructWithContext(ctx, s,
		validation.Field(&s.Types),
		validation.Field(&s.Tags, validation.Each(cbc.InKeyDefs(r.Tags))),
		validation.Field(&s.Name),
		validation.Field(&s.Note),
		validation.Field(&s.Codes),
		validation.Field(&s.Ext),
	)
	return err
}
