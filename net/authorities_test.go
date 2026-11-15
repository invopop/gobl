package net

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/note"
	"github.com/invopop/gobl/org"
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
// endpoints), with optional per-URL errors.
type mapFetcher struct {
	data map[string][]byte
	errs map[string]error
}

func (m *mapFetcher) Fetch(_ context.Context, url string, _ http.Header) ([]byte, error) {
	if err, ok := m.errs[url]; ok {
		return nil, err
	}
	if d, ok := m.data[url]; ok {
		return d, nil
	}
	if strings.HasSuffix(url, "/keys/"+testSubjectKey.ID()) {
		return json.Marshal(testSubjectKey.Public())
	}
	if strings.HasSuffix(url, "/keys/"+testVerifierKey.ID()) {
		return json.Marshal(testVerifierKey.Public())
	}
	return nil, ErrFetchFailed
}

func (m *mapFetcher) Post(_ context.Context, _ string, _ []byte) error {
	return ErrFetchFailed
}

// testSubjectKey signs subject envelopes in authority tests; any
// address's key endpoint serves it under its kid.
var testSubjectKey = dsig.NewES256Key()

// regExp returns the exp an authority countersignature needs.
func regExp() head.SignOption {
	return head.WithExpiration(time.Now().Add(30 * 24 * time.Hour))
}

// regCnf returns the cnf an authority countersignature needs, binding
// testSubjectKey.
func regCnf() head.SignOption {
	return head.WithConfirmation(testSubjectKey.JKT())
}

// testVerifier countersigns subjects in authority tests; any address's
// key endpoint serves testVerifierKey under its kid.
var (
	testVerifier    = Address("verify.test.example")
	testVerifierKey = dsig.NewES256Key()
)

// regVerifier names testVerifier on an authority countersignature.
func regVerifier() head.SignOption {
	return head.WithVerifier(testVerifier.String())
}

// addVerification countersigns env as testVerifier for subject, binding
// testSubjectKey.
func addVerification(t *testing.T, env *gobl.Envelope, subject Address) {
	t.Helper()
	require.NoError(t, env.Sign(testVerifierKey, regCnf(),
		head.WithIssuer(testVerifier.String()),
		head.WithAudience(subject.String())))
}

// signPair signs env with testSubjectKey for its from and to.
func signPair(t *testing.T, env *gobl.Envelope) {
	t.Helper()
	from, err := ParseAddress(env.Head.From.Opaque())
	require.NoError(t, err)
	to, err := ParseAddress(env.Head.To.Opaque())
	require.NoError(t, err)
	require.NoError(t, Sign(env, testSubjectKey, from, to))
}

