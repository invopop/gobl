package en16931_test

import (
	"testing"
	"time"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/addons/eu/en16931"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoiceValidation(t *testing.T) {
	t.Run("valid invoice", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "380", inv.Tax.Ext.Get(untdid.ExtKeyDocumentType).String())
		err := rules.Validate(inv)
		assert.NoError(t, err)
	})
	t.Run("missing tax", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Type = bill.InvoiceTypeOther
		require.NoError(t, inv.Calculate())
		inv.Tax = nil // not sure why this would happen...
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "tax details are required")
	})
	t.Run("missing tax document type", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Type = bill.InvoiceTypeOther
		inv.Tax = &bill.Tax{PricesInclude: "VAT"}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "document type extension is required")
	})
	t.Run("supplier and customer with addresses", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.NoError(t, err)
	})
	t.Run("supplier with no address", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.Addresses = nil
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "addresses are required")
	})
	t.Run("customer with no address", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.Addresses = nil
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "addresses are required")
	})
	t.Run("nil customer", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer = nil
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.NoError(t, err)
	})
	t.Run("credit note", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Type = bill.InvoiceTypeCreditNote
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.NoError(t, err)
	})
}

func TestExemptionNoteValidation(t *testing.T) {
	t.Run("exempt with matching note", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Lines = []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item:     &org.Item{Name: "Exempt item", Price: num.NewAmount(100, 2)},
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Key:      tax.KeyExempt,
						Ext: tax.ExtensionsOf(tax.ExtMap{
							"cef-vatex": "VATEX-EU-132",
						}),
					},
				},
			},
		}
		inv.Tax = &bill.Tax{
			Notes: []*tax.Note{
				{Category: tax.CategoryVAT, Key: "exempt", Text: "Exempt under Article 132"},
			},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("exempt without note or vatex", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Lines = []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item:     &org.Item{Name: "Exempt item", Price: num.NewAmount(100, 2)},
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Key:      tax.KeyExempt,
					},
				},
			},
		}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "exempt tax categories require either a VATEX code or an exemption note")
	})

	t.Run("exempt with vatex code no note needed", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Lines = []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item:     &org.Item{Name: "Exempt item", Price: num.NewAmount(100, 2)},
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Key:      tax.KeyExempt,
						Ext: tax.ExtensionsOf(tax.ExtMap{
							"cef-vatex": "VATEX-EU-132",
						}),
					},
				},
			},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("nil note in notes slice", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Lines = []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item:     &org.Item{Name: "Exempt item", Price: num.NewAmount(100, 2)},
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Key:      tax.KeyExempt,
						Ext: tax.ExtensionsOf(tax.ExtMap{
							"cef-vatex": "VATEX-EU-132",
						}),
					},
				},
			},
		}
		var nilNote *tax.Note
		inv.Tax = &bill.Tax{
			Notes: []*tax.Note{
				nilNote,
				{Category: tax.CategoryVAT, Key: "exempt", Text: "Exempt under Article 132"},
			},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("nil tax note normalization", func(t *testing.T) {
		ad := tax.AddonForKey(en16931.V2017)
		var n *tax.Note
		assert.NotPanics(t, func() {
			ad.Normalizer(n)
		})
	})

	t.Run("non-VAT note skips normalization", func(t *testing.T) {
		ad := tax.AddonForKey(en16931.V2017)
		n := &tax.Note{
			Category: "IGIC",
			Key:      "exempt",
			Text:     "Some IGIC exemption",
		}
		ad.Normalizer(n)
		assert.False(t, n.Ext.Has(untdid.ExtKeyTaxCategory))
	})

	t.Run("note normalization derives key from ext", func(t *testing.T) {
		ad := tax.AddonForKey(en16931.V2017)
		n := &tax.Note{
			Category: tax.CategoryVAT,
			Text:     "Exempt under Article 132",
			Ext: tax.ExtensionsOf(tax.ExtMap{
				untdid.ExtKeyTaxCategory: "E",
			}),
		}
		ad.Normalizer(n)
		assert.Equal(t, tax.KeyExempt, n.Key)
	})

	t.Run("note normalization derives key for reverse charge", func(t *testing.T) {
		ad := tax.AddonForKey(en16931.V2017)
		n := &tax.Note{
			Category: tax.CategoryVAT,
			Text:     "Reverse charge applies",
			Ext: tax.ExtensionsOf(tax.ExtMap{
				untdid.ExtKeyTaxCategory: "AE",
			}),
		}
		ad.Normalizer(n)
		assert.Equal(t, tax.KeyReverseCharge, n.Key)
	})

	t.Run("note normalization does not override existing key", func(t *testing.T) {
		ad := tax.AddonForKey(en16931.V2017)
		n := &tax.Note{
			Category: tax.CategoryVAT,
			Key:      tax.KeyExempt,
			Text:     "Exempt under Article 132",
			Ext: tax.ExtensionsOf(tax.ExtMap{
				untdid.ExtKeyTaxCategory: "AE",
			}),
		}
		ad.Normalizer(n)
		assert.Equal(t, tax.KeyExempt, n.Key, "should not override existing key")
	})

	t.Run("note normalization adds tax category", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Lines = []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item:     &org.Item{Name: "Exempt item", Price: num.NewAmount(100, 2)},
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Key:      tax.KeyExempt,
					},
				},
			},
		}
		inv.Tax = &bill.Tax{
			Notes: []*tax.Note{
				{Category: tax.CategoryVAT, Key: "exempt", Text: "Exempt under Article 132"},
			},
		}
		require.NoError(t, inv.Calculate())
		// After calculation, the tax note should have the UNTDID tax category set
		require.Len(t, inv.Tax.Notes, 1)
		assert.Equal(t, "E", inv.Tax.Notes[0].Ext.Get(untdid.ExtKeyTaxCategory).String())
	})

	t.Run("standard rate does not need note", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})
}

