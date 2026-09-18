package verifactu

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeTaxCombo(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyStandard,
			Rate:     tax.RateGeneral,
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "S1", tc.Ext.Get(ExtKeyOpClass).String())
		assert.Equal(t, "01", tc.Ext.Get(ExtKeyRegime).String())
	})
	t.Run("valid - no key", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateGeneral,
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "S1", tc.Ext.Get(ExtKeyOpClass).String())
		assert.Equal(t, "01", tc.Ext.Get(ExtKeyRegime).String())
	})

	t.Run("valid with country", func(t *testing.T) {
		tc := &tax.Combo{
			Country:  l10n.ES.Tax(),
			Category: tax.CategoryVAT,
			Rate:     tax.RateSuperReduced,
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "S1", tc.Ext.Get(ExtKeyOpClass).String())
		assert.Equal(t, "01", tc.Ext.Get(ExtKeyRegime).String())
		assert.Empty(t, tc.Ext.Get(ExtKeyExempt))
	})

	t.Run("exempt", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyExempt,
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "01", tc.Ext.Get(ExtKeyRegime).String())
		assert.Equal(t, "E1", tc.Ext.Get(ExtKeyExempt).String())
		assert.Empty(t, tc.Ext.Get(ExtKeyOpClass))
	})
	t.Run("export", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyExport,
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "02", tc.Ext.Get(ExtKeyRegime).String())
		assert.Equal(t, "E2", tc.Ext.Get(ExtKeyExempt).String())
		assert.Empty(t, tc.Ext.Get(ExtKeyOpClass))
	})
	t.Run("surcharge", func(t *testing.T) {
		tc := &tax.Combo{
			Category:  tax.CategoryVAT,
			Rate:      tax.RateGeneral.With(es.TaxRateEquivalence),
			Percent:   num.NewPercentage(210, 3),
			Surcharge: num.NewPercentage(50, 3),
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "18", tc.Ext.Get(ExtKeyRegime).String())
		assert.Equal(t, "S1", tc.Ext.Get(ExtKeyOpClass).String())
		assert.Empty(t, tc.Ext.Get(ExtKeyExempt))
	})
	t.Run("intra-community", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyIntraCommunity,
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "01", tc.Ext.Get(ExtKeyRegime).String())
		assert.Equal(t, "E5", tc.Ext.Get(ExtKeyExempt).String())
		assert.Empty(t, tc.Ext.Get(ExtKeyOpClass))
	})
	t.Run("outside scope", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyOutsideScope,
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "01", tc.Ext.Get(ExtKeyRegime).String())
		assert.Equal(t, "N2", tc.Ext.Get(ExtKeyOpClass).String())
		assert.Empty(t, tc.Ext.Get(ExtKeyExempt))
	})

	t.Run("outside scope N1", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyOutsideScope,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				ExtKeyOpClass: "N1",
			}),
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "01", tc.Ext.Get(ExtKeyRegime).String())
		assert.Equal(t, "N1", tc.Ext.Get(ExtKeyOpClass).String())
		assert.Empty(t, tc.Ext.Get(ExtKeyExempt))
	})
	t.Run("reverse charge", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyReverseCharge,
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "01", tc.Ext.Get(ExtKeyRegime).String())
		assert.Equal(t, "S2", tc.Ext.Get(ExtKeyOpClass).String())
		assert.Empty(t, tc.Ext.Get(ExtKeyExempt))
	})

	t.Run("foreign country", func(t *testing.T) {
		tc := &tax.Combo{
			Country:  "FR",
			Category: tax.CategoryVAT,
			Rate:     tax.RateGeneral,
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, cbc.Code("N2"), tc.Ext.Get(ExtKeyOpClass))
		assert.Empty(t, tc.Ext.Get(ExtKeyExempt))
	})

	t.Run("with tax regime", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateGeneral,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				ExtKeyRegime: "03",
			}),
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "03", tc.Ext.Get(ExtKeyRegime).String())
	})

	t.Run("with exempt code set", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				ExtKeyExempt: "E6",
			}),
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "E6", tc.Ext.Get(ExtKeyExempt).String())
		assert.Equal(t, "01", tc.Ext.Get(ExtKeyRegime).String())
		assert.Empty(t, tc.Ext.Get(ExtKeyOpClass).String())
		assert.Equal(t, tax.KeyExempt, tc.Key)
	})
	t.Run("with export code set", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				ExtKeyExempt: "E2",
				ExtKeyRegime: "02",
			}),
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "E2", tc.Ext.Get(ExtKeyExempt).String())
		assert.Equal(t, "02", tc.Ext.Get(ExtKeyRegime).String())
		assert.Equal(t, tax.KeyExport, tc.Key)
	})
	t.Run("with reverse-charge", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				ExtKeyOpClass: "S2",
				ExtKeyRegime:  "01",
			}),
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "S2", tc.Ext.Get(ExtKeyOpClass).String())
		assert.Equal(t, "01", tc.Ext.Get(ExtKeyRegime).String())
		assert.Equal(t, tax.KeyReverseCharge, tc.Key)
	})
	t.Run("with outside-scope", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				ExtKeyOpClass: "N1",
				ExtKeyRegime:  "01",
			}),
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "N1", tc.Ext.Get(ExtKeyOpClass).String())
		assert.Equal(t, "01", tc.Ext.Get(ExtKeyRegime).String())
		assert.Equal(t, tax.KeyOutsideScope, tc.Key)
	})
	t.Run("with intra-community", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				ExtKeyExempt: "E5",
				ExtKeyRegime: "01",
			}),
		}
		normalizeTaxCombo(tc)
		assert.Equal(t, "E5", tc.Ext.Get(ExtKeyExempt).String())
		assert.Equal(t, "01", tc.Ext.Get(ExtKeyRegime).String())
		assert.Equal(t, tax.KeyIntraCommunity, tc.Key)
	})
}