func TestVerifyEndorsement(t *testing.T) {
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

	t.Run("rejects an authority countersignature naming no verifier", func(t *testing.T) {
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		env.Head.From = subjectAddr.URI()
		env.Head.To = testRequester.URI()
		signPair(t, env)
		// Subject's self-signature (not an authority) + authority countersignature.
		require.NoError(t, env.Sign(subjKey, head.WithIssuer(subjectAddr.String())))
		require.NoError(t, env.Sign(authKey, regCnf(), regExp(),
			head.WithIssuer(authorityAddr.String()),
			head.WithAudience(subjectAddr.String())))

		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
			}}),
		)
		_, err = c.VerifyEndorsement(ctx, env)
		require.ErrorIs(t, err, ErrVerifyFailed)
		assert.Contains(t, err.Error(), "names no verifier")
	})

	t.Run("rejects an envelope its subject did not sign for the pair", func(t *testing.T) {
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		env.Head.From = subjectAddr.URI()
		env.Head.To = testRequester.URI()
		require.NoError(t, env.Sign(authKey, regCnf(), regExp(), regVerifier(),
			head.WithIssuer(authorityAddr.String()),
			head.WithAudience(subjectAddr.String())))
		addVerification(t, env, subjectAddr)

		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
			}}),
		)
		_, err = c.VerifyEndorsement(ctx, env)
		require.ErrorIs(t, err, ErrVerifyFailed)
		assert.ErrorContains(t, err, "no valid signature from")
	})

	t.Run("confirms a verifier named by the authority", func(t *testing.T) {
		verifierAddr := Address("verify.example.net")
		verifKey := dsig.NewES256Key()
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		env.Head.From = subjectAddr.URI()
		env.Head.To = testRequester.URI()
		signPair(t, env)
		require.NoError(t, env.Sign(subjKey, head.WithIssuer(subjectAddr.String())))
		require.NoError(t, env.Sign(authKey, regCnf(), regExp(),
			head.WithIssuer(authorityAddr.String()),
			head.WithAudience(subjectAddr.String()),
			head.WithVerifier(verifierAddr.String())))
		require.NoError(t, env.Sign(verifKey, regCnf(),
			head.WithIssuer(verifierAddr.String()),
			head.WithAudience(subjectAddr.String())))

		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
				verifierAddr.KeyURL(verifKey.ID()): jwkOf(verifKey),
			}}),
		)
		end, err := c.VerifyEndorsement(ctx, env)
		require.NoError(t, err)
		assert.Equal(t, authorityAddr, end.Authority)
		assert.Equal(t, verifierAddr, end.Verifier)
	})

	t.Run("authority cannot verify itself", func(t *testing.T) {
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		env.Head.From = subjectAddr.URI()
		env.Head.To = testRequester.URI()
		signPair(t, env)
		require.NoError(t, env.Sign(authKey, regCnf(), regExp(),
			head.WithAudience(subjectAddr.String()),
			head.WithIssuer(authorityAddr.String()),
			head.WithVerifier(authorityAddr.String())))

		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
			}}),
		)
		_, err = c.VerifyEndorsement(ctx, env)
		require.ErrorIs(t, err, ErrVerifyFailed)
		assert.Contains(t, err.Error(), "names itself")
	})

	t.Run("rejects when the verifier signature is missing", func(t *testing.T) {
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		env.Head.From = subjectAddr.URI()
		env.Head.To = testRequester.URI()
		signPair(t, env)
		require.NoError(t, env.Sign(authKey, regCnf(), regExp(),
			head.WithAudience(subjectAddr.String()),
			head.WithIssuer(authorityAddr.String()),
			head.WithVerifier("verify.example.net")))

		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
			}}),
		)
		_, err = c.VerifyEndorsement(ctx, env)
		require.ErrorIs(t, err, ErrVerifyFailed)
		assert.Contains(t, err.Error(), "no signature from")
	})

	t.Run("rejects when the verifier signature is expired", func(t *testing.T) {
		verifierAddr := Address("verify.example.net")
		verifKey := dsig.NewES256Key()
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		env.Head.From = subjectAddr.URI()
		env.Head.To = testRequester.URI()
		signPair(t, env)
		require.NoError(t, env.Sign(authKey, regCnf(), regExp(),
			head.WithAudience(subjectAddr.String()),
			head.WithIssuer(authorityAddr.String()),
			head.WithVerifier(verifierAddr.String())))
		require.NoError(t, env.Sign(verifKey, regCnf(), head.WithAudience(subjectAddr.String()),
			head.WithIssuer(verifierAddr.String()),
			head.WithExpiration(time.Now().Add(-time.Hour))))

		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
				verifierAddr.KeyURL(verifKey.ID()): jwkOf(verifKey),
			}}),
		)
		_, err = c.VerifyEndorsement(ctx, env)
		require.ErrorIs(t, err, ErrSignatureExpired)
	})

	t.Run("rejects when the verifier signature does not verify", func(t *testing.T) {
		// The verifier countersigned, but its published key is a
		// different one, so the crypto check fails.
		verifierAddr := Address("verify.example.net")
		verifKey := dsig.NewES256Key()
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		env.Head.From = subjectAddr.URI()
		env.Head.To = testRequester.URI()
		signPair(t, env)
		require.NoError(t, env.Sign(authKey, regCnf(), regExp(),
			head.WithAudience(subjectAddr.String()),
			head.WithIssuer(authorityAddr.String()),
			head.WithVerifier(verifierAddr.String())))
		require.NoError(t, env.Sign(verifKey, regCnf(), head.WithAudience(subjectAddr.String()), head.WithIssuer(verifierAddr.String())))

		other := dsig.NewES256Key()
		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
				verifierAddr.KeyURL(verifKey.ID()): jwkOf(other),
			}}),
		)
		_, err = c.VerifyEndorsement(ctx, env)
		require.ErrorIs(t, err, ErrVerifyFailed)
	})

	t.Run("rejects a malformed verifier claim", func(t *testing.T) {
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		env.Head.From = subjectAddr.URI()
		env.Head.To = testRequester.URI()
		signPair(t, env)
		require.NoError(t, env.Sign(authKey, regCnf(), regExp(),
			head.WithAudience(subjectAddr.String()),
			head.WithIssuer(authorityAddr.String()),
			head.WithVerifier("https://verify.example.net")))

		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
			}}),
		)
		_, err = c.VerifyEndorsement(ctx, env)
		require.ErrorIs(t, err, ErrVerifyFailed)
		assert.Contains(t, err.Error(), "invalid verifier")
	})

	t.Run("rejects an expired countersignature", func(t *testing.T) {
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		env.Head.From = subjectAddr.URI()
		env.Head.To = testRequester.URI()
		signPair(t, env)
		require.NoError(t, env.Sign(authKey, regCnf(),
			head.WithAudience(subjectAddr.String()),
			head.WithIssuer(authorityAddr.String()),
			head.WithExpiration(time.Now().Add(-time.Hour))))

		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
			}}),
		)
		_, err = c.VerifyEndorsement(ctx, env)
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
		env.Head.From = subjectAddr.URI()
		env.Head.To = testRequester.URI()
		signPair(t, env)
		require.NoError(t, env.Sign(authKey, regCnf(), regExp(), regVerifier(),
			head.WithAudience(subjectAddr.String()),
			head.WithIssuer("KYC.Example.COM.")))
		addVerification(t, env, subjectAddr)

		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
			}}),
		)
		_, err = c.VerifyEndorsement(ctx, env)
		assert.NoError(t, err)
	})

	t.Run("rejects a pre-epoch expiry", func(t *testing.T) {
		// A negative exp is a valid NumericDate before 1970 and must be
		// enforced, not treated as an absent claim.
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		env.Head.From = subjectAddr.URI()
		env.Head.To = testRequester.URI()
		signPair(t, env)
		require.NoError(t, env.Sign(authKey, regCnf(),
			head.WithAudience(subjectAddr.String()),
			head.WithIssuer(authorityAddr.String()),
			head.WithExpiration(time.Unix(-1000, 0))))

		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
			}}),
		)
		_, err = c.VerifyEndorsement(ctx, env)
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
		env.Head.From = subjectAddr.URI()
		env.Head.To = testRequester.URI()
		signPair(t, env)
		require.NoError(t, env.Sign(authKey, regCnf(),
			head.WithAudience(subjectAddr.String()),
			head.WithIssuer(authorityAddr.String()),
			head.WithExpiration(time.Now().Add(-time.Hour))))
		require.NoError(t, env.Sign(otherKey, regCnf(), regExp(),
			head.WithAudience(subjectAddr.String()), head.WithIssuer(otherAuth.String())))

		c := NewClient(
			WithAuthorities(authorityAddr, otherAuth),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				// Only the first authority's key resolves.
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
			}}),
		)
		_, err = c.VerifyEndorsement(ctx, env)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrSignatureExpired))
	})

	t.Run("accepts a countersignature within its exp window", func(t *testing.T) {
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		env.Head.From = subjectAddr.URI()
		env.Head.To = testRequester.URI()
		signPair(t, env)
		require.NoError(t, env.Sign(authKey, regCnf(), regVerifier(),
			head.WithAudience(subjectAddr.String()),
			head.WithIssuer(authorityAddr.String()),
			head.WithExpiration(time.Now().Add(60*24*time.Hour))))
		addVerification(t, env, subjectAddr)

		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
			}}),
		)
		end, err := c.VerifyEndorsement(ctx, env)
		require.NoError(t, err)
		assert.Equal(t, authorityAddr, end.Authority)
	})

	t.Run("rejects an envelope with only a self-signature", func(t *testing.T) {
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		env.Head.From = subjectAddr.URI()
		env.Head.To = testRequester.URI()
		signPair(t, env)
		require.NoError(t, env.Sign(subjKey, head.WithIssuer(subjectAddr.String())))

		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{}}),
		)
		_, err = c.VerifyEndorsement(ctx, env)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrNotEndorsed))
	})

	t.Run("skips undecodable signature payloads", func(t *testing.T) {
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		env.Head.From = subjectAddr.URI()
		env.Head.To = testRequester.URI()
		signPair(t, env)
		require.NoError(t, env.Sign(authKey, regCnf(), regExp(), regVerifier(),
			head.WithAudience(subjectAddr.String()),
			head.WithIssuer(authorityAddr.String())))
		addVerification(t, env, subjectAddr)
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
		_, err = c.VerifyEndorsement(ctx, env)
		assert.NoError(t, err)
	})

	t.Run("skips an invalid iss FQDN", func(t *testing.T) {
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		env.Head.From = subjectAddr.URI()
		env.Head.To = testRequester.URI()
		signPair(t, env)
		require.NoError(t, env.Sign(authKey, regCnf(), regExp(),
			head.WithAudience(subjectAddr.String()), head.WithIssuer("localhost")))

		c := NewClient(WithAuthorities(authorityAddr))
		_, err = c.VerifyEndorsement(ctx, env)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrNotEndorsed))
	})

	t.Run("rejects a tampered envelope", func(t *testing.T) {
		// The countersignature fetches its key fine but no longer
		// matches the (mutated) header digest.
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		env.Head.From = subjectAddr.URI()
		env.Head.To = testRequester.URI()
		signPair(t, env)
		require.NoError(t, env.Sign(authKey, regCnf(), regExp(),
			head.WithAudience(subjectAddr.String()), head.WithIssuer(authorityAddr.String())))
		env.Head.Digest = dsig.NewSHA256Digest([]byte("tampered"))

		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
			}}),
		)
		_, err = c.VerifyEndorsement(ctx, env)
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
		env.Head.From = subjectAddr.URI()
		env.Head.To = testRequester.URI()
		signPair(t, env)
		require.NoError(t, env.Sign(authKey, regCnf(), regExp(),
			head.WithAudience(subjectAddr.String()), head.WithIssuer(authorityAddr.String())))
		env.Head = nil

		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
			}}),
		)
		_, err = c.VerifyEndorsement(ctx, env)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrVerifyFailed))
	})

	t.Run("rejects an envelope with no signatures at all", func(t *testing.T) {
		msg := &note.Message{Content: "x"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		env.Head.From = subjectAddr.URI()
		env.Head.To = testRequester.URI()

		c := NewClient(WithAuthorities(authorityAddr))
		_, err = c.VerifyEndorsement(ctx, env)
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
		env.Head.From = subjectAddr.URI()
		env.Head.To = testRequester.URI()
		signPair(t, env)
		require.NoError(t, env.Sign(authKey, regCnf(), regExp(),
			head.WithAudience(subjectAddr.String()), head.WithIssuer(authorityAddr.String())))

		c := NewClient() // empty authorities
		_, err = c.VerifyEndorsement(ctx, env)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrNotEndorsed))
	})

	t.Run("trusts the default authority out of the box", func(t *testing.T) {
		lookupKey := dsig.NewES256Key()
		lookup := Address("lookup.gobl.org")
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		env.Head.From = subjectAddr.URI()
		env.Head.To = testRequester.URI()
		signPair(t, env)
		require.NoError(t, env.Sign(lookupKey, regCnf(), regExp(), regVerifier(),
			head.WithAudience(subjectAddr.String()), head.WithIssuer(lookup.String())))
		addVerification(t, env, subjectAddr)

		c := NewClient(WithFetcher(&mapFetcher{data: map[string][]byte{
			lookup.KeyURL(lookupKey.ID()): jwkOf(lookupKey),
		}}))
		end, err := c.VerifyEndorsement(ctx, env)
		require.NoError(t, err)
		assert.Equal(t, lookup, end.Authority)
	})

	t.Run("rejects when authority signature does not verify against published key", func(t *testing.T) {
		// The envelope claims to be signed by the authority, but the
		// fetcher returns a different key, so the crypto verification
		// fails. ErrVerifyFailed is returned (not ErrNotEndorsed).
		msg := &note.Message{Content: "x"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		env.Head.From = subjectAddr.URI()
		env.Head.To = testRequester.URI()
		signPair(t, env)
		require.NoError(t, env.Sign(authKey, regCnf(), regExp(),
			head.WithAudience(subjectAddr.String()), head.WithIssuer(authorityAddr.String())))

		other := dsig.NewES256Key()
		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				// Wrong key for the authority's claimed kid.
				authorityAddr.KeyURL(authKey.ID()): jwkOf(other),
			}}),
		)
		_, err = c.VerifyEndorsement(ctx, env)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrVerifyFailed))
	})

	t.Run("skips signatures with non-gobl iss schemes", func(t *testing.T) {
		// A signature whose iss isn't a gobl: URI cannot be an
		// authority signature — VerifyEndorsement just steps past it.
		msg := &note.Message{Content: "x"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		env.Head.From = subjectAddr.URI()
		env.Head.To = testRequester.URI()
		signPair(t, env)
		require.NoError(t, env.Sign(subjKey, head.WithIssuer("mailto:a@b")))

		c := NewClient(WithAuthorities(authorityAddr))
		_, err = c.VerifyEndorsement(ctx, env)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrNotEndorsed))
	})
}

