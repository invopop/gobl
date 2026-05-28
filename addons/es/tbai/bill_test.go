package tbai

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoiceNormalization(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		ad := tax.AddonForKey(V1)
		var inv *bill.Invoice
		assert.NotPanics(t, func() {
			ad.Normalizer(inv)
		})
	})

	t.Run("standard invoice, no address", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Tax = nil
		require.NoError(t, inv.Calculate())
		assert.Nil(t, inv.Tax)
	})

	t.Run("standard invoice in Vizcaya", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Tax = nil
		inv.Supplier.Addresses = append(inv.Supplier.Addresses, &org.Address{
			Region: "Vizcaya",
		})
		require.NoError(t, inv.Calculate())
		assert.Equal(t, ExtValueRegionBI, inv.Tax.Ext.Get(ExtKeyRegion))
	})

	t.Run("standard invoice in Gipuzkoa", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Tax = nil
		inv.Supplier.Addresses = append(inv.Supplier.Addresses, &org.Address{
			Region: "Gipuzkoa",
		})
		require.NoError(t, inv.Calculate())
		assert.Equal(t, ExtValueRegionSS, inv.Tax.Ext.Get(ExtKeyRegion))
	})

	t.Run("standard invoice in Álava (accent)", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Tax = nil
		inv.Supplier.Addresses = append(inv.Supplier.Addresses, &org.Address{
			Region: "Álava",
		})
		require.NoError(t, inv.Calculate())
		assert.Equal(t, ExtValueRegionVI, inv.Tax.Ext.Get(ExtKeyRegion))
	})

	t.Run("standard invoice in Araba", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Tax = nil
		inv.Supplier.Addresses = append(inv.Supplier.Addresses, &org.Address{
			Region: "Araba",
		})
		require.NoError(t, inv.Calculate())
		assert.Equal(t, ExtValueRegionVI, inv.Tax.Ext.Get(ExtKeyRegion))
	})

	t.Run("standard invoice in Araba", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Tax = nil
		inv.Supplier.Addresses = append(inv.Supplier.Addresses, &org.Address{
			Region: "Madrid",
		})
		require.NoError(t, inv.Calculate())
		assert.Nil(t, inv.Tax)
	})

	t.Run("with existing region", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.Addresses = append(inv.Supplier.Addresses, &org.Address{
			Region: "Araba",
		})
		inv.Tax = &bill.Tax{
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				ExtKeyRegion: ExtValueRegionBI, // not Alaba
			}),
		}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, ExtValueRegionBI, inv.Tax.Ext.Get(ExtKeyRegion))
	})

	t.Run("regime defaults to 01", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Lines[0].Taxes[0] = &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyStandard,
			Rate:     tax.RateGeneral,
		}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, cbc.Code("01"), inv.Lines[0].Taxes[0].Ext.Get(ExtKeyRegime))
	})

	t.Run("regime 51 with equivalence surcharge", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Lines[0].Taxes[0] = &tax.Combo{
			Category: tax.CategoryVAT,
			Rate:     tax.RateGeneral.With(es.TaxRateEquivalence),
		}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, cbc.Code("51"), inv.Lines[0].Taxes[0].Ext.Get(ExtKeyRegime))
	})

	t.Run("regime 52 with simplified-scheme tag", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.SetTags(es.TagSimplifiedScheme)
		inv.Lines[0].Taxes[0] = &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyStandard,
			Rate:     tax.RateGeneral,
		}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, cbc.Code("52"), inv.Lines[0].Taxes[0].Ext.Get(ExtKeyRegime))
	})

	t.Run("regime 02 with export key", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Lines[0].Taxes[0] = &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyExport,
		}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, cbc.Code("02"), inv.Lines[0].Taxes[0].Ext.Get(ExtKeyRegime))
	})

	t.Run("regime explicit override is preserved", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.SetTags(es.TagSimplifiedScheme)
		inv.Lines[0].Taxes[0] = &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      tax.KeyStandard,
			Rate:     tax.RateGeneral,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				ExtKeyRegime: "07",
			}),
		}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, cbc.Code("07"), inv.Lines[0].Taxes[0].Ext.Get(ExtKeyRegime))
	})

	t.Run("regime applied to invoice-level charges and discounts", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Charges = []*bill.Charge{{
			Reason: "handling",
			Amount: num.MakeAmount(100, 2),
			Taxes: tax.Set{
				{Category: tax.CategoryVAT, Rate: tax.RateGeneral},
			},
		}}
		inv.Discounts = []*bill.Discount{{
			Reason: "loyalty",
			Amount: num.MakeAmount(50, 2),
			Taxes: tax.Set{
				{Category: tax.CategoryVAT, Rate: tax.RateGeneral},
			},
		}}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, cbc.Code("01"), inv.Charges[0].Taxes[0].Ext.Get(ExtKeyRegime))
		assert.Equal(t, cbc.Code("01"), inv.Discounts[0].Taxes[0].Ext.Get(ExtKeyRegime))
	})

	t.Run("simplified tag sets es-tbai-simplified=S", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer = nil
		inv.SetTags(tax.TagSimplified)
		require.NoError(t, inv.Calculate())
		assert.Equal(t, ExtValueSimplifiedYes, inv.Tax.Ext.Get(ExtKeySimplified))
	})

	t.Run("no simplified tag leaves es-tbai-simplified unset", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		assert.Empty(t, inv.Tax.Ext.Get(ExtKeySimplified).String())
	})
}

