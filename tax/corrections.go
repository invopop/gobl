package tax

import (
	"slices"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

// CorrectionNormalize is a callback invoked during the Correct() process
// to allow addons and regimes to perform custom normalization. The document
// is passed as `any` to avoid circular imports; implementations should
// type-assert to the concrete type (e.g., *bill.Invoice).
// When called, the document's Preceding field is already set and the
// correction options are available via the document's accessor method.
type CorrectionNormalize func(doc any)

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
	// Normalize is an optional callback invoked during Correct() to allow
	// addon/regime-specific logic to route extensions between the document
	// and the preceding reference.
	Normalize CorrectionNormalize `json:"-"`
}

func correctionDefinitionRules() *rules.Set {
	return rules.For(new(CorrectionDefinition),
		rules.Field("schema",
			rules.Assert("01", "schema is required", is.Present),
		),
	)

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

// Merge combines two correction definitions into a new definition without
// mutating either input.
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
	// Chain normalizers so both run in sequence.
	var norm CorrectionNormalize
	switch {
	case cd.Normalize != nil && other.Normalize != nil:
		first, second := cd.Normalize, other.Normalize
		norm = func(doc any) {
			first(doc)
			second(doc)
		}
	case other.Normalize != nil:
		norm = other.Normalize
	default:
		norm = cd.Normalize
	}
	return &CorrectionDefinition{
		Schema:         cd.Schema,
		Types:          slices.Concat(cd.Types, other.Types),
		Extensions:     slices.Concat(cd.Extensions, other.Extensions),
		ReasonRequired: cd.ReasonRequired || other.ReasonRequired,
		Stamps:         slices.Concat(cd.Stamps, other.Stamps),
		CopyTax:        cd.CopyTax || other.CopyTax,
		Normalize:      norm,
	}
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
