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
			name: "empty",
			code: cbc.Code(""),
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
			name:    "too long",
			code:    cbc.Code("12345678901234567890ABCDE"),
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
