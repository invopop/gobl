package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"testing"

	"github.com/invopop/gobl/dsig"
	"github.com/stretchr/testify/assert"
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
			err: "code=400, message=public key required",
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
			err: "code=422, message=key mismatch",
		}
	})
	tests.Add("draft", func(t *testing.T) interface{} {
		return tt{
			in:  testFileReader(t, "testdata/draft.json"),
			key: publicKey,
			err: "code=422, message=document is a draft",
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
