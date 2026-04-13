package org

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/uuid"
)

// Image describes a logo or photo that represents an entity. Most
// details except the URL are optional, but are potentially useful
// for validation if that's a requirement for the use case.
type Image struct {
	uuid.Identify
	// Label to help identify the image.
	Label string `json:"label,omitempty" jsonschema:"title=Label"`
	// URL of the image
	URL string `json:"url,omitempty" jsonschema:"title=URL"`
	// As an alternative to the URL and only when the source data is small,
	// like an SVG, the raw data may be provided using Base64 encoding.
	Data []byte `json:"data,omitempty" jsonschema:"title=Data"`
	// Format of the image.
	MIME string `json:"mime,omitempty" jsonschema:"title=MIME"`
	// Details of what the image represents.
	Description string `json:"description,omitempty" jsonschema:"title=Description"`
	// Alternative text if the image cannot be shown.
	Alt string `json:"alt,omitempty" jsonschema:"title=Alt"`
	// Height of the image in pixels.
	Height int32 `json:"height,omitempty" jsonschema:"title=Height"`
	// Width of the image in pixels.
	Width int32 `json:"width,omitempty" jsonschema:"title=Width"`
	// Digest can be used to ensure the image contained at the URL
	// is the same one as originally intended.
	Digest *dsig.Digest `json:"digest,omitempty" jsonschema:"title=Digest"`
	// Meta contains additional information about the image.
	Meta *cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

func imageRules() *rules.Set {
	return rules.For(new(Image),
		rules.Field("url",
			rules.AssertIfPresent("01", "image URL must be valid", is.URL),
		),
		rules.Field("height",
			rules.Assert("02", "image height must be between 64 and 2048 pixels",
				is.Min(int32(64)),
				is.Max(int32(2048)),
			),
		),
		rules.Field("width",
			rules.Assert("03", "image width must be between 64 and 2048 pixels",
				is.Min(int32(64)),
				is.Max(int32(2048)),
			),
		),
	)
}
