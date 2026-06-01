package ops

import (
	"context"
	"io"

	jsonyaml "github.com/invopop/yaml"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/net"
)

// Verify reads a GOBL document from in, and returns an error if there are any
// validation errors.
func Verify(ctx context.Context, in io.Reader, key *dsig.PublicKey) error {
	body, err := io.ReadAll(cancelableReader(ctx, in))
	if err != nil {
		return gobl.ErrInput.WithCause(err)
	}
	env := new(gobl.Envelope)
	if err := jsonyaml.Unmarshal(body, env); err != nil {
		return gobl.ErrInput.WithCause(err)
	}
	if err := env.Validate(); err != nil {
		return gobl.ErrValidation.WithCause(err)
	}
	if key == nil {
		return gobl.ErrInput.WithReason("public key required")
	}
	if !env.Signed() {
		return gobl.ErrSignature.WithReason("envelope is not signed")
	}
	if err := env.Signatures[0].VerifyPayload(key, env); err != nil {
		return gobl.ErrSignature.WithCause(err)
	}
	return nil
}

// VerifyRemote reads a GOBL envelope and verifies it using remote
// JWKS discovery via the GOBL Net client.
func VerifyRemote(ctx context.Context, in io.Reader, client *net.Client, addr net.Address) error {
	body, err := io.ReadAll(cancelableReader(ctx, in))
	if err != nil {
		return gobl.ErrInput.WithCause(err)
	}
	env := new(gobl.Envelope)
	if err := jsonyaml.Unmarshal(body, env); err != nil {
		return gobl.ErrInput.WithCause(err)
	}
	if err := env.Validate(); err != nil {
		return gobl.ErrValidation.WithCause(err)
	}
	issuer, err := client.VerifyEnvelope(ctx, env, "")
	if err != nil {
		return gobl.ErrValidation.WithCause(err)
	}
	if addr != "" && issuer != addr {
		return gobl.ErrValidation.WithReason("envelope signed by %s, expected %s", issuer, addr)
	}
	return nil
}
