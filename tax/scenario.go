package tax

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

// ScenarioSet is a collection of tax scenarios for a given schema that can be used to
// determine special codes or notes that need to be included in the final document.
type ScenarioSet struct {
	// Partial or complete schema URL for the document type
	Schema string `json:"schema" jsonschema:"title=Schema"`

	// List of scenarios for the schema
	List []*Scenario `json:"list" jsonschema:"title=List"`
}

// ScenarioDocument is used to determine if scenarios can be applied to a document.
type ScenarioDocument interface {
	// GetType returns a type associated with the document.
	GetType() cbc.Key
	// GetTags returns a list of the tags used in the document.
	GetTags() []cbc.Key
	// GetExtensions an array of extensions that used in the document.
	GetExtensions() []Extensions
	// GetTaxCategories returns a list of tax categories used in the document's lines,
	// charges, and discounts.
	GetTaxCategories() []cbc.Code
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

	// Description of the scenario for documentation purposes.
	Desc i18n.String `json:"desc,omitempty" jsonschema:"title=Description"`

	/* Filters */

	// Type of document, if present.
	Types []cbc.Key `json:"type,omitempty" jsonschema:"title=Type"`

	// Array of tags that have been applied to the document.
	Tags []cbc.Key `json:"tags,omitempty" jsonschema:"title=Tags"`

	// Categories is an optional list of tax category codes that acts as a filter.
	// When set, at least one of the specified categories must be present in the
	// document's line taxes for the scenario to match.
	Categories []cbc.Code `json:"cat,omitempty" jsonschema:"title=Tax Categories"`

	// Extension key that must be present in the document.
	ExtKey cbc.Key `json:"ext_key,omitempty" jsonschema:"title=Extension Key"`

	// Extension code that along side the key must be present for a match
	// to happen. This cannot be used without an `cbc.Code`. The value will
	// be copied to the note code if needed.
	ExtCode cbc.Code `json:"ext_code,omitempty" jsonschema:"title=Extension Code"`

	// Filter defines a custom filter method for when the regular basic filters
	// are not sufficient.
	Filter func(doc any) bool `json:"-"`

	/* Outputs */

	// A note to be added to the document if the scenario is applied.
	Note *Note `json:"note,omitempty" jsonschema:"title=Note"`

	// Codes is used to define additional codes for regime specific
	// situations.
	Codes cbc.CodeMap `json:"codes,omitempty" jsonschema:"title=Codes"`

	// Ext represents a set of tax extensions that should be applied to
	// the document in the appropriate "tax" context.
	Ext Extensions `json:"ext,omitzero" jsonschema:"title=Extensions"`
}

// ScenarioSummary is the result after running through a set of
// scenarios and determining which combinations of Notes, Codes, Meta,
// and extensions are viable.
type ScenarioSummary struct {
	Notes []*Note
	Codes cbc.CodeMap
	Ext   Extensions
}

// NewScenarioSet creates a new scenario set with the given schema.
func NewScenarioSet(schema string) *ScenarioSet {
	return &ScenarioSet{
		Schema: schema,
		List:   make([]*Scenario, 0),
	}
}

func scenarioSetRules() *rules.Set {
	return rules.For(new(ScenarioSet),
		rules.Field("schema",
			rules.Assert("01", "schema is required", is.Present),
		),
		rules.Field("list",
			rules.Assert("02", "at least one scenario is required", is.Present),
		),
	)
}

// Merge appends the scenarios from the other set to the current set.
func (ss *ScenarioSet) Merge(other []*ScenarioSet) {
	for _, os := range other {
		if os.Schema != ss.Schema {
			return
		}
		ss.List = append(ss.List, os.List...)
	}
}

// ExtensionKeys extracts all the possible extension keys that could be applied to a
// document.
func (ss *ScenarioSet) ExtensionKeys() []cbc.Key {
	keys := make([]cbc.Key, 0)
	for _, row := range ss.List {
		for _, k := range row.Ext.Keys() {
			if !k.In(keys...) {
				keys = append(keys, k)
			}
		}
	}
	return keys
}

// Notes extracts all the possible notes that could be applied to a document.
func (ss *ScenarioSet) Notes() []*Note {
	notes := make([]*Note, 0)
	for _, row := range ss.List {
		if row.Note != nil {
			notes = append(notes, row.Note)
		}
	}
	return notes
}

// SummaryFor returns a summary by applying the scenarios to the
// supplied document.
func (ss *ScenarioSet) SummaryFor(doc ScenarioDocument) *ScenarioSummary {
	summary := &ScenarioSummary{
		Notes: make([]*Note, 0),
		Codes: make(cbc.CodeMap),
		Ext:   MakeExtensions(),
	}
	for _, s := range ss.List {
		if s.match(doc) {
			if s.Note != nil {
				summary.addNote(s.Note)
			}
			for k, v := range s.Codes {
				summary.Codes[k] = v
			}
			summary.Ext = summary.Ext.Merge(s.Ext)
		}
	}
	return summary
}

func (ss *ScenarioSummary) addNote(note *Note) {
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
func (s *Scenario) match(doc ScenarioDocument) bool {
	if len(s.Types) > 0 {
		if !s.hasType(doc.GetType()) {
			return false
		}
	}
	if len(s.Tags) > 0 {
		if !s.hasTags(doc.GetTags()) {
			return false
		}
	}
	if len(s.Categories) > 0 {
		if !s.hasCategories(doc.GetTaxCategories()) {
			return false
		}
	}
	if s.ExtKey != cbc.KeyEmpty {
		// For extensions we need to find a complete match
		// and reject if none found. We intentionally don't try
		// to combine extensions from the document.
		for _, ext := range doc.GetExtensions() {
			if !ext.Has(s.ExtKey) {
				continue // try next extension
			}
			v := ext.Get(s.ExtKey)
			if s.ExtCode != "" {
				if v == s.ExtCode {
					return true
				}
			} else {
				return true
			}
		}
		return false
	}
	if s.Filter != nil {
		if !s.Filter(doc) {
			return false
		}
	}
	return true
}

// hasType returns true if the scenario has the specified document type.
func (s *Scenario) hasType(docType cbc.Key) bool {
	return docType.In(s.Types...)
}

// hasCategories returns true if at least one of the scenario's categories
// is present in the document's categories.
func (s *Scenario) hasCategories(docCats []cbc.Code) bool {
	for _, c := range s.Categories {
		if c.In(docCats...) {
			return true
		}
	}
	return false
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
