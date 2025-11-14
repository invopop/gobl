package choruspro_test

import (
	"testing"

	"github.com/invopop/gobl/addons/fr/choruspro"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validInvoice() *bill.Invoice {
	return &bill.Invoice{
		Series: "TEST",
		Code:   "0002",
		Regime: tax.WithRegime("FR"),
		Addons: tax.WithAddons(choruspro.V1),
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "DE",
				Code:    "44732829320",
			},
			Addresses: []*org.Address{
				{
					Country: "DE",
				},
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "FR",
				Code:    "391838042",
			},
			Identities: []*org.Identity{
				{
					Type: fr.IdentityTypeSIRET,
					Code: "39183804212345",
				},
			},
			Addresses: []*org.Address{
				{
					Country: "FR",
				},
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
						Rate:     "general",
					},
				},
			},
		},
		Payment: &bill.PaymentDetails{
			Terms: &pay.Terms{
				Notes: "Please pay in 10 days",
			},
		},
	}
}

func TestValidateInvoice(t *testing.T) {
	addon := tax.AddonForKey(choruspro.V1)

	t.Run("valid invoice", func(t *testing.T) {
		inv := validInvoice()
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("missing customer", func(t *testing.T) {
		inv := validInvoice()
		require.NoError(t, inv.Calculate())
		inv.Customer = nil
		err := addon.Validator(inv)
		assert.NoError(t, err) // Customer validation only runs if customer exists
	})

	t.Run("customer with SIRET with wrong scheme", func(t *testing.T) {
		inv := validInvoice()
		require.NoError(t, inv.Calculate())
		inv.Customer.Ext = tax.Extensions{
			choruspro.ExtKeyScheme: "2",
		}
		err := addon.Validator(inv)
		assert.ErrorContains(t, err, "invalid value")
	})

	t.Run("customer with SIRET with correct scheme", func(t *testing.T) {
		inv := validInvoice()
		require.NoError(t, inv.Calculate())
		inv.Customer.Identities = []*org.Identity{
			{
				Type: fr.IdentityTypeSIRET,
				Code: "12345678901234",
			},
		}
		inv.Customer.Ext = tax.Extensions{
			choruspro.ExtKeyScheme: "1",
		}
		err := addon.Validator(inv)
		assert.NoError(t, err)
	})

	t.Run("customer with no SIRET", func(t *testing.T) {
		inv := validInvoice()
		require.NoError(t, inv.Calculate())
		inv.Customer.Identities = nil
		err := addon.Validator(inv)
		assert.ErrorContains(t, err, "identities: cannot be blank")
	})

	t.Run("missing framework extension", func(t *testing.T) {
		inv := validInvoice()
		require.NoError(t, inv.Calculate())
		// Remove framework extension
		if inv.Tax != nil && inv.Tax.Ext != nil {
			delete(inv.Tax.Ext, choruspro.ExtKeyFramework)
		}
		err := addon.Validator(inv)
		assert.ErrorContains(t, err, "required")
	})

	t.Run("framework A2 with unpaid invoice", func(t *testing.T) {
		inv := validInvoice()
		require.NoError(t, inv.Calculate())

		// Set framework to A2 (already paid)
		if inv.Tax == nil {
			inv.Tax = &bill.Tax{}
		}
		if inv.Tax.Ext == nil {
			inv.Tax.Ext = make(tax.Extensions)
		}
		inv.Tax.Ext[choruspro.ExtKeyFramework] = "A2"

		err := addon.Validator(inv)
		assert.ErrorContains(t, err, "totals: must be paid in full for framework 'A2'.")
	})

	t.Run("framework A2 with paid invoice", func(t *testing.T) {
		inv := validInvoice()
		inv.Tax = &bill.Tax{
			Ext: tax.Extensions{
				choruspro.ExtKeyFramework: "A2",
			},
		}
		inv.Payment = &bill.PaymentDetails{
			Advances: []*pay.Advance{
				{
					Percent: num.NewPercentage(100, 2),
				},
			},
		}
		require.NoError(t, inv.Calculate())
		err := addon.Validator(inv)
		assert.NoError(t, err)
		assert.Equal(t, cbc.Code("A2"), inv.Tax.Ext.Get(choruspro.ExtKeyFramework))
	})

}

func TestNormalizeInvoice(t *testing.T) {
	t.Run("sets default framework A1", func(t *testing.T) {
		inv := validInvoice()
		inv.Tax = nil

		require.NoError(t, inv.Calculate())

		assert.NotNil(t, inv.Tax)
		assert.NotNil(t, inv.Tax.Ext)
		assert.Equal(t, cbc.Code("A1"), inv.Tax.Ext.Get(choruspro.ExtKeyFramework))
	})

	t.Run("preserves existing framework", func(t *testing.T) {
		inv := validInvoice()
		inv.Tax = &bill.Tax{
			Ext: tax.Extensions{
				choruspro.ExtKeyFramework: "A3",
			},
		}

		require.NoError(t, inv.Calculate())

		assert.Equal(t, cbc.Code("A3"), inv.Tax.Ext.Get(choruspro.ExtKeyFramework))
	})

	addon := tax.AddonForKey(choruspro.V1)

	t.Run("with nil tax", func(t *testing.T) {
		inv := validInvoice()
		inv.Tax = nil

		addon.Normalizer(inv)

		assert.NotNil(t, inv.Tax)
		assert.NotNil(t, inv.Tax.Ext)
		assert.Equal(t, cbc.Code("A1"), inv.Tax.Ext.Get(choruspro.ExtKeyFramework))
	})

	t.Run("with nil invoice", func(t *testing.T) {
		var inv *bill.Invoice
		addon.Normalizer(inv)
		assert.Nil(t, inv)
	})
}

func TestNormalizeBillLines(t *testing.T) {
	addon := tax.AddonForKey(choruspro.V1)
	t.Run("Remove decimal places", func(t *testing.T) {

		line := &bill.Line{
			Quantity: num.MakeAmount(1000000, 6),
		}

		addon.Normalizer(line)

		assert.Equal(t, "1.0000", line.Quantity.String())
	})

	t.Run("Remove trailing decimals", func(t *testing.T) {
		line := &bill.Line{
			Quantity: num.MakeAmount(1530000, 6),
		}

		addon.Normalizer(line)

		assert.Equal(t, "1.5300", line.Quantity.String())
	})

	t.Run("Remove trailing decimals", func(t *testing.T) {
		line := &bill.Line{
			Quantity: num.MakeAmount(13342423, 6),
		}

		addon.Normalizer(line)

		assert.Equal(t, "13.3424", line.Quantity.String())
	})

	t.Run("Empty quantity", func(t *testing.T) {
		line := &bill.Line{
			Quantity: num.Amount{},
		}

		addon.Normalizer(line)

		assert.Equal(t, "0", line.Quantity.String())
	})

	t.Run("with nil line", func(t *testing.T) {
		var line *bill.Line
		addon.Normalizer(line)
		assert.Nil(t, line)
	})
}
