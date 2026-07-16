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
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// buildPartyEnvelope returns the JSON of a party envelope self-signed
// by subject, optionally countersigned by an authority at scope.
func buildPartyEnvelope(t *testing.T, subjKey *dsig.PrivateKey, subject Address, authKey *dsig.PrivateKey, authority Address, scope cbc.Key) []byte {
	t.Helper()
	party := &org.Party{Name: "Test Subject"}
	party.SetUUID(uuid.V7())
	env, err := gobl.Envelop(party)
	require.NoError(t, err)
	require.NoError(t, env.Sign(subjKey, head.WithIssuer(subject.URI())))
	if authKey != nil {
		require.NoError(t, env.Sign(authKey,
			head.WithIssuer(authority.URI()),
			head.WithAudience(subject.URI()),
			head.WithScope(scope)))
	}
	data, err := json.Marshal(env)
	require.NoError(t, err)
	return data
}

func TestWho(t *testing.T) {
	ctx := context.Background()
	subject := Address("supplier.example.com")
	subjKey := dsig.NewES256Key()

	jwkOf := func(k *dsig.PrivateKey) []byte {
		t.Helper()
		out, err := json.Marshal(k.Public())
		require.NoError(t, err)
		return out
	}

	t.Run("fetches and verifies a self-signed identity", func(t *testing.T) {
		c := NewClient(WithFetcher(&mapFetcher{data: map[string][]byte{
			subject.WhoURL():             buildPartyEnvelope(t, subjKey, subject, nil, "", ""),
			subject.KeyURL(subjKey.ID()): jwkOf(subjKey),
		}}))
		env, err := c.Who(ctx, subject)
		require.NoError(t, err)
		party, ok := env.Extract().(*org.Party)
		require.True(t, ok)
		assert.Equal(t, "Test Subject", party.Name)
	})

	t.Run("rejects an issuer that does not match the address", func(t *testing.T) {
		other := Address("other.example.com")
		c := NewClient(WithFetcher(&mapFetcher{data: map[string][]byte{
			// Served from subject's who URL, but signed as other.
			subject.WhoURL():           buildPartyEnvelope(t, subjKey, other, nil, "", ""),
			other.KeyURL(subjKey.ID()): jwkOf(subjKey),
		}}))
		_, err := c.Who(ctx, subject)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrVerifyFailed))
		assert.Contains(t, err.Error(), "does not match address")
	})

	t.Run("rejects a non-party document", func(t *testing.T) {
		msg := &note.Message{Content: "not a party"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		require.NoError(t, env.Sign(subjKey, head.WithIssuer(subject.URI())))
		data, err := json.Marshal(env)
		require.NoError(t, err)

		c := NewClient(WithFetcher(&mapFetcher{data: map[string][]byte{
			subject.WhoURL():             data,
			subject.KeyURL(subjKey.ID()): jwkOf(subjKey),
		}}))
		_, err = c.Who(ctx, subject)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrPartyMissing))
	})

	t.Run("propagates a no-content response", func(t *testing.T) {
		c := NewClient(WithFetcher(&mockFetcher{err: ErrNoContent}))
		_, err := c.Who(ctx, subject)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrNoContent))
	})

	t.Run("rejects an unsigned envelope", func(t *testing.T) {
		party := &org.Party{Name: "Unsigned"}
		party.SetUUID(uuid.V7())
		env, err := gobl.Envelop(party)
		require.NoError(t, err)
		data, err := json.Marshal(env)
		require.NoError(t, err)

		c := NewClient(WithFetcher(&mapFetcher{data: map[string][]byte{
			subject.WhoURL(): data,
		}}))
		_, err = c.Who(ctx, subject)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrVerifyFailed))
	})

	t.Run("rejects a body that is not an envelope", func(t *testing.T) {
		c := NewClient(WithFetcher(&mapFetcher{data: map[string][]byte{
			subject.WhoURL(): []byte("not json"),
		}}))
		_, err := c.Who(ctx, subject)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrFetchFailed))
	})
}

func TestVerifySender(t *testing.T) {
	ctx := context.Background()
	subject := Address("supplier.example.com")
	authority := Address("kyc.example.com")
	subjKey := dsig.NewES256Key()
	authKey := dsig.NewES256Key()

	jwkOf := func(k *dsig.PrivateKey) []byte {
		t.Helper()
		out, err := json.Marshal(k.Public())
		require.NoError(t, err)
		return out
	}

	endorsed := map[string][]byte{
		subject.WhoURL():               buildPartyEnvelope(t, subjKey, subject, authKey, authority, head.ScopeVerified),
		subject.KeyURL(subjKey.ID()):   jwkOf(subjKey),
		authority.KeyURL(authKey.ID()): jwkOf(authKey),
	}

	t.Run("accepts an authority-endorsed sender", func(t *testing.T) {
		c := NewClient(
			WithAuthorities(authority),
			WithFetcher(&mapFetcher{data: endorsed}),
		)
		party, err := c.VerifySender(ctx, subject, head.ScopeVerified)
		require.NoError(t, err)
		assert.Equal(t, "Test Subject", party.Name)
	})

	t.Run("rejects a self-signed sender", func(t *testing.T) {
		c := NewClient(
			WithAuthorities(authority),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				subject.WhoURL():             buildPartyEnvelope(t, subjKey, subject, nil, "", ""),
				subject.KeyURL(subjKey.ID()): jwkOf(subjKey),
			}}),
		)
		_, err := c.VerifySender(ctx, subject, "")
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrUnknownAuthority))
	})

	t.Run("rejects an endorsement below the minimum scope", func(t *testing.T) {
		c := NewClient(
			WithAuthorities(authority),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				subject.WhoURL():               buildPartyEnvelope(t, subjKey, subject, authKey, authority, head.ScopeRegistered),
				subject.KeyURL(subjKey.ID()):   jwkOf(subjKey),
				authority.KeyURL(authKey.ID()): jwkOf(authKey),
			}}),
		)
		_, err := c.VerifySender(ctx, subject, head.ScopeVerified)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrScopeInsufficient))
	})

	t.Run("rejects a receive-only account", func(t *testing.T) {
		c := NewClient(
			WithAuthorities(authority),
			WithFetcher(&mockFetcher{err: ErrNoContent}),
		)
		_, err := c.VerifySender(ctx, subject, head.ScopeVerified)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrNoContent))
	})
}
