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

// ShortSchemaContract is the short schema name for a legal contract.
const ShortSchemaContract = "legal/contract"

func init() {
	schema.Register(schema.GOBL.Add("legal"),
		Contract{},
		Chapter{},
		Section{},
	)
	rules.Register(
		"legal",
		rules.GOBL.Add("LEGAL"),
		contractRules(),
		chapterRules(),
		sectionRules(),
	)
}
