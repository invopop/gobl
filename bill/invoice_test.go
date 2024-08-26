package bill_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	_ "github.com/invopop/gobl" // load regions
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/internal"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoiceCodeRegexp(t *testing.T) {
	tests := []struct {
		code string
		ok   bool
	}{
		// Good
		{"1", true},
		{"A", true},
		{"123", true},
		{"123TEST", true},
		{"123-TEST", true},
		{"FR F-01/37", true},
		{"MultiCase", true},
		{"F.01_21", true},
		// Bad
		{"F101-", false},
		{" 123 ", false},
		{"F--01", false},
		{"\n", false},
		{"FOO\n", false},
	}
	for _, ts := range tests {
		t.Run(ts.code, func(t *testing.T) {
			assert.Equal(t, ts.ok, bill.InvoiceCodeRegexp.MatchString(ts.code))
		})
	}
}

func TestInvoiceRegimeCurrency(t *testing.T) {
	lines := []*bill.Line{
		{
			Quantity: num.MakeAmount(1, 0),
			Item: &org.Item{
				Name:  "Test Item",
				Price: num.MakeAmount(10, 0),
			},
			Taxes: tax.Set{
				{
					Category: "VAT",
					Rate:     tax.RateStandard,
				},
			},
		},
	}
	i := baseInvoice(t, lines...)

	require.NoError(t, i.Calculate())

	assert.Equal(t, currency.EUR, i.Currency, "should set currency automatically")
	assert.Equal(t, "10.00", i.Lines[0].Item.Price.String(), "should update price precision")
	i.Lines[0].Item.Price = num.MakeAmount(10000, 3)
	require.NoError(t, i.Calculate())
	assert.Equal(t, "10.000", i.Lines[0].Item.Price.String(), "should not update price precision")
}

func TestInvoiceRegimeCurrencyCLP(t *testing.T) {
	lines := []*bill.Line{
		{
			Quantity: num.MakeAmount(1, 0),
			Item: &org.Item{
				Name:  "Test Item",
				Price: num.MakeAmount(10, 0),
			},
		},
	}
	i := baseInvoice(t, lines...)
	i.Currency = currency.CLP
	require.NoError(t, i.Calculate())
	assert.Equal(t, currency.CLP, i.Currency, "should honor currency")
	assert.Equal(t, "10", i.Lines[0].Item.Price.String(), "should not update price precision")
}

func TestInvoiceInvalidCurrency(t *testing.T) {
	lines := []*bill.Line{
		{
			Quantity: num.MakeAmount(1, 0),
			Item: &org.Item{
				Name:  "Test Item",
				Price: num.MakeAmount(10, 0),
			},
		},
	}
	i := baseInvoice(t, lines...)
	i.Currency = "MX"
	require.NoError(t, i.Calculate())
	assert.Equal(t, currency.EUR, i.Currency, "should correct currency")
}

func TestInvoiceRegimeCurrencyWithDiscounts(t *testing.T) {
	lines := []*bill.Line{
		{
			Quantity: num.MakeAmount(1, 0),
			Item: &org.Item{
				Name:  "Test Item",
				Price: num.MakeAmount(10, 0),
			},
		},
	}
	i := baseInvoice(t, lines...)
	i.Lines[0].Discounts = []*bill.LineDiscount{
		{
			Reason: "Testing",
			Amount: num.MakeAmount(10, 0),
		},
	}
	i.Lines[0].Charges = []*bill.LineCharge{
		{
			Reason: "Testing",
			Amount: num.MakeAmount(20, 0),
		},
	}
	require.NoError(t, i.Calculate())

	assert.Equal(t, "10.00", i.Lines[0].Discounts[0].Amount.String(), "should update discount precision")
	assert.Equal(t, "20.00", i.Lines[0].Charges[0].Amount.String(), "should update charges precision")
}