func TestValidateTaxCombo(t *testing.T) {
	ruleSet := taxComboRules()

	t.Run("valid", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateGeneral,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				ExtKeyOpClass: "S1",
				ExtKeyRegime:  "01",
			}),
		}
		err := ruleSet.Validate(tc)
		assert.NoError(t, err)
	})

	t.Run("not in category", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryGST,
			Rate:     tax.RateGeneral,
		}
		err := ruleSet.Validate(tc)
		assert.NoError(t, err)
	})

	t.Run("exempt with valid reason", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyExempt,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				ExtKeyRegime: "01",
				ExtKeyExempt: "E1",
			}),
		}
		err := ruleSet.Validate(tc)
		assert.NoError(t, err)
	})

	t.Run("excludes E2 exemption code with regime 01", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				ExtKeyRegime: "01",
				ExtKeyExempt: "E2",
			}),
		}
		err := ruleSet.Validate(tc)
		assert.ErrorContains(t, err, "E2")
	})

	t.Run("excludes E3 exemption code with regime 01", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				ExtKeyRegime: "01",
				ExtKeyExempt: "E3",
			}),
		}
		err := ruleSet.Validate(tc)
		assert.ErrorContains(t, err, "E3")
	})

	t.Run("allows E2 exemption code with non-01 regime", func(t *testing.T) {
		tc := &tax.Combo{
			Category: tax.CategoryVAT,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				ExtKeyRegime: "02",
				ExtKeyExempt: "E2",
			}),
		}
		err := ruleSet.Validate(tc)
		assert.NoError(t, err)
	})

	t.Run("excludes E2 exemption code with regime 01 and IGIC category", func(t *testing.T) {
		tc := &tax.Combo{
			Category: es.TaxCategoryIGIC,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				ExtKeyRegime: "01",
				ExtKeyExempt: "E2",
			}),
		}
		err := ruleSet.Validate(tc)
		assert.ErrorContains(t, err, "E2")
	})
}

