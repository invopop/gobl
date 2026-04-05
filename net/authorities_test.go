package net

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/go-jose/go-jose/v4"
	"github.com/invopop/gobl/dsig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVerifyKeySetPinned(t *testing.T) {
	ctx := context.Background()
	senderKey := dsig.NewES256Key()
	authorityKey := dsig.NewES256Key()

	jwk := makeJWK(t, senderKey)
	ks, err := NewKeySet(jwk)
	require.NoError(t, err)

	err = ks.Sign(authorityKey, dsig.WithGN("authority.gobl.net"))
	require.NoError(t, err)

	// Round-trip to ensure signatures parse correctly
	data, err := json.Marshal(ks)
	require.NoError(t, err)
	ks = new(KeySet)
	require.NoError(t, json.Unmarshal(data, ks))

	// Client with pinned authority key
	c := NewClient(
		WithFetcher(&mockFetcher{}), // should not be called
		WithAuthorities(authorityKey.Public()),
	)

	result, err := c.VerifyKeySet(ctx, ks)
	require.NoError(t, err)
	assert.True(t, result.Signed)
	assert.True(t, result.Pinned)
	assert.Empty(t, result.Authority)
}

func TestVerifyKeySetRemote(t *testing.T) {
	ctx := context.Background()
	senderKey := dsig.NewES256Key()
	authorityKey := dsig.NewES256Key()

	jwk := makeJWK(t, senderKey)
	ks, err := NewKeySet(jwk)
	require.NoError(t, err)

	err = ks.Sign(authorityKey, dsig.WithGN("authority.gobl.net"))
	require.NoError(t, err)

	// Round-trip
	data, err := json.Marshal(ks)
	require.NoError(t, err)
	ks = new(KeySet)
	require.NoError(t, json.Unmarshal(data, ks))

	// Build the authority's JWKS response
	authorityJWKS := jwksFromKey(t, authorityKey)

	// Client with NO pinned keys — will fetch from authority's address
	c := NewClient(WithFetcher(&mockFetcher{data: authorityJWKS}))

	result, err := c.VerifyKeySet(ctx, ks)
	require.NoError(t, err)
	assert.True(t, result.Signed)
	assert.False(t, result.Pinned)
	assert.Equal(t, Address("authority.gobl.net"), result.Authority)
}

func TestVerifyKeySetUnsigned(t *testing.T) {
	ctx := context.Background()
	senderKey := dsig.NewES256Key()

	jwk := makeJWK(t, senderKey)
	ks := &KeySet{Keys: []jose.JSONWebKey{jwk}}

	c := NewClient()

	result, err := c.VerifyKeySet(ctx, ks)
	require.NoError(t, err)
	assert.False(t, result.Signed)
}

func TestVerifyKeySetNoMatch(t *testing.T) {
	ctx := context.Background()
	senderKey := dsig.NewES256Key()
	authorityKey := dsig.NewES256Key()
	wrongKey := dsig.NewES256Key()

	jwk := makeJWK(t, senderKey)
	ks, err := NewKeySet(jwk)
	require.NoError(t, err)

	err = ks.Sign(authorityKey, dsig.WithGN("authority.gobl.net"))
	require.NoError(t, err)

	// Round-trip
	data, err := json.Marshal(ks)
	require.NoError(t, err)
	ks = new(KeySet)
	require.NoError(t, json.Unmarshal(data, ks))

	// Serve the WRONG authority's keys
	wrongJWKS := jwksFromKey(t, wrongKey)
	c := NewClient(WithFetcher(&mockFetcher{data: wrongJWKS}))

	_, err = c.VerifyKeySet(ctx, ks)
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrVerifyFailed))
}
