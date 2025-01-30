package cbc

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

// Source is used to identify a specific source of data. Typically this is used
// as part of other structures to identify where the data came from.
type Source struct {
	// Title stores the name of the source of information.
	Title i18n.String `json:"title,omitempty" jsonschema:"title=Title"`

	// URL is the location of the source of information.
	URL string `json:"url" jsonschema:"title=URL,format=uri"`

	// ContentType of the information expected at the URL.
	ContentType string `json:"content_type,omitempty" jsonschema:"title=Content Type"`

	// At is the date and time the information was retrieved.
	At *cal.DateTime `json:"at,omitempty" jsonschema:"title=At"`
}

// Validate ensures that the source object looks valid.
func (src *Source) Validate() error {
	return validation.ValidateStruct(src,
		validation.Field(&src.Title),
		validation.Field(&src.URL, validation.Required, is.URL),
		validation.Field(&src.ContentType),
	)
}
