package net

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/note"
	"github.com/invopop/gobl/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterAuthority(t *testing.T) {
	original := Authorities
	t.Cleanup(func() { Authorities = original })
	Authorities = nil

	RegisterAuthority("kyc.example.com")
	RegisterAuthority("auth.example.org")

	assert.Equal(t, []Address{"kyc.example.com", "auth.example.org"}, Authorities)
}

func TestNewClientAuthoritiesIndependent(t *testing.T) {
	original := Authorities
	t.Cleanup(func() { Authorities = original })
	Authorities = []Address{"kyc.example.com"}

	c := NewClient()
	assert.Equal(t, []Address{"kyc.example.com"}, c.authorities)

	// Mutating the global after construction must not affect the client.
	Authorities = append(Authorities, "auth.example.org")
	assert.Equal(t, []Address{"kyc.example.com"}, c.authorities)
}

func TestWithAuthorities(t *testing.T) {
	original := Authorities
	t.Cleanup(func() { Authorities = original })
	Authorities = nil

	c := NewClient(WithAuthorities("kyc.example.com", "auth.example.org"))
	assert.Equal(t, []Address{"kyc.example.com", "auth.example.org"}, c.authorities)
}

// mapFetcher is a URL-keyed Fetcher for tests that need to serve
// different blobs for different URLs (e.g. multiple per-key
// endpoints).
type mapFetcher struct {
	data map[string][]byte
}

func (m *mapFetcher) Fetch(_ context.Context, url string) ([]byte, error) {
	if d, ok := m.data[url]; ok {
		return d, nil
	}
	return nil, ErrFetchFailed
}

func TestVerifyAuthority(t *testing.T) {
	ctx := context.Background()
	authorityAddr := Address("kyc.example.com")
	subjectAddr := Address("subject.example")
	authKey := dsig.NewES256Key()
	subjKey := dsig.NewES256Key()

	jwkOf := func(k *dsig.PrivateKey) []byte {
		t.Helper()
		out, err := json.Marshal(k.Public())
		require.NoError(t, err)
		return out
	}

	t.Run("verifies an authority-signed envelope", func(t *testing.T) {
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		// Subject's self-signature (not an authority) + authority countersignature.
		require.NoError(t, env.Sign(subjKey, subjectAddr.URI(), ""))
		require.NoError(t, env.Sign(authKey,
			authorityAddr.URI(),
			subjectAddr.URI(),
			head.WithScope(head.ScopeVerified)))

		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
			}}),
		)
		assert.NoError(t, c.VerifyAuthority(ctx, env))
	})

	t.Run("rejects an envelope with only a self-signature", func(t *testing.T) {
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		require.NoError(t, env.Sign(subjKey, subjectAddr.URI(), ""))

		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{}}),
		)
		err = c.VerifyAuthority(ctx, env)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrUnknownAuthority))
	})

	t.Run("rejects an envelope with no signatures at all", func(t *testing.T) {
		msg := &note.Message{Content: "x"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)

		c := NewClient(WithAuthorities(authorityAddr))
		err = c.VerifyAuthority(ctx, env)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrVerifyFailed))
	})

	t.Run("rejects when client has no authorities registered", func(t *testing.T) {
		msg := &note.Message{Content: "x"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		require.NoError(t, env.Sign(authKey, authorityAddr.URI(), ""))

		c := NewClient() // empty authorities
		err = c.VerifyAuthority(ctx, env)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrUnknownAuthority))
	})

	t.Run("rejects when authority signature does not verify against published key", func(t *testing.T) {
		// The envelope claims to be signed by the authority, but the
		// fetcher returns a different key, so the crypto verification
		// fails. ErrVerifyFailed is returned (not ErrUnknownAuthority).
		msg := &note.Message{Content: "x"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		require.NoError(t, env.Sign(authKey, authorityAddr.URI(), ""))

		other := dsig.NewES256Key()
		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				// Wrong key for the authority's claimed kid.
				authorityAddr.KeyURL(authKey.ID()): jwkOf(other),
			}}),
		)
		err = c.VerifyAuthority(ctx, env)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrVerifyFailed))
	})

	t.Run("skips signatures with non-gobl iss schemes", func(t *testing.T) {
		// A signature whose iss isn't a gobl: URI cannot be an
		// authority signature — VerifyAuthority just steps past it.
		msg := &note.Message{Content: "x"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		require.NoError(t, env.Sign(subjKey, cbc.URI("mailto:a@b"), ""))

		c := NewClient(WithAuthorities(authorityAddr))
		err = c.VerifyAuthority(ctx, env)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrUnknownAuthority))
	})
}
