package net

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/invopop/gobl"
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

func (m *mapFetcher) Fetch(_ context.Context, url string, _ http.Header) ([]byte, error) {
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
		require.NoError(t, env.Sign(subjKey, head.WithIssuer(subjectAddr.String())))
		require.NoError(t, env.Sign(authKey,
			head.WithIssuer(authorityAddr.String()),
			head.WithAudience(subjectAddr.String())))

		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
			}}),
		)
		end, err := c.VerifyAuthority(ctx, env)
		require.NoError(t, err)
		assert.Equal(t, authorityAddr, end.Authority)
		assert.False(t, end.Verified(), "no verifier claim means registered only")
	})

	t.Run("confirms a verifier named by the authority", func(t *testing.T) {
		verifierAddr := Address("verify.example.net")
		verifKey := dsig.NewES256Key()
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		require.NoError(t, env.Sign(subjKey, head.WithIssuer(subjectAddr.String())))
		require.NoError(t, env.Sign(authKey,
			head.WithIssuer(authorityAddr.String()),
			head.WithAudience(subjectAddr.String()),
			head.WithVerifier(verifierAddr.String())))
		require.NoError(t, env.Sign(verifKey,
			head.WithIssuer(verifierAddr.String()),
			head.WithAudience(subjectAddr.String())))

		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
				verifierAddr.KeyURL(verifKey.ID()): jwkOf(verifKey),
			}}),
		)
		end, err := c.VerifyAuthority(ctx, env)
		require.NoError(t, err)
		assert.Equal(t, authorityAddr, end.Authority)
		assert.Equal(t, verifierAddr, end.Verifier)
		assert.True(t, end.Verified())
	})

	t.Run("authority naming itself needs no second signature", func(t *testing.T) {
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		require.NoError(t, env.Sign(authKey,
			head.WithIssuer(authorityAddr.String()),
			head.WithVerifier(authorityAddr.String())))

		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
			}}),
		)
		end, err := c.VerifyAuthority(ctx, env)
		require.NoError(t, err)
		assert.Equal(t, authorityAddr, end.Verifier)
		assert.True(t, end.Verified())
	})

	t.Run("degrades to registered when the verifier signature is missing", func(t *testing.T) {
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		require.NoError(t, env.Sign(authKey,
			head.WithIssuer(authorityAddr.String()),
			head.WithVerifier("verify.example.net")))

		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
			}}),
		)
		end, err := c.VerifyAuthority(ctx, env)
		require.NoError(t, err)
		assert.Equal(t, authorityAddr, end.Authority)
		assert.False(t, end.Verified())
	})

	t.Run("degrades when the verifier signature is expired", func(t *testing.T) {
		verifierAddr := Address("verify.example.net")
		verifKey := dsig.NewES256Key()
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		require.NoError(t, env.Sign(authKey,
			head.WithIssuer(authorityAddr.String()),
			head.WithVerifier(verifierAddr.String())))
		require.NoError(t, env.Sign(verifKey,
			head.WithIssuer(verifierAddr.String()),
			head.WithExpiration(time.Now().Add(-time.Hour))))

		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
				verifierAddr.KeyURL(verifKey.ID()): jwkOf(verifKey),
			}}),
		)
		end, err := c.VerifyAuthority(ctx, env)
		require.NoError(t, err)
		assert.False(t, end.Verified())
	})

	t.Run("degrades when the verifier signature does not verify", func(t *testing.T) {
		// The verifier countersigned, but its published key is a
		// different one, so the crypto check fails.
		verifierAddr := Address("verify.example.net")
		verifKey := dsig.NewES256Key()
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		require.NoError(t, env.Sign(authKey,
			head.WithIssuer(authorityAddr.String()),
			head.WithVerifier(verifierAddr.String())))
		require.NoError(t, env.Sign(verifKey, head.WithIssuer(verifierAddr.String())))

		other := dsig.NewES256Key()
		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
				verifierAddr.KeyURL(verifKey.ID()): jwkOf(other),
			}}),
		)
		end, err := c.VerifyAuthority(ctx, env)
		require.NoError(t, err)
		assert.False(t, end.Verified())
	})

	t.Run("ignores a malformed verifier claim", func(t *testing.T) {
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		require.NoError(t, env.Sign(authKey,
			head.WithIssuer(authorityAddr.String()),
			head.WithVerifier("https://verify.example.net")))

		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
			}}),
		)
		end, err := c.VerifyAuthority(ctx, env)
		require.NoError(t, err)
		assert.False(t, end.Verified())
	})

	t.Run("rejects an expired countersignature", func(t *testing.T) {
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		require.NoError(t, env.Sign(authKey,
			head.WithIssuer(authorityAddr.String()),
			head.WithExpiration(time.Now().Add(-time.Hour))))

		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
			}}),
		)
		_, err = c.VerifyAuthority(ctx, env)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrSignatureExpired))
	})

	t.Run("canonicalizes a non-canonical authority iss", func(t *testing.T) {
		// The countersignature names the authority in mixed case with a
		// trailing dot; it must still match the registered authority.
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		require.NoError(t, env.Sign(authKey,
			head.WithIssuer("KYC.Example.COM.")))

		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
			}}),
		)
		_, err = c.VerifyAuthority(ctx, env)
		assert.NoError(t, err)
	})

	t.Run("rejects a pre-epoch expiry", func(t *testing.T) {
		// A negative exp is a valid NumericDate before 1970 and must be
		// enforced, not treated as an absent claim.
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		require.NoError(t, env.Sign(authKey,
			head.WithIssuer(authorityAddr.String()),
			head.WithExpiration(time.Unix(-1000, 0))))

		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
			}}),
		)
		_, err = c.VerifyAuthority(ctx, env)
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
			head.WithIssuer(authorityAddr.String()),
			head.WithExpiration(time.Now().Add(-time.Hour))))
		require.NoError(t, env.Sign(otherKey, head.WithIssuer(otherAuth.String())))

		c := NewClient(
			WithAuthorities(authorityAddr, otherAuth),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				// Only the first authority's key resolves.
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
			}}),
		)
		_, err = c.VerifyAuthority(ctx, env)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrSignatureExpired))
	})

	t.Run("accepts a countersignature within its exp window", func(t *testing.T) {
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		require.NoError(t, env.Sign(authKey,
			head.WithIssuer(authorityAddr.String()),
			head.WithExpiration(time.Now().Add(90*24*time.Hour))))

		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
			}}),
		)
		end, err := c.VerifyAuthority(ctx, env)
		require.NoError(t, err)
		assert.Equal(t, authorityAddr, end.Authority)
	})

	t.Run("rejects an envelope with only a self-signature", func(t *testing.T) {
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		require.NoError(t, env.Sign(subjKey, head.WithIssuer(subjectAddr.String())))

		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{}}),
		)
		_, err = c.VerifyAuthority(ctx, env)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrUnknownAuthority))
	})

	t.Run("skips undecodable signature payloads", func(t *testing.T) {
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		require.NoError(t, env.Sign(authKey,
			head.WithIssuer(authorityAddr.String())))
		// Prepend a signature whose payload does not decode; it must be
		// stepped past, not fail the whole check.
		bad, err := dsig.NewSignature(authKey, "not-an-object")
		require.NoError(t, err)
		env.Signatures = append([]*dsig.Signature{bad}, env.Signatures...)

		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
			}}),
		)
		_, err = c.VerifyAuthority(ctx, env)
		assert.NoError(t, err)
	})

	t.Run("skips an invalid iss FQDN", func(t *testing.T) {
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		require.NoError(t, env.Sign(authKey, head.WithIssuer("localhost")))

		c := NewClient(WithAuthorities(authorityAddr))
		_, err = c.VerifyAuthority(ctx, env)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrUnknownAuthority))
	})

	t.Run("rejects a tampered envelope", func(t *testing.T) {
		// The countersignature fetches its key fine but no longer
		// matches the (mutated) header digest.
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		require.NoError(t, env.Sign(authKey, head.WithIssuer(authorityAddr.String())))
		env.Head.Digest = dsig.NewSHA256Digest([]byte("tampered"))

		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
			}}),
		)
		_, err = c.VerifyAuthority(ctx, env)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrVerifyFailed))
	})

	t.Run("rejects signatures without a header", func(t *testing.T) {
		// A malformed envelope carrying signatures but no head must be
		// rejected, not panic.
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		require.NoError(t, env.Sign(authKey, head.WithIssuer(authorityAddr.String())))
		env.Head = nil

		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
			}}),
		)
		_, err = c.VerifyAuthority(ctx, env)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrVerifyFailed))
	})

	t.Run("rejects an envelope with no signatures at all", func(t *testing.T) {
		msg := &note.Message{Content: "x"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)

		c := NewClient(WithAuthorities(authorityAddr))
		_, err = c.VerifyAuthority(ctx, env)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrVerifyFailed))
	})

	t.Run("rejects when client has no authorities registered", func(t *testing.T) {
		original := Authorities
		t.Cleanup(func() { Authorities = original })
		Authorities = nil

		msg := &note.Message{Content: "x"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		require.NoError(t, env.Sign(authKey, head.WithIssuer(authorityAddr.String())))

		c := NewClient() // empty authorities
		_, err = c.VerifyAuthority(ctx, env)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrUnknownAuthority))
	})

	t.Run("trusts the default authority out of the box", func(t *testing.T) {
		lookupKey := dsig.NewES256Key()
		lookup := Address("lookup.gobl.org")
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		require.NoError(t, env.Sign(lookupKey, head.WithIssuer(lookup.String())))

		c := NewClient(WithFetcher(&mapFetcher{data: map[string][]byte{
			lookup.KeyURL(lookupKey.ID()): jwkOf(lookupKey),
		}}))
		end, err := c.VerifyAuthority(ctx, env)
		require.NoError(t, err)
		assert.Equal(t, lookup, end.Authority)
	})

	t.Run("rejects when authority signature does not verify against published key", func(t *testing.T) {
		// The envelope claims to be signed by the authority, but the
		// fetcher returns a different key, so the crypto verification
		// fails. ErrVerifyFailed is returned (not ErrUnknownAuthority).
		msg := &note.Message{Content: "x"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		require.NoError(t, env.Sign(authKey, head.WithIssuer(authorityAddr.String())))

		other := dsig.NewES256Key()
		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				// Wrong key for the authority's claimed kid.
				authorityAddr.KeyURL(authKey.ID()): jwkOf(other),
			}}),
		)
		_, err = c.VerifyAuthority(ctx, env)
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
		_, err = c.VerifyAuthority(ctx, env)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrUnknownAuthority))
	})
}