func testInvoiceStandard(t *testing.T) *bill.Invoice {
	t.Helper()
	inv := &bill.Invoice{
		Regime:   tax.WithRegime("DE"),
		Addons:   tax.WithAddons(en16931.V2017),
		Type:     "standard",
		Currency: "EUR",
		Series:   "2024",
		Code:     "1000",
		Supplier: &org.Party{
			Name: "Cursor AG",
			TaxID: &tax.Identity{
				Country: "DE",
				Code:    "505898911",
			},
			People: []*org.Person{
				{
					Name: &org.Name{
						Given:   "Peter",
						Surname: "Cursorstone",
					},
				},
			},
			Addresses: []*org.Address{
				{
					Street:   "Dietmar-Hopp-Allee",
					Locality: "Walldorf",
					Code:     "69190",
					Country:  "DE",
				},
			},
		},
		Customer: &org.Party{
			Name: "Sample Consumer",
			TaxID: &tax.Identity{
				Country: "DE",
				Code:    "449674701",
			},
			People: []*org.Person{
				{
					Name: &org.Name{
						Given:   "Max",
						Surname: "Musterman",
					},
				},
			},
			Addresses: []*org.Address{
				{
					Street:   "Werner-Heisenberg-Allee",
					Locality: "München",
					Code:     "80939",
					Country:  "DE",
				},
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(10, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(10000, 2),
					Unit:  "item",
				},
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Rate:     "general",
					},
				},
			},
		},
		Payment: &bill.PaymentDetails{
			Terms: &pay.Terms{
				DueDates: []*pay.DueDate{
					{
						Date:   cal.NewDate(2025, time.January, 1),
						Amount: num.MakeAmount(1000, 2),
					},
				},
			},
		},
	}
	return inv
}

func TestNormalizeBillLineDiscount(t *testing.T) {
	ad := tax.AddonForKey(en16931.V2017)
	t.Run("with key", func(t *testing.T) {
		l := &bill.LineDiscount{
			Key:    "sample",
			Reason: "Product sample",
			Amount: num.MakeAmount(100, 2),
		}
		ad.Normalizer(l)
		assert.Equal(t, "67", l.Ext.Get(untdid.ExtKeyAllowance).String())
	})
	t.Run("without key", func(t *testing.T) {
		l := &bill.LineDiscount{
			Reason: "Product sample",
			Amount: num.MakeAmount(100, 2),
		}
		ad.Normalizer(l)
		assert.True(t, l.Ext.IsZero())
	})
}

