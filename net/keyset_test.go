package net

import (
	"encoding/json"
	"testing"

	"github.com/go-jose/go-jose/v4"
	"github.com/invopop/gobl/dsig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeJWK(t *testing.T, key *dsig.PrivateKey) jose.JSONWebKey {
	t.Helper()
	pub := key.Public()
	data, err := json.Marshal(pub)
	require.NoError(t, err)
	var jwk jose.JSONWebKey
	require.NoError(t, json.Unmarshal(data, &jwk))
	return jwk
}

func TestKeySetDigest(t *testing.T) {
	key := dsig.NewES256Key()
	jwk := makeJWK(t, key)

	ks, err := NewKeySet(jwk)
	require.NoError(t, err)
	require.NotNil(t, ks.Digest)
	assert.Equal(t, dsig.DigestSHA256, ks.Digest.Algorithm)
	assert.NotEmpty(t, ks.Digest.Value)

	// Digest should be stable
	ks2, err := NewKeySet(jwk)
	require.NoError(t, err)
	assert.Equal(t, ks.Digest.Value, ks2.Digest.Value)
}

func TestKeySetCalculate(t *testing.T) {
	key1 := dsig.NewES256Key()
	key2 := dsig.NewES256Key()
	jwk1 := makeJWK(t, key1)
	jwk2 := makeJWK(t, key2)

	ks, err := NewKeySet(jwk1)
	require.NoError(t, err)
	origDigest := ks.Digest.Value

	// Add a key and recalculate
	ks.Keys = append(ks.Keys, jwk2)
	require.NoError(t, ks.Calculate())
	assert.NotEqual(t, origDigest, ks.Digest.Value)
}

func TestKeySetSign(t *testing.T) {
	senderKey := dsig.NewES256Key()
	authorityKey := dsig.NewES256Key()

	jwk := makeJWK(t, senderKey)
	ks, err := NewKeySet(jwk)
	require.NoError(t, err)

	assert.False(t, ks.Signed())

	// Sign the KeySet
	err = ks.Sign(authorityKey, dsig.WithGN("authority.gobl.net"))
	require.NoError(t, err)
	assert.True(t, ks.Signed())
	assert.Len(t, ks.Signatures, 1)

	// Verify with the authority's public key
	err = ks.Verify(authorityKey.Public())
	require.NoError(t, err)

	// Verify fails with wrong key
	wrongKey := dsig.NewES256Key()
	err = ks.Verify(wrongKey.Public())
	require.Error(t, err)
}

func TestKeySetSignRoundTrip(t *testing.T) {
	senderKey := dsig.NewES256Key()
	authorityKey := dsig.NewES256Key()

	jwk := makeJWK(t, senderKey)
	ks, err := NewKeySet(jwk)
	require.NoError(t, err)

	err = ks.Sign(authorityKey, dsig.WithGN("authority.gobl.net"))
	require.NoError(t, err)

	// Round-trip through JSON
	data, err := json.Marshal(ks)
	require.NoError(t, err)

	ks2 := new(KeySet)
	require.NoError(t, json.Unmarshal(data, ks2))

	assert.True(t, ks2.Signed())
	assert.Len(t, ks2.Keys, 1)
	assert.NotNil(t, ks2.Digest)

	// Verify the round-tripped KeySet
	err = ks2.Verify(authorityKey.Public())
	require.NoError(t, err)
}

func TestKeySetUnsigned(t *testing.T) {
	// Standard JWKS without dig/sigs
	data := []byte(`{"keys":[]}`)
	ks := new(KeySet)
	require.NoError(t, json.Unmarshal(data, ks))

	assert.False(t, ks.Signed())
	assert.Nil(t, ks.Digest)
	assert.Empty(t, ks.Signatures)
}

func TestKeySetVerifyNoSignatures(t *testing.T) {
	key := dsig.NewES256Key()
	jwk := makeJWK(t, key)
	ks := &KeySet{Keys: []jose.JSONWebKey{jwk}}

	err := ks.Verify()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no signatures")
}

func TestKeySetVerifyDigestTampered(t *testing.T) {
	senderKey := dsig.NewES256Key()
	authorityKey := dsig.NewES256Key()

	jwk := makeJWK(t, senderKey)
	ks, err := NewKeySet(jwk)
	require.NoError(t, err)

	err = ks.Sign(authorityKey)
	require.NoError(t, err)

	// Tamper with the keys after signing
	extraKey := dsig.NewES256Key()
	ks.Keys = append(ks.Keys, makeJWK(t, extraKey))

	// Verification should fail because the digest no longer matches
	err = ks.Verify(authorityKey.Public())
	require.Error(t, err)
}

func TestKeySetKey(t *testing.T) {
	key1 := dsig.NewES256Key()
	key2 := dsig.NewES256Key()

	ks := &KeySet{
		Keys: []jose.JSONWebKey{makeJWK(t, key1), makeJWK(t, key2)},
	}

	t.Run("found first", func(t *testing.T) {
		pk, err := ks.Key(key1.ID())
		require.NoError(t, err)
		assert.Equal(t, key1.ID(), pk.ID())
	})

	t.Run("found second", func(t *testing.T) {
		pk, err := ks.Key(key2.ID())
		require.NoError(t, err)
		assert.Equal(t, key2.ID(), pk.ID())
	})

	t.Run("not found", func(t *testing.T) {
		_, err := ks.Key("nonexistent")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrKeyNotFound)
	})
}