func TestVerifySignatureByEdgeCases(t *testing.T) {
	ctx := context.Background()
	authorityAddr := Address("kyc.example.com")
	verifierAddr := Address("verify.example.net")
	authKey := dsig.NewES256Key()
	verifKey := dsig.NewES256Key()

	jwkOf := func(k *dsig.PrivateKey) []byte {
		out, err := json.Marshal(k.Public())
		require.NoError(t, err)
		return out
	}

	t.Run("steps past undecodable payloads and rejects a kid-spoofed key", func(t *testing.T) {
		// The verifier's published JWK claims its kid but holds other
		// key material, so the fetch succeeds and the crypto fails —
		// the endorsement degrades to registered. A non-object payload
		// signature is stepped past on the way.
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		require.NoError(t, env.Sign(authKey,
			head.WithIssuer(authorityAddr.String()),
			head.WithVerifier(verifierAddr.String())))
		require.NoError(t, env.Sign(verifKey, head.WithIssuer(verifierAddr.String())))
		bad, err := dsig.NewSignature(verifKey, "not-an-object")
		require.NoError(t, err)
		env.Signatures = append(env.Signatures, bad)

		other := dsig.NewES256Key()
		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
				verifierAddr.KeyURL(verifKey.ID()): spoofKID(t, other, verifKey.ID()),
			}}),
		)
		end, err := c.VerifyAuthority(ctx, env)
		require.NoError(t, err)
		assert.False(t, end.Verified())
	})

	t.Run("prefers a verified endorsement regardless of signature order", func(t *testing.T) {
		// Two trusted authorities countersign: the first is registered
		// only, the second names itself as verifier. The verified
		// endorsement must win even though it appears later.
		secondAddr := Address("auth.example.org")
		secondKey := dsig.NewES256Key()
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		require.NoError(t, env.Sign(authKey, head.WithIssuer(authorityAddr.String())))
		require.NoError(t, env.Sign(secondKey,
			head.WithIssuer(secondAddr.String()),
			head.WithVerifier(secondAddr.String())))

		c := NewClient(
			WithAuthorities(authorityAddr, secondAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
				secondAddr.KeyURL(secondKey.ID()):  jwkOf(secondKey),
			}}),
		)
		end, err := c.VerifyAuthority(ctx, env)
		require.NoError(t, err)
		assert.True(t, end.Verified())
		assert.Equal(t, secondAddr, end.Authority)
		assert.Equal(t, secondAddr, end.Verifier)
	})
}