func TestInvoicePartyNormalization(t *testing.T) {
	t.Run("regular Spanish customer", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID = &tax.Identity{
			Country: "ES",
			Code:    "B12345678",
		}
		require.NoError(t, inv.Calculate())
	})

	t.Run("Spanish customer with identities is not normalized", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID = &tax.Identity{
			Country: "ES",
			Code:    "B12345678",
		}
		inv.Customer.Identities = []*org.Identity{
			{
				Key:  org.IdentityKeyPassport,
				Code: "AA123456",
			},
		}
		require.NoError(t, inv.Calculate())
		assert.True(t, inv.Customer.Identities[0].Ext.IsZero())
	})

	t.Run("customer without identities", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.Identities = nil
		require.NoError(t, inv.Calculate())
	})

	t.Run("passport identity normalization", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID = nil
		inv.Customer.Identities = []*org.Identity{
			{
				Key:  org.IdentityKeyPassport,
				Code: "AA123456",
			},
		}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, ExtCodeIdentityTypePassport, inv.Customer.Identities[0].Ext.Get(ExtKeyIdentityType))
	})

	t.Run("foreign identity normalization", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID = nil
		inv.Customer.Identities = []*org.Identity{
			{
				Key:  org.IdentityKeyForeign,
				Code: "FOR123456",
			},
		}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, ExtCodeIdentityTypeForeign, inv.Customer.Identities[0].Ext.Get(ExtKeyIdentityType))
	})

	t.Run("resident identity normalization", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID = nil
		inv.Customer.Identities = []*org.Identity{
			{
				Key:  org.IdentityKeyResident,
				Code: "RES123456",
			},
		}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, ExtCodeIdentityTypeResident, inv.Customer.Identities[0].Ext.Get(ExtKeyIdentityType))
	})

	t.Run("other identity normalization", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID = nil
		inv.Customer.Identities = []*org.Identity{
			{
				Key:  org.IdentityKeyOther,
				Code: "OTH123456",
			},
		}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, ExtCodeIdentityTypeOther, inv.Customer.Identities[0].Ext.Get(ExtKeyIdentityType))
	})

	t.Run("unknown identity key not normalized", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID = nil
		inv.Customer.Identities = []*org.Identity{
			{
				Key:  "unknown",
				Code: "UNK123456",
			},
		}
		require.NoError(t, inv.Calculate())
		assert.True(t, inv.Customer.Identities[0].Ext.IsZero())
	})

	t.Run("explicit extension on unkeyed identity preserved", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID = nil
		inv.Customer.Identities = []*org.Identity{
			{
				Code: "AA123456",
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					ExtKeyIdentityType: ExtCodeIdentityTypeOther,
				}),
			},
		}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, ExtCodeIdentityTypeOther, inv.Customer.Identities[0].Ext.Get(ExtKeyIdentityType))
	})

	t.Run("multiple identities only normalizes first", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID = nil
		inv.Customer.Identities = []*org.Identity{
			{
				Key:  org.IdentityKeyPassport,
				Code: "AA123456",
			},
			{
				Key:  org.IdentityKeyForeign,
				Code: "FOR123456",
			},
		}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, ExtCodeIdentityTypePassport, inv.Customer.Identities[0].Ext.Get(ExtKeyIdentityType))
		assert.True(t, inv.Customer.Identities[1].Ext.IsZero())
	})
}

