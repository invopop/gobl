package br_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/br"
	"github.com/invopop/gobl/rules"
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
			err: "[GOBL-BR-ORG-PARTY-03] ($.addresses[0].code) Brazilian postal code must match the valid format",
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
			err: "[GOBL-BR-ORG-PARTY-02] ($.addresses[0].state) Brazilian state must be one of the valid states",
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
			err: "[GOBL-BR-ORG-PARTY-03] ($.addresses[0].code) Brazilian postal code must match the valid format",
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
		{
			name: "with municipality code",
			party: &org.Party{
				TaxID: &tax.Identity{
					Country: "BR",
				},
				Ext: tax.Extensions{
					br.ExtKeyMunicipality: "3550308", // São Paulo city code
				},
			},
		},
		{
			name: "with invalid municipality code",
			party: &org.Party{
				TaxID: &tax.Identity{
					Country: "BR",
				},
				Ext: tax.Extensions{
					br.ExtKeyMunicipality: "00", // Invalid city code
				},
			},
			err: "[GOBL-BR-ORG-PARTY-01] ($.ext) Brazilian party ext must define a valid 'br-ibge-municipality' code",
		},
	}

	br := tax.RegimeContext(br.CountryCode)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := rules.Validate(tt.party, br)
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
			err:  "[GOBL-BR-ORG-PARTY-03] ($.addresses[0].code) Brazilian postal code must match the valid format",
		},
		{
			name: "too long",
			code: "123456789",
			err:  "GOBL-BR-ORG-PARTY-03",
		},
		{
			name: "invalid chars",
			code: "12345-678a",
			err:  "GOBL-BR-ORG-PARTY-03",
		},
		{
			name: "dash in wrong place",
			code: "1234-5678",
			err:  "GOBL-BR-ORG-PARTY-03",
		},
	}
	br := tax.RegimeContext(br.CountryCode)
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
			err := rules.Validate(party, br)
			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.err)
			}
		})
	}
}
