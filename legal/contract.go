package legal

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/uuid"
)

// Contract kind keys classify the legal effect intended for a document.
const (
	ContractKindAgreement cbc.Key = "agreement"
	ContractKindDeed      cbc.Key = "deed"
	ContractKindAmendment cbc.Key = "amendment"
	ContractKindTerms     cbc.Key = "terms"
)

// Contract represents the structured body of a legal document.
type Contract struct {
	uuid.Identify

	// Agreement identifies the enduring legal relationship. Each amendment or
	// restatement has its own UUID while retaining this agreement UUID.
	Agreement uuid.UUID `json:"agreement" jsonschema:"title=Agreement UUID"`
	// Kind describes the intended legal form of this document.
	Kind cbc.Key `json:"kind,omitempty" jsonschema:"title=Contract Kind"`
	// Human or system assigned identifier for this contract.
	Code cbc.Code `json:"code,omitempty" jsonschema:"title=Code"`
	// Version label supplied by the author or contract management system.
	Version cbc.Code `json:"version,omitempty" jsonschema:"title=Version"`
	// Primary language of the authoritative text.
	Language i18n.Lang `json:"language" jsonschema:"title=Language"`
	// Title of the document
	Title string `json:"title" jsonschema:"title=Title"`
	// Sub-title
	SubTitle string `json:"sub_title,omitempty" jsonschema:"title=Sub Title"`
	// Brief summary reflecting on the contents.
	Description string `json:"description,omitempty" jsonschema:"title=Description"`
	// Date on which this version was issued.
	IssueDate *cal.Date `json:"issue_date,omitempty" jsonschema:"title=Issue Date"`
	// Conditions and dates governing when this contract has effect.
	Effect *Effect `json:"effect,omitempty" jsonschema:"title=Effect"`
	// Parties and the roles in which they enter this contract.
	Parties []*Party `json:"parties,omitempty" jsonschema:"title=Parties"`
	// Background facts and purposes used to interpret the operative terms.
	Recitals []*Recital `json:"recitals,omitempty" jsonschema:"title=Recitals"`
	// Defined terms used by the document.
	Definitions []*Definition `json:"definitions,omitempty" jsonschema:"title=Definitions"`
	// Set of chapters that make up the contract.
	Chapters []*Chapter `json:"chapters,omitempty" jsonschema:"title=Chapters"`
	// Law that the parties intend to govern the agreement.
	Law *GoverningLaw `json:"law,omitempty" jsonschema:"title=Governing Law"`
	// Agreed process and forum for disputes.
	Disputes *DisputeResolution `json:"disputes,omitempty" jsonschema:"title=Dispute Resolution"`
	// Requirements for legally executing this version.
	Execution *Execution `json:"execution,omitempty" jsonschema:"title=Execution"`
	// Earlier versions, amendments, or other related legal documents.
	Preceding []*Reference `json:"preceding,omitempty" jsonschema:"title=Preceding Documents"`
	// Digest-bound material incorporated into or supporting the agreement.
	Resources []*Resource `json:"resources,omitempty" jsonschema:"title=Resources"`
	// Additional non-essential metadata.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Metadata"`
}

// Calculate prepares the contract for use as a GOBL document payload.
// Intrinsic normalization is registered by type; this method is only the
// standard document-level integration hook used by schema.Object.
func (c *Contract) Calculate() error {
	norm.Normalize(c)
	return nil
}