func TestInvoiceCurrencyValidation(t *testing.T) {
	lines := []*bill.Line{
		{
			Quantity: num.MakeAmount(1, 0),
			Item: &org.Item{
				Name:  "Test Item",
				Price: num.MakeAmount(10, 0),
			},
		},
	}
	inv := baseInvoice(t, lines...)
	inv.Currency = currency.USD
	require.NoError(t, inv.Calculate())

	assert.ErrorContains(t, inv.Validate(), "currency: no exchange rate defined for 'USD' to 'EUR'")

	inv.ExchangeRates = []*currency.ExchangeRate{
		{
			From:   currency.USD,
			To:     currency.EUR,
			Amount: num.MakeAmount(875967, 6),
		},
	}
	assert.NoError(t, inv.Validate())
}

func TestInvoiceAutoSetIssueDate(t *testing.T) {
	lines := []*bill.Line{
		{
			Quantity: num.MakeAmount(1, 0),
			Item: &org.Item{
				Name:  "Test Item",
				Price: num.MakeAmount(10, 0),
			},
		},
	}
	i := baseInvoice(t, lines...)
	i.IssueDate = cal.Date{} // zero

	require.NoError(t, i.Calculate())

	loc, err := time.LoadLocation("Europe/Madrid")
	require.NoError(t, err)
	dn := cal.TodayIn(loc)
	assert.Equal(t, dn.String(), i.IssueDate.String(), "should set issue date to today")
}

func TestRemoveIncludedTax(t *testing.T) {
	lines := []*bill.Line{
		{
			Quantity: num.MakeAmount(1, 0),
			Item: &org.Item{
				Name:  "Test Item",
				Price: num.MakeAmount(100000, 2),
			},
			Taxes: tax.Set{
				{
					Category: "VAT",
					Percent:  num.NewPercentage(21, 2),
				},
			},
			Discounts: []*bill.LineDiscount{
				{
					Reason:  "Testing",
					Percent: num.NewPercentage(10, 2),
				},
			},
		},
	}
	i := baseInvoice(t, lines...)

	require.NoError(t, i.Calculate())

	i2, err := i.RemoveIncludedTaxes()
	require.NoError(t, err)

	assert.Equal(t, "1000.00", i.Lines[0].Item.Price.String())

	assert.Empty(t, i2.Tax.PricesInclude)
	l0 := i2.Lines[0]
	assert.Equal(t, "826.4463", l0.Item.Price.String())
	assert.Equal(t, "826.4463", l0.Sum.String())
	assert.Equal(t, "826.4463", l0.Sum.String())
	assert.Equal(t, "82.6446", l0.Discounts[0].Amount.String())
	assert.Equal(t, "743.8017", l0.Total.String())

	assert.Equal(t, "743.80", i2.Totals.Sum.String())
	assert.Equal(t, i.Totals.Total.String(), i2.Totals.Total.String())
	assert.Equal(t, i.Totals.Tax.String(), i2.Totals.Tax.String())
	assert.Equal(t, i.Totals.Payable.String(), i2.Totals.Payable.String())
}

func TestRemoveIncludedTax2(t *testing.T) {
	i := &bill.Invoice{
		Code: "123TEST",
		Tax: &bill.Tax{
			PricesInclude: tax.CategoryVAT,
		},
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "B98602642",
			},
		},
		Customer: &org.Party{
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "54387763P",
			},
		},
		IssueDate: cal.MakeDate(2022, 6, 13),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Item",
					Price: num.MakeAmount(4320, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Percent:  num.NewPercentage(6, 2),
					},
				},
			},
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Item 2",
					Price: num.MakeAmount(259, 2),
				},
			},
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Item 3",
					Price: num.MakeAmount(300, 2),
				},
			},
		},
	}

	require.NoError(t, i.Calculate())

	i2, err := i.RemoveIncludedTaxes()
	require.NoError(t, err)
	l0 := i2.Lines[0]
	assert.Equal(t, "40.7547", l0.Item.Price.String())
	assert.Equal(t, "40.7547", l0.Total.String())

	assert.Equal(t, "46.34", i2.Totals.Sum.String())
	assert.Equal(t, i.Totals.Total.String(), i2.Totals.Total.String())
	assert.Equal(t, i.Totals.Tax.String(), i2.Totals.Tax.String())
	assert.Equal(t, i.Totals.Payable.String(), i2.Totals.Payable.String())
}

