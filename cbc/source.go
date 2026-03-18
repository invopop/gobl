package cbc

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
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

func sourceRules() *rules.Set {
	return rules.For(new(Source),
		rules.Field("url",
			rules.Assert("01", "url is required and must be a URL",
				rules.Present,
				is.URL,
			),
		),
	)
}
