package cli

import (
	"context"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/flimzy/testy"
)

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
			err: "code=422, message=yaml: line 2: found unexpected end of stream",
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
