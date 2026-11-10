package net

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

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

// TokenClaims is the signed payload of a Who or Send request token. Iss is the
// requester, which may differ from an envelope signer, and Aud is the target.
type TokenClaims struct {
	Iss Address `json:"iss"`
	Aud Address `json:"aud"`
	Iat int64   `json:"iat"`
	Exp int64   `json:"exp"`
	JTI string  `json:"jti,omitempty"`
}

// NewToken creates a compact ES256 request token from iss to aud. A zero TTL
// uses one minute; other values are clamped to 30 seconds through five minutes.
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
		Iss: iss,
		Aud: aud,
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

// VerifyToken verifies a request token's signature, audience, key validity, and
// freshness, then returns the canonical requester address.
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
	issuer, err := ParseAddress(string(claims.Iss))
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrTokenInvalid, err)
	}
	kid := sig.KeyID()
	if kid == "" {
		return "", fmt.Errorf("%w: token has no key ID", ErrTokenInvalid)
	}
	key, err := c.FetchKey(ctx, issuer, kid)
	if err != nil {
		if errors.Is(err, ErrUnavailable) {
			// The issuer's key endpoint could not be reached: a
			// transient condition, not an invalid token. Servers
			// respond 503, not 401.
			return "", err
		}
		return "", fmt.Errorf("%w: %v", ErrTokenInvalid, err)
	}
	verified := new(TokenClaims)
	if err := sig.VerifyPayload(key, verified); err != nil {
		return "", fmt.Errorf("%w: %v", ErrTokenInvalid, err)
	}
	if got, err := ParseAddress(string(verified.Aud)); err != nil || got != aud {
		return "", fmt.Errorf("%w: audience mismatch (got %q, want %q)", ErrTokenInvalid, verified.Aud, aud)
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

// VerifyAuthorization verifies a Bearer Authorization header for aud and
// returns the canonical requester address.
func (c *Client) VerifyAuthorization(ctx context.Context, header string, aud Address) (Address, error) {
	scheme, token, found := strings.Cut(strings.TrimSpace(header), " ")
	if !found || !strings.EqualFold(scheme, "Bearer") {
		return "", fmt.Errorf("%w: authorization header is not a bearer token", ErrTokenInvalid)
	}
	return c.VerifyToken(ctx, strings.TrimSpace(token), aud)
}
