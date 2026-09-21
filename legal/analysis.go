package legal

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/uuid"
)

// Annotation kind keys classify optional machine interpretations of prose.
const (
	AnnotationDefinition     cbc.Key = "definition"
	AnnotationObligation     cbc.Key = "obligation"
	AnnotationPermission     cbc.Key = "permission"
	AnnotationProhibition    cbc.Key = "prohibition"
	AnnotationCondition      cbc.Key = "condition"
	AnnotationRepresentation cbc.Key = "representation"
	AnnotationWarranty       cbc.Key = "warranty"
	AnnotationRemedy         cbc.Key = "remedy"
	AnnotationTermination    cbc.Key = "termination"
)

// Analysis is a replaceable, non-authoritative interpretation of an exact
// contract version produced by an LLM, rules engine, or human reviewer.
type Analysis struct {
	uuid.Identify
	Contract    *Reference    `json:"contract" jsonschema:"title=Contract"`
	ProducedAt  cal.Timestamp `json:"produced_at" jsonschema:"title=Produced At"`
	Method      cbc.Code      `json:"method" jsonschema:"title=Method"`
	Annotations []*Annotation `json:"annotations" jsonschema:"title=Annotations"`
}

// Annotation links a semantic observation to an anchored fragment. It is an
// interpretation and never overrides the signed natural-language text.
type Annotation struct {
	Anchor  string   `json:"anchor" jsonschema:"title=Anchor"`
	Kind    cbc.Key  `json:"kind" jsonschema:"title=Kind"`
	Subject cbc.Key  `json:"subject,omitempty" jsonschema:"title=Subject Party"`
	Object  cbc.Key  `json:"object,omitempty" jsonschema:"title=Object Party"`
	Summary string   `json:"summary" jsonschema:"title=Summary"`
	Data    cbc.Meta `json:"data,omitempty" jsonschema:"title=Data"`
}

// Calculate prepares the analysis for use as a GOBL document payload.
func (a *Analysis) Calculate() error {
	norm.Normalize(a)
	return nil
}