func TestInvoiceValidation(t *testing.T) {
	t.Run("standard invoice", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("with services", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Lines[0].Taxes[0].Ext = inv.Lines[0].Taxes[0].Ext.Set(ExtKeyProduct, "services")
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("missing customer", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer = nil
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "customer is required for non-simplified invoices")
	})

	t.Run("missing customer tax ID and identity", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID = nil
		inv.Customer.Identities = nil
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "must have a tax_id or an identity with ext 'es-tbai-identity-type'")
	})

	t.Run("customer with identity-type extension but no tax ID is valid", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID = nil
		inv.Customer.Identities = []*org.Identity{
			{
				Country: "CH",
				Code:    "CH-OTHER-9001",
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					ExtKeyIdentityType: ExtCodeIdentityTypeOther,
				}),
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("customer with identity but no identity-type extension is rejected", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.TaxID = nil
		inv.Customer.Identities = []*org.Identity{
			{
				Country: "CH",
				Code:    "CH-XYZ-001",
				Key:     "unknown",
			},
		}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "must have a tax_id or an identity with ext 'es-tbai-identity-type'")
	})

	t.Run("simplified invoice without customer", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.SetTags(tax.TagSimplified)
		inv.Customer = nil
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("simplified invoice with customer without tax ID", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.SetTags(tax.TagSimplified)
		inv.Customer.TaxID = nil
		inv.Customer.Identities = nil
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("simplified invoice with customer tax ID", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.SetTags(tax.TagSimplified)
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("with exemption reason", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Lines[0].Taxes[0].Ext = tax.Extensions{}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.NoError(t, err)
	})

	t.Run("without series", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Series = ""
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("without notes", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Notes = nil
		assertValidationError(t, inv, "with key 'general' missing")
	})

	t.Run("correction", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Correct(
			bill.Credit,
			bill.WithExtension(ExtKeyCorrection, "R4"),
		))
		assert.Len(t, inv.Preceding, 1)
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("BI individual missing activity", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.TaxID = &tax.Identity{Country: "ES", Code: "12345678Z"}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "GOBL-ES-TBAI-BILL-INVOICE-10")
		assert.ErrorContains(t, err, "es-tbai-bi-activity")
	})

	t.Run("BI individual with valid activity", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.TaxID = &tax.Identity{Country: "ES", Code: "12345678Z"}
		inv.Supplier.Ext = tax.ExtensionsOf(cbc.CodeMap{
			ExtKeyBIActivity: "722300",
		})
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("BI persona jurídica without activity", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.TaxID = &tax.Identity{Country: "ES", Code: "B64847106"}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("VI individual without activity", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Tax.Ext = inv.Tax.Ext.Set(ExtKeyRegion, ExtValueRegionVI)
		inv.Supplier.TaxID = &tax.Identity{Country: "ES", Code: "12345678Z"}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("SS individual without activity", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Tax.Ext = inv.Tax.Ext.Set(ExtKeyRegion, ExtValueRegionSS)
		inv.Supplier.TaxID = &tax.Identity{Country: "ES", Code: "12345678Z"}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("BI individual with non-numeric activity", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.TaxID = &tax.Identity{Country: "ES", Code: "12345678Z"}
		inv.Supplier.Ext = tax.ExtensionsOf(cbc.CodeMap{
			ExtKeyBIActivity: "abc",
		})
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "es-tbai-bi-activity")
	})

	t.Run("BI individual with too-long activity", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.TaxID = &tax.Identity{Country: "ES", Code: "12345678Z"}
		inv.Supplier.Ext = tax.ExtensionsOf(cbc.CodeMap{
			ExtKeyBIActivity: "12345678",
		})
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "es-tbai-bi-activity")
	})

	t.Run("No tax", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Tax = nil
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "tax is required")
	})
}

func TestBillLineNormalization(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		ad := tax.AddonForKey(V1)
		var line *bill.Line
		assert.NotPanics(t, func() {
			ad.Normalizer(line)
		})
	})
	t.Run("with standard invoice, set default", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "services", inv.Lines[0].Taxes[0].Ext.Get(ExtKeyProduct).String())
	})
	t.Run("with standard invoice, set override for goods", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Lines[0].Item.Key = org.ItemKeyGoods
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "goods", inv.Lines[0].Taxes[0].Ext.Get(ExtKeyProduct).String())
	})
	t.Run("with standard invoice, set override for resale", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Lines[0].Item.Key = org.ItemKeyGoods
		inv.Lines[0].Taxes[0].Ext = inv.Lines[0].Taxes[0].Ext.Set(ExtKeyProduct, "resale")
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "resale", inv.Lines[0].Taxes[0].Ext.Get(ExtKeyProduct).String())
	})
}

