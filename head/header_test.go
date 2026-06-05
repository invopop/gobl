package head_test

import (
	"testing"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHeader(t *testing.T) {
	h := head.NewHeader()
	assert.False(t, h.UUID.IsZero())
	assert.Equal(t, uuid.Version(7), h.UUID.Version())
	assert.NotPanics(t, func() {
		h.Meta["foo"] = "bar"
		h.Tags = append(h.Tags, "foo")
		h.Stamps = append(h.Stamps, &head.Stamp{})
	}, "header and meta hash should have been initialized")
}

func TestHeaderValidation(t *testing.T) {
	t.Run("missing digest", func(t *testing.T) {
		h := head.NewHeader()
		err := rules.Validate(h)
		assert.ErrorContains(t, err, "header must have a digest")
	})

	t.Run("valid header", func(t *testing.T) {
		h := head.NewHeader()
		h.Digest = dsig.NewSHA256Digest([]byte("testing"))
		assert.NoError(t, rules.Validate(h))
	})

	t.Run("duplicate stamps", func(t *testing.T) {
		h := head.NewHeader()
		h.Digest = dsig.NewSHA256Digest([]byte("testing"))
		h.Stamps = []*head.Stamp{
			{Provider: "foo", Value: "bar"},
			{Provider: "foo", Value: "bar"},
		}
		err := rules.Validate(h)
		assert.ErrorContains(t, err, "duplicate stamp providers are not allowed")
	})

	t.Run("duplicate links", func(t *testing.T) {
		h := head.NewHeader()
		h.Digest = dsig.NewSHA256Digest([]byte("testing"))
		h.Links = []*head.Link{
			{Key: "foo", URL: "https://example.com"},
			{Key: "foo", URL: "https://example.com/2"},
		}
		err := rules.Validate(h)
		assert.ErrorContains(t, err, "duplicate link keys are not allowed")
	})

	t.Run("valid from/to", func(t *testing.T) {
		h := head.NewHeader()
		h.Digest = dsig.NewSHA256Digest([]byte("testing"))
		h.From = cbc.URI("gobl:samlown.example.com")
		h.To = cbc.URI("iso6523-actorid-upis::9920:x3157928m")
		assert.NoError(t, rules.Validate(h))
	})

	t.Run("empty from/to allowed", func(t *testing.T) {
		h := head.NewHeader()
		h.Digest = dsig.NewSHA256Digest([]byte("testing"))
		assert.NoError(t, rules.Validate(h))
	})

	t.Run("invalid to", func(t *testing.T) {
		h := head.NewHeader()
		h.Digest = dsig.NewSHA256Digest([]byte("testing"))
		h.To = cbc.URI("no scheme")
		err := rules.Validate(h)
		assert.ErrorContains(t, err, "valid absolute URI")
	})
}

func TestHeaderRulesContext(t *testing.T) {
	// The accumulator on rules.Context is unexported, so behaviour is
	// covered end-to-end in the envelope tests; here we only guard the
	// nil / empty cases against panics.
	t.Run("nil header is a no-op", func(t *testing.T) {
		var h *head.Header
		assert.NotPanics(t, func() { h.RulesContext()(new(rules.Context)) })
	})

	t.Run("empty ignore list is a no-op", func(t *testing.T) {
		h := &head.Header{}
		assert.NotPanics(t, func() { h.RulesContext()(new(rules.Context)) })
	})

	t.Run("populated ignore list yields a usable option", func(t *testing.T) {
		h := &head.Header{Ignore: []rules.Code{"GOBL-NOTE-MESSAGE-01"}}
		assert.NotPanics(t, func() { h.RulesContext()(new(rules.Context)) })
	})
}

func TestHeaderAddStamp(t *testing.T) {
	h := head.NewHeader()
	h.AddStamp(&head.Stamp{Provider: "foo", Value: "bar"})
	assert.Len(t, h.Stamps, 1)
	h.AddStamp(&head.Stamp{Provider: "foo", Value: "baz"})
	assert.Len(t, h.Stamps, 1)
	assert.Equal(t, "baz", h.Stamps[0].Value)
	h.AddStamp(&head.Stamp{Provider: "bar", Value: "bax"})
	assert.Len(t, h.Stamps, 2)
	assert.Equal(t, "bax", h.Stamps[1].Value)
}

func TestHeaderAddLink(t *testing.T) {
	h := head.NewHeader()
	assert.Len(t, h.Links, 0)
	h.AddLink(&head.Link{Key: "foo", URL: "bar.com"})
	assert.Len(t, h.Links, 1)
}

func TestHeaderLink(t *testing.T) {
	t.Run("without category", func(t *testing.T) {
		h := head.NewHeader()
		h.AddLink(&head.Link{Key: "foo", URL: "bar.com"})
		l := h.Link("", "foo")
		assert.NotNil(t, l)
		assert.Equal(t, "bar.com", l.URL)
		l = h.Link("", "baa")
		assert.Nil(t, l)
	})
	t.Run("with category", func(t *testing.T) {
		h := head.NewHeader()
		h.AddLink(&head.Link{Category: head.LinkCategoryKeyPortal, Key: "foo", URL: "bar.com"})
		l := h.Link(head.LinkCategoryKeyPortal, "foo")
		assert.NotNil(t, l)
		assert.Equal(t, "bar.com", l.URL)
		l = h.Link(head.LinkCategoryKeyPortal, "baa")
		assert.Nil(t, l)
	})
}

func TestHeaderStamp(t *testing.T) {
	h := head.NewHeader()
	h.AddStamp(&head.Stamp{Provider: "foo", Value: "boo"})
	h.AddStamp(&head.Stamp{Provider: "foo2", Value: "bling"})
	st := h.GetStamp("foo")
	assert.NotNil(t, st)
	assert.Equal(t, "boo", st.Value)
	st = h.GetStamp("bad")
	assert.Nil(t, st)
}

func TestHeaderContains(t *testing.T) {
	h1 := head.NewHeader()
	h2 := head.NewHeader()

	// UUID
	assert.False(t, h1.Contains(h2))
	h2.UUID = h1.UUID
	assert.True(t, h1.Contains(h2))

	// Digest
	h1.Digest = dsig.NewSHA256Digest([]byte("testing"))
	h2.Digest = dsig.NewSHA256Digest([]byte("testing 2"))
	assert.False(t, h1.Contains(h2))
	h2.Digest = h1.Digest
	assert.True(t, h1.Contains(h2))

	// Stamps
	h1.AddStamp(&head.Stamp{Provider: "foo", Value: "boo"})
	h1.AddStamp(&head.Stamp{Provider: "foo2", Value: "bling"})
	assert.True(t, h1.Contains(h2))
	h2.AddStamp(&head.Stamp{Provider: "foo", Value: "boo"})
	assert.True(t, h1.Contains(h2))
	h2.AddStamp(&head.Stamp{Provider: "foo2", Value: "bling"})
	assert.True(t, h1.Contains(h2))
	h2.AddStamp(&head.Stamp{Provider: "foo3", Value: "bow"})
	assert.False(t, h1.Contains(h2))
	h1.AddStamp(&head.Stamp{Provider: "foo3", Value: "bow"})
	assert.True(t, h1.Contains(h2))

	// Links
	h1.AddLink(&head.Link{Key: "foo", URL: "bar.com"})
	h1.AddLink(&head.Link{Key: "foo2", URL: "bar2.com"})
	assert.True(t, h1.Contains(h2))
	h2.AddLink(&head.Link{Key: "foo", URL: "bar.com"})
	assert.True(t, h1.Contains(h2))
	h2.AddLink(&head.Link{Key: "foo2", URL: "bar2.com"})
	assert.True(t, h1.Contains(h2))
	h2.AddLink(&head.Link{Key: "foo3", URL: "bar3.com"})
	assert.False(t, h1.Contains(h2))
	h1.AddLink(&head.Link{Key: "foo3", URL: "bar3.com"})

	// Tags
	h1.Tags = append(h1.Tags, "foo")
	assert.True(t, h1.Contains(h2))
	h2.Tags = append(h2.Tags, "foo")
	assert.True(t, h1.Contains(h2))
	h2.Tags = append(h2.Tags, "foo2")
	assert.False(t, h1.Contains(h2))
	h1.Tags = append(h1.Tags, "foo2")
	assert.True(t, h1.Contains(h2))

	// Meta
	h1.Meta["foo"] = "bang"
	assert.True(t, h1.Contains(h2))
	h2.Meta["foo"] = "bang"
	assert.True(t, h1.Contains(h2))
	h2.Meta["foo2"] = "bo2"
	assert.False(t, h1.Contains(h2))
	h1.Meta["foo2"] = "bo2"
	assert.True(t, h1.Contains(h2))

	// Notes
	h1.Notes = "test notes"
	assert.True(t, h1.Contains(h2))
	h2.Notes = "bad notes"
	assert.False(t, h1.Contains(h2))
	h2.Notes = h1.Notes
	assert.True(t, h1.Contains(h2))
}

func TestHeaderSignNoJKU(t *testing.T) {
	// jku is intentionally not auto-stamped — generic JWT verifiers
	// resolve keys via <iss>/.well-known/jwks.json.
	priv := dsig.NewES256Key()
	h := head.NewHeader()
	h.UUID = uuid.V7()
	h.Digest = dsig.NewSHA256Digest([]byte(`{"x":1}`))

	sig, err := h.Sign(priv, cbc.URI("gobl:acme.example"), cbc.URI(""))
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	parsed, err := dsig.ParseSignature(sig.String())
	if err != nil {
		t.Fatalf("ParseSignature: %v", err)
	}
	assert.Equal(t, "", parsed.JKU())
}

func TestHeaderVerifyEnforcesValidityWindow(t *testing.T) {
	priv := dsig.NewES256Key()
	pub := priv.Public()

	h := head.NewHeader()
	h.UUID = uuid.V7()
	h.Digest = dsig.NewSHA256Digest([]byte(`{"x":1}`))

	sig, err := h.Sign(priv, cbc.URI(""), cbc.URI(""))
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}

	// No window on the key: verify succeeds.
	assert.NoError(t, h.Verify(sig, pub))

	// Window in the future: signed ts (now) is before valid_from.
	future := mustParseTS(t, "2099-01-01T00:00:00Z")
	pub.ValidFrom = &future
	pub.ValidUntil = nil
	err = h.Verify(sig, pub)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "before key's valid_from")

	// Window in the past: signed ts is after valid_until.
	past := mustParseTS(t, "2000-01-01T00:00:00Z")
	pub.ValidFrom = nil
	pub.ValidUntil = &past
	err = h.Verify(sig, pub)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "after key's valid_until")
}

