package cli

import (
	"context"
	"net/http"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/schema"
)

// ReplicateOptions define all the basic options required to build a replicated
// document from the input.
type ReplicateOptions struct {
	*ParseOptions
}

// Replicate takes a base document as input and builds a replicated document
// for the output.
func Replicate(ctx context.Context, opts *ReplicateOptions) (interface{}, error) {
	res, err := replicate(ctx, opts)
	if err != nil {
		return nil, wrapError(http.StatusUnprocessableEntity, err)
	}
	return res, nil
}

func replicate(ctx context.Context, opts *ReplicateOptions) (interface{}, error) {
	obj, err := parseGOBLData(ctx, opts.ParseOptions)
	if err != nil {
		return nil, err
	}

	if env, ok := obj.(*gobl.Envelope); ok {
		e2, err := env.Replicate()
		if err != nil {
			return nil, err
		}
		if err = e2.Validate(); err != nil {
			return nil, err
		}
		return e2, nil
	}

	if doc, ok := obj.(*schema.Object); ok {
		// Documents are updated in place
		if err := doc.Replicate(); err != nil {
			return nil, err
		}
		if err = doc.Calculate(); err != nil {
			return nil, err
		}
		if err = doc.Validate(); err != nil {
			return nil, err
		}
		return doc, nil
	}

	panic("input must be either an envelope or a document")
}
