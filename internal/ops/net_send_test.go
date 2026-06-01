package ops

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/net"
	"github.com/invopop/gobl/note"
	"github.com/invopop/gobl/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func signedNoteEnvelope(t *testing.T, content string) []byte {
	t.Helper()
	msg := &note.Message{Content: content}
	msg.SetUUID(uuid.V7())
	env, err := gobl.Envelop(msg)
	require.NoError(t, err)
	require.NoError(t, env.Sign(testPeerKey, net.Address(testPeerDomain).URI(), net.Address(testServeDomain).URI()))
	body, err := json.Marshal(env)
	require.NoError(t, err)
	return body
}

func TestNetSendSuccess(t *testing.T) {
	body := signedNoteEnvelope(t, "round trip")

	var received []byte
	var receivedContentType string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, net.InboxPath, r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		receivedContentType = r.Header.Get("Content-Type")
		received, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	to := hostFromURL(t, srv.URL)
	err := NetSend(context.Background(), &NetSendOptions{
		Input:    bytes.NewReader(body),
		To:       net.Address(to),
		Insecure: true,
	})
	require.NoError(t, err)
	assert.Equal(t, "application/json", receivedContentType)
	assert.JSONEq(t, string(body), string(received))
}

func TestNetSendRejectsBadEnvelope(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	err := NetSend(context.Background(), &NetSendOptions{
		Input:    bytes.NewReader([]byte("not json")),
		To:       net.Address(hostFromURL(t, srv.URL)),
		Insecure: true,
	})
	require.Error(t, err)
}

func TestNetSendNon202(t *testing.T) {
	body := signedNoteEnvelope(t, "bad sig path")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "no thanks", http.StatusUnauthorized)
	}))
	defer srv.Close()

	err := NetSend(context.Background(), &NetSendOptions{
		Input:    bytes.NewReader(body),
		To:       net.Address(hostFromURL(t, srv.URL)),
		Insecure: true,
	})
	require.Error(t, err)
	assert.True(t, errors.Is(err, net.ErrInboxRejected))
}

func TestNetSendInsecureURL(t *testing.T) {
	got, err := inboxURL("localhost:8080", true)
	require.NoError(t, err)
	assert.Equal(t, "http://localhost:8080/.well-known/gobl/inbox", got)
}

func TestNetSendSecureURL(t *testing.T) {
	got, err := inboxURL("example.com", false)
	require.NoError(t, err)
	assert.Equal(t, "https://example.com/.well-known/gobl/inbox", got)
}

// hostFromURL strips the scheme from an httptest.NewServer URL so it
// can be used as a `host:port`-form GOBL Net address in --insecure mode.
func hostFromURL(t *testing.T, raw string) string {
	t.Helper()
	u, err := url.Parse(raw)
	require.NoError(t, err)
	return strings.TrimPrefix(u.Host, "")
}

func TestNetSendInboxURLErrors(t *testing.T) {
	t.Run("insecure empty address", func(t *testing.T) {
		_, err := inboxURL("", true)
		require.Error(t, err)
		assert.True(t, errors.Is(err, net.ErrAddressEmpty))
	})
	t.Run("secure invalid FQDN", func(t *testing.T) {
		_, err := inboxURL("localhost", false)
		require.Error(t, err)
	})
}

func TestNetSendInvalidEnvelopeJSON(t *testing.T) {
	// Looks like JSON but fails to unmarshal into Envelope.
	err := NetSend(context.Background(), &NetSendOptions{
		Input:    bytes.NewReader([]byte("[1,2,3]")),
		To:       net.Address("example.com"),
		Insecure: true,
	})
	require.Error(t, err)
}

func TestNetSendInputReadError(t *testing.T) {
	err := NetSend(context.Background(), &NetSendOptions{
		Input:    errReader{},
		To:       net.Address("example.com"),
		Insecure: true,
	})
	require.Error(t, err)
}

type errReader struct{}

func (errReader) Read([]byte) (int, error) { return 0, io.ErrUnexpectedEOF }

func TestNetSendInboxURLError(t *testing.T) {
	body := signedNoteEnvelope(t, "x")
	err := NetSend(context.Background(), &NetSendOptions{
		Input:    bytes.NewReader(body),
		To:       net.Address("localhost"), // single label fails FQDN validation
		Insecure: false,
	})
	require.Error(t, err)
}

func TestNetSendTransportError(t *testing.T) {
	body := signedNoteEnvelope(t, "x")
	// Send to a closed loopback port so client.Do returns an error.
	err := NetSend(context.Background(), &NetSendOptions{
		Input:    bytes.NewReader(body),
		To:       net.Address("127.0.0.1:1"),
		Insecure: true,
		Client:   &http.Client{Timeout: 100 * 1000 * 1000}, // 100ms
	})
	require.Error(t, err)
}

func TestNetSendRoundTrip(t *testing.T) {
	srv, inboxDir := setupNetServer(t)

	body := signedNoteEnvelope(t, "round trip via serve")
	err := NetSend(context.Background(), &NetSendOptions{
		Input:    bytes.NewReader(body),
		To:       net.Address(hostFromURL(t, srv.URL)),
		Insecure: true,
	})
	require.NoError(t, err)

	// Confirm the envelope landed in the inbox directory.
	files, err := readDirNames(inboxDir)
	require.NoError(t, err)
	require.Len(t, files, 1)
	assert.True(t, strings.HasSuffix(files[0], ".json"))
}
