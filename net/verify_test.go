package net

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/note"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// partyEnvelopeFor envelopes a party declaring addr as its gobl:
// endpoint, signed with key and the given options.
func partyEnvelopeFor(t *testing.T, key *dsig.PrivateKey, addr Address, opts ...head.SignOption) *gobl.Envelope {
	t.Helper()
	party := &org.Party{
		Name:      "Test Party",
		Endpoints: []*org.Endpoint{{URI: addr.URI()}},
	}
	party.SetUUID(uuid.V7())
	env, err := gobl.Envelop(party)
	require.NoError(t, err)
	require.NoError(t, env.Sign(key, opts...))

	// Round-trip through JSON to simulate realistic usage.
	data, err := json.Marshal(env)
	require.NoError(t, err)
	env = new(gobl.Envelope)
	require.NoError(t, json.Unmarshal(data, env))
	return env
}

// buildTestEnvelope envelopes a plain document (not a party), signed
// with key and the given iss/aud — the document-delivery shape.
func buildTestEnvelope(t *testing.T, key *dsig.PrivateKey, iss, aud string) *gobl.Envelope {
	t.Helper()

	msg := &note.Message{Content: "test message content"}
	msg.SetUUID(uuid.V7())

	env, err := gobl.Envelop(msg)
	require.NoError(t, err)
	opts := []head.SignOption{head.WithIssuer(iss)}
	if aud != "" {
		opts = append(opts, head.WithAudience(aud))
	}
	require.NoError(t, env.Sign(key, opts...))

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

func TestVerifyParty(t *testing.T) {
	ctx := context.Background()
	addr := Address("billing.invopop.com")

	t.Run("bearer envelope", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := partyEnvelopeFor(t, key, addr, head.WithIssuer(addr.String()))

		mock := &mockFetcher{data: jwkFromKey(t, key)}
		c := NewClient(WithFetcher(mock))
		subject, err := c.VerifyParty(ctx, env)
		require.NoError(t, err)
		assert.Equal(t, addr, subject)
		assert.Equal(t, addr.KeyURL(key.ID()), mock.url)
	})

	t.Run("signature order is not significant", func(t *testing.T) {
		// A countersignature serialized before the self-signature
		// changes nothing: the subject comes from the document.
		key := dsig.NewES256Key()
		env := partyEnvelopeFor(t, key, addr, head.WithIssuer(addr.String()))
		other := dsig.NewES256Key()
		require.NoError(t, env.Sign(other,
			head.WithIssuer("kyc.example.com"),
			head.WithAudience(addr.String())))
		env.Signatures[0], env.Signatures[1] = env.Signatures[1], env.Signatures[0]

		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
		subject, err := c.VerifyParty(ctx, env)
		require.NoError(t, err)
		assert.Equal(t, addr, subject)
	})

	t.Run("legacy audience-bound self-signature", func(t *testing.T) {
		// Envelopes signed under the previous revision carry an
		// audience on the self-signature; VerifyParty accepts any
		// valid self-signature.
		key := dsig.NewES256Key()
		env := partyEnvelopeFor(t, key, addr,
			head.WithIssuer(addr.String()),
			head.WithAudience("lookup.example.com"))

		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
		subject, err := c.VerifyParty(ctx, env)
		require.NoError(t, err)
		assert.Equal(t, addr, subject)
	})

	t.Run("non-canonical endpoint is canonicalized", func(t *testing.T) {
		key := dsig.NewES256Key()
		party := &org.Party{
			Name:      "Canon",
			Endpoints: []*org.Endpoint{{URI: "gobl:Billing.Invopop.COM."}},
		}
		party.SetUUID(uuid.V7())
		env, err := gobl.Envelop(party)
		require.NoError(t, err)
		require.NoError(t, env.Sign(key, head.WithIssuer(addr.String())))

		mock := &mockFetcher{data: jwkFromKey(t, key)}
		c := NewClient(WithFetcher(mock))
		subject, err := c.VerifyParty(ctx, env)
		require.NoError(t, err)
		assert.Equal(t, addr, subject)
		assert.Equal(t, addr.KeyURL(key.ID()), mock.url, "key URL uses the canonical form")
	})

	t.Run("not signed", func(t *testing.T) {
		c := NewClient()
		_, err := c.VerifyParty(ctx, new(gobl.Envelope))
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrVerifyFailed))
	})

	t.Run("not a party", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := buildTestEnvelope(t, key, addr.String(), "")
		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
		_, err := c.VerifyParty(ctx, env)
		require.ErrorIs(t, err, ErrPartyMissing)
	})

	t.Run("party endpoint is not a valid address", func(t *testing.T) {
		key := dsig.NewES256Key()
		party := &org.Party{
			Name:      "Bad Endpoint",
			Endpoints: []*org.Endpoint{{URI: "gobl:localhost"}},
		}
		party.SetUUID(uuid.V7())
		env, err := gobl.Envelop(party)
		require.NoError(t, err)
		require.NoError(t, env.Sign(key, head.WithIssuer(addr.String())))

		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
		_, err = c.VerifyParty(ctx, env)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrVerifyFailed))
		assert.Contains(t, err.Error(), "not a valid address")
	})

	t.Run("party without a gobl endpoint", func(t *testing.T) {
		key := dsig.NewES256Key()
		party := &org.Party{Name: "No Endpoint"}
		party.SetUUID(uuid.V7())
		env, err := gobl.Envelop(party)
		require.NoError(t, err)
		require.NoError(t, env.Sign(key, head.WithIssuer(addr.String())))

		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
		_, err = c.VerifyParty(ctx, env)
		require.ErrorIs(t, err, ErrPartyMissing)
	})

	t.Run("declared endpoint differs from the signer", func(t *testing.T) {
		// Alice signs an envelope whose party document claims Bob's
		// address: the subject is Bob, and there is no valid
		// self-signature by Bob.
		key := dsig.NewES256Key()
		env := partyEnvelopeFor(t, key, "bob.example.com",
			head.WithIssuer("alice.example.com"))

		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
		_, err := c.VerifyParty(ctx, env)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrVerifyFailed))
		assert.Contains(t, err.Error(), "bob.example.com")
	})

	t.Run("unknown signer key", func(t *testing.T) {
		// The per-key endpoint returns a JWK whose kid does not match
		// the signature's kid, so the self-signature never verifies.
		key := dsig.NewES256Key()
		other := dsig.NewES256Key()
		env := partyEnvelopeFor(t, key, addr, head.WithIssuer(addr.String()))
		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, other)}))
		_, err := c.VerifyParty(ctx, env)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrVerifyFailed))
	})

	t.Run("subject key unavailable", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := partyEnvelopeFor(t, key, addr, head.WithIssuer(addr.String()))
		c := NewClient(WithFetcher(&mapFetcher{errs: map[string]error{
			addr.KeyURL(key.ID()): fmt.Errorf("%w: HTTP 503", ErrUnavailable),
		}}))
		_, err := c.VerifyParty(ctx, env)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrUnavailable))
		assert.False(t, errors.Is(err, ErrVerifyFailed))
	})

	t.Run("nil signature entries are skipped", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := partyEnvelopeFor(t, key, addr, head.WithIssuer(addr.String()))
		env.Signatures = append([]*dsig.Signature{nil}, env.Signatures...)

		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
		subject, err := c.VerifyParty(ctx, env)
		require.NoError(t, err)
		assert.Equal(t, addr, subject)
	})

	t.Run("signing time before valid_from", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := partyEnvelopeFor(t, key, addr, head.WithIssuer(addr.String()))
		future := cal.TimestampOf(time.Now().Add(24 * time.Hour))
		c := NewClient(WithFetcher(&mockFetcher{data: publishedJWK(t, key, &future, nil)}))
		_, err := c.VerifyParty(ctx, env)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrVerifyFailed))
	})

	t.Run("signing time after valid_until", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := partyEnvelopeFor(t, key, addr, head.WithIssuer(addr.String()))
		past := cal.TimestampOf(time.Now().Add(-24 * time.Hour))
		c := NewClient(WithFetcher(&mockFetcher{data: publishedJWK(t, key, nil, &past)}))
		_, err := c.VerifyParty(ctx, env)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrVerifyFailed))
	})

	t.Run("signing time inside window", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := partyEnvelopeFor(t, key, addr, head.WithIssuer(addr.String()))
		from := cal.TimestampOf(time.Now().Add(-time.Hour))
		until := cal.TimestampOf(time.Now().Add(time.Hour))
		c := NewClient(WithFetcher(&mockFetcher{data: publishedJWK(t, key, &from, &until)}))
		subject, err := c.VerifyParty(ctx, env)
		require.NoError(t, err)
		assert.Equal(t, addr, subject)
	})

	t.Run("rejects an envelope with too many signatures", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := partyEnvelopeFor(t, key, addr, head.WithIssuer(addr.String()))
		for len(env.Signatures) <= maxEnvelopeSignatures {
			env.Signatures = append(env.Signatures, env.Signatures[0])
		}
		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
		_, err := c.VerifyParty(ctx, env)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrVerifyFailed))
		assert.Contains(t, err.Error(), "signatures")
	})
}

