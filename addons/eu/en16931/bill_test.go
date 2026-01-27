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
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoiceValidation(t *testing.T) {
	ad := tax.AddonForKey(en16931.V2017)
	t.Run("valid invoice", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "380", inv.Tax.Ext[untdid.ExtKeyDocumentType].String())
		err := inv.Validate()
		assert.NoError(t, err)
	})
	t.Run("missing tax", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Type = bill.InvoiceTypeOther
		require.NoError(t, inv.Calculate())
		inv.Tax = nil // not sure why this would happen...
		err := ad.Validator(inv)
		assert.ErrorContains(t, err, "tax: cannot be blank")
	})
	t.Run("missing tax document type", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Type = bill.InvoiceTypeOther
		inv.Tax = &bill.Tax{PricesInclude: "VAT"}
		require.NoError(t, inv.Calculate())
		err := ad.Validator(inv)
		assert.ErrorContains(t, err, "tax: (ext: (untdid-document-type: required.).)")
	})
	t.Run("supplier and customer with addresses", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})
	t.Run("supplier with no address", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.Addresses = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "addresses: cannot be blank")
	})
	t.Run("customer with no address", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.Addresses = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "addresses: cannot be blank")
	})
	t.Run("nil customer", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})
	t.Run("credit note", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Type = bill.InvoiceTypeCreditNote
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
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
		assert.Equal(t, "67", l.Ext[untdid.ExtKeyAllowance].String())
	})
	t.Run("without key", func(t *testing.T) {
		l := &bill.LineDiscount{
			Reason: "Product sample",
			Amount: num.MakeAmount(100, 2),
		}
		ad.Normalizer(l)
		assert.Nil(t, l.Ext)
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
		assert.Equal(t, "67", l.Ext[untdid.ExtKeyAllowance].String())
	})
	t.Run("without key", func(t *testing.T) {
		l := &bill.Discount{
			Reason: "Product sample",
			Amount: num.MakeAmount(100, 2),
		}
		ad.Normalizer(l)
		assert.Nil(t, l.Ext)
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
		assert.Equal(t, "AAE", l.Ext[untdid.ExtKeyCharge].String())
	})
	t.Run("without key", func(t *testing.T) {
		l := &bill.LineCharge{
			Reason: "Additional costs",
			Amount: num.MakeAmount(3000, 2),
		}
		ad.Normalizer(l)
		assert.Nil(t, l.Ext)
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
		assert.Equal(t, "AAE", l.Ext[untdid.ExtKeyCharge].String())
	})
	t.Run("without key", func(t *testing.T) {
		l := &bill.Charge{
			Reason: "Additional costs",
			Amount: num.MakeAmount(3000, 2),
		}
		ad.Normalizer(l)
		assert.Nil(t, l.Ext)
	})
}

func TestValidateBillDiscount(t *testing.T) {
	ad := tax.AddonForKey(en16931.V2017)
	t.Run("with reason", func(t *testing.T) {
		l := &bill.Discount{
			Reason: "Product sample",
			Amount: num.MakeAmount(100, 2),
			Taxes: tax.Set{
				{
					Category: tax.CategoryVAT,
					Rate:     "standard",
				},
			},
		}
		err := ad.Validator(l)
		assert.NoError(t, err)
	})

	t.Run("with extension", func(t *testing.T) {
		l := &bill.Discount{
			Ext: tax.Extensions{
				untdid.ExtKeyAllowance: "67",
			},
			Amount: num.MakeAmount(100, 2),
			Taxes: tax.Set{
				{
					Category: tax.CategoryVAT,
					Rate:     "standard",
				},
			},
		}
		err := ad.Validator(l)
		assert.NoError(t, err)
	})

	t.Run("without reason or extension", func(t *testing.T) {
		l := &bill.Discount{
			Amount: num.MakeAmount(100, 2),
			Taxes: tax.Set{
				{
					Category: tax.CategoryVAT,
					Rate:     "standard",
				},
			},
		}
		err := ad.Validator(l)
		assert.ErrorContains(t, err, "either a reason or an allowance type extension is required")
	})

	t.Run("with reason and extension", func(t *testing.T) {
		l := &bill.Discount{
			Reason: "Product sample",
			Ext: tax.Extensions{
				untdid.ExtKeyAllowance: "67",
			},
			Amount: num.MakeAmount(100, 2),
			Taxes: tax.Set{
				{
					Category: tax.CategoryVAT,
					Rate:     "standard",
				},
			},
		}
		err := ad.Validator(l)
		assert.NoError(t, err)
	})

	t.Run("without taxes (BR-32)", func(t *testing.T) {
		l := &bill.Discount{
			Reason: "Product sample",
			Amount: num.MakeAmount(100, 2),
		}
		err := ad.Validator(l)
		assert.ErrorContains(t, err, "taxes are required (BR-32)")
	})
}

