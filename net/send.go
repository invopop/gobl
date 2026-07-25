package net

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/invopop/gobl"
)

// Send delivers a signed envelope to the address's inbox with a POST
// on its well-known inbox endpoint. The request carries a bearer
// request token minted from the client's identity (WithIdentity),
// which may name a different party than the envelope's signer — a
// trusted intermediary transmitting on the signer's behalf.
//
// A 202 response means the inbox has persisted the envelope. Any
// other 4xx — except 429 — returns ErrInboxRejected: the inbox has
// decided, do not retry. Transient conditions (429, 5xx, transport
// failures) return ErrUnavailable, and delivery is idempotent on the
// envelope's uuid and digest, so "retry on ErrUnavailable until 202"
// is the correct recovery strategy. Other errors are permanent input
// or configuration failures.
func (c *Client) Send(ctx context.Context, addr Address, env *gobl.Envelope) error {
	if env == nil {
		return fmt.Errorf("%w: envelope is nil", ErrFetchFailed)
	}
	addr, err := ParseAddress(string(addr))
	if err != nil {
		return err
	}
	body, err := json.Marshal(env)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFetchFailed, err)
	}
	header, err := c.authHeader(addr)
	if err != nil {
		return err
	}
	return c.fetcher.Post(ctx, addr.InboxURL(), body, header)
}
