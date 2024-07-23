package cli

import (
	"context"
	"io"
	"net/http"

	jsonyaml "github.com/invopop/yaml"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/internal/iotools"
)

// Verify reads a GOBL document from in, and returns an error if there are any
// validation errors.
func Verify(ctx context.Context, in io.Reader, key *dsig.PublicKey) error {
	body, err := io.ReadAll(iotools.CancelableReader(ctx, in))
	if err != nil {
		return wrapError(StatusBadRequest, err)
	}
	env := new(gobl.Envelope)
	if err := jsonyaml.Unmarshal(body, env); err != nil {
		return wrapError(StatusBadRequest, err)
	}
	if err := env.Validate(); err != nil {
		return wrapError(StatusUnprocessableEntity, err)
	}
	if key == nil {
		return wrapErrorf(StatusBadRequest, "public key required")
	}
	if env.Head.Draft {
		return wrapErrorf(http.StatusUnprocessableEntity, "document is a draft")
	}
	if err := env.Signatures[0].VerifyPayload(key, env); err != nil {
		return wrapError(http.StatusUnprocessableEntity, err)
	}
	return nil
}
