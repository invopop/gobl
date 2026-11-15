package net

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/head"
)

// Sign sets env's header from and to fields and signs it with key as from
// for to, the pair a GOBL Net envelope must carry a signature for. Options
// such as head.WithExpiration or head.WithVerifier are added to the
// signature.
func Sign(env *gobl.Envelope, key *dsig.PrivateKey, from, to Address, opts ...head.SignOption) error {
	if env == nil || env.Head == nil {
		return fmt.Errorf("%w: envelope has no header", ErrVerifyFailed)
	}
	from, err := ParseAddress(string(from))
	if err != nil {
		return err
	}
	to, err = ParseAddress(string(to))
	if err != nil {
		return err
	}
	env.Head.From, env.Head.To = from.URI(), to.URI()
	opts = append(opts, head.WithIssuer(from.String()), head.WithAudience(to.String()))
	return env.Sign(key, opts...)
}

// SignKeyChange signs env as from for to with current and next, the
// current key's signature naming next as its bound key: a registration
// request changing the subject's key from current to next (README.md
// §7.3). Options are added to both signatures.
func SignKeyChange(env *gobl.Envelope, current, next *dsig.PrivateKey, from, to Address, opts ...head.SignOption) error {
	if current == nil || next == nil {
		return fmt.Errorf("%w: a key change needs the current and the next key", ErrVerifyFailed)
	}
	named := append(opts[:len(opts):len(opts)], head.WithConfirmation(next.JKT()))
	if err := Sign(env, current, from, to, named...); err != nil {
		return err
	}
	return Sign(env, next, from, to, opts...)
}

// Send posts env to the inbox of its to address. A 202 response succeeds,
// definitive 4xx responses return ErrInboxRejected, and transient failures
// return ErrUnavailable.
func (c *Client) Send(ctx context.Context, env *gobl.Envelope) error {
	if env == nil || env.Head == nil {
		return fmt.Errorf("%w: envelope is nil or has no header", ErrFetchFailed)
	}
	_, to, err := routing(env)
	if err != nil {
		return err
	}
	body, err := json.Marshal(env)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFetchFailed, err)
	}
	return c.fetcher.Post(ctx, to.InboxURL(), body)
}
