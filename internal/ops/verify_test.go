package ops

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"testing"

	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/net"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/flimzy/testy"
)

func signedDoc(t *testing.T) []byte {
	t.Helper()
	env, err := Sign(context.Background(), &SignOptions{
		ParseOptions: &ParseOptions{
			Input: testFileReader(t, "testdata/invoice-es-es.env.yaml"),
			SetFile: map[string]string{
				"doc": "testdata/invoice-es-es.yaml",
			},
		},
		PrivateKey: privateKey,
	})
	if err != nil {
		t.Fatal(err)
	}
	out, err := json.Marshal(env)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func TestVerify(t *testing.T) {
	type tt struct {
		in  io.Reader
		key *dsig.PublicKey
		err string
	}

	tests := testy.NewTable()
	tests.Add("validation pass", func(t *testing.T) interface{} {
		return tt{
			in:  bytes.NewReader(signedDoc(t)),
			key: publicKey,
		}
	})
	tests.Add("missing key", func(t *testing.T) interface{} {
		return tt{
			in:  testFileReader(t, "testdata/success.json"),
			err: "input: public key required",
		}
	})
	tests.Add("wrong public key", func(t *testing.T) interface{} {
		wrongKey := new(dsig.PublicKey)
		if err := json.Unmarshal([]byte(`{"use":"sig","kty":"EC","kid":"0f8726f9-5e1c-43e2-a0c8-55e696550800","crv":"P-256","alg":"ES256","x":"rmBvuoO4LJCnuo0htUaAhSgF7OdCEO6GAuI_o1dUP8s","y":"7qXu5-LHoQM52C_HJJ4eYoGlQnesFLcT76hIlDs-ZfU"}`), wrongKey); err != nil {
			panic(err)
		}

		return tt{
			in:  testFileReader(t, "testdata/success.json"),
			key: wrongKey,
			err: "signature: key mismatch",
		}
	})
	tests.Add("draft", func(t *testing.T) interface{} {
		return tt{
			in:  testFileReader(t, "testdata/draft.json"),
			key: publicKey,
			err: "signature: envelope is not signed",
		}
	})

	tests.Run(t, func(t *testing.T, tt tt) {
		t.Parallel()
		err := Verify(context.Background(), tt.in, tt.key)
		if tt.err == "" {
			assert.Nil(t, err)
		} else {
			assert.EqualError(t, err, tt.err)
		}
	})
}

func TestVerifyInvalidYAML(t *testing.T) {
	err := Verify(context.Background(), bytes.NewReader([]byte("\t\t\t\n@@:")), publicKey)
	require.Error(t, err)
}

func TestVerifyValidationFail(t *testing.T) {
	// Envelope missing digest -> validation fails (not a signature error).
	err := Verify(context.Background(), bytes.NewReader([]byte(`{}`)), publicKey)
	require.Error(t, err)
}

func TestVerifyReadError(t *testing.T) {
	err := Verify(context.Background(), errReader{}, publicKey)
	require.Error(t, err)
}

func TestVerifyRemote(t *testing.T) {
	addr := net.Address("billing.invopop.com")
	// Sign with iss set so VerifyRemote can resolve the issuer.
	env, err := Sign(context.Background(), &SignOptions{
		ParseOptions: &ParseOptions{
			Input: testFileReader(t, "testdata/invoice-es-es.env.yaml"),
			SetFile: map[string]string{
				"doc": "testdata/invoice-es-es.yaml",
			},
		},
		PrivateKey: privateKey,
		Iss:        addr.URI(),
	})
	require.NoError(t, err)
	body, err := json.Marshal(env)
	require.NoError(t, err)

	// Build a fetcher that serves the matching public key for whichever
	// kid is requested.
	pkBytes, err := json.Marshal(privateKey.Public())
	require.NoError(t, err)
	c := net.NewClient(net.WithFetcher(&mapFetcher{data: map[string][]byte{
		addr.KeyURL(privateKey.ID()): pkBytes,
	}}))

	t.Run("success no address pin", func(t *testing.T) {
		assert.NoError(t, VerifyRemote(context.Background(), bytes.NewReader(body), c, ""))
	})

	t.Run("address pin matches", func(t *testing.T) {
		assert.NoError(t, VerifyRemote(context.Background(), bytes.NewReader(body), c, addr))
	})

	t.Run("address pin mismatch", func(t *testing.T) {
		err := VerifyRemote(context.Background(), bytes.NewReader(body), c, net.Address("other.example"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "expected other.example")
	})

	t.Run("invalid JSON", func(t *testing.T) {
		err := VerifyRemote(context.Background(), bytes.NewReader([]byte("\t\t\t\n@@:")), c, "")
		require.Error(t, err)
	})

	t.Run("read error", func(t *testing.T) {
		err := VerifyRemote(context.Background(), errReader{}, c, "")
		require.Error(t, err)
	})

	t.Run("validation fails", func(t *testing.T) {
		err := VerifyRemote(context.Background(), bytes.NewReader([]byte(`{}`)), c, "")
		require.Error(t, err)
	})

	t.Run("verify fails", func(t *testing.T) {
		// Tamper with the envelope so VerifyEnvelope rejects it.
		var env map[string]any
		require.NoError(t, json.Unmarshal(body, &env))
		// Remove signatures so Signed() returns false.
		delete(env, "sigs")
		tampered, err := json.Marshal(env)
		require.NoError(t, err)
		err = VerifyRemote(context.Background(), bytes.NewReader(tampered), c, "")
		require.Error(t, err)
	})
}