func assertValidationError(t *testing.T, inv *bill.Invoice, expected string) {
	t.Helper()
	require.NoError(t, inv.Calculate())
	err := rules.Validate(inv)
	require.ErrorContains(t, err, expected)
}

func testInvoiceStandard(t *testing.T) *bill.Invoice {
	t.Helper()
	return &bill.Invoice{
		Addons: tax.WithAddons(V1),
		Series: "ABC",
		Code:   "123",
		Tax: &bill.Tax{
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				ExtKeyRegion: ExtValueRegionBI,
			}),
		},
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "B98602642",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "NL",
				Code:    "000099995B57",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "bogus",
					Price: num.NewAmount(10000, 2),
					Unit:  org.UnitPackage,
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Key:      "exempt",
						Ext: tax.ExtensionsOf(cbc.CodeMap{
							ExtKeyExempt: "E1",
						}),
					},
				},
			},
		},
		Notes: []*org.Note{
			{
				Key:  org.NoteKeyGeneral,
				Text: "This is a test invoice",
			},
		},
	}
}

func TestNormalizeInvoiceRegimeDefensive(t *testing.T) {
	t.Run("nil line/charge/discount is skipped", func(t *testing.T) {
		inv := &bill.Invoice{
			Lines:     []*bill.Line{nil},
			Charges:   []*bill.Charge{nil},
			Discounts: []*bill.Discount{nil},
		}
		assert.NotPanics(t, func() { normalizeInvoiceRegime(inv) })
	})

	t.Run("nil tax combo is skipped", func(t *testing.T) {
		inv := &bill.Invoice{
			Lines:     []*bill.Line{{Taxes: tax.Set{nil}}},
			Charges:   []*bill.Charge{{Taxes: tax.Set{nil}}},
			Discounts: []*bill.Discount{{Taxes: tax.Set{nil}}},
		}
		assert.NotPanics(t, func() { normalizeInvoiceRegime(inv) })
	})

	t.Run("non VAT/IGIC category is skipped", func(t *testing.T) {
		tc := &tax.Combo{Category: tax.CategoryGST}
		inv := &bill.Invoice{
			Lines: []*bill.Line{{Taxes: tax.Set{tc}}},
		}
		normalizeInvoiceRegime(inv)
		assert.True(t, tc.Ext.IsZero())
	})

	t.Run("charge/discount VAT combos get the default 01", func(t *testing.T) {
		ch := &tax.Combo{Category: tax.CategoryVAT}
		d := &tax.Combo{Category: tax.CategoryVAT}
		inv := &bill.Invoice{
			Charges:   []*bill.Charge{{Taxes: tax.Set{ch}}},
			Discounts: []*bill.Discount{{Taxes: tax.Set{d}}},
		}
		normalizeInvoiceRegime(inv)
		assert.Equal(t, cbc.Code("01"), ch.Ext.Get(ExtKeyRegime))
		assert.Equal(t, cbc.Code("01"), d.Ext.Get(ExtKeyRegime))
	})
}

