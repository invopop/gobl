package net

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/invopop/gobl"
)

// Send posts an envelope to addr's inbox. WithIdentity adds a request token;
// its issuer may differ from the envelope signer. A 202 response succeeds,
// definitive 4xx responses return ErrInboxRejected, and transient failures
// return ErrUnavailable.
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
