package net

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/invopop/gobl/cal"

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
		assert.Equal(t, requester, claims.Iss)
		assert.Equal(t, target, claims.Aud)
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
		assert.Equal(t, requester, claims.Iss)
		assert.Equal(t, target, claims.Aud)
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

	t.Run("rejects a non-address issuer", func(t *testing.T) {
		now := time.Now().UTC()
		token := signClaims(t, key, &TokenClaims{
			Iss: "https://sender.example.com",
			Aud: target,
			Iat: now.Unix(),
			Exp: now.Add(time.Minute).Unix(),
		})
		_, err := client.VerifyToken(ctx, token, target)
		assert.True(t, errors.Is(err, ErrTokenInvalid))
	})

	t.Run("rejects missing iat or exp", func(t *testing.T) {
		token := signClaims(t, key, &TokenClaims{
			Iss: requester,
			Aud: target,
		})
		_, err := client.VerifyToken(ctx, token, target)
		assert.True(t, errors.Is(err, ErrTokenInvalid))
		assert.Contains(t, err.Error(), "iat and exp are required")
	})

	t.Run("rejects a passed exp", func(t *testing.T) {
		now := time.Now().UTC()
		token := signClaims(t, key, &TokenClaims{
			Iss: requester,
			Aud: target,
			Iat: now.Add(-2 * time.Minute).Unix(),
			Exp: now.Add(-time.Minute).Unix(),
		})
		_, err := client.VerifyToken(ctx, token, target)
		assert.True(t, errors.Is(err, ErrTokenExpired))
	})

	t.Run("rejects an iat older than the ceiling even with a live exp", func(t *testing.T) {
		now := time.Now().UTC()
		token := signClaims(t, key, &TokenClaims{
			Iss: requester,
			Aud: target,
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
			Iss: requester,
			Aud: target,
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
		assert.Equal(t, self, claims.Iss)
		assert.Equal(t, subject, claims.Aud)

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

// signWithoutKID builds a compact ES256 JWS with no kid header, which
// dsig cannot produce but a hostile client could.
func signWithoutKID(t *testing.T, payload any) string {
	t.Helper()
	raw, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.ES256, Key: raw}, nil)
	require.NoError(t, err)
	data, err := json.Marshal(payload)
	require.NoError(t, err)
	jws, err := signer.Sign(data)
	require.NoError(t, err)
	token, err := jws.CompactSerialize()
	require.NoError(t, err)
	return token
}

// spoofKID re-labels a published JWK with another signature's kid so
// the per-key fetch succeeds but the crypto check must fail.
func spoofKID(t *testing.T, key *dsig.PrivateKey, kid string) []byte {
	t.Helper()
	data, err := json.Marshal(key.Public())
	require.NoError(t, err)
	m := map[string]any{}
	require.NoError(t, json.Unmarshal(data, &m))
	m["kid"] = kid
	out, err := json.Marshal(m)
	require.NoError(t, err)
	return out
}

func TestNewTokenEdgeCases(t *testing.T) {
	key := dsig.NewES256Key()

	t.Run("rejects an invalid audience", func(t *testing.T) {
		_, err := NewToken(key, "sender.example.com", "not valid!", 0)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrAddressInvalid))
	})

	t.Run("rejects an unusable key", func(t *testing.T) {
		_, err := NewToken(new(dsig.PrivateKey), "sender.example.com", "receiver.example.com", 0)
		require.Error(t, err)
	})
}

func TestVerifyTokenEdgeCases(t *testing.T) {
	ctx := context.Background()
	key := dsig.NewES256Key()
	requester := Address("sender.example.com")
	target := Address("receiver.example.com")

	keyData, err := json.Marshal(key.Public())
	require.NoError(t, err)
	client := NewClient(WithFetcher(&mapFetcher{data: map[string][]byte{
		requester.KeyURL(key.ID()): keyData,
	}}))

	t.Run("rejects an invalid audience argument", func(t *testing.T) {
		token, err := NewToken(key, requester, target, 0)
		require.NoError(t, err)
		_, err = client.VerifyToken(ctx, token, "not valid!")
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrAddressInvalid))
	})

	t.Run("rejects a non-object payload", func(t *testing.T) {
		sig, err := dsig.NewSignature(key, "not-an-object")
		require.NoError(t, err)
		_, err = client.VerifyToken(ctx, sig.String(), target)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrTokenInvalid))
	})

	t.Run("rejects a token without a key ID", func(t *testing.T) {
		now := time.Now().UTC()
		token := signWithoutKID(t, &TokenClaims{
			Iss: requester,
			Aud: target,
			Iat: now.Unix(),
			Exp: now.Add(time.Minute).Unix(),
		})
		_, err := client.VerifyToken(ctx, token, target)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrTokenInvalid))
		assert.Contains(t, err.Error(), "no key ID")
	})

	t.Run("rejects a signature that fails against the published key", func(t *testing.T) {
		// The published JWK claims the signer's kid but holds another
		// key's material, so the fetch succeeds and the crypto fails.
		other := dsig.NewES256Key()
		c := NewClient(WithFetcher(&mapFetcher{data: map[string][]byte{
			requester.KeyURL(key.ID()): spoofKID(t, other, key.ID()),
		}}))
		token, err := NewToken(key, requester, target, 0)
		require.NoError(t, err)
		_, err = c.VerifyToken(ctx, token, target)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrTokenInvalid))
	})

	t.Run("rejects a token outside the key's validity window", func(t *testing.T) {
		// The published key retired an hour ago; a token minted now
		// carries an iat the key no longer allows.
		data, err := json.Marshal(key.Public())
		require.NoError(t, err)
		pk := new(dsig.PublicKey)
		require.NoError(t, json.Unmarshal(data, pk))
		until := cal.TimestampOf(time.Now().Add(-time.Hour))
		pk.ValidUntil = &until
		retired, err := json.Marshal(pk)
		require.NoError(t, err)

		c := NewClient(WithFetcher(&mapFetcher{data: map[string][]byte{
			requester.KeyURL(key.ID()): retired,
		}}))
		token, err := NewToken(key, requester, target, 0)
		require.NoError(t, err)
		_, err = c.VerifyToken(ctx, token, target)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrTokenInvalid))
	})
}
