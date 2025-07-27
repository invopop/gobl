package cbc_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCodeIn(t *testing.T) {
	c := cbc.Code("FOO")

	assert.True(t, c.In("BAR", "FOO", "DOM"))
	assert.False(t, c.In("BAR", "DOM"))
}

func TestCodeEmpty(t *testing.T) {
	assert.Equal(t, cbc.Code(""), cbc.CodeEmpty)
	assert.True(t, cbc.Code("").IsEmpty())
}

func TestCodeJoin(t *testing.T) {
	t.Run("basic join", func(t *testing.T) {
		c := cbc.Code("BAR")
		assert.Equal(t, "BAR-FOO", c.Join("FOO").String())
	})
	t.Run("empty base join", func(t *testing.T) {
		c := cbc.Code("")
		assert.Equal(t, "FOO", c.Join("FOO").String())
	})
	t.Run("empty postfix join", func(t *testing.T) {
		c := cbc.Code("BAR")
		assert.Equal(t, "BAR", c.Join("").String())
	})
	t.Run("custom separator", func(t *testing.T) {
		c := cbc.Code("BAR")
		assert.Equal(t, "BAR|FOO", c.JoinWith("|", "FOO").String())
	})
}

func TestNormalizeCode(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		want cbc.Code
	}{
		{
			name: "uppercase",
			code: cbc.Code("FOO"),
			want: cbc.Code("FOO"),
		},
		{
			name: "lowercase",
			code: cbc.Code("foo"),
			want: cbc.Code("foo"),
		},
		{
			name: "mixed case",
			code: cbc.Code("Foo"),
			want: cbc.Code("Foo"),
		},
		{
			name: "with spaces",
			code: cbc.Code("FOO BAR"),
			want: cbc.Code("FOO BAR"),
		},
		{
			name: "empty",
			code: cbc.Code(""),
			want: cbc.Code(""),
		},
		{
			name: "underscore",
			code: cbc.Code("FOO_BAR"),
			want: cbc.Code("FOO_BAR"),
		},
		{
			name: "whitespace",
			code: cbc.Code(" foo-bar1  "),
			want: cbc.Code("foo-bar1"),
		},
		{
			name: "invalid chars",
			code: cbc.Code("f$oo-bar1!"),
			want: cbc.Code("foo-bar1"),
		},
		{
			name: "multiple spaces",
			code: cbc.Code("foo bar dome"),
			want: cbc.Code("foo bar dome"),
		},
		{
			name: "multiple symbols 1",
			code: cbc.Code("foo- bar-$dome"),
			want: cbc.Code("foo-bar-dome"),
		},
		{
			name: "multiple symbols 2",
			code: cbc.Code("FOO  BAR--DOME"),
			want: cbc.Code("FOO BAR-DOME"),
		},
		{
			name: "colons",
			code: cbc.Code("0088:1234567891234"), // peppol example
			want: cbc.Code("0088:1234567891234"),
		},
		{
			name: "commas",
			code: cbc.Code("FL-C 64-3,5"),
			want: cbc.Code("FL-C 64-3,5"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, cbc.NormalizeCode(tt.code))
		})
	}
}

func TestNormalizeAlphanumericalCode(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		want cbc.Code
	}{
		{
			name: "uppercase",
			code: cbc.Code("FOO"),
			want: cbc.Code("FOO"),
		},
		{
			name: "lowercase",
			code: cbc.Code("foo"),
			want: cbc.Code("FOO"),
		},
		{
			name: "mixed case",
			code: cbc.Code("Foo"),
			want: cbc.Code("FOO"),
		},
		{
			name: "with spaces",
			code: cbc.Code("FOO BAR"),
			want: cbc.Code("FOOBAR"),
		},
		{
			name: "empty",
			code: cbc.Code(""),
			want: cbc.Code(""),
		},
		{
			name: "underscore",
			code: cbc.Code("FOO_BAR"),
			want: cbc.Code("FOOBAR"),
		},
		{
			name: "whitespace",
			code: cbc.Code(" foo-bar1  "),
			want: cbc.Code("FOOBAR1"),
		},
		{
			name: "invalid chars",
			code: cbc.Code("f$oo-bar1!"),
			want: cbc.Code("FOOBAR1"),
		},
		{
			name: "multiple spaces",
			code: cbc.Code("foo bar dome"),
			want: cbc.Code("FOOBARDOME"),
		},
		{
			name: "multiple symbols 1",
			code: cbc.Code("foo- bar-$dome"),
			want: cbc.Code("FOOBARDOME"),
		},
		{
			name: "multiple symbols 2",
			code: cbc.Code("FOO  BAR--DOME"),
			want: cbc.Code("FOOBARDOME"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, cbc.NormalizeAlphanumericalCode(tt.code))
		})
	}
}

