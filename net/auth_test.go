package net

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/invopop/gobl/dsig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// tokenPayload extracts the unverified claims from a compact token.
func tokenPayload(t *testing.T, token string) *TokenClaims {
	t.Helper()
	sig, err := dsig.ParseSignature(token)
	require.NoError(t, err)
	claims := new(TokenClaims)
	require.NoError(t, sig.UnsafePayload(claims))
	return claims
}

// signClaims mints a token with arbitrary claims, bypassing NewToken's
// clamping, for freshness edge cases.
func signClaims(t *testing.T, key *dsig.PrivateKey, claims *TokenClaims) string {
	t.Helper()
	sig, err := dsig.NewSignature(key, claims)
	require.NoError(t, err)
	return sig.String()
}

func TestNewToken(t *testing.T) {
	key := dsig.NewES256Key()
	requester := Address("sender.example.com")
	target := Address("receiver.example.com")

	t.Run("mints default window claims", func(t *testing.T) {
		token, err := NewToken(key, requester, target, 0)
		require.NoError(t, err)
		claims := tokenPayload(t, token)
		assert.Equal(t, requester.URI(), claims.Iss)
		assert.Equal(t, target.URI(), claims.Aud)
		assert.NotEmpty(t, claims.JTI)
		assert.Equal(t, int64(60), claims.Exp-claims.Iat)
	})

	t.Run("clamps the ttl to the allowed window", func(t *testing.T) {
		token, err := NewToken(key, requester, target, time.Second)
		require.NoError(t, err)
		claims := tokenPayload(t, token)
		assert.Equal(t, int64(30), claims.Exp-claims.Iat)

		token, err = NewToken(key, requester, target, time.Hour)
		require.NoError(t, err)
		claims = tokenPayload(t, token)
		assert.Equal(t, int64(300), claims.Exp-claims.Iat)
	})

	t.Run("canonicalizes addresses", func(t *testing.T) {
		token, err := NewToken(key, "Sender.Example.COM.", "RECEIVER.example.com", 0)
		require.NoError(t, err)
		claims := tokenPayload(t, token)
		assert.Equal(t, requester.URI(), claims.Iss)
		assert.Equal(t, target.URI(), claims.Aud)
	})

	t.Run("rejects an invalid address", func(t *testing.T) {
		_, err := NewToken(key, "not valid!", target, 0)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrAddressInvalid))
	})
}