func TestNormalizeInvoiceCustomerRates(t *testing.T) {
	makeInvoice := func() *bill.Invoice {
		return &bill.Invoice{
			Lines: []*bill.Line{
				{Taxes: tax.Set{{Category: tax.CategoryVAT, Key: tax.KeyStandard}}},
			},
			Charges: []*bill.Charge{
				{Taxes: tax.Set{{Category: tax.CategoryVAT, Key: tax.KeyStandard}}},
			},
			Discounts: []*bill.Discount{
				{Taxes: tax.Set{{Category: tax.CategoryVAT, Key: tax.KeyStandard}}},
			},
		}
	}

	t.Run("tag not set is a no-op", func(t *testing.T) {
		inv := makeInvoice()
		normalizeInvoiceCustomerRates(inv)
		assert.True(t, inv.Lines[0].Taxes[0].Ext.IsZero())
	})

	t.Run("standard combos pick up RL by default", func(t *testing.T) {
		inv := makeInvoice()
		inv.SetTags(tax.TagCustomerRates)
		normalizeInvoiceCustomerRates(inv)
		for _, tc := range []*tax.Combo{
			inv.Lines[0].Taxes[0],
			inv.Charges[0].Taxes[0],
			inv.Discounts[0].Taxes[0],
		} {
			assert.Equal(t, cbc.Code("RL"), tc.Ext.Get(ExtKeyExempt))
		}
	})

	t.Run("existing exempt code is preserved", func(t *testing.T) {
		inv := makeInvoice()
		inv.SetTags(tax.TagCustomerRates)
		inv.Lines[0].Taxes[0].Ext = tax.ExtensionsOf(cbc.CodeMap{ExtKeyExempt: "IE"})
		normalizeInvoiceCustomerRates(inv)
		assert.Equal(t, cbc.Code("IE"), inv.Lines[0].Taxes[0].Ext.Get(ExtKeyExempt))
	})

	t.Run("non-standard key is untouched", func(t *testing.T) {
		inv := makeInvoice()
		inv.Lines[0].Taxes[0].Key = tax.KeyReverseCharge
		inv.SetTags(tax.TagCustomerRates)
		normalizeInvoiceCustomerRates(inv)
		assert.True(t, inv.Lines[0].Taxes[0].Ext.IsZero())
	})

	t.Run("non-VAT/IGIC/IPSI combos are skipped", func(t *testing.T) {
		inv := makeInvoice()
		inv.SetTags(tax.TagCustomerRates)
		inv.Lines[0].Taxes[0].Category = tax.CategoryGST
		normalizeInvoiceCustomerRates(inv)
		assert.True(t, inv.Lines[0].Taxes[0].Ext.IsZero())
	})

	t.Run("nil line/charge/discount is skipped", func(t *testing.T) {
		inv := &bill.Invoice{
			Lines:     []*bill.Line{nil},
			Charges:   []*bill.Charge{nil},
			Discounts: []*bill.Discount{nil},
		}
		inv.SetTags(tax.TagCustomerRates)
		assert.NotPanics(t, func() { normalizeInvoiceCustomerRates(inv) })
	})

	t.Run("nil combo is skipped", func(t *testing.T) {
		inv := &bill.Invoice{
			Lines:     []*bill.Line{{Taxes: tax.Set{nil}}},
			Charges:   []*bill.Charge{{Taxes: tax.Set{nil}}},
			Discounts: []*bill.Discount{{Taxes: tax.Set{nil}}},
		}
		inv.SetTags(tax.TagCustomerRates)
		assert.NotPanics(t, func() { normalizeInvoiceCustomerRates(inv) })
	})
}

func TestNormalizeBillLineNoVAT(t *testing.T) {
	line := &bill.Line{
		Quantity: num.MakeAmount(1, 0),
		Item:     &org.Item{Name: "x", Price: num.NewAmount(100, 2)},
		Taxes: tax.Set{
			{Category: tax.CategoryGST},
		},
	}
	assert.NotPanics(t, func() { normalizeBillLine(line) })
	assert.True(t, line.Taxes[0].Ext.IsZero())
}

func TestNotesHasGeneralKeyWrongType(t *testing.T) {
	assert.False(t, notesHasGeneralKey("not a slice"))
	assert.False(t, notesHasGeneralKey(nil))
}

func TestNotesHasGeneralKeyNoGeneralNote(t *testing.T) {
	notes := []*org.Note{
		{Key: org.NoteKeyLegal, Text: "legal"},
	}
	assert.False(t, notesHasGeneralKey(notes))
}

func TestNormalizeInvoicePartyIdentityNilCustomer(t *testing.T) {
	assert.NotPanics(t, func() { normalizeInvoicePartyIdentity(nil) })
}

func TestNormalizeInvoicePartyIdentityUnkeyedNoExt(t *testing.T) {
	cus := &org.Party{
		Identities: []*org.Identity{
			{Code: "X"},
		},
	}
	normalizeInvoicePartyIdentity(cus)
	assert.True(t, cus.Identities[0].Ext.IsZero())
}

func TestNormalizeInvoicePartyIdentitySpanishNIFShortCircuits(t *testing.T) {
	cus := &org.Party{
		TaxID: &tax.Identity{Country: "ES", Code: "B12345678"},
		Identities: []*org.Identity{
			{Key: org.IdentityKeyPassport, Code: "AA"},
		},
	}
	normalizeInvoicePartyIdentity(cus)
	assert.True(t, cus.Identities[0].Ext.IsZero())
}

func TestNormalizeInvoicePartyIdentityEmptyIdentities(t *testing.T) {
	cus := &org.Party{}
	assert.NotPanics(t, func() { normalizeInvoicePartyIdentity(cus) })
}

func TestNormalizeInvoiceNil(t *testing.T) {
	assert.NotPanics(t, func() { normalizeInvoice(nil) })
}

func TestIsBizkaiaIndividualWrongType(t *testing.T) {
	assert.False(t, isBizkaiaIndividual("not an invoice"))
	assert.False(t, isBizkaiaIndividual(nil))
}
