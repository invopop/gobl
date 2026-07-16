package net

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

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
		require.NoError(t, env.Sign(subjKey, head.WithIssuer(subjectAddr.URI())))
		require.NoError(t, env.Sign(authKey,
			head.WithIssuer(authorityAddr.URI()),
			head.WithAudience(subjectAddr.URI()),
			head.WithScope(head.ScopeVerified)))

		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
			}}),
		)
		assert.NoError(t, c.VerifyAuthority(ctx, env))
	})

	t.Run("enforces a minimum scope", func(t *testing.T) {
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		require.NoError(t, env.Sign(subjKey, head.WithIssuer(subjectAddr.URI())))
		require.NoError(t, env.Sign(authKey,
			head.WithIssuer(authorityAddr.URI()),
			head.WithAudience(subjectAddr.URI()),
			head.WithScope(head.ScopeRegistered)))

		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
			}}),
		)
		// registered satisfies registered but not verified.
		assert.NoError(t, c.VerifyAuthorityWithScope(ctx, env, head.ScopeRegistered))
		err = c.VerifyAuthorityWithScope(ctx, env, head.ScopeVerified)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrScopeInsufficient))
	})

	t.Run("verified satisfies a registered minimum", func(t *testing.T) {
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		require.NoError(t, env.Sign(authKey,
			head.WithIssuer(authorityAddr.URI()),
			head.WithScope(head.ScopeVerified)))

		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
			}}),
		)
		assert.NoError(t, c.VerifyAuthorityWithScope(ctx, env, head.ScopeRegistered))
	})

	t.Run("scopeless authority signature fails any named minimum", func(t *testing.T) {
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		require.NoError(t, env.Sign(authKey, head.WithIssuer(authorityAddr.URI())))

		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
			}}),
		)
		assert.NoError(t, c.VerifyAuthority(ctx, env))
		err = c.VerifyAuthorityWithScope(ctx, env, head.ScopeRegistered)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrScopeInsufficient))
	})

	t.Run("custom scope requires an exact match", func(t *testing.T) {
		custom := cbc.Key("premium")
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		require.NoError(t, env.Sign(authKey,
			head.WithIssuer(authorityAddr.URI()),
			head.WithScope(custom)))

		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
			}}),
		)
		assert.NoError(t, c.VerifyAuthorityWithScope(ctx, env, custom))
		err = c.VerifyAuthorityWithScope(ctx, env, head.ScopeRegistered)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrScopeInsufficient))
	})

	t.Run("rejects an expired countersignature", func(t *testing.T) {
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		require.NoError(t, env.Sign(authKey,
			head.WithIssuer(authorityAddr.URI()),
			head.WithScope(head.ScopeVerified),
			head.WithExpiration(time.Now().Add(-time.Hour))))

		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
			}}),
		)
		err = c.VerifyAuthority(ctx, env)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrSignatureExpired))
	})

	t.Run("rejects a pre-epoch expiry", func(t *testing.T) {
		// A negative exp is a valid NumericDate before 1970 and must be
		// enforced, not treated as an absent claim.
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		require.NoError(t, env.Sign(authKey,
			head.WithIssuer(authorityAddr.URI()),
			head.WithScope(head.ScopeVerified),
			head.WithExpiration(time.Unix(-1000, 0))))

		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
			}}),
		)
		err = c.VerifyAuthority(ctx, env)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrSignatureExpired))
	})

	t.Run("expiry outlives later candidate failures", func(t *testing.T) {
		// An expired countersignature followed by a candidate whose key
		// cannot be fetched must still surface ErrSignatureExpired, not
		// the generic ErrVerifyFailed from the later failure.
		otherAuth := Address("auth.example.org")
		otherKey := dsig.NewES256Key()
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		require.NoError(t, env.Sign(authKey,
			head.WithIssuer(authorityAddr.URI()),
			head.WithScope(head.ScopeVerified),
			head.WithExpiration(time.Now().Add(-time.Hour))))
		require.NoError(t, env.Sign(otherKey, head.WithIssuer(otherAuth.URI())))

		c := NewClient(
			WithAuthorities(authorityAddr, otherAuth),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				// Only the first authority's key resolves.
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
			}}),
		)
		err = c.VerifyAuthority(ctx, env)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrSignatureExpired))
	})

	t.Run("accepts a countersignature within its exp window", func(t *testing.T) {
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		require.NoError(t, env.Sign(authKey,
			head.WithIssuer(authorityAddr.URI()),
			head.WithScope(head.ScopeVerified),
			head.WithExpiration(time.Now().Add(90*24*time.Hour))))

		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
			}}),
		)
		assert.NoError(t, c.VerifyAuthorityWithScope(ctx, env, head.ScopeVerified))
	})

	t.Run("rejects an envelope with only a self-signature", func(t *testing.T) {
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		require.NoError(t, env.Sign(subjKey, head.WithIssuer(subjectAddr.URI())))

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
		require.NoError(t, env.Sign(authKey, head.WithIssuer(authorityAddr.URI())))

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
		require.NoError(t, env.Sign(authKey, head.WithIssuer(authorityAddr.URI())))

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
		require.NoError(t, env.Sign(subjKey, head.WithIssuer("mailto:a@b")))

		c := NewClient(WithAuthorities(authorityAddr))
		err = c.VerifyAuthority(ctx, env)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrUnknownAuthority))
	})
}
