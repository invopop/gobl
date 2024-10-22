package pl_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/pl"
	"github.com/stretchr/testify/assert"
)

func TestValidateTaxNumber(t *testing.T) {
	tests := []struct {
		name     string
		identity *org.Identity
		wantErr  bool
	}{
		{
			name: "Valid PESEL",
			identity: &org.Identity{
				Key:  pl.IdentityKeyTaxNumber,
				Code: cbc.Code("44051401359"), // Replace with an actual valid PESEL number
			},
			wantErr: false,
		},
		{
			name: "Invalid length PESEL",
			identity: &org.Identity{
				Key:  pl.IdentityKeyTaxNumber,
				Code: cbc.Code("1234567890"), // Invalid PESEL with less than 11 digits
			},
			wantErr: true,
		},
		{
			name: "Invalid checksum PESEL",
			identity: &org.Identity{
				Key:  pl.IdentityKeyTaxNumber,
				Code: cbc.Code("44051401358"), // Incorrect checksum
			},
			wantErr: true,
		},
		{
			name: "Empty PESEL code",
			identity: &org.Identity{
				Key:  pl.IdentityKeyTaxNumber,
				Code: cbc.Code(""),
			},
			wantErr: false,
		},
		{
			name: "Wrong Key Identity",
			identity: &org.Identity{
				Key:  cbc.Key("wrong-key"),
				Code: cbc.Code("44051401359"),
			},
			wantErr: false,
		},
		{
			name:     "Nil Identity",
			identity: nil,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pl.Validate(tt.identity)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