func TestNormalizeBillDiscount(t *testing.T) {
	ad := tax.AddonForKey(en16931.V2017)
	t.Run("with key", func(t *testing.T) {
		l := &bill.Discount{
			Key:    "sample",
			Reason: "Product sample",
			Amount: num.MakeAmount(100, 2),
		}
		ad.Normalizer(l)
		assert.Equal(t, "67", l.Ext.Get(untdid.ExtKeyAllowance).String())
	})
	t.Run("without key", func(t *testing.T) {
		l := &bill.Discount{
			Reason: "Product sample",
			Amount: num.MakeAmount(100, 2),
		}
		ad.Normalizer(l)
		assert.True(t, l.Ext.IsZero())
	})
}

func TestNormalizeBillLineCharge(t *testing.T) {
	ad := tax.AddonForKey(en16931.V2017)
	t.Run("with key", func(t *testing.T) {
		l := &bill.LineCharge{
			Key:    "outlay",
			Reason: "Notary costs",
			Amount: num.MakeAmount(1000, 2),
		}
		ad.Normalizer(l)
		assert.Equal(t, "AAE", l.Ext.Get(untdid.ExtKeyCharge).String())
	})
	t.Run("without key", func(t *testing.T) {
		l := &bill.LineCharge{
			Reason: "Additional costs",
			Amount: num.MakeAmount(3000, 2),
		}
		ad.Normalizer(l)
		assert.True(t, l.Ext.IsZero())
	})
}

func TestNormalizeBillCharge(t *testing.T) {
	ad := tax.AddonForKey(en16931.V2017)
	t.Run("with key", func(t *testing.T) {
		l := &bill.Charge{
			Key:    "outlay",
			Reason: "Notary costs",
			Amount: num.MakeAmount(1000, 2),
		}
		ad.Normalizer(l)
		assert.Equal(t, "AAE", l.Ext.Get(untdid.ExtKeyCharge).String())
	})
	t.Run("without key", func(t *testing.T) {
		l := &bill.Charge{
			Reason: "Additional costs",
			Amount: num.MakeAmount(3000, 2),
		}
		ad.Normalizer(l)
		assert.True(t, l.Ext.IsZero())
	})
}

func TestValidateBillDiscount(t *testing.T) {
	t.Run("with reason", func(t *testing.T) {
		l := &bill.Discount{
			Reason: "Product sample",
			Amount: num.MakeAmount(100, 2),
			Taxes: tax.Set{
				{
					Category: tax.CategoryVAT,
					Rate:     "standard",
					Ext: tax.ExtensionsOf(tax.ExtMap{
						untdid.ExtKeyTaxCategory: en16931.TaxCategoryStandard,
					}),
				},
			},
		}
		err := rules.Validate(l, tax.AddonContext(en16931.V2017))
		assert.NoError(t, err)
	})

	t.Run("with extension", func(t *testing.T) {
		l := &bill.Discount{
			Ext: tax.ExtensionsOf(tax.ExtMap{
				untdid.ExtKeyAllowance: "67",
			}),
			Amount: num.MakeAmount(100, 2),
			Taxes: tax.Set{
				{
					Category: tax.CategoryVAT,
					Rate:     "standard",
					Ext: tax.ExtensionsOf(tax.ExtMap{
						untdid.ExtKeyTaxCategory: en16931.TaxCategoryStandard,
					}),
				},
			},
		}
		err := rules.Validate(l, tax.AddonContext(en16931.V2017))
		assert.NoError(t, err)
	})

	t.Run("without reason or extension", func(t *testing.T) {
		l := &bill.Discount{
			Amount: num.MakeAmount(100, 2),
			Taxes: tax.Set{
				{
					Category: tax.CategoryVAT,
					Rate:     "standard",
					Ext: tax.ExtensionsOf(tax.ExtMap{
						untdid.ExtKeyTaxCategory: en16931.TaxCategoryStandard,
					}),
				},
			},
		}
		err := rules.Validate(l, tax.AddonContext(en16931.V2017))
		assert.ErrorContains(t, err, "either a reason or an allowance type extension is required")
	})

	t.Run("with reason and extension", func(t *testing.T) {
		l := &bill.Discount{
			Reason: "Product sample",
			Ext: tax.ExtensionsOf(tax.ExtMap{
				untdid.ExtKeyAllowance: "67",
			}),
			Amount: num.MakeAmount(100, 2),
			Taxes: tax.Set{
				{
					Category: tax.CategoryVAT,
					Rate:     "standard",
					Ext: tax.ExtensionsOf(tax.ExtMap{
						untdid.ExtKeyTaxCategory: en16931.TaxCategoryStandard,
					}),
				},
			},
		}
		err := rules.Validate(l, tax.AddonContext(en16931.V2017))
		assert.NoError(t, err)
	})

	t.Run("without taxes (BR-32)", func(t *testing.T) {
		l := &bill.Discount{
			Reason: "Product sample",
			Amount: num.MakeAmount(100, 2),
		}
		err := rules.Validate(l, tax.AddonContext(en16931.V2017))
		assert.ErrorContains(t, err, "taxes are required (BR-32)")
	})
}

