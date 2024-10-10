package org

import (
	"context"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

// DocumentRef is used to describe an existing document or a specific part of it's contents.
type DocumentRef struct {
	uuid.Identify
	// Type of the document referenced.
	Type cbc.Key `json:"type,omitempty" jsonschema:"title=Type"`
	// IssueDate reflects the date the document was issued.
	IssueDate *cal.Date `json:"issue_date,omitempty" jsonschema:"title=Issue Date"`
	// Series the referenced document belongs to.
	Series cbc.Code `json:"series,omitempty" jsonschema:"title=Series"`
	// Source document's code or other identifier.
	Code cbc.Code `json:"code" jsonschema:"title=Code" en16931:"BT-122"`
	// Line index number inside the document, if relevant.
	Line int `json:"line,omitempty" jsonschema:"title=Line"`
	// List of additional codes, IDs, or SKUs which can be used to identify the document or its contents, agreed upon by the supplier and customer.
	Identities []*Identity `json:"identities,omitempty" jsonschema:"title=Identities"`
	// Tax period in which the referred document had an effect required by some tax regimes and formats.
	Period *cal.Period `json:"period,omitempty" jsonschema:"title=Period"`
	// Human readable description on why this reference is here or needs to be used.
	Reason string `json:"reason,omitempty" jsonschema:"title=Reason"`
	// Additional details about the document.
	Description string `json:"description,omitempty" jsonschema:"title=Description" en16931:"BT-123"`
	// Seals of approval from other organisations that may need to be listed.
	Stamps []*head.Stamp `json:"stamps,omitempty" jsonschema:"title=Stamps"`
	// Link to the source document.
	URL string `json:"url,omitempty" jsonschema:"title=URL,format=uri" en16931:"BT-124"`
	// Extensions for additional codes that may be required.
	Ext tax.Extensions `json:"ext,omitempty" jsonschemaL:"title=Extensions"`
	// Meta contains additional information about the document.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Normalize attempts to clean and normalize the DocumentRef.
func (dr *DocumentRef) Normalize(normalizers tax.Normalizers) {
	if dr == nil {
		return
	}
	dr.Ext = tax.CleanExtensions(dr.Ext)
	dr.Series = cbc.NormalizeCode(dr.Series)
	dr.Code = cbc.NormalizeCode(dr.Code)
	normalizers.Each(dr)
	tax.Normalize(normalizers, dr.Identities)
}

// Validate ensures the Document looks correct.
func (dr *DocumentRef) Validate() error {
	return dr.ValidateWithContext(context.Background())
}

// ValidateWithContext ensures the Document looks correct within the provided context.
func (dr *DocumentRef) ValidateWithContext(ctx context.Context) error {
	return tax.ValidateStructWithContext(ctx, dr,
		validation.Field(&dr.UUID),
		validation.Field(&dr.Type),
		validation.Field(&dr.IssueDate, cal.DateNotZero()),
		validation.Field(&dr.Series),
		validation.Field(&dr.Code,
			validation.Match(cbc.CodePatternRegexp),
			validation.Required,
		),
		validation.Field(&dr.URL, is.URL),
		validation.Field(&dr.Stamps),
		validation.Field(&dr.Period),
		validation.Field(&dr.Ext),
		validation.Field(&dr.Meta),
	)
}
