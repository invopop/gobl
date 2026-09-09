// Package legal provides models for structured, signable legal document bodies.
//
// The package is experimental. A Contract currently describes the canonical
// contents and hierarchy of a document; it does not by itself model formation,
// party authority, governing law, or the lifecycle of an agreement.
package legal

import (
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/schema"
)

// Short schema names for legal document payloads.
const (
	ShortSchemaContract = "legal/contract"
	ShortSchemaAssent   = "legal/assent"
	ShortSchemaAnalysis = "legal/analysis"
)

func init() {
	schema.Register(schema.GOBL.Add("legal"),
		Analysis{},
		Annotation{},
		Assent{},
		Contract{},
		Chapter{},
		Definition{},
		DisputeResolution{},
		Effect{},
		Execution{},
		GoverningLaw{},
		Party{},
		Recital{},
		Reference{},
		Resource{},
		Section{},
		Signatory{},
		SignatureRequirement{},
	)
	rules.Register(
		"legal",
		rules.GOBL.Add("LEGAL"),
		contractRules(),
		chapterRules(),
		sectionRules(),
		partyRules(),
		signatoryRules(),
		recitalRules(),
		definitionRules(),
		effectRules(),
		governingLawRules(),
		disputeResolutionRules(),
		executionRules(),
		signatureRequirementRules(),
		referenceRules(),
		resourceRules(),
		assentRules(),
		analysisRules(),
		annotationRules(),
	)
}