func TestNormalizeTaxComboIPSI(t *testing.T) {
	tests := []struct {
		name    string
		combo   *tax.Combo
		key     cbc.Key
		opClass cbc.Code
		exempt  cbc.Code
		regime  cbc.Code
	}{
		{
			name:    "percent only",
			combo:   &tax.Combo{Category: es.TaxCategoryIPSI, Percent: num.NewPercentage(4, 2)},
			key:     tax.KeyStandard,
			opClass: "S1",
			regime:  "01",
		},
		{
			name: "explicit regime preserved",
			combo: &tax.Combo{
				Category: es.TaxCategoryIPSI,
				Percent:  num.NewPercentage(4, 2),
				Ext:      tax.ExtensionsOf(cbc.CodeMap{ExtKeyRegime: "11"}),
			},
			key:     tax.KeyStandard,
			opClass: "S1",
			regime:  "11",
		},
		{
			name: "exempt with explicit regime preserved",
			combo: &tax.Combo{
				Category: es.TaxCategoryIPSI,
				Key:      tax.KeyExempt,
				Ext:      tax.ExtensionsOf(cbc.CodeMap{ExtKeyRegime: "01"}),
			},
			key:    tax.KeyExempt,
			exempt: "E1",
			regime: "01",
		},
		{
			name:    "standard",
			combo:   &tax.Combo{Category: es.TaxCategoryIPSI, Key: tax.KeyStandard, Percent: num.NewPercentage(4, 2)},
			key:     tax.KeyStandard,
			opClass: "S1",
		},
		{
			name:    "zero",
			combo:   &tax.Combo{Category: es.TaxCategoryIPSI, Key: tax.KeyZero, Percent: num.NewPercentage(0, 2)},
			key:     tax.KeyZero,
			opClass: "S1",
		},
		{
			name:    "reverse charge",
			combo:   &tax.Combo{Category: es.TaxCategoryIPSI, Key: tax.KeyReverseCharge},
			key:     tax.KeyReverseCharge,
			opClass: "S2",
		},
		{
			name:    "outside scope",
			combo:   &tax.Combo{Category: es.TaxCategoryIPSI, Key: tax.KeyOutsideScope},
			key:     tax.KeyOutsideScope,
			opClass: "N2",
		},
		{
			name: "outside scope N1 preserved",
			combo: &tax.Combo{
				Category: es.TaxCategoryIPSI,
				Key:      tax.KeyOutsideScope,
				Ext:      tax.ExtensionsOf(cbc.CodeMap{ExtKeyOpClass: "N1"}),
			},
			key:     tax.KeyOutsideScope,
			opClass: "N1",
		},
		{
			name:   "exempt",
			combo:  &tax.Combo{Category: es.TaxCategoryIPSI, Key: tax.KeyExempt},
			key:    tax.KeyExempt,
			exempt: "E1",
		},
		{
			name: "exempt E6 preserved",
			combo: &tax.Combo{
				Category: es.TaxCategoryIPSI,
				Key:      tax.KeyExempt,
				Ext:      tax.ExtensionsOf(cbc.CodeMap{ExtKeyExempt: "E6"}),
			},
			key:    tax.KeyExempt,
			exempt: "E6",
		},
		{
			name: "exempt removes stale op class",
			combo: &tax.Combo{
				Category: es.TaxCategoryIPSI,
				Key:      tax.KeyExempt,
				Ext:      tax.ExtensionsOf(cbc.CodeMap{ExtKeyOpClass: "S1"}),
			},
			key:    tax.KeyExempt,
			exempt: "E1",
		},
		{
			name: "standard removes stale exempt code",
			combo: &tax.Combo{
				Category: es.TaxCategoryIPSI,
				Key:      tax.KeyStandard,
				Percent:  num.NewPercentage(4, 2),
				Ext:      tax.ExtensionsOf(cbc.CodeMap{ExtKeyExempt: "E1"}),
			},
			key:     tax.KeyStandard,
			opClass: "S1",
		},
		{
			name:   "export",
			combo:  &tax.Combo{Category: es.TaxCategoryIPSI, Key: tax.KeyExport},
			key:    tax.KeyExport,
			exempt: "E2",
		},
		{
			name: "export E3 preserved",
			combo: &tax.Combo{
				Category: es.TaxCategoryIPSI,
				Key:      tax.KeyExport,
				Ext:      tax.ExtensionsOf(cbc.CodeMap{ExtKeyExempt: "E3"}),
			},
			key:    tax.KeyExport,
			exempt: "E3",
		},
		{
			name: "legacy explicit S1 with percent",
			combo: &tax.Combo{
				Category: es.TaxCategoryIPSI,
				Percent:  num.NewPercentage(4, 2),
				Ext:      tax.ExtensionsOf(cbc.CodeMap{ExtKeyOpClass: "S1"}),
			},
			key:     tax.KeyStandard,
			opClass: "S1",
		},
		{
			name: "legacy explicit N1 without key",
			combo: &tax.Combo{
				Category: es.TaxCategoryIPSI,
				Ext:      tax.ExtensionsOf(cbc.CodeMap{ExtKeyOpClass: "N1"}),
			},
			key:     tax.KeyOutsideScope,
			opClass: "N1",
		},
		{
			name: "legacy explicit E6 without key",
			combo: &tax.Combo{
				Category: es.TaxCategoryIPSI,
				Ext:      tax.ExtensionsOf(cbc.CodeMap{ExtKeyExempt: "E6"}),
			},
			key:    tax.KeyExempt,
			exempt: "E6",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			normalizeTaxCombo(tt.combo)
			assert.Equal(t, tt.key, tt.combo.Key)
			assert.Equal(t, tt.opClass, tt.combo.Ext.Get(ExtKeyOpClass), "op class")
			assert.Equal(t, tt.exempt, tt.combo.Ext.Get(ExtKeyExempt), "exempt")
			regime := tt.regime
			if regime == "" {
				// Default regime: 19 for E1 exemptions, 01 otherwise.
				regime = "01"
				if tt.exempt == "E1" {
					regime = "19"
				}
			}
			assert.Equal(t, regime, tt.combo.Ext.Get(ExtKeyRegime), "regime")
		})
	}
}

