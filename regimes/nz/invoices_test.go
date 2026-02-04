package nz_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validInvoice() *bill.Invoice {
	return &bill.Invoice{
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Code:    "49091850",
				Country: "NZ",
			},
			Name: "Kiwi Software Solutions Ltd",
			Addresses: []*org.Address{
				{
					Street:  "Queen Street",
					Code:    "1010",
					Country: l10n.NZ.ISO(),
				},
			},
		},
		Customer: &org.Party{
			Name: "Wellington Tech Enterprises",
			Addresses: []*org.Address{
				{
					Street:  "Lambton Quay",
					Code:    "6011",
					Country: l10n.NZ.ISO(),
				},
			},
		},
		Code:     "0001",
		Currency: "NZD",
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(10, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(100, 0),
				},
				Taxes: tax.Set{
					{
						Category: tax.CategoryGST,
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
	assert.Equal(t, "1150.00", inv.Totals.TotalWithTax.String())
	require.NoError(t, inv.Validate())
}

func TestBoundaryExactly200NoSupplierTaxCode(t *testing.T) {
	inv := validInvoice()
	inv.Lines[0].Quantity = num.MakeAmount(1, 0)
	inv.Lines[0].Item.Price = num.NewAmount(200, 0)
	inv.Lines[0].Taxes = tax.Set{{Category: tax.CategoryGST, Rate: tax.RateZero, Key: tax.KeyZero}}
	inv.Supplier.TaxID = &tax.Identity{Country: "NZ"}
	require.NoError(t, inv.Calculate())
	assert.Equal(t, "200.00", inv.Totals.TotalWithTax.String())
	require.NoError(t, inv.Validate(), "≤$200 should not require supplier GST number")
}

func TestBoundaryJustOver200RequiresSupplierTaxCode(t *testing.T) {
	inv := validInvoice()
	inv.Lines[0].Quantity = num.MakeAmount(1, 0)
	inv.Lines[0].Item.Price = num.NewAmount(20001, 2) // $200.01
	inv.Lines[0].Taxes = tax.Set{{Category: tax.CategoryGST, Rate: tax.RateZero, Key: tax.KeyZero}}
	inv.Supplier.TaxID = &tax.Identity{Country: "NZ"}
	require.NoError(t, inv.Calculate())
	assert.Equal(t, "200.01", inv.Totals.TotalWithTax.String())
	err := inv.Validate()
	require.Error(t, err, ">$200 should require supplier GST number")
	assert.Contains(t, err.Error(), "supplier")
}

func TestBoundaryExactly1000NoCustomer(t *testing.T) {
	inv := validInvoice()
	inv.Lines[0].Quantity = num.MakeAmount(1, 0)
	inv.Lines[0].Item.Price = num.NewAmount(1000, 0)
	inv.Lines[0].Taxes = tax.Set{{Category: tax.CategoryGST, Rate: tax.RateZero, Key: tax.KeyZero}}
	inv.Customer = nil
	require.NoError(t, inv.Calculate())
	assert.Equal(t, "1000.00", inv.Totals.TotalWithTax.String())
	require.NoError(t, inv.Validate(), "≤$1,000 should not require customer")
}

func TestBoundaryJustOver1000RequiresCustomer(t *testing.T) {
	inv := validInvoice()
	inv.Lines[0].Quantity = num.MakeAmount(1, 0)
	inv.Lines[0].Item.Price = num.NewAmount(100001, 2) // $1,000.01
	inv.Lines[0].Taxes = tax.Set{{Category: tax.CategoryGST, Rate: tax.RateZero, Key: tax.KeyZero}}
	inv.Customer = nil
	require.NoError(t, inv.Calculate())
	assert.Equal(t, "1000.01", inv.Totals.TotalWithTax.String())
	err := inv.Validate()
	require.Error(t, err, ">$1,000 should require customer")
	assert.Contains(t, err.Error(), "customer")
}

func TestLowValueInvoiceNoSupplierTaxCode(t *testing.T) {
	inv := validInvoice()
	inv.Lines[0].Quantity = num.MakeAmount(1, 0)
	inv.Lines[0].Item.Price = num.NewAmount(100, 0)
	inv.Supplier.TaxID = &tax.Identity{Country: "NZ"}
	require.NoError(t, inv.Calculate())
	assert.Equal(t, "115.00", inv.Totals.TotalWithTax.String())
	require.NoError(t, inv.Validate())
}

func TestLowValueInvoiceNoCustomer(t *testing.T) {
	inv := validInvoice()
	inv.Lines[0].Quantity = num.MakeAmount(1, 0)
	inv.Lines[0].Item.Price = num.NewAmount(100, 0)
	inv.Customer = nil
	require.NoError(t, inv.Calculate())
	assert.Equal(t, "115.00", inv.Totals.TotalWithTax.String())
	require.NoError(t, inv.Validate())
}

func TestMidValueInvoiceRequiresSupplierTaxCode(t *testing.T) {
	inv := validInvoice()
	inv.Lines[0].Quantity = num.MakeAmount(1, 0)
	inv.Lines[0].Item.Price = num.NewAmount(500, 0)
	inv.Supplier.TaxID = &tax.Identity{Country: "NZ"}
	require.NoError(t, inv.Calculate())
	// Total = $500 + $75 GST = $575 (mid-value)
	assert.Equal(t, "575.00", inv.Totals.TotalWithTax.String())
	err := inv.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "supplier")
}

