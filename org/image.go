package org

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/uuid"
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
	URL string `json:"url" jsonschema:"title=URL"`
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
	return validation.ValidateStruct(i,
		validation.Field(&i.UUID),
		validation.Field(&i.URL, validation.Required, is.URL),
		validation.Field(&i.Height, validation.Min(1), validation.Max(2048)),
		validation.Field(&i.Width, validation.Min(1), validation.Max(2048)),
		validation.Field(&i.Digest),
	)
}