func TestVerifySignatureByEdgeCases(t *testing.T) {
	ctx := context.Background()
	subjectAddr := Address("subject.example")
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
		// the endorsement is rejected. A non-object payload
		// signature is stepped past on the way.
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		env.Head.From = subjectAddr.URI()
		env.Head.To = testRequester.URI()
		signPair(t, env)
		require.NoError(t, env.Sign(authKey, regCnf(), regExp(),
			head.WithAudience(subjectAddr.String()),
			head.WithIssuer(authorityAddr.String()),
			head.WithVerifier(verifierAddr.String())))
		require.NoError(t, env.Sign(verifKey, regCnf(), head.WithAudience(subjectAddr.String()), head.WithIssuer(verifierAddr.String())))
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
		_, err = c.VerifyEndorsement(ctx, env)
		require.ErrorIs(t, err, ErrVerifyFailed)
	})

	t.Run("skips an authority countersignature naming no verifier", func(t *testing.T) {
		// Two trusted authorities countersign: the first names no
		// verifier, while the second names an independent one.
		secondAddr := Address("auth.example.org")
		secondKey := dsig.NewES256Key()
		verifierAddr := Address("verify.example.org")
		verifierKey := dsig.NewES256Key()
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		env.Head.From = subjectAddr.URI()
		env.Head.To = testRequester.URI()
		signPair(t, env)
		require.NoError(t, env.Sign(authKey, regCnf(), regExp(),
			head.WithAudience(subjectAddr.String()), head.WithIssuer(authorityAddr.String())))
		require.NoError(t, env.Sign(secondKey, regCnf(), regExp(),
			head.WithAudience(subjectAddr.String()),
			head.WithIssuer(secondAddr.String()),
			head.WithVerifier(verifierAddr.String())))
		require.NoError(t, env.Sign(verifierKey, regCnf(), head.WithAudience(subjectAddr.String()),
			head.WithIssuer(verifierAddr.String())))

		c := NewClient(
			WithAuthorities(authorityAddr, secondAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()):    jwkOf(authKey),
				secondAddr.KeyURL(secondKey.ID()):     jwkOf(secondKey),
				verifierAddr.KeyURL(verifierKey.ID()): jwkOf(verifierKey),
			}}),
		)
		end, err := c.VerifyEndorsement(ctx, env)
		require.NoError(t, err)
		assert.Equal(t, secondAddr, end.Authority)
		assert.Equal(t, verifierAddr, end.Verifier)
	})
}

