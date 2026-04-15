package org

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/uuid"
)

// Website describes what is expected for a web address.
type Website struct {
	uuid.Identify
	// Label for the website to show alongside the URL.
	Label string `json:"label,omitempty" jsonschema:"title=Label"`
	// Title of the website to help distinguish between this and other links.
	Title string `json:"title,omitempty" jsonschema:"title=Title"`
	// URL for the website.
	URL string `json:"url" jsonschema:"title=URL,format=uri"`
}

// Normalize will try to remove any unnecessary whitespace from the website fields.
func (w *Website) Normalize() {
	if w == nil {
		return
	}
	uuid.Normalize(&w.UUID)
	w.Label = cbc.NormalizeString(w.Label)
	w.Title = cbc.NormalizeString(w.Title)
	w.URL = cbc.NormalizeString(w.URL)
}

func websiteRules() *rules.Set {
	return rules.For(new(Website),
		rules.Field("url",
			rules.Assert("01", "website URL is required and must be valid", is.Present, is.URL),
		),
	)
}