func TestValidateBillCharge(t *testing.T) {
	t.Run("with reason", func(t *testing.T) {
		l := &bill.Charge{
			Reason: "Product sample",
			Amount: num.MakeAmount(100, 2),
		}
		err := rules.Validate(l, tax.AddonContext(en16931.V2017))
		assert.NoError(t, err)
	})

	t.Run("with extension", func(t *testing.T) {
		l := &bill.Charge{
			Ext: tax.ExtensionsOf(tax.ExtMap{
				untdid.ExtKeyCharge: "AAE",
			}),
			Amount: num.MakeAmount(100, 2),
		}
		err := rules.Validate(l, tax.AddonContext(en16931.V2017))
		assert.NoError(t, err)
	})

	t.Run("without reason or extension", func(t *testing.T) {
		l := &bill.Charge{
			Amount: num.MakeAmount(100, 2),
		}
		err := rules.Validate(l, tax.AddonContext(en16931.V2017))
		assert.ErrorContains(t, err, "either a reason or a charge type extension is required")
	})

	t.Run("with reason and extension", func(t *testing.T) {
		l := &bill.Charge{
			Reason: "Product sample",
			Ext: tax.ExtensionsOf(tax.ExtMap{
				untdid.ExtKeyCharge: "AAE",
			}),
			Amount: num.MakeAmount(100, 2),
		}
		err := rules.Validate(l, tax.AddonContext(en16931.V2017))
		assert.NoError(t, err)
	})
}

func TestValidateBillLineDiscount(t *testing.T) {
	t.Run("Discount with reason", func(t *testing.T) {
		d := &bill.LineDiscount{
			Reason: "Product sample",
			Amount: num.MakeAmount(100, 2),
		}
		err := rules.Validate(d, tax.AddonContext(en16931.V2017))
		assert.NoError(t, err)
	})

	t.Run("Discount with extension", func(t *testing.T) {
		d := &bill.LineDiscount{
			Ext: tax.ExtensionsOf(tax.ExtMap{
				untdid.ExtKeyAllowance: "67",
			}),
			Amount: num.MakeAmount(100, 2),
		}
		err := rules.Validate(d, tax.AddonContext(en16931.V2017))
		assert.NoError(t, err)
	})

	t.Run("Discount without reason or extension", func(t *testing.T) {
		d := &bill.LineDiscount{
			Amount: num.MakeAmount(100, 2),
		}
		err := rules.Validate(d, tax.AddonContext(en16931.V2017))
		assert.ErrorContains(t, err, "either a reason or an allowance type extension is required")
	})

	t.Run("Discount with reason and extension", func(t *testing.T) {
		d := &bill.LineDiscount{
			Reason: "Product sample",
			Ext: tax.ExtensionsOf(tax.ExtMap{
				untdid.ExtKeyAllowance: "67",
			}),
			Amount: num.MakeAmount(100, 2),
		}
		err := rules.Validate(d, tax.AddonContext(en16931.V2017))
		assert.NoError(t, err)
	})

	t.Run("nil discount", func(t *testing.T) {
		var d *bill.LineDiscount
		err := rules.Validate(d, tax.AddonContext(en16931.V2017))
		assert.NoError(t, err)
	})
}

