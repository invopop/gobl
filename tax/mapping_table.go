package tax

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/validation"
)

// MappingTable is used to define the mapping between semantic specifications of document
// types and GOBL fields.
type MappingTable struct {
	// Unique identifier for the mapping table to use for combining sources.
	Key cbc.Key `json:"key" jsonschema:"title=Key"`

	// Name of the mapping table
	Name string `json:"name" jsonschema:"title=Name"`

	// Description provides additional information about the mapping table, Markdown syntax
	// is recommended.
	Description string `json:"description,omitempty" jsonschema:"title=Description"`

	// Terms defines the list of mapping rows.
	Terms []*MappingTerm `json:"terms" jsonschema:"title=Terms"`
}

// MappingTerm contains the details of a mapping between a concept or "term" and
// a single or set of paths in the GOBL document.
type MappingTerm struct {
	// ID contains a single code that will be used to map to the defined paths.
	// If multiple concepts are required, they should be split into separate rows.
	ID cbc.Code `json:"id" jsonschema:"title=ID"`

	// Business term name.
	Name string `json:"name" jsonschema:"title=Name"`

	// Paths represents an array of JSONPaths to the GOBL fields inside an envelope.
	Paths []string `json:"paths" jsonschema:"title=Paths"`

	// Notes provides additional information about the mapping.
	Notes string `json:"notes,omitempty" jsonschema:"title=Notes"`

	// Sub-terms provides additional mappings within the same group.
	Terms []*MappingTerm `json:"terms,omitempty" jsonschema:"title=Terms"`
}

// Validate ensures the table is correctly defined.
func (mt *MappingTable) Validate() error {
	return validation.ValidateStruct(mt,
		validation.Field(&mt.Key, validation.Required),
		validation.Field(&mt.Name, validation.Required),
		validation.Field(&mt.Terms, validation.Required),
	)
}

// Validate ensures that the term definition looks correct.
func (mt *MappingTerm) Validate() error {
	return validation.ValidateStruct(mt,
		validation.Field(&mt.ID, validation.Required),
		validation.Field(&mt.Name, validation.Required),
		validation.Field(&mt.Paths),
		validation.Field(&mt.Terms),
	)
}
