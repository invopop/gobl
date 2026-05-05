package zatca_test

import (
	"testing"

	"github.com/invopop/gobl/addons/eu/en16931"
	"github.com/invopop/gobl/addons/sa/zatca"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/cef"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeInvoice(t *testing.T) {
	ad := tax.AddonForKey(zatca.V1)

	t.Run("nil invoice does not panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			ad.Normalizer((*bill.Invoice)(nil))
		})
	})

	t.Run("sets rounding to currency", func(t *testing.T) {
		inv := validStandardInvoice()
		ad.Normalizer(inv)
		assert.Equal(t, tax.RoundingRuleCurrency, inv.Tax.Rounding)
	})

	t.Run("creates tax object when nil", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Tax = nil
		ad.Normalizer(inv)
		require.NotNil(t, inv.Tax)
		assert.Equal(t, tax.RoundingRuleCurrency, inv.Tax.Rounding)
	})

	t.Run("creates issue time when nil", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.IssueTime = nil
		ad.Normalizer(inv)
		require.NotNil(t, inv.IssueTime)
	})

	t.Run("outside scope line gets zero percent", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Lines = []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Out of scope item",
					Price: num.NewAmount(100, 0),
				},
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Key:      tax.KeyOutsideScope,
					},
				},
			},
		}
		ad.Normalizer(inv)
		vat := inv.Lines[0].Taxes.Get(tax.CategoryVAT)
		require.NotNil(t, vat.Percent)
		assert.True(t, vat.Percent.IsZero())
	})

	t.Run("line without VAT combo is skipped", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Lines = []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "No tax item",
					Price: num.NewAmount(50, 0),
				},
				Taxes: tax.Set{},
			},
		}
		assert.NotPanics(t, func() {
			ad.Normalizer(inv)
		})
	})
}

