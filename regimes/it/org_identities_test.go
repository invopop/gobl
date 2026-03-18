package it_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/it"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestIdentityNormalization(t *testing.T) {
	r := tax.RegimeDefFor("IT")

	t.Run("normalize codice fiscale", func(t *testing.T) {
		p1 := &org.Identity{
			Key:  it.IdentityKeyFiscalCode,
			Code: "RSS.mra-74D22-A00 . 1Q",
		}
		r.NormalizeObject(p1)
		assert.Equal(t, "RSSMRA74D22A001Q", p1.Code.String())
	})
}

func TestIdentityValidation(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "good 1", code: "RSSGNN60R30H501U"}, // Technical code specs
		{name: "good 2", code: "RSSMRA74D22A001Q"}, // https://www.studiolegalemetta.com/en/italian-tax-code-codice-fiscale/
		{name: "good 3", code: "FOOBRR80C04H146T"}, // Generated at https://www.codicefiscale.com/calcolo-completata.php
		{name: "good 4", code: "LWNSML81L16F205A"}, // ..
		{
			name: "good company 1",
			code: "12345678903",
		},
		{
			name: "good company 2",
			code: "10182640150",
		},
		{
			name: "bad company",
			code: "12345678901",
			err:  "IDENTITY-03",
		},
		{
			name: "empty",
			code: "",
			err:  "IDENTITY-01",
		},
		{
			name: "too long",
			code: "RSSGNN60R30H501U1",
			err:  "IDENTITY-02",
		},
		{
			name: "too short",
			code: "RSSGNN60R30H501",
			err:  "IDENTITY-02",
		},
		{
			name: "not normalized",
			code: "RSS GNN60R30 H501U",
			err:  "IDENTITY-02",
		},
		{
			name: "incorrect format",
			code: "AYSGNN60R30H50UU",
			err:  "IDENTITY-02",
		},
		{
			name: "invalid check digit",
			code: "RSXGNN60R30H501U",
			err:  "IDENTITY-03",
		},
		{
			name: "invalid check digit 2",
			code: "RSSGNN60R30H502U",
			err:  "IDENTITY-03",
		},
	}

	opts := []rules.WithContext{
		tax.RegimeContext(it.CountryCode),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &org.Identity{Key: it.IdentityKeyFiscalCode, Code: tt.code}
			err := rules.Validate(tID, opts...)
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
