package it_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/it"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeTaxIdentity(t *testing.T) {
	tests := []struct {
		code     cbc.Code
		expected cbc.Code
		typ      cbc.Key
	}{
		{
			code:     "12345678901",
			expected: "12345678901",
			typ:      it.TaxIdentityTypeBusiness,
		},
		{
			code:     "123-456-789-01",
			expected: "12345678901",
			typ:      it.TaxIdentityTypeBusiness,
		},
		{
			code:     "123456 789 01",
			expected: "12345678901",
			typ:      it.TaxIdentityTypeBusiness,
		},
		{
			code:     "IT 12345678901",
			expected: "12345678901",
			typ:      it.TaxIdentityTypeBusiness,
		},
		{
			code:     "RSSMRA74D22A001Q",
			expected: "RSSMRA74D22A001Q",
			typ:      it.TaxIdentityTypeIndividual,
		},
		{
			code:     " RSS-MRA 74D22 A00 1Q ",
			expected: "RSSMRA74D22A001Q",
			typ:      it.TaxIdentityTypeIndividual,
		},
	}
	for _, ts := range tests {
		tID := &tax.Identity{Country: l10n.IT, Code: ts.code}
		err := it.Calculate(tID)
		assert.NoError(t, err)
		assert.Equal(t, ts.expected, tID.Code)
		assert.Equal(t, ts.typ, tID.Type)
	}
}

func TestValidateBusinessTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		zone l10n.Code
		err  string
	}{
		{name: "good 1", code: "12345678903"},
		{name: "good 2", code: "13029381004"},
		{name: "good 3", code: "10182640150"},
		{
			name: "empty",
			code: "",
			err:  "",
		},
		{
			name: "too long",
			code: "123456789001",
			err:  "invalid length",
		},
		{
			name: "too short",
			code: "1234567890",
			err:  "invalid length",
		},
		{
			name: "not normalized",
			code: "12.449.965-439",
			err:  "contains invalid characters",
		},
		{
			name: "includes non-numeric characters",
			code: "A764352056Z",
			err:  "contains invalid characters",
		},
		{
			name: "invalid check digit",
			code: "12345678901",
			err:  "invalid check digit",
		},
		{
			name: "invalid check digit",
			code: "13029381009",
			err:  "invalid check digit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: l10n.IT, Type: it.TaxIdentityTypeBusiness, Code: tt.code}
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

func TestValidateIndividualTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		zone l10n.Code
		err  string
	}{
		{name: "good 1", code: "RSSGNN60R30H501U"}, // Technical code specs
		{name: "good 2", code: "RSSMRA74D22A001Q"}, // https://www.studiolegalemetta.com/en/italian-tax-code-codice-fiscale/
		{name: "good 3", code: "FOOBRR80C04H146T"}, // Generated at https://www.codicefiscale.com/calcolo-completato.php
		{name: "good 4", code: "LWNSML81L16F205A"}, // ..
		{
			name: "empty",
			code: "",
			err:  "",
		},
		{
			name: "too long",
			code: "RSSGNN60R30H501U1",
			err:  "invalid format",
		},
		{
			name: "too short",
			code: "RSSGNN60R30H501",
			err:  "invalid format",
		},
		{
			name: "not normalized",
			code: "RSS GNN60R30 H501U",
			err:  "must be in a valid format",
		},
		{
			name: "incorrect format",
			code: "AYSGNN60R30H50UU",
			err:  "invalid format",
		},
		{
			name: "invalid check digit",
			code: "RSXGNN60R30H501U",
			err:  "invalid check digit",
		},
		{
			name: "invalid check digit 2",
			code: "RSSGNN60R30H502U",
			err:  "invalid check digit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: l10n.IT, Type: it.TaxIdentityTypeIndividual, Code: tt.code}
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

func TestTaxIdentityValidateGeneralCases(t *testing.T) {
	tests := []struct {
		name string
		tID  *tax.Identity
		err  string
	}{
		{
			name: "just country",
			tID:  &tax.Identity{Country: l10n.IT},
			err:  "",
		},
		{
			name: "no type, assume biz",
			tID:  &tax.Identity{Country: l10n.IT, Code: "12345678903"},
			err:  "",
		},
		{
			name: "biz code",
			tID:  &tax.Identity{Country: l10n.IT, Type: it.TaxIdentityTypeBusiness, Code: "12345678903"},
			err:  "",
		},
		{
			name: "admin code",
			tID:  &tax.Identity{Country: l10n.IT, Type: it.TaxIdentityTypeGovernment, Code: "12345678903"},
			err:  "",
		},
		{
			name: "no type, with individual",
			tID:  &tax.Identity{Country: l10n.IT, Code: "RSSGNN60R30H501U"},
			err:  "invalid characters",
		},
		{
			name: "with type for individual",
			tID:  &tax.Identity{Country: l10n.IT, Type: it.TaxIdentityTypeIndividual, Code: "RSSGNN60R30H501U"},
			err:  "",
		},
		{
			name: "with type for individual with typo",
			tID:  &tax.Identity{Country: l10n.IT, Type: it.TaxIdentityTypeIndividual, Code: "RSSGNN60R30H501Z"},
			err:  "invalid check digit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := it.Validate(tt.tID)
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