func TestValidateBillCharge(t *testing.T) {
	ad := tax.AddonForKey(en16931.V2017)
	t.Run("with reason", func(t *testing.T) {
		l := &bill.Charge{
			Reason: "Product sample",
			Amount: num.MakeAmount(100, 2),
		}
		err := ad.Validator(l)
		assert.NoError(t, err)
	})

	t.Run("with extension", func(t *testing.T) {
		l := &bill.Charge{
			Ext: tax.Extensions{
				untdid.ExtKeyCharge: "AAE",
			},
			Amount: num.MakeAmount(100, 2),
		}
		err := ad.Validator(l)
		assert.NoError(t, err)
	})

	t.Run("without reason or extension", func(t *testing.T) {
		l := &bill.Charge{
			Amount: num.MakeAmount(100, 2),
		}
		err := ad.Validator(l)
		assert.ErrorContains(t, err, "either a reason or a charge type extension is required")
	})

	t.Run("with reason and extension", func(t *testing.T) {
		l := &bill.Charge{
			Reason: "Product sample",
			Ext: tax.Extensions{
				untdid.ExtKeyCharge: "AAE",
			},
			Amount: num.MakeAmount(100, 2),
		}
		err := ad.Validator(l)
		assert.NoError(t, err)
	})
}

func TestValidateBillLine(t *testing.T) {
	ad := tax.AddonForKey(en16931.V2017)
	t.Run("Discount with reason", func(t *testing.T) {
		l := &bill.Line{
			Discounts: []*bill.LineDiscount{
				{
					Reason: "Product sample",
					Amount: num.MakeAmount(100, 2),
				},
			},
		}
		err := ad.Validator(l)
		assert.NoError(t, err)
	})

	t.Run("Discount with extension", func(t *testing.T) {
		l := &bill.Line{
			Discounts: []*bill.LineDiscount{
				{
					Ext: tax.Extensions{
						untdid.ExtKeyAllowance: "67",
					},
					Amount: num.MakeAmount(100, 2),
				},
			},
		}
		err := ad.Validator(l)
		assert.NoError(t, err)
	})

	t.Run("Discount without reason or extension", func(t *testing.T) {
		l := &bill.Line{
			Discounts: []*bill.LineDiscount{
				{
					Amount: num.MakeAmount(100, 2),
				},
			},
		}
		err := ad.Validator(l)
		assert.ErrorContains(t, err, "either a reason or an allowance type extension is required")
	})

	t.Run("Discount with reason and extension", func(t *testing.T) {
		l := &bill.Line{
			Discounts: []*bill.LineDiscount{
				{
					Reason: "Product sample",
					Ext: tax.Extensions{
						untdid.ExtKeyAllowance: "67",
					},
					Amount: num.MakeAmount(100, 2),
				},
			},
		}
		err := ad.Validator(l)
		assert.NoError(t, err)
	})

	t.Run("Charge with reason", func(t *testing.T) {
		l := &bill.Line{
			Charges: []*bill.LineCharge{
				{
					Reason: "Product sample",
					Amount: num.MakeAmount(100, 2),
				},
			},
		}
		err := ad.Validator(l)
		assert.NoError(t, err)
	})

	t.Run("Charge with extension", func(t *testing.T) {
		l := &bill.Line{
			Charges: []*bill.LineCharge{
				{
					Ext: tax.Extensions{
						untdid.ExtKeyCharge: "AAE",
					},
					Amount: num.MakeAmount(100, 2),
				},
			},
		}
		err := ad.Validator(l)
		assert.NoError(t, err)
	})

	t.Run("Charge without reason or extension", func(t *testing.T) {
		l := &bill.Line{
			Charges: []*bill.LineCharge{
				{
					Amount: num.MakeAmount(100, 2),
				},
			},
		}
		err := ad.Validator(l)
		assert.ErrorContains(t, err, "either a reason or a charge type extension is required")
	})

	t.Run("Charge with reason and extension", func(t *testing.T) {
		l := &bill.Line{
			Charges: []*bill.LineCharge{
				{
					Reason: "Product sample",
					Ext: tax.Extensions{
						untdid.ExtKeyCharge: "AAE",
					},
					Amount: num.MakeAmount(100, 2),
				},
			},
		}
		err := ad.Validator(l)
		assert.NoError(t, err)
	})

	t.Run("Line with nil charge and discount", func(t *testing.T) {
		l := &bill.Line{
			Discounts: []*bill.LineDiscount{nil},
			Charges:   []*bill.LineCharge{nil},
		}
		err := ad.Validator(l)
		assert.NoError(t, err)
	})
}

func TestValidateBillPayment(t *testing.T) {
	ad := tax.AddonForKey(en16931.V2017)
	t.Run("with terms", func(t *testing.T) {
		p := &bill.PaymentDetails{
			Terms: &pay.Terms{
				DueDates: []*pay.DueDate{
					{
						Date:   cal.NewDate(2025, time.January, 1),
						Amount: num.MakeAmount(1000, 2),
					},
				},
			},
		}
		err := ad.Validator(p)
		assert.NoError(t, err)
	})

	t.Run("without terms", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Payment.Terms = nil
		err := inv.Calculate()
		require.NoError(t, err)
		err = inv.Validate()
		assert.ErrorContains(t, err, " payment terms are required when amount is due (BR-CO-25)")
	})

	t.Run("with nil payment details", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Payment = nil
		err := inv.Calculate()
		require.NoError(t, err)
		err = inv.Validate()
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
		err = inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("with payment details but no terms when due", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Payment = &bill.PaymentDetails{} // payment details exist but no terms
		err := inv.Calculate()
		require.NoError(t, err)
		err = inv.Validate()
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
		err = inv.Validate()
		assert.NoError(t, err) // Should pass because no amount is due
	})

}