func TestVerifyEndorsementSubjectBinding(t *testing.T) {
	ctx := context.Background()
	authorityAddr := Address("kyc.example.com")
	verifierAddr := Address("verify.example.net")
	subjectAddr := Address("subject.example")
	authKey := dsig.NewES256Key()
	verifKey := dsig.NewES256Key()

	jwkOf := func(k *dsig.PrivateKey) []byte {
		out, err := json.Marshal(k.Public())
		require.NoError(t, err)
		return out
	}
	client := func() *Client {
		return NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
				verifierAddr.KeyURL(verifKey.ID()): jwkOf(verifKey),
			}}),
		)
	}
	newEnv := func(t *testing.T) *gobl.Envelope {
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		env.Head.From = subjectAddr.URI()
		env.Head.To = testRequester.URI()
		signPair(t, env)
		return env
	}
	register := func(t *testing.T, env *gobl.Envelope, aud Address) {
		require.NoError(t, env.Sign(authKey, regCnf(), regExp(),
			head.WithIssuer(authorityAddr.String()),
			head.WithAudience(aud.String()),
			head.WithVerifier(verifierAddr.String())))
	}

	t.Run("ignores a registration bound to another address", func(t *testing.T) {
		env := newEnv(t)
		register(t, env, "someone-else.example")
		_, err := client().VerifyEndorsement(ctx, env)
		assert.ErrorIs(t, err, ErrNotEndorsed)
	})

	t.Run("rejects an envelope without a from address", func(t *testing.T) {
		env := newEnv(t)
		env.Head.From = ""
		register(t, env, subjectAddr)
		_, err := client().VerifyEndorsement(ctx, env)
		assert.ErrorIs(t, err, ErrVerifyFailed)
	})

	t.Run("ignores a registration without an audience", func(t *testing.T) {
		env := newEnv(t)
		require.NoError(t, env.Sign(authKey, regCnf(), regExp(), head.WithIssuer(authorityAddr.String())))
		_, err := client().VerifyEndorsement(ctx, env)
		assert.ErrorIs(t, err, ErrNotEndorsed)
	})

	t.Run("a verification serves only its own address", func(t *testing.T) {
		// A copied envelope registered at another address does not
		// keep the verification, so the registration fails.
		env := newEnv(t)
		require.NoError(t, env.Sign(verifKey, regCnf(),
			head.WithIssuer(verifierAddr.String()),
			head.WithAudience(subjectAddr.String())))
		register(t, env, subjectAddr)
		register(t, env, "copier.example")

		_, err := client().VerifyEndorsement(ctx, env)
		require.NoError(t, err)

		env.Head.From = "gobl:copier.example"
		signPair(t, env)
		_, err = client().VerifyEndorsement(ctx, env)
		require.ErrorIs(t, err, ErrVerifyFailed)
	})

	t.Run("each address needs its own verification", func(t *testing.T) {
		env := newEnv(t)
		for _, addr := range []Address{subjectAddr, "second.example"} {
			require.NoError(t, env.Sign(verifKey, regCnf(),
				head.WithIssuer(verifierAddr.String()),
				head.WithAudience(addr.String())))
			register(t, env, addr)
		}
		for _, addr := range []Address{subjectAddr, "second.example"} {
			env.Head.From = addr.URI()
			signPair(t, env)
			_, err := client().VerifyEndorsement(ctx, env)
			require.NoError(t, err, addr)
		}
	})

	t.Run("verifier signature without a cnf is rejected", func(t *testing.T) {
		env := newEnv(t)
		require.NoError(t, env.Sign(verifKey,
			head.WithIssuer(verifierAddr.String()),
			head.WithAudience(subjectAddr.String())))
		register(t, env, subjectAddr)
		_, err := client().VerifyEndorsement(ctx, env)
		require.ErrorIs(t, err, ErrVerifyFailed)
	})

	t.Run("verifier signature without an audience is rejected", func(t *testing.T) {
		env := newEnv(t)
		require.NoError(t, env.Sign(verifKey, regCnf(), head.WithIssuer(verifierAddr.String())))
		register(t, env, subjectAddr)
		_, err := client().VerifyEndorsement(ctx, env)
		require.ErrorIs(t, err, ErrVerifyFailed)
	})

	t.Run("registration without iat is rejected", func(t *testing.T) {
		env := newEnv(t)
		sig, err := dsig.NewSignature(authKey, &head.SigningPayload{
			UUID:      env.Head.UUID,
			Digest:    env.Head.Digest,
			Iss:       authorityAddr.String(),
			Aud:       subjectAddr.String(),
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		})
		require.NoError(t, err)
		env.Signatures = append(env.Signatures, sig)
		_, err = client().VerifyEndorsement(ctx, env)
		require.ErrorIs(t, err, ErrVerifyFailed)
		assert.Contains(t, err.Error(), "no iat")
	})

	t.Run("registration without exp is rejected", func(t *testing.T) {
		env := newEnv(t)
		require.NoError(t, env.Sign(authKey, regCnf(),
			head.WithIssuer(authorityAddr.String()),
			head.WithAudience(subjectAddr.String())))
		_, err := client().VerifyEndorsement(ctx, env)
		require.ErrorIs(t, err, ErrVerifyFailed)
		assert.Contains(t, err.Error(), "no exp")
	})

	t.Run("registration lasting more than 90 days is rejected", func(t *testing.T) {
		env := newEnv(t)
		require.NoError(t, env.Sign(authKey, regCnf(),
			head.WithIssuer(authorityAddr.String()),
			head.WithAudience(subjectAddr.String()),
			head.WithExpiration(time.Now().Add(91*24*time.Hour))))
		_, err := client().VerifyEndorsement(ctx, env)
		require.ErrorIs(t, err, ErrVerifyFailed)
		assert.Contains(t, err.Error(), "more than 90 days")
	})

	t.Run("registration without cnf is rejected", func(t *testing.T) {
		env := newEnv(t)
		require.NoError(t, env.Sign(authKey, regExp(),
			head.WithIssuer(authorityAddr.String()),
			head.WithAudience(subjectAddr.String())))
		_, err := client().VerifyEndorsement(ctx, env)
		require.ErrorIs(t, err, ErrVerifyFailed)
		assert.Contains(t, err.Error(), "no cnf")
	})

	t.Run("registration bound to another key is rejected", func(t *testing.T) {
		// A new holder of the address signs with its own key: the
		// registration it presents binds the previous holder's.
		env := newEnv(t)
		register(t, env, subjectAddr)
		env.Signatures = env.Signatures[1:] // drop testSubjectKey's pair signature
		newKey := dsig.NewES256Key()
		require.NoError(t, Sign(env, newKey, subjectAddr, testRequester))
		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
				subjectAddr.KeyURL(newKey.ID()):    jwkOf(newKey),
			}}),
		)
		_, err := c.VerifyEndorsement(ctx, env)
		require.ErrorIs(t, err, ErrVerifyFailed)
		assert.Contains(t, err.Error(), "another key")
	})

	t.Run("endorsement reports the bound key", func(t *testing.T) {
		env := newEnv(t)
		require.NoError(t, env.Sign(verifKey, regCnf(),
			head.WithIssuer(verifierAddr.String()),
			head.WithAudience(subjectAddr.String())))
		register(t, env, subjectAddr)
		end, err := client().VerifyEndorsement(ctx, env)
		require.NoError(t, err)
		assert.Equal(t, testSubjectKey.JKT(), end.Key)
	})

	t.Run("checks a bounded number of signatures per authority", func(t *testing.T) {
		// Signatures under unknown kids each cost a key fetch: only
		// the first few are tried.
		env := newEnv(t)
		for range 2 * maxIssuerCandidates {
			require.NoError(t, env.Sign(dsig.NewES256Key(), regExp(), regCnf(),
				head.WithIssuer(authorityAddr.String()),
				head.WithAudience(subjectAddr.String())))
		}
		f := &countingFetcher{inner: &mapFetcher{}, n: map[string]int{}}
		c := NewClient(WithAuthorities(authorityAddr), WithFetcher(f))
		_, err := c.VerifyEndorsement(ctx, env)
		require.Error(t, err)
		assert.Len(t, f.n, maxIssuerCandidates)
	})

	t.Run("verifier signature without iat is rejected", func(t *testing.T) {
		env := newEnv(t)
		sig, err := dsig.NewSignature(verifKey, &head.SigningPayload{
			UUID:   env.Head.UUID,
			Digest: env.Head.Digest,
			Iss:    verifierAddr.String(),
			Aud:    subjectAddr.String(),

			Confirmation: &head.Confirmation{JKT: testSubjectKey.JKT()},
		})
		require.NoError(t, err)
		env.Signatures = append(env.Signatures, sig)
		register(t, env, subjectAddr)
		_, err = client().VerifyEndorsement(ctx, env)
		require.ErrorIs(t, err, ErrVerifyFailed)
	})

	t.Run("verifier signature issued in the future is rejected", func(t *testing.T) {
		env := newEnv(t)
		sig, err := dsig.NewSignature(verifKey, &head.SigningPayload{
			UUID:     env.Head.UUID,
			Digest:   env.Head.Digest,
			Iss:      verifierAddr.String(),
			Aud:      subjectAddr.String(),
			IssuedAt: time.Now().Add(time.Hour).Unix(),

			Confirmation: &head.Confirmation{JKT: testSubjectKey.JKT()},
		})
		require.NoError(t, err)
		env.Signatures = append(env.Signatures, sig)
		register(t, env, subjectAddr)
		_, err = client().VerifyEndorsement(ctx, env)
		require.ErrorIs(t, err, ErrVerifyFailed)
	})
}

