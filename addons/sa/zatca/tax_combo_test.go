package zatca_test

import (
	"testing"

	"github.com/invopop/gobl/addons/eu/en16931"
	"github.com/invopop/gobl/addons/sa/zatca"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/cef"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// invoiceWithTaxCombo returns a valid standard invoice whose first line's
// tax set is replaced with the provided combo. Useful for exercising the
// taxComboRules in isolation while keeping all other invoice fields valid.
func invoiceWithTaxCombo(combo *tax.Combo) *bill.Invoice {
	inv := validStandardInvoice()
	inv.Lines[0].Taxes = tax.Set{combo}
	return inv
}

// --- Rule 01 (BR-KSA-CL-04): exempt/zero/outside-scope must have a valid SA VATEX ---

func TestTaxComboRule01_VATEXRequired(t *testing.T) {
	t.Run("exempt with valid VATEX-SA-29 (financial services) is valid", func(t *testing.T) {
		inv := invoiceWithTaxCombo(&tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyExempt,
			Ext: tax.ExtensionsOf(tax.ExtMap{
				cef.ExtKeyVATEX: zatca.VatexFinancialServices,
			}),
		})
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("exempt with valid VATEX-SA-30 (real-estate) is valid", func(t *testing.T) {
		inv := invoiceWithTaxCombo(&tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyExempt,
			Ext: tax.ExtensionsOf(tax.ExtMap{
				cef.ExtKeyVATEX: zatca.VatexRealEstate,
			}),
		})
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("exempt without any VATEX fails", func(t *testing.T) {
		inv := invoiceWithTaxCombo(&tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyExempt,
		})
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "VATEX exemption code must be present and valid for Z/E/O categories, and must not be set for Standard")
	})

	t.Run("exempt with VATEX from another category (VATEX-SA-32 belongs to Z) fails", func(t *testing.T) {
		inv := invoiceWithTaxCombo(&tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyExempt,
			Ext: tax.ExtensionsOf(tax.ExtMap{
				cef.ExtKeyVATEX: zatca.VatexExportGoods,
			}),
		})
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "VATEX exemption code must be present and valid for Z/E/O categories, and must not be set for Standard",
			"VATEX-SA-32 belongs to category Z, must not be accepted on exempt (E)")
	})

	t.Run("outside scope with VATEX-SA-OOS is valid", func(t *testing.T) {
		inv := invoiceWithTaxCombo(&tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyOutsideScope,
			Ext: tax.ExtensionsOf(tax.ExtMap{
				cef.ExtKeyVATEX: zatca.VatexOutOfScope,
			}),
		})
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("outside scope without VATEX fails", func(t *testing.T) {
		inv := invoiceWithTaxCombo(&tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyOutsideScope,
		})
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "VATEX exemption code must be present and valid for Z/E/O categories, and must not be set for Standard")
	})

	t.Run("outside scope with mismatched VATEX (Vatex29 belongs to E) fails", func(t *testing.T) {
		inv := invoiceWithTaxCombo(&tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyOutsideScope,
			Ext: tax.ExtensionsOf(tax.ExtMap{
				cef.ExtKeyVATEX: zatca.VatexFinancialServices,
			}),
		})
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "VATEX exemption code must be present and valid for Z/E/O categories, and must not be set for Standard")
	})

	// NOTE: zero-rated (Z) + VATEX combinations cannot currently be tested at
	// the invoice level because EN16931 rule 07 (BR-S-10/BR-Z-10) prohibits
	// any VATEX on Z-rated lines and runs first. See ISSUE TRACKER #18.
}

// --- Rule 02 (BR-KSA-CL-04): standard rate must NOT have a VATEX code ---

func TestTaxComboRule02_StandardNoVATEX(t *testing.T) {
	t.Run("standard rate without VATEX is valid", func(t *testing.T) {
		inv := invoiceWithTaxCombo(&tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateGeneral,
		})
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("standard rate with VATEX fails", func(t *testing.T) {
		inv := invoiceWithTaxCombo(&tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateGeneral,
			Ext: tax.ExtensionsOf(tax.ExtMap{
				cef.ExtKeyVATEX: zatca.VatexFinancialServices,
			}),
		})
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "VATEX exemption code must be present and valid for Z/E/O categories, and must not be set for Standard")
	})
}

// --- Rule 03 (BR-KSA-18): VAT category code must be one of S, Z, E, O ---

