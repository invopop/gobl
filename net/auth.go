package net

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/uuid"
)

// Request token limits: clients keep issuing windows short, and
// verifiers enforce a hard ceiling on token age regardless of the
// token's own exp claim.
const (
	minTokenTTL     = 30 * time.Second
	maxTokenTTL     = 5 * time.Minute
	defaultTokenTTL = time.Minute
	tokenClockSkew  = 30 * time.Second
)

// TokenClaims is the signed payload of a request token presented in
// the Authorization header of who and inbox requests. Iss names the
// party making the request — possibly a trusted intermediary distinct
// from any document signer inside the request body — and Aud names
// the destination address the token is bound to.
type TokenClaims struct {
	Iss cbc.URI `json:"iss"`
	Aud cbc.URI `json:"aud"`
	Iat int64   `json:"iat"`
	Exp int64   `json:"exp"`
	JTI string  `json:"jti,omitempty"`
}

// NewToken mints a request token: a compact ES256 JWS asserting that
// iss is making a request to aud. A ttl of zero uses the one minute
// default; other values are clamped to the 30 second – 5 minute
// window the protocol allows. A jti is stamped automatically for
// audit correlation.
func NewToken(key *dsig.PrivateKey, iss, aud Address, ttl time.Duration) (string, error) {
	iss, err := ParseAddress(string(iss))
	if err != nil {
		return "", err
	}
	aud, err = ParseAddress(string(aud))
	if err != nil {
		return "", err
	}
	switch {
	case ttl == 0:
		ttl = defaultTokenTTL
	case ttl < minTokenTTL:
		ttl = minTokenTTL
	case ttl > maxTokenTTL:
		ttl = maxTokenTTL
	}
	now := time.Now().UTC()
	claims := &TokenClaims{
		Iss: iss.URI(),
		Aud: aud.URI(),
		Iat: now.Unix(),
		Exp: now.Add(ttl).Unix(),
		JTI: uuid.V7().String(),
	}
	sig, err := dsig.NewSignature(key, claims)
	if err != nil {
		return "", fmt.Errorf("net: signing request token: %w", err)
	}
	return sig.String(), nil
}

// VerifyToken verifies a request token received on an inbound who or
// inbox request. The token's signature is checked against the
// issuer's published key, its aud claim must equal the given address,
// and its freshness window must include the current time. The
// verified requester address is returned.
func (c *Client) VerifyToken(ctx context.Context, token string, aud Address) (Address, error) {
	if token == "" {
		return "", fmt.Errorf("%w: token is empty", ErrTokenInvalid)
	}
	aud, err := ParseAddress(string(aud))
	if err != nil {
		return "", err
	}
	sig, err := dsig.ParseSignature(token)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrTokenInvalid, err)
	}
	claims := new(TokenClaims)
	if err := sig.UnsafePayload(claims); err != nil {
		return "", fmt.Errorf("%w: %v", ErrTokenInvalid, err)
	}
	if claims.Iss.Scheme() != Scheme {
		return "", fmt.Errorf("%w: iss %q is not a gobl address", ErrTokenInvalid, claims.Iss)
	}
	issuer, err := ParseAddress(claims.Iss.Opaque())
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrTokenInvalid, err)
	}
	kid := sig.KeyID()
	if kid == "" {
		return "", fmt.Errorf("%w: token has no key ID", ErrTokenInvalid)
	}
	key, err := c.FetchKey(ctx, issuer, kid)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrTokenInvalid, err)
	}
	verified := new(TokenClaims)
	if err := sig.VerifyPayload(key, verified); err != nil {
		return "", fmt.Errorf("%w: %v", ErrTokenInvalid, err)
	}
	if verified.Aud != aud.URI() {
		return "", fmt.Errorf("%w: audience mismatch (got %q, want %q)", ErrTokenInvalid, verified.Aud, aud.URI())
	}
	if verified.Iat == 0 || verified.Exp == 0 {
		return "", fmt.Errorf("%w: iat and exp are required", ErrTokenInvalid)
	}
	now := time.Now().UTC()
	iat := time.Unix(verified.Iat, 0)
	if iat.After(now.Add(tokenClockSkew)) {
		return "", fmt.Errorf("%w: issued in the future", ErrTokenExpired)
	}
	if now.After(time.Unix(verified.Exp, 0).Add(tokenClockSkew)) {
		return "", fmt.Errorf("%w: exp has passed", ErrTokenExpired)
	}
	if now.Sub(iat) > maxTokenTTL+tokenClockSkew {
		return "", fmt.Errorf("%w: iat is too old", ErrTokenExpired)
	}
	if err := key.Allows(iat); err != nil {
		return "", fmt.Errorf("%w: %v", ErrTokenInvalid, err)
	}
	return issuer, nil
}

// VerifyAuthorization verifies the value of an Authorization header,
// stripping the "Bearer" scheme before passing the token to
// VerifyToken. Servers call this with the raw header of an inbound
// who or inbox request and their own address.
func (c *Client) VerifyAuthorization(ctx context.Context, header string, aud Address) (Address, error) {
	scheme, token, found := strings.Cut(strings.TrimSpace(header), " ")
	if !found || !strings.EqualFold(scheme, "Bearer") {
		return "", fmt.Errorf("%w: authorization header is not a bearer token", ErrTokenInvalid)
	}
	return c.VerifyToken(ctx, strings.TrimSpace(token), aud)
}