func TestRemoveIncludedTax3(t *testing.T) {
	i := &bill.Invoice{
		Code: "123TEST",
		Tax: &bill.Tax{
			PricesInclude: tax.CategoryVAT,
		},
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "B98602642",
			},
		},
		Customer: &org.Party{
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "54387763P",
			},
		},
		IssueDate: cal.MakeDate(2022, 6, 13),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Item",
					Price: num.MakeAmount(23666, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Percent:  num.NewPercentage(6, 2),
					},
				},
			},
			{
				Quantity: num.MakeAmount(2, 0),
				Item: &org.Item{
					Name:  "Item 2",
					Price: num.MakeAmount(23667, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Percent:  num.NewPercentage(6, 2),
					},
				},
			},
			{
				Quantity: num.MakeAmount(12, 0),
				Item: &org.Item{
					Name:  "Item 3",
					Price: num.MakeAmount(1000, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Percent:  num.NewPercentage(13, 2),
					},
				},
			},
			{
				Quantity: num.MakeAmount(18, 0),
				Item: &org.Item{
					Name:  "Local tax",
					Price: num.MakeAmount(150, 2),
				},
			},
		},
	}

	require.NoError(t, i.Calculate())

	i2, err := i.RemoveIncludedTaxes()
	require.NoError(t, err)
	assert.Equal(t, "223.2642", i2.Lines[0].Total.String())
	assert.Equal(t, "106.1952", i2.Lines[2].Total.String())

	/*
		data, _ := json.Marshal(i.Lines)
		t.Logf("LINES: %v", string(data))
		data, _ = json.Marshal(i.Totals)
		t.Logf("TOTALS: %v", string(data))
		data, _ = json.Marshal(i2.Lines)
		t.Logf("Lines: %v", string(data))
		data, _ = json.Marshal(i2.Totals)
		t.Logf("TOTALS: %v", string(data))
	*/

	assert.Equal(t, "803.01", i2.Totals.Sum.String())
	assert.Equal(t, i.Totals.Total.String(), i2.Totals.Total.String())
	assert.Equal(t, i.Totals.Tax.String(), i2.Totals.Tax.String())
	assert.Equal(t, i.Totals.Payable.String(), i2.Totals.Payable.String())
}

