package cli

import (
	"context"
	"encoding/json"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/flimzy/testy"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/note"
)

var (
	privateKey    = new(dsig.PrivateKey)
	publicKey     = new(dsig.PublicKey)
	verifyKeyText string
)

const signingKeyFile = "testdata/private.jwk"

func init() {
	data, err := os.ReadFile(signingKeyFile)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(data, privateKey); err != nil {
		panic(err)
	}
	publicKey = privateKey.Public()
	pub, err := json.Marshal(publicKey)
	if err != nil {
		panic(err)
	}
	verifyKeyText = string(pub)
}

func Test_parseSets(t *testing.T) {
	tests := []struct {
		name string
		opts *ParseOptions
		err  string
	}{
		{
			name: "invalid yaml",
			opts: &ParseOptions{
				SetYAML: map[string]string{
					"foo": "[bar",
				},
			},
			err: `code=422, message=yaml: line 1: did not find expected ',' or ']'`,
		},
		{
			name: "valid yaml",
			opts: &ParseOptions{
				SetYAML: map[string]string{
					"sring":  "bar",
					"number": "1234",
					"bool":   "true",
					"array":  "[1,2,3]",
					"object": `{"foo":"bar"}`,
				},
			},
		},
		{
			name: "root key",
			opts: &ParseOptions{
				SetYAML: map[string]string{
					".": `{"foo":"bar"}`,
				},
			},
		},
		{
			name: "literal period",
			opts: &ParseOptions{
				SetYAML: map[string]string{
					"\\.": `foo`,
				},
			},
		},
		{
			name: "period",
			opts: &ParseOptions{
				SetYAML: map[string]string{
					"foo.bar": "baz",
				},
			},
		},
		{
			name: "anchored at root",
			opts: &ParseOptions{
				SetYAML: map[string]string{
					".foo": "bar",
				},
			},
		},
		{
			name: "unmergable",
			opts: &ParseOptions{
				SetYAML: map[string]string{
					".": "foo",
				},
			},
			err: "code=422, message=src and dst must be of same type",
		},
		{
			name: "explicit string",
			opts: &ParseOptions{
				SetString: map[string]string{
					"foo": "1234",
				},
			},
		},
		{
			name: "root string",
			opts: &ParseOptions{
				SetString: map[string]string{
					".": "1234",
				},
			},
			err: "src and dst must be of same type",
		},
		{
			name: "missing file",
			opts: &ParseOptions{
				SetFile: map[string]string{
					"foo": "notfound.yaml",
				},
			},
			err: `open notfound.yaml: no such file or directory`,
		},
		{
			name: "invalid file",
			opts: &ParseOptions{
				SetFile: map[string]string{
					"foo": "testdata/invalid.yaml",
				},
			},
			err: `yaml: line 2: found unexpected end of stream`,
		},
		{
			name: "unmergable",
			opts: &ParseOptions{
				SetFile: map[string]string{
					".": "testdata/unmergable.yaml",
				},
			},
			err: `src and dst must be of same type`,
		},
		{
			name: "valid file",
			opts: &ParseOptions{
				SetFile: map[string]string{
					"foo": "testdata/valid.yaml",
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			opts := tt.opts
			got, err := parseSets(opts)
			if tt.err == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, tt.err)
			}
			if err != nil {
				return
			}
			if d := testy.DiffInterface(testy.Snapshot(t), got); d != nil {
				t.Error(d)
			}
		})
	}
}

