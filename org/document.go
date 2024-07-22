package org

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

// Document is used to describe an existing document.
type Document struct {
	uuid.Identify
	// IssueDate reflects the date the document was issued.
	IssueDate *cal.Date `json:"issue_date,omitempty" jsonschema:"title=Issue Date"`
	// Series the referenced document belongs to.
	Series cbc.Code `json:"series,omitempty" jsonschema:"title=Series"`
	// Source document's code or other identifier.
	Code cbc.Code `json:"code,omitempty" jsonschema:"title=Code"`
	// Title or name of the document
	Title string `json:"title,omitempty" jsonschema:"title=Title"`
	// Additional details about the document.
	Description string `json:"description,omitempty" jsonschema:"title=Description"`
	// Link to the source document.
	URL string `json:"url,omitempty" jsonschema:"title=URL,format=uri"`
	// Meta contains additional information about the document.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Validate ensures the Document looks correct.
func (d *Document) Validate() error {
	return validation.ValidateStruct(d,
		validation.Field(&d.UUID),
		validation.Field(&d.Series),
		validation.Field(&d.Code),
		validation.Field(&d.URL, is.URL),
	)
}