func TestVerifyDelivery(t *testing.T) {
	ctx := context.Background()
	sender := Address("supplier.example.com")
	inbox := Address("customer.example.com")

	t.Run("bound delivery", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := buildTestEnvelope(t, key, sender.String(), inbox.String())

		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
		got, err := c.VerifyDelivery(ctx, env, inbox)
		require.NoError(t, err)
		assert.Equal(t, sender, got)
	})

	t.Run("binding may sit on any signature", func(t *testing.T) {
		// The sender signs once per recipient; this inbox's binding is
		// the second signature.
		key := dsig.NewES256Key()
		env := buildTestEnvelope(t, key, sender.String(), "other.example.com")
		require.NoError(t, env.Sign(key,
			head.WithIssuer(sender.String()),
			head.WithAudience(inbox.String())))

		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
		got, err := c.VerifyDelivery(ctx, env, inbox)
		require.NoError(t, err)
		assert.Equal(t, sender, got)
	})

	t.Run("no binding", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := buildTestEnvelope(t, key, sender.String(), "other.example.com")

		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
		_, err := c.VerifyDelivery(ctx, env, inbox)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrVerifyFailed))
	})

	t.Run("unsigned audience does not bind", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := buildTestEnvelope(t, key, sender.String(), "")

		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
		_, err := c.VerifyDelivery(ctx, env, inbox)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrVerifyFailed))
	})

	t.Run("two issuers binding is ambiguous", func(t *testing.T) {
		keyA, keyB := dsig.NewES256Key(), dsig.NewES256Key()
		other := Address("relay.example.com")
		env := buildTestEnvelope(t, keyA, sender.String(), inbox.String())
		require.NoError(t, env.Sign(keyB,
			head.WithIssuer(other.String()),
			head.WithAudience(inbox.String())))

		c := NewClient(WithFetcher(&mapFetcher{data: map[string][]byte{
			sender.KeyURL(keyA.ID()): jwkFromKey(t, keyA),
			other.KeyURL(keyB.ID()):  jwkFromKey(t, keyB),
		}}))
		_, err := c.VerifyDelivery(ctx, env, inbox)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrVerifyFailed))
		assert.Contains(t, err.Error(), "ambiguous")
	})

	t.Run("same issuer binding twice is fine", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := buildTestEnvelope(t, key, sender.String(), inbox.String())
		require.NoError(t, env.Sign(key,
			head.WithIssuer(sender.String()),
			head.WithAudience(inbox.String())))

		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
		got, err := c.VerifyDelivery(ctx, env, inbox)
		require.NoError(t, err)
		assert.Equal(t, sender, got)
	})

	t.Run("binding signer key unavailable", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := buildTestEnvelope(t, key, sender.String(), inbox.String())
		c := NewClient(WithFetcher(&mapFetcher{errs: map[string]error{
			sender.KeyURL(key.ID()): fmt.Errorf("%w: HTTP 503", ErrUnavailable),
		}}))
		_, err := c.VerifyDelivery(ctx, env, inbox)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrUnavailable))
	})

	t.Run("forged binding does not verify", func(t *testing.T) {
		// The binding signature's published key is different material
		// under the same kid: the crypto check fails.
		key := dsig.NewES256Key()
		env := buildTestEnvelope(t, key, sender.String(), inbox.String())
		other := dsig.NewES256Key()
		var jwk map[string]any
		require.NoError(t, json.Unmarshal(jwkFromKey(t, other), &jwk))
		jwk["kid"] = key.ID()
		forged, err := json.Marshal(jwk)
		require.NoError(t, err)

		c := NewClient(WithFetcher(&mapFetcher{data: map[string][]byte{
			sender.KeyURL(key.ID()): forged,
		}}))
		_, err = c.VerifyDelivery(ctx, env, inbox)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrVerifyFailed))
	})

	t.Run("binding with an invalid issuer is skipped", func(t *testing.T) {
		// A signature bound to this inbox whose iss is not a valid
		// address can never resolve a key and does not bind.
		key := dsig.NewES256Key()
		env := buildTestEnvelope(t, key, "localhost", inbox.String())

		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
		_, err := c.VerifyDelivery(ctx, env, inbox)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrVerifyFailed))
	})

	t.Run("invalid self address", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := buildTestEnvelope(t, key, sender.String(), inbox.String())
		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
		_, err := c.VerifyDelivery(ctx, env, "not valid!")
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrVerifyFailed))
	})

	t.Run("not signed", func(t *testing.T) {
		c := NewClient()
		_, err := c.VerifyDelivery(ctx, new(gobl.Envelope), inbox)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrVerifyFailed))
	})

	t.Run("rejects an envelope with too many signatures", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := buildTestEnvelope(t, key, sender.String(), inbox.String())
		for len(env.Signatures) <= maxEnvelopeSignatures {
			env.Signatures = append(env.Signatures, env.Signatures[0])
		}
		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
		_, err := c.VerifyDelivery(ctx, env, inbox)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrVerifyFailed))
		assert.Contains(t, err.Error(), "signatures")
	})

	t.Run("nil signature entries are skipped", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := buildTestEnvelope(t, key, sender.String(), inbox.String())
		env.Signatures = append([]*dsig.Signature{nil}, env.Signatures...)

		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
		got, err := c.VerifyDelivery(ctx, env, inbox)
		require.NoError(t, err)
		assert.Equal(t, sender, got)
	})
}

func TestWhoPublicationSignatureUnavailable(t *testing.T) {
	// The publication (audience-free) self-signature uses a second key
	// whose endpoint is down while another self-signature verifies:
	// Who reports the transient condition instead of rejecting.
	ctx := context.Background()
	subject := Address("issuer.example.com")
	keyA, keyB := dsig.NewES256Key(), dsig.NewES256Key()

	env := partyEnvelopeFor(t, keyA, subject,
		head.WithIssuer(subject.String()),
		head.WithAudience("lookup.example.com"))
	require.NoError(t, env.Sign(keyB, head.WithIssuer(subject.String())))
	data, err := json.Marshal(env)
	require.NoError(t, err)

	c := NewClient(WithFetcher(&mapFetcher{
		data: map[string][]byte{
			subject.WhoURL():          data,
			subject.KeyURL(keyA.ID()): jwkFromKey(t, keyA),
		},
		errs: map[string]error{
			subject.KeyURL(keyB.ID()): fmt.Errorf("%w: HTTP 503", ErrUnavailable),
		},
	}))
	_, err = c.Who(ctx, subject)
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrUnavailable))
}
