package cli

import (
	"context"
	"io"
	"net/http"

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
		return wrapError(StatusUnprocessableEntity, err)
	}

	if env, ok := obj.(*gobl.Envelope); ok {
		if err := env.Validate(); err != nil {
			return wrapError(http.StatusUnprocessableEntity, err)
		}
		return nil
	}

	if doc, ok := obj.(*schema.Object); ok {
		if err := doc.Validate(); err != nil {
			return wrapError(http.StatusUnprocessableEntity, err)
		}
		return nil
	}

	return wrapErrorf(http.StatusUnprocessableEntity, "invalid document type")
}