func TestValidateBillLineCharge(t *testing.T) {
	t.Run("Charge with reason", func(t *testing.T) {
		c := &bill.LineCharge{
			Reason: "Product sample",
			Amount: num.MakeAmount(100, 2),
		}
		err := rules.Validate(c, tax.AddonContext(en16931.V2017))
		assert.NoError(t, err)
	})

	t.Run("Charge with extension", func(t *testing.T) {
		c := &bill.LineCharge{
			Ext: tax.ExtensionsOf(tax.ExtMap{
				untdid.ExtKeyCharge: "AAE",
			}),
			Amount: num.MakeAmount(100, 2),
		}
		err := rules.Validate(c, tax.AddonContext(en16931.V2017))
		assert.NoError(t, err)
	})

	t.Run("Charge without reason or extension", func(t *testing.T) {
		c := &bill.LineCharge{
			Amount: num.MakeAmount(100, 2),
		}
		err := rules.Validate(c, tax.AddonContext(en16931.V2017))
		assert.ErrorContains(t, err, "either a reason or a charge type extension is required")
	})

	t.Run("Charge with reason and extension", func(t *testing.T) {
		c := &bill.LineCharge{
			Reason: "Product sample",
			Ext: tax.ExtensionsOf(tax.ExtMap{
				untdid.ExtKeyCharge: "AAE",
			}),
			Amount: num.MakeAmount(100, 2),
		}
		err := rules.Validate(c, tax.AddonContext(en16931.V2017))
		assert.NoError(t, err)
	})

	t.Run("nil charge", func(t *testing.T) {
		var c *bill.LineCharge
		err := rules.Validate(c, tax.AddonContext(en16931.V2017))
		assert.NoError(t, err)
	})
}

func TestValidateBillPayment(t *testing.T) {
	t.Run("with terms", func(t *testing.T) {
		p := &pay.Terms{
			DueDates: []*pay.DueDate{
				{
					Date:   cal.NewDate(2025, time.January, 1),
					Amount: num.MakeAmount(1000, 2),
				},
			},
		}
		err := rules.Validate(p, tax.AddonContext(en16931.V2017))
		assert.NoError(t, err)
	})

	t.Run("without terms", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Payment.Terms = nil
		err := inv.Calculate()
		require.NoError(t, err)
		err = rules.Validate(inv)
		assert.ErrorContains(t, err, " payment terms are required when amount is due (BR-CO-25)")
	})

	t.Run("with nil payment details", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Payment = nil
		err := inv.Calculate()
		require.NoError(t, err)
		err = rules.Validate(inv)
		assert.ErrorContains(t, err, "payment details are required when amount is due")
	})

	t.Run("with due amount zero", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		advances := []*pay.Advance{
			{
				Percent:     num.NewPercentage(100, 2),
				Description: "Advance payment",
			},
		}
		inv.Payment.Advances = advances
		inv.Payment.Terms = nil

		err := inv.Calculate()
		require.NoError(t, err)
		err = rules.Validate(inv)
		assert.NoError(t, err)
	})

	t.Run("with payment details but no terms when due", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Payment = &bill.PaymentDetails{} // payment details exist but no terms
		err := inv.Calculate()
		require.NoError(t, err)
		err = rules.Validate(inv)
		assert.ErrorContains(t, err, " payment terms are required when amount is due (BR-CO-25)")
	})

	t.Run("no payment details and no amount due", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		// Add advance payment to make due amount zero
		advances := []*pay.Advance{
			{
				Percent:     num.NewPercentage(100, 2),
				Description: "Full advance payment",
			},
		}
		inv.Payment = &bill.PaymentDetails{Advances: advances}
		err := inv.Calculate()
		require.NoError(t, err)
		// Remove payment details after calculation
		inv.Payment = nil
		err = rules.Validate(inv)
		assert.NoError(t, err) // Should pass because no amount is due
	})

}
