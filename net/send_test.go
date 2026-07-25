package net

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockPoster is a Fetcher+Poster that records the POST it receives.
type mockPoster struct {
	mapFetcher
	url    string
	body   []byte
	header http.Header
	err    error
}

func (m *mockPoster) Post(_ context.Context, url string, body []byte, header http.Header) error {
	m.url = url
	m.body = body
	m.header = header
	return m.err
}

func buildSignedEnvelope(t *testing.T, key *dsig.PrivateKey, iss, aud Address) *gobl.Envelope {
	t.Helper()
	party := &org.Party{Name: "Test Sender"}
	party.SetUUID(uuid.V7())
	env, err := gobl.Envelop(party)
	require.NoError(t, err)
	require.NoError(t, env.Sign(key,
		head.WithIssuer(iss.String()),
		head.WithAudience(aud.String())))
	return env
}

func TestSend(t *testing.T) {
	ctx := context.Background()
	sender := Address("sender.example.com")
	receiver := Address("receiver.example.com")
	key := dsig.NewES256Key()
	env := buildSignedEnvelope(t, key, sender, receiver)

	t.Run("posts the envelope with a bearer token", func(t *testing.T) {
		poster := new(mockPoster)
		c := NewClient(WithFetcher(poster), WithIdentity(sender, key))
		require.NoError(t, c.Send(ctx, receiver, env))

		assert.Equal(t, receiver.InboxURL(), poster.url)
		posted := new(gobl.Envelope)
		require.NoError(t, json.Unmarshal(poster.body, posted))
		assert.Equal(t, env.Head.UUID, posted.Head.UUID)

		auth := poster.header.Get("Authorization")
		require.NotEmpty(t, auth)
		claims := tokenPayload(t, auth[len("Bearer "):])
		assert.Equal(t, sender, claims.Iss)
		assert.Equal(t, receiver, claims.Aud)
	})

	t.Run("canonicalizes the address", func(t *testing.T) {
		poster := new(mockPoster)
		c := NewClient(WithFetcher(poster), WithIdentity(sender, key))
		require.NoError(t, c.Send(ctx, "RECEIVER.example.com.", env))
		assert.Equal(t, receiver.InboxURL(), poster.url)
	})

	t.Run("propagates a rejection", func(t *testing.T) {
		poster := &mockPoster{err: ErrInboxRejected}
		c := NewClient(WithFetcher(poster), WithIdentity(sender, key))
		err := c.Send(ctx, receiver, env)
		assert.True(t, errors.Is(err, ErrInboxRejected))
	})

	t.Run("rejects a fetcher without POST support", func(t *testing.T) {
		c := NewClient(WithFetcher(&mapFetcher{}))
		err := c.Send(ctx, receiver, env)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrFetchFailed))
		assert.Contains(t, err.Error(), "does not support POST")
	})

	t.Run("rejects an invalid address", func(t *testing.T) {
		c := NewClient(WithFetcher(new(mockPoster)))
		err := c.Send(ctx, "not valid!", env)
		assert.True(t, errors.Is(err, ErrAddressInvalid))
	})
}

func TestHTTPFetcherPost(t *testing.T) {
	ctx := context.Background()

	t.Run("202 is success and headers are forwarded", func(t *testing.T) {
		var gotAuth, gotCT string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			gotAuth = r.Header.Get("Authorization")
			gotCT = r.Header.Get("Content-Type")
			w.WriteHeader(http.StatusAccepted)
		}))
		defer srv.Close()

		h := http.Header{}
		h.Set("Authorization", "Bearer token")
		err := newHTTPFetcher(true).Post(ctx, srv.URL, []byte(`{}`), h)
		require.NoError(t, err)
		assert.Equal(t, "Bearer token", gotAuth)
		assert.Equal(t, "application/json", gotCT)
	})

	t.Run("4xx reports ErrInboxRejected", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "nope", http.StatusUnauthorized)
		}))
		defer srv.Close()

		err := newHTTPFetcher(true).Post(ctx, srv.URL, []byte(`{}`), nil)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrInboxRejected))
		assert.Contains(t, err.Error(), "HTTP 401")
	})

	t.Run("429 is retryable, not a rejection", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "later", http.StatusTooManyRequests)
		}))
		defer srv.Close()

		err := newHTTPFetcher(true).Post(ctx, srv.URL, []byte(`{}`), nil)
		require.Error(t, err)
		assert.False(t, errors.Is(err, ErrInboxRejected))
		assert.True(t, errors.Is(err, ErrFetchFailed))
	})

	t.Run("5xx reports ErrFetchFailed", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "boom", http.StatusInternalServerError)
		}))
		defer srv.Close()

		err := newHTTPFetcher(true).Post(ctx, srv.URL, []byte(`{}`), nil)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrFetchFailed))
	})

	t.Run("transport error", func(t *testing.T) {
		// Unreachable address (port 1 is privileged + usually closed).
		f := &HTTPFetcher{Client: &http.Client{Timeout: 100 * time.Millisecond}}
		err := f.Post(ctx, "http://127.0.0.1:1/x", []byte(`{}`), nil)
		require.Error(t, err)
		assert.True(t, errors.Is(err, ErrFetchFailed))
	})
}

func TestSendNilEnvelope(t *testing.T) {
	c := NewClient(WithFetcher(new(mockPoster)))
	err := c.Send(context.Background(), "receiver.example.com", nil)
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrFetchFailed))
	assert.Contains(t, err.Error(), "envelope is nil")
}
