package bill_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	_ "github.com/invopop/gobl" // load regions
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/internal"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
					Rate:     common.TaxRateStandard,
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

	i2 := i.RemoveIncludedTaxes()
	require.NoError(t, i2.Calculate())

	assert.Equal(t, "1000.00", i.Lines[0].Item.Price.String())

	assert.Empty(t, i2.Tax.PricesInclude)
	l0 := i2.Lines[0]
	assert.Equal(t, "826.4463", l0.Item.Price.String())
	assert.Equal(t, "826.4463", l0.Sum.String())
	assert.Equal(t, "826.4463", l0.Sum.String())
	assert.Equal(t, "82.6446", l0.Discounts[0].Amount.String())
	assert.Equal(t, "743.8017", l0.Total.String())

	assert.Equal(t, "743.8017", i2.Totals.Sum.String())
	assert.Equal(t, i.Totals.Total.String(), i2.Totals.Total.String())
	assert.Equal(t, i.Totals.Tax.String(), i2.Totals.Tax.String())
	assert.Equal(t, i.Totals.Payable.String(), i2.Totals.Payable.String())
}

func TestRemoveIncludedTax2(t *testing.T) {
	i := &bill.Invoice{
		Code: "123TEST",
		Tax: &bill.Tax{
			PricesInclude: common.TaxCategoryVAT,
		},
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Country: l10n.ES,
				Code:    "B98602642",
			},
		},
		Customer: &org.Party{
			TaxID: &tax.Identity{
				Country: l10n.ES,
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

	i2 := i.RemoveIncludedTaxes()
	require.NoError(t, i2.Calculate())
	l0 := i2.Lines[0]
	assert.Equal(t, "40.7547", l0.Item.Price.String())
	assert.Equal(t, "40.7547", l0.Total.String())

	assert.Equal(t, "46.3447", i2.Totals.Sum.String())
	assert.Equal(t, i.Totals.Total.String(), i2.Totals.Total.String())
	assert.Equal(t, i.Totals.Tax.String(), i2.Totals.Tax.String())
	assert.Equal(t, i.Totals.Payable.String(), i2.Totals.Payable.String())
}

func TestRemoveIncludedTax3(t *testing.T) {
	i := &bill.Invoice{
		Code: "123TEST",
		Tax: &bill.Tax{
			PricesInclude: common.TaxCategoryVAT,
		},
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Country: l10n.ES,
				Code:    "B98602642",
			},
		},
		Customer: &org.Party{
			TaxID: &tax.Identity{
				Country: l10n.ES,
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

	i2 := i.RemoveIncludedTaxes()
	require.NoError(t, i2.Calculate())
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
	assert.Equal(t, "803.0066", i2.Totals.Sum.String())
	assert.Equal(t, i.Totals.Total.String(), i2.Totals.Total.String())
	assert.Equal(t, i.Totals.Tax.String(), i2.Totals.Tax.String())
	assert.Equal(t, i.Totals.Payable.String(), i2.Totals.Payable.String())
}

func TestRemoveIncludedTax4(t *testing.T) {
	i := &bill.Invoice{
		Code: "123TEST",
		Tax: &bill.Tax{
			PricesInclude: common.TaxCategoryVAT,
		},
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Country: l10n.ES,
				Code:    "B98602642",
			},
		},
		Customer: &org.Party{
			TaxID: &tax.Identity{
				Country: l10n.ES,
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

	i2 := i.RemoveIncludedTaxes()
	require.NoError(t, i2.Calculate())

	data, _ := json.Marshal(i2.Lines)
	t.Logf("TOTALS: %v", string(data))
	assert.Equal(t, "4268.8209", i2.Totals.Sum.String())
	assert.Equal(t, i.Totals.Total.String(), i2.Totals.Total.String())
	assert.Equal(t, i.Totals.Tax.String(), i2.Totals.Tax.String())
	assert.Equal(t, i.Totals.Payable.String(), i2.Totals.Payable.String())
}

func TestRemoveIncludedTaxQuantity(t *testing.T) {
	i := &bill.Invoice{
		Code: "123TEST",
		Tax: &bill.Tax{
			PricesInclude: common.TaxCategoryVAT,
		},
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Country: l10n.ES,
				Code:    "B98602642",
			},
		},
		Customer: &org.Party{
			TaxID: &tax.Identity{
				Country: l10n.ES,
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

	i2 := i.RemoveIncludedTaxes()
	require.NoError(t, i2.Calculate())

	assert.Empty(t, i2.Tax.PricesInclude)
	l0 := i2.Lines[0]
	assert.Equal(t, "8.26446", l0.Item.Price.String())
	assert.Equal(t, "826.44600", l0.Sum.String())
	assert.Equal(t, "82.64460", l0.Discounts[0].Amount.String())
	assert.Equal(t, "743.80140", l0.Total.String())
	assert.Equal(t, "10.00", i.Lines[0].Item.Price.String())

	assert.Equal(t, i.Totals.Total.String(), i2.Totals.Total.String())
	assert.Equal(t, i.Totals.Tax.String(), i2.Totals.Tax.String())
	assert.Equal(t, i.Totals.Payable.String(), i2.Totals.Payable.String())
}

func TestRemoveIncludedTaxDeep(t *testing.T) {
	i := &bill.Invoice{
		Code: "123TEST",
		Tax: &bill.Tax{
			PricesInclude: common.TaxCategoryVAT,
		},
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Country: l10n.ES,
				Code:    "B98602642",
			},
		},
		Customer: &org.Party{
			TaxID: &tax.Identity{
				Country: l10n.ES,
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

	i2 := i.RemoveIncludedTaxes()

	require.NoError(t, i2.Calculate())

	//data, _ := json.MarshalIndent(i2, "", "  ")
	//t.Log(string(data))

	assert.Empty(t, i2.Tax.PricesInclude)
	l0 := i2.Lines[0]
	assert.Equal(t, "48.84906", l0.Item.Price.String()) // note extra digit!
	assert.Equal(t, "17781.05784", l0.Sum.String())
	l1 := i2.Lines[1]
	assert.Equal(t, "49.1321", l1.Item.Price.String())
	assert.Equal(t, "49.1321", l1.Sum.String())

	assert.Equal(t, "17830.19", i2.Totals.Total.String())
	assert.Equal(t, i.Totals.Total.String(), i2.Totals.Total.String())
	assert.Equal(t, i.Totals.Tax.String(), i2.Totals.Tax.String())
	assert.Equal(t, i.Totals.Payable.String(), i2.Totals.Payable.String())
}

func TestRemoveIncludedTaxDeep2(t *testing.T) {
	i := &bill.Invoice{
		Code: "123TEST",
		Tax: &bill.Tax{
			PricesInclude: common.TaxCategoryVAT,
		},
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Country: l10n.ES,
				Code:    "B98602642",
			},
		},
		Customer: &org.Party{
			TaxID: &tax.Identity{
				Country: l10n.ES,
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

	i2 := i.RemoveIncludedTaxes()

	require.NoError(t, i2.Calculate())

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

func TestCalculate(t *testing.T) {
	i := &bill.Invoice{
		Code: "123TEST",
		Tax: &bill.Tax{
			PricesInclude: common.TaxCategoryVAT,
		},
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Country: l10n.ES,
				Code:    "B98602642",
			},
		},
		Customer: &org.Party{
			TaxID: &tax.Identity{
				Country: l10n.ES,
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

func TestValidation(t *testing.T) {
	inv := &bill.Invoice{
		Currency:  currency.EUR,
		IssueDate: cal.MakeDate(2022, 6, 13),
		Tax: &bill.Tax{
			PricesInclude: common.TaxCategoryVAT,
		},
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: l10n.ES,
				Code:    "B98602642",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: l10n.ES,
				Code:    "54387763P",
			},
		},
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
			},
		},
	}
	ctx := context.Background()
	require.NoError(t, inv.Calculate())
	err := inv.ValidateWithContext(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "code: cannot be blank")
	ctx = context.WithValue(ctx, internal.KeyDraft, true)
	assert.NoError(t, inv.ValidateWithContext(ctx))
}

func baseInvoice(t *testing.T, lines ...*bill.Line) *bill.Invoice {
	t.Helper()
	i := &bill.Invoice{
		Code: "123TEST",
		Tax: &bill.Tax{
			PricesInclude: common.TaxCategoryVAT,
		},
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Country: l10n.ES,
				Code:    "B98602642",
			},
		},
		Customer: &org.Party{
			TaxID: &tax.Identity{
				Country: l10n.ES,
				Code:    "54387763P",
			},
		},
		IssueDate: cal.MakeDate(2022, 6, 13),
		Lines:     lines,
	}
	return i
}
