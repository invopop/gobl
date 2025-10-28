package br_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/br"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestValidateAddresses(t *testing.T) {
	tests := []struct {
		name  string
		party *org.Party
		err   string
	}{
		{
			name:  "nil party",
			party: nil,
		},
		{
			name:  "empty party",
			party: &org.Party{},
		},
		{
			name: "empty address",
			party: &org.Party{
				Addresses: []*org.Address{},
			},
		},
		{
			name: "empty Brazilian address",
			party: &org.Party{
				Addresses: []*org.Address{
					{
						Country: "BR",
					},
				},
			},
		},
		{
			name: "valid Brazilian address",
			party: &org.Party{
				Addresses: []*org.Address{
					{
						Country: "BR",
						Code:    "12345-678",
						State:   "SP",
					},
				},
			},
		},
		{
			name: "invalid Brazilian post code",
			party: &org.Party{
				Addresses: []*org.Address{
					{
						Country: "BR",
						Code:    "12345",
					},
				},
			},
			err: "code: must be in a valid format",
		},
		{
			name: "invalid Brazilian state",
			party: &org.Party{
				Addresses: []*org.Address{
					{
						Country: "BR",
						State:   "XX",
					},
				},
			},
			err: "state: must be a valid value.",
		},
		{
			name: "invalid Brazilian address with tax country only",
			party: &org.Party{
				TaxID: &tax.Identity{
					Country: "BR",
				},
				Addresses: []*org.Address{
					{
						Code: "12345",
					},
				},
			},
			err: "code: must be in a valid format",
		},
		{
			name: "non-Brazilian address",
			party: &org.Party{
				TaxID: &tax.Identity{
					Country: "BR",
				},
				Addresses: []*org.Address{
					{
						Country: "US",
						Code:    "123",
						State:   "NY",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := br.Validate(tt.party)
			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.err)
			}
		})
	}
}

func TestValidatePostCodes(t *testing.T) {
	tests := []struct {
		name string
		code string
		err  string
	}{
		{
			name: "valid",
			code: "12345-678",
		},
		{
			name: "valid without dash",
			code: "12345678",
		},
		{
			name: "too short",
			code: "12345",
			err:  "code: must be in a valid format",
		},
		{
			name: "too long",
			code: "123456789",
			err:  "code: must be in a valid format",
		},
		{
			name: "invalid chars",
			code: "12345-678a",
			err:  "code: must be in a valid format",
		},
		{
			name: "dash in wrong place",
			code: "1234-5678",
			err:  "code: must be in a valid format",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			party := &org.Party{
				Addresses: []*org.Address{
					{
						Country: "BR",
						Code:    cbc.Code(tt.code),
					},
				},
			}
			err := br.Validate(party)
			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.err)
			}
		})
	}
}
