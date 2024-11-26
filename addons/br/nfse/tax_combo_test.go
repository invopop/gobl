package nfse_test

import (
	"testing"

	"github.com/invopop/gobl/addons/br/nfse"
	"github.com/invopop/gobl/regimes/br"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestTaxComboValidation(t *testing.T) {
	addon := tax.AddonForKey(nfse.V1)

	tests := []struct {
		name string
		tc   *tax.Combo
		err  string
	}{
		{
			name: "valid ISS tax combo",
			tc: &tax.Combo{
				Category: br.TaxCategoryISS,
				Ext: tax.Extensions{
					nfse.ExtKeyISSLiability: "1",
				},
			},
		},
		{
			name: "valid non-ISS tax combo",
			tc: &tax.Combo{
				Category: br.TaxCategoryPIS,
			},
		},
		{
			name: "missing ISS liability",
			tc: &tax.Combo{
				Category: br.TaxCategoryISS,
			},
			err: "br-nfse-iss-liability: required",
		},
	}
	for _, ts := range tests {
		t.Run(ts.name, func(t *testing.T) {
			err := addon.Validator(ts.tc)
			if ts.err == "" {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.ErrorContains(t, err, ts.err)
				}
			}
		})
	}
}

func TestTaxComboNormalization(t *testing.T) {
	addon := tax.AddonForKey(nfse.V1)

	tests := []struct {
		name string
		tc   *tax.Combo
		out  tax.ExtValue
	}{
		{
			name: "no tax combo",
			tc:   nil,
		},
		{
			name: "sets default ISS liability",
			tc: &tax.Combo{
				Category: br.TaxCategoryISS,
			},
			out: "1",
		},
		{
			name: "does not override ISS liability",
			tc: &tax.Combo{
				Category: br.TaxCategoryISS,
				Ext: tax.Extensions{
					nfse.ExtKeyISSLiability: "2",
				},
			},
			out: "2",
		},
		{
			name: "non-ISS tax combo",
			tc: &tax.Combo{
				Category: br.TaxCategoryPIS,
			},
		},
	}
	for _, ts := range tests {
		t.Run(ts.name, func(t *testing.T) {
			addon.Normalizer(ts.tc)
			if ts.tc == nil {
				assert.Nil(t, ts.tc)
			} else {
				assert.NotNil(t, ts.tc)
				assert.Equal(t, ts.out, ts.tc.Ext[nfse.ExtKeyISSLiability])
			}
		})
	}

}
