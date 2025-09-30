package org

import (
	"context"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/num"
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
	Code cbc.Code `json:"code" jsonschema:"title=Code"`
	// Currency used in the document, if different from the parent's currency.
	Currency currency.Code `json:"currency,omitempty" jsonschema:"title=Currency"`
	// Line index numbers inside the document, if relevant.
	Lines []int `json:"lines,omitempty" jsonschema:"title=Lines"`
	// List of additional codes, IDs, or SKUs which can be used to identify the document or its contents, agreed upon by the supplier and customer.
	Identities []*Identity `json:"identities,omitempty" jsonschema:"title=Identities"`
	// Tax period in which the referred document had an effect required by some tax regimes and formats.
	Period *cal.Period `json:"period,omitempty" jsonschema:"title=Period"`
	// Human readable description on why this reference is here or needs to be used.
	Reason string `json:"reason,omitempty" jsonschema:"title=Reason"`
	// Additional details about the document.
	Description string `json:"description,omitempty" jsonschema:"title=Description"`
	// Seals of approval from other organizations that may need to be listed.
	Stamps []*head.Stamp `json:"stamps,omitempty" jsonschema:"title=Stamps"`
	// Link to the source document.
	URL string `json:"url,omitempty" jsonschema:"title=URL,format=uri"`
	// Tax total breakdown from the original document in the provided currency. Should
	// only be included if required by a specific tax regime or addon.
	Tax *tax.Total `json:"tax,omitempty" jsonschema:"title=Tax"`
	// Payable is the total amount that is payable in the referenced document. Only needed
	// for specific tax regimes or addons. This may also be used in some scenarios
	// to determine the proportion of the referenced document that has been paid, and
	// calculate the remaining amount due and taxes.
	Payable *num.Amount `json:"payable,omitempty" jsonschema:"title=Payable"`
	// Extensions for additional codes that may be required.
	Ext tax.Extensions `json:"ext,omitempty" jsonschema:"title=Extensions"`
	// Meta contains additional information about the document.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Normalize attempts to clean and normalize the DocumentRef.
func (dr *DocumentRef) Normalize(normalizers tax.Normalizers) {
	if dr == nil {
		return
	}
	uuid.Normalize(&dr.UUID)
	dr.Series = cbc.NormalizeCode(dr.Series)
	dr.Code = cbc.NormalizeCode(dr.Code)
	dr.Reason = cbc.NormalizeString(dr.Reason)
	dr.URL = cbc.NormalizeString(dr.URL)
	dr.Ext = tax.CleanExtensions(dr.Ext)

	normalizers.Each(dr)
	tax.Normalize(normalizers, dr.Identities)
	tax.Normalize(normalizers, dr.Tax)
}

// Calculate will ensure the tax total is recalculated according to the
// rounding rule and currency precision provided. Users of this should first
// check the optional currency property of the document ref to see if that
// should be used instead.
func (dr *DocumentRef) Calculate(cur currency.Code, rr cbc.Key) {
	if dr == nil || dr.Tax == nil {
		return
	}
	dr.Tax.Calculate(cur, rr)
	dr.Tax.Round(cur.Def().Zero())
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
		validation.Field(&dr.Currency),
		validation.Field(&dr.URL, is.URL),
		validation.Field(&dr.Stamps),
		validation.Field(&dr.Period),
		validation.Field(&dr.Tax),
		validation.Field(&dr.Payable),
		validation.Field(&dr.Ext),
		validation.Field(&dr.Meta),
	)
}