func TestNormalizeNumericalCode(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		want cbc.Code
	}{
		{
			name: "letters",
			code: cbc.Code("FOO"),
			want: cbc.Code(""),
		},
		{
			name: "numbers",
			code: cbc.Code("1234"),
			want: cbc.Code("1234"),
		},
		{
			name: "mixed case",
			code: cbc.Code("Foo1234"),
			want: cbc.Code("1234"),
		},
		{
			name: "with spaces",
			code: cbc.Code("12 34"),
			want: cbc.Code("1234"),
		},
		{
			name: "empty",
			code: cbc.Code(""),
			want: cbc.Code(""),
		},
		{
			name: "underscore",
			code: cbc.Code("12_34"),
			want: cbc.Code("1234"),
		},
		{
			name: "whitespace",
			code: cbc.Code(" 345  "),
			want: cbc.Code("345"),
		},
		{
			name: "invalid chars",
			code: cbc.Code("f$oo-bar1!"),
			want: cbc.Code("1"),
		},
		{
			name: "multiple spaces",
			code: cbc.Code("1 2 3 4"),
			want: cbc.Code("1234"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, cbc.NormalizeNumericalCode(tt.code))
		})
	}
}

func TestCode_Validate(t *testing.T) {
	tests := []struct {
		name    string
		code    cbc.Code
		wantErr string
	}{
		{
			name: "valid 1",
			code: cbc.Code("ABC123"),
		},
		{
			name: "valid 2",
			code: cbc.Code("12345678901234567890ABCD"),
		},
		{
			name: "valid with lower",
			code: cbc.Code("ABC abc/123"),
		},
		{
			name: "valid with dot",
			code: cbc.Code("B3.12"),
		},
		{
			name: "valid with dash",
			code: cbc.Code("B3-12"),
		},
		{
			name: "valid with multiple dots",
			code: cbc.Code("B3.1.2"),
		},
		{
			name: "valid with multiple dashes",
			code: cbc.Code("B3-1-2"),
		},
		{
			name: "valid with slash",
			code: cbc.Code("B3/12"),
		},
		{
			name: "valid with space",
			code: cbc.Code("FR 12/BX"),
		},
		{
			name: "valid with colon",
			code: cbc.Code("FR:12/BX"),
		},
		{
			name: "empty",
			code: cbc.Code(""),
		},
		{
			name: "almost too long",
			code: cbc.Code("123456789012345678901234567890AB"),
		},
		{
			name:    "dot at start",
			code:    cbc.Code(".B123"),
			wantErr: "valid format",
		},
		{
			name:    "dot at end",
			code:    cbc.Code("B123."),
			wantErr: "valid format",
		},
		{
			name:    "dash at start",
			code:    cbc.Code("-B123"),
			wantErr: "valid format",
		},
		{
			name:    "dash at end",
			code:    cbc.Code("B123-"),
			wantErr: "valid format",
		},
		{
			name:    "multiple symbols",
			code:    cbc.Code("AB/.CD"),
			wantErr: "valid format",
		},
		{
			name:    "character return",
			code:    cbc.Code("AB\nCD"),
			wantErr: "valid format",
		},
		{
			name:    "character return",
			code:    cbc.Code("\n"),
			wantErr: "valid format",
		},
		{
			name:    "multi-dash",
			code:    cbc.Code("AB--CD"),
			wantErr: "valid format",
		},
		{
			name:    "too long",
			code:    cbc.Code("123456789012345678901234567890ABC123456789012345678901234567890ABC"),
			wantErr: "length must be between",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.code.Validate()
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), tt.wantErr)
				}
			}
		})
	}
}