func TestVerifyToken(t *testing.T) {
	ctx := context.Background()
	key := dsig.NewES256Key()
	requester := Address("sender.example.com")
	target := Address("receiver.example.com")

	keyData, err := json.Marshal(key.Public())
	require.NoError(t, err)
	client := NewClient(WithFetcher(&mapFetcher{data: map[string][]byte{
		requester.KeyURL(key.ID()): keyData,
	}}))

	t.Run("verifies a fresh token", func(t *testing.T) {
		token, err := NewToken(key, requester, target, 0)
		require.NoError(t, err)
		iss, err := client.VerifyToken(ctx, token, target)
		require.NoError(t, err)
		assert.Equal(t, requester, iss)
	})

	t.Run("rejects an empty token", func(t *testing.T) {
		_, err := client.VerifyToken(ctx, "", target)
		assert.True(t, errors.Is(err, ErrTokenInvalid))
	})

	t.Run("rejects garbage", func(t *testing.T) {
		_, err := client.VerifyToken(ctx, "not.a.token", target)
		assert.True(t, errors.Is(err, ErrTokenInvalid))
	})

	t.Run("rejects an audience mismatch", func(t *testing.T) {
		token, err := NewToken(key, requester, "other.example.com", 0)
		require.NoError(t, err)
		_, err = client.VerifyToken(ctx, token, target)
		assert.True(t, errors.Is(err, ErrTokenInvalid))
		assert.Contains(t, err.Error(), "audience mismatch")
	})

	t.Run("rejects a signer whose key is not published", func(t *testing.T) {
		otherKey := dsig.NewES256Key()
		token, err := NewToken(otherKey, requester, target, 0)
		require.NoError(t, err)
		_, err = client.VerifyToken(ctx, token, target)
		assert.True(t, errors.Is(err, ErrTokenInvalid))
	})

	t.Run("rejects a key substitution", func(t *testing.T) {
		// Signed with another key but claiming the published kid: the
		// signature must not verify against the published key.
		otherKey := dsig.NewES256Key()
		otherData, err := json.Marshal(otherKey.Public())
		require.NoError(t, err)
		c := NewClient(WithFetcher(&mapFetcher{data: map[string][]byte{
			requester.KeyURL(otherKey.ID()): otherData,
			requester.KeyURL(key.ID()):      otherData,
		}}))
		token, err := NewToken(key, requester, target, 0)
		require.NoError(t, err)
		_, err = c.VerifyToken(ctx, token, target)
		assert.True(t, errors.Is(err, ErrTokenInvalid))
	})

	t.Run("rejects a non-gobl issuer", func(t *testing.T) {
		now := time.Now().UTC()
		token := signClaims(t, key, &TokenClaims{
			Iss: "https://sender.example.com",
			Aud: target.URI(),
			Iat: now.Unix(),
			Exp: now.Add(time.Minute).Unix(),
		})
		_, err := client.VerifyToken(ctx, token, target)
		assert.True(t, errors.Is(err, ErrTokenInvalid))
	})

	t.Run("rejects missing iat or exp", func(t *testing.T) {
		token := signClaims(t, key, &TokenClaims{
			Iss: requester.URI(),
			Aud: target.URI(),
		})
		_, err := client.VerifyToken(ctx, token, target)
		assert.True(t, errors.Is(err, ErrTokenInvalid))
		assert.Contains(t, err.Error(), "iat and exp are required")
	})

	t.Run("rejects a passed exp", func(t *testing.T) {
		now := time.Now().UTC()
		token := signClaims(t, key, &TokenClaims{
			Iss: requester.URI(),
			Aud: target.URI(),
			Iat: now.Add(-2 * time.Minute).Unix(),
			Exp: now.Add(-time.Minute).Unix(),
		})
		_, err := client.VerifyToken(ctx, token, target)
		assert.True(t, errors.Is(err, ErrTokenExpired))
	})

	t.Run("rejects an iat older than the ceiling even with a live exp", func(t *testing.T) {
		now := time.Now().UTC()
		token := signClaims(t, key, &TokenClaims{
			Iss: requester.URI(),
			Aud: target.URI(),
			Iat: now.Add(-10 * time.Minute).Unix(),
			Exp: now.Add(10 * time.Minute).Unix(),
		})
		_, err := client.VerifyToken(ctx, token, target)
		assert.True(t, errors.Is(err, ErrTokenExpired))
		assert.Contains(t, err.Error(), "too old")
	})

	t.Run("rejects an iat in the future", func(t *testing.T) {
		now := time.Now().UTC()
		token := signClaims(t, key, &TokenClaims{
			Iss: requester.URI(),
			Aud: target.URI(),
			Iat: now.Add(5 * time.Minute).Unix(),
			Exp: now.Add(6 * time.Minute).Unix(),
		})
		_, err := client.VerifyToken(ctx, token, target)
		assert.True(t, errors.Is(err, ErrTokenExpired))
		assert.Contains(t, err.Error(), "future")
	})
}