func TestBuild(t *testing.T) {
	type tt struct {
		opts *BuildOptions
		err  string
	}

	tests := testy.NewTable()
	tests.Add("success", func(t *testing.T) interface{} {
		return tt{
			opts: &BuildOptions{
				ParseOptions: &ParseOptions{
					Input: testFileReader(t, "testdata/nototals.json"),
				},
			},
		}
	})
	tests.Add("merge YAML", func(t *testing.T) interface{} {
		return tt{
			opts: &BuildOptions{
				ParseOptions: &ParseOptions{
					Input: testFileReader(t, "testdata/nototals.json"),
					SetYAML: map[string]string{
						"doc.supplier.name": "Other Company",
					},
				},
			},
		}
	})
	tests.Add("invalid type", tt{
		opts: &BuildOptions{
			ParseOptions: &ParseOptions{
				Input: strings.NewReader(`{
				"$schema": "https://gobl.org/draft-0/envelope",
				"head": {
					"uuid": "9d8eafd5-77be-11ec-b485-5405db9a3e49",
					"dig": {
						"alg": "sha256",
						"val": "dce3bc3c8bf28f3d209f783917b3082ddc0339a66e9ba3aa63849e4357db1422"
					}
				},
				doc: {
					"$schema": "https://example.com/duck",
					"walk": "like a duck",
					"talk": "like a duck",
					"look": "like a duck"
				}
			}`),
			},
		},
		err: `code=400, message=unknown-schema`,
	})
	tests.Add("with template", func(t *testing.T) interface{} {
		return tt{
			opts: &BuildOptions{
				ParseOptions: &ParseOptions{
					Template: strings.NewReader(`{"doc":{"supplier":{"name": "Other Company"}}}`),
					Input:    testFileReader(t, "testdata/noname.json"),
				},
			},
		}
	})
	tests.Add("template with empty input", func(t *testing.T) interface{} {
		return tt{
			opts: &BuildOptions{
				ParseOptions: &ParseOptions{
					Template: testFileReader(t, "testdata/nosig.json"),
					Input:    strings.NewReader("{}"),
				},
			},
		}
	})
	tests.Add("with signature", func(t *testing.T) interface{} {
		return tt{
			opts: &BuildOptions{
				ParseOptions: &ParseOptions{
					Template: testFileReader(t, "testdata/signed.json"),
					Input:    strings.NewReader("{}"),
				},
			},
		}
	})
	tests.Add("explicit type", func(t *testing.T) interface{} {
		return tt{
			opts: &BuildOptions{
				ParseOptions: &ParseOptions{
					Input:   testFileReader(t, "testdata/notype.json"),
					DocType: "bill.Invoice",
				},
			},
		}
	})
	tests.Add("draft", func(t *testing.T) interface{} {
		return tt{
			opts: &BuildOptions{
				ParseOptions: &ParseOptions{
					Input: testFileReader(t, "testdata/draft.json"),
				},
			},
		}
	})
	tests.Add("do not envelop", func(t *testing.T) interface{} {
		return tt{
			opts: &BuildOptions{
				ParseOptions: &ParseOptions{
					Input: testFileReader(t, "testdata/invoice.json"),
				},
			},
		}
	})
	tests.Add("envelop", func(t *testing.T) interface{} {
		return tt{
			opts: &BuildOptions{
				ParseOptions: &ParseOptions{
					Input:   testFileReader(t, "testdata/invoice.json"),
					Envelop: true,
				},
			},
		}
	})

	tests.Run(t, func(t *testing.T, tt tt) {
		t.Parallel()
		opts := tt.opts
		got, err := Build(context.Background(), opts)
		if tt.err == "" {
			assert.Nil(t, err)
		} else {
			assert.EqualError(t, err, tt.err)
		}
		if err != nil {
			return
		}
		replacements := []testy.Replacement{
			{
				Regexp:      regexp.MustCompile(`(?s)"sigs": \[.*\]`),
				Replacement: `"sigs": ["signature data"]`,
			},
			{
				Regexp:      regexp.MustCompile(`"uuid":.?"[^\"]+"`),
				Replacement: `"uuid":"00000000-0000-0000-0000-000000000000"`,
			},
		}
		if d := testy.DiffAsJSON(testy.Snapshot(t), got, replacements...); d != nil {
			t.Error(d)
		}
	})
}

func TestBuildWithPartialEnvelope(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		opts := &BuildOptions{
			ParseOptions: &ParseOptions{
				Input: testFileReader(t, "testdata/message.env.yaml"),
			},
		}
		got, err := Build(context.Background(), opts)
		require.NoError(t, err)
		env, ok := got.(*gobl.Envelope)
		require.True(t, ok)

		assert.NotEmpty(t, env.Head.UUID.String())
		assert.Empty(t, env.Signatures)

		msg, ok := env.Extract().(*note.Message)
		require.True(t, ok)
		assert.Equal(t, "https://gobl.org/draft-0/note/message", env.Document.Schema.String())
		assert.Equal(t, "Test Message", msg.Title)
		assert.Equal(t, "We hope you like this test message!", msg.Content)
	})
}
