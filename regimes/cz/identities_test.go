package cz_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/cz"
	"github.com/stretchr/testify/assert"
)

func TestValidateIdentity(t *testing.T) {
	tests := []struct {
		name string
		key  cbc.Key
		code cbc.Code
		err  string
	}{
		{name: "valid ICO - Skoda Auto", key: cz.IdentityKeyICO, code: "00177041"},
		{name: "valid ICO - CEZ", key: cz.IdentityKeyICO, code: "45274649"},
		{name: "valid ICO - Komercni banka", key: cz.IdentityKeyICO, code: "45317054"},
		{
			name: "invalid ICO checksum",
			key:  cz.IdentityKeyICO,
			code: "00177042",
			err:  "checksum mismatch",
		},
		{
			name: "invalid ICO too short",
			key:  cz.IdentityKeyICO,
			code: "0017704",
			err:  "invalid format",
		},
		{
			name: "non-ICO identity ignored",
			key:  "other",
			code: "anything",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &org.Identity{Key: tt.key, Code: tt.code}
			err := cz.Validate(id)
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
