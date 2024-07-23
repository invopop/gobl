package cli

import (
	"context"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/schema"
)

// BuildOptions are the options used for building and validating GOBL data.
type BuildOptions struct {
	*ParseOptions
}

// Build builds and validates GOBL data. Only structured errors are returned,
// which is a break from regular Go convention and replicated on all the main
// internal CLI functions. The object is to ensure that errors are always
// structured in a consistent manner.
func Build(ctx context.Context, opts *BuildOptions) (any, error) {
	obj, err := parseGOBLData(ctx, opts.ParseOptions)
	if err != nil {
		return nil, wrapError(StatusUnprocessableEntity, err)
	}

	if env, ok := obj.(*gobl.Envelope); ok {
		// 2024-04-05: Remove previous signatures. Assume the user knows what
		// they are doing and remove previous signatures as they're unlikely
		// to be useful.
		env.Signatures = nil

		if err := env.Calculate(); err != nil {
			return nil, wrapError(StatusUnprocessableEntity, err)
		}

		if err := env.Validate(); err != nil {
			return nil, wrapError(StatusUnprocessableEntity, err)
		}

		return env, nil
	}

	if doc, ok := obj.(*schema.Object); ok {
		if c, ok := doc.Instance().(schema.Calculable); ok {
			if err := c.Calculate(); err != nil {
				err = gobl.ErrCalculation.WithCause(err)
				return nil, wrapError(StatusUnprocessableEntity, err)
			}
		}

		if err := doc.Validate(); err != nil {
			err = gobl.ErrValidation.WithCause(err)
			return nil, wrapError(StatusUnprocessableEntity, err)
		}

		return doc, nil
	}

	panic("parsed data must be either an envelope or a document")
}
