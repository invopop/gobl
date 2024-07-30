package it_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/it"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIdentityNormalization(t *testing.T) {
	r := tax.RegimeFor("IT")

	t.Run("normalize codice fiscale", func(t *testing.T) {
		p1 := &org.Identity{
			Key:  it.IdentityKeyFiscalCode,
			Code: "RSS.mra-74D22-A00 . 1Q",
		}
		require.NoError(t, r.CalculateObject(p1))
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
		{name: "good 3", code: "FOOBRR80C04H146T"}, // Generated at https://www.codicefiscale.com/calcolo-completato.php
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
			err:  "invalid check digit",
		},
		{
			name: "empty",
			code: "",
			err:  "code: cannot be blank",
		},
		{
			name: "too long",
			code: "RSSGNN60R30H501U1",
			err:  "code: invalid format",
		},
		{
			name: "too short",
			code: "RSSGNN60R30H501",
			err:  "code: invalid format",
		},
		{
			name: "not normalized",
			code: "RSS GNN60R30 H501U",
			err:  "code: invalid format",
		},
		{
			name: "incorrect format",
			code: "AYSGNN60R30H50UU",
			err:  "code: invalid format",
		},
		{
			name: "invalid check digit",
			code: "RSXGNN60R30H501U",
			err:  "code: invalid check digit",
		},
		{
			name: "invalid check digit 2",
			code: "RSSGNN60R30H502U",
			err:  "invalid check digit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &org.Identity{Key: it.IdentityKeyFiscalCode, Code: tt.code}
			err := it.Validate(tID)
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