func TestRemoveIncludedTax4(t *testing.T) {
	i := &bill.Invoice{
		Code: "123TEST",
		Tax: &bill.Tax{
			PricesInclude: tax.CategoryVAT,
		},
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "B98602642",
			},
		},
		Customer: &org.Party{
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "54387763P",
			},
		},
		IssueDate: cal.MakeDate(2022, 6, 13),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(20, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.MakeAmount(400, 2),
				},
			},
			{Quantity: num.MakeAmount(1, 0), Item: &org.Item{Name: "X", Price: num.MakeAmount(40365, 2)}, Taxes: tax.Set{{Category: "VAT", Percent: num.NewPercentage(6, 2)}}},
			{Quantity: num.MakeAmount(1, 0), Item: &org.Item{Name: "X", Price: num.MakeAmount(40365, 2)}, Taxes: tax.Set{{Category: "VAT", Percent: num.NewPercentage(6, 2)}}},
			{Quantity: num.MakeAmount(1, 0), Item: &org.Item{Name: "X", Price: num.MakeAmount(40365, 2)}, Taxes: tax.Set{{Category: "VAT", Percent: num.NewPercentage(6, 2)}}},
			{Quantity: num.MakeAmount(1, 0), Item: &org.Item{Name: "X", Price: num.MakeAmount(40365, 2)}, Taxes: tax.Set{{Category: "VAT", Percent: num.NewPercentage(6, 2)}}},
			{Quantity: num.MakeAmount(1, 0), Item: &org.Item{Name: "X", Price: num.MakeAmount(40365, 2)}, Taxes: tax.Set{{Category: "VAT", Percent: num.NewPercentage(6, 2)}}},
			{Quantity: num.MakeAmount(1, 0), Item: &org.Item{Name: "X", Price: num.MakeAmount(40365, 2)}, Taxes: tax.Set{{Category: "VAT", Percent: num.NewPercentage(6, 2)}}},
			{Quantity: num.MakeAmount(1, 0), Item: &org.Item{Name: "X", Price: num.MakeAmount(40365, 2)}, Taxes: tax.Set{{Category: "VAT", Percent: num.NewPercentage(6, 2)}}},
			{Quantity: num.MakeAmount(1, 0), Item: &org.Item{Name: "X", Price: num.MakeAmount(40365, 2)}, Taxes: tax.Set{{Category: "VAT", Percent: num.NewPercentage(6, 2)}}},
			{Quantity: num.MakeAmount(1, 0), Item: &org.Item{Name: "X", Price: num.MakeAmount(40365, 2)}, Taxes: tax.Set{{Category: "VAT", Percent: num.NewPercentage(6, 2)}}},
			{Quantity: num.MakeAmount(1, 0), Item: &org.Item{Name: "X", Price: num.MakeAmount(40365, 2)}, Taxes: tax.Set{{Category: "VAT", Percent: num.NewPercentage(6, 2)}}},
			{Quantity: num.MakeAmount(1, 0), Item: &org.Item{Name: "X", Price: num.MakeAmount(40365, 2)}, Taxes: tax.Set{{Category: "VAT", Percent: num.NewPercentage(6, 2)}}},
		},
	}

	require.NoError(t, i.Calculate())

	i2, err := i.RemoveIncludedTaxes()
	require.NoError(t, err)

	data, _ := json.Marshal(i2.Lines)
	t.Logf("TOTALS: %v", string(data))
	assert.Equal(t, "4268.82", i2.Totals.Sum.String())
	assert.Equal(t, i.Totals.Total.String(), i2.Totals.Total.String())
	assert.Equal(t, i.Totals.Tax.String(), i2.Totals.Tax.String())
	assert.Equal(t, i.Totals.Payable.String(), i2.Totals.Payable.String())
}

func TestRemoveIncludedTax5(t *testing.T) {
	lines := []*bill.Line{
		{
			Quantity: num.MakeAmount(32, 0),
			Item: &org.Item{
				Name:  "Test Item",
				Price: num.MakeAmount(4375, 2),
			},
			Taxes: tax.Set{
				{
					Category: "VAT",
					Percent:  num.NewPercentage(6, 2),
				},
			},
		},
	}
	i := baseInvoice(t, lines...)
	require.NoError(t, i.Calculate())
	i2, err := i.RemoveIncludedTaxes()
	require.NoError(t, err)

	assert.Empty(t, i2.Tax.PricesInclude)
	l0 := i2.Lines[0]
	assert.Equal(t, "41.2736", l0.Item.Price.String())

	/*
		data, _ := json.Marshal(i.Lines)
		t.Logf("LINES: %v", string(data))
		data, _ = json.Marshal(i.Totals)
		t.Logf("TOTALS: %v", string(data))
		data, _ = json.Marshal(i2.Lines)
		t.Logf("Lines: %v", string(data))
		data, _ = json.Marshal(i2.Totals)
		t.Logf("TOTALS: %v", string(data))
	*/

	assert.Equal(t, "1320.76", i2.Totals.Sum.String())
	// in this case the total is different, but that's acceptable as long
	// as the payable total is the same
	//assert.Equal(t, i.Totals.Total.String(), i2.Totals.Total.String())
	assert.Equal(t, i.Totals.Tax.String(), i2.Totals.Tax.String())
	assert.Equal(t, i.Totals.Payable.String(), i2.Totals.Payable.String())
}

