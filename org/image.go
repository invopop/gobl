package org

import (
	"context"

	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

// Image describes a logo or photo that represents an entity. Most
// details except the URL are optional, but are potentially useful
// for validation if that's a requirement for the use case.
type Image struct {
	// Unique ID of the image
	UUID *uuid.UUID `json:"uuid,omitempty" jsonschema:"title=UUID"`
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
}

// Validate ensures the details on the image look okay.
func (i *Image) Validate() error {
	return i.ValidateWithContext(context.Background())
}

// ValidateWithContext ensures the details on the image look okay inside the provided context.
func (i *Image) ValidateWithContext(ctx context.Context) error {
	return tax.ValidateStructWithRegime(ctx, i,
		validation.Field(&i.UUID),
		validation.Field(&i.URL, is.URL),
		validation.Field(&i.Height, validation.Min(1), validation.Max(2048)),
		validation.Field(&i.Width, validation.Min(1), validation.Max(2048)),
		validation.Field(&i.Digest),
	)
}