func TestMidValueInvoiceNoCustomer(t *testing.T) {
	inv := validInvoice()
	inv.Lines[0].Quantity = num.MakeAmount(1, 0)
	inv.Lines[0].Item.Price = num.NewAmount(500, 0)
	inv.Customer = nil
	require.NoError(t, inv.Calculate())
	assert.Equal(t, "575.00", inv.Totals.TotalWithTax.String())
	require.NoError(t, inv.Validate())
}

func TestMidValueInvoiceMissingSupplierTaxCode(t *testing.T) {
	inv := validInvoice()
	inv.Lines[0].Quantity = num.MakeAmount(1, 0)
	inv.Lines[0].Item.Price = num.NewAmount(500, 0)
	inv.Supplier.TaxID = &tax.Identity{
		Country: "NZ",
	}
	require.NoError(t, inv.Calculate())
	err := inv.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "supplier")
}

func TestHighValueInvoiceRequiresCustomer(t *testing.T) {
	inv := validInvoice()
	inv.Customer = nil
	require.NoError(t, inv.Calculate())
	assert.Equal(t, "1150.00", inv.Totals.TotalWithTax.String())
	err := inv.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "customer")
}

func TestHighValueInvoiceRequiresCustomerName(t *testing.T) {
	inv := validInvoice()
	inv.Customer.Name = ""
	require.NoError(t, inv.Calculate())
	err := inv.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "name")
}

func TestHighValueInvoiceCustomerWithEmail(t *testing.T) {
	inv := validInvoice()
	inv.Customer = &org.Party{
		Name: "Wellington Tech Enterprises",
		Emails: []*org.Email{
			{Address: "info@wellingtontech.nz"},
		},
	}
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())
}

func TestHighValueInvoiceCustomerWithAddress(t *testing.T) {
	inv := validInvoice()
	inv.Customer = &org.Party{
		Name: "Wellington Tech Enterprises",
		Addresses: []*org.Address{
			{
				Street:  "Lambton Quay",
				Code:    "6011",
				Country: l10n.NZ.ISO(),
			},
		},
	}
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())
}

func TestHighValueInvoiceCustomerNoIdentifier(t *testing.T) {
	inv := validInvoice()
	inv.Customer = &org.Party{
		Name: "Wellington Tech Enterprises",
	}
	require.NoError(t, inv.Calculate())
	err := inv.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "addresses")
}

func TestHighValueInvoiceCustomerWithNZBN(t *testing.T) {
	inv := validInvoice()
	inv.Customer = &org.Party{
		Name: "Wellington Tech Enterprises",
		Identities: []*org.Identity{
			{
				Key:  org.IdentityKeyGLN,
				Code: "9429041234563",
			},
		},
	}
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())
}

func TestHighValueInvoiceCustomerWithPhone(t *testing.T) {
	inv := validInvoice()
	inv.Customer = &org.Party{
		Name: "Wellington Tech Enterprises",
		Telephones: []*org.Telephone{
			{Number: "+64 4 123 4567"},
		},
	}
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())
}

func TestHighValueInvoiceCustomerWithWebsite(t *testing.T) {
	inv := validInvoice()
	inv.Customer = &org.Party{
		Name: "Wellington Tech Enterprises",
		Websites: []*org.Website{
			{URL: "https://wellingtontech.nz"},
		},
	}
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())
}

func TestHighValueInvoiceRequiresSupplierTaxID(t *testing.T) {
	inv := validInvoice()
	inv.Supplier.TaxID = &tax.Identity{Country: "NZ"}
	require.NoError(t, inv.Calculate())
	assert.Equal(t, "1150.00", inv.Totals.TotalWithTax.String())
	err := inv.Validate()
	require.Error(t, err, ">$1,000 should also require supplier GST number")
	assert.Contains(t, err.Error(), "supplier")
}

func TestExportInvoiceValid(t *testing.T) {
	inv := validInvoice()
	inv.SetTags(tax.TagExport)
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())
}

func TestExportInvoiceRequiresCustomer(t *testing.T) {
	inv := validInvoice()
	inv.SetTags(tax.TagExport)
	inv.Customer = nil
	// Use low value to isolate tag validation from threshold
	inv.Lines[0].Quantity = num.MakeAmount(1, 0)
	inv.Lines[0].Item.Price = num.NewAmount(50, 0)
	require.NoError(t, inv.Calculate())
	err := inv.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "customer")
}

func TestExportInvoiceRequiresCustomerAddress(t *testing.T) {
	inv := validInvoice()
	inv.SetTags(tax.TagExport)
	inv.Customer = &org.Party{
		Name: "Export Customer",
	}
	inv.Lines[0].Quantity = num.MakeAmount(1, 0)
	inv.Lines[0].Item.Price = num.NewAmount(50, 0)
	require.NoError(t, inv.Calculate())
	err := inv.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "addresses")
}

func TestExportInvoiceRequiresCustomerName(t *testing.T) {
	inv := validInvoice()
	inv.SetTags(tax.TagExport)
	inv.Customer = &org.Party{
		Addresses: []*org.Address{
			{
				Street:  "123 Export St",
				Code:    "1010",
				Country: l10n.NZ.ISO(),
			},
		},
	}
	inv.Lines[0].Quantity = num.MakeAmount(1, 0)
	inv.Lines[0].Item.Price = num.NewAmount(50, 0)
	require.NoError(t, inv.Calculate())
	err := inv.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "name")
}