func TestCode_IsValidLuhnChecksum(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		want bool
	}{
		// Valid Luhn checksums - pure numeric
		{
			name: "single digit zero",
			code: cbc.Code("0"),
			want: true,
		},
		{
			name: "two digit valid",
			code: cbc.Code("18"),
			want: true,
		},
		{
			name: "multiple valid cases",
			code: cbc.Code("26"),
			want: true,
		},
		{
			name: "valid case 34",
			code: cbc.Code("34"),
			want: true,
		},
		{
			name: "valid case 42",
			code: cbc.Code("42"),
			want: true,
		},
		{
			name: "valid case 59",
			code: cbc.Code("59"),
			want: true,
		},
		{
			name: "valid case 67",
			code: cbc.Code("67"),
			want: true,
		},
		{
			name: "valid case 75",
			code: cbc.Code("75"),
			want: true,
		},
		{
			name: "valid case 83",
			code: cbc.Code("83"),
			want: true,
		},
		{
			name: "valid case 91",
			code: cbc.Code("91"),
			want: true,
		},
		{
			name: "long valid number",
			code: cbc.Code("79927398713"),
			want: true,
		},
		{
			name: "valid credit card number",
			code: cbc.Code("4111111111111111"),
			want: true,
		},
		// Valid Luhn checksums with non-numeric characters (should be stripped)
		{
			name: "valid with letters",
			code: cbc.Code("ABC18DEF"),
			want: true,
		},
		{
			name: "valid with spaces",
			code: cbc.Code("1 8"),
			want: true,
		},
		{
			name: "valid with dashes",
			code: cbc.Code("2-6"),
			want: true,
		},
		{
			name: "valid with mixed separators",
			code: cbc.Code("4111-1111_1111\\1111"),
			want: true,
		},
		{
			name: "valid with dots",
			code: cbc.Code("7.9.9.2.7.3.9.8.7.1.3"),
			want: true,
		},
		{
			name: "valid with alphanumeric mix",
			code: cbc.Code("ID-4111111111111111-CHECK"),
			want: true,
		},
		{
			name: "valid with underscores",
			code: cbc.Code("CODE_18_END"),
			want: true,
		},
		{
			name: "valid with slashes",
			code: cbc.Code("REF/18/END"),
			want: true,
		},
		{
			name: "valid with colons",
			code: cbc.Code("PREFIX:34:SUFFIX"),
			want: true,
		},
		{
			name: "valid with commas",
			code: cbc.Code("A,B,C,42,X,Y,Z"),
			want: true,
		},
		// Invalid Luhn checksums - pure numeric
		{
			name: "single digit one",
			code: cbc.Code("1"),
			want: false,
		},
		{
			name: "invalid two digit",
			code: cbc.Code("19"),
			want: false,
		},
		{
			name: "invalid three digit",
			code: cbc.Code("123"),
			want: false,
		},
		{
			name: "invalid long number",
			code: cbc.Code("79927398714"),
			want: false,
		},
		{
			name: "invalid credit card number",
			code: cbc.Code("4111111111111112"),
			want: false,
		},
		// Invalid Luhn checksums with non-numeric characters
		{
			name: "invalid with letters",
			code: cbc.Code("ABC19DEF"),
			want: false,
		},
		{
			name: "invalid with spaces",
			code: cbc.Code("1 9"),
			want: false,
		},
		{
			name: "invalid with dashes",
			code: cbc.Code("1-2-3"),
			want: false,
		},
		{
			name: "invalid with mixed separators",
			code: cbc.Code("4111-1111-1111-1112"),
			want: false,
		},
		// Edge cases
		{
			name: "empty string",
			code: cbc.Code(""),
			want: false,
		},
		{
			name: "only letters",
			code: cbc.Code("ABCDEF"),
			want: false,
		},
		{
			name: "only separators",
			code: cbc.Code("-_.\\/"),
			want: false,
		},
		{
			name: "mixed separators no digits",
			code: cbc.Code("A-B.C_D/E:F,G H"),
			want: false,
		},
		{
			name: "single valid digit with noise",
			code: cbc.Code("XXX0YYY"),
			want: true,
		},
		{
			name: "single invalid digit with noise",
			code: cbc.Code("XXX1YYY"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tt.code.IsValidLuhnChecksum()
			assert.Equal(t, tt.want, result, "IsValidLuhnChecksum() for code %q", tt.code)
		})
	}
}