func TestRemoveIncludedTax6WithDiscount(t *testing.T) {
	lines := []*bill.Line{
		{
			Quantity: num.MakeAmount(1, 0),
			Item: &org.Item{
				Name:  "Test Item",
				Price: num.MakeAmount(9338, 2),
			},
			Discounts: []*bill.LineDiscount{
				{
					Percent: num.NewPercentage(40009, 5),
				},
			},
			Taxes: tax.Set{
				{
					Category: "VAT",
					Percent:  num.NewPercentage(23, 2),
				},
			},
		},
	}
	i := baseInvoice(t, lines...)
	require.NoError(t, i.Calculate())

	assert.Equal(t, "56.02", i.Totals.Sum.String())
	assert.Equal(t, i.Totals.Sum.String(), i.Totals.Payable.String())

	/*
		data, _ := json.Marshal(i.Lines)
		t.Logf("LINES: %v", string(data))
		data, _ = json.Marshal(i.Totals)
		t.Logf("TOTALS: %v", string(data))
	*/

	i2, err := i.RemoveIncludedTaxes()
	require.NoError(t, err)

	assert.Empty(t, i2.Tax.PricesInclude)
	l0 := i2.Lines[0]
	assert.Equal(t, "75.9187", l0.Item.Price.String())

	/*
		data, _ = json.Marshal(i2.Lines)
		t.Logf("Lines: %v", string(data))
		data, _ = json.Marshal(i2.Totals)
		t.Logf("TOTALS: %v", string(data))
	*/

	assert.Equal(t, "45.54", i2.Totals.Sum.String())
	assert.Equal(t, i.Totals.Total.String(), i2.Totals.Total.String())
	assert.Equal(t, i.Totals.Tax.String(), i2.Totals.Tax.String())
	assert.Equal(t, i.Totals.Payable.String(), i2.Totals.Payable.String())
}

func TestRemoveIncludedTaxQuantity(t *testing.T) {
	i := &bill.Invoice{
		Code: "123TEST",
		Tax: &bill.Tax{
			PricesInclude: tax.CategoryVAT,
		},
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "B98602642",
			},
		},
		Customer: &org.Party{
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "54387763P",
			},
		},
		IssueDate: cal.MakeDate(2022, 6, 13),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(100, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.MakeAmount(1000, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Percent:  num.NewPercentage(21, 2),
					},
				},
				Discounts: []*bill.LineDiscount{
					{
						Reason:  "Testing",
						Percent: num.NewPercentage(10, 2),
					},
				},
			},
		},
	}

	require.NoError(t, i.Calculate())

	i2, err := i.RemoveIncludedTaxes()
	require.NoError(t, err)
	require.NotNil(t, i2)

	assert.Empty(t, i2.Tax.PricesInclude)
	l0 := i2.Lines[0]
	assert.Equal(t, "8.2645", l0.Item.Price.String())
	assert.Equal(t, "826.4500", l0.Sum.String())
	assert.Equal(t, "82.6450", l0.Discounts[0].Amount.String())
	assert.Equal(t, "743.8050", l0.Total.String())
	assert.Equal(t, "10.00", i.Lines[0].Item.Price.String())

	/*
		data, _ := json.Marshal(i.Lines)
		t.Logf("LINES: %v", string(data))
		data, _ = json.Marshal(i.Totals)
		t.Logf("TOTALS: %v", string(data))
		data, _ = json.Marshal(i2.Lines)
		t.Logf("Lines: %v", string(data))
		data, _ = json.Marshal(i2.Totals)
		t.Logf("TOTALS: %v", string(data))
	*/

	// Total changes slightly
	//assert.Equal(t, i.Totals.Total.String(), i2.Totals.Total.String())
	assert.Equal(t, "743.81", i2.Totals.Total.String())
	assert.Equal(t, i.Totals.Tax.String(), i2.Totals.Tax.String())
	assert.Equal(t, i.Totals.Payable.String(), i2.Totals.Payable.String())
}

