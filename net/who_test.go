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

// testRequester is the identity test clients present to /who; party
// envelopes served to it carry it as their to address.
var (
	testRequester    = Address("requester.example.com")
	testRequesterKey = dsig.NewES256Key()
)

// buildPartyEnvelope returns the JSON of a party envelope from subject
// to testRequester, signed for it by subject, optionally countersigned by
// an authority naming a verifier.
func buildPartyEnvelope(t *testing.T, subjKey *dsig.PrivateKey, subject Address, authKey *dsig.PrivateKey, authority, verifier Address) []byte {
	t.Helper()
	party := &org.Party{
		Name:      "Test Subject",
		Endpoints: []*org.Endpoint{{URI: subject.URI()}},
	}
	party.SetUUID(uuid.V7())
	env, err := gobl.Envelop(party)
	require.NoError(t, err)
	env.Head.From = subject.URI()
	env.Head.To = testRequester.URI()
	require.NoError(t, env.Sign(subjKey,
		head.WithIssuer(subject.String()),
		head.WithAudience(testRequester.String())))
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
		c := NewClient(WithIdentity(testRequester, testRequesterKey), WithFetcher(&mapFetcher{data: map[string][]byte{
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
		c := NewClient(WithIdentity(testRequester, testRequesterKey), WithFetcher(&mapFetcher{data: map[string][]byte{
			subject.WhoURL():             buildPartyEnvelope(t, subjKey, subject, nil, "", ""),
			subject.KeyURL(subjKey.ID()): jwkOf(subjKey),
		}}))
		env, err := c.Who(ctx, Address("SUPPLIER.Example.COM."))
		require.NoError(t, err)
		assert.NotNil(t, env)
	})

	t.Run("rejects a response from another address", func(t *testing.T) {
		other := Address("other.example.com")
		c := NewClient(WithIdentity(testRequester, testRequesterKey), WithFetcher(&mapFetcher{data: map[string][]byte{
			// Served from subject's who URL, but from and signed by other.
			subject.WhoURL():           buildPartyEnvelope(t, subjKey, other, nil, "", ""),
			other.KeyURL(subjKey.ID()): jwkOf(subjKey),
		}}))
		_, err := c.Who(ctx, subject)
		require.ErrorIs(t, err, ErrVerifyFailed)
		assert.Contains(t, err.Error(), "is from")
	})

	t.Run("rejects a response not signed for the requester", func(t *testing.T) {
		// Signatures without an aud, or for another audience, do not
		// confirm the from and to pair.
		for _, aud := range []string{"", "caller.example.com"} {
			party := &org.Party{Name: "Bound"}
			party.SetUUID(uuid.V7())
			env, err := gobl.Envelop(party)
			require.NoError(t, err)
			env.Head.From = subject.URI()
			env.Head.To = testRequester.URI()
			opts := []head.SignOption{head.WithIssuer(subject.String())}
			if aud != "" {
				opts = append(opts, head.WithAudience(aud))
			}
			require.NoError(t, env.Sign(subjKey, opts...))
			data, err := json.Marshal(env)
			require.NoError(t, err)

			c := NewClient(WithIdentity(testRequester, testRequesterKey), WithFetcher(&mapFetcher{data: map[string][]byte{
				subject.WhoURL():             data,
				subject.KeyURL(subjKey.ID()): jwkOf(subjKey),
			}}))
			_, err = c.Who(ctx, subject)
			require.ErrorIs(t, err, ErrVerifyFailed)
			assert.Contains(t, err.Error(), "no valid signature")
		}
	})

	t.Run("accepts other signatures alongside the requester's", func(t *testing.T) {
		// A published endorsed envelope: the subject's registration
		// signature is bound to the authority, an authority
		// countersignature is aboard, and the response is signed for
		// the requester. Signature order is not significant.
		party := &org.Party{
			Name:      "Endorsed",
			Endpoints: []*org.Endpoint{{URI: subject.URI()}},
		}
		party.SetUUID(uuid.V7())
		env, err := gobl.Envelop(party)
		require.NoError(t, err)
		env.Head.From = subject.URI()
		env.Head.To = testRequester.URI()
		require.NoError(t, env.Sign(subjKey,
			head.WithIssuer(subject.String()),
			head.WithAudience("lookup.example.com")))
		authKey := dsig.NewES256Key()
		require.NoError(t, env.Sign(authKey,
			head.WithIssuer("lookup.example.com"),
			head.WithAudience(subject.String())))
		require.NoError(t, env.Sign(subjKey,
			head.WithIssuer(subject.String()),
			head.WithAudience(testRequester.String())))
		data, err := json.Marshal(env)
		require.NoError(t, err)

		c := NewClient(WithIdentity(testRequester, testRequesterKey), WithFetcher(&mapFetcher{data: map[string][]byte{
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
		env.Head.From = subject.URI()
		env.Head.To = testRequester.URI()
		require.NoError(t, env.Sign(subjKey,
			head.WithIssuer(subject.String()),
			head.WithAudience(testRequester.String())))
		data, err := json.Marshal(env)
		require.NoError(t, err)

		c := NewClient(WithIdentity(testRequester, testRequesterKey), WithFetcher(&mapFetcher{data: map[string][]byte{
			subject.WhoURL():             data,
			subject.KeyURL(subjKey.ID()): jwkOf(subjKey),
		}}))
		_, err = c.Who(ctx, subject)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrPartyMissing))
	})

	t.Run("rejects a response for another requester", func(t *testing.T) {
		env := new(gobl.Envelope)
		require.NoError(t, json.Unmarshal(buildPartyEnvelope(t, subjKey, subject, nil, "", ""), env))
		env.Head.To = "gobl:someone-else.example.com"
		data, err := json.Marshal(env)
		require.NoError(t, err)

		c := NewClient(WithIdentity(testRequester, testRequesterKey), WithFetcher(&mapFetcher{data: map[string][]byte{
			subject.WhoURL():             data,
			subject.KeyURL(subjKey.ID()): jwkOf(subjKey),
		}}))
		_, err = c.Who(ctx, subject)
		require.ErrorIs(t, err, ErrVerifyFailed)
		assert.Contains(t, err.Error(), "is not")
	})

	t.Run("rejects a response without a to address", func(t *testing.T) {
		env := new(gobl.Envelope)
		require.NoError(t, json.Unmarshal(buildPartyEnvelope(t, subjKey, subject, nil, "", ""), env))
		env.Head.To = ""
		data, err := json.Marshal(env)
		require.NoError(t, err)

		c := NewClient(WithIdentity(testRequester, testRequesterKey), WithFetcher(&mapFetcher{data: map[string][]byte{
			subject.WhoURL():             data,
			subject.KeyURL(subjKey.ID()): jwkOf(subjKey),
		}}))
		_, err = c.Who(ctx, subject)
		require.ErrorIs(t, err, ErrVerifyFailed)
	})

	t.Run("requires a client identity", func(t *testing.T) {
		c := NewClient(WithFetcher(&mapFetcher{}))
		_, err := c.Who(ctx, subject)
		require.ErrorIs(t, err, ErrFetchFailed)
	})

	t.Run("rejects an invalid address", func(t *testing.T) {
		c := NewClient()
		_, err := c.Who(ctx, Address("localhost"))
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrAddressInvalid))
	})

	t.Run("propagates a no-content response", func(t *testing.T) {
		c := NewClient(WithIdentity(testRequester, testRequesterKey), WithFetcher(&mockFetcher{err: ErrNoContent}))
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

		c := NewClient(WithIdentity(testRequester, testRequesterKey), WithFetcher(&mapFetcher{data: map[string][]byte{
			subject.WhoURL(): data,
		}}))
		_, err = c.Who(ctx, subject)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrVerifyFailed))
	})

	t.Run("rejects a body that is not an envelope", func(t *testing.T) {
		c := NewClient(WithIdentity(testRequester, testRequesterKey), WithFetcher(&mapFetcher{data: map[string][]byte{
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
			WithIdentity(testRequester, testRequesterKey),
			WithAuthorities(authority),
			WithFetcher(&mapFetcher{data: endorsed}),
		)
		party, err := c.VerifySender(ctx, subject, true)
		require.NoError(t, err)
		assert.Equal(t, "Test Subject", party.Name)
	})

	t.Run("rejects a self-signed sender", func(t *testing.T) {
		c := NewClient(
			WithIdentity(testRequester, testRequesterKey),
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
			WithIdentity(testRequester, testRequesterKey),
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
			WithIdentity(testRequester, testRequesterKey),
			WithAuthorities(authority),
			WithFetcher(&mockFetcher{err: ErrNoContent}),
		)
		_, err := c.VerifySender(ctx, subject, true)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrNoContent))
	})
}