func TestNormalizeInvoiceExemptionNotes(t *testing.T) {
	ad := tax.AddonForKey(zatca.V1)

	t.Run("exempt line gets tax note", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Lines = []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Exempt item",
					Price: num.NewAmount(100, 0),
				},
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Key:      tax.KeyExempt,
						Ext: tax.ExtensionsOf(cbc.CodeMap{
							cef.ExtKeyVATEX:          "VATEX-SA-29",
							untdid.ExtKeyTaxCategory: en16931.TaxCategoryExempt,
						}),
					},
				},
			},
		}
		ad.Normalizer(inv)
		require.Len(t, inv.Tax.Notes, 1)
		n := inv.Tax.Notes[0]
		assert.Equal(t, tax.CategoryVAT, n.Category)
		assert.Equal(t, tax.KeyExempt, n.Key)
		assert.Equal(t, "Financial services mentioned in Article 29 of the VAT Regulations", n.Text)
		assert.Equal(t, en16931.TaxCategoryExempt, n.Ext.Get(untdid.ExtKeyTaxCategory))
	})

	t.Run("zero-rated line gets tax note", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Lines = []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Export item",
					Price: num.NewAmount(100, 0),
				},
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Key:      tax.KeyZero,
						Ext: tax.ExtensionsOf(cbc.CodeMap{
							cef.ExtKeyVATEX:          "VATEX-SA-32",
							untdid.ExtKeyTaxCategory: en16931.TaxCategoryZero,
						}),
					},
				},
			},
		}
		ad.Normalizer(inv)
		require.Len(t, inv.Tax.Notes, 1)
		n := inv.Tax.Notes[0]
		assert.Equal(t, tax.CategoryVAT, n.Category)
		assert.Equal(t, tax.KeyZero, n.Key)
		assert.Equal(t, "Export of goods", n.Text)
		assert.Equal(t, en16931.TaxCategoryZero, n.Ext.Get(untdid.ExtKeyTaxCategory))
	})

	t.Run("outside-scope line gets tax note", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Lines = []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "OOS item",
					Price: num.NewAmount(100, 0),
				},
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Key:      tax.KeyOutsideScope,
						Ext: tax.ExtensionsOf(cbc.CodeMap{
							cef.ExtKeyVATEX:          "VATEX-SA-OOS",
							untdid.ExtKeyTaxCategory: en16931.TaxCategoryOutsideScope,
						}),
					},
				},
			},
		}
		ad.Normalizer(inv)
		require.Len(t, inv.Tax.Notes, 1)
		n := inv.Tax.Notes[0]
		assert.Equal(t, tax.CategoryVAT, n.Category)
		assert.Equal(t, tax.KeyOutsideScope, n.Key)
		assert.Equal(t, "Reason is free text, to be provided by the taxpayer on case to case basis", n.Text)
		assert.Equal(t, en16931.TaxCategoryOutsideScope, n.Ext.Get(untdid.ExtKeyTaxCategory))
	})

	t.Run("standard VAT line does not add note", func(t *testing.T) {
		inv := validStandardInvoice()
		ad.Normalizer(inv)
		assert.Empty(t, inv.Tax.Notes)
	})

	t.Run("existing note for same category not duplicated", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Tax = &bill.Tax{
			Notes: []*tax.Note{
				{
					Category: tax.CategoryVAT,
					Key:      tax.KeyExempt,
					Text:     "Existing exemption note",
					Ext:      tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyTaxCategory: en16931.TaxCategoryExempt}),
				},
			},
		}
		inv.Lines = []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Exempt item",
					Price: num.NewAmount(100, 0),
				},
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Key:      tax.KeyExempt,
						Ext: tax.ExtensionsOf(cbc.CodeMap{
							cef.ExtKeyVATEX:          "VATEX-SA-29",
							untdid.ExtKeyTaxCategory: en16931.TaxCategoryExempt,
						}),
					},
				},
			},
		}
		ad.Normalizer(inv)
		assert.Len(t, inv.Tax.Notes, 1)
		assert.Equal(t, "Existing exemption note", inv.Tax.Notes[0].Text)
	})

	t.Run("unknown VATEX code does not add note", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Lines = []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Unknown item",
					Price: num.NewAmount(100, 0),
				},
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Key:      tax.KeyExempt,
						Ext: tax.ExtensionsOf(cbc.CodeMap{
							cef.ExtKeyVATEX:          "VATEX-XX-99",
							untdid.ExtKeyTaxCategory: en16931.TaxCategoryExempt,
						}),
					},
				},
			},
		}
		ad.Normalizer(inv)
		assert.Empty(t, inv.Tax.Notes)
	})

	t.Run("multiple lines with different categories get separate notes", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Lines = []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Exempt item",
					Price: num.NewAmount(100, 0),
				},
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Key:      tax.KeyExempt,
						Ext: tax.ExtensionsOf(cbc.CodeMap{
							cef.ExtKeyVATEX:          "VATEX-SA-29",
							untdid.ExtKeyTaxCategory: en16931.TaxCategoryExempt,
						}),
					},
				},
			},
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Export item",
					Price: num.NewAmount(200, 0),
				},
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Key:      tax.KeyZero,
						Ext: tax.ExtensionsOf(cbc.CodeMap{
							cef.ExtKeyVATEX:          "VATEX-SA-32",
							untdid.ExtKeyTaxCategory: en16931.TaxCategoryZero,
						}),
					},
				},
			},
		}
		ad.Normalizer(inv)
		require.Len(t, inv.Tax.Notes, 2)
		assert.Equal(t, en16931.TaxCategoryExempt, inv.Tax.Notes[0].Ext.Get(untdid.ExtKeyTaxCategory))
		assert.Equal(t, en16931.TaxCategoryZero, inv.Tax.Notes[1].Ext.Get(untdid.ExtKeyTaxCategory))
	})
}

func TestBillDiscountRules(t *testing.T) {
	t.Run("discount with taxes is valid", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Discounts = []*bill.Discount{
			{
				Reason: "Loyalty discount",
				Amount: num.MakeAmount(50, 0),
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Rate:     tax.RateGeneral,
					},
				},
			},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("discount without taxes fails", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Discounts = []*bill.Discount{
			{
				Reason: "Loyalty discount",
				Amount: num.MakeAmount(50, 0),
			},
		}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "taxes are required (BR-32)")
	})
}
