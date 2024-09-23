package cbc_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/stretchr/testify/assert"
)

func TestCodeIn(t *testing.T) {
	c := cbc.Code("FOO")

	assert.True(t, c.In("BAR", "FOO", "DOM"))
	assert.False(t, c.In("BAR", "DOM"))
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
			want: cbc.Code("FOO-BAR"),
		},
		{
			name: "whitespace",
			code: cbc.Code(" foo-bar1  "),
			want: cbc.Code("FOO-BAR1"),
		},
		{
			name: "invalid chars",
			code: cbc.Code("f$oo-bar1!"),
			want: cbc.Code("FOO-BAR1"),
		},
		{
			name: "multiple spaces",
			code: cbc.Code("foo bar dome"),
			want: cbc.Code("FOO BAR DOME"),
		},
		{
			name: "multiple symbols 1",
			code: cbc.Code("foo- bar-$dome"),
			want: cbc.Code("FOO-BAR-DOME"),
		},
		{
			name: "multiple symbols 2",
			code: cbc.Code("foo  bar--dome"),
			want: cbc.Code("FOO BAR-DOME"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, cbc.NormalizeCode(tt.code))
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
			name:    "lower case",
			code:    cbc.Code("ab"),
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
			code:    cbc.Code("123456789012345678901234567890ABC"),
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
