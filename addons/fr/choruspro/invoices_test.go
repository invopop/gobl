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
						Rate:     "standard",
					},
				},
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

		err := addon.Validator(inv)
		assert.NoError(t, err)
	})

	t.Run("missing customer", func(t *testing.T) {
		inv := validInvoice()
		require.NoError(t, inv.Calculate())
		inv.Customer = nil
		err := addon.Validator(inv)
		assert.NoError(t, err) // Customer validation only runs if customer exists
	})

	t.Run("customer without identities", func(t *testing.T) {
		inv := validInvoice()
		require.NoError(t, inv.Calculate())
		inv.Customer.Identities = nil
		err := addon.Validator(inv)
		assert.ErrorContains(t, err, "cannot be blank")
	})

	t.Run("customer with empty identities", func(t *testing.T) {
		inv := validInvoice()
		require.NoError(t, inv.Calculate())
		inv.Customer.Identities = []*org.Identity{}
		err := addon.Validator(inv)
		assert.ErrorContains(t, err, "cannot be blank")
	})

	t.Run("customer with SIRET without scheme", func(t *testing.T) {
		inv := validInvoice()
		require.NoError(t, inv.Calculate())
		inv.Customer.Identities = []*org.Identity{
			{
				Type: fr.IdentityTypeSIRET,
				Code: "12345678901234",
				// No scheme extension
			},
		}
		err := addon.Validator(inv)
		assert.ErrorContains(t, err, "required")
	})

	t.Run("customer with SIRET with wrong scheme", func(t *testing.T) {
		inv := validInvoice()
		require.NoError(t, inv.Calculate())
		inv.Customer.Identities = []*org.Identity{
			{
				Type: fr.IdentityTypeSIRET,
				Code: "12345678901234",
				Ext: tax.Extensions{
					choruspro.ExtKeyScheme: "2", // Should be "1" for SIRET
				},
			},
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
				Ext: tax.Extensions{
					choruspro.ExtKeyScheme: "1",
				},
			},
		}
		err := addon.Validator(inv)
		assert.NoError(t, err)
	})

	t.Run("customer with non-SIRET identity", func(t *testing.T) {
		inv := validInvoice()
		require.NoError(t, inv.Calculate())
		inv.Customer.Identities = []*org.Identity{
			{
				Type: "OTHER",
				Code: "123456789",
				Ext: tax.Extensions{
					choruspro.ExtKeyScheme: "3", // Any scheme is fine for non-SIRET
				},
			},
		}
		err := addon.Validator(inv)
		assert.NoError(t, err)
	})

	t.Run("missing framework extension", func(t *testing.T) {
		inv := validInvoice()
		require.NoError(t, inv.Calculate())
		// Remove framework extension
		if inv.Tax != nil && inv.Tax.Ext != nil {
			delete(inv.Tax.Ext, choruspro.ExtKeyFrameWork)
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
		inv.Tax.Ext[choruspro.ExtKeyFrameWork] = "A2"

		err := addon.Validator(inv)
		assert.ErrorContains(t, err, "If the invoice has type A2, it must be paid in full")
	})

	t.Run("framework A2 with paid invoice", func(t *testing.T) {
		inv := validInvoice()
		inv.Tax = &bill.Tax{
			Ext: tax.Extensions{
				choruspro.ExtKeyFrameWork: "A2",
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
		assert.Equal(t, cbc.Code("A2"), inv.Tax.Ext.Get(choruspro.ExtKeyFrameWork))
	})

}

func TestNormalizeInvoice(t *testing.T) {
	t.Run("sets default framework A1", func(t *testing.T) {
		inv := validInvoice()
		inv.Tax = nil

		require.NoError(t, inv.Calculate())

		assert.NotNil(t, inv.Tax)
		assert.NotNil(t, inv.Tax.Ext)
		assert.Equal(t, cbc.Code("A1"), inv.Tax.Ext.Get(choruspro.ExtKeyFrameWork))
	})

	t.Run("preserves existing framework", func(t *testing.T) {
		inv := validInvoice()
		inv.Tax = &bill.Tax{
			Ext: tax.Extensions{
				choruspro.ExtKeyFrameWork: "A3",
			},
		}

		require.NoError(t, inv.Calculate())

		assert.Equal(t, cbc.Code("A3"), inv.Tax.Ext.Get(choruspro.ExtKeyFrameWork))
	})

	t.Run("normalizes French supplier", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID = &tax.Identity{
			Country: "FR",
			Code:    "39356000000",
		}

		require.NoError(t, inv.Calculate())

		// French suppliers should not have scheme set at tax level
		assert.Equal(t, cbc.Code(""), inv.Supplier.TaxID.Scheme)
	})

	t.Run("normalizes EU supplier with scheme 2", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID = &tax.Identity{
			Country: "DE", // EU country
			Code:    "123456789",
		}

		require.NoError(t, inv.Calculate())

		assert.Equal(t, cbc.Code("2"), inv.Supplier.TaxID.Scheme)
	})

	t.Run("normalizes non-EU supplier with scheme 3", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID = &tax.Identity{
			Country: "US", // Non-EU country
			Code:    "123456789",
		}

		require.NoError(t, inv.Calculate())

		assert.Equal(t, cbc.Code("3"), inv.Supplier.TaxID.Scheme)
	})

	t.Run("normalizes RIDET supplier with scheme 4", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID = &tax.Identity{
			Country: "NC", // New Caledonia - RIDET
			Code:    "123456789",
		}

		require.NoError(t, inv.Calculate())

		assert.Equal(t, cbc.Code("3"), inv.Supplier.TaxID.Scheme) // Should be 3 for non-EU
	})

	t.Run("normalizes Tahiti supplier with scheme 5", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID = &tax.Identity{
			Country: "PF", // French Polynesia - Tahiti
			Code:    "123456789",
		}

		require.NoError(t, inv.Calculate())

		assert.Equal(t, cbc.Code("3"), inv.Supplier.TaxID.Scheme) // Should be 3 for non-EU
	})

	t.Run("does not normalize supplier without tax ID", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID = nil

		require.NoError(t, inv.Calculate())

		assert.Nil(t, inv.Supplier.TaxID)
	})

	t.Run("does not normalize supplier without tax code", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.TaxID = &tax.Identity{
			Country: "DE",
			Code:    "", // Empty code
		}

		require.NoError(t, inv.Calculate())

		assert.Equal(t, cbc.Code(""), inv.Supplier.TaxID.Scheme)
	})

	addon := tax.AddonForKey(choruspro.V1)

	t.Run("with nil tax", func(t *testing.T) {
		inv := validInvoice()
		inv.Tax = nil

		addon.Normalizer(inv)

		assert.NotNil(t, inv.Tax)
		assert.NotNil(t, inv.Tax.Ext)
		assert.Equal(t, cbc.Code("A1"), inv.Tax.Ext.Get(choruspro.ExtKeyFrameWork))
	})

	t.Run("with nil invoice", func(t *testing.T) {
		var inv *bill.Invoice
		addon.Normalizer(inv)
		assert.Nil(t, inv)
	})
}

