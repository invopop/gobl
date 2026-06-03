package dsig

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/invopop/gobl/cal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPublicKeyFromJWK(t *testing.T) {
	priv := NewES256Key()
	pk := new(PublicKey)
	raw, err := json.Marshal(priv.Public())
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(raw, pk))
	pub, err := NewPublicKey(*pk.jwk)
	require.NoError(t, err)
	require.NotNil(t, pub)
	assert.Equal(t, priv.ID(), pub.ID())
}

func TestNewPublicKeyInvalid(t *testing.T) {
	// A zero-value JOSE JSONWebKey has .Valid()==false and no kid, so
	// NewPublicKey rejects it.
	_, err := NewPublicKey(jose.JSONWebKey{})
	require.Error(t, err)
}

func TestPrivateKeySignatureAlgorithmUnknown(t *testing.T) {
	priv := NewES256Key()
	priv.jwk.Key = "not an ecdsa key"
	_, err := priv.signatureAlgorithm()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unrecognized")
}

func TestKeyThumbprintError(t *testing.T) {
	// A JWK with a nil underlying Key fails Thumbprint; the helper
	// returns an empty string.
	priv := NewES256Key()
	priv.jwk.Key = nil
	assert.Equal(t, "", priv.Thumbprint())
}

func TestPublicKeyMarshalWithWindow(t *testing.T) {
	priv := NewES256Key()
	pk := new(PublicKey)
	raw, err := json.Marshal(priv.Public())
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(raw, pk))

	t.Run("plain JWK when no window", func(t *testing.T) {
		data, err := json.Marshal(pk)
		require.NoError(t, err)
		var m map[string]any
		require.NoError(t, json.Unmarshal(data, &m))
		_, ok := m["valid_from"]
		assert.False(t, ok)
	})

	t.Run("round-trip with window", func(t *testing.T) {
		from := cal.TimestampOf(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
		until := cal.TimestampOf(time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC))
		pk.ValidFrom = &from
		pk.ValidUntil = &until
		data, err := json.Marshal(pk)
		require.NoError(t, err)
		out := new(PublicKey)
		require.NoError(t, json.Unmarshal(data, out))
		require.NotNil(t, out.ValidFrom)
		require.NotNil(t, out.ValidUntil)
		assert.True(t, from.Time.Equal(out.ValidFrom.Time))
		assert.True(t, until.Time.Equal(out.ValidUntil.Time))
	})
}

func TestPublicKeyAllows(t *testing.T) {
	from := cal.TimestampOf(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	until := cal.TimestampOf(time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC))
	pk := &PublicKey{ValidFrom: &from, ValidUntil: &until}

	before := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
	inside := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	after := time.Date(2027, 2, 1, 0, 0, 0, 0, time.UTC)

	assert.NoError(t, pk.Allows(inside))
	assert.Error(t, pk.Allows(before))
	assert.Error(t, pk.Allows(after))
	assert.NoError(t, pk.Allows(time.Time{}), "zero-value t (no iat) is permitted")

	// Unbounded key accepts any t.
	assert.NoError(t, (&PublicKey{}).Allows(inside))
}

func TestPublicKeyUnmarshalBad(t *testing.T) {
	pk := new(PublicKey)
	err := json.Unmarshal([]byte("not json"), pk)
	require.Error(t, err)
}

func TestPublicKeyUnmarshalBadValidFrom(t *testing.T) {
	pk := new(PublicKey)
	// valid_from is a number — fails to decode as cal.Timestamp.
	err := json.Unmarshal([]byte(`{"kty":"EC","crv":"P-256","kid":"x","x":"AA","y":"BB","valid_from":12345}`), pk)
	require.Error(t, err)
}
