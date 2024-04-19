package org

import (
	"context"

	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

// Website describes what is expected for a web address.
type Website struct {
	// Unique identity code
	UUID uuid.UUID `json:"uuid,omitempty" jsonschema:"title=UUID"`
	// Identifier for this number.
	Label string `json:"label,omitempty" jsonschema:"title=Label"`
	// Title of the website to help distinguish between this and other links.
	Title string `json:"title,omitempty" jsonschema:"title=Title"`
	// URL for the website.
	URL string `json:"url" jsonschema:"title=URL,format=uri"`
}

// Validate checks the website objects URL to ensure it looks correct.
func (w *Website) Validate() error {
	return w.ValidateWithContext(context.Background())
}

// ValidateWithContext checks the website objects URL to ensure it looks correct inside the provided context.
func (w *Website) ValidateWithContext(ctx context.Context) error {
	return tax.ValidateStructWithRegime(ctx, w,
		validation.Field(&w.URL, validation.Required, is.URL),
	)
}