func TestHeaderVerifyNoKeys(t *testing.T) {
	priv := dsig.NewES256Key()
	h := head.NewHeader()
	h.UUID = uuid.V7()
	h.Digest = dsig.NewSHA256Digest([]byte(`{"x":1}`))

	sig, err := h.Sign(priv, cbc.URI(""), cbc.URI(""))
	require.NoError(t, err)

	// No keys: signature is decoded with UnsafePayload and matched
	// against the header without cryptographic verification.
	assert.NoError(t, h.Verify(sig))

	// A header whose UUID doesn't match what was signed.
	other := head.NewHeader()
	other.UUID = uuid.V7()
	other.Digest = h.Digest
	err = other.Verify(sig)
	assert.ErrorIs(t, err, head.ErrSignatureMismatch)
}

func TestHeaderVerifyKeyMismatch(t *testing.T) {
	signer := dsig.NewES256Key()
	other := dsig.NewES256Key()

	h := head.NewHeader()
	h.UUID = uuid.V7()
	h.Digest = dsig.NewSHA256Digest([]byte(`{"x":1}`))

	sig, err := h.Sign(signer, cbc.URI(""), cbc.URI(""))
	require.NoError(t, err)

	// Verifying against a key that didn't sign it: no key matches.
	err = h.Verify(sig, other.Public())
	assert.ErrorIs(t, err, head.ErrSignatureKeyMismatch)
}

