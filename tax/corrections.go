package tax

import (
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/validation"
)

// CorrectionSet defines a set of correction definitions for
// a selection of schemas.
type CorrectionSet []*CorrectionDefinition

// CorrectionDefinition contains details about what can be defined in .
type CorrectionDefinition struct {
	// Partial or complete schema URL for the document type supported by correction.
	Schema string `json:"schema" jsonschema:"title=Schema"`
	// The types of sub-documents supported by the regime
	Types []cbc.Key `json:"types,omitempty" jsonschema:"title=Types"`
	// Extension keys that can be included
	Extensions []cbc.Key `json:"extensions,omitempty" jsonschema:"title=Extensions"`
	// ReasonRequired when true implies that a reason must be provided
	ReasonRequired bool `json:"reason_required,omitempty" jsonschema:"title=Reason Required"`
	// Stamps that must be copied from the preceding document.
	Stamps []cbc.Key `json:"stamps,omitempty" jsonschema:"title=Stamps"`
	// Copy tax from the preceding document to the document ref.
	CopyTax bool `json:"copy_tax,omitempty" jsonschema:"title=Copy Tax Totals"`
}

// Def provides the correction definition in the set for the
// schema provided.
func (cs CorrectionSet) Def(schema string) *CorrectionDefinition {
	if cs == nil {
		return nil
	}
	for _, cd := range cs {
		if strings.HasSuffix(schema, cd.Schema) {
			return cd
		}
	}
	return nil
}

// Merge combines two correction definitions into a single one.
func (cd *CorrectionDefinition) Merge(other *CorrectionDefinition) *CorrectionDefinition {
	if cd == nil {
		return other
	}
	if other == nil {
		return cd
	}
	if cd.Schema != other.Schema {
		return cd
	}
	if other.CopyTax {
		cd.CopyTax = other.CopyTax
	}
	cd = &CorrectionDefinition{
		Schema:         cd.Schema,
		Types:          append(cd.Types, other.Types...),
		Extensions:     append(cd.Extensions, other.Extensions...),
		ReasonRequired: cd.ReasonRequired || other.ReasonRequired,
		Stamps:         append(cd.Stamps, other.Stamps...),
		CopyTax:        cd.CopyTax,
	}
	return cd
}

// HasType returns true if the correction definition has a type that matches the one provided.
func (cd *CorrectionDefinition) HasType(t cbc.Key) bool {
	if cd == nil {
		return false // no preceding definitions
	}
	return t.In(cd.Types...)
}

// HasExtension returns true if the correction definition has the change key provided.
func (cd *CorrectionDefinition) HasExtension(key cbc.Key) bool {
	if cd == nil {
		return false // no correction definitions
	}
	return key.In(cd.Extensions...)
}

// Validate ensures the key definition looks correct in the context of the regime.
func (cd *CorrectionDefinition) Validate() error {
	err := validation.ValidateStruct(cd,
		validation.Field(&cd.Schema, validation.Required),
		validation.Field(&cd.Types),
		validation.Field(&cd.Stamps),
		validation.Field(&cd.Extensions),
	)
	return err
}