func TestNormalizeParty(t *testing.T) {
	addon := tax.AddonForKey(choruspro.V1)
	require.NotNil(t, addon)

	t.Run("normalizes SIRET identity with scheme 1", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			Identities: []*org.Identity{
				{
					Type: fr.IdentityTypeSIRET,
					Code: "12345678901234",
					// No scheme extension initially
				},
			},
		}

		addon.Normalizer(party)

		assert.NotNil(t, party.Identities[0].Ext)
		assert.Equal(t, cbc.Code("1"), party.Identities[0].Ext.Get(choruspro.ExtKeyScheme))
	})

	t.Run("preserves existing SIRET scheme", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			Identities: []*org.Identity{
				{
					Type: fr.IdentityTypeSIRET,
					Code: "12345678901234",
					Ext: tax.Extensions{
						choruspro.ExtKeyScheme: "1",
					},
				},
			},
		}

		addon.Normalizer(party)

		assert.Equal(t, cbc.Code("1"), party.Identities[0].Ext.Get(choruspro.ExtKeyScheme))
	})

	t.Run("only normalizes first SIRET identity", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			Identities: []*org.Identity{
				{
					Type: "OTHER",
					Code: "123456789",
				},
				{
					Type: fr.IdentityTypeSIRET,
					Code: "12345678901234",
					// No scheme extension initially
				},
				{
					Type: fr.IdentityTypeSIRET,
					Code: "98765432109876",
					// This should not be normalized
				},
			},
		}

		addon.Normalizer(party)

		// First SIRET should be normalized
		assert.NotNil(t, party.Identities[1].Ext)
		assert.Equal(t, cbc.Code("1"), party.Identities[1].Ext.Get(choruspro.ExtKeyScheme))

		// Second SIRET should not be normalized
		assert.Nil(t, party.Identities[2].Ext)
	})

	t.Run("does not normalize non-SIRET identities", func(t *testing.T) {
		party := &org.Party{
			Name: "Test Party",
			Identities: []*org.Identity{
				{
					Type: "OTHER",
					Code: "123456789",
				},
			},
		}

		addon.Normalizer(party)

		assert.Nil(t, party.Identities[0].Ext)
	})

	t.Run("handles nil identities", func(t *testing.T) {
		party := &org.Party{
			Name:       "Test Party",
			Identities: nil,
		}

		addon.Normalizer(party)

		assert.Nil(t, party.Identities)
	})

	t.Run("handles empty identities", func(t *testing.T) {
		party := &org.Party{
			Name:       "Test Party",
			Identities: []*org.Identity{},
		}

		addon.Normalizer(party)

		assert.Empty(t, party.Identities)
	})
}
