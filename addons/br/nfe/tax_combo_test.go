package nfe_test

import (
	"testing"

	"github.com/invopop/gobl/addons/br/nfe"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/br"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestTaxComboValidation(t *testing.T) {
	tests := []struct {
		name string
		tc   *tax.Combo
		err  string
	}{
		{
			name: "valid ICMS with CST",
			tc: &tax.Combo{
				Category: br.TaxCategoryICMS,
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					nfe.ExtKeyICMSCST:    "00",
					nfe.ExtKeyICMSOrigin: "0",
				}),
			},
		},
		{
			name: "valid ICMS with CSOSN",
			tc: &tax.Combo{
				Category: br.TaxCategoryICMS,
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					nfe.ExtKeyICMSCSOSN:  "102",
					nfe.ExtKeyICMSOrigin: "0",
				}),
			},
		},
		{
			name: "ICMS missing situation code",
			tc: &tax.Combo{
				Category: br.TaxCategoryICMS,
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					nfe.ExtKeyICMSOrigin: "0",
				}),
			},
			err: "ICMS tax combo requires 'br-nfe-icms-cst' or 'br-nfe-icms-csosn' extension",
		},
		{
			name: "ICMS missing origin",
			tc: &tax.Combo{
				Category: br.TaxCategoryICMS,
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					nfe.ExtKeyICMSCST: "00",
				}),
			},
			err: "ICMS tax combo requires 'br-nfe-icms-origin' extension",
		},
		{
			name: "valid PIS",
			tc: &tax.Combo{
				Category: br.TaxCategoryPIS,
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					nfe.ExtKeyPISCST: "01",
				}),
			},
		},
		{
			name: "PIS missing CST",
			tc: &tax.Combo{
				Category: br.TaxCategoryPIS,
			},
			err: "PIS tax combo requires 'br-nfe-pis-cst' extension",
		},
		{
			name: "valid COFINS",
			tc: &tax.Combo{
				Category: br.TaxCategoryCOFINS,
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					nfe.ExtKeyCOFINSCST: "01",
				}),
			},
		},
		{
			name: "COFINS missing CST",
			tc: &tax.Combo{
				Category: br.TaxCategoryCOFINS,
			},
			err: "COFINS tax combo requires 'br-nfe-cofins-cst' extension",
		},
		{
			name: "unrelated category is not constrained",
			tc: &tax.Combo{
				Category: br.TaxCategoryIPI,
			},
		},
	}
	for _, ts := range tests {
		t.Run(ts.name, func(t *testing.T) {
			err := rules.Validate(ts.tc, withAddonContext())
			if ts.err == "" {
				assert.NoError(t, err)
			} else if assert.Error(t, err) {
				assert.ErrorContains(t, err, ts.err)
			}
		})
	}
}

func TestTaxComboNormalization(t *testing.T) {
	tests := []struct {
		name   string
		tc     *tax.Combo
		expect map[cbc.Key]cbc.Code // expected codes; "" means the key must be absent
	}{
		{
			name: "ICMS sets default CST and origin",
			tc:   &tax.Combo{Category: br.TaxCategoryICMS},
			expect: map[cbc.Key]cbc.Code{
				nfe.ExtKeyICMSCST:    "00",
				nfe.ExtKeyICMSOrigin: "0",
			},
		},
		{
			name: "ICMS keeps CSOSN and does not add CST",
			tc: &tax.Combo{
				Category: br.TaxCategoryICMS,
				Ext:      tax.ExtensionsOf(cbc.CodeMap{nfe.ExtKeyICMSCSOSN: "102"}),
			},
			expect: map[cbc.Key]cbc.Code{
				nfe.ExtKeyICMSCSOSN:  "102",
				nfe.ExtKeyICMSCST:    "",
				nfe.ExtKeyICMSOrigin: "0",
			},
		},
		{
			name: "ICMS does not override CST or origin",
			tc: &tax.Combo{
				Category: br.TaxCategoryICMS,
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					nfe.ExtKeyICMSCST:    "40",
					nfe.ExtKeyICMSOrigin: "2",
				}),
			},
			expect: map[cbc.Key]cbc.Code{
				nfe.ExtKeyICMSCST:    "40",
				nfe.ExtKeyICMSOrigin: "2",
			},
		},
		{
			name:   "PIS sets default CST",
			tc:     &tax.Combo{Category: br.TaxCategoryPIS},
			expect: map[cbc.Key]cbc.Code{nfe.ExtKeyPISCST: "01"},
		},
		{
			name: "PIS does not override CST",
			tc: &tax.Combo{
				Category: br.TaxCategoryPIS,
				Ext:      tax.ExtensionsOf(cbc.CodeMap{nfe.ExtKeyPISCST: "49"}),
			},
			expect: map[cbc.Key]cbc.Code{nfe.ExtKeyPISCST: "49"},
		},
		{
			name:   "COFINS sets default CST",
			tc:     &tax.Combo{Category: br.TaxCategoryCOFINS},
			expect: map[cbc.Key]cbc.Code{nfe.ExtKeyCOFINSCST: "01"},
		},
		{
			name: "COFINS does not override CST",
			tc: &tax.Combo{
				Category: br.TaxCategoryCOFINS,
				Ext:      tax.ExtensionsOf(cbc.CodeMap{nfe.ExtKeyCOFINSCST: "49"}),
			},
			expect: map[cbc.Key]cbc.Code{nfe.ExtKeyCOFINSCST: "49"},
		},
		{
			name:   "unrelated category is untouched",
			tc:     &tax.Combo{Category: br.TaxCategoryIPI},
			expect: map[cbc.Key]cbc.Code{nfe.ExtKeyICMSCST: ""},
		},
	}
	for _, ts := range tests {
		t.Run(ts.name, func(t *testing.T) {
			norm.Normalize(ts.tc, tax.AddonContext(nfe.V4))
			for k, v := range ts.expect {
				assert.Equal(t, v, ts.tc.Ext.Get(k), "ext %s", k)
			}
		})
	}
}

func TestInvoiceRegimeNormalization(t *testing.T) {
	t.Run("defaults supplier regime to normal", func(t *testing.T) {
		inv := &bill.Invoice{
			Addons:   tax.WithAddons(nfe.V4),
			Supplier: &org.Party{Name: "Test Supplier"},
		}
		norm.Normalize(inv, tax.AddonContext(nfe.V4))
		assert.Equal(t, cbc.Code("3"), inv.Supplier.Ext.Get(nfe.ExtKeyRegime))
	})

	t.Run("does not override an existing supplier regime", func(t *testing.T) {
		inv := &bill.Invoice{
			Addons: tax.WithAddons(nfe.V4),
			Supplier: &org.Party{
				Name: "Test Supplier",
				Ext:  tax.ExtensionsOf(cbc.CodeMap{nfe.ExtKeyRegime: "1"}),
			},
		}
		norm.Normalize(inv, tax.AddonContext(nfe.V4))
		assert.Equal(t, cbc.Code("1"), inv.Supplier.Ext.Get(nfe.ExtKeyRegime))
	})

	t.Run("nil supplier is a no-op", func(t *testing.T) {
		inv := &bill.Invoice{Addons: tax.WithAddons(nfe.V4)}
		assert.NotPanics(t, func() {
			norm.Normalize(inv, tax.AddonContext(nfe.V4))
		})
	})
}