func TestVerifyAuthorization(t *testing.T) {
	ctx := context.Background()
	key := dsig.NewES256Key()
	requester := Address("sender.example.com")
	target := Address("receiver.example.com")

	keyData, err := json.Marshal(key.Public())
	require.NoError(t, err)
	client := NewClient(WithFetcher(&mapFetcher{data: map[string][]byte{
		requester.KeyURL(key.ID()): keyData,
	}}))

	token, err := NewToken(key, requester, target, 0)
	require.NoError(t, err)

	t.Run("accepts a bearer header", func(t *testing.T) {
		iss, err := client.VerifyAuthorization(ctx, "Bearer "+token, target)
		require.NoError(t, err)
		assert.Equal(t, requester, iss)
	})

	t.Run("scheme is case-insensitive", func(t *testing.T) {
		iss, err := client.VerifyAuthorization(ctx, "bearer "+token, target)
		require.NoError(t, err)
		assert.Equal(t, requester, iss)
	})

	t.Run("rejects a missing scheme", func(t *testing.T) {
		_, err := client.VerifyAuthorization(ctx, token, target)
		assert.True(t, errors.Is(err, ErrTokenInvalid))
	})

	t.Run("rejects another scheme", func(t *testing.T) {
		_, err := client.VerifyAuthorization(ctx, "Basic dXNlcjpwdw==", target)
		assert.True(t, errors.Is(err, ErrTokenInvalid))
	})

	t.Run("rejects an empty header", func(t *testing.T) {
		_, err := client.VerifyAuthorization(ctx, "", target)
		assert.True(t, errors.Is(err, ErrTokenInvalid))
	})
}

// headerRecorder wraps a Fetcher and records the headers sent per URL.
type headerRecorder struct {
	inner Fetcher
	byURL map[string]http.Header
}

func (h *headerRecorder) Fetch(ctx context.Context, url string, header http.Header) ([]byte, error) {
	h.byURL[url] = header
	return h.inner.Fetch(ctx, url, header)
}

func TestWhoAuthorization(t *testing.T) {
	ctx := context.Background()
	subject := Address("supplier.example.com")
	subjKey := dsig.NewES256Key()
	self := Address("customer.example.com")
	selfKey := dsig.NewES256Key()

	subjKeyData, err := json.Marshal(subjKey.Public())
	require.NoError(t, err)

	t.Run("attaches a bearer token when the client has an identity", func(t *testing.T) {
		rec := &headerRecorder{
			inner: &mapFetcher{data: map[string][]byte{
				subject.WhoURL():             buildPartyEnvelope(t, subjKey, subject, nil, "", ""),
				subject.KeyURL(subjKey.ID()): subjKeyData,
			}},
			byURL: map[string]http.Header{},
		}
		c := NewClient(WithFetcher(rec), WithIdentity(self, selfKey))
		_, err := c.Who(ctx, subject)
		require.NoError(t, err)

		auth := rec.byURL[subject.WhoURL()].Get("Authorization")
		require.NotEmpty(t, auth)
		require.True(t, len(auth) > 7 && auth[:7] == "Bearer ")
		claims := tokenPayload(t, auth[7:])
		assert.Equal(t, self.URI(), claims.Iss)
		assert.Equal(t, subject.URI(), claims.Aud)

		// Key fetches go out bare.
		assert.Empty(t, rec.byURL[subject.KeyURL(subjKey.ID())].Get("Authorization"))
	})

	t.Run("sends no token without an identity", func(t *testing.T) {
		rec := &headerRecorder{
			inner: &mapFetcher{data: map[string][]byte{
				subject.WhoURL():             buildPartyEnvelope(t, subjKey, subject, nil, "", ""),
				subject.KeyURL(subjKey.ID()): subjKeyData,
			}},
			byURL: map[string]http.Header{},
		}
		c := NewClient(WithFetcher(rec))
		_, err := c.Who(ctx, subject)
		require.NoError(t, err)
		assert.Empty(t, rec.byURL[subject.WhoURL()].Get("Authorization"))
	})

	t.Run("202 reports ErrPending", func(t *testing.T) {
		c := NewClient(WithFetcher(&mockFetcher{err: ErrPending}))
		_, err := c.Who(ctx, subject)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrPending))
	})
}