func TestRemoveIncludedTaxDeep(t *testing.T) {
	i := &bill.Invoice{
		Code: "123TEST",
		Tax: &bill.Tax{
			PricesInclude: tax.CategoryVAT,
		},
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "B98602642",
			},
		},
		Customer: &org.Party{
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "54387763P",
			},
		},
		IssueDate: cal.MakeDate(2022, 6, 13),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(364, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.MakeAmount(5178, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Percent:  num.NewPercentage(6, 2),
					},
				},
			},
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Test Item 2",
					Price: num.MakeAmount(5208, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Percent:  num.NewPercentage(6, 2),
					},
				},
			},
		},
	}

	require.NoError(t, i.Calculate())

	i2, err := i.RemoveIncludedTaxes()
	require.NoError(t, err)

	assert.Empty(t, i2.Tax.PricesInclude)
	l0 := i2.Lines[0]
	assert.Equal(t, "48.8491", l0.Item.Price.String())
	assert.Equal(t, "17781.0724", l0.Sum.String())
	l1 := i2.Lines[1]
	assert.Equal(t, "49.1321", l1.Item.Price.String())
	assert.Equal(t, "49.1321", l1.Sum.String())

	/*
		data, _ := json.Marshal(i.Lines)
		t.Logf("LINES: %v", string(data))
		data, _ = json.Marshal(i.Totals)
		t.Logf("TOTALS: %v", string(data))
		data, _ = json.Marshal(i2.Lines)
		t.Logf("Lines: %v", string(data))
		data, _ = json.Marshal(i2.Totals)
		t.Logf("TOTALS: %v", string(data))
	*/

	assert.Equal(t, "17830.20", i2.Totals.Total.String())
	// assert.Equal(t, i.Totals.Total.String(), i2.Totals.Total.String())
	assert.Equal(t, "17830.20", i2.Totals.Total.String())
	assert.Equal(t, "-0.02", i2.Totals.Rounding.String())
	assert.Equal(t, i.Totals.Tax.String(), i2.Totals.Tax.String())
	assert.Equal(t, i.Totals.Payable.String(), i2.Totals.Payable.String())
}

func TestRemoveIncludedTaxDeep2(t *testing.T) {
	i := &bill.Invoice{
		Code: "123TEST",
		Tax: &bill.Tax{
			PricesInclude: tax.CategoryVAT,
		},
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "B98602642",
			},
		},
		Customer: &org.Party{
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "54387763P",
			},
		},
		IssueDate: cal.MakeDate(2022, 6, 13),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(99999, 3),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.MakeAmount(5178, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Percent:  num.NewPercentage(6, 2),
					},
				},
			},
		},
	}

	require.NoError(t, i.Calculate())

	i2, err := i.RemoveIncludedTaxes()
	require.NoError(t, err)

	//data, _ := json.MarshalIndent(i2, "", "  ")
	//t.Log(string(data))

	assert.Empty(t, i2.Tax.PricesInclude)
	l0 := i2.Lines[0]
	assert.Equal(t, "48.8491", l0.Item.Price.String())
	assert.Equal(t, "4884.8612", l0.Sum.String())

	assert.Equal(t, "4884.86", i2.Totals.Total.String())
	assert.Equal(t, i.Totals.Total.String(), i2.Totals.Total.String())
	assert.Equal(t, i.Totals.Tax.String(), i2.Totals.Tax.String())
	assert.Equal(t, i.Totals.Payable.String(), i2.Totals.Payable.String())
}

