package jp_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/jp"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		// Known valid corporate numbers
		{name: "valid sony corporate number", code: "5010401067252"}, // https://www.houjin-bangou.nta.go.jp/henkorireki-johoto.html?selHouzinNo=1130001011420
		{name: "valid nintendo corporate number", code: "1130001011420"},
		// Invalid format
		{
			name: "too short",
			code: "123456789",
			err:  "invalid format",
		},
		{
			name: "too long",
			code: "12345678901234",
			err:  "invalid format",
		},
		{
			name: "leading zero",
			code: "0123456789012",
			err:  "invalid format",
		},
		{
			name: "contains letters",
			code: "A123456789012",
			err:  "invalid format",
		},
		// Invalid checksum
		{
			name: "bad checksum",
			code: "2010401067252",
			err:  "checksum mismatch",
		},
		{
			name: "bad checksum 2",
			code: "1130001011421",
			err:  "checksum mismatch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "JP", Code: tt.code}
			err := jp.Validate(tID)
			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), tt.err)
				}
			}
		})
	}
}