func TestCodeMap(t *testing.T) {
	cm := cbc.CodeMap{
		"foo": cbc.Code("01"),
		"bar": cbc.Code("02"),
	}
	t.Run("Has", func(t *testing.T) {
		assert.True(t, cm.Has("foo"))
		assert.True(t, cm.Has("foo", "bar"))
		assert.False(t, cm.Has("dom"))
		assert.False(t, cm.Has("foo", "dom"))
	})

	t.Run("validation", func(t *testing.T) {
		assert.NoError(t, cm.Validate())
		cm2 := cbc.CodeMap{
			"Invalid": cbc.Code("01"),
		}
		assert.ErrorContains(t, cm2.Validate(), "Invalid: must be in a valid format")
	})
}

func TestCodeMap_Equals(t *testing.T) {
	cm := cbc.CodeMap{
		"foo": cbc.Code("01"),
		"bar": cbc.Code("02"),
	}
	cm2 := cbc.CodeMap{
		"foo": cbc.Code("01"),
		"bar": cbc.Code("02"),
	}
	cm3 := cbc.CodeMap{
		"foo": cbc.Code("01"),
		"bar": cbc.Code("03"),
	}
	cm4 := cbc.CodeMap{
		"foo": cbc.Code("01"),
		"bar": cbc.Code("02"),
		"dom": cbc.Code("03"),
	}
	cm5 := cbc.CodeMap{
		"foo": cbc.Code("01"),
	}
	cm6 := cbc.CodeMap{
		"foo": cbc.Code("01"),
		"dom": cbc.Code("02"),
	}
	assert.True(t, cm.Equals(cm2))
	assert.False(t, cm.Equals(cm3))
	assert.False(t, cm.Equals(cm4))
	assert.False(t, cm.Equals(cm5))
	assert.False(t, cm.Equals(cm6))
}

func TestCodeMapHas(t *testing.T) {
	cm := cbc.CodeMap{
		"foo": cbc.Code("01"),
		"bar": cbc.Code("02"),
	}
	err := validation.Validate(cm, cbc.CodeMapHas("foo", "bar"))
	assert.NoError(t, err)
	assert.ErrorContains(t, validation.Validate(cm, cbc.CodeMapHas("foo", "dom")), "dom: required.")
	err = validation.Validate(nil, cbc.CodeMapHas("foo"))
	assert.NoError(t, err)
}

func TestCodeJSONSchema(t *testing.T) {
	s := cbc.Code("").JSONSchema()
	assert.Equal(t, "string", s.Type)
	assert.Equal(t, "Code", s.Title)
	assert.Equal(t, uint64(1), *s.MinLength)
	assert.Equal(t, uint64(64), *s.MaxLength)
}

func TestCodeMapJSONSchemaExtend(t *testing.T) {
	eg := `{
			"type": "object",
			"additionalProperties": {
				"$ref": "https://gobl.org/draft-0/cbc/code"
			},
			"description": "CodeMap is a map of keys to specific codes, useful to determine regime specific codes from their key counterparts."
	}`
	js := new(jsonschema.Schema)
	require.NoError(t, json.Unmarshal([]byte(eg), js))

	cm := cbc.CodeMap{}
	cm.JSONSchemaExtend(js)

	assert.Nil(t, js.AdditionalProperties)
	assert.Equal(t, 1, len(js.PatternProperties))
	assert.Equal(t, "https://gobl.org/draft-0/cbc/code", js.PatternProperties[cbc.KeyPattern].Ref)
}
