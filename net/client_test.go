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

type mockFetcher struct {
	data []byte
	err  error
	url  string // records the URL that was fetched
}

func (m *mockFetcher) Fetch(_ context.Context, url string) ([]byte, error) {
	m.url = url
	return m.data, m.err
}

func TestFetchKeySet(t *testing.T) {
	ctx := context.Background()
	key := dsig.NewES256Key()
	pub := key.Public()

	// Build a valid JWKS response
	pubData, err := json.Marshal(pub)
	require.NoError(t, err)
	jwksData := []byte(`{"keys":[` + string(pubData) + `]}`)

	t.Run("success", func(t *testing.T) {
		mock := &mockFetcher{data: jwksData}
		c := NewClient(WithFetcher(mock))

		ks, err := c.FetchKeySet(ctx, Address("billing.invopop.com"))
		require.NoError(t, err)
		assert.NotNil(t, ks)
		assert.Equal(t, "https://billing.invopop.com/.well-known/gobl/jwks.json", mock.url)
		assert.Len(t, ks.Keys, 1)
		assert.Equal(t, key.ID(), ks.Keys[0].KeyID)
	})

	t.Run("fetch error", func(t *testing.T) {
		mock := &mockFetcher{err: ErrFetchFailed}
		c := NewClient(WithFetcher(mock))

		_, err := c.FetchKeySet(ctx, Address("billing.invopop.com"))
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrFetchFailed))
	})

	t.Run("invalid JSON", func(t *testing.T) {
		mock := &mockFetcher{data: []byte("not json")}
		c := NewClient(WithFetcher(mock))

		_, err := c.FetchKeySet(ctx, Address("billing.invopop.com"))
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrFetchFailed))
	})

	t.Run("invalid address", func(t *testing.T) {
		mock := &mockFetcher{data: jwksData}
		c := NewClient(WithFetcher(mock))

		_, err := c.FetchKeySet(ctx, Address(""))
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrAddressEmpty))
	})
}

func TestFetchPublicKey(t *testing.T) {
	ctx := context.Background()
	key := dsig.NewES256Key()
	pub := key.Public()

	pubData, err := json.Marshal(pub)
	require.NoError(t, err)
	jwksData := []byte(`{"keys":[` + string(pubData) + `]}`)

	t.Run("found", func(t *testing.T) {
		mock := &mockFetcher{data: jwksData}
		c := NewClient(WithFetcher(mock))

		pk, err := c.FetchPublicKey(ctx, Address("billing.invopop.com"), key.ID())
		require.NoError(t, err)
		assert.Equal(t, key.ID(), pk.ID())
	})

	t.Run("not found", func(t *testing.T) {
		mock := &mockFetcher{data: jwksData}
		c := NewClient(WithFetcher(mock))

		_, err := c.FetchPublicKey(ctx, Address("billing.invopop.com"), "nonexistent-kid")
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrKeyNotFound))
	})
}

func TestKeySetKeyLookup(t *testing.T) {
	key := dsig.NewES256Key()
	pub := key.Public()

	pubData, err := json.Marshal(pub)
	require.NoError(t, err)

	var jwk jose.JSONWebKey
	require.NoError(t, json.Unmarshal(pubData, &jwk))

	ks := &KeySet{Keys: []jose.JSONWebKey{jwk}}

	t.Run("found", func(t *testing.T) {
		pk, err := ks.Key(key.ID())
		require.NoError(t, err)
		assert.Equal(t, key.ID(), pk.ID())
	})

	t.Run("not found", func(t *testing.T) {
		_, err := ks.Key("wrong-kid")
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrKeyNotFound))
	})
}
