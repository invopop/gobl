package net

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/note"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// buildPartyEnvelope returns the JSON of a party envelope self-signed
// by subject, optionally countersigned by an authority naming a
// verifier.
func buildPartyEnvelope(t *testing.T, subjKey *dsig.PrivateKey, subject Address, authKey *dsig.PrivateKey, authority, verifier Address) []byte {
	t.Helper()
	party := &org.Party{
		Name:      "Test Subject",
		Endpoints: []*org.Endpoint{{URI: subject.URI()}},
	}
	party.SetUUID(uuid.V7())
	env, err := gobl.Envelop(party)
	require.NoError(t, err)
	require.NoError(t, env.Sign(subjKey, head.WithIssuer(subject.String())))
	if authKey != nil {
		opts := []head.SignOption{
			head.WithIssuer(authority.String()),
			head.WithAudience(subject.String()),
		}
		if verifier != "" {
			opts = append(opts, head.WithVerifier(verifier.String()))
		}
		require.NoError(t, env.Sign(authKey, opts...))
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

	t.Run("canonicalizes a non-canonical address", func(t *testing.T) {
		// Mixed case and a trailing dot must resolve to the same
		// canonical who URL and match the signed issuer.
		c := NewClient(WithFetcher(&mapFetcher{data: map[string][]byte{
			subject.WhoURL():             buildPartyEnvelope(t, subjKey, subject, nil, "", ""),
			subject.KeyURL(subjKey.ID()): jwkOf(subjKey),
		}}))
		env, err := c.Who(ctx, Address("SUPPLIER.Example.COM."))
		require.NoError(t, err)
		assert.NotNil(t, env)
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
		assert.Contains(t, err.Error(), "identity of")
	})

	t.Run("rejects an audience-bound response", func(t *testing.T) {
		// A who response is a public document; an envelope signed with
		// an aud is caller-bound and non-conforming.
		party := &org.Party{
			Name:      "Bound",
			Endpoints: []*org.Endpoint{{URI: subject.URI()}},
		}
		party.SetUUID(uuid.V7())
		env, err := gobl.Envelop(party)
		require.NoError(t, err)
		require.NoError(t, env.Sign(subjKey,
			head.WithIssuer(subject.String()),
			head.WithAudience("caller.example.com")))
		data, err := json.Marshal(env)
		require.NoError(t, err)

		c := NewClient(WithFetcher(&mapFetcher{data: map[string][]byte{
			subject.WhoURL():             data,
			subject.KeyURL(subjKey.ID()): jwkOf(subjKey),
		}}))
		_, err = c.Who(ctx, subject)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrVerifyFailed))
		assert.Contains(t, err.Error(), "audience-free")
	})

	t.Run("accepts hop-bound self-signatures alongside the publication one", func(t *testing.T) {
		// A published endorsed envelope: the subject's registration
		// signature is audience-bound to the registry, an authority
		// countersignature is aboard, and an audience-free
		// self-signature asserts the public identity. Signature order
		// is not significant.
		party := &org.Party{
			Name:      "Endorsed",
			Endpoints: []*org.Endpoint{{URI: subject.URI()}},
		}
		party.SetUUID(uuid.V7())
		env, err := gobl.Envelop(party)
		require.NoError(t, err)
		require.NoError(t, env.Sign(subjKey,
			head.WithIssuer(subject.String()),
			head.WithAudience("lookup.example.com")))
		authKey := dsig.NewES256Key()
		require.NoError(t, env.Sign(authKey,
			head.WithIssuer("lookup.example.com"),
			head.WithAudience(subject.String())))
		require.NoError(t, env.Sign(subjKey, head.WithIssuer(subject.String())))
		data, err := json.Marshal(env)
		require.NoError(t, err)

		c := NewClient(WithFetcher(&mapFetcher{data: map[string][]byte{
			subject.WhoURL():             data,
			subject.KeyURL(subjKey.ID()): jwkOf(subjKey),
		}}))
		got, err := c.Who(ctx, subject)
		require.NoError(t, err)
		require.Len(t, got.Signatures, 3, "countersignatures preserved for VerifyAuthority")
	})

	t.Run("rejects a non-party document", func(t *testing.T) {
		msg := &note.Message{Content: "not a party"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		require.NoError(t, env.Sign(subjKey, head.WithIssuer(subject.String())))
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

	t.Run("rejects an invalid address", func(t *testing.T) {
		c := NewClient()
		_, err := c.Who(ctx, Address("localhost"))
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrAddressInvalid))
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
	verifier := Address("verify.example.com")
	subjKey := dsig.NewES256Key()
	authKey := dsig.NewES256Key()
	verifierKey := dsig.NewES256Key()

	jwkOf := func(k *dsig.PrivateKey) []byte {
		t.Helper()
		out, err := json.Marshal(k.Public())
		require.NoError(t, err)
		return out
	}

	endorsedEnvelope := new(gobl.Envelope)
	require.NoError(t, json.Unmarshal(
		buildPartyEnvelope(t, subjKey, subject, authKey, authority, verifier),
		endorsedEnvelope,
	))
	require.NoError(t, endorsedEnvelope.Sign(verifierKey,
		head.WithIssuer(verifier.String()),
		head.WithAudience(subject.String())))
	endorsedData, err := json.Marshal(endorsedEnvelope)
	require.NoError(t, err)

	endorsed := map[string][]byte{
		subject.WhoURL():                  endorsedData,
		subject.KeyURL(subjKey.ID()):      jwkOf(subjKey),
		authority.KeyURL(authKey.ID()):    jwkOf(authKey),
		verifier.KeyURL(verifierKey.ID()): jwkOf(verifierKey),
	}

	t.Run("accepts an authority-endorsed sender", func(t *testing.T) {
		c := NewClient(
			WithAuthorities(authority),
			WithFetcher(&mapFetcher{data: endorsed}),
		)
		party, err := c.VerifySender(ctx, subject, true)
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
		_, err := c.VerifySender(ctx, subject, false)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrUnknownAuthority))
	})

	t.Run("rejects an unverified sender when verification is required", func(t *testing.T) {
		registered := map[string][]byte{
			subject.WhoURL():               buildPartyEnvelope(t, subjKey, subject, authKey, authority, ""),
			subject.KeyURL(subjKey.ID()):   jwkOf(subjKey),
			authority.KeyURL(authKey.ID()): jwkOf(authKey),
		}
		c := NewClient(
			WithAuthorities(authority),
			WithFetcher(&mapFetcher{data: registered}),
		)
		_, err := c.VerifySender(ctx, subject, true)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrNotVerified))

		// Registration alone satisfies a caller that does not require
		// verification.
		party, err := c.VerifySender(ctx, subject, false)
		require.NoError(t, err)
		assert.Equal(t, "Test Subject", party.Name)
	})

	t.Run("rejects a receive-only account", func(t *testing.T) {
		c := NewClient(
			WithAuthorities(authority),
			WithFetcher(&mockFetcher{err: ErrNoContent}),
		)
		_, err := c.VerifySender(ctx, subject, true)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrNoContent))
	})
}
