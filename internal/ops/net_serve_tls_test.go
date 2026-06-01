package ops

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io"
	"math/big"
	stdnet "net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/invopop/gobl/net"
	"github.com/invopop/gobl/org"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// writeSelfSignedCert generates an ECDSA self-signed cert valid for
// localhost and 127.0.0.1, and writes the cert + key to disk. Returns
// the file paths.
func writeSelfSignedCert(t *testing.T, dir string) (certPath, keyPath string) {
	t.Helper()

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "localhost"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:     []string{"localhost"},
		IPAddresses:  []stdnet.IP{stdnet.ParseIP("127.0.0.1"), stdnet.ParseIP("::1")},
	}
	der, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	require.NoError(t, err)

	certPath = filepath.Join(dir, "cert.pem")
	keyPath = filepath.Join(dir, "key.pem")

	certOut, err := os.Create(certPath)
	require.NoError(t, err)
	defer certOut.Close() //nolint:errcheck
	require.NoError(t, pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: der}))

	keyOut, err := os.Create(keyPath)
	require.NoError(t, err)
	defer keyOut.Close() //nolint:errcheck
	keyDER, err := x509.MarshalECPrivateKey(priv)
	require.NoError(t, err)
	require.NoError(t, pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER}))

	return certPath, keyPath
}

// runServeOnListeners spins up NetServe with the given options against
// the supplied listeners and returns a stop function plus the listener
// addresses for client use.
func runServeOnListeners(t *testing.T, opts *NetServeOptions, tlsConfig *tls.Config) (httpAddr, httpsAddr string, stop func()) {
	t.Helper()

	handler, err := NetServeHandler(opts)
	require.NoError(t, err)

	httpHandler := handler
	httpsHandler := handler

	httpLn, err := stdnet.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	var httpsLn stdnet.Listener
	if tlsConfig != nil {
		httpsLn, err = stdnet.Listen("tcp", "127.0.0.1:0")
		require.NoError(t, err)
		httpsAddr = httpsLn.Addr().String()
	}
	httpAddr = httpLn.Addr().String()

	if opts.Out == nil {
		opts.Out = io.Discard
	}

	ctx, cancel := context.WithCancel(context.Background())
	doneCh := make(chan struct{})
	go func() {
		_ = serveOnListeners(ctx, opts, httpHandler, httpsHandler, tlsConfig, httpLn, httpsLn)
		close(doneCh)
	}()

	// Tiny wait so the goroutine's Serve calls are accepting before tests fire.
	time.Sleep(20 * time.Millisecond)

	return httpAddr, httpsAddr, func() {
		cancel()
		<-doneCh
	}
}

func TestNetServeFileTLS(t *testing.T) {
	dir := t.TempDir()
	certPath, keyPath := writeSelfSignedCert(t, dir)

	// Reuse the party + keys setup from net_serve_test.go.
	partyFile := filepath.Join(dir, "party.json")
	keysDir := filepath.Join(dir, "keys")
	privFile := filepath.Join(dir, "private.jwk")
	inboxDir := filepath.Join(dir, "inbox")

	signKey := privateKey
	writeRawParty(t, partyFile, &org.Party{Name: "TLS Party"})
	writeKey(t, keysDir, signKey)
	writePrivate(t, privFile, signKey)

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	require.NoError(t, err)
	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}

	opts := &NetServeOptions{
		PartyFile:      partyFile,
		KeysDir:        keysDir,
		PrivateKeyFile: privFile,
		InboxDir:       inboxDir,
		CertFile:       certPath,
		KeyFile:        keyPath,
	}

	httpAddr, httpsAddr, stop := runServeOnListeners(t, opts, tlsConfig)
	defer stop()

	// HTTP path: plain request works.
	httpResp, err := http.Get("http://" + httpAddr + net.KeyPath(signKey.ID()))
	require.NoError(t, err)
	defer httpResp.Body.Close() //nolint:errcheck
	assert.Equal(t, http.StatusOK, httpResp.StatusCode)

	// HTTPS path: must accept the self-signed cert.
	tlsClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec // test only
		},
	}
	httpsResp, err := tlsClient.Get("https://" + httpsAddr + net.KeyPath(signKey.ID()))
	require.NoError(t, err)
	defer httpsResp.Body.Close() //nolint:errcheck
	assert.Equal(t, http.StatusOK, httpsResp.StatusCode)
}

func TestListenTCPGenericError(t *testing.T) {
	// Bind the same port twice — the second listen returns an error
	// that is NOT EACCES, exercising the bare-error wrap branch.
	first, err := listenTCP(0)
	require.NoError(t, err)
	defer first.Close() //nolint:errcheck
	port := first.Addr().(*stdnet.TCPAddr).Port
	_, err = listenTCP(port)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "net serve: listen")
}

func TestListenTCPEACCES(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("EACCES semantics differ on Windows")
	}
	if os.Geteuid() == 0 {
		t.Skip("test must run as a non-root user")
	}
	// Port 1 is privileged on Unix-like systems; non-root cannot bind.
	_, err := listenTCP(1)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "permission denied")
	assert.Contains(t, err.Error(), "--http-port")
	assert.Contains(t, err.Error(), "setcap")
}

func TestNewAutocertManagerLive(t *testing.T) {
	dir := t.TempDir()
	m := newAutocertManager(&NetServeOptions{ACMELive: true, CertDir: dir, ACMEEmail: "ops@example.com"}, []string{"example.com"})
	require.NotNil(t, m)
	assert.Equal(t, "ops@example.com", m.Email)
	assert.Nil(t, m.Client, "live mode uses autocert's default LE production directory")
	// Verify HostPolicy accepts the configured domain and rejects others.
	assert.NoError(t, m.HostPolicy(context.Background(), "example.com"))
	assert.Error(t, m.HostPolicy(context.Background(), "other.example.com"))
	// DirCache stores under the supplied path; round-trip via Put/Get to confirm.
	require.NoError(t, m.Cache.Put(context.Background(), "probe", []byte("ok")))
	got, err := m.Cache.Get(context.Background(), "probe")
	require.NoError(t, err)
	assert.Equal(t, "ok", string(got))
}

func TestNewAutocertManagerTest(t *testing.T) {
	m := newAutocertManager(&NetServeOptions{ACMETest: true}, []string{"example.com"})
	require.NotNil(t, m.Client, "test mode must override the Client to point at LE staging")
	assert.True(t, strings.HasPrefix(m.Client.DirectoryURL, "https://acme-staging-v02.api.letsencrypt.org"))
}
