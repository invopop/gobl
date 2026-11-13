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
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/note"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// partyEnvelopeFor envelopes a party from addr to testRequester,
// declaring addr as its gobl: endpoint, signed with key and the given
// options.
func partyEnvelopeFor(t *testing.T, key *dsig.PrivateKey, addr Address, opts ...head.SignOption) *gobl.Envelope {
	t.Helper()
	party := &org.Party{
		Name:      "Test Party",
		Endpoints: []*org.Endpoint{{URI: addr.URI()}},
	}
	party.SetUUID(uuid.V7())
	env, err := gobl.Envelop(party)
	require.NoError(t, err)
	env.Head.From = addr.URI()
	env.Head.To = testRequester.URI()
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

// buildDelivery envelopes a plain document from one address to another,
// signed with key by from for to.
func buildDelivery(t *testing.T, key *dsig.PrivateKey, from, to Address) *gobl.Envelope {
	t.Helper()
	env := buildTestEnvelope(t, key, from.String(), to.String())
	env.Head.From, env.Head.To = from.URI(), to.URI()
	return env
}

func TestVerifyDelivery(t *testing.T) {
	ctx := context.Background()
	sender := Address("supplier.example.com")
	inbox := Address("customer.example.com")

	t.Run("bound delivery", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := buildDelivery(t, key, sender, inbox)

		mock := &mockFetcher{data: jwkFromKey(t, key)}
		c := NewClient(WithFetcher(mock))
		require.NoError(t, c.VerifyDelivery(ctx, env, inbox))
		assert.Equal(t, sender.KeyURL(key.ID()), mock.url)
	})

	t.Run("party envelope", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := partyEnvelopeFor(t, key, sender,
			head.WithIssuer(sender.String()),
			head.WithAudience(testRequester.String()))

		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
		require.NoError(t, c.VerifyDelivery(ctx, env, testRequester))
	})

	t.Run("non-canonical routing is canonicalized", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := buildDelivery(t, key, sender, inbox)
		env.Head.From = "gobl:Supplier.Example.COM."
		env.Head.To = "gobl:CUSTOMER.example.com"

		mock := &mockFetcher{data: jwkFromKey(t, key)}
		c := NewClient(WithFetcher(mock))
		require.NoError(t, c.VerifyDelivery(ctx, env, "Customer.Example.com."))
		assert.Equal(t, sender.KeyURL(key.ID()), mock.url, "key URL uses the canonical form")
	})

	t.Run("rejects a missing from", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := buildDelivery(t, key, sender, inbox)
		env.Head.From = ""

		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
		err := c.VerifyDelivery(ctx, env, inbox)
		require.ErrorIs(t, err, ErrVerifyFailed)
		assert.Contains(t, err.Error(), "no from address")
	})

	t.Run("rejects a missing to", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := buildDelivery(t, key, sender, inbox)
		env.Head.To = ""

		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
		err := c.VerifyDelivery(ctx, env, inbox)
		require.ErrorIs(t, err, ErrVerifyFailed)
		assert.Contains(t, err.Error(), "no to address")
	})

	t.Run("rejects a from on another transport", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := buildDelivery(t, key, sender, inbox)
		env.Head.From = "iso6523-actorid-upis::9920:b123123123"

		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
		err := c.VerifyDelivery(ctx, env, inbox)
		require.ErrorIs(t, err, ErrVerifyFailed)
		assert.Contains(t, err.Error(), "not a gobl: address")
	})

	t.Run("rejects a to on another transport", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := buildDelivery(t, key, sender, inbox)
		env.Head.To = "mailto:billing@example.com"

		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
		err := c.VerifyDelivery(ctx, env, inbox)
		require.ErrorIs(t, err, ErrVerifyFailed)
		assert.Contains(t, err.Error(), "not a gobl: address")
	})

	t.Run("scheme is case-insensitive", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := buildDelivery(t, key, sender, inbox)
		env.Head.From = "GOBL:supplier.example.com"
		env.Head.To = "Gobl:customer.example.com"

		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
		require.NoError(t, c.VerifyDelivery(ctx, env, inbox))
	})

	t.Run("rejects a query or fragment", func(t *testing.T) {
		for _, uri := range []cbc.URI{
			"gobl:supplier.example.com?x=1",
			"gobl:supplier.example.com?",
			"gobl:supplier.example.com#f",
		} {
			key := dsig.NewES256Key()
			env := buildDelivery(t, key, sender, inbox)
			env.Head.From = uri

			c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
			err := c.VerifyDelivery(ctx, env, inbox)
			require.ErrorIs(t, err, ErrVerifyFailed, uri)
			assert.Contains(t, err.Error(), "query or fragment")
		}
	})

	t.Run("rejects a malformed gobl from", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := buildDelivery(t, key, sender, inbox)
		env.Head.From = "gobl:localhost"

		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
		require.ErrorIs(t, c.VerifyDelivery(ctx, env, inbox), ErrVerifyFailed)
	})

	t.Run("rejects a to naming another address", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := buildDelivery(t, key, sender, inbox)
		env.Head.To = "gobl:other.example.com"

		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
		err := c.VerifyDelivery(ctx, env, inbox)
		require.ErrorIs(t, err, ErrVerifyFailed)
		assert.Contains(t, err.Error(), "is not")
	})

	t.Run("rejects a from no signature carries", func(t *testing.T) {
		// Signed by the sender, but the header claims somebody else.
		key := dsig.NewES256Key()
		env := buildDelivery(t, key, sender, inbox)
		env.Head.From = "gobl:someone-else.example.com"

		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
		err := c.VerifyDelivery(ctx, env, inbox)
		require.ErrorIs(t, err, ErrVerifyFailed)
		assert.Contains(t, err.Error(), "someone-else.example.com")
	})

	t.Run("rejects a signature by from for another audience", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := buildTestEnvelope(t, key, sender.String(), "other.example.com")
		env.Head.From, env.Head.To = sender.URI(), inbox.URI()

		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
		require.ErrorIs(t, c.VerifyDelivery(ctx, env, inbox), ErrVerifyFailed)
	})

	t.Run("rejects an audience-free signature by from", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := buildTestEnvelope(t, key, sender.String(), "")
		env.Head.From, env.Head.To = sender.URI(), inbox.URI()

		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
		require.ErrorIs(t, c.VerifyDelivery(ctx, env, inbox), ErrVerifyFailed)
	})

	t.Run("rejects a signature for to by another issuer", func(t *testing.T) {
		// Forwarding is not supported: a relay signing for the inbox
		// does not make the original sender's envelope deliverable.
		relay := Address("relay.example.com")
		key := dsig.NewES256Key()
		env := buildDelivery(t, key, relay, inbox)
		env.Head.From = sender.URI()

		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
		require.ErrorIs(t, c.VerifyDelivery(ctx, env, inbox), ErrVerifyFailed)
	})

	t.Run("binding may sit on any signature", func(t *testing.T) {
		// Signature order is not significant: this pair's binding is
		// the second signature.
		key := dsig.NewES256Key()
		env := buildTestEnvelope(t, key, sender.String(), "other.example.com")
		env.Head.From, env.Head.To = sender.URI(), inbox.URI()
		require.NoError(t, env.Sign(key,
			head.WithIssuer(sender.String()),
			head.WithAudience(inbox.String())))

		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
		require.NoError(t, c.VerifyDelivery(ctx, env, inbox))
	})

	t.Run("one envelope, several pairs", func(t *testing.T) {
		// One party envelope published at two addresses: each pair is
		// confirmed by its own signature.
		key := dsig.NewES256Key()
		env := partyEnvelopeFor(t, key, sender,
			head.WithIssuer(sender.String()),
			head.WithAudience(testRequester.String()))
		other := Address("invoices.example.org")
		otherKey := dsig.NewES256Key()
		require.NoError(t, env.Sign(otherKey,
			head.WithIssuer(other.String()),
			head.WithAudience(testRequester.String())))

		c := NewClient(WithFetcher(&mapFetcher{data: map[string][]byte{
			sender.KeyURL(key.ID()):     jwkFromKey(t, key),
			other.KeyURL(otherKey.ID()): jwkFromKey(t, otherKey),
		}}))
		for _, a := range []Address{sender, other} {
			env.Head.From = a.URI()
			require.NoError(t, c.VerifyDelivery(ctx, env, testRequester))
		}
	})

	t.Run("unknown signer key", func(t *testing.T) {
		// The per-key endpoint returns a JWK whose kid does not match
		// the signature's kid, so the signature never verifies.
		key := dsig.NewES256Key()
		env := buildDelivery(t, key, sender, inbox)
		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, dsig.NewES256Key())}))
		require.ErrorIs(t, c.VerifyDelivery(ctx, env, inbox), ErrVerifyFailed)
	})

	t.Run("signer key unavailable", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := buildDelivery(t, key, sender, inbox)
		c := NewClient(WithFetcher(&mapFetcher{errs: map[string]error{
			sender.KeyURL(key.ID()): fmt.Errorf("%w: HTTP 503", ErrUnavailable),
		}}))
		err := c.VerifyDelivery(ctx, env, inbox)
		require.ErrorIs(t, err, ErrUnavailable)
		assert.False(t, errors.Is(err, ErrVerifyFailed))
	})

	t.Run("forged binding does not verify", func(t *testing.T) {
		// The published key is different material under the same kid:
		// the crypto check fails.
		key := dsig.NewES256Key()
		env := buildDelivery(t, key, sender, inbox)
		var jwk map[string]any
		require.NoError(t, json.Unmarshal(jwkFromKey(t, dsig.NewES256Key()), &jwk))
		jwk["kid"] = key.ID()
		forged, err := json.Marshal(jwk)
		require.NoError(t, err)

		c := NewClient(WithFetcher(&mapFetcher{data: map[string][]byte{
			sender.KeyURL(key.ID()): forged,
		}}))
		require.ErrorIs(t, c.VerifyDelivery(ctx, env, inbox), ErrVerifyFailed)
	})

	t.Run("signing time before valid_from", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := buildDelivery(t, key, sender, inbox)
		future := cal.TimestampOf(time.Now().Add(24 * time.Hour))
		c := NewClient(WithFetcher(&mockFetcher{data: publishedJWK(t, key, &future, nil)}))
		require.ErrorIs(t, c.VerifyDelivery(ctx, env, inbox), ErrVerifyFailed)
	})

	t.Run("signing time after valid_until", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := buildDelivery(t, key, sender, inbox)
		past := cal.TimestampOf(time.Now().Add(-24 * time.Hour))
		c := NewClient(WithFetcher(&mockFetcher{data: publishedJWK(t, key, nil, &past)}))
		require.ErrorIs(t, c.VerifyDelivery(ctx, env, inbox), ErrVerifyFailed)
	})

	t.Run("signing time inside window", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := buildDelivery(t, key, sender, inbox)
		from := cal.TimestampOf(time.Now().Add(-time.Hour))
		until := cal.TimestampOf(time.Now().Add(time.Hour))
		c := NewClient(WithFetcher(&mockFetcher{data: publishedJWK(t, key, &from, &until)}))
		require.NoError(t, c.VerifyDelivery(ctx, env, inbox))
	})

	// pairSignature appends a signature by sender for inbox with the
	// given iat, bypassing Sign's own timestamp.
	pairSignature := func(t *testing.T, env *gobl.Envelope, key *dsig.PrivateKey, iat int64) {
		t.Helper()
		sig, err := dsig.NewSignature(key, &head.SigningPayload{
			UUID:     env.Head.UUID,
			Digest:   env.Head.Digest,
			Iss:      sender.String(),
			Aud:      inbox.String(),
			IssuedAt: iat,
		})
		require.NoError(t, err)
		env.Signatures = []*dsig.Signature{sig}
	}

	t.Run("rejects a signature without iat", func(t *testing.T) {
		// Without an iat a retired key's validity window would never
		// be checked.
		key := dsig.NewES256Key()
		env := buildDelivery(t, key, sender, inbox)
		pairSignature(t, env, key, 0)
		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
		require.ErrorIs(t, c.VerifyDelivery(ctx, env, inbox), ErrVerifyFailed)
	})

	t.Run("rejects a signature issued in the future", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := buildDelivery(t, key, sender, inbox)
		pairSignature(t, env, key, time.Now().Add(time.Hour).Unix())
		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
		require.ErrorIs(t, c.VerifyDelivery(ctx, env, inbox), ErrVerifyFailed)
	})

	t.Run("rejects a key retired since signing", func(t *testing.T) {
		// Signed while the key was valid, but it has since been
		// retired: a live receiver does not accept it.
		key := dsig.NewES256Key()
		env := buildDelivery(t, key, sender, inbox)
		pairSignature(t, env, key, time.Now().Add(-2*time.Hour).Unix())
		until := cal.TimestampOf(time.Now().Add(-time.Hour))
		c := NewClient(WithFetcher(&mockFetcher{data: publishedJWK(t, key, nil, &until)}))
		require.ErrorIs(t, c.VerifyDelivery(ctx, env, inbox), ErrVerifyFailed)
	})

	t.Run("checks a bounded number of signatures", func(t *testing.T) {
		// Signatures under unknown kids each cost a key fetch from the
		// sender: only the first few are tried.
		env := buildDelivery(t, dsig.NewES256Key(), sender, inbox)
		for range 2 * maxIssuerCandidates {
			require.NoError(t, Sign(env, dsig.NewES256Key(), sender, inbox))
		}
		f := &countingFetcher{inner: &mapFetcher{}, n: map[string]int{}}
		c := NewClient(WithFetcher(f))
		require.ErrorIs(t, c.VerifyDelivery(ctx, env, inbox), ErrVerifyFailed)
		assert.Len(t, f.n, maxIssuerCandidates)
	})

	t.Run("invalid self address", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := buildDelivery(t, key, sender, inbox)
		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
		require.ErrorIs(t, c.VerifyDelivery(ctx, env, "not valid!"), ErrVerifyFailed)
	})

	t.Run("not signed", func(t *testing.T) {
		c := NewClient()
		require.ErrorIs(t, c.VerifyDelivery(ctx, new(gobl.Envelope), inbox), ErrVerifyFailed)
	})

	t.Run("rejects an envelope with too many signatures", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := buildDelivery(t, key, sender, inbox)
		for len(env.Signatures) <= maxEnvelopeSignatures {
			env.Signatures = append(env.Signatures, env.Signatures[0])
		}
		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
		err := c.VerifyDelivery(ctx, env, inbox)
		require.ErrorIs(t, err, ErrVerifyFailed)
		assert.Contains(t, err.Error(), "signatures")
	})

	t.Run("nil signature entries are skipped", func(t *testing.T) {
		key := dsig.NewES256Key()
		env := buildDelivery(t, key, sender, inbox)
		env.Signatures = append([]*dsig.Signature{nil}, env.Signatures...)

		c := NewClient(WithFetcher(&mockFetcher{data: jwkFromKey(t, key)}))
		require.NoError(t, c.VerifyDelivery(ctx, env, inbox))
	})
}

func TestWhoSignatureUnavailable(t *testing.T) {
	// The self-signature for the requester uses a second key whose
	// endpoint is down while another self-signature verifies: Who
	// reports the transient condition instead of rejecting.
	ctx := context.Background()
	subject := Address("issuer.example.com")
	keyA, keyB := dsig.NewES256Key(), dsig.NewES256Key()

	env := partyEnvelopeFor(t, keyA, subject,
		head.WithIssuer(subject.String()),
		head.WithAudience("lookup.example.com"))
	require.NoError(t, env.Sign(keyB,
		head.WithIssuer(subject.String()),
		head.WithAudience(testRequester.String())))
	data, err := json.Marshal(env)
	require.NoError(t, err)

	c := NewClient(WithIdentity(testRequester, testRequesterKey), WithFetcher(&mapFetcher{
		data: map[string][]byte{
			subject.WhoURL():          data,
			subject.KeyURL(keyA.ID()): jwkFromKey(t, keyA),
		},
		errs: map[string]error{
			subject.KeyURL(keyB.ID()): fmt.Errorf("%w: HTTP 503", ErrUnavailable),
		},
	}))
	_, err = c.Who(ctx, subject)
	require.ErrorIs(t, err, ErrUnavailable)
}