func TestTaxComboRule03_ValidCategoryCodes(t *testing.T) {
	validCategories := map[string]struct {
		rate     cbc.Key
		key      cbc.Key
		expected cbc.Code
		ext      tax.Extensions
	}{
		"standard maps to S": {
			rate:     tax.RateGeneral,
			expected: en16931.TaxCategoryStandard,
		},
		"exempt maps to E": {
			key:      tax.KeyExempt,
			expected: en16931.TaxCategoryExempt,
			ext:      tax.ExtensionsOf(tax.ExtMap{cef.ExtKeyVATEX: zatca.VatexFinancialServices}),
		},
		"outside scope maps to O": {
			key:      tax.KeyOutsideScope,
			expected: en16931.TaxCategoryOutsideScope,
			ext:      tax.ExtensionsOf(tax.ExtMap{cef.ExtKeyVATEX: zatca.VatexOutOfScope}),
		},
	}

	for name, tc := range validCategories {
		t.Run(name, func(t *testing.T) {
			inv := invoiceWithTaxCombo(&tax.Combo{
				Category: tax.CategoryVAT,
				Key:      tc.key,
				Rate:     tc.rate,
				Ext:      tc.ext,
			})
			require.NoError(t, inv.Calculate())
			// Confirm the category extension was set correctly during normalization
			assert.Equal(t, tc.expected, inv.Lines[0].Taxes[0].Ext.Get(untdid.ExtKeyTaxCategory))
			require.NoError(t, rules.Validate(inv))
		})
	}

	t.Run("category code outside the SA-allowed subset fails", func(t *testing.T) {
		// Manually inject an EN16931-valid but ZATCA-disallowed category (e.g. AE
		// = reverse charge). This exercises rule 03 directly without relying on
		// EN16931 normalisation (which would refuse to set AE for KeyExempt).
		inv := invoiceWithTaxCombo(&tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyExempt,
			Ext: tax.ExtensionsOf(tax.ExtMap{
				cef.ExtKeyVATEX: zatca.VatexFinancialServices,
			}),
		})
		require.NoError(t, inv.Calculate())
		// Override the normaliser's choice to an out-of-subset value.
		inv.Lines[0].Taxes[0].Ext = inv.Lines[0].Taxes[0].Ext.Set(untdid.ExtKeyTaxCategory, en16931.TaxCategoryReverseCharge)
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "VAT category code must contain one of the values (S, Z, E, O)")
	})
}

// --- Rule 04 (BR-KSA-54): tax category must be 'VAT' ---

func TestTaxComboRule04_CategoryMustBeVAT(t *testing.T) {
	t.Run("VAT category is valid", func(t *testing.T) {
		inv := invoiceWithTaxCombo(&tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateGeneral,
		})
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("non-VAT category fails", func(t *testing.T) {
		// Use an arbitrary non-VAT category (the rule only inspects "cat",
		// regardless of whether the regime defines it).
		inv := invoiceWithTaxCombo(&tax.Combo{
			Category: cbc.Code("OTHER"),
			Rate:     tax.RateGeneral,
		})
		_ = inv.Calculate()
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "tax schema id must be 'VAT'")
	})
}

// --- Cross-cutting: VATEX must match the combo's tax category ---

func TestTaxComboVATEX_PerCategoryRestrictions(t *testing.T) {
	cases := []struct {
		name    string
		key     cbc.Key
		vatex   cbc.Code
		wantErr bool
	}{
		{"E+Vatex29 (financial) ok", tax.KeyExempt, zatca.VatexFinancialServices, false},
		{"E+Vatex29_7 (life insurance) ok", tax.KeyExempt, zatca.VatexLifeInsurance, false},
		{"E+Vatex30 (real estate) ok", tax.KeyExempt, zatca.VatexRealEstate, false},
		{"E+Vatex32 (Z code) rejected", tax.KeyExempt, zatca.VatexExportGoods, true},
		{"E+VatexEdu (Z code) rejected", tax.KeyExempt, zatca.VatexPrivateEducation, true},
		{"E+VatexOutOfScope (O code) rejected", tax.KeyExempt, zatca.VatexOutOfScope, true},

		{"O+VatexOutOfScope ok", tax.KeyOutsideScope, zatca.VatexOutOfScope, false},
		{"O+Vatex29 (E code) rejected", tax.KeyOutsideScope, zatca.VatexFinancialServices, true},
		{"O+Vatex32 (Z code) rejected", tax.KeyOutsideScope, zatca.VatexExportGoods, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			inv := invoiceWithTaxCombo(&tax.Combo{
				Category: tax.CategoryVAT,
				Key:      tc.key,
				Ext: tax.ExtensionsOf(tax.ExtMap{
					cef.ExtKeyVATEX: tc.vatex,
				}),
			})
			require.NoError(t, inv.Calculate())
			err := rules.Validate(inv)
			if tc.wantErr {
				assert.ErrorContains(t, err, "VATEX exemption code must be present and valid for Z/E/O categories, and must not be set for Standard")
			} else {
				require.NoError(t, err, "%s: expected combo to validate", tc.name)
			}
		})
	}
}