func TestSignedPayload(t *testing.T) {
	priv := dsig.NewES256Key()
	h := head.NewHeader()
	h.UUID = uuid.V7()
	h.Digest = dsig.NewSHA256Digest([]byte(`{"x":1}`))

	sig, err := h.Sign(priv,
		cbc.URI("gobl:alice.example"),
		cbc.URI("gobl:bob.example"))
	require.NoError(t, err)

	p, err := head.SignedPayload(sig)
	require.NoError(t, err)
	assert.Equal(t, h.UUID, p.UUID)
	assert.Equal(t, cbc.URI("gobl:alice.example"), p.Iss)
	assert.Equal(t, cbc.URI("gobl:bob.example"), p.Aud)
	assert.NotZero(t, p.IssuedAt)
}

func TestSignedPayloadDecodeError(t *testing.T) {
	// A signature whose payload is a JSON string (rather than an
	// object) fails to decode into a SigningPayload struct, hitting
	// SignedPayload's error branch.
	priv := dsig.NewES256Key()
	sig, err := dsig.NewSignature(priv, "not-an-object")
	require.NoError(t, err)
	_, err = head.SignedPayload(sig)
	require.Error(t, err)
	assert.ErrorIs(t, err, head.ErrSignaturePayload)
}

func TestHeaderVerifyPayloadDecodeError(t *testing.T) {
	priv := dsig.NewES256Key()
	sig, err := dsig.NewSignature(priv, "not-an-object")
	require.NoError(t, err)

	// Keyless path: UnsafePayload fails to decode into a SigningPayload.
	h := head.NewHeader()
	h.UUID = uuid.V7()
	h.Digest = dsig.NewSHA256Digest([]byte(`{"x":1}`))
	err = h.Verify(sig)
	require.Error(t, err)
	assert.ErrorIs(t, err, head.ErrSignaturePayload)
}