func TestCalculateTotalsWithFractions(t *testing.T) {
	i := &bill.Invoice{
		Code: "123TEST",
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "B98602642",
			},
		},
		Customer: &org.Party{
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "54387763P",
			},
		},
		IssueDate: cal.MakeDate(2022, 6, 13),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(2010, 2),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.MakeAmount(305, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Percent:  num.NewPercentage(230, 3),
					},
				},
			},
			{
				Quantity: num.MakeAmount(2010, 2),
				Item: &org.Item{
					Name:  "Test Item 2",
					Price: num.MakeAmount(305, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Percent:  num.NewPercentage(230, 3),
					},
				},
			},
		},
	}

	require.NoError(t, i.Calculate())

	//data, _ := json.MarshalIndent(i, "", "  ")
	//t.Log(string(data))

	l0 := i.Lines[0]
	assert.Equal(t, "3.05", l0.Item.Price.String())
	assert.Equal(t, "61.31", l0.Sum.String())
	assert.Equal(t, "122.61", i.Totals.Total.String())
}

func TestCalculate(t *testing.T) {
	i := &bill.Invoice{
		Code: "123TEST",
		Tax: &bill.Tax{
			PricesInclude: tax.CategoryVAT,
		},
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "B98602642",
			},
		},
		Customer: &org.Party{
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "54387763P",
			},
		},
		IssueDate: cal.MakeDate(2022, 6, 13),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(10, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.MakeAmount(10000, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "standard",
					},
				},
				Discounts: []*bill.LineDiscount{
					{
						Reason:  "Testing",
						Percent: num.NewPercentage(10, 2),
					},
				},
				Charges: []*bill.LineCharge{
					{
						Reason:  "Testing Charge",
						Percent: num.NewPercentage(5, 2),
					},
				},
			},
		},
		Outlays: []*bill.Outlay{
			{
				Description: "Something paid in advance",
				Amount:      num.MakeAmount(1000, 2),
			},
		},
		Payment: &bill.Payment{
			Advances: []*pay.Advance{
				{
					Description: "Test Advance",
					Percent:     num.NewPercentage(30, 2), // 30%
				},
			},
		},
	}

	require.NoError(t, i.Calculate())
	assert.Equal(t, i.Totals.Sum.String(), "950.00")
	assert.Equal(t, i.Totals.Total.String(), "785.12")
	assert.Equal(t, i.Totals.Tax.String(), "164.88")
	assert.Equal(t, i.Totals.TotalWithTax.String(), "950.00")
	assert.Equal(t, i.Payment.Advances[0].Amount.String(), "285.00")
	assert.Equal(t, i.Totals.Advances.String(), "285.00")
	assert.Equal(t, i.Totals.Payable.String(), "960.00")
	assert.Equal(t, i.Totals.Due.String(), "675.00")
}

func TestCalculateInverted(t *testing.T) {
	i := &bill.Invoice{
		Code: "123TEST",
		Tax: &bill.Tax{
			PricesInclude: tax.CategoryVAT,
		},
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "B98602642",
			},
		},
		Customer: &org.Party{
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "54387763P",
			},
		},
		IssueDate: cal.MakeDate(2022, 6, 13),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(10, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.MakeAmount(10000, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "standard",
					},
				},
				Discounts: []*bill.LineDiscount{
					{
						Reason: "Testing",
						Amount: num.MakeAmount(10000, 2),
					},
				},
				Charges: []*bill.LineCharge{
					{
						Reason: "Testing Charge",
						Amount: num.MakeAmount(5000, 2),
					},
				},
			},
		},
		Outlays: []*bill.Outlay{
			{
				Description: "Something paid in advance",
				Amount:      num.MakeAmount(1000, 2),
			},
		},
		Payment: &bill.Payment{
			Advances: []*pay.Advance{
				{
					Description: "Test Advance",
					Amount:      num.MakeAmount(25000, 2),
				},
			},
		},
	}

	require.NoError(t, i.Calculate())
	assert.Equal(t, i.Totals.Sum.String(), "950.00")
	assert.Equal(t, i.Totals.Due.String(), "710.00")

	require.NoError(t, i.Invert())
	assert.Equal(t, i.Totals.Sum.String(), "-950.00")
	assert.Equal(t, i.Totals.Due.String(), "-710.00")
}

