package nz_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/nz"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func parsePrice(s string) *num.Amount {
	amt, err := num.AmountFromString(s)
	if err != nil {
		panic(err)
	}
	return &amt
}

func validInvoice(total string) *bill.Invoice {
	return &bill.Invoice{
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Code:    "123456789",
				Country: "NZ",
			},
			Name: "Test Supplier Ltd",
			Addresses: []*org.Address{
				{
					Street:   "123 Test Street",
					Code:     "1010",
					Locality: "Auckland",
					Country:  l10n.NZ.ISO(),
				},
			},
		},
		Customer: &org.Party{
			Name: "Test Customer Ltd",
			Addresses: []*org.Address{
				{
					Street:   "456 Customer Road",
					Code:     "2010",
					Locality: "Wellington",
					Country:  l10n.NZ.ISO(),
				},
			},
		},
		Code:     "INV-001",
		Currency: "NZD",
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: parsePrice(total),
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

func TestValidInvoiceUnder200(t *testing.T) {
	inv := validInvoice("150.00")
	require.NoError(t, inv.Calculate())
	inv.Supplier.TaxID = nil
	err := nz.Validate(inv)
	assert.NoError(t, err, "invoices under $200 should not require GST number")
}

func TestValidInvoice200To1000(t *testing.T) {
	inv := validInvoice("500.00")
	require.NoError(t, inv.Calculate())
	err := nz.Validate(inv)
	assert.NoError(t, err, "valid invoice with GST number should pass")
}

func TestInvoice200WithoutGSTNumber(t *testing.T) {
	inv := validInvoice("200.00")
	require.NoError(t, inv.Calculate())
	inv.Supplier.TaxID = nil
	err := nz.Validate(inv)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "supplier must have GST number for invoices > $200")
}

func TestInvoice500WithoutGSTNumber(t *testing.T) {
	inv := validInvoice("500.00")
	require.NoError(t, inv.Calculate())
	inv.Supplier.TaxID = nil
	err := nz.Validate(inv)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "supplier must have GST number for invoices > $200")
}

func TestValidInvoiceOver1000(t *testing.T) {
	inv := validInvoice("1500.00")
	require.NoError(t, inv.Calculate())
	err := nz.Validate(inv)
	assert.NoError(t, err, "valid invoice over $1000 with customer details should pass")
}

func TestInvoiceOver1000WithoutCustomerName(t *testing.T) {
	inv := validInvoice("1500.00")
	inv.Customer.Name = ""
	require.NoError(t, inv.Calculate())
	err := nz.Validate(inv)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "customer name required for invoices > $1,000")
}

func TestInvoiceOver1000WithoutCustomerIdentifier(t *testing.T) {
	inv := validInvoice("1500.00")
	inv.Customer.Addresses = nil
	require.NoError(t, inv.Calculate())
	err := nz.Validate(inv)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "customer must have at least one identifier")
}

func TestInvoiceOver1000WithEmail(t *testing.T) {
	inv := validInvoice("1500.00")
	inv.Customer.Addresses = nil
	inv.Customer.Emails = []*org.Email{
		{Address: "customer@example.com"},
	}
	require.NoError(t, inv.Calculate())
	err := nz.Validate(inv)
	assert.NoError(t, err, "email should satisfy identifier requirement")
}

func TestInvoiceOver1000WithPhone(t *testing.T) {
	inv := validInvoice("1500.00")
	inv.Customer.Addresses = nil
	inv.Customer.Telephones = []*org.Telephone{
		{Number: "+64 21 123 4567"},
	}
	require.NoError(t, inv.Calculate())
	err := nz.Validate(inv)
	assert.NoError(t, err, "phone should satisfy identifier requirement")
}

func TestInvoiceOver1000WithTaxID(t *testing.T) {
	inv := validInvoice("1500.00")
	inv.Customer.Addresses = nil
	inv.Customer.TaxID = &tax.Identity{
		Code:    "987654321",
		Country: "NZ",
	}
	require.NoError(t, inv.Calculate())
	err := nz.Validate(inv)
	assert.NoError(t, err, "tax ID should satisfy identifier requirement")
}

func TestInvoiceOver1000WithIdentity(t *testing.T) {
	inv := validInvoice("1500.00")
	inv.Customer.Addresses = nil
	inv.Customer.Identities = []*org.Identity{
		{
			Key:  cbc.Key("nzbn"),
			Code: "9429000000000",
		},
	}
	require.NoError(t, inv.Calculate())
	err := nz.Validate(inv)
	assert.NoError(t, err, "NZBN identity should satisfy identifier requirement")
}

func TestInvoiceExactly1000(t *testing.T) {
	// $869.57 * 1.15 = $1000.01 (just over threshold)
	// $869.56 * 1.15 = $999.99 (just under threshold)
	inv := validInvoice("869.56")
	inv.Customer.Addresses = nil
	require.NoError(t, inv.Calculate())
	err := nz.Validate(inv)
	assert.NoError(t, err, "invoice with total under $1000 should not require customer identifier")
}

func TestInvoiceExactly1001(t *testing.T) {
	// $869.57 * 1.15 = $1000.01 (just over threshold)
	inv := validInvoice("869.57")
	inv.Customer.Addresses = nil
	require.NoError(t, inv.Calculate())
	err := nz.Validate(inv)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "customer must have at least one identifier")
}

func TestInvoiceExactly200(t *testing.T) {
	inv := validInvoice("200.00")
	require.NoError(t, inv.Calculate())
	err := nz.Validate(inv)
	assert.NoError(t, err, "invoice at exactly $200 with GST number should pass")
}

func TestInvoiceExactly199(t *testing.T) {
	// $173.90 * 1.15 = $199.985 which rounds to $199.99 (under $200)
	inv := validInvoice("173.90")
	require.NoError(t, inv.Calculate())
	inv.Supplier.TaxID = nil
	err := nz.Validate(inv)
	assert.NoError(t, err, "invoice under $200 should not require GST number")
}

func TestNilInvoice(t *testing.T) {
	var inv *bill.Invoice
	err := nz.Validate(inv)
	assert.NoError(t, err, "nil invoice should not cause error")
}

func TestInvoiceWithNilSupplier(t *testing.T) {
	inv := validInvoice("500.00")
	require.NoError(t, inv.Calculate())
	inv.Supplier = nil
	err := nz.Validate(inv)
	assert.NoError(t, err, "nil supplier should not cause panic")
}

func TestInvoiceWithNilCustomer(t *testing.T) {
	inv := validInvoice("1500.00")
	inv.Customer = nil
	require.NoError(t, inv.Calculate())
	err := nz.Validate(inv)
	assert.NoError(t, err, "nil customer should not cause panic")
}

func TestInvoiceWithEmptyTaxID(t *testing.T) {
	inv := validInvoice("500.00")
	inv.Supplier.TaxID = &tax.Identity{
		Code:    "",
		Country: "NZ",
	}
	require.NoError(t, inv.Calculate())
	err := nz.Validate(inv)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be blank")
}

func TestValidateInvoiceSupplierTotalsNil(t *testing.T) {
	inv := validInvoice("500.00")
	inv.Totals = nil
	err := nz.Validate(inv)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invoice totals must be calculated before NZ regime validation")
}

func TestValidateInvoiceCustomerTotalsNil(t *testing.T) {
	inv := validInvoice("1500.00")
	inv.Totals = nil
	err := nz.Validate(inv)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invoice totals must be calculated before NZ regime validation")
}