func TestWithAuthoritiesReplacesDefault(t *testing.T) {
	original := Authorities
	t.Cleanup(func() { Authorities = original })
	Authorities = []Address{"lookup.gobl.org"}

	// A closed deployment states its full trust list: the default
	// authority must not be implicitly trusted alongside it.
	c := NewClient(WithAuthorities("authority.corp.example"))
	assert.Equal(t, []Address{"authority.corp.example"}, c.authorities)
}

func TestVerifyEndorsementUnavailability(t *testing.T) {
	ctx := context.Background()
	subjectAddr := Address("subject.example")
	authorityAddr := Address("kyc.example.com")
	verifierAddr := Address("verify.example.net")
	authKey := dsig.NewES256Key()
	verifKey := dsig.NewES256Key()

	jwkOf := func(k *dsig.PrivateKey) []byte {
		out, err := json.Marshal(k.Public())
		require.NoError(t, err)
		return out
	}

	buildVerified := func(t *testing.T) *gobl.Envelope {
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		env.Head.From = subjectAddr.URI()
		env.Head.To = testRequester.URI()
		signPair(t, env)
		require.NoError(t, env.Sign(authKey, regCnf(), regExp(),
			head.WithAudience(subjectAddr.String()),
			head.WithIssuer(authorityAddr.String()),
			head.WithVerifier(verifierAddr.String())))
		require.NoError(t, env.Sign(verifKey, regCnf(), head.WithAudience(subjectAddr.String()), head.WithIssuer(verifierAddr.String())))
		return env
	}

	t.Run("unreachable verifier key surfaces instead of degrading", func(t *testing.T) {
		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{
				data: map[string][]byte{
					authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
				},
				errs: map[string]error{
					verifierAddr.KeyURL(verifKey.ID()): fmt.Errorf("%w: HTTP 503", ErrUnavailable),
				},
			}),
		)
		_, err := c.VerifyEndorsement(ctx, buildVerified(t))
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrUnavailable))
	})

	t.Run("unreachable authority key surfaces instead of rejecting", func(t *testing.T) {
		// The outcome is indeterminate while the authority's key
		// endpoint is down: callers must see the transient condition
		// (retry) rather than a definitive verification failure.
		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{
				errs: map[string]error{
					authorityAddr.KeyURL(authKey.ID()): fmt.Errorf("%w: HTTP 503", ErrUnavailable),
				},
			}),
		)
		_, err := c.VerifyEndorsement(ctx, buildVerified(t))
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrUnavailable))
		assert.False(t, errors.Is(err, ErrVerifyFailed))
	})

	t.Run("removed verifier key is rejected", func(t *testing.T) {
		// A definitive 404 means the key is gone — revoked.
		c := NewClient(
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkOf(authKey),
			}}),
		)
		_, err := c.VerifyEndorsement(ctx, buildVerified(t))
		require.ErrorIs(t, err, ErrVerifyFailed)
	})

	t.Run("rejects an envelope with too many signatures", func(t *testing.T) {
		env := buildVerified(t)
		for len(env.Signatures) <= maxEnvelopeSignatures {
			env.Signatures = append(env.Signatures, env.Signatures[0])
		}
		c := NewClient(WithAuthorities(authorityAddr), WithFetcher(&mapFetcher{}))
		_, err := c.VerifyEndorsement(ctx, env)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrVerifyFailed))
		assert.Contains(t, err.Error(), "signatures")
	})
}