func TestInvoiceForUnknownRegime(t *testing.T) {
	lines := []*bill.Line{
		{
			Quantity: num.MakeAmount(32, 0),
			Item: &org.Item{
				Name:  "Test Item",
				Price: num.MakeAmount(4375, 2),
			},
			Taxes: tax.Set{
				{
					Category: "VAT",
					Percent:  num.NewPercentage(6, 2),
				},
			},
		},
	}
	inv := baseInvoice(t, lines...)

	// Set an undefined regime
	inv.Supplier.TaxID.Country = l10n.AD.Tax()
	assert.Nil(t, tax.RegimeFor(l10n.AD), "if Andorra is defined, change this to another country")

	assert.ErrorContains(t, inv.Calculate(), "currency: missing")
	inv.Currency = currency.USD
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())
}

func TestValidation(t *testing.T) {
	t.Run("basic validation", func(t *testing.T) {
		inv := baseInvoiceWithLines(t)
		inv.Code = ""
		ctx := context.Background()
		require.NoError(t, inv.Calculate())
		err := inv.ValidateWithContext(ctx)
		assert.ErrorContains(t, err, "code: cannot be blank")
		ctx = context.WithValue(ctx, internal.KeyDraft, true)
		assert.NoError(t, inv.ValidateWithContext(ctx))
	})

	t.Run("supplier name", func(t *testing.T) {
		inv := baseInvoiceWithLines(t)
		inv.Supplier.Name = ""
		inv.Customer.TaxID.Code = "" // simplified
		inv.Customer.Name = ""       // so this is okay
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "supplier: (name: cannot be blank.).")
		assert.NotContains(t, err.Error(), "customer")
	})

	t.Run("simplified", func(t *testing.T) {
		inv := baseInvoiceWithLines(t)
		inv.Customer = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "customer: cannot be blank.")

		inv.Tax = &bill.Tax{
			Tags: []cbc.Key{
				tax.TagSimplified,
			},
		}
		require.NoError(t, inv.Calculate())
		err = inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("simplified without customer name", func(t *testing.T) {
		inv := baseInvoiceWithLines(t)
		inv.Tax = &bill.Tax{
			Tags: []cbc.Key{
				tax.TagSimplified,
			},
		}
		inv.Customer.TaxID = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("implied simplified without customer tax ID", func(t *testing.T) {
		inv := baseInvoiceWithLines(t)
		inv.Customer.TaxID = nil
		inv.Customer.Name = ""
		inv.Customer.Emails = append(inv.Customer.Emails, &org.Email{
			Address: "foo@example.com",
		})
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("not simplified without customer name", func(t *testing.T) {
		inv := baseInvoiceWithLines(t)
		inv.Customer.Name = ""
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "customer: (name: cannot be blank.).")
	})
}

func baseInvoiceWithLines(t *testing.T) *bill.Invoice {
	inv := baseInvoice(t,
		&bill.Line{
			Quantity: num.MakeAmount(10, 0),
			Item: &org.Item{
				Name:  "Test Item",
				Price: num.MakeAmount(10000, 2),
			},
			Taxes: tax.Set{
				{
					Category: "VAT",
					Rate:     "standard",
				},
			},
		},
	)
	return inv
}

func baseInvoice(t *testing.T, lines ...*bill.Line) *bill.Invoice {
	t.Helper()
	i := &bill.Invoice{
		Series:    "TEST",
		Code:      "00123",
		IssueDate: cal.MakeDate(2022, 6, 13),
		Tax: &bill.Tax{
			PricesInclude: tax.CategoryVAT,
		},
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "B98602642",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "54387763P",
			},
		},
		Lines: lines,
	}
	return i
}
