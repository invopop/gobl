package net

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/note"
	"github.com/invopop/gobl/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildTestEnvelope(t *testing.T, key *dsig.PrivateKey, addr string) *gobl.Envelope {
	t.Helper()

	msg := &note.Message{
		Content: "test message content",
	}
	msg.SetUUID(uuid.V7())

	env, err := gobl.Envelop(msg)
	require.NoError(t, err)

	var opts []dsig.SignerOption
	if addr != "" {
		opts = append(opts, dsig.WithGN(addr))
	}
	err = env.Sign(key, opts...)
	require.NoError(t, err)

	// Round-trip through JSON to simulate realistic usage.
	// JWS extra headers are only accessible after parse.
	data, err := json.Marshal(env)
	require.NoError(t, err)
	env = new(gobl.Envelope)
	require.NoError(t, json.Unmarshal(data, env))

	return env
}

func jwksFromKey(t *testing.T, key *dsig.PrivateKey) []byte {
	t.Helper()
	pub := key.Public()
	pubData, err := json.Marshal(pub)
	require.NoError(t, err)
	return []byte(`{"keys":[` + string(pubData) + `]}`)
}

func TestVerifyEnvelope(t *testing.T) {
	ctx := context.Background()
	addr := Address("billing.invopop.com")

	t.Run("success with gn header", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := buildTestEnvelope(t, key, addr.String())

		mock := &mockFetcher{data: jwksFromKey(t, key)}
		c := NewClient(WithFetcher(mock))

		err := c.VerifyEnvelope(ctx, env, "")
		require.NoError(t, err)
		assert.Equal(t, addr.JWKSURL(), mock.url)
	})

	t.Run("success with explicit address", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := buildTestEnvelope(t, key, "") // no gn header

		mock := &mockFetcher{data: jwksFromKey(t, key)}
		c := NewClient(WithFetcher(mock))

		err := c.VerifyEnvelope(ctx, env, addr)
		require.NoError(t, err)
	})

	t.Run("not signed", func(t *testing.T) {
		env := new(gobl.Envelope)
		c := NewClient()

		err := c.VerifyEnvelope(ctx, env, addr)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrVerifyFailed))
	})

	t.Run("no gn header and no address", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := buildTestEnvelope(t, key, "") // no gn header

		c := NewClient()

		err := c.VerifyEnvelope(ctx, env, "")
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrNoGNHeader))
	})

	t.Run("key not found", func(t *testing.T) {
		key := dsig.NewES256Key()
		otherKey := dsig.NewES256Key()
		env := buildTestEnvelope(t, key, addr.String())

		// Serve a JWKS with a different key
		mock := &mockFetcher{data: jwksFromKey(t, otherKey)}
		c := NewClient(WithFetcher(mock))

		err := c.VerifyEnvelope(ctx, env, "")
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrVerifyFailed))
	})

	t.Run("wrong key same kid", func(t *testing.T) {
		key := dsig.NewES256Key()
		otherKey := dsig.NewES256Key()
		env := buildTestEnvelope(t, key, addr.String())

		// Build JWKS with the other key but using the signing key's ID
		pub := otherKey.Public()
		pubData, err := json.Marshal(pub)
		require.NoError(t, err)

		var m map[string]any
		require.NoError(t, json.Unmarshal(pubData, &m))
		m["kid"] = key.ID()
		fakeData, err := json.Marshal(m)
		require.NoError(t, err)
		jwks := []byte(`{"keys":[` + string(fakeData) + `]}`)

		mock := &mockFetcher{data: jwks}
		c := NewClient(WithFetcher(mock))

		err = c.VerifyEnvelope(ctx, env, "")
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrVerifyFailed))
	})
}