func TestWithSandbox(t *testing.T) {
	original := Authorities
	t.Cleanup(func() { Authorities = original })
	Authorities = []Address{"lookup.gobl.org"}

	t.Run("switches to the sandbox trust list", func(t *testing.T) {
		c := NewClient(WithSandbox())
		assert.Equal(t, SandboxAuthorities, c.authorities)
	})

	t.Run("live default excludes sandbox authorities", func(t *testing.T) {
		c := NewClient()
		for _, sandbox := range SandboxAuthorities {
			assert.NotContains(t, c.authorities, sandbox,
				"live clients must never trust sandbox endorsements")
		}
	})

	t.Run("last trust option wins", func(t *testing.T) {
		c := NewClient(WithSandbox(), WithAuthorities("authority.corp.example"))
		assert.Equal(t, []Address{"authority.corp.example"}, c.authorities)
	})
}

func TestFindVerifications(t *testing.T) {
	ctx := context.Background()
	subjectAddr := Address("subject.example")
	verifierAddr := Address("verify.example.net")
	verifKey := dsig.NewES256Key()
	c := NewClient(WithFetcher(&mapFetcher{data: map[string][]byte{
		verifierAddr.KeyURL(verifKey.ID()): jwkFromKey(t, verifKey),
	}}))
	newEnv := func(t *testing.T) *gobl.Envelope {
		env, err := gobl.Envelop(&org.Party{Name: "Subject"})
		require.NoError(t, err)
		env.Head.From = "gobl:lookup.example"
		env.Head.To = subjectAddr.URI()
		return env
	}

	t.Run("returns the verifier's cnf and times", func(t *testing.T) {
		env := newEnv(t)
		exp := time.Now().Add(365 * 24 * time.Hour).Truncate(time.Second)
		require.NoError(t, env.Sign(verifKey,
			head.WithIssuer(verifierAddr.String()),
			head.WithAudience(subjectAddr.String()),
			head.WithConfirmation(testSubjectKey.JKT()),
			head.WithExpiration(exp)))
		vs, err := c.FindVerifications(ctx, env, verifierAddr, subjectAddr)
		require.NoError(t, err)
		require.Len(t, vs, 1)
		v := vs[0]
		assert.Equal(t, verifierAddr, v.Verifier)
		assert.Equal(t, testSubjectKey.JKT(), v.Key)
		assert.WithinDuration(t, time.Now(), v.IssuedAt, time.Minute)
		assert.True(t, exp.Equal(v.ExpiresAt))
	})

	t.Run("ignores a verification for another address", func(t *testing.T) {
		env := newEnv(t)
		require.NoError(t, env.Sign(verifKey, regCnf(),
			head.WithIssuer(verifierAddr.String()),
			head.WithAudience("other.example")))
		_, err := c.FindVerifications(ctx, env, verifierAddr, subjectAddr)
		require.ErrorIs(t, err, ErrVerifyFailed)
	})

	t.Run("rejects a verification without a cnf", func(t *testing.T) {
		env := newEnv(t)
		require.NoError(t, env.Sign(verifKey,
			head.WithIssuer(verifierAddr.String()),
			head.WithAudience(subjectAddr.String())))
		_, err := c.FindVerifications(ctx, env, verifierAddr, subjectAddr)
		require.ErrorIs(t, err, ErrVerifyFailed)
		assert.Contains(t, err.Error(), "no cnf")
	})

	t.Run("rejects an expired verification", func(t *testing.T) {
		env := newEnv(t)
		require.NoError(t, env.Sign(verifKey, regCnf(),
			head.WithIssuer(verifierAddr.String()),
			head.WithAudience(subjectAddr.String()),
			head.WithExpiration(time.Now().Add(-time.Hour))))
		_, err := c.FindVerifications(ctx, env, verifierAddr, subjectAddr)
		require.ErrorIs(t, err, ErrSignatureExpired)
	})

	t.Run("rejects an unsigned envelope", func(t *testing.T) {
		_, err := c.FindVerifications(ctx, newEnv(t), verifierAddr, subjectAddr)
		require.ErrorIs(t, err, ErrVerifyFailed)
		_, err = c.FindVerifications(ctx, nil, verifierAddr, subjectAddr)
		require.ErrorIs(t, err, ErrVerifyFailed)
	})

	t.Run("rejects too many signatures", func(t *testing.T) {
		env := newEnv(t)
		for range maxEnvelopeSignatures + 1 {
			require.NoError(t, env.Sign(verifKey))
		}
		_, err := c.FindVerifications(ctx, env, verifierAddr, subjectAddr)
		require.ErrorIs(t, err, ErrVerifyFailed)
		assert.ErrorContains(t, err, "signatures (max")
	})

	t.Run("rejects invalid addresses", func(t *testing.T) {
		env := newEnv(t)
		require.NoError(t, env.Sign(verifKey, regCnf(),
			head.WithIssuer(verifierAddr.String()),
			head.WithAudience(subjectAddr.String())))
		_, err := c.FindVerifications(ctx, env, "not an address", subjectAddr)
		require.ErrorIs(t, err, ErrAddressInvalid)
		_, err = c.FindVerifications(ctx, env, verifierAddr, "not an address")
		require.ErrorIs(t, err, ErrAddressInvalid)
	})

	t.Run("reports a verifier key that cannot be fetched", func(t *testing.T) {
		env := newEnv(t)
		require.NoError(t, env.Sign(dsig.NewES256Key(), regCnf(),
			head.WithIssuer(verifierAddr.String()),
			head.WithAudience(subjectAddr.String())))
		_, err := c.FindVerifications(ctx, env, verifierAddr, subjectAddr)
		require.ErrorIs(t, err, ErrFetchFailed)
	})

	t.Run("checks a bounded number of verifications", func(t *testing.T) {
		env := newEnv(t)
		for range maxIssuerCandidates + 1 {
			require.NoError(t, env.Sign(verifKey, regCnf(),
				head.WithIssuer(verifierAddr.String()),
				head.WithAudience(subjectAddr.String())))
		}
		vs, err := c.FindVerifications(ctx, env, verifierAddr, subjectAddr)
		require.NoError(t, err)
		assert.Len(t, vs, maxIssuerCandidates)
	})

	t.Run("returns every valid verification", func(t *testing.T) {
		env := newEnv(t)
		other := dsig.NewES256Key()
		for _, k := range []*dsig.PrivateKey{testSubjectKey, other} {
			require.NoError(t, env.Sign(verifKey, head.WithConfirmation(k.JKT()),
				head.WithIssuer(verifierAddr.String()),
				head.WithAudience(subjectAddr.String())))
		}
		vs, err := c.FindVerifications(ctx, env, verifierAddr, subjectAddr)
		require.NoError(t, err)
		require.Len(t, vs, 2)
		assert.Equal(t, testSubjectKey.JKT(), vs[0].Key)
		assert.Equal(t, other.JKT(), vs[1].Key)
	})
}

