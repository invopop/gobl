package org

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/uuid"
)

// An Attachment provides a structure to be used to attach documents
// inside a GOBL document, either as a reference via a URL, or directly
// as a base64 encoded string.
//
// Attachments must not be used to include alternative versions of the
// same document, but rather include supporting documents or additional
// information that is not included in the main document.
//
// While it is possible to include the data directly in the JSON document,
// it is recommended to use the URL field to link to the document instead.
//
// Deprecation: attachments are being deprecated in favor of using Links
// in the `head` package defined inside the envelope. These allow for
// alternative versions or formats of the same source document to be provided
// without breaking signatures.
type Attachment struct {
	uuid.Identify

	// Key used to identify the attachment inside the document.
	Key cbc.Key `json:"key,omitempty" jsonschema:"title=Key"`

	// Code used to identify the payload of the attachment.
	Code cbc.Code `json:"code,omitempty" jsonschema:"title=Code"`

	// Filename of the attachment, will override name retrieved from URL.
	Name string `json:"name,omitempty" jsonschema:"title=Name"`

	// Details of why the attachment is being included and details on
	// what it contains.
	Description string `json:"description,omitempty" jsonschema:"title=Description"`

	// URL of where to find the attachment.
	URL string `json:"url" jsonschema:"title=URL,format=uri"`

	// Digest is used to verify the integrity of the attachment
	// when downloaded from the URL.
	Digest *dsig.Digest `json:"digest,omitempty" jsonschema:"title=Digest"`

	// MIME type of the attachment.
	MIME string `json:"mime,omitempty" jsonschema:"title=MIME Type"`
}

// Normalize will try to clean the attachment information.
func (a *Attachment) Normalize() {
	if a == nil {
		return
	}
	uuid.Normalize(&a.UUID)
	a.Code = cbc.NormalizeCode(a.Code)
	a.Name = cbc.NormalizeString(a.Name)
	a.Description = cbc.NormalizeString(a.Description)
	a.URL = cbc.NormalizeString(a.URL)
	a.MIME = cbc.NormalizeString(a.MIME)
}

func attachmentRules() *rules.Set {
	return rules.For(new(Attachment),
		rules.Field("url",
			rules.Assert("01", "attachment URL must be valid",
				is.Present,
				is.URL,
			),
		),
		rules.Field("mime",
			rules.AssertIfPresent("02", "attachment MIME type must be of an expected types defined in EN16931-1:2017",
				// MIME types as defined by EN 16931-1:2017 to be used as attachments.
				is.In(
					"application/pdf",
					"image/jpeg",
					"image/png",
					"text/csv",
					"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
					"application/vnd.oasis.opendocument.spreadsheet",
				),
			),
		),
	)
}
