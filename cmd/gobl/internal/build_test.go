package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/flimzy/testy"
)

func Test_parseSets(t *testing.T) {
	tests := []struct {
		name string
		opts BuildOptions
		err  string
	}{
		{
			name: "invalid yaml",
			opts: BuildOptions{
				SetYAML: map[string]string{
					"foo": "[bar",
				},
			},
			err: `yaml: line 1: did not find expected ',' or ']'`,
		},
		{
			name: "valid yaml",
			opts: BuildOptions{
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
			opts: BuildOptions{
				SetYAML: map[string]string{
					".": `{"foo":"bar"}`,
				},
			},
		},
		{
			name: "literal period",
			opts: BuildOptions{
				SetYAML: map[string]string{
					"\\.": `foo`,
				},
			},
		},
		{
			name: "period",
			opts: BuildOptions{
				SetYAML: map[string]string{
					"foo.bar": "baz",
				},
			},
		},
		{
			name: "anchored at root",
			opts: BuildOptions{
				SetYAML: map[string]string{
					".foo": "bar",
				},
			},
		},
		{
			name: "unmergable",
			opts: BuildOptions{
				SetYAML: map[string]string{
					".": "foo",
				},
			},
			err: "src and dst must be of same type",
		},
		{
			name: "explicit string",
			opts: BuildOptions{
				SetString: map[string]string{
					"foo": "1234",
				},
			},
		},
		{
			name: "root string",
			opts: BuildOptions{
				SetString: map[string]string{
					".": "1234",
				},
			},
			err: "src and dst must be of same type",
		},
		{
			name: "missing file",
			opts: BuildOptions{
				SetFile: map[string]string{
					"foo": "notfound.yaml",
				},
			},
			err: `open notfound.yaml: no such file or directory`,
		},
		{
			name: "invalid file",
			opts: BuildOptions{
				SetFile: map[string]string{
					"foo": "testdata/invalid.yaml",
				},
			},
			err: `yaml: line 2: found unexpected end of stream`,
		},
		{
			name: "unmergable",
			opts: BuildOptions{
				SetFile: map[string]string{
					".": "testdata/unmergable.yaml",
				},
			},
			err: `src and dst must be of same type`,
		},
		{
			name: "valid file",
			opts: BuildOptions{
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
			got, err := parseSets(tt.opts)
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
