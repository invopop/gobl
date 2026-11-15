package legal

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/uuid"
)

// Reference relation keys describe why another legal artifact is relevant.
const (
	ReferenceRelationAgreement  cbc.Key = "agreement"
	ReferenceRelationAmends     cbc.Key = "amends"
	ReferenceRelationRestates   cbc.Key = "restates"
	ReferenceRelationSupersedes cbc.Key = "supersedes"
	ReferenceRelationEvidence   cbc.Key = "evidence"
)

// Reference identifies another immutable legal artifact. Digest should be
// supplied whenever the exact referenced bytes matter.
type Reference struct {
	Relation cbc.Key      `json:"relation" jsonschema:"title=Relation"`
	UUID     uuid.UUID    `json:"uuid,omitempty" jsonschema:"title=UUID"`
	Code     cbc.Code     `json:"code,omitempty" jsonschema:"title=Code"`
	Title    string       `json:"title,omitempty" jsonschema:"title=Title"`
	URL      string       `json:"url,omitempty" jsonschema:"title=URL,format=uri"`
	Digest   *dsig.Digest `json:"digest,omitempty" jsonschema:"title=Digest"`
}

// Resource is supporting or incorporated material. Incorporated resources
// require a digest so a signature commits to the referenced content, not just
// to a mutable URL.
type Resource struct {
	Key          cbc.Key      `json:"key" jsonschema:"title=Key"`
	Title        string       `json:"title" jsonschema:"title=Title"`
	URL          string       `json:"url" jsonschema:"title=URL,format=uri"`
	MIME         string       `json:"mime,omitempty" jsonschema:"title=MIME Type,format=mime"`
	Digest       *dsig.Digest `json:"digest,omitempty" jsonschema:"title=Digest"`
	Incorporated bool         `json:"incorporated,omitempty" jsonschema:"title=Incorporated"`
}