func TestWithVerifiers(t *testing.T) {
	assert.Empty(t, NewClient().verifiers, "no verifiers are trusted by default")
	c := NewClient(WithVerifiers("verify.example.net"))
	assert.Equal(t, []Address{"verify.example.net"}, c.verifiers)
}

func TestVerifyEndorsementDirect(t *testing.T) {
	ctx := context.Background()
	subjectAddr := Address("subject.example")
	authorityAddr := Address("kyc.example.com")
	authKey := dsig.NewES256Key()
	newEnv := func(t *testing.T) *gobl.Envelope {
		msg := &note.Message{Content: "party doc"}
		msg.SetUUID(uuid.V7())
		env, err := gobl.Envelop(msg)
		require.NoError(t, err)
		env.Head.From = subjectAddr.URI()
		env.Head.To = testRequester.URI()
		signPair(t, env)
		return env
	}
	client := func(opts ...ClientOption) *Client {
		return NewClient(append([]ClientOption{
			WithAuthorities(authorityAddr),
			WithFetcher(&mapFetcher{data: map[string][]byte{
				authorityAddr.KeyURL(authKey.ID()): jwkFromKey(t, authKey),
			}}),
		}, opts...)...)
	}

	t.Run("accepts a trusted verifier without an authority", func(t *testing.T) {
		env := newEnv(t)
		addVerification(t, env, subjectAddr)
		end, err := client(WithVerifiers(testVerifier)).VerifyEndorsement(ctx, env)
		require.NoError(t, err)
		assert.Empty(t, end.Authority)
		assert.Equal(t, testVerifier, end.Verifier)
		assert.Equal(t, testSubjectKey.JKT(), end.Key)
	})

	t.Run("accepts a trusted verifier with no authorities configured", func(t *testing.T) {
		env := newEnv(t)
		addVerification(t, env, subjectAddr)
		end, err := client(WithAuthorities(), WithVerifiers(testVerifier)).VerifyEndorsement(ctx, env)
		require.NoError(t, err)
		assert.Equal(t, testVerifier, end.Verifier)
	})

	t.Run("ignores a verifier the client does not trust", func(t *testing.T) {
		env := newEnv(t)
		addVerification(t, env, subjectAddr)
		_, err := client().VerifyEndorsement(ctx, env)
		require.ErrorIs(t, err, ErrNotEndorsed)
	})

	t.Run("prefers the authority", func(t *testing.T) {
		env := newEnv(t)
		addVerification(t, env, subjectAddr)
		require.NoError(t, env.Sign(authKey, regCnf(), regExp(), regVerifier(),
			head.WithIssuer(authorityAddr.String()),
			head.WithAudience(subjectAddr.String())))
		end, err := client(WithVerifiers(testVerifier)).VerifyEndorsement(ctx, env)
		require.NoError(t, err)
		assert.Equal(t, authorityAddr, end.Authority)
		assert.Equal(t, testVerifier, end.Verifier)
	})

	t.Run("rejects a verification bound to another key", func(t *testing.T) {
		// With no authority keeping a lineage, the verifier's cnf must
		// name the key confirming the envelope.
		env := newEnv(t)
		require.NoError(t, env.Sign(testVerifierKey,
			head.WithConfirmation(dsig.NewES256Key().JKT()),
			head.WithIssuer(testVerifier.String()),
			head.WithAudience(subjectAddr.String())))
		_, err := client(WithVerifiers(testVerifier)).VerifyEndorsement(ctx, env)
		require.ErrorIs(t, err, ErrVerifyFailed)
		assert.Contains(t, err.Error(), "another key")
	})

	t.Run("finds the verification for the current key", func(t *testing.T) {
		// After a key change the envelope may carry the verifier's
		// earlier countersignature too.
		env := newEnv(t)
		require.NoError(t, env.Sign(testVerifierKey,
			head.WithConfirmation(dsig.NewES256Key().JKT()),
			head.WithIssuer(testVerifier.String()),
			head.WithAudience(subjectAddr.String())))
		addVerification(t, env, subjectAddr)
		end, err := client(WithVerifiers(testVerifier)).VerifyEndorsement(ctx, env)
		require.NoError(t, err)
		assert.Equal(t, testSubjectKey.JKT(), end.Key)
	})

	t.Run("rejects an expired verification", func(t *testing.T) {
		env := newEnv(t)
		require.NoError(t, env.Sign(testVerifierKey, regCnf(),
			head.WithExpiration(time.Now().Add(-time.Hour)),
			head.WithIssuer(testVerifier.String()),
			head.WithAudience(subjectAddr.String())))
		_, err := client(WithVerifiers(testVerifier)).VerifyEndorsement(ctx, env)
		require.ErrorIs(t, err, ErrSignatureExpired)
	})
}
