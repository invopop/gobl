package cli

import (
	"context"
	"io"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/schema"
)

// Validate asserts the contents of the envelope and document are correct.
func Validate(ctx context.Context, r io.Reader) error {
	opts := &ParseOptions{
		Input: r,
	}
	obj, err := parseGOBLData(ctx, opts)
	if err != nil {
		return gobl.ErrInput.WithCause(err)
	}

	if env, ok := obj.(*gobl.Envelope); ok {
		if err := env.Validate(); err != nil {
			return gobl.ErrValidation.WithCause(err)
		}
		return nil
	}

	if doc, ok := obj.(*schema.Object); ok {
		if err := doc.Validate(); err != nil {
			return gobl.ErrValidation.WithCause(err)
		}
		return nil
	}

	return gobl.ErrInput.WithReason("invalid document type")
}
