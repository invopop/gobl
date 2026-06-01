package ops

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/head"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/flimzy/testy"
)

const noteMessageJSON = `{"$schema":"https://gobl.org/draft-0/note/message","content":"hi"}`

// signIss signs the note message with the given iss/aud and round-trips
// the envelope through JSON (the signed payload is read after parse).
func signIss(t *testing.T, iss, aud cbc.URI) *gobl.Envelope {
	t.Helper()
	env, err := Sign(context.Background(), &SignOptions{
		ParseOptions: &ParseOptions{Input: strings.NewReader(noteMessageJSON)},
		PrivateKey:   privateKey,
		Iss:          iss,
		Aud:          aud,
	})
	require.NoError(t, err)
	data, err := json.Marshal(env)
	require.NoError(t, err)
	out := new(gobl.Envelope)
	require.NoError(t, json.Unmarshal(data, out))
	return out
}

func TestSignWithIss(t *testing.T) {
	env := signIss(t, "gobl:billing.invopop.com", "gobl:acme.example")
	require.True(t, env.Signed())
	p, err := head.SignedPayload(env.Signatures[0])
	require.NoError(t, err)
	assert.Equal(t, cbc.URI("gobl:billing.invopop.com"), p.Iss)
	assert.Equal(t, cbc.URI("gobl:acme.example"), p.Aud)
}

func TestSignWithoutIss(t *testing.T) {
	env := signIss(t, "", "")
	require.True(t, env.Signed())
	p, err := head.SignedPayload(env.Signatures[0])
	require.NoError(t, err)
	assert.Empty(t, p.Iss)
}

func TestSignInvalidKey(t *testing.T) {
	// Sign should propagate env.Sign's error when the private key is
	// zero-valued (no underlying JWK).
	_, err := Sign(context.Background(), &SignOptions{
		ParseOptions: &ParseOptions{Input: strings.NewReader(noteMessageJSON)},
		PrivateKey:   &dsig.PrivateKey{},
	})
	require.Error(t, err)
}

// TestSignSetsTimestamp asserts that signing automatically stamps the
// signed payload with a UTC timestamp close to "now".
func TestSignSetsTimestamp(t *testing.T) {
	before := time.Now().UTC()
	env := signIss(t, "gobl:a.example", "gobl:b.example")
	after := time.Now().UTC()

	p, err := head.SignedPayload(env.Signatures[0])
	require.NoError(t, err)
	require.NotNil(t, p.Ts, "signing must stamp ts automatically")
	// cal.Timestamp truncates to millisecond precision, so widen the
	// window by 1ms on each side to avoid spurious failures.
	assert.True(t,
		!p.Ts.Time.Before(before.Add(-time.Millisecond)) &&
			!p.Ts.Time.After(after.Add(time.Millisecond)),
		"ts %v should fall within [%v, %v]", p.Ts.Time, before, after,
	)
}

func TestSign(t *testing.T) {
	type tt struct {
		opts *SignOptions
		err  string
	}

	tests := testy.NewTable()
	tests.Add("success", func(t *testing.T) interface{} {
		return tt{
			opts: &SignOptions{
				ParseOptions: &ParseOptions{
					Input: testFileReader(t, "testdata/nototals.json"),
				},
			},
		}
	})
	tests.Add("with signature", func(t *testing.T) interface{} {
		return tt{
			opts: &SignOptions{
				ParseOptions: &ParseOptions{
					Input: testFileReader(t, "testdata/signed.json"),
				},
			},
		}
	})
	tests.Add("draft envelope", func(t *testing.T) interface{} {
		return tt{
			opts: &SignOptions{
				ParseOptions: &ParseOptions{
					Input: testFileReader(t, "testdata/draft.json"),
				},
			},
		}
	})
	tests.Add("invalid document", func(_ *testing.T) interface{} {
		return tt{
			opts: &SignOptions{
				ParseOptions: &ParseOptions{
					SetFile: map[string]string{
						"foo": "testdata/invalid.yaml",
					},
				},
			},
			err: "internal: yaml: line 2: found unexpected end of stream",
		}
	})

	// TODO: Add test that asserts errors returned by `env.Calculate` are
	// returned.

	tests.Run(t, func(t *testing.T, tt tt) {
		t.Parallel()

		opts := tt.opts
		if opts.PrivateKey == nil {
			opts.PrivateKey = privateKey
		}

		got, err := Sign(context.Background(), opts)

		if tt.err == "" {
			require.Nil(t, err)
		} else {
			require.EqualError(t, err, tt.err)
		}

		re := testy.Replacement{
			Regexp:      regexp.MustCompile(`(?s)"sigs": \[.+\]`),
			Replacement: `"sigs": ["signature data"]`,
		}
		if d := testy.DiffAsJSON(testy.Snapshot(t), got, re); d != nil {
			t.Error(d)
		}
	})
}