func TestValidateTaxComboIPSI(t *testing.T) {
	ruleSet := taxComboRules()

	t.Run("taxed with op class", func(t *testing.T) {
		tc := &tax.Combo{
			Category: es.TaxCategoryIPSI,
			Percent:  num.NewPercentage(4, 2),
			Ext:      tax.ExtensionsOf(cbc.CodeMap{ExtKeyOpClass: "S1", ExtKeyRegime: "01"}),
		}
		assert.NoError(t, ruleSet.Validate(tc))
	})

	t.Run("exempt with regime 19", func(t *testing.T) {
		tc := &tax.Combo{
			Category: es.TaxCategoryIPSI,
			Key:      tax.KeyExempt,
			Ext:      tax.ExtensionsOf(cbc.CodeMap{ExtKeyExempt: "E1", ExtKeyRegime: "19"}),
		}
		assert.NoError(t, ruleSet.Validate(tc))
	})

	t.Run("export E2 allowed with regime 01", func(t *testing.T) {
		// The E2/E3 restriction on regime 01 only applies to VAT and IGIC.
		tc := &tax.Combo{
			Category: es.TaxCategoryIPSI,
			Key:      tax.KeyExport,
			Ext:      tax.ExtensionsOf(cbc.CodeMap{ExtKeyExempt: "E2", ExtKeyRegime: "01"}),
		}
		assert.NoError(t, ruleSet.Validate(tc))
	})

	t.Run("regime required", func(t *testing.T) {
		tc := &tax.Combo{
			Category: es.TaxCategoryIPSI,
			Key:      tax.KeyExempt,
			Ext:      tax.ExtensionsOf(cbc.CodeMap{ExtKeyExempt: "E1"}),
		}
		err := ruleSet.Validate(tc)
		assert.ErrorContains(t, err, "es-verifactu-regime' is required")
	})

	t.Run("regime not valid for IPSI", func(t *testing.T) {
		tc := &tax.Combo{
			Category: es.TaxCategoryIPSI,
			Key:      tax.KeyExport,
			Ext:      tax.ExtensionsOf(cbc.CodeMap{ExtKeyExempt: "E2", ExtKeyRegime: "02"}),
		}
		err := ruleSet.Validate(tc)
		assert.ErrorContains(t, err, "es-verifactu-regime' for IPSI must be one of")
	})

	t.Run("taxed without op class", func(t *testing.T) {
		tc := &tax.Combo{
			Category: es.TaxCategoryIPSI,
			Percent:  num.NewPercentage(4, 2),
			Ext:      tax.ExtensionsOf(cbc.CodeMap{ExtKeyRegime: "01"}),
		}
		err := ruleSet.Validate(tc)
		assert.ErrorContains(t, err, "es-verifactu-op-class' is required for taxed operations")
	})

	t.Run("op class and exempt together", func(t *testing.T) {
		tc := &tax.Combo{
			Category: es.TaxCategoryIPSI,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				ExtKeyOpClass: "S1",
				ExtKeyExempt:  "E1",
				ExtKeyRegime:  "01",
			}),
		}
		err := ruleSet.Validate(tc)
		assert.ErrorContains(t, err, "cannot use both")
	})
}

