package sa_test

import (
	"testing"

	"github.com/invopop/gobl/org"
	_ "github.com/invopop/gobl/regimes/sa"
	"github.com/invopop/gobl/rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCustomerTaxIDOrIdentity(t *testing.T) {
	t.Run("customer with VAT and no identity is valid", func(t *testing.T) {
		inv := validInvoice()
		inv.Customer.Identities = nil
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("customer without VAT but with valid identity is valid", func(t *testing.T) {
		inv := validInvoice()
		inv.Customer.TaxID = nil
		inv.Customer.Identities = []*org.Identity{
			{
				Type: "TIN",
				Code: "123456789012345",
			},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("customer without VAT and without identity is invalid", func(t *testing.T) {
		inv := validInvoice()
		inv.Customer.TaxID = nil
		inv.Customer.Identities = nil
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "customer must have a valid tax ID code or identity")
	})

	t.Run("customer without VAT and with empty identity code is invalid", func(t *testing.T) {
		inv := validInvoice()
		inv.Customer.TaxID = nil
		inv.Customer.Identities = []*org.Identity{
			{
				Type: "TIN",
				Code: "",
			},
		}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "identity code must be provided")
	})
}

func TestSupplierSingleIdentity(t *testing.T) {
	t.Run("supplier with one identity is valid", func(t *testing.T) {
		inv := validInvoice()
		inv.Supplier.Identities = []*org.Identity{
			{
				Type: "CRN",
				Code: "1234567890",
			},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})
}

func TestCustomerSingleIdentity(t *testing.T) {
	t.Run("customer with one identity is valid", func(t *testing.T) {
		inv := validInvoice()
		inv.Customer.Identities = []*org.Identity{
			{
				Type: "TIN",
				Code: "123456789012345",
			},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})
}

func TestCustomerIdentityTIN(t *testing.T) {
	t.Run("valid TIN identity", func(t *testing.T) {
		inv := validInvoice()
		inv.Customer.Identities = []*org.Identity{
			{
				Type: "TIN",
				Code: "123456789012345",
			},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("invalid TIN identity too short", func(t *testing.T) {
		inv := validInvoice()
		inv.Customer.Identities = []*org.Identity{
			{
				Type: "TIN",
				Code: "1234567890",
			},
		}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "identity code for type TIN must be valid")
	})

	t.Run("invalid TIN identity with letters", func(t *testing.T) {
		inv := validInvoice()
		inv.Customer.Identities = []*org.Identity{
			{
				Type: "TIN",
				Code: "12345678901234A",
			},
		}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "identity code for type TIN must be valid")
	})
}
