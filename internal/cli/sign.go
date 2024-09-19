package cli

import (
	"context"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/dsig"
)

// SignOptions are the options used for signing a GOBL document.
type SignOptions struct {
	*ParseOptions
	PrivateKey *dsig.PrivateKey
}

// Sign parses a GOBL document into an envelope, performs calculations,
// validates it, and finally signs its headers. The parsed envelope *must* be a
// draft, or else an error is returned.
func Sign(ctx context.Context, opts *SignOptions) (*gobl.Envelope, error) {
	// Always envelop incoming data.
	opts.Envelop = true

	obj, err := parseGOBLData(ctx, opts.ParseOptions)
	if err != nil {
		return nil, wrapError(StatusUnprocessableEntity, err)
	}

	env, ok := obj.(*gobl.Envelope)
	if !ok {
		panic("parsed sign data must be an envelope")
	}

	if err := env.Calculate(); err != nil {
		return nil, wrapError(StatusUnprocessableEntity, err)
	}

	// Sign envelope headers. Validation is done transparently in `Sign`.
	if err := env.Sign(opts.PrivateKey); err != nil {
		return nil, wrapError(StatusUnprocessableEntity, err)
	}

	return env, nil
}
