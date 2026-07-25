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
// other 4xx — except 429, which is retryable — returns
// ErrInboxRejected. Delivery is idempotent on the envelope's uuid and
// digest, so "retry until 202" is the correct recovery strategy.
func (c *Client) Send(ctx context.Context, addr Address, env *gobl.Envelope) error {
	addr, err := ParseAddress(string(addr))
	if err != nil {
		return err
	}
	poster, ok := c.fetcher.(Poster)
	if !ok {
		return fmt.Errorf("%w: fetcher does not support POST", ErrFetchFailed)
	}
	body, err := json.Marshal(env)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFetchFailed, err)
	}
	header, err := c.authHeader(addr)
	if err != nil {
		return err
	}
	return poster.Post(ctx, addr.InboxURL(), body, header)
}
