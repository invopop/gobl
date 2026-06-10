package net

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/note"
	"github.com/invopop/gobl/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildTestEnvelope(t *testing.T, key *dsig.PrivateKey, iss, aud cbc.URI) *gobl.Envelope {
	t.Helper()

	msg := &note.Message{Content: "test message content"}
	msg.SetUUID(uuid.V7())

	env, err := gobl.Envelop(msg)
	require.NoError(t, err)
	require.NoError(t, env.Sign(key, head.WithIssuer(iss), head.WithAudience(aud)))

	// Round-trip through JSON to simulate realistic usage.
	data, err := json.Marshal(env)
	require.NoError(t, err)
	env = new(gobl.Envelope)
	require.NoError(t, json.Unmarshal(data, env))
	return env
}

// jwkFromKey returns the single-JWK JSON bytes served at the per-key
// endpoint for this key.
func jwkFromKey(t *testing.T, key *dsig.PrivateKey) []byte {
	t.Helper()
	pubData, err := json.Marshal(key.Public())
	require.NoError(t, err)
	return pubData
}

func TestVerifyEnvelope(t *testing.T) {
	ctx := context.Background()
	addr := Address("billing.invopop.com")

	t.Run("success", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := buildTestEnvelope(t, key, addr.URI(), "")

		mock := &mockFetcher{data: jwkFromKey(t, key)}
		c := NewClient(WithFetcher(mock))

		issuer, err := c.VerifyEnvelope(ctx, env, "")
		require.NoError(t, err)
		assert.Equal(t, addr, issuer)
		assert.Equal(t, addr.KeyURL(key.ID()), mock.url)
	})

	t.Run("audience match", func(t *testing.T) {
		key := dsig.NewES256Key()
		aud := Address("recipient.example.com")
		env := buildTestEnvelope(t, key, addr.URI(), aud.URI())

		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
		issuer, err := c.VerifyEnvelope(ctx, env, aud.URI())
		require.NoError(t, err)
		assert.Equal(t, addr, issuer)
	})

	t.Run("audience mismatch", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := buildTestEnvelope(t, key, addr.URI(), Address("a.example").URI())

		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
		_, err := c.VerifyEnvelope(ctx, env, Address("b.example").URI())
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrVerifyFailed))
	})

	t.Run("not signed", func(t *testing.T) {
		c := NewClient()
		_, err := c.VerifyEnvelope(ctx, new(gobl.Envelope), "")
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrVerifyFailed))
	})

	t.Run("no iss", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := buildTestEnvelope(t, key, "", "") // signed without an iss
		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
		_, err := c.VerifyEnvelope(ctx, env, "")
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrVerifyFailed))
	})

	t.Run("key not found", func(t *testing.T) {
		// The per-key endpoint returns a JWK whose kid does not match the
		// signature's kid, so the client rejects the response.
		key := dsig.NewES256Key()
		other := dsig.NewES256Key()
		env := buildTestEnvelope(t, key, addr.URI(), "")
		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, other)}))
		_, err := c.VerifyEnvelope(ctx, env, "")
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrVerifyFailed))
	})

	t.Run("signing time before valid_from", func(t *testing.T) {
		// Sign now, then publish the key with a valid_from in the future.
		key := dsig.NewES256Key()
		env := buildTestEnvelope(t, key, addr.URI(), "")
		future := cal.TimestampOf(time.Now().Add(24 * time.Hour))
		c := NewClient(WithFetcher(&mockFetcher{data: publishedJWK(t, key, &future, nil)}))
		_, err := c.VerifyEnvelope(ctx, env, "")
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrVerifyFailed))
		assert.Contains(t, err.Error(), "before key's valid_from")
	})

	t.Run("signing time after valid_until", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := buildTestEnvelope(t, key, addr.URI(), "")
		past := cal.TimestampOf(time.Now().Add(-24 * time.Hour))
		c := NewClient(WithFetcher(&mockFetcher{data: publishedJWK(t, key, nil, &past)}))
		_, err := c.VerifyEnvelope(ctx, env, "")
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrVerifyFailed))
		assert.Contains(t, err.Error(), "after key's valid_until")
	})

	t.Run("signing time inside window", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := buildTestEnvelope(t, key, addr.URI(), "")
		from := cal.TimestampOf(time.Now().Add(-time.Hour))
		until := cal.TimestampOf(time.Now().Add(time.Hour))
		c := NewClient(WithFetcher(&mockFetcher{data: publishedJWK(t, key, &from, &until)}))
		issuer, err := c.VerifyEnvelope(ctx, env, "")
		require.NoError(t, err)
		assert.Equal(t, addr, issuer)
	})
}

func TestVerifyEnvelopePayloadErrors(t *testing.T) {
	ctx := context.Background()

	t.Run("non-gobl iss scheme", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := buildTestEnvelope(t, key, cbc.URI("mailto:a@example.com"), "")
		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
		_, err := c.VerifyEnvelope(ctx, env, "")
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrVerifyFailed))
		assert.Contains(t, err.Error(), "is not a gobl address")
	})

	t.Run("invalid iss FQDN", func(t *testing.T) {
		key := dsig.NewES256Key()
		// "localhost" is a single label — fails FQDN validation.
		env := buildTestEnvelope(t, key, cbc.URI("gobl:localhost"), "")
		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
		_, err := c.VerifyEnvelope(ctx, env, "")
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrVerifyFailed))
	})
}

// publishedJWK marshals key.Public() as a dsig.PublicKey JSON object
// with the supplied (optional) validity bounds.
func publishedJWK(t *testing.T, key *dsig.PrivateKey, from, until *cal.Timestamp) []byte {
	t.Helper()
	data, err := json.Marshal(key.Public())
	require.NoError(t, err)
	pk := new(dsig.PublicKey)
	require.NoError(t, json.Unmarshal(data, pk))
	pk.ValidFrom = from
	pk.ValidUntil = until
	out, err := json.Marshal(pk)
	require.NoError(t, err)
	return out
}
