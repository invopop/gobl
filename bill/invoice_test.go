package bill_test

import (
	"testing"

	_ "github.com/invopop/gobl" // load regions
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/regions/common"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRemoveIncludedTax(t *testing.T) {
	i := &bill.Invoice{
		Code: "123TEST",
		Tax: &bill.Tax{
			PricesInclude: common.TaxCategoryVAT,
		},
		Supplier: &org.Party{
			TaxID: &org.TaxIdentity{
				Country: l10n.ES,
				Code:    "B98602642",
			},
		},
		Customer: &org.Party{
			TaxID: &org.TaxIdentity{
				Country: l10n.ES,
				Code:    "54387763P",
			},
		},
		IssueDate: cal.MakeDate(2022, 6, 13),
		Lines: bill.Lines{
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
			},
		},
	}

	require.NoError(t, i.Calculate())

	i2 := i.RemoveIncludedTaxes(2)

	assert.Empty(t, i2.Tax.PricesInclude)
	l0 := i2.Lines[0]
	assert.Equal(t, "82.6446", l0.Item.Price.String())
	assert.Equal(t, "826.4463", l0.Sum.String())
	assert.Equal(t, "826.4463", l0.Sum.String())
	assert.Equal(t, "82.6446", l0.Discounts[0].Amount.String())
	assert.Equal(t, "743.8017", l0.Total.String())

	assert.Equal(t, "100.00", i.Lines[0].Item.Price.String())
}

func TestCalculate(t *testing.T) {
	i := &bill.Invoice{
		Code: "123TEST",
		Tax: &bill.Tax{
			PricesInclude: common.TaxCategoryVAT,
		},
		Supplier: &org.Party{
			TaxID: &org.TaxIdentity{
				Country: l10n.ES,
				Code:    "B98602642",
			},
		},
		Customer: &org.Party{
			TaxID: &org.TaxIdentity{
				Country: l10n.ES,
				Code:    "54387763P",
			},
		},
		IssueDate: cal.MakeDate(2022, 6, 13),
		Lines: bill.Lines{
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
