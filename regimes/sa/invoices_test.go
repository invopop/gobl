package sa_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/sa"
	"github.com/invopop/gobl/tax"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validInvoice() *bill.Invoice {
	return &bill.Invoice{
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Code:    "300075588700003",
				Country: "SA",
			},
			Name: "Test Supplier",
			Identities: []*org.Identity{
				{
					Type: sa.IdentityTypeCRN,
					Code: "1234567890",
				},
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Code:    "300075588700003",
				Country: "SA",
			},
		},
		Code:     "INV-001",
		Currency: "SAR",
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(100, 0),
				},
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Rate:     tax.RateGeneral,
					},
				},
			},
		},
	}
}

func TestValidInvoice(t *testing.T) {
	inv := validInvoice()
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())
}

func TestValidInvoiceWithCRN(t *testing.T) {
	inv := validInvoice()
	inv.Supplier.Identities = []*org.Identity{
		{
			Type: sa.IdentityTypeCRN,
			Code: "1234567890",
		},
	}
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())
}

func TestMissingSupplierName(t *testing.T) {
	inv := validInvoice()
	inv.Supplier.Name = ""
	require.Error(t, inv.Validate())
}

func TestStandardInvoiceRequiresCustomer(t *testing.T) {
	inv := validInvoice()
	inv.Customer = nil
	require.NoError(t, inv.Calculate())
	err := inv.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "customer")
}

func TestStandardInvoiceRequiresCustomerName(t *testing.T) {
	inv := validInvoice()
	inv.Customer = &org.Party{
		Name: "",
		TaxID: &tax.Identity{
			Code:    "300075588700003",
			Country: "SA",
		},
	}
	require.NoError(t, inv.Calculate())
	err := inv.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "name")
}

func TestStandardInvoiceWithCustomerTaxID(t *testing.T) {
	inv := validInvoice()
	inv.Customer = &org.Party{
		Name: "Buyer Co.",
		TaxID: &tax.Identity{
			Code:    "300075588700003",
			Country: "SA",
		},
	}
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())
}

func TestStandardInvoiceWithCustomerIdentity(t *testing.T) {
	inv := validInvoice()
	inv.Customer = &org.Party{
		Name: "Buyer Co.",
		Identities: []*org.Identity{
			{
				Type: sa.IdentityTypeNAT,
				Code: "1234567890",
			},
		},
	}
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())
}

func TestStandardInvoiceCustomerRequiresIdentification(t *testing.T) {
	inv := validInvoice()
	inv.Customer = &org.Party{
		Name: "Buyer Co.",
	}
	require.NoError(t, inv.Calculate())
	err := inv.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "identities")
}

func TestSimplifiedInvoiceWithoutCustomer(t *testing.T) {
	inv := validInvoice()
	inv.Tags = tax.WithTags(tax.TagSimplified)
	inv.Customer = nil
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())
}

func TestSimplifiedInvoiceWithCustomer(t *testing.T) {
	inv := validInvoice()
	inv.Tags = tax.WithTags(tax.TagSimplified)
	inv.Customer = &org.Party{
		Name: "Walk-in Customer",
	}
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())
}

func TestInvoiceWithZeroRatedLine(t *testing.T) {
	inv := validInvoice()
	inv.Lines = []*bill.Line{
		{
			Quantity: num.MakeAmount(1, 0),
			Item: &org.Item{
				Name:  "Exported Service",
				Price: num.NewAmount(1000, 0),
			},
			Taxes: tax.Set{
				{
					Category: tax.CategoryVAT,
					Key:      tax.KeyZero,
				},
			},
		},
	}
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())

	// KeyZero should auto-assign 0% via global VAT keys
	line := inv.Lines[0]
	require.NotNil(t, line.Taxes[0].Percent)
	assert.True(t, line.Taxes[0].Percent.IsZero())
}

func TestInvoiceWithExemptLine(t *testing.T) {
	inv := validInvoice()
	inv.Lines = []*bill.Line{
		{
			Quantity: num.MakeAmount(1, 0),
			Item: &org.Item{
				Name:  "Healthcare Service",
				Price: num.NewAmount(500, 0),
			},
			Taxes: tax.Set{
				{
					Category: tax.CategoryVAT,
					Key:      tax.KeyExempt,
				},
			},
		},
	}
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())

	// KeyExempt should result in nil percent (no tax)
	line := inv.Lines[0]
	assert.Nil(t, line.Taxes[0].Percent)
}
