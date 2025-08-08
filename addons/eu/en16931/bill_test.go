package en16931_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/addons/eu/en16931"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
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
