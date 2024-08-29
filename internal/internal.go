// Package internal contains internal objects that may be used for reference
// inside GOBL but are not intended for use outside of the library.
package internal

import (
	"context"
)

type contextKey string

const (
	// KeySigned is the context key used to indicate that the envelope is signed.
	KeySigned contextKey = "signed"
)

// IsSigned returns true if the context indicates we're working with an
// envelope with signatures.
func IsSigned(ctx context.Context) bool {
	if ctx == nil {
		return false
	}
	return ctx.Value(KeySigned) == true
}

// SignedContext returns a new context with the signed flag set.
func SignedContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, KeySigned, true)
}