func TestInvoiceIPSI(t *testing.T) {
	newInvoice := func(taxes ...*tax.Combo) *bill.Invoice {
		lines := make([]*bill.Line, len(taxes))
		for i, tc := range taxes {
			lines[i] = &bill.Line{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Service",
					Price: num.NewAmount(10000, 2),
				},
				Taxes: tax.Set{tc},
			}
		}
		return &bill.Invoice{
			Addons: tax.WithAddons(V1),
			Code:   "IPSI-1",
			Supplier: &org.Party{
				Name: "Ceuta Supplier",
				TaxID: &tax.Identity{
					Country: "ES",
					Code:    "B98602642",
				},
			},
			Customer: &org.Party{
				Name: "Customer",
				TaxID: &tax.Identity{
					Country: "ES",
					Code:    "54387763P",
				},
			},
			Lines: lines,
		}
	}

	t.Run("acceptance cases", func(t *testing.T) {
		inv := newInvoice(
			&tax.Combo{Category: es.TaxCategoryIPSI, Percent: num.NewPercentage(4, 2)},
			&tax.Combo{Category: es.TaxCategoryIPSI, Key: tax.KeyExempt},
			&tax.Combo{
				Category: es.TaxCategoryIPSI,
				Key:      tax.KeyOutsideScope,
				Ext:      tax.ExtensionsOf(cbc.CodeMap{ExtKeyOpClass: "N1"}),
			},
			&tax.Combo{Category: es.TaxCategoryIPSI, Key: tax.KeyOutsideScope},
			&tax.Combo{Category: es.TaxCategoryIPSI, Key: tax.KeyZero, Percent: num.NewPercentage(0, 2)},
		)
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))

		expected := []struct {
			key     cbc.Key
			percent string
			ext     cbc.CodeMap
		}{
			{tax.KeyStandard, "4%", cbc.CodeMap{ExtKeyOpClass: "S1", ExtKeyRegime: "01"}},
			{tax.KeyExempt, "", cbc.CodeMap{ExtKeyExempt: "E1", ExtKeyRegime: "19"}},
			{tax.KeyOutsideScope, "", cbc.CodeMap{ExtKeyOpClass: "N1", ExtKeyRegime: "01"}},
			{tax.KeyOutsideScope, "", cbc.CodeMap{ExtKeyOpClass: "N2", ExtKeyRegime: "01"}},
			{tax.KeyZero, "0%", cbc.CodeMap{ExtKeyOpClass: "S1", ExtKeyRegime: "01"}},
		}
		for i, exp := range expected {
			tc := inv.Lines[i].Taxes[0]
			assert.Equal(t, exp.key, tc.Key, "line %d key", i+1)
			if exp.percent == "" {
				assert.Nil(t, tc.Percent, "line %d percent", i+1)
			} else {
				require.NotNil(t, tc.Percent, "line %d percent", i+1)
				assert.Equal(t, exp.percent, tc.Percent.String(), "line %d percent", i+1)
			}
			assert.Equal(t, tax.ExtensionsOf(exp.ext), tc.Ext, "line %d ext", i+1)
		}

		// Totals must carry the extensions through to the rate breakdown.
		cat := inv.Totals.Taxes.Category(es.TaxCategoryIPSI)
		require.NotNil(t, cat)
		for _, r := range cat.Rates {
			assert.Contains(t, ipsiRegimeCodes, r.Ext.Get(ExtKeyRegime))
			assert.True(t, r.Ext.Has(ExtKeyOpClass) || r.Ext.Has(ExtKeyExempt))
		}
	})

	t.Run("intra-community rejected", func(t *testing.T) {
		inv := newInvoice(&tax.Combo{Category: es.TaxCategoryIPSI, Key: tax.KeyIntraCommunity})
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "tax combo key not valid for category in regime")
	})

	t.Run("percent required without key", func(t *testing.T) {
		inv := newInvoice(&tax.Combo{Category: es.TaxCategoryIPSI})
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "percent required or invalid")
	})
}
