// Package internal contains internal objects that may be used for reference
// inside GOBL but are not intended for use outside of the library.
package internal

import (
	"context"
)

type contextKey string

const (
	// KeyDraft is used for extract the draft status from the context.
	KeyDraft contextKey = "draft"
)

// IsDraft returns true if the context indicates we're working with an
// envelope with the header draft status marked as true.
func IsDraft(ctx context.Context) bool {
	if ctx == nil {
		return false
	}
	return ctx.Value(KeyDraft) == true
}