func TestHeaderMatchPayloadBothDigestNil(t *testing.T) {
	priv := dsig.NewES256Key()

	// Sign a header with nil digest; both sides end up with nil
	// digest in matchPayload, exercising the "both nil → ok" branch.
	h := head.NewHeader()
	h.UUID = uuid.V7()
	h.Digest = nil

	sig, err := h.Sign(priv, cbc.URI(""), cbc.URI(""))
	require.NoError(t, err)

	assert.NoError(t, h.Verify(sig, priv.Public()))
}

func TestHeaderMatchPayloadDigestEqualsFail(t *testing.T) {
	priv := dsig.NewES256Key()
	h := head.NewHeader()
	h.UUID = uuid.V7()
	h.Digest = dsig.NewSHA256Digest([]byte(`{"x":1}`))

	sig, err := h.Sign(priv, cbc.URI(""), cbc.URI(""))
	require.NoError(t, err)

	// Mutate the digest value after signing; Digest.Equals now reports a
	// mismatch rather than the field-nil branch.
	swapped := *h
	swapped.Digest = dsig.NewSHA256Digest([]byte(`{"x":2}`))
	err = swapped.Verify(sig, priv.Public())
	assert.ErrorIs(t, err, head.ErrSignatureMismatch)
}

func TestHeaderMatchPayloadDigestNil(t *testing.T) {
	priv := dsig.NewES256Key()

	// Sign a header that has a digest; then try to verify against a
	// header whose digest is nil — the mismatch hits the
	// "h.Digest == nil || actual.Digest == nil" branch.
	h := head.NewHeader()
	h.UUID = uuid.V7()
	h.Digest = dsig.NewSHA256Digest([]byte(`{"x":1}`))

	sig, err := h.Sign(priv, cbc.URI(""), cbc.URI(""))
	require.NoError(t, err)

	stripped := *h
	stripped.Digest = nil
	err = stripped.Verify(sig, priv.Public())
	assert.ErrorIs(t, err, head.ErrSignatureMismatch)
}

func mustParseTS(t *testing.T, s string) cal.Timestamp {
	t.Helper()
	ts, err := cal.ParseTimestamp(s)
	if err != nil {
		t.Fatalf("ParseTimestamp(%q): %v", s, err)
	}
	return ts
}
